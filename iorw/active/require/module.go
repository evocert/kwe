package require

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"
	"syscall"
	"text/template"

	js "github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/evocert/kwe/iorw/parsing"
)

type ModuleLoader func(*js.Runtime, *js.Object)

// SourceLoader represents a function that returns a file data at a given path.
// The function should return ModuleFileDoesNotExistError if the file either doesn't exist or is a directory.
// This error will be ignored by the resolver and the search will continue. Any other errors will be propagated.
type SourceLoader func(path string) ([]byte, error)

var (
	InvalidModuleError     = errors.New("Invalid module")
	IllegalModuleNameError = errors.New("Illegal module name")

	ModuleFileDoesNotExistError = errors.New("module file does not exist")
)

var native map[string]ModuleLoader

// Registry contains a cache of compiled modules which can be used by multiple Runtimes
type Registry struct {
	sync.Mutex
	native        map[string]ModuleLoader
	compiled      map[string]*js.Program
	parsed        map[string]*parsing.Parsing
	srcLoader     SourceLoader
	globalFolders []string
	Actv          parsing.AltActiveAPI
}

type RequireModule struct {
	r           *Registry
	runtime     *js.Runtime
	modules     map[string]*js.Object
	nodeModules map[string]*js.Object
	parsings    map[*parsing.Parsing]*parsing.Parsing
	Lstprsng    *parsing.Parsing
}

func NewRegistry(opts ...Option) *Registry {
	r := &Registry{}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func NewRegistryWithLoader(srcLoader SourceLoader) *Registry {
	return NewRegistry(WithLoader(srcLoader))
}

type Option func(*Registry)

// WithLoader sets a function which will be called by the require() function in order to get a source code for a
// module at the given path. The same function will be used to get external source maps.
// Note, this only affects the modules loaded by the require() function. If you need to use it as a source map
// loader for code parsed in a different way (such as runtime.RunString() or eval()), use (*Runtime).SetParserOptions()
func WithLoader(srcLoader SourceLoader) Option {
	return func(r *Registry) {
		r.srcLoader = srcLoader
	}
}

// WithGlobalFolders appends the given paths to the registry's list of
// global folders to search if the requested module is not found
// elsewhere.  By default, a registry's global folders list is empty.
// In the reference Node.js implementation, the default global folders
// list is $NODE_PATH, $HOME/.node_modules, $HOME/.node_libraries and
// $PREFIX/lib/node, see
// https://nodejs.org/api/modules.html#modules_loading_from_the_global_folders.
func WithGlobalFolders(globalFolders ...string) Option {
	return func(r *Registry) {
		r.globalFolders = globalFolders
	}
}

func (r *Registry) Dispose() {
	if r != nil {
		if r.compiled != nil || r.native != nil || r.parsed != nil {
			func() {
				r.Lock()
				defer r.Unlock()
				if r.compiled != nil {

					if len(r.compiled) > 0 {
						for k := range r.compiled {
							r.compiled[k] = nil
							delete(r.compiled, k)
						}
					}
					r.compiled = nil
				}
				if r.native != nil {
					func() {
						r.Lock()
						defer r.Unlock()
						if len(r.native) > 0 {
							for k := range r.native {
								r.native[k] = nil
								delete(r.native, k)
							}
						}
					}()
					r.compiled = nil
				}
				if r.parsed != nil {
					if len(r.parsed) > 0 {
						for k := range r.parsed {
							r.parsed[k].Dispose()
							r.parsed[k] = nil
							delete(r.parsed, k)
						}
					}
					r.parsed = nil
				}
			}()
		}
		if r.globalFolders != nil {
			r.globalFolders = nil
		}
		if r.srcLoader != nil {
			r.srcLoader = nil
		}
		if r.Actv != nil {
			r.Actv = nil
		}
		r = nil
	}
}

// Enable adds the require() function to the specified runtime.
func (r *Registry) Enable(runtime *js.Runtime) *RequireModule {
	rrt := &RequireModule{
		r:           r,
		runtime:     runtime,
		modules:     make(map[string]*js.Object),
		nodeModules: make(map[string]*js.Object),
	}

	runtime.Set("require", rrt.require)
	return rrt
}

func (r *Registry) RegisterNativeModule(name string, loader ModuleLoader) {
	r.Lock()
	defer r.Unlock()

	if r.native == nil {
		r.native = make(map[string]ModuleLoader)
	}
	name = filepathClean(name)
	r.native[name] = loader
}

// DefaultSourceLoader is used if none was set (see WithLoader()). It simply loads files from the host's filesystem.
func DefaultSourceLoader(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath.FromSlash(filename))
	if err != nil {
		if os.IsNotExist(err) || errors.Is(err, syscall.EISDIR) {
			err = ModuleFileDoesNotExistError
		}
	}
	return data, err
}

func (r *Registry) getSource(p string) ([]byte, error) {
	srcLoader := r.srcLoader
	if srcLoader == nil {
		srcLoader = DefaultSourceLoader
	}
	return srcLoader(p)
}

func (r *Registry) getCompiledSource(p string) (*js.Program, *parsing.Parsing, error) {
	r.Lock()
	defer r.Unlock()
	prsng := r.parsed[p]
	prg := r.compiled[p]
	if prg == nil {
		if buf, err := r.getSource(p); len(buf) > 0 {
			if err != nil {
				return nil, nil, err
			}
			s := string(buf)

			if prsng == nil {
				prsng = parsing.NextParsing(r.Actv, nil, nil, nil, p)
				if prsrngerr := parsing.EvalParsing(prsng, nil, nil, p, true, true, s, func(prsng *parsing.Parsing) (err error) {

					return
				}); prsrngerr == nil {
					s = parsing.Code(prsng)
				} else {
					return nil, nil, prsrngerr
				}
			}

			if path.Ext(p) == ".json" {
				s = "module.exports = JSON.parse('" + template.JSEscapeString(s) + "')"
			}

			source := "(function(exports, require, module) {" + s + "\n})"
			parsed, err := js.Parse(p, source, parser.WithSourceMapLoader(
				func(path string) (bytes []byte, byteserr error) {
					if bytes, byteserr = r.srcLoader(path); byteserr == nil {
						if len(bytes) > 0 {
							if prsng == nil {
								prsng = parsing.NextParsing(r.Actv, nil, nil, nil, path)
								parsing.EvalParsing(prsng, nil, nil, path, false, true, string(bytes), func(prsng *parsing.Parsing) (err error) {

									return
								})
							}
						}
					}
					return
				}))
			if err != nil {
				return nil, nil, err
			}
			prg, err = js.CompileAST(parsed, false)
			if err == nil {
				if r.compiled == nil {
					r.compiled = make(map[string]*js.Program)
				}
				if r.parsed == nil {
					r.parsed = make(map[string]*parsing.Parsing)
				}
				r.parsed[p] = prsng
				r.compiled[p] = prg
			}
			return prg, nil, err
		}
		return nil, nil, InvalidModuleError
	}
	return prg, prsng, nil
}

func (r *RequireModule) Dispose() {
	if r != nil {
		if r.modules != nil {
			for k := range r.modules {
				r.modules[k] = nil
				delete(r.modules, k)
			}
		}
		if r.nodeModules != nil {
			for k := range r.nodeModules {
				r.nodeModules[k] = nil
				delete(r.nodeModules, k)
			}
		}
		if r.runtime != nil {
			r.runtime = nil
		}
		if r.r != nil {
			r.r = nil
		}
	}
}

func (r *RequireModule) substring(offsets int64, offsete int64) string {
	if r.Lstprsng != nil {
		return r.Lstprsng.SubString(offsets, offsete)
	}
	return ""
}

func (r *RequireModule) require(call js.FunctionCall) js.Value {
	ret, err := r.Require(call.Argument(0).String())
	if err != nil {
		if _, ok := err.(*js.Exception); !ok {
			panic(r.runtime.NewGoError(err))
		}
		panic(err)
	}
	return ret
}

func filepathClean(p string) string {
	return path.Clean(p)
}

// Require can be used to import modules from Go source (similar to JS require() function).
func (r *RequireModule) Require(p string) (ret js.Value, err error) {
	module, err := r.resolve(p)
	if err != nil {
		return
	}
	ret = module.Get("exports")
	return
}

func Require(runtime *js.Runtime, name string) js.Value {
	if r, ok := js.AssertFunction(runtime.Get("require")); ok {
		mod, err := r(js.Undefined(), runtime.ToValue(name))
		if err != nil {
			panic(err)
		}
		return mod
	}
	panic(runtime.NewTypeError("Please enable require for this runtime using new(require.Require).Enable(runtime)"))
}

func RegisterNativeModule(name string, loader ModuleLoader) {
	if native == nil {
		native = make(map[string]ModuleLoader)
	}
	name = filepathClean(name)
	native[name] = loader
}

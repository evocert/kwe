package require

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"text/template"

	"github.com/lnksnk/ecma/es6"
)

type ModuleLoader func(*es6.Runtime, *es6.Object)

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
	native   map[string]ModuleLoader
	compiled map[string]*es6.Program

	srcLoader     SourceLoader
	globalFolders []string
}

type RequireModule struct {
	r           *Registry
	runtime     *es6.Runtime
	modules     map[string]*es6.Object
	nodeModules map[string]*es6.Object
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
// https://nodees6.org/api/modules.html#modules_loading_from_the_global_folders.
func WithGlobalFolders(globalFolders ...string) Option {
	return func(r *Registry) {
		r.globalFolders = globalFolders
	}
}

// Enable adds the require() function to the specified runtime.
func (r *Registry) Enable(runtime *es6.Runtime) *RequireModule {
	rrt := &RequireModule{
		r:           r,
		runtime:     runtime,
		modules:     make(map[string]*es6.Object),
		nodeModules: make(map[string]*es6.Object),
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
	data, err := ioutil.ReadFile(filename)
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

func (r *Registry) getCompiledSource(p string) (prg *es6.Program, err error) {
	r.Lock()
	defer r.Unlock()

	prg = r.compiled[p]
	if prg == nil {
		if buf, err1 := r.getSource(p); err1 == nil {
			s := string(buf)

			if filepath.Ext(p) == ".json" {
				s = "module.exports = JSON.parse('" + template.JSEscapeString(s) + "')"
			}

			source := "(function(exports, require, module) {" + s + "\n})"
			prg, err = es6.Compile(p, source, false)
			if err == nil {
				if r.compiled == nil {
					r.compiled = make(map[string]*es6.Program)
				}
				r.compiled[p] = prg
			}
		} else {
			err = err1
		}
	}
	return
}

func (r *RequireModule) require(call es6.FunctionCall) es6.Value {
	ret, err := r.Require(call.Argument(0).String())
	if err != nil {
		if _, ok := err.(*es6.Exception); !ok {
			panic(r.runtime.NewGoError(err))
		}
		panic(err)
	}
	return ret
}

func filepathClean(p string) string {
	return filepath.Clean(p)
}

// Require can be used to import modules from Go source (similar to JS require() function).
func (r *RequireModule) Require(p string) (ret es6.Value, err error) {
	module, err := r.resolve(p)
	if err != nil {
		return
	}
	ret = module.Get("exports")
	return
}

func Require(runtime *es6.Runtime, name string) es6.Value {
	if r, ok := es6.AssertFunction(runtime.Get("require")); ok {
		mod, err := r(es6.Undefined(), runtime.ToValue(name))
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

package active

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja/parser"

	"github.com/dop251/goja"

	"github.com/evocert/kwe/ecma/jsext"
	"github.com/evocert/kwe/iorw/active/require"
	"github.com/evocert/kwe/iorw/parsing"

	"github.com/evocert/kwe/iorw"
)

//Active - struct
type Active struct {
	Print          func(a ...interface{})
	Println        func(a ...interface{})
	BinWrite       func(b ...byte) (n int, err error)
	FPrint         func(w io.Writer, a ...interface{})
	FPrintLn       func(w io.Writer, a ...interface{})
	FBinWrite      func(w io.Writer, b ...byte) (n int, err error)
	Seek           func(offset int64, whence int) (n int64, err error)
	Readln         func() (string, error)
	ReadLines      func() ([]string, error)
	ReadAll        func() (string, error)
	BinRead        func(size int) (b []byte, err error)
	FSeek          func(io.Reader, int64, int) (int64, error)
	FReadln        func(io.Reader) (string, error)
	FReadLines     func(io.Reader) ([]string, error)
	FReadAll       func(io.Reader) (string, error)
	FBinRead       func(io.Reader, int) ([]byte, error)
	LookupTemplate func(string, ...interface{}) (io.Reader, error)
	ObjectMapRef   func() map[string]interface{}
	CleanupValue   func(vali interface{}, valt reflect.Type)
	//lckprnt        *sync.Mutex
	InterruptVM func(v interface{})
	*atvruntime
}

func (atv *Active) AltLookupTemplate(path string, a ...interface{}) (r io.Reader, err error) {
	if atv != nil && atv.LookupTemplate != nil {
		r, err = atv.LookupTemplate(path, a...)
	}
	return
}

func (atv *Active) AltPrint(w io.Writer, a ...interface{}) {
	if atv != nil {
		atv.print(w, a...)
	}
	return
}

func (atv *Active) AltPrintln(w io.Writer, a ...interface{}) {
	if atv != nil {
		atv.println(w, a...)
	}
	return
}

func (atv *Active) AltBinWrite(w io.Writer, b ...byte) (n int, err error) {
	if atv != nil {
		n, err = atv.binwrite(w, b...)
	}
	return
}

func (atv *Active) AltReadln(r io.Reader) (ln string, err error) {
	if atv != nil {
		ln, err = atv.readln(r)
	}
	return
}

func (atv *Active) AltSeek(r io.Reader, offset int64, whence int) (n int64, err error) {
	if atv != nil {
		n, err = atv.seek(r, offset, whence)
	}
	return
}

func (atv *Active) AltReadlines(r io.Reader) (lines []string, err error) {
	if atv != nil {
		lines, err = atv.readlines(r)
	}
	return
}

func (atv *Active) AltReadAll(r io.Reader) (s string, err error) {
	if atv != nil {
		s, err = atv.readAll(r)
	}
	return
}

func (atv *Active) AltBinRead(r io.Reader, size int) (b []byte, err error) {
	if atv != nil {
		b, err = atv.binread(r, size)
	}
	return
}

func (atv *Active) LockPrint() {
	/*if atv != nil && atv.lckprnt != nil {
		atv.lckprnt.Lock()
	}*/
}

func (atv *Active) UnlockPrint() {
	/*if atv != nil && atv.lckprnt != nil {
		atv.lckprnt.Unlock()
	}*/
}

//InvokeFunction ivoke *Acive.actvruntime function
func (atv *Active) InvokeFunction(functocall interface{}, args ...interface{}) (result interface{}) {
	if atv != nil {
		result = atv.atvruntime.InvokeFunction(functocall, args...)
	}
	return
}

//ExtractGlobals extract globals from atv.atvruntime
func (atv *Active) ExtractGlobals(extrglbs map[string]interface{}) {
	if atv.atvruntime != nil {
		if extrglbs != nil {
			if gbl := atv.atvruntime.vm.GlobalObject(); gbl != nil {
				for _, k := range gbl.Keys() {
					glbv := gbl.Get(k)
					if t := glbv.ExportType(); t != nil {
						if expv := glbv.Export(); expv == nil {
							extrglbs[k] = glbv
						} else {
							extrglbs[k] = expv
						}
					}
				}
				gbl = nil
			}
		}
	}
}

//ImportGlobals import globals into atv.atvruntime
func (atv *Active) ImportGlobals(imprtglbs map[string]interface{}) {
	if atv.atvruntime != nil {
		if len(imprtglbs) > 0 {
			if vm := atv.atvruntime.lclvm(); vm != nil {
				if gbl := vm.GlobalObject(); gbl != nil {
					for k, kv := range imprtglbs {
						if gjv, gjvok := kv.(goja.Value); gjvok {
							if expv := gjv.Export(); expv == nil {
								gbl.Set(k, gjv)
							} else {
								gbl.Set(k, expv)
							}
						} else {
							gbl.Set(k, kv)
						}
					}
					gbl = nil
				}
			}
		}
	}
}

//NewActive - instance
func NewActive() (atv *Active) {
	atv = &Active{atvruntime: nil}
	atv.atvruntime, _ = newatvruntime(atv)
	return
}

func (atv *Active) print(w io.Writer, a ...interface{}) {
	if prntr, prntrok := w.(iorw.Printer); prntrok {
		prntr.Print(a...)
	} else {
		if atv.Print != nil {
			if len(a) > 0 {
				func() {
					atv.Print(a...)
				}()
			}
		} else {
			if atv.FPrint != nil && w != nil {
				if len(a) > 0 {
					atv.FPrint(w, a...)
				}
			} else if w != nil {
				if len(a) > 0 {
					if prntr, prntrok := w.(iorw.Printer); prntrok {
						prntr.Print(a...)
					} else {

						iorw.Fprint(w, a...)
					}
				}
			}
		}
	}
}

func (atv *Active) println(w io.Writer, a ...interface{}) {
	if prntr, prntrok := w.(iorw.Printer); prntrok {
		prntr.Println(a...)
	} else {
		if atv.Println != nil {
			if len(a) > 0 {
				func() {
					atv.Println(a...)
				}()
			}
		} else if atv.FPrintLn != nil && w != nil {
			atv.FPrintLn(w, a...)
		} else if w != nil {
			if prntr, prntrok := w.(iorw.Printer); prntrok {
				prntr.Println(a...)
			} else {
				if len(a) > 0 {
					fmt.Fprint(w, a...)
				}
				fmt.Fprintln(w)
			}
		}
	}
}

func (atv *Active) binwrite(w io.Writer, b ...byte) (n int, err error) {
	if prntr, prntrok := w.(iorw.Printer); prntrok {
		n, err = prntr.Write(b)
	} else {
		if atv.BinWrite != nil {
			if len(b) > 0 {
				n, err = atv.BinWrite(b...)
			}
		} else if atv.FBinWrite != nil && w != nil {
			n, err = atv.FBinWrite(w, b...)
		} else if w != nil {
			if prntr, prntrok := w.(iorw.Printer); prntrok {
				n, err = prntr.Write(b)
			} else {
				if len(b) > 0 {
					n, err = w.Write(b)
				}
			}
		}
	}
	return
}

func (atv *Active) binread(r io.Reader, size int) (b []byte, err error) {
	if rdr, rdrok := r.(iorw.Reader); rdrok {
		if size > 0 {
			p := make([]byte, size)
			pn, perr := rdr.Read(p)
			if pn > 0 {
				b = make([]byte, pn)
				copy(b, p[0:pn])
			}
			err = perr
		}
	} else {
		if atv.BinWrite != nil {
			b, err = atv.BinRead(size)
		} else if atv.FBinRead != nil && r != nil {
			b, err = atv.FBinRead(r, size)
		} else if r != nil {
			atv.UnlockPrint()
			if size > 0 {
				p := make([]byte, size)
				pn, perr := r.Read(p)
				if pn > 0 {
					b = make([]byte, pn)
					copy(b, p[0:pn])
				}
				err = perr
			}
		}
	}
	return
}

func (atv *Active) seek(r io.Reader, offset int64, whence int) (n int64, err error) {
	if rds, rdsok := r.(io.Seeker); rdsok {
		n, err = rds.Seek(offset, whence)
	} else {
		if atv.Seek != nil {
			n, err = atv.Seek(offset, whence)
		} else if atv.FSeek != nil && r != nil {
			func() {
				n, err = atv.FSeek(r, offset, whence)
			}()
		} else if r != nil {
			if rds, rdsok := r.(io.Seeker); rdsok {
				n, err = rds.Seek(offset, whence)
			}
		}
	}
	return
}

func (atv *Active) readln(r io.Reader) (ln string, err error) {
	if rdr, rdrok := r.(iorw.Reader); rdrok {
		ln, err = rdr.Readln()
	} else {
		if atv.Readln != nil {
			ln, err = atv.Readln()
		} else if atv.FReadln != nil && r != nil {
			func() {
				ln, err = atv.FReadln(r)
			}()
		} else if r != nil {
			if rdr, rdrok := r.(iorw.Reader); rdrok {
				ln, err = rdr.Readln()
			} else {
				ln, err = iorw.ReadLine(r)
			}
		}
	}
	return
}

func (atv *Active) readlines(r io.Reader) (lines []string, err error) {
	if rdr, rdrok := r.(iorw.Reader); rdrok {
		lines, err = rdr.Readlines()
	} else {
		if atv.ReadLines != nil {
			lines, err = atv.ReadLines()
		} else if atv.FReadLines != nil && r != nil {
			func() {
				lines, err = atv.FReadLines(r)
			}()
		} else if r != nil {
			if rdr, rdrok := r.(iorw.Reader); rdrok {
				lines, err = rdr.Readlines()
			} else {
				lines, err = iorw.ReadLines(r)
			}
		}
	}
	return
}

func (atv *Active) readAll(r io.Reader) (s string, err error) {
	if rdr, rdrok := r.(iorw.Reader); rdrok {
		s, err = rdr.ReadAll()
	} else {
		if atv.ReadAll != nil {
			s, err = atv.ReadAll()
		} else if atv.FReadAll != nil && r != nil {
			func() {
				s, err = atv.FReadAll(r)
			}()
		} else if r != nil {
			if rdr, rdrok := r.(iorw.Reader); rdrok {
				s, err = rdr.ReadAll()
			} else {
				s, err = iorw.ReaderToString(r)
			}
		}
	}
	return
}

//InvokeVM invoke vm
func (atv *Active) InvokeVM(callback func(vm *goja.Runtime) error) {
	if callback != nil {
		callback(atv.lclvm())
	}
}

func (atv *Active) vm() (vm *goja.Runtime) {
	if atv != nil && atv.atvruntime != nil && atv.atvruntime.vm != nil {
		vm = atv.atvruntime.vm
	}
	return
}

type CodeException struct {
	cde      string
	err      error
	execpath string
}

func codeException(cde string, execpath string, err error) (cdeerr *CodeException) {
	cdeerr = &CodeException{err: err, execpath: execpath, cde: cde}
	return
}

func (cdeerr *CodeException) Error() (s string) {
	s = "err:" + cdeerr.err.Error() + "\r\n"
	s += "path:" + cdeerr.execpath + "\r\n"
	s += cdeerr.cde
	return
}

func (cdeerr *CodeException) Code() string {
	return cdeerr.cde
}

func (cdeerr *CodeException) ExecPath() string {
	return cdeerr.execpath
}

func (atv *Active) atvrun(prsng *parsing.Parsing) (err error) {
	if prsng != nil {
		if atv.atvruntime == nil {
			atv.atvruntime, err = newatvruntime(atv)
		}
		if atv.atvruntime != nil {
			atv.atvruntime.prsng = prsng
			_, err = atv.atvruntime.run()
		}
	}
	return
}

//Eval - parse a ...interface{} arguments, execute if neaded and output to wou io.Writer
func (atv *Active) Eval(wout io.Writer, rin io.Reader, initpath string, invertactpsv bool, a ...interface{}) (err error) {
	var prsng = parsing.NextParsing(atv, nil, rin, wout, initpath)
	defer prsng.Dispose()
	if len(a) > 0 {
		if invertactpsv {
			a = append(append([]interface{}{"<@"}, a...), "@>")
		}
	}
	err = parsing.ParsePrsng(prsng, true, atv.processParsing, a...)
	return
}

func (atv *Active) processParsing(prsng *parsing.Parsing) (err error) {
	err = atv.atvrun(prsng)
	return
}

//Close - refer to  io.Closer
func (atv *Active) Close() (err error) {
	//putActive(atv)
	err = atv.Dispose()
	return
}

//Dispose
func (atv *Active) Dispose() (err error) {
	/*if atv.lckprnt != nil {
		atv.lckprnt = nil
	}*/
	if atv.atvruntime != nil {
		atv.atvruntime.dispose(atv.CleanupValue)
		atv.atvruntime = nil
	}
	return
}

//Clear
func (atv *Active) Clear() (err error) {
	if atv.atvruntime != nil {
		atv.atvruntime.dispose(atv.CleanupValue, true)
	}
	return
}

//Interrupt - Active processing
func (atv *Active) Interrupt() {
	if atv.InterruptVM != nil {
		atv.InterruptVM("exit")
	}
}

type atvruntime struct {
	prsng         *parsing.Parsing
	atv           *Active
	vm            *goja.Runtime
	vmregister    *require.Registry
	vmreq         *require.RequireModule
	intrnbuffs    map[*iorw.Buffer]*iorw.Buffer
	includedpgrms map[string]*goja.Program
	rntmeche      map[int]map[string]interface{}
	//glblobjstoremove []string
}

func (atvrntme *atvruntime) decrntimecache() {
	if mpl := len(atvrntme.rntmeche); mpl > 0 {
		mp := atvrntme.rntmeche[mpl-1]
		for k := range mp {
			mp[k] = nil
			delete(mp, k)
		}
		delete(atvrntme.rntmeche, mpl-1)
	}
}

func (atvrntme *atvruntime) incrntimecache() map[string]interface{} {
	if atvrntme.rntmeche == nil {
		atvrntme.rntmeche = map[int]map[string]interface{}{}
	}
	mpi := len(atvrntme.rntmeche)
	atvrntme.rntmeche[mpi] = map[string]interface{}{}
	return atvrntme.rntmeche[mpi]
}

func (atvrntme *atvruntime) rntimecache() map[string]interface{} {
	if mpl := len(atvrntme.rntmeche); mpl > 0 {
		return atvrntme.rntmeche[mpl-1]
	}
	return atvrntme.incrntimecache()
}

func (atvrntme *atvruntime) InvokeFunction(functocall interface{}, args ...interface{}) (result interface{}) {
	if functocall != nil {
		if atvrntme.vm != nil {
			var fnccallargs []goja.Value = nil
			var argsn = 0

			for argsn < len(args) {
				if fnccallargs == nil {
					fnccallargs = make([]goja.Value, len(args))
				}
				fnccallargs[argsn] = atvrntme.vm.ToValue(args[argsn])
				argsn++
			}
			if atvfunc, atvfuncok := functocall.(func(goja.FunctionCall) goja.Value); atvfuncok {
				if len(fnccallargs) == 0 || fnccallargs == nil {
					fnccallargs = []goja.Value{}
				}
				var funccll = goja.FunctionCall{This: goja.Undefined(), Arguments: fnccallargs}
				if rsltval := atvfunc(funccll); rsltval != nil {
					result = rsltval.Export()
				}
			}
		}
	}
	return result
}

func (atvrntme *atvruntime) run() (val interface{}, err error) {
	var objmapref map[string]interface{} = nil
	if atvrntme.atv != nil && atvrntme.atv.ObjectMapRef != nil {
		objmapref = atvrntme.atv.ObjectMapRef()
	}
	val, err = atvrntme.corerun(parsing.Code(atvrntme.prsng), objmapref)
	return
}

func (atvrntme *atvruntime) corerun(code string, objmapref map[string]interface{}, includelibs ...string) (val interface{}, err error) {
	if code != "" {
		if atvrntme.vm != nil {
			atvrntme.vm.ClearInterrupt()
		}
		if len(objmapref) > 0 {
			atvrntme.lclvm(objmapref)
		}
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("%v", r)
				}
			}()
			isrequired := false
			if code, isrequired, err = transformCode(code, nil); err == nil {
				if isrequired {
					code = `function ` + `_vmrequire(args) { return require([args]);}` + code
				}
				prsd, prsderr := parser.ParseFile(nil, "", code, 0)
				if prsderr != nil {
					err = prsderr
				}

				if err == nil {
					if len(includelibs) > 0 {
						for _, incllib := range includelibs {
							if incllib == "require.js" || incllib == "require.min.js" {
								if _, included := atvrntme.includedpgrms[incllib]; included {
									continue
								} else {
									if requirejsprgm != nil {
										if _, err = atvrntme.vm.RunProgram(requirejsprgm); err != nil {
											break
										} else {
											atvrntme.includedpgrms[incllib] = requirejsprgm
										}
									}
								}
							}
						}
					}
					if p, perr := goja.CompileAST(prsd, false); perr == nil {
						_, err = atvrntme.lclvm(objmapref).RunProgram(p)
					} else {
						err = perr
					}
				}
			}
		}()
		if err != nil {
			if errs := err.Error(); errs != "" && !strings.HasPrefix(errs, "exit at <eval>:") {
				//fmt.Println(err.Error())
				//fmt.Println(code)
				//err = nil
			}

		}
	} else {
		if err != nil {
			if errs := err.Error(); errs != "" && !strings.HasPrefix(errs, "exit at <eval>:") {
				//fmt.Println(err.Error())
				//fmt.Println(code)
				//err = nil
			}
		}
	}
	if err != nil {
		cde := ""
		excpath := ""
		if atvrntme.prsng != nil {
			excpath = atvrntme.prsng.Prsvpth
		}
		for cdn, cd := range strings.Split(code, "\n") {
			cde += fmt.Sprintf("%d:%s\r\n", (cdn + 1), strings.TrimSpace(cd))
		}
		cdeerr := codeException(cde, excpath, err)
		err = cdeerr
	}
	return
}

func transformCode(code string, opts map[string]interface{}) (trsnfrmdcde string, isrequired bool, err error) {
	trsnfrmdcde = code
	isrequired = strings.Contains(code, "require(\"")
	if isrequired {
		//trsnfrmdcde = strings.Replace(trsnfrmdcde, "require(\"", "_vmrequire(\"", -1)
	}
	return
}

func (atvrntme *atvruntime) parseEval(forceCode bool, a ...interface{}) (val interface{}, err error) {
	return parsing.ParseEval(atvrntme.prsng, forceCode, atvrntme.corerun, a...)
}

func (atvrntme *atvruntime) removeBuffer(buff *iorw.Buffer) {
	if len(atvrntme.intrnbuffs) > 0 {
		if bf, bfok := atvrntme.intrnbuffs[buff]; bfok && bf == buff {
			atvrntme.intrnbuffs[buff] = nil
			delete(atvrntme.intrnbuffs, buff)
		}
	}
}

func (atvrntme *atvruntime) passiveout(i int) {
	if atvrntme != nil && atvrntme.prsng != nil {
		parsing.Passiveout(atvrntme.prsng, i)
	}
}

func (atvrntme *atvruntime) dispose(cleanupVal func(vali interface{}, valt reflect.Type), clear ...bool) {
	if atvrntme != nil {
		var clearonly = len(clear) > 0 && clear[0]
		if atvrntme.prsng != nil {
			atvrntme.prsng.Dispose()
			atvrntme.prsng = nil
		}
		if atvrntme.atv != nil {
			if !clearonly {
				atvrntme.atv = nil
			}
		}
		if atvrntme.vm != nil {
			resetvm(atvrntme.vm, cleanupVal)
			if !clearonly {
				vmpool.Put(atvrntme.vm)
				atvrntme.vm = nil
			}
		}
		if atvrntme.vmregister != nil {
			if !clearonly {
				atvrntme.vmregister = nil
			}
		}
		if atvrntme.vmreq != nil {
			atvrntme.vmreq = nil
		}
		if atvrntme.includedpgrms != nil {
			if il := len(atvrntme.includedpgrms); il > 0 {
				includedpgrms := make([]string, il)
				incldsi := 0
				for include := range atvrntme.includedpgrms {
					includedpgrms[incldsi] = include
					incldsi++
				}
				for len(includedpgrms) > 0 {
					atvrntme.includedpgrms[includedpgrms[0]] = nil
					delete(atvrntme.includedpgrms, includedpgrms[0])
					includedpgrms = includedpgrms[1:]
				}
			}
			if !clearonly {
				atvrntme.includedpgrms = nil
			}
		}
		if atvrntme.intrnbuffs != nil {
			if il := len(atvrntme.intrnbuffs); il > 0 {
				bfs := make([]*iorw.Buffer, il)
				bfsi := 0
				for bf := range atvrntme.intrnbuffs {
					bfs[bfsi] = bf
					bfsi++
				}
				for len(bfs) > 0 {
					bf := bfs[0]
					bf.Close()
					bf = nil
					bfs = bfs[1:]
				}
			}
			if !clearonly {
				atvrntme.intrnbuffs = nil
			}
		}
		if !clearonly {
			atvrntme = nil
		}
	}
}

func defaultAtvRuntimeInternMap(atvrntme *atvruntime) (internmapref map[string]interface{}) {
	internmapref = map[string]interface{}{
		"buffer": func() (buff *iorw.Buffer) {
			buff = iorw.NewBuffer()
			buff.OnClose = atvrntme.removeBuffer
			atvrntme.intrnbuffs[buff] = buff
			return
		},
		"sleep": func(mils int64) {
			time.Sleep(time.Millisecond * time.Duration(mils))
		},
		"_psvout": func(i int) {
			atvrntme.passiveout(i)
		},
		"_cache": func() map[string]interface{} {
			return atvrntme.rntimecache()
		},
		"_inccache": func() map[string]interface{} {
			return atvrntme.incrntimecache()
		},
		"_deccache": func() map[string]interface{} {
			return atvrntme.incrntimecache()
		},
		"_parseEval": func(a ...interface{}) (val interface{}, err error) {
			return atvrntme.parseEval(true, a...)
		},
		"parseEval": func(a ...interface{}) (val interface{}, err error) {
			return atvrntme.parseEval(false, a...)
		},
		//WRITER
		"incprint": func(w io.Writer) {
			if atvrntme.prsng != nil {
				atvrntme.prsng.IncPrint(w)
			}
		},
		"resetprint": func() {
			if atvrntme.prsng != nil {
				atvrntme.prsng.ResetPrint()
			}
		},
		"decprint": func() {
			if atvrntme.prsng != nil {
				atvrntme.prsng.DecPrint()
			}
		},
		"print": func(a ...interface{}) {
			if atvrntme.prsng != nil {
				atvrntme.prsng.Print(a...)
			} else if atvrntme.atv != nil {
				atvrntme.atv.print(nil, a...)
			}
		},
		"println": func(a ...interface{}) {
			if atvrntme.prsng != nil {
				atvrntme.prsng.Println(a...)
			} else if atvrntme.atv != nil {
				atvrntme.atv.println(nil, a...)
			}
		},
		"binwrite": func(b ...byte) (n int, err error) {
			if atvrntme.prsng != nil {
				n, err = atvrntme.prsng.BinWrite(b...)
			} else if atvrntme.atv != nil {
				n, err = atvrntme.atv.binwrite(nil, b...)
			}
			return
		},
		//READER
		"incread": func(r io.Reader) {
			if atvrntme.prsng != nil {
				atvrntme.prsng.IncRead(r)
			}
		},
		"resetread": func() {
			if atvrntme.prsng != nil {
				atvrntme.prsng.ResetRead()
			}
		},
		"decread": func() {
			if atvrntme.prsng != nil {
				atvrntme.prsng.DecRead()
			}
		},
		"seek": func(offset int64, whence int) (n int64, err error) {
			if atvrntme.prsng != nil {
				n, err = atvrntme.prsng.Seek(offset, whence)
			} else if atvrntme.atv != nil {
				n, err = atvrntme.atv.seek(nil, offset, whence)
			}
			return
		},
		"readln": func() (ln string, err error) {
			if atvrntme.prsng != nil {
				ln, err = atvrntme.prsng.ReadLn()
			} else if atvrntme.atv != nil {
				ln, err = atvrntme.atv.readln(nil)
			}
			return
		},
		"readLines": func() (lines []string, err error) {
			if atvrntme.prsng != nil {
				lines, err = atvrntme.prsng.ReadLines()
			} else if atvrntme.atv != nil {
				lines, err = atvrntme.atv.readlines(nil)
			}
			return
		}, "readAll": func() (s string, err error) {
			if atvrntme.prsng != nil {
				s, err = atvrntme.prsng.ReadAll()
			} else if atvrntme.atv != nil {
				s, err = atvrntme.atv.readAll(nil)
			}
			return
		}, "binread": func(size int) (b []byte, err error) {
			if atvrntme.prsng != nil {
				b, err = atvrntme.prsng.BinRead(size)
			} else if atvrntme.atv != nil {
				b, err = atvrntme.atv.binread(nil, size)
			}
			return
		}, "_scriptinclude": func(url string, a ...interface{}) (src interface{}, srcerr error) {
			if atvrntme.prsng != nil && atvrntme.prsng.AtvActv != nil {
				if lkpr, lkprerr := atvrntme.prsng.AtvActv.AltLookupTemplate(url, a...); lkpr != nil && lkprerr == nil {
					if s, _ := iorw.ReaderToString(lkpr); s != "" {
						src = strings.TrimSpace(s)
					} else {
						src = s
					}
				} else if lkprerr != nil {
					srcerr = lkprerr
				}
			}
			if src == nil {
				src = ""
			}
			return
		},
		"script": atvrntme}
	return
}

func newatvruntime(atv *Active) (atvrntme *atvruntime, err error) {
	atvrntme = &atvruntime{atv: atv, includedpgrms: map[string]*goja.Program{}, intrnbuffs: map[*iorw.Buffer]*iorw.Buffer{}}
	atvrntme.atv.InterruptVM = func(v interface{}) {
		atvrntme.lclvm().Interrupt(v)
	}
	return
}

type fieldmapper struct {
	fldmppr goja.FieldNameMapper
}

// FieldName returns a JavaScript name for the given struct field in the given type.
// If this method returns "" the field becomes hidden.
func (fldmppr *fieldmapper) FieldName(t reflect.Type, f reflect.StructField) (fldnme string) {
	if f.Tag != "" {
		fldnme = f.Tag.Get("json")
	} else {
		fldnme = uncapitalize(t.Name()) // fldmppr.fldmppr.FieldName(t, f)
	}
	return
}

// MethodName returns a JavaScript name for the given method in the given type.
// If this method returns "" the method becomes hidden.
func (fldmppr *fieldmapper) MethodName(t reflect.Type, m reflect.Method) (mthdnme string) {
	mthdnme = uncapitalize(m.Name)
	return
}

func uncapitalize(s string) (nme string) {
	if sl := len(s); sl > 0 {
		var nrxtsr = rune(0)
		for sn, sr := range s {
			if 'A' <= sr && sr <= 'Z' {
				sr += 'a' - 'A'
				nme += string(sr)
			} else {
				nme += string(sr)
			}
			if sn <= (sl-1)-1 {
				nrxtsr = rune(s[sn+1])
			} else {
				nrxtsr = rune(0)
			}
			if 'a' <= nrxtsr && nrxtsr <= 'z' {
				nme += s[sn+1:]
				break
			}
		}
	}
	return nme
}

var vmpool = &sync.Pool{New: func() interface{} { return newvm() }}

func newvm() (vm *goja.Runtime) {
	vm = goja.New()
	var fldmppr = &fieldmapper{fldmppr: goja.UncapFieldNameMapper()}
	vm.SetFieldNameMapper(fldmppr)
	return
}

func resetvm(vm *goja.Runtime, cleanupVal func(vali interface{}, valt reflect.Type)) {
	if vm != nil {
		if vmgbl := vm.GlobalObject(); vmgbl != nil {
			var ks = vmgbl.Keys()
			var rsetcode = ""

			if len(ks) > 0 {
				for _, k := range ks {
					if vmgblval := vm.GlobalObject().Get(k); vmgblval != nil {
						var vali = vmgblval.Export()
						var valt = vmgblval.ExportType()
						if vali != nil && valt != nil && cleanupVal != nil {
							cleanupVal(vali, valt)
						}
					}
					vm.GlobalObject().Delete(k)
					rsetcode += k + "=undefined;\n"
				}
				vm.RunString(rsetcode)
			}
		}
	}
}

func (atvrntme *atvruntime) lclvm(objmapref ...map[string]interface{}) (vm *goja.Runtime) {
	if atvrntme != nil {
		if atvrntme.vm == nil {
			vm, _ := vmpool.Get().(*goja.Runtime)
			if vm == nil {
				vm = newvm()
			}
			var dne = make(chan bool, 1)
			if atvrntme.vmregister == nil {
				vmregister := require.NewRegistryWithLoader(func(path string) (sourcebytes []byte, sourceerr error) {
					if atvrntme != nil && atvrntme.atv != nil && atvrntme.atv.LookupTemplate != nil {
						if lkprdr, lkprdrerr := atvrntme.atv.LookupTemplate(path); lkprdr != nil || lkprdrerr != nil {
							if lkprdrerr == nil && lkprdr != nil {
								buf := new(bytes.Buffer)
								_, sourceerr = buf.ReadFrom(lkprdr)
								if sourcebytes = buf.Bytes(); len(sourcebytes) > 0 && (sourceerr == nil || sourceerr == io.EOF) {
									sourceerr = nil
								} else if len(sourcebytes) == 0 && sourceerr == nil {
									sourcebytes = nil
									sourceerr = require.InvalidModuleError
								}
								if sourceerr != nil && sourceerr != io.EOF {
									sourcebytes = nil
								}
								return sourcebytes, sourceerr
							} else if lkprdr == nil && lkprdrerr == nil {
								return require.DefaultSourceLoader(path)
							}
						} else {
							return nil, lkprdrerr
						}
					}
					return nil, require.ModuleFileDoesNotExistError
				})
				atvrntme.vmregister = vmregister
			}
			if atvrntme.vmreq == nil {
				vmreq := atvrntme.vmregister.Enable(vm)
				atvrntme.vmreq = vmreq
			}
			//go func(vm *goja.Runtime) {
			//	defer func() { dne <- true }()
			jsext.Register(vm)
			atvrntme.vm = vm
			//}(atvrntme.vm)
			//<-dne
			if definternmapref := defaultAtvRuntimeInternMap(atvrntme); len(definternmapref) > 0 {
				if len(definternmapref) > 0 {
					for k, ref := range definternmapref {
						atvrntme.vm.Set(k, ref)
					}
				}
			}
			if len(objmapref) > 0 && objmapref[0] != nil {
				for k, ref := range objmapref[0] {
					atvrntme.vm.Set(k, ref)
				}
			}
			var modstoload = RetrieveModule()
			modstoload = append([]*goja.Program{requirejsprgm}, modstoload...)
			go loadInternalModules(dne, atvrntme.vm, modstoload...)
			<-dne
			defer close(dne)
		} else {
			if len(objmapref) > 0 && objmapref[0] != nil {
				for k, ref := range objmapref[0] {
					atvrntme.vm.Set(k, ref)
				}
			}
		}
		vm = atvrntme.vm
	}
	return
}

func loadInternalModules(dne chan bool, vm *goja.Runtime, prgrms ...*goja.Program) {
	defer func() { dne <- true }()
	if len(prgrms) > 0 {
		for _, prgm := range prgrms {
			if prgm != nil {
				if _, err := vm.RunProgram(prgm); err != nil {

				}
			}
		}
	}
}

var requirejsprgm *goja.Program = nil

var globalModules map[string]*goja.Program

//var globalModuleslck *sync.RWMutex

func RetrieveModule(modulepath ...string) (modules []*goja.Program) {

	if len(globalModules) > 0 {
		if len(modulepath) > 0 {
			var glblmodpths []string = nil
			func() {
				//globalModuleslck.RLock()
				//defer globalModuleslck.RUnlock()
				if glblmdpthsl := len(globalModules); glblmdpthsl > 0 {
					glblmodpths = make([]string, glblmdpthsl)
					var glblmodpthsi = 0
					for mdpth := range globalModules {
						glblmodpths[glblmodpthsi] = mdpth
						glblmodpthsi++
					}
				}
			}()
			if len(glblmodpths) > 0 {
				var modpthsi = 0
				var modpthsl = len(modulepath)
				for modpthsi < modpthsl {
					for _, glgmdpth := range glblmodpths {
						if modulepath[modpthsi] != glgmdpth {
							modulepath = append(modulepath[0:modpthsi], modulepath[modpthsi])
							modpthsl--
						} else {
							modpthsi++
						}
					}
				}

				if modpthsl > 0 {
					func() {
						//globalModuleslck.RLock()
						//defer globalModuleslck.RUnlock()
						modules = make([]*goja.Program, modpthsl)
						var modulesi = 0
						for _, mdpth := range modulepath {
							modules[modulesi] = globalModules[mdpth]
							modulesi++
						}
					}()
				}
			}
		} else {
			func() {
				if glblmdsl := len(globalModules); glblmdsl > 0 {
					//globalModuleslck.RLock()
					//defer globalModuleslck.RUnlock()
					modules = make([]*goja.Program, glblmdsl)
					var modulesi = 0
					for mdpth := range globalModules {
						modules[modulesi] = globalModules[mdpth]
						modulesi++
					}
				}
			}()
		}
	}
	return
}

func LoadGlobalModule(modulepath string, a ...interface{}) (err error) {
	if modulepath != "" && len(a) > 0 {
		var bufcde = iorw.NewBuffer()
		defer bufcde.Close()
		bufcde.Print(a...)
		if bufcde.Size() > 0 {
			var modulepgrm *goja.Program = nil
			if modulepgrm, err = goja.Compile("", bufcde.String(), false); modulepgrm != nil && err == nil {
				func() {
					//globalModuleslck.Lock()
					//defer globalModuleslck.Unlock()
					if globalModules[modulepath] != nil {
						globalModules[modulepath] = nil
					}
					globalModules[modulepath] = modulepgrm
				}()
			}
		}
	}
	return
}

func UnloadGlobalModule(modulepath string) (existed bool) {
	if modulepath != "" {
		func() {
			//globalModuleslck.Lock()
			//defer globalModuleslck.Unlock()
			if existed = globalModules[modulepath] != nil; existed {
				globalModules[modulepath] = nil
				delete(globalModules, modulepath)
			}
		}()
	}
	return
}

func init() {
	globalModules = map[string]*goja.Program{}
	//globalModuleslck = &sync.RWMutex{}
	/*var errpgrm error = nil
	if requirejsprgm, errpgrm = goja.Compile("", requirejs.RequireJSString(), false); errpgrm != nil {
		fmt.Println(errpgrm.Error())
	}*/
}

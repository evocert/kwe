package active

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/dop251/goja/parser"

	"github.com/dop251/goja"

	//"github.com/evocert/kwe/ecma/es51"
	//"github.com/evocert/kwe/ecma/es51/parser"

	"github.com/evocert/kwe/ecma/jsext"
	"github.com/evocert/kwe/requirejs"

	"github.com/evocert/kwe/iorw"
)

//Active - struct
type Active struct {
	Namespace      string
	Print          func(a ...interface{})
	Println        func(a ...interface{})
	FPrint         func(w io.Writer, a ...interface{})
	FPrintLn       func(w io.Writer, a ...interface{})
	LookupTemplate func(string, ...interface{}) (io.Reader, error)
	ObjectMapRef   func() map[string]interface{}
	lckprnt        *sync.Mutex
	InterruptVM    func(v interface{})
	*atvruntime
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
		if imprtglbs != nil && len(imprtglbs) > 0 {
			if gbl := atv.atvruntime.vm.GlobalObject(); gbl != nil {
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

//NewActive - instance
func NewActive(namespace ...string) (atv *Active) {
	atv = &Active{lckprnt: &sync.Mutex{}, Namespace: "", atvruntime: nil}
	atv.atvruntime, _ = newatvruntime(atv, nil)
	if len(namespace) == 1 && namespace[0] != "" {
		atv.Namespace = namespace[0] + "."
	}
	return
}

func (atv *Active) namespace() string {
	if atv.Namespace != "" {
		return atv.Namespace
	}
	return ""
}

func (atv *Active) print(w io.Writer, a ...interface{}) {
	if prntr, prntrok := w.(iorw.Printer); prntrok {
		prntr.Print(a...)
	} else {
		if atv.Print != nil {
			if len(a) > 0 {
				atv.lckprnt.Lock()
				defer atv.lckprnt.Unlock()
				atv.Print(a...)
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
				atv.lckprnt.Lock()
				defer atv.lckprnt.Unlock()
				atv.Println(a...)
			}
		} else if atv.FPrintLn != nil && w != nil {
			atv.lckprnt.Lock()
			defer atv.lckprnt.Unlock()
			atv.FPrint(w, a...)
		} else if w != nil {
			if prntr, prntrok := w.(iorw.Printer); prntrok {
				prntr.Println(a...)
			} else {
				if len(a) > 0 {
					atv.lckprnt.Lock()
					defer atv.lckprnt.Unlock()
					fmt.Fprint(w, a...)
				}
				fmt.Fprintln(w)
			}
		}
	}
}

//InvokeVM invoke vm
func (atv *Active) InvokeVM(callback func(vm *goja.Runtime) error) {
	if callback != nil {
		callback(atv.vm())
	}
}

func (atv *Active) vm() (vm *goja.Runtime) {
	if atv != nil && atv.atvruntime != nil && atv.atvruntime.vm != nil {
		vm = atv.atvruntime.vm
	}
	return
}

func (atv *Active) atvrun(prsng *parsing) {
	if prsng != nil {
		if atv.atvruntime == nil {
			atv.atvruntime, _ = newatvruntime(atv, prsng)
		} else {
			if atv.prsng != prsng {
				atv.prsng.dispose()
				atv.prsng = nil
				atv.prsng = prsng
			}
		}
		if atv.atvruntime != nil {
			atv.atvruntime.run()
		}
	}
}

//Eval - parse rin io.Reader, execute if neaded and output to wou io.Writer
func (atv *Active) Eval(wout io.Writer, rin io.Reader) {
	var parsing = nextparsing(atv, nil, wout)
	var rnr io.RuneReader = nil
	var bfr *bufio.Reader = nil
	if rr, rrok := rin.(io.RuneReader); rrok {
		rnr = rr
	} else {
		bfr = bufio.NewReader(rin)
		rnr = bfr
	}
	parseprsngrunerdr(parsing, rnr, true)

	parsing.Close()
}

//Close - refer to  io.Closer
func (atv *Active) Close() (err error) {
	if atv.lckprnt != nil {
		atv.lckprnt = nil
	}
	if atv.atvruntime != nil {
		atv.atvruntime.dispose()
		atv.atvruntime = nil
	}
	return
}

//Interrupt - Active processing
func (atv *Active) Interrupt() {
	if atv.InterruptVM != nil {
		atv.InterruptVM("exit")
	}
}

var prslbl = [][]rune{[]rune("<@"), []rune("@>")}

type parsing struct {
	*iorw.Buffer
	atv *Active
	//atvrntme       *atvruntime
	wout           io.Writer
	prntrs         []io.Writer
	prslbli        []int
	prslblprv      []rune
	prntprsng      *parsing
	foundcde       bool
	hascde         bool
	cdetxt         rune
	cdeoffsetstart int64
	cdeoffsetend   int64
	cdemap         map[int][]int64
	cder           []rune
	cderi          int
	psvoffsetstart int64
	psvoffsetend   int64
	psvmap         map[int][]int64
	psvr           []rune
	psvri          int
	psvctrl        *passivectrl
	prvpsvctrls    map[*passivectrl]*passivectrl
}

func (prsng *parsing) print(a ...interface{}) {
	if prsng.atv != nil {
		w := prsng.wout
		if pl := len(prsng.prntrs); pl > 0 {
			w = prsng.prntrs[pl-1]
		}
		prsng.atv.print(w, a...)
	}
}

func (prsng *parsing) println(a ...interface{}) {
	if prsng.atv != nil {
		w := prsng.wout
		if pl := len(prsng.prntrs); pl > 0 {
			w = prsng.prntrs[pl-1]
		}
		prsng.atv.println(w, a...)
	}
}

func (prsng *parsing) incprint(w io.Writer) {
	if prsng != nil {
		prsng.prntrs = append(prsng.prntrs, w)
	}
}

func (prsng *parsing) resetprint() {
	if prsng.prntrs != nil {
		for len(prsng.prntrs) > 0 {
			prsng.prntrs[len(prsng.prntrs)-1] = nil
			prsng.prntrs = prsng.prntrs[:len(prsng.prntrs)-1]
		}
	}
}

func (prsng *parsing) decprint() {
	if prsng.prntrs != nil {
		if len(prsng.prntrs) > 0 {
			prsng.prntrs[len(prsng.prntrs)-1] = nil
			prsng.prntrs = prsng.prntrs[:len(prsng.prntrs)-1]
		}
	}
}

func (prsng *parsing) dispose() {
	if prsng != nil {
		if prsng.cdemap != nil {
			if len(prsng.cdemap) > 0 {
				var cdeks = make([]int, len(prsng.cdemap))
				var cdeksi = 0
				for cdek := range prsng.cdemap {
					cdeks[cdeksi] = cdek
					prsng.cdemap[cdek] = nil
				}
				for len(cdeks) > 0 {
					delete(prsng.cdemap, cdeks[0])
					cdeks = cdeks[1:]
				}
				cdeks = nil
			}
			prsng.cdemap = nil
		}
		if prsng.prntrs != nil {
			for len(prsng.prntrs) > 0 {
				prsng.prntrs[len(prsng.prntrs)-1] = nil
				prsng.prntrs = prsng.prntrs[:len(prsng.prntrs)-1]
			}
			prsng.prntrs = nil
		}
		if prsng.atv != nil {
			prsng.atv = nil
		}
		if prsng.prntprsng != nil {
			prsng.prntprsng = nil
		}
		if prsng.prslbli != nil {
			prsng.prslbli = nil
		}
		if prsng.prslblprv != nil {
			prsng.prslblprv = nil
		}
		if prsng.Buffer != nil {
			prsng.Buffer.Close()
			prsng.Buffer = nil
		}
		if prsng.wout != nil {
			prsng.wout = nil
		}
		if prsng.psvctrl != nil {
			prsng.psvctrl.close()
			prsng.psvctrl = nil
		}
		prsng = nil
	}
}

func (prsng *parsing) topprsng() *parsing {
	if prsng.prntprsng == nil {
		return prsng
	}
	return prsng.prntprsng.topprsng()
}

func (prsng *parsing) setcdepos(startoffset int64, endoffset int64) {
	if prsng.cdemap == nil {
		prsng.cdemap = map[int][]int64{}
	}
	prsng.cdemap[len(prsng.cdemap)] = []int64{startoffset, endoffset}
}

func (prsng *parsing) setpsvpos(startoffset int64, endoffset int64) (pos int) {
	if prsng.psvmap == nil {
		prsng.psvmap = map[int][]int64{}
	}
	pos = len(prsng.psvmap)
	prsng.psvmap[pos] = []int64{startoffset, endoffset}
	return
}

func (prsng *parsing) writePsv(p []rune) (err error) {
	if pl := len(p); pl > 0 {
		if prsng.foundCode() {
			if prsng.psvoffsetstart == -1 {
				prsng.psvoffsetstart = prsng.Size()
			}
			err = prsng.WriteRunes(p[:pl]...)
		} else {
			_, err = prsng.wout.Write([]byte(string(p[:pl])))
		}
	}
	return
}

func (prsng *parsing) writeCde(p []rune) (err error) {
	if pl := len(p); pl > 0 {
		if prsng.cdeoffsetstart == -1 {
			prsng.cdeoffsetstart = prsng.Size()
		}
		err = prsng.WriteRunes(p[:pl]...)
	}
	return
}

func (prsng *parsing) foundCode() bool {
	if prsng.foundcde {
		return true
	}
	return false
}

func (prsng *parsing) flushPsv() (err error) {
	if pi := prsng.psvri; pi > 0 {
		prsng.psvri = 0
		if prsng.psvctrl != nil && prsng.psvctrl.lastElmType == ElemStart {
			err = prsng.psvctrl.cachedbuf().WriteRunes(prsng.psvr[:pi]...)
		} else {
			err = prsng.writePsv(prsng.psvr[:pi])
		}
	}
	if (prsng.psvctrl == nil || prsng.psvctrl.lastElmType == elemnone) && prsng.foundCode() {
		if psvoffsetstart := prsng.psvoffsetstart; psvoffsetstart > -1 {
			prsng.psvoffsetstart = -1
			pos := prsng.setpsvpos(psvoffsetstart, prsng.Size())
			err = parseatvrunes(prsng, []rune(fmt.Sprintf("_passiveout(%d);", pos)))
		}
	}
	return
}

func parsepsvrunes(prsng *parsing, p []rune) (err error) {
	if len(p) > 0 {
		for _, rn := range p {
			if err = parsepsvrune(prsng, rn); err != nil {
				break
			}
		}
	}
	return
}

func parseatvrunes(prsng *parsing, p []rune) (err error) {
	if len(p) > 0 {
		for _, rn := range p {
			if err = parseatvrune(prsng, rn); err != nil {
				break
			}
		}
	}
	return
}

func (prsng *parsing) flushCde() (err error) {
	if pi := prsng.cderi; pi > 0 {
		prsng.cderi = 0
		prsng.writeCde(prsng.cder[:pi])
	}
	if cdeoffsetstart := prsng.cdeoffsetstart; cdeoffsetstart > -1 {
		prsng.cdeoffsetstart = -1
		prsng.setcdepos(cdeoffsetstart, prsng.Size())
	}
	return
}

func parseprsngrunerdr(prsng *parsing, rnr io.RuneReader, canexec bool) (err error) {
	var crunes = make([]rune, 4096)
	var crunesi = 0
	for err == nil {
		r, rsize, rerr := rnr.ReadRune()
		if rsize > 0 {
			crunes[crunesi] = r
			crunesi++
			if crunesi == len(crunes) {
				cl := crunesi
				crunesi = 0
				for _, cr := range crunes[:cl] {
					parseprsng(prsng, prsng.prslbli, prsng.prslblprv, cr)
				}
			}
		}
		if rerr != nil {
			err = rerr
		}
	}
	if err == io.EOF || err == nil {
		if crunesi > 0 {
			cl := crunesi
			crunesi = 0
			for _, cr := range crunes[:cl] {
				parseprsng(prsng, prsng.prslbli, prsng.prslblprv, cr)
			}
		}
		prsng.flushPsv()
		prsng.flushCde()
		if canexec {
			if prsng.foundCode() {
				prsng.atv.atvrun(prsng)
			} else {
				if rdr := prsng.Reader(); rdr != nil {
					io.Copy(prsng.wout, rdr)
					rdr.Close()
					rdr = nil
				}
			}
		}
	}
	return
}

func parseprsng(prsng *parsing, prslbli []int, prslblprv []rune, pr rune) (err error) {
	if prsng.psvctrl != nil && prsng.psvctrl.lastElmType == ElemStart {
		err = prsng.psvctrl.processrn(pr)
	} else {
		if prsng.cdetxt == rune(0) && prslbli[1] == 0 && prslbli[0] < len(prslbl[0]) {
			if prslbli[0] > 0 && prslbl[0][prslbli[0]-1] == prslblprv[0] && prslbl[0][prslbli[0]] != pr {
				if psvl := prslbli[0]; psvl > 0 {
					prslbli[0] = 0
					prslblprv[0] = 0
					err = parsepsvrunes(prsng, prslbl[0][0:psvl])
				}
			}
			if prslbl[0][prslbli[0]] == pr {
				prslbli[0]++
				if prslbli[0] == len(prslbl[0]) {

					prslblprv[0] = 0
				} else {
					prslblprv[0] = pr
				}
			} else {
				if psvl := prslbli[0]; psvl > 0 {
					prslbli[0] = 0
					prslblprv[0] = 0
					err = parsepsvrunes(prsng, prslbl[0][0:psvl])
				}
				err = parsepsvrune(prsng, pr)
				prslblprv[0] = pr
			}
		} else if prslbli[0] == len(prslbl[0]) && prslbli[1] < len(prslbl[1]) {
			if prsng.cdetxt == rune(0) && prslbl[1][prslbli[1]] == pr {
				prslbli[1]++
				if prslbli[1] == len(prslbl[1]) {
					prslbli[0] = 0
					prslblprv[1] = 0
					prslbli[1] = 0
				} else {
					prslblprv[1] = pr
				}
			} else {
				if prsl := prslbli[1]; prsl > 0 {
					prslbli[1] = 0
					err = parseatvrunes(prsng, prslbl[1][:prsl])
				}
				err = parseatvrune(prsng, pr)
				prslblprv[1] = pr
			}
		}
	}
	return
}

func nextparsing(atv *Active, prntprsng *parsing, wout io.Writer) (prsng *parsing) {
	prsng = &parsing{Buffer: iorw.NewBuffer(), wout: wout, prntprsng: prntprsng, atv: atv, cdetxt: rune(0), prslbli: []int{0, 0}, prslblprv: []rune{0, 0}, cdeoffsetstart: -1, cdeoffsetend: -1, psvoffsetstart: -1, psvoffsetend: -1, psvr: make([]rune, 8192), cder: make([]rune, 8192), prntrs: []io.Writer{}}
	return
}

type atvruntime struct {
	prsng         *parsing
	atv           *Active
	vm            *goja.Runtime
	intrnbuffs    map[*iorw.Buffer]*iorw.Buffer
	includedpgrms map[string]*goja.Program
	//glblobjstoremove []string
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
	val, err = atvrntme.corerun(atvrntme.code(), objmapref)
	return
}

func (atvrntme *atvruntime) corerun(code string, objmapref map[string]interface{}, includelibs ...string) (val interface{}, err error) {
	if code != "" {
		atvrntme.vm.ClearInterrupt()
		if objmapref != nil && len(objmapref) > 0 {
			for k, ref := range objmapref {
				atvrntme.vm.Set(k, ref)
			}
		}
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("%v", r)
				}
			}()
			isrequired := false
			if code, isrequired, err = transformCode(code, atvrntme.atv.namespace(), nil); err == nil {
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
						_, err = atvrntme.vm.RunProgram(p)
						if err != nil {
							fmt.Println(err.Error())
						}
					} else {
						err = perr
					}
				}
			}
		}()
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(code)
		}
	} else {
		if err != nil {
			fmt.Println(err.Error())
			//fmt.Println(code)
		}
	}
	return
}

/*var (
	defaultOpts = map[string]interface{}{
		"presets":       []interface{}{"es2015"},
		"ast":           false,
		"sourceMaps":    false,
		"babelrc":       false,
		"compact":       false,
		"retainLines":   true,
		"highlightCode": false,
	}
)*/

func transformCode(code string, namespace string, opts map[string]interface{}) (trsnfrmdcde string, isrequired bool, err error) {
	trsnfrmdcde = code
	isrequired = strings.Index(code, "require(\"") > -1
	if isrequired {
		trsnfrmdcde = strings.Replace(trsnfrmdcde, "require(\"", "_vmrequire(\"", -1)
	}
	return
}

func (atvrntme *atvruntime) parseEval(forceCode bool, a ...interface{}) (val interface{}, err error) {
	var prsng = atvrntme.prsng
	if forceCode {
		prsng.prslbli[0] = len(prslbl[0])
		prsng.prslbli[1] = 0
		prsng.prslblprv[0] = 0
		prsng.prslblprv[1] = 0
	}
	var cdecoords []int64 = nil
	orgcdemapl := len(prsng.cdemap)
	pr, pw := io.Pipe()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(args ...interface{}) {
		defer func() {
			pw.Close()
		}()
		wg.Done()
		if len(args) > 0 {
			iorw.Fprint(pw, args...)
		}
	}(a...)
	rnr := bufio.NewReader(pr)
	wg.Wait()
	var crunes = make([]rune, 4096)
	var crunesi = 0
	var cdetestbuf = iorw.NewBuffer()
	func() {
		defer func() {
			if cdetestbuf != nil {
				if cdetestbuf.Size() > 0 {
					fmt.Println(cdetestbuf.String())
				}
				cdetestbuf.Close()
				cdetestbuf = nil
			}
		}()
		canprcss := false
		for err == nil {
			r, rsize, rerr := rnr.ReadRune()
			if rsize > 0 {
				if !canprcss {
					canprcss = true
				}
				crunes[crunesi] = r
				crunesi++
				if crunesi == len(crunes) {
					cl := crunesi
					crunesi = 0
					cdetestbuf.Print(string(crunes[:cl]))
					for _, cr := range crunes[:cl] {
						parseprsng(prsng, prsng.prslbli, prsng.prslblprv, cr)
					}
				}
			}
			if rerr != nil {
				err = rerr
			}
		}
		if err == io.EOF || err == nil {
			if err == io.EOF {
				err = nil
			}
			if canprcss {
				if crunesi > 0 {
					cl := crunesi
					crunesi = 0
					cdetestbuf.Print(string(crunes[:cl]))
					for _, cr := range crunes[:cl] {
						parseprsng(prsng, prsng.prslbli, prsng.prslblprv, cr)
					}
				}

				prsng.flushPsv()
				prsng.flushCde()
				if prsng.foundCode() {
					if cdemapl := len(prsng.cdemap); cdemapl > orgcdemapl {
						cdecoords = []int64{prsng.cdemap[orgcdemapl][0], prsng.cdemap[cdemapl-1][1]}
					}
					cde := atvrntme.code(cdecoords...)
					val, err = atvrntme.corerun(cde, nil)
				} else {
					if rdr := prsng.Reader(); rdr != nil {
						io.Copy(prsng.wout, rdr)
						rdr.Close()
						rdr = nil
					}
				}
				cdetestbuf.Close()
				cdetestbuf = nil
			}
		}
	}()
	return
}

func (atvrntme *atvruntime) removeBuffer(buff *iorw.Buffer) {
	if len(atvrntme.intrnbuffs) > 0 {
		if bf, bfok := atvrntme.intrnbuffs[buff]; bfok && bf == buff {
			atvrntme.intrnbuffs[buff] = nil
			delete(atvrntme.intrnbuffs, buff)
		}
	}
}

func (atvrntme *atvruntime) code(coords ...int64) (c string) {
	if atvrntme != nil && atvrntme.prsng != nil {
		if cdel := len(atvrntme.prsng.cdemap); cdel > 0 {
			var cdei = 0
			var rdr = atvrntme.prsng.Reader()
			if len(coords) == 0 {
				coords = []int64{atvrntme.prsng.cdemap[cdei][0], atvrntme.prsng.cdemap[cdel-1][1]}
			}
			for cdei < cdel {
				cdecoors := atvrntme.prsng.cdemap[cdei]
				cdei++
				if cdecoors[0] <= coords[1] {
					if cdecoors[0] == coords[1] {
						break
					}
					if cdecoors[1] <= coords[1] {
						if cdecoors[0] < coords[0] {
							cdecoors[0] = coords[0]
						}
					} else if cdecoors[1] > coords[1] {
						if cdecoors[0] >= coords[0] {
							cdecoors[0] = coords[0]
							cdecoors[1] = coords[0]
						} else {
							break
						}
					}
				}

				if cdecoors[1] > cdecoors[0] {
					rdr.MaxRead = -1
					rdr.Seek(cdecoors[0], 0)
					rdr.MaxRead = cdecoors[1] - cdecoors[0]
					if rs, rserr := iorw.ReaderToString(rdr); rs != "" {
						c += rs
					} else if rserr != nil {
						break
					}
				}
			}
			rdr.Close()
		}
	}
	return c
}

func (atvrntme *atvruntime) passiveout(i int) {
	if atvrntme != nil && atvrntme.prsng != nil {
		if psvl := len(atvrntme.prsng.psvmap); psvl > 0 && i >= 0 && i < psvl {
			psvcoors := atvrntme.prsng.psvmap[i]
			if psvcoors[1] > psvcoors[0] {
				rdr := atvrntme.prsng.Reader()
				rdr.Seek(psvcoors[0], 0)
				io.CopyN(atvrntme.prsng.wout, rdr, psvcoors[1]-psvcoors[0])
			}
		}
	}
}

func (atvrntme *atvruntime) dispose() {
	if atvrntme != nil {
		if atvrntme.prsng != nil {
			atvrntme.prsng.dispose()
			atvrntme.prsng = nil
		}
		if atvrntme.atv != nil {
			atvrntme.atv = nil
		}
		if atvrntme.vm != nil {
			if vmgbl := atvrntme.vm.GlobalObject(); vmgbl != nil {
				var ks = vmgbl.Keys()
				if len(ks) > 0 {
					for _, k := range ks {
						atvrntme.vm.GlobalObject().Delete(k)
					}
				}
			}
			atvrntme.vm = nil
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
			atvrntme.includedpgrms = nil
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
			atvrntme.intrnbuffs = nil
		}
		atvrntme = nil
	}
}

func defaultAtvRuntimeInternMap(atvrntme *atvruntime) (internmapref map[string]interface{}) {
	internmapref = map[string]interface{}{
		atvrntme.atv.namespace() + "newbuffer": func() (buff *iorw.Buffer) {
			buff = iorw.NewBuffer()
			buff.OnClose = atvrntme.removeBuffer
			atvrntme.intrnbuffs[buff] = buff
			return
		},
		"_passiveout": func(i int) {
			atvrntme.passiveout(i)
		},
		"_parseEval": func(a ...interface{}) (val interface{}, err error) {
			return atvrntme.parseEval(true, a...)
		},
		atvrntme.atv.namespace() + "parseEval": func(a ...interface{}) (val interface{}, err error) {
			return atvrntme.parseEval(false, a...)
		},
		atvrntme.atv.namespace() + "incprint": func(w io.Writer) {
			if atvrntme.prsng != nil {
				atvrntme.prsng.incprint(w)
			}
		},
		atvrntme.atv.namespace() + "resetprint": func() {
			if atvrntme.prsng != nil {
				atvrntme.prsng.resetprint()
			}
		},
		atvrntme.atv.namespace() + "decprint": func() {
			if atvrntme.prsng != nil {
				atvrntme.prsng.decprint()
			}
		},
		atvrntme.atv.namespace() + "print": func(a ...interface{}) {
			if atvrntme.prsng != nil {
				atvrntme.prsng.print(a...)
			}
		},
		atvrntme.atv.namespace() + "println": func(a ...interface{}) {
			if atvrntme.prsng != nil {
				atvrntme.prsng.println(a...)
			}
		}, "_scriptinclude": func(url string, a ...interface{}) (src interface{}, srcerr error) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				if lkpr, lkprerr := atvrntme.prsng.atv.LookupTemplate(url, a...); lkpr != nil && lkprerr == nil {
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

func newatvruntime(atv *Active, parsing *parsing) (atvrntme *atvruntime, err error) {
	atvrntme = &atvruntime{atv: atv, prsng: parsing, vm: goja.New(), includedpgrms: map[string]*goja.Program{}, intrnbuffs: map[*iorw.Buffer]*iorw.Buffer{}}
	atvrntme.atv.InterruptVM = func(v interface{}) {
		atvrntme.vm.Interrupt(v)
	}
	jsext.Register(atvrntme.vm)
	if definternmapref := defaultAtvRuntimeInternMap(atvrntme); definternmapref != nil && len(definternmapref) > 0 {
		if definternmapref != nil && len(definternmapref) > 0 {
			for k, ref := range definternmapref {
				atvrntme.vm.Set(k, ref)
			}
		}
	}
	if requirejsprgm != nil {
		if _, err = atvrntme.vm.RunProgram(requirejsprgm); err != nil {
			//atvrntme.includedpgrms[incllib] = requirejsprgm
		} else {
			//atvrntme.includedpgrms[incllib] = requirejsprgm
		}
	}
	return
}

var requirejsprgm *goja.Program = nil

//var GlobelModules map[string]*

func init() {
	var errpgrm error = nil
	if requirejsprgm, errpgrm = goja.Compile("", requirejs.RequireJSString(), false); errpgrm != nil {
		fmt.Println(errpgrm.Error())
	}
}

package active

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/evocert/kwe/babeljs"
	"github.com/evocert/kwe/ecma/es51"
	"github.com/evocert/kwe/ecma/es51/parser"
	"github.com/evocert/kwe/ecma/jsext"
	"github.com/evocert/kwe/requirejs"

	"github.com/evocert/kwe/iorw"
)

//Active - struct
type Active struct {
	Print          func(a ...interface{})
	Println        func(a ...interface{})
	FPrint         func(w io.Writer, a ...interface{})
	FPrintLn       func(w io.Writer, a ...interface{})
	LookupTemplate func(string, ...interface{}) (io.Reader, error)
	ObjectMapRef   func() map[string]interface{}
	lckprnt        *sync.Mutex
	InterruptVM    func(v interface{})
}

//NewActive - instance
func NewActive() (atv *Active) {
	atv = &Active{lckprnt: &sync.Mutex{}}
	return
}

func (atv *Active) print(w io.Writer, a ...interface{}) {
	if atv.Print != nil {
		if len(a) > 0 {
			atv.lckprnt.Lock()
			defer atv.lckprnt.Unlock()
			atv.Print(a...)
		}
	} else if atv.FPrint != nil && w != nil {
		if len(a) > 0 {
			atv.FPrint(w, a...)
		}
	} else if w != nil {
		if len(a) > 0 {
			iorw.Fprint(w, a...)
		}
	}
}

func (atv *Active) println(w io.Writer, a ...interface{}) {
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
		if len(a) > 0 {
			atv.lckprnt.Lock()
			defer atv.lckprnt.Unlock()
			fmt.Fprint(w, a...)
		}
		fmt.Fprintln(w)
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
	atv            *Active
	atvrntme       *atvruntime
	wout           io.Writer
	prslbli        []int
	prslblprv      []rune
	prntprsng      *parsing
	foundcde       bool
	hascde         bool
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
		prsng.atv.print(prsng.wout, a...)
	}
}

func (prsng *parsing) println(a ...interface{}) {
	if prsng.atv != nil {
		prsng.atv.println(prsng.wout, a...)
	}
}

func (prsng *parsing) close() {
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
		if prsng.atvrntme != nil {
			prsng.atvrntme.close()
			prsng.atvrntme = nil
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
			err = prsng.WriteRunes(p[:pl])
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
		err = prsng.WriteRunes(p[:pl])
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
			err = prsng.psvctrl.cachedbuf().WriteRunes(prsng.psvr[:pi])
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
				prsng.atvrntme = newatvruntime(prsng.atv, prsng)
				prsng.atvrntme.run()
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
		if prslbli[1] == 0 && prslbli[0] < len(prslbl[0]) {
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
				prslblprv[0] = pr
				err = parsepsvrune(prsng, pr)
			}
		} else if prslbli[0] == len(prslbl[0]) && prslbli[1] < len(prslbl[1]) {
			if prslbl[1][prslbli[1]] == pr {
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
				prslblprv[1] = pr
				err = parseatvrune(prsng, pr)
			}
		}
	}
	return
}

func nextparsing(atv *Active, prntprsng *parsing, wout io.Writer) (prsng *parsing) {
	prsng = &parsing{Buffer: iorw.NewBuffer(), wout: wout, prntprsng: prntprsng, atv: atv, prslbli: []int{0, 0}, prslblprv: []rune{0, 0}, cdeoffsetstart: -1, cdeoffsetend: -1, psvoffsetstart: -1, psvoffsetend: -1, psvr: make([]rune, 8192), cder: make([]rune, 8192)}
	return
}

type atvruntime struct {
	*parsing
	atv        *Active
	vm         *es51.Runtime
	intrnbuffs map[*iorw.Buffer]*iorw.Buffer
}

func (atvrntme *atvruntime) InvokeFunction(functocall interface{}, args ...interface{}) (result interface{}) {
	if functocall != nil {
		if atvrntme.vm != nil {
			var fnccallargs []es51.Value = nil
			var argsn = 0

			for argsn < len(args) {
				if fnccallargs == nil {
					fnccallargs = make([]es51.Value, len(args))
				}
				fnccallargs[argsn] = atvrntme.vm.ToValue(args[argsn])
				argsn++
			}
			if atvfunc, atvfuncok := functocall.(func(es51.FunctionCall) es51.Value); atvfuncok {
				var funccll = es51.FunctionCall{This: es51.Undefined(), Arguments: fnccallargs}
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
	val, err = atvrntme.corerun(atvrntme.code(), objmapref, map[string]interface{}{
		"newbuffer": func() (buff *iorw.Buffer) {
			buff = iorw.NewBuffer()
			buff.OnClose = atvrntme.removeBuffer
			atvrntme.intrnbuffs[buff] = buff
			return
		},
		"_passiveout": func(i int) {
			atvrntme.passiveout(i)
		},
		"parseEval": func(a ...interface{}) (val interface{}, err error) {
			return atvrntme.parseEval(a...)
		},
		"print": func(a ...interface{}) {
			if atvrntme.parsing != nil {
				atvrntme.parsing.print(a...)
			}
		},
		"println": func(a ...interface{}) {
			if atvrntme.parsing != nil {
				atvrntme.parsing.println(a...)
			}
		}, "scriptinclude": func(url string, a ...interface{}) (src interface{}, srcerr error) {
			if atvrntme.parsing != nil && atvrntme.parsing.atv != nil {
				if lkpr, lkprerr := atvrntme.parsing.atv.LookupTemplate(url, a...); lkpr != nil && lkprerr == nil {
					bufr := bufio.NewReader(lkpr)
					rnrs := make([]rune, 1024)
					inclsrc := ""
					rnrsi := 0
					for {
						rn, rns, rnerr := bufr.ReadRune()
						if rns > 0 {
							rnrs[rnrsi] = rn
							rnrsi++
							if rnrsi == len(rnrs) {
								inclsrc += string(rnrs[:])
								rnrsi = 0
							}
						}
						if rnerr != nil {
							break
						}
					}
					if rnrsi > 0 {
						inclsrc += string(rnrs[:rnrsi])
						rnrsi = 0
					}
					src = inclsrc
				} else if lkprerr != nil {
					srcerr = lkprerr
				}
			}
			if src == nil {
				src = ""
			}
			return
		},
		"script": atvrntme}, "require.js")
	return
}

func (atvrntme *atvruntime) corerun(code string, objmapref map[string]interface{}, internmapref map[string]interface{}, includelibs ...string) (val interface{}, err error) {
	if code != "" {
		atvrntme.vm.ClearInterrupt()
		jsext.Register(atvrntme.vm)
		atvrntme.atv.InterruptVM = func(v interface{}) {
			atvrntme.vm.Interrupt(v)
		}
		var glblobjstoremove []string = nil
		if objmapref != nil && len(objmapref) > 0 {
			if glblobjstoremove == nil {
				glblobjstoremove = []string{}
			}
			for k, ref := range objmapref {
				glblobjstoremove = append(glblobjstoremove, k)
				atvrntme.vm.Set(k, ref)
			}
		}
		if internmapref != nil && len(internmapref) > 0 {
			if glblobjstoremove == nil {
				glblobjstoremove = []string{}
			}
			for k, ref := range internmapref {
				glblobjstoremove = append(glblobjstoremove, k)
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
			if code, isrequired, err = transformCode(code, nil); err == nil {
				if isrequired {
					code = `function vmrequire(args) { return require([args]);}` + code
				}
				prsd, prsderr := parser.ParseFile(nil, "", code, 0)
				if prsderr != nil {
					err = prsderr
				}

				if err == nil {
					if len(includelibs) > 0 {
						for _, incllib := range includelibs {
							if incllib == "require.js" || incllib == "require.min.js" {
								if requirejsprgm != nil {
									if _, err = atvrntme.vm.RunProgram(requirejsprgm); err != nil {
										break
									}
								}
							} else if incllib == "babel.js" || incllib == "babel.min.js" {
								if babeljsprgm != nil {
									if _, err = atvrntme.vm.RunProgram(babeljsprgm); err != nil {
										break
									}
								}
							}
						}
					}
					if p, perr := es51.CompileAST(prsd, false); perr == nil {
						_, err = atvrntme.vm.RunProgram(p)
						if err != nil {
							fmt.Println(err.Error())
						}
						if len(glblobjstoremove) > 0 {
							if gbl := atvrntme.vm.GlobalObject(); gbl != nil {
								for _, k := range glblobjstoremove {
									gbl.Delete(k)
								}
								gbl = nil
							}
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
		}
	}
	return
}

var (
	defaultOpts = map[string]interface{}{
		"presets":       []interface{}{"es2015"},
		"ast":           true,
		"sourceMaps":    false,
		"babelrc":       false,
		"compact":       false,
		"retainLines":   true,
		"highlightCode": false,
	}
)

func transformCode(code string, opts map[string]interface{}) (trsnfrmdcde string, isrequired bool, err error) {
	vm := es51.New()
	isrequired = strings.IndexAny(code, "import ") > -1

	_, err = vm.RunProgram(babeljsprgm)

	if err != nil {
		err = fmt.Errorf("unable to load babel.js: %s", err)
	} else {
		var transform es51.Callable
		babel := vm.Get("Babel")
		if err := vm.ExportTo(babel.ToObject(vm).Get("transform"), &transform); err != nil {
			err = fmt.Errorf("unable to export transform fn: %s", err)
		} else {
			if opts == nil {
				opts = defaultOpts
			}
			if v, verr := transform(babel, vm.ToValue(code), vm.ToValue(opts)); verr != nil {
				err = fmt.Errorf("unable to export transform fn: %s", verr)
				fmt.Println(err.Error())
			} else {
				trsnfrmdcde = v.ToObject(vm).Get("code").String()
				if isrequired {
					trsnfrmdcde = strings.Replace(trsnfrmdcde, "require(\"", "vmrequire(\"", -1)
				}
			}
		}
	}
	//} else {
	//	trsnfrmdcde = code
	//}
	return
}

func (atvrntme *atvruntime) parseEval(a ...interface{}) (val interface{}, err error) {
	var prsng = atvrntme.parsing
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
			iorw.Fprint(pw, "<@")
			iorw.Fprint(pw, args...)
			iorw.Fprint(pw, "@>")
		}
	}(a...)
	rnr := bufio.NewReader(pr)
	wg.Wait()
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
		//if canexec {
		if prsng.foundCode() {
			if cdemapl := len(prsng.cdemap); cdemapl > orgcdemapl {
				cdecoords = []int64{prsng.cdemap[orgcdemapl][0], prsng.cdemap[cdemapl-1][1]}
			}
			val, err = atvrntme.corerun(atvrntme.code(cdecoords...), nil, nil)
		} else {
			if rdr := prsng.Reader(); rdr != nil {
				io.Copy(prsng.wout, rdr)
				rdr.Close()
				rdr = nil
			}
		}
	}
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
	if atvrntme != nil && atvrntme.parsing != nil {
		if cdel := len(atvrntme.parsing.cdemap); cdel > 0 {
			var cdei = 0
			var rdr = atvrntme.parsing.Reader()
			if len(coords) == 0 {
				coords = []int64{atvrntme.parsing.cdemap[cdei][0], atvrntme.parsing.cdemap[cdel-1][1]}
			}
			for cdei < cdel {
				cdecoors := atvrntme.parsing.cdemap[cdei]
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
					rdr.Seek(cdecoors[0], 0)
					var p = make([]byte, cdecoors[1]-cdecoors[0])
					rdr.Read(p)
					c += string(p)
				}
			}
		}
	}
	return c
}

func (atvrntme *atvruntime) passiveout(i int) {
	if atvrntme != nil && atvrntme.parsing != nil {
		if psvl := len(atvrntme.parsing.psvmap); psvl > 0 && i >= 0 && i < psvl {
			psvcoors := atvrntme.parsing.psvmap[i]
			if psvcoors[1] > psvcoors[0] {
				rdr := atvrntme.parsing.Reader()
				rdr.Seek(psvcoors[0], 0)
				io.CopyN(atvrntme.parsing.wout, rdr, psvcoors[1]-psvcoors[0])
			}
		}
	}
}

func (atvrntme *atvruntime) close() {
	if atvrntme != nil {
		if atvrntme.parsing != nil {
			atvrntme.parsing = nil
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

func newatvruntime(atv *Active, parsing *parsing) (atvrntme *atvruntime) {
	atvrntme = &atvruntime{atv: atv, parsing: parsing, vm: es51.New(), intrnbuffs: map[*iorw.Buffer]*iorw.Buffer{}}
	return
}

var requirejsprgm *es51.Program = nil
var babeljsprgm *es51.Program = nil

//var GlobelModules map[string]*

func init() {
	var errpgrm error = nil
	if requirejsprgm, errpgrm = es51.Compile("", requirejs.RequireMinJSString(), false); errpgrm != nil {
		fmt.Println(errpgrm.Error())
	}
	if babeljsprgm, errpgrm = es51.Compile("", babeljs.BabelJSString(), true); errpgrm != nil {
		fmt.Println(errpgrm.Error())
	}
}

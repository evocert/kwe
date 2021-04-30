package active

import (
	"bufio"
	"fmt"
	"io"
	"runtime"
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
		if len(imprtglbs) > 0 {
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

func activeFinalize(atv *Active) {
	if atv != nil {
		atv.dispose()
		atv = nil
	}
}

//NewActive - instance
func NewActive(namespace ...string) (atv *Active) {
	atv = &Active{lckprnt: &sync.Mutex{}, Namespace: "", atvruntime: nil}
	atv.atvruntime, _ = newatvruntime(atv, nil)
	if len(namespace) == 1 && namespace[0] != "" {
		atv.Namespace = namespace[0] + "."
	}
	runtime.SetFinalizer(atv, activeFinalize)
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
			if atv.prsng == nil || atv.prsng != prsng {
				if atv.prsng != nil {
					atv.prsng.dispose()
					atv.prsng = nil
				}
				atv.prsng = prsng
			}
		}
		if atv.atvruntime != nil {
			atv.atvruntime.run()
		}
	}
}

//Eval - parse rin io.Reader, execute if neaded and output to wou io.Writer
func (atv *Active) Eval(wout io.Writer, rin io.Reader, evalstngs ...interface{}) {
	lck := &sync.RWMutex{}
	lck.RLock()
	defer lck.RUnlock()
	var parsing = nextparsing(atv, nil, wout, evalstngs...)
	defer parsing.Close()
	var rnr io.RuneReader = nil
	var bfr *bufio.Reader = nil
	if rr, rrok := rin.(io.RuneReader); rrok {
		rnr = rr
	} else {
		bfr = bufio.NewReader(rin)
		rnr = bfr
	}
	parseprsngrunerdr(parsing, rnr, true)
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
var elmlbl = [][]rune{[]rune("<#"), []rune(">"), []rune("</#"), []rune(">"), []rune("<#"), []rune("/>")}
var phrslbl = [][]rune{[]rune("{#"), []rune("#}")}

type parsing struct {
	*iorw.Buffer
	tmpltbuf       *iorw.Buffer
	tmpltmap       map[string][]int64
	atv            *Active
	wout           io.Writer
	prntrs         []io.Writer
	prslbli        []int
	prslblprv      []rune
	prntprsng      *parsing
	foundcde       bool
	hascde         bool
	cdetxt         rune
	cdebuf         *iorw.Buffer
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
	//psvsection
	tmpbuf    *iorw.Buffer
	elmlbli   []int
	elmoffset int
	elmprvrns []rune
	//elmType     elemtype
	prvelmrn    rune
	crntpsvsctn *psvsection
	headpsvsctn *psvsection
	tailpsvsctn *psvsection
	//psvctrl        *passivectrl
	//prvpsvctrls    map[*passivectrl]*passivectrl
	prsvpth string
}

func (prsng *parsing) tempbuf() *iorw.Buffer {
	if prsng.tmpbuf == nil {
		prsng.tmpbuf = iorw.NewBuffer()
	}
	return prsng.tmpbuf
}

func (prsng *parsing) cdeBuff() *iorw.Buffer {
	if prsng.cdebuf == nil {
		prsng.cdebuf = iorw.NewBuffer()
	}
	return prsng.cdebuf
}

func (prsng *parsing) tmpltBuf() *iorw.Buffer {
	if prsng != nil {
		if prsng.tmpltbuf == nil {
			prsng.tmpltbuf = iorw.NewBuffer()
		}
		return prsng.tmpltbuf
	}
	return nil
}

func (prsng *parsing) tmpltMap() map[string][]int64 {
	if prsng != nil {
		if prsng.tmpltmap == nil {
			prsng.tmpltmap = map[string][]int64{}
		}
		return prsng.tmpltmap
	}
	return nil
}

func (prsng *parsing) tmpltrdr(tmpltnme string) (rdr *iorw.BuffReader, mxlen int64) {
	mxlen = -1
	if prsng != nil && prsng.tmpltmap != nil && prsng.tmpltbuf != nil {
		if strtend, strtendok := prsng.tmpltmap[tmpltnme]; strtendok {
			if s := prsng.tmpltbuf.Size(); s > 0 && len(strtend) > 0 && strtend[0] >= 0 && strtend[1] <= s {
				if mxlen = (strtend[1] - strtend[0]); mxlen > 0 {
					rdr = prsng.tmpltbuf.Reader()
					rdr.Seek(strtend[0], io.SeekStart)
					rdr.MaxRead = mxlen
				}
			}
		}
	}
	return
}

func (prsng *parsing) print(a ...interface{}) {
	if prsng.atv != nil {
		if pl := len(prsng.prntrs); pl > 0 {
			prsng.atv.print(prsng.prntrs[pl-1], a...)
		} else {
			prsng.atv.print(prsng.wout, a...)
		}
	}
}

func (prsng *parsing) println(a ...interface{}) {
	if prsng.atv != nil {
		if pl := len(prsng.prntrs); pl > 0 {
			prsng.atv.println(prsng.prntrs[pl-1], a...)
		} else {
			prsng.atv.println(prsng.wout, a...)
		}
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
		if prsng.cdebuf != nil {
			prsng.cdebuf.Close()
			prsng.cdebuf = nil
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
		if prsng.tmpltmap != nil {
			for k := range prsng.tmpltmap {
				prsng.tmpltmap[k] = nil
				delete(prsng.tmpltmap, k)
			}
			prsng.tmpltmap = nil
		}
		if prsng.tmpltbuf != nil {
			prsng.tmpltbuf.Close()
			prsng.tmpltbuf = nil
		}
		if prsng.tmpbuf != nil {
			prsng.tmpbuf.Close()
			prsng.tmpbuf = nil
		}
		for prsng.tailpsvsctn != nil {
			prsng.tailpsvsctn.dispose()
		}
		if prsng.headpsvsctn != nil {
			prsng.headpsvsctn.dispose()
		}
		if prsng.crntpsvsctn != nil {
			prsng.crntpsvsctn = nil
		}
		if prsng.elmlbli != nil {
			prsng.elmlbli = nil
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

func (prsng *parsing) writePsv(p ...rune) (err error) {
	if pl := len(p); pl > 0 {
		if prsng.crntpsvsctn == nil {
			if prsng.foundCode() {
				if prsng.psvoffsetstart == -1 {
					prsng.psvoffsetstart = prsng.Size()
				}
				err = prsng.WriteRunes(p[:pl]...)
			} else {
				if bs := iorw.RunesToUTF8(p[:pl]); len(bs) > 0 {
					_, err = prsng.wout.Write(bs)
				}
			}
		} else {
			if prsng.crntpsvsctn.canphrs {
				for _, phrsr := range p[:pl] {
					if err = parsepsvphrase(prsng, prsng.crntpsvsctn, prsng.crntpsvsctn.phrslbli, phrsr); err != nil {
						return
					}
				}
			} else {
				err = prsng.crntpsvsctn.CachedBuf().WriteRunes(p[:pl]...)
			}
		}
	}
	return
}

func (prsng *parsing) writeCde(p []rune) (err error) {
	if pl := len(p); pl > 0 {
		if prsng.cdeoffsetstart == -1 {
			prsng.cdeoffsetstart = prsng.cdeBuff().Size()
		}
		err = prsng.cdeBuff().WriteRunes(p[:pl]...)
	}
	return
}

func (prsng *parsing) foundCode() bool {
	return prsng.foundcde
}

func (prsng *parsing) flushPsv() (err error) {
	if pi := prsng.psvri; pi > 0 {
		prsng.psvri = 0
		err = prsng.writePsv(prsng.psvr[:pi]...)
	}
	if prsng.crntpsvsctn == nil && prsng.foundCode() {
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
		prsng.setcdepos(cdeoffsetstart, prsng.cdeBuff().Size())
	}
	return
}

func parseprsngrunerdr(prsng *parsing, rnr io.RuneReader, canexec bool) (err error) {
	for err == nil {
		r, rsize, rerr := rnr.ReadRune()
		if rsize > 0 {
			if err = parseprsng(prsng, prsng.prslbli, prsng.prslblprv, r); err != nil {
				break
			}
		}
		if rerr != nil {
			err = rerr
		}
	}
	if err == io.EOF || err == nil {
		prsng.flushPsv()
		if canexec {
			prsng.flushCde()
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
				if prsng.crntpsvsctn != nil {
					if err = prsng.crntpsvsctn.CachedBuf().WriteRunes(prslbl[0]...); err != nil {
						return
					}
				}
				prslblprv[0] = 0
			} else {
				prslblprv[0] = pr
			}
		} else {
			if psvl := prslbli[0]; psvl > 0 {
				prslbli[0] = 0
				prslblprv[0] = 0
				if err = parsepsvrunes(prsng, prslbl[0][0:psvl]); err != nil {
					return
				}
			}
			err = parsepsvrune(prsng, pr)
			prslblprv[0] = pr
		}
	} else if prslbli[0] == len(prslbl[0]) && prslbli[1] < len(prslbl[1]) {
		if prsng.cdetxt == rune(0) && prslbl[1][prslbli[1]] == pr {
			prslbli[1]++
			if prslbli[1] == len(prslbl[1]) {
				if prsng.crntpsvsctn != nil {
					if err = prsng.crntpsvsctn.CachedBuf().WriteRunes(prslbl[1]...); err != nil {
						return
					}
				}
				prslbli[0] = 0
				prslblprv[1] = 0
				prslbli[1] = 0
			} else {
				prslblprv[1] = pr
			}
		} else {
			if prsl := prslbli[1]; prsl > 0 {
				prslbli[1] = 0
				if prsng.crntpsvsctn != nil {
					if err = prsng.crntpsvsctn.CachedBuf().WriteRunes(prslbl[1][:prsl]...); err != nil {
						return
					}
				} else {
					if err = parseatvrunes(prsng, prslbl[1][:prsl]); err != nil {
						return
					}
				}
			}
			if prsng.crntpsvsctn != nil {
				err = prsng.crntpsvsctn.CachedBuf().WriteRune(pr)
			} else {
				err = parseatvrune(prsng, pr)
			}
			prslblprv[1] = pr
		}
	}
	//}
	return
}

func parsingFinalize(prsng *parsing) {
	if prsng != nil {
		prsng.dispose()
		prsng = nil
	}
}

func nextparsing(atv *Active, prntprsng *parsing, wout io.Writer, prsrstngs ...interface{}) (prsng *parsing) {
	prsng = &parsing{Buffer: iorw.NewBuffer(), wout: wout, prntprsng: prntprsng, atv: atv, cdetxt: rune(0), prslbli: []int{0, 0}, prslblprv: []rune{0, 0}, cdeoffsetstart: -1, cdeoffsetend: -1, psvoffsetstart: -1, psvoffsetend: -1, psvr: make([]rune, 8192), cder: make([]rune, 8192), prntrs: []io.Writer{},
		crntpsvsctn: nil, prvelmrn: rune(0), elmoffset: -1, elmlbli: []int{0, 0}, elmprvrns: []rune{rune(0), rune(0)}}
	if len(prsrstngs) > 0 {
		for _, d := range prsrstngs {
			if mp, mpok := d.(map[string]interface{}); mpok && len(mp) > 0 {
				for k, v := range mp {
					if k == "init-path" {
						if s, _ := v.(string); s != "" {
							if prsng.prsvpth == "" {
								prsng.prsvpth = s
							}
						}
					}
				}
			} else if s, sok := d.(string); sok && s != "" {
				if prsng.prsvpth == "" {
					prsng.prsvpth = s
				}
			}
		}
	}
	//runtime.SetFinalizer(prsng, parsingFinalize)
	return
}

type atvruntime struct {
	prsng         *parsing
	atv           *Active
	vm            *goja.Runtime
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
	val, err = atvrntme.corerun(atvrntme.code(), objmapref)
	return
}

func (atvrntme *atvruntime) corerun(code string, objmapref map[string]interface{}, includelibs ...string) (val interface{}, err error) {
	if code != "" {
		atvrntme.vm.ClearInterrupt()
		if len(objmapref) > 0 {
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
			if errs := err.Error(); errs != "" && !strings.HasPrefix(errs, "exit at <eval>:") {
				fmt.Println(err.Error())
				fmt.Println(code)
				err = nil
			}

		}
	} else {
		if err != nil {
			if errs := err.Error(); errs != "" && !strings.HasPrefix(errs, "exit at <eval>:") {
				fmt.Println(err.Error())
				fmt.Println(code)
				err = nil
			}
		}
	}
	return
}

func transformCode(code string, namespace string, opts map[string]interface{}) (trsnfrmdcde string, isrequired bool, err error) {
	trsnfrmdcde = code
	isrequired = strings.Contains(code, "require(\"")
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
	} else {
		prsng.prslbli[0] = 0
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
			if crunesi > 0 {
				cl := crunesi
				crunesi = 0
				cdetestbuf.Print(string(crunes[:cl]))
				for _, cr := range crunes[:cl] {
					parseprsng(prsng, prsng.prslbli, prsng.prslblprv, cr)
				}
			}
			prsng.flushPsv()
			if canprcss {
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
			var rdr *iorw.BuffReader = nil
			if len(coords) == 0 {
				coords = []int64{atvrntme.prsng.cdemap[cdei][0], atvrntme.prsng.cdemap[cdel-1][1]}
			}
			var mxdcde int64 = int64(0)
			if len(coords) == 2 && coords[0] <= coords[1] {
				mxdcde = coords[1] - coords[0]
			}
			if mxdcde > 0 {
				for cdei < cdel && mxdcde > 0 {
					if cdecrds, cdecrdsok := atvrntme.prsng.cdemap[cdei]; cdecrdsok && (cdecrds[0] <= coords[1] && cdecrds[1] >= coords[0]) {
						if strti := cdecrds[0]; strti >= 0 {
							if strti < coords[0] {
								strti = coords[0]
							}
							if endi := cdecrds[1]; endi >= 0 {
								if endi > coords[1] {
									endi = coords[1]
								}
								if endi-strti > 0 {
									if rdr == nil {
										rdr = atvrntme.prsng.cdeBuff().Reader()
									}
									rdr.Seek(strti, 0)
									rdr.MaxRead = endi - strti
									if rs, rserr := iorw.ReaderToString(rdr); rs != "" {
										c += rs
										mxdcde -= (endi - strti)
									} else if rserr != nil {
										break
									}
								}
							}
						}
					}
					cdei++
				}

			}
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
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				atvrntme.prsng.print(a...)
			} else if atvrntme.atv != nil {
				atvrntme.atv.print(nil, a...)
			}
		},
		atvrntme.atv.namespace() + "println": func(a ...interface{}) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				atvrntme.prsng.println(a...)
			} else if atvrntme.atv != nil {
				atvrntme.atv.println(nil, a...)
			}
		},
		atvrntme.atv.namespace() + "println": func(a ...interface{}) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				atvrntme.prsng.println(a...)
			} else if atvrntme.atv != nil {
				atvrntme.atv.println(nil, a...)
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

func atvruntimeFinalize(atvrntme *atvruntime) {
	if atvrntme != nil {
		atvrntme.dispose()
		atvrntme = nil
	}
}

func newatvruntime(atv *Active, parsing *parsing) (atvrntme *atvruntime, err error) {
	atvrntme = &atvruntime{atv: atv, prsng: parsing, vm: goja.New(), includedpgrms: map[string]*goja.Program{}, intrnbuffs: map[*iorw.Buffer]*iorw.Buffer{}}
	atvrntme.atv.InterruptVM = func(v interface{}) {
		atvrntme.vm.Interrupt(v)
	}
	jsext.Register(atvrntme.vm)
	if definternmapref := defaultAtvRuntimeInternMap(atvrntme); len(definternmapref) > 0 {
		if len(definternmapref) > 0 {
			for k, ref := range definternmapref {
				atvrntme.vm.Set(k, ref)
			}
		}
	}
	if requirejsprgm != nil {
		_, err = atvrntme.vm.RunProgram(requirejsprgm)
	}
	//runtime.SetFinalizer(atvrntme, atvruntimeFinalize)
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

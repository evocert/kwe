package active

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	//"github.com/evocert/kwe/ecma/es6"
	"github.com/dop251/goja"
	"github.com/evocert/kwe/iorw"
)

//Active - struct
type Active struct {
	Print    func(a ...interface{})
	Println  func(a ...interface{})
	FPrint   func(w io.Writer, a ...interface{})
	FPrintLn func(w io.Writer, a ...interface{})
	lckprnt  *sync.Mutex
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
			atv.lckprnt.Lock()
			defer atv.lckprnt.Unlock()
			atv.FPrint(w, a...)
		}
	} else if w != nil {
		if len(a) > 0 {
			atv.lckprnt.Lock()
			defer atv.lckprnt.Unlock()
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
	} else if atv.FPrint != nil && w != nil {
		if len(a) > 0 {
			atv.lckprnt.Lock()
			defer atv.lckprnt.Unlock()
			atv.FPrint(w, a...)
		}
	} else if w != nil {
		if len(a) > 0 {
			atv.lckprnt.Lock()
			defer atv.lckprnt.Unlock()
			fmt.Fprint(w, a...)
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
	parseprsngrunerdr(parsing, rnr)

	parsing.Close()
}

//Close - refer to  io.Closer
func (atv *Active) Close() (err error) {
	if atv.lckprnt != nil {
		atv.lckprnt = nil
	}
	return
}

type parsing struct {
	*iorw.Buffer
	atv            *Active
	atvrntme       *atvruntime
	wout           io.Writer
	prslbl         [][]rune
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
		if prsng.prslbl != nil {
			prsng.prslbl = nil
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
		prsng.writePsv(prsng.psvr[:pi])
	}
	if prsng.foundCode() {
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

func parsepsvrune(prsng *parsing, rn rune) (err error) {
	prsng.flushCde()
	if prsng.hascde {
		prsng.hascde = false
	}
	prsng.psvr[prsng.psvri] = rn
	prsng.psvri++
	if prsng.psvri == len(prsng.psvr) {
		prsng.psvri = 0
		err = prsng.writePsv(prsng.psvr)
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

func parseatvrune(prsng *parsing, rn rune) (err error) {
	if !prsng.hascde {
		prsng.flushPsv()
		if !prsng.foundcde {
			prsng.foundcde = true
		}
		if strings.TrimSpace(string(rn)) != "" {
			prsng.hascde = true
		}
	}
	if prsng.hascde {
		prsng.cder[prsng.cderi] = rn
		prsng.cderi++
		if prsng.cderi == len(prsng.cder) {
			prsng.cderi = 0
			err = prsng.writeCde(prsng.cder)
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

func parseprsngrunerdr(prsng *parsing, rnr io.RuneReader) (err error) {
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
					parseprsng(prsng, prsng.prslbl, prsng.prslbli, prsng.prslblprv, cr)
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
				parseprsng(prsng, prsng.prslbl, prsng.prslbli, prsng.prslblprv, cr)
			}
		}
		prsng.flushPsv()
		prsng.flushCde()
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
	return
}

func parseprsng(prsng *parsing, prslbl [][]rune, prslbli []int, prslblprv []rune, pr rune) {
	if prslbli[1] == 0 && prslbli[0] < len(prslbl[0]) {
		if prslbli[0] > 0 && prslbl[0][prslbli[0]-1] == prslblprv[0] && prslbl[0][prslbli[0]] != pr {
			if psvl := prslbli[0]; psvl > 0 {
				prslbli[0] = 0
				prslblprv[0] = 0
				parsepsvrunes(prsng, prslbl[0][0:psvl])
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
				parsepsvrunes(prsng, prslbl[0][0:psvl])
			}
			prslblprv[0] = pr
			parsepsvrune(prsng, pr)
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
				parseatvrunes(prsng, prslbl[1][:prsl])
			}
			prslblprv[1] = pr
			parseatvrune(prsng, pr)
		}
	}
}

func nextparsing(atv *Active, prntprsng *parsing, wout io.Writer) (prsng *parsing) {
	prsng = &parsing{Buffer: iorw.NewBuffer(), wout: wout, prntprsng: prntprsng, atv: atv, prslbl: [][]rune{[]rune("<@"), []rune("@>")}, prslbli: []int{0, 0}, prslblprv: []rune{0, 0}, cdeoffsetstart: -1, cdeoffsetend: -1, psvoffsetstart: -1, psvoffsetend: -1, psvr: make([]rune, 8192), cder: make([]rune, 8192)}
	return
}

type atvruntime struct {
	*parsing
	atv *Active
	vm  *goja.Runtime
}

func (atvrntme *atvruntime) run() (err error) {
	if cde := atvrntme.code(); cde != "" {
		atvrntme.vm.Set("_passiveout", func(i int) {
			atvrntme.passiveout(i)
		})
		atvrntme.vm.Set("print", func(a ...interface{}) {
			if atvrntme.parsing != nil {
				atvrntme.parsing.print(a...)
			}
		})
		atvrntme.vm.Set("println", func(a ...interface{}) {
			if atvrntme.parsing != nil {
				atvrntme.parsing.println(a...)
			}
		})
		p, err := goja.Compile("", cde, false)
		if err == nil {
			_, err = atvrntme.vm.RunProgram(p)
		}
		if gbl := atvrntme.vm.GlobalObject(); gbl != nil {
			if ks := gbl.Keys(); len(ks) > 0 {
				for _, k := range ks {
					gbl.Delete(k)
				}
				ks = nil
			}
			gbl = nil
		}
	}
	return
}

func (atvrntme *atvruntime) code() (c string) {
	if atvrntme != nil && atvrntme.parsing != nil {
		if cdel := len(atvrntme.parsing.cdemap); cdel > 0 {
			var cdei = 0
			var rdr = atvrntme.parsing.Reader()
			for cdei < cdel {
				cdecoors := atvrntme.parsing.cdemap[cdei]
				cdei++
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
					/*for _, k := range ks {
						atvrntme.vm.GlobalObject().Delete(k)
					}*/
				}
			}
			atvrntme.vm = nil
		}
	}
}

func newatvruntime(atv *Active, parsing *parsing) (atvrntme *atvruntime) {
	atvrntme = &atvruntime{atv: atv, parsing: parsing, vm: goja.New()}
	return
}

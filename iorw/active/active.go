package active

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja/parser"

	"github.com/dop251/goja"

	"github.com/evocert/kwe/ecma/jsext"
	"github.com/evocert/kwe/enumeration"
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
	Readln         func() (string, error)
	ReadLines      func() ([]string, error)
	ReadAll        func() (string, error)
	FReadln        func(r io.Reader) (string, error)
	FReadLines     func(r io.Reader) ([]string, error)
	FReadAll       func(r io.Reader) (string, error)
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
				func() {
					atv.lckprnt.Lock()
					defer atv.lckprnt.Unlock()
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

func (atv *Active) readln(r io.Reader) (ln string, err error) {
	if rdr, rdrok := r.(iorw.Reader); rdrok {
		ln, err = rdr.Readln()
	} else {
		if atv.Readln != nil {
			ln, err = atv.Readln()
		} else if atv.FReadln != nil && r != nil {
			func() {
				atv.lckprnt.Lock()
				defer atv.lckprnt.Unlock()
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
				atv.lckprnt.Lock()
				defer atv.lckprnt.Unlock()
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
				atv.lckprnt.Lock()
				defer atv.lckprnt.Unlock()
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
		callback(atv.vm())
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

func (atv *Active) atvrun(prsng *parsing) (err error) {
	if prsng != nil {
		if atv.atvruntime == nil {
			atv.atvruntime, err = newatvruntime(atv, prsng)
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
			_, err = atv.atvruntime.run()
		}
	}
	return
}

//Eval - parse a ...interface{} arguments, execute if neaded and output to wou io.Writer
func (atv *Active) Eval(wout io.Writer, rin io.Reader, initpath string, a ...interface{}) (err error) {
	var parsing = nextparsing(atv, nil, rin, wout, initpath)
	defer parsing.dispose()
	err = parseprsng(parsing, true, a...)
	return
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
	rnrdrsbeingparsed    *enumeration.List
	crntrnrdrbeingparsed io.RuneReader
	*iorw.Buffer
	tmpltbuf       *iorw.Buffer
	tmpltmap       map[string][]int64
	atv            *Active
	wout           io.Writer
	rin            io.Reader
	prntrs         []io.Writer
	rdrs           []io.Reader
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

func (prsng *parsing) readLn() (ln string, err error) {
	if prsng.atv != nil {
		if pl := len(prsng.rdrs); pl > 0 {
			ln, err = prsng.atv.readln(prsng.rdrs[pl-1])
		} else {
			ln, err = prsng.atv.readln(prsng.rin)
		}
	}
	return
}

func (prsng *parsing) readLines() (lines []string, err error) {
	if prsng.atv != nil {
		if pl := len(prsng.rdrs); pl > 0 {
			lines, err = prsng.atv.readlines(prsng.rdrs[pl-1])
		} else {
			lines, err = prsng.atv.readlines(prsng.rin)
		}
	}
	return
}

func (prsng *parsing) readAll() (s string, err error) {
	if prsng.atv != nil {
		if pl := len(prsng.rdrs); pl > 0 {
			s, err = prsng.atv.readAll(prsng.rdrs[pl-1])
		} else {
			s, err = prsng.atv.readAll(prsng.rin)
		}
	}
	return
}

func (prsng *parsing) incread(r io.Reader) {
	if prsng != nil {
		prsng.rdrs = append(prsng.rdrs, r)
	}
}

func (prsng *parsing) resetread() {
	if prsng.rdrs != nil {
		for len(prsng.rdrs) > 0 {
			prsng.rdrs[len(prsng.rdrs)-1] = nil
			prsng.rdrs = prsng.rdrs[:len(prsng.rdrs)-1]
		}
	}
}

func (prsng *parsing) decread() {
	if prsng.rdrs != nil {
		if len(prsng.rdrs) > 0 {
			prsng.rdrs[len(prsng.rdrs)-1] = nil
			prsng.rdrs = prsng.rdrs[:len(prsng.rdrs)-1]
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
		if prsng.crntrnrdrbeingparsed != nil {
			prsng.crntrnrdrbeingparsed = nil
		}
		if prsng.rnrdrsbeingparsed != nil {
			prsng.rnrdrsbeingparsed.Dispose(
				func(n *enumeration.Node, i interface{}) {

				},
				func(n *enumeration.Node, i interface{}) {
					if rc, _ := i.(io.Closer); rc != nil {
						rc.Close()
					}
				})
			prsng.rnrdrsbeingparsed = nil
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
			err = parseatvrunes(prsng, []rune(fmt.Sprintf("_psvout(%d);", pos)))
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

func loadRuneReaders(prsng *parsing, a ...interface{}) {
	if al := len(a); al > 0 {
		var tmpa []interface{}

		tampardr := func() {
			if len(tmpa) > 0 {
				pr, pw := io.Pipe()
				ctx, ctxcancel := context.WithCancel(context.Background())
				go func(ga ...interface{}) {
					defer pw.Close()
					ctxcancel()
					iorw.Fprint(pw, ga...)
				}(tmpa...)
				tmpa = nil
				<-ctx.Done()
				ctx = nil
				prsng.rnrdrsbeingparsed.Push(nil, nil, iorw.NewEOFCloseSeekReader(pr))
			}
			return
		}
		for ai := 0; ai < al; ai++ {
			if d := a[0]; d != nil {
				ai++
				a = a[1:]
				if r, rok := d.(io.Reader); rok {
					tampardr()
					if rnr, rnrok := r.(io.RuneReader); rnrok {
						prsng.rnrdrsbeingparsed.Push(nil, nil, rnr)
					} else {
						prsng.rnrdrsbeingparsed.Push(nil, nil, iorw.NewEOFCloseSeekReader(r, false))
					}
				} else if buf, bufok := d.(*iorw.Buffer); bufok {
					tampardr()
					prsng.rnrdrsbeingparsed.Push(nil, nil, buf.Reader())
				} else {
					if tmpa == nil {
						tmpa = []interface{}{}
					}
					tmpa = append(tmpa, d)
				}
			} else {
				ai++
				a = a[1:]
			}
		}
		tampardr()
	}
}

func parseprsng(prsng *parsing, canexec bool, a ...interface{}) (err error) {
	loadRuneReaders(prsng, a...)
	for err == nil {
		r, rsize, rerr := prsng.ReadRune()
		if rsize > 0 {
			if err = parseprsngrune(prsng, prsng.prslbli, prsng.prslblprv, r); err != nil {
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
				err = prsng.atv.atvrun(prsng)
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

func parseprsngrune(prsng *parsing, prslbli []int, prslblprv []rune, pr rune) (err error) {
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
	return
}

func (prsng *parsing) ReadRune() (r rune, size int, err error) {
	if prsng != nil {
		if prsng.crntrnrdrbeingparsed == nil {
			if prsng.rnrdrsbeingparsed != nil && prsng.rnrdrsbeingparsed.Length() > 0 {
				prsng.rnrdrsbeingparsed.Tail().Dispose(nil, func(nde *enumeration.Node, val interface{}) {
					prsng.crntrnrdrbeingparsed, _ = val.(io.RuneReader)
				})
			}
			if prsng.crntrnrdrbeingparsed == nil {
				err = io.EOF
				return
			}
		}
		r, size, err = prsng.crntrnrdrbeingparsed.ReadRune()
		if err == io.EOF {
			prsng.crntrnrdrbeingparsed = nil
		}
	} else {
		err = io.EOF
	}
	return
}

func nextparsing(atv *Active, prntprsng *parsing, rin io.Reader, wout io.Writer, initpath string) (prsng *parsing) {
	prsng = &parsing{Buffer: iorw.NewBuffer(), prsvpth: initpath, rnrdrsbeingparsed: enumeration.NewList(), crntrnrdrbeingparsed: nil, rin: rin, wout: wout, prntprsng: prntprsng, atv: atv, cdetxt: rune(0), prslbli: []int{0, 0}, prslblprv: []rune{0, 0}, cdeoffsetstart: -1, cdeoffsetend: -1, psvoffsetstart: -1, psvoffsetend: -1, psvr: make([]rune, 8192), cder: make([]rune, 8192), prntrs: []io.Writer{},
		crntpsvsctn: nil, prvelmrn: rune(0), elmoffset: -1, elmlbli: []int{0, 0}, elmprvrns: []rune{rune(0), rune(0)}}

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
						/*if err != nil {
							fmt.Println(err.Error())
						}*/
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
			excpath = atvrntme.prsng.prsvpth
		}
		for cdn, cd := range strings.Split(code, "\n") {
			cde += fmt.Sprintf("%d:%s\r\n", (cdn + 1), strings.TrimSpace(cd))
		}
		cdeerr := codeException(cde, excpath, err)
		err = cdeerr
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
	ctx, ctxcancel := context.WithCancel(context.Background())
	go func(args ...interface{}) {
		defer func() {
			pw.Close()
		}()
		ctxcancel()
		if len(args) > 0 {
			iorw.Fprint(pw, args...)
		}
	}(a...)
	rnr := bufio.NewReader(pr)
	<-ctx.Done()
	ctx = nil
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
						parseprsngrune(prsng, prsng.prslbli, prsng.prslblprv, cr)
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
					parseprsngrune(prsng, prsng.prslbli, prsng.prslblprv, cr)
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
		atvrntme.atv.namespace() + "parseEval": func(a ...interface{}) (val interface{}, err error) {
			return atvrntme.parseEval(false, a...)
		},
		//WRITER
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
		//READER
		atvrntme.atv.namespace() + "incread": func(r io.Reader) {
			if atvrntme.prsng != nil {
				atvrntme.prsng.incread(r)
			}
		},
		atvrntme.atv.namespace() + "resetread": func() {
			if atvrntme.prsng != nil {
				atvrntme.prsng.resetread()
			}
		},
		atvrntme.atv.namespace() + "decread": func() {
			if atvrntme.prsng != nil {
				atvrntme.prsng.decread()
			}
		},
		atvrntme.atv.namespace() + "print": func(a ...interface{}) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				atvrntme.prsng.print(a...)
			} else if atvrntme.atv != nil {
				atvrntme.atv.print(nil, a...)
			}
		},
		atvrntme.atv.namespace() + "readln": func() (ln string, err error) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				ln, err = atvrntme.prsng.readLn()
			} else if atvrntme.atv != nil {
				ln, err = atvrntme.atv.readln(nil)
			}
			return
		},
		atvrntme.atv.namespace() + "readLines": func() (lines []string, err error) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				lines, err = atvrntme.prsng.readLines()
			} else if atvrntme.atv != nil {
				lines, err = atvrntme.atv.readlines(nil)
			}
			return
		}, atvrntme.atv.namespace() + "readAll": func() (s string, err error) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				s, err = atvrntme.prsng.readAll()
			} else if atvrntme.atv != nil {
				s, err = atvrntme.atv.readAll(nil)
			}
			return
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

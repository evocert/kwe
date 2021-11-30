package active

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/dop251/goja/parser"

	"github.com/dop251/goja"

	"github.com/evocert/kwe/ecma/jsext"
	"github.com/evocert/kwe/requirejs"

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
	//lckprnt        *sync.Mutex
	InterruptVM func(v interface{})
	*atvruntime
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

func (atv *Active) atvrun(prsng *parsing) (err error) {
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
	var parsing = nextparsing(atv, nil, rin, wout, initpath)
	defer parsing.dispose()
	if len(a) > 0 {
		if invertactpsv {
			a = append(append([]interface{}{"<@"}, a...), "@>")
		}
	}
	err = parseprsng(parsing, true, a...)
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
		atv.atvruntime.dispose()
		atv.atvruntime = nil
	}
	return
}

//Clear
func (atv *Active) Clear() (err error) {
	if atv.atvruntime != nil {
		atv.atvruntime.dispose(true)
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
	woutbytes      []byte
	woutbytesi     int
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

func (prsng *parsing) binwrite(b ...byte) (n int, err error) {
	if prsng.atv != nil {
		if pl := len(prsng.prntrs); pl > 0 {
			n, err = prsng.atv.binwrite(prsng.prntrs[pl-1], b...)
		} else {
			n, err = prsng.atv.binwrite(prsng.wout, b...)
		}
	}
	return
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

func (prsng *parsing) seek(offset int64, whence int) (n int64, err error) {
	if prsng.atv != nil {
		if pl := len(prsng.rdrs); pl > 0 {
			n, err = prsng.atv.seek(prsng.rdrs[pl-1], offset, whence)
		} else {
			n, err = prsng.atv.seek(prsng.rin, offset, whence)
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

func (prsng *parsing) binread(size int) (b []byte, err error) {
	if prsng.atv != nil {
		if pl := len(prsng.rdrs); pl > 0 {
			b, err = prsng.atv.binread(prsng.rdrs[pl-1], size)
		} else {
			b, err = prsng.atv.binread(prsng.rin, size)
		}
	}
	return
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
		if prsng.woutbytes != nil {
			prsng.woutbytes = nil
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

func (prsng *parsing) flushWritePsv() (err error) {
	if prsng != nil && prsng.woutbytesi > 0 {
		_, err = prsng.wout.Write(prsng.woutbytes[0:prsng.woutbytesi])
		prsng.woutbytesi = 0
	}
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
					for _, bsb := range bs {
						prsng.woutbytes[prsng.woutbytesi] = bsb
						prsng.woutbytesi++
						if prsng.woutbytesi == len(prsng.woutbytes) {
							if err = prsng.flushWritePsv(); err != nil {
								break
							}
						}
					}
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
	if err == nil {
		err = prsng.flushWritePsv()
	}
	if err == nil && prsng.crntpsvsctn == nil && prsng.foundCode() {
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

func parseprsng(prsng *parsing, canexec bool, a ...interface{}) (err error) {
	var rnr = iorw.NewMultiArgsReader(a...)
	func() {
		defer rnr.Close()
		for err == nil {
			r, rsize, rerr := rnr.ReadRune()
			if rsize > 0 {
				if err = parseprsngrune(prsng, prsng.prslbli, prsng.prslblprv, r); err != nil {
					break
				}
			}
			if rerr != nil {
				err = rerr
			}
		}
	}()
	if err == io.EOF || err == nil {
		if err == io.EOF {
			err = nil
		}
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

func nextparsing(atv *Active, prntprsng *parsing, rin io.Reader, wout io.Writer, initpath string) (prsng *parsing) {
	prsng = &parsing{Buffer: iorw.NewBuffer(), prsvpth: initpath, rin: rin, wout: wout, woutbytes: make([]byte, 8192), woutbytesi: 0, prntprsng: prntprsng, atv: atv, cdetxt: rune(0), prslbli: []int{0, 0}, prslblprv: []rune{0, 0}, cdeoffsetstart: -1, cdeoffsetend: -1, psvoffsetstart: -1, psvoffsetend: -1, psvr: make([]rune, 8192), cder: make([]rune, 8192), prntrs: []io.Writer{},
		crntpsvsctn: nil, prvelmrn: rune(0), elmoffset: -1, elmlbli: []int{0, 0}, elmprvrns: []rune{rune(0), rune(0)}}
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

func transformCode(code string, opts map[string]interface{}) (trsnfrmdcde string, isrequired bool, err error) {
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

func (atvrntme *atvruntime) dispose(clear ...bool) {
	if atvrntme != nil {
		var clearonly = len(clear) > 0 && clear[0]
		if atvrntme.prsng != nil {
			atvrntme.prsng.dispose()
			atvrntme.prsng = nil
		}
		if atvrntme.atv != nil {
			if !clearonly {
				atvrntme.atv = nil
			}
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
			if !clearonly {
				atvrntme.vm = nil
			}
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
				atvrntme.prsng.incprint(w)
			}
		},
		"resetprint": func() {
			if atvrntme.prsng != nil {
				atvrntme.prsng.resetprint()
			}
		},
		"decprint": func() {
			if atvrntme.prsng != nil {
				atvrntme.prsng.decprint()
			}
		},
		"print": func(a ...interface{}) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				atvrntme.prsng.print(a...)
			} else if atvrntme.atv != nil {
				atvrntme.atv.print(nil, a...)
			}
		},
		"println": func(a ...interface{}) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				atvrntme.prsng.println(a...)
			} else if atvrntme.atv != nil {
				atvrntme.atv.println(nil, a...)
			}
		},
		"binwrite": func(b ...byte) (n int, err error) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				n, err = atvrntme.prsng.binwrite(b...)
			} else if atvrntme.atv != nil {
				n, err = atvrntme.atv.binwrite(nil, b...)
			}
			return
		},
		//READER
		"incread": func(r io.Reader) {
			if atvrntme.prsng != nil {
				atvrntme.prsng.incread(r)
			}
		},
		"resetread": func() {
			if atvrntme.prsng != nil {
				atvrntme.prsng.resetread()
			}
		},
		"decread": func() {
			if atvrntme.prsng != nil {
				atvrntme.prsng.decread()
			}
		},
		"seek": func(offset int64, whence int) (n int64, err error) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				n, err = atvrntme.prsng.seek(offset, whence)
			} else if atvrntme.atv != nil {
				n, err = atvrntme.atv.seek(nil, offset, whence)
			}
			return
		},
		"readln": func() (ln string, err error) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				ln, err = atvrntme.prsng.readLn()
			} else if atvrntme.atv != nil {
				ln, err = atvrntme.atv.readln(nil)
			}
			return
		},
		"readLines": func() (lines []string, err error) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				lines, err = atvrntme.prsng.readLines()
			} else if atvrntme.atv != nil {
				lines, err = atvrntme.atv.readlines(nil)
			}
			return
		}, "readAll": func() (s string, err error) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				s, err = atvrntme.prsng.readAll()
			} else if atvrntme.atv != nil {
				s, err = atvrntme.atv.readAll(nil)
			}
			return
		}, "binread": func(size int) (b []byte, err error) {
			if atvrntme.prsng != nil && atvrntme.prsng.atv != nil {
				b, err = atvrntme.prsng.binread(size)
			} else if atvrntme.atv != nil {
				b, err = atvrntme.atv.binread(nil, size)
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

func (atvrntme *atvruntime) lclvm(objmapref ...map[string]interface{}) (vm *goja.Runtime) {
	if atvrntme != nil {
		if atvrntme.vm == nil {
			atvrntme.vm = goja.New()
			var fldmppr = &fieldmapper{fldmppr: goja.UncapFieldNameMapper()}
			atvrntme.vm.SetFieldNameMapper(fldmppr)
			var dne = make(chan bool, 1)
			go func(vm *goja.Runtime) {
				defer func() { dne <- true }()
				jsext.Register(vm)
			}(atvrntme.vm)
			<-dne
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
	var errpgrm error = nil
	if requirejsprgm, errpgrm = goja.Compile("", requirejs.RequireJSString(), false); errpgrm != nil {
		fmt.Println(errpgrm.Error())
	}
}

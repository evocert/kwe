package parsing

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/evocert/kwe/iorw"
)

var prslbl = [][]rune{[]rune("<@"), []rune("@>")}
var elmlbl = [][]rune{[]rune("<#"), []rune(">"), []rune("</#"), []rune(">"), []rune("<#"), []rune("/>")}
var phrslbl = [][]rune{[]rune("{#"), []rune("#}")}

type AltActiveAPI interface {
	AltLookupTemplate(string, ...interface{}) (io.Reader, error)
	AltPrint(w io.Writer, a ...interface{})
	AltPrintln(w io.Writer, a ...interface{})
	AltBinWrite(w io.Writer, b ...byte) (n int, err error)
	AltReadln(r io.Reader) (ln string, err error)
	AltSeek(r io.Reader, offset int64, whence int) (n int64, err error)
	AltReadlines(r io.Reader) (lines []string, err error)
	AltReadAll(r io.Reader) (s string, err error)
	AltBinRead(r io.Reader, size int) (b []byte, err error)
	AltObjectRef() map[string]interface{}
	ProcessParsing(prsng *Parsing) (err error)
}

type Parsing struct {
	AtvActv AltActiveAPI
	*iorw.Buffer
	tmpltbuf       *iorw.Buffer
	tmpltmap       map[string][]int64
	wout           io.Writer
	woutbytes      []byte
	woutbytesi     int
	rin            io.Reader
	prntrs         []io.Writer
	rdrs           []io.Reader
	prslbli        []int
	prslblprv      []rune
	prntprsng      *Parsing
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
	//psvmap         map[int][]int64
	psvr  []rune
	psvri int
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
	Prsvpth string
}

func EvalParsing(prsng *Parsing, atv AltActiveAPI, wout io.Writer, rin io.Reader, initpath string, canexec bool, invertactpsv bool, a ...interface{}) (err error) {
	func() {
		if prsng == nil {
			prsng = NextParsing(atv, nil, rin, wout, initpath)
			defer prsng.Dispose()
		}
		if len(a) > 0 {
			if invertactpsv {
				a = append(append([]interface{}{"<@"}, a...), "@>")
			}
		}
		var prcssprsng func(prsng *Parsing) (err error) = nil
		if al := len(a); al > 0 {
			var ai = 0

			for ai < al {
				if d := a[ai]; d != nil {
					if dprcssprsng, _ := d.(func(prsng *Parsing) (err error)); dprcssprsng != nil {
						if prcssprsng == nil {
							prcssprsng = dprcssprsng
						}
						a = append(a[:ai], a[ai+1:]...)
						al--
						continue
					}
				}
				ai++
			}
		}
		if atv != nil {
			prcssprsng = atv.ProcessParsing
		}
		err = ParsePrsng(prsng, canexec, prcssprsng, a...)
	}()
	return
}
func ParseEval(prsng *Parsing, forceCode bool, callcode func(string, map[string]interface{}, ...string) (val interface{}, err error), a ...interface{}) (val interface{}, err error) {
	if prsng != nil {
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
						cde := Code(prsng, cdecoords...)
						if callcode != nil {
							val, err = callcode(cde, nil)
						}
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
	}
	return
}

func Code(prsng *Parsing, coords ...int64) (c string) {
	if prsng != nil {
		if cdel := len(prsng.cdemap); cdel > 0 {
			var cdei = 0
			var rdr *iorw.BuffReader = nil
			if len(coords) == 0 {
				coords = []int64{prsng.cdemap[cdei][0], prsng.cdemap[cdel-1][1]}
			}
			var mxdcde int64 = int64(0)
			if len(coords) == 2 && coords[0] <= coords[1] {
				mxdcde = coords[1] - coords[0]
			}
			if mxdcde > 0 {
				for cdei < cdel && mxdcde > 0 {
					if cdecrds, cdecrdsok := prsng.cdemap[cdei]; cdecrdsok && (cdecrds[0] <= coords[1] && cdecrds[1] >= coords[0]) {
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
										rdr = prsng.cdeBuff().Reader()
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

func (prsng *Parsing) tempbuf() *iorw.Buffer {
	if prsng.tmpbuf == nil {
		prsng.tmpbuf = iorw.NewBuffer()
	}
	return prsng.tmpbuf
}

func (prsng *Parsing) cdeBuff() *iorw.Buffer {
	if prsng.cdebuf == nil {
		prsng.cdebuf = iorw.NewBuffer()
	}
	return prsng.cdebuf
}

func (prsng *Parsing) tmpltBuf() *iorw.Buffer {
	if prsng != nil {
		if prsng.tmpltbuf == nil {
			prsng.tmpltbuf = iorw.NewBuffer()
		}
		return prsng.tmpltbuf
	}
	return nil
}

func (prsng *Parsing) tmpltMap() map[string][]int64 {
	if prsng != nil {
		if prsng.tmpltmap == nil {
			prsng.tmpltmap = map[string][]int64{}
		}
		return prsng.tmpltmap
	}
	return nil
}

func (prsng *Parsing) tmpltrdr(tmpltnme string) (rdr *iorw.BuffReader, mxlen int64) {
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

func (prsng *Parsing) Print(a ...interface{}) {
	if prsng.AtvActv != nil {
		if pl := len(prsng.prntrs); pl > 0 {
			prsng.AtvActv.AltPrint(prsng.prntrs[pl-1], a...)
		} else {
			prsng.AtvActv.AltPrint(prsng.wout, a...)
		}
	}
}

func (prsng *Parsing) Println(a ...interface{}) {
	if prsng.AtvActv != nil {
		if pl := len(prsng.prntrs); pl > 0 {
			prsng.AtvActv.AltPrintln(prsng.prntrs[pl-1], a...)
		} else {
			prsng.AtvActv.AltPrintln(prsng.wout, a...)
		}
	}
}

func (prsng *Parsing) BinWrite(b ...byte) (n int, err error) {
	if prsng.AtvActv != nil {
		if pl := len(prsng.prntrs); pl > 0 {
			n, err = prsng.AtvActv.AltBinWrite(prsng.prntrs[pl-1], b...)
		} else {
			n, err = prsng.AtvActv.AltBinWrite(prsng.wout, b...)
		}
	}
	return
}

func (prsng *Parsing) IncPrint(w io.Writer) {
	if prsng != nil {
		prsng.prntrs = append(prsng.prntrs, w)
	}
}

func (prsng *Parsing) ResetPrint() {
	if prsng.prntrs != nil {
		for len(prsng.prntrs) > 0 {
			prsng.prntrs[len(prsng.prntrs)-1] = nil
			prsng.prntrs = prsng.prntrs[:len(prsng.prntrs)-1]
		}
	}
}

func (prsng *Parsing) DecPrint() {
	if prsng.prntrs != nil {
		if len(prsng.prntrs) > 0 {
			prsng.prntrs[len(prsng.prntrs)-1] = nil
			prsng.prntrs = prsng.prntrs[:len(prsng.prntrs)-1]
		}
	}
}

func (prsng *Parsing) ReadLn() (ln string, err error) {
	if prsng.AtvActv != nil {
		if pl := len(prsng.rdrs); pl > 0 {
			ln, err = prsng.AtvActv.AltReadln(prsng.rdrs[pl-1])
		} else {
			ln, err = prsng.AtvActv.AltReadln(prsng.rin)
		}
	}
	return
}

func (prsng *Parsing) Seek(offset int64, whence int) (n int64, err error) {
	if prsng.AtvActv != nil {
		if pl := len(prsng.rdrs); pl > 0 {
			n, err = prsng.AtvActv.AltSeek(prsng.rdrs[pl-1], offset, whence)
		} else {
			n, err = prsng.AtvActv.AltSeek(prsng.rin, offset, whence)
		}
	}
	return
}

func (prsng *Parsing) ReadLines() (lines []string, err error) {
	if prsng.AtvActv != nil {
		if pl := len(prsng.rdrs); pl > 0 {
			lines, err = prsng.AtvActv.AltReadlines(prsng.rdrs[pl-1])
		} else {
			lines, err = prsng.AtvActv.AltReadlines(prsng.rin)
		}
	}
	return
}

func (prsng *Parsing) ReadAll() (s string, err error) {
	if prsng.AtvActv != nil {
		if pl := len(prsng.rdrs); pl > 0 {
			s, err = prsng.AtvActv.AltReadAll(prsng.rdrs[pl-1])
		} else {
			s, err = prsng.AtvActv.AltReadAll(prsng.rin)
		}
	}
	return
}

func (prsng *Parsing) IncRead(r io.Reader) {
	if prsng != nil {
		prsng.rdrs = append(prsng.rdrs, r)
	}
}

func (prsng *Parsing) BinRead(size int) (b []byte, err error) {
	if prsng.AtvActv != nil {
		if pl := len(prsng.rdrs); pl > 0 {
			b, err = prsng.AtvActv.AltBinRead(prsng.rdrs[pl-1], size)
		} else {
			b, err = prsng.AtvActv.AltBinRead(prsng.rin, size)
		}
	}
	return
}

func (prsng *Parsing) ResetRead() {
	if prsng.rdrs != nil {
		for len(prsng.rdrs) > 0 {
			prsng.rdrs[len(prsng.rdrs)-1] = nil
			prsng.rdrs = prsng.rdrs[:len(prsng.rdrs)-1]
		}
	}
}

func (prsng *Parsing) DecRead() {
	if prsng.rdrs != nil {
		if len(prsng.rdrs) > 0 {
			prsng.rdrs[len(prsng.rdrs)-1] = nil
			prsng.rdrs = prsng.rdrs[:len(prsng.rdrs)-1]
		}
	}
}

func (prsng *Parsing) Dispose() {
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
		if prsng.AtvActv != nil {
			prsng.AtvActv = nil
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

func (prsng *Parsing) TopPrsng() *Parsing {
	if prsng.prntprsng == nil {
		return prsng
	}
	return prsng.prntprsng.TopPrsng()
}

func (prsng *Parsing) setcdepos(startoffset int64, endoffset int64) {
	if prsng.cdemap == nil {
		prsng.cdemap = map[int][]int64{}
	}
	prsng.cdemap[len(prsng.cdemap)] = []int64{startoffset, endoffset}
}

func (prsng *Parsing) flushWritePsv() (err error) {
	if prsng != nil && prsng.woutbytesi > 0 {
		_, err = prsng.wout.Write(prsng.woutbytes[0:prsng.woutbytesi])
		prsng.woutbytesi = 0
	}
	return
}

func (prsng *Parsing) writePsv(p ...rune) (err error) {
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

func (prsng *Parsing) writeCde(p []rune) (err error) {
	if pl := len(p); pl > 0 {
		if prsng.cdeoffsetstart == -1 {
			prsng.cdeoffsetstart = prsng.cdeBuff().Size()
		}
		err = prsng.cdeBuff().WriteRunes(p[:pl]...)
	}
	return
}

func (prsng *Parsing) foundCode() bool {
	return prsng.foundcde
}

func (prsng *Parsing) flushPsv() (err error) {
	if pi := prsng.psvri; pi > 0 {
		prsng.psvri = 0
		err = prsng.writePsv(prsng.psvr[:pi]...)
	}
	if err == nil {
		err = prsng.flushWritePsv()
	}
	if err == nil && prsng.crntpsvsctn == nil && prsng.foundCode() {
		if psvoffsetstart := prsng.psvoffsetstart; psvoffsetstart > -1 {
			psvoffsetend := prsng.Size()
			prsng.psvoffsetstart = -1
			//err = parseatvrunes(prsng, []rune(fmt.Sprintf("print(_psvsub(%d,%d));", psvoffsetstart, psvoffsetend)))
			//pos := prsng.setpsvpos(psvoffsetstart, prsng.Size())
			//err = parseatvrunes(prsng, []rune(fmt.Sprintf("_psvout(%d);", pos)))
			//if psvouts := PassiveoutS(prsng, pos); psvouts != "" {
			//	err = parseatvrunes(prsng, []rune(fmt.Sprintf("print(`%s`);", psvouts)))
			//}

			if psvouts := PassiveoutSubString(prsng, psvoffsetstart, psvoffsetend); psvouts != "" {
				err = parseatvrunes(prsng, []rune(fmt.Sprintf("print(`%s`);", psvouts)))
			}

		}
	}
	return
}

func parsepsvrunes(prsng *Parsing, p []rune) (err error) {
	if len(p) > 0 {
		for _, rn := range p {
			if err = parsepsvrune(prsng, rn); err != nil {
				break
			}
		}
	}
	return
}

func parseatvrunes(prsng *Parsing, p []rune) (err error) {
	if len(p) > 0 {
		for _, rn := range p {
			if err = parseatvrune(prsng, rn); err != nil {
				break
			}
		}
	}
	return
}

func (prsng *Parsing) flushCde() (err error) {
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

func ParsePrsng(prsng *Parsing, canexec bool, performParsing func(prsng *Parsing) (err error), a ...interface{}) (err error) {
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
				if performParsing != nil {
					err = performParsing(prsng)
				}
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

func parseprsngrune(prsng *Parsing, prslbli []int, prslblprv []rune, pr rune) (err error) {
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

func NextParsing(atvActv AltActiveAPI, prntprsng *Parsing, rin io.Reader, wout io.Writer, initpath string) (prsng *Parsing) {
	prsng = &Parsing{
		AtvActv: atvActv, Buffer: iorw.NewBuffer(), Prsvpth: initpath, rin: rin, wout: wout, woutbytes: make([]byte, 8192), woutbytesi: 0, prntprsng: prntprsng, cdetxt: rune(0), prslbli: []int{0, 0}, prslblprv: []rune{0, 0}, cdeoffsetstart: -1, cdeoffsetend: -1, psvoffsetstart: -1, psvoffsetend: -1, psvr: make([]rune, 8192), cder: make([]rune, 8192), prntrs: []io.Writer{},
		crntpsvsctn: nil, prvelmrn: rune(0), elmoffset: -1, elmlbli: []int{0, 0}, elmprvrns: []rune{rune(0), rune(0)}}
	return
}

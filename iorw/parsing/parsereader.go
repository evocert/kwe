package parsing

import (
	"io"
	"strings"
	"sync"

	"github.com/evocert/kwe/iorw"
)

type ParsingReader struct {
	mltiargsrdr *iorw.MultiArgsReader
	prsng       *Parsing
	tmplts      map[string]*iorw.Buffer
	includes    map[string]*iorw.Buffer
	//
	intlbli   []int
	frasebuf  *iorw.Buffer
	finrunes  []rune
	finrunesi int
	cchdr     rune
	cchdsize  int
	errrns    error
	prvpr     rune
	bufoutwtr *iorw.Buffer
	bufoutrdr *iorw.BuffReader
}

func (prsngrdr *ParsingReader) flushfin() {
	if prsngrdr.finrunesi > 0 {
		prsngrdr.bufoutwtr.WriteRunes(prsngrdr.finrunes[:prsngrdr.finrunesi]...)
		prsngrdr.finrunesi = 0
	}
}

func (prsngrdr *ParsingReader) fincrn(frs ...rune) {
	for _, fr := range frs {
		prsngrdr.finrunes[prsngrdr.finrunesi] = fr
		prsngrdr.finrunesi++
		if prsngrdr.finrunesi == len(prsngrdr.finrunes) {
			prsngrdr.flushfin()
		}
	}
}

func (prsngrdr *ParsingReader) intparsechar(pr rune) {
	if prsngrdr.intlbli[1] == 0 && prsngrdr.intlbli[0] < len(intlbl[0]) {
		if prsngrdr.intlbli[0] > 0 && intlbl[0][prsngrdr.intlbli[0]-1] == prsngrdr.prvpr && intlbl[0][prsngrdr.intlbli[0]] != pr {
			prsngrdr.fincrn(intlbl[0][:prsngrdr.intlbli[0]]...)
			prsngrdr.intlbli[0] = 0
		}
		if intlbl[0][prsngrdr.intlbli[0]] == pr {
			prsngrdr.intlbli[0]++
			if prsngrdr.intlbli[0] == len(intlbl[0]) {
				prsngrdr.prvpr = rune(0)
			} else {
				prsngrdr.prvpr = pr
			}
		} else {
			if prsngrdr.intlbli[0] > 0 {
				prsngrdr.fincrn(intlbl[0][:prsngrdr.intlbli[0]]...)
				prsngrdr.intlbli[0] = 0
			}
			prsngrdr.prvpr = pr
			prsngrdr.fincrn(pr)
		}
	} else if prsngrdr.intlbli[0] == len(intlbl[0]) && prsngrdr.intlbli[1] < len(intlbl[1]) {
		if intlbl[1][prsngrdr.intlbli[1]] == pr {
			prsngrdr.intlbli[1]++
			if prsngrdr.intlbli[1] == len(intlbl[1]) {
				prsngrdr.intlbli[0] = 0
				prsngrdr.intlbli[1] = 0
				if prsngrdr.frasebuf.Size() > 0 {
					prsngrdr.flushfin()
					processPreParsingFrase(prsngrdr.bufoutwtr, prsngrdr.prsng, prsngrdr, prsngrdr.frasebuf)
					prsngrdr.frasebuf.Clear()
				} else {
					prsngrdr.fincrn(intlbl[0]...)
					prsngrdr.fincrn(intlbl[1]...)
				}
			}
		} else {
			if prsngrdr.intlbli[1] > 0 {
				for _, fr := range intlbl[1][:prsngrdr.intlbli[1]] {
					prsngrdr.frasebuf.WriteRune(fr)
				}
				prsngrdr.intlbli[1] = 0
			}
			prsngrdr.frasebuf.WriteRune(pr)
		}
	}
}

func NewParseReader(prsng *Parsing, a ...interface{}) (prsngrdr *ParsingReader) {
	if len(a) > 0 && prsng != nil {
		if prsngrdr, _ = prsgnrdrpool.Get().(*ParsingReader); prsngrdr == nil {
			prsngrdr = newprsngrdr()
		}

		if prsngrdr.prsng != prsng {
			prsngrdr.prsng = prsng
		}
		prsngrdr.mltiargsrdr = iorw.NewMultiArgsReader(a...)

		if prsngrdr.bufoutrdr != nil {
			prsngrdr.bufoutrdr.Close()
			prsngrdr.bufoutrdr = nil
		}
		prsngrdr.bufoutrdr = prsngrdr.bufoutwtr.Reader()
	}
	return
}

var prsgnrdrpool = &sync.Pool{New: func() interface{} { return newprsngrdr() }}

func newprsngrdr() (prsngrdr *ParsingReader) {
	prsngrdr = &ParsingReader{
		intlbli:   []int{0, 0},
		frasebuf:  iorw.NewBuffer(),
		finrunes:  make([]rune, 4096),
		finrunesi: 0, cchdr: rune(0), cchdsize: 0, prvpr: rune(0),
		bufoutwtr: iorw.NewBuffer()}
	return
}

func resetprsngrdr(prsngrdr *ParsingReader) {
	if prsngrdr.mltiargsrdr != nil {
		prsngrdr.mltiargsrdr.Close()
		prsngrdr.mltiargsrdr = nil
	}
	if prsngrdr.prsng != nil {
		prsngrdr.prsng = nil
	}
	if prsngrdr.includes != nil {
		for incld := range prsngrdr.includes {
			prsngrdr.includes[incld].Close()
			prsngrdr.includes[incld] = nil
			delete(prsngrdr.includes, incld)
		}
		prsngrdr.includes = nil
	}
	if prsngrdr.tmplts != nil {
		for tmpl := range prsngrdr.tmplts {
			prsngrdr.tmplts[tmpl].Close()
			prsngrdr.tmplts[tmpl] = nil
			delete(prsngrdr.tmplts, tmpl)
		}
		prsngrdr.tmplts = nil
	}
	if prsngrdr.bufoutwtr != nil {
		prsngrdr.bufoutwtr.Clear()
	}
	if prsngrdr.bufoutrdr != nil {
		prsngrdr.bufoutrdr.Close()
		prsngrdr.bufoutrdr = nil
	}
	prsngrdr.intlbli[0] = 0
	prsngrdr.intlbli[1] = 0
	prsngrdr.prvpr = rune(0)
	prsngrdr.frasebuf.Clear()
	prsngrdr.finrunesi = 0
	prsngrdr.cchdr = rune(0)
	prsngrdr.cchdsize = 0
}

var intlbl = [][]rune{[]rune("[@"), []rune("@]")}

func processPreParsingFrase(bufrnwtrr *iorw.Buffer, prsng *Parsing, prsngrdr *ParsingReader, frasebuf *iorw.Buffer) {
	if prsng != nil {
		if frasebuf.HasPrefix("include ") {
			if bufr := frasebuf.Reader(); bufr != nil {
				func() {
					defer bufr.Close()
					var ln = len([]byte("include "))
					var frases = strings.TrimSpace(bufr.SubString(int64(ln), frasebuf.Size()))
					if frases != "" {
						if inclbuf, _ := prsngrdr.includes[frases]; inclbuf != nil {
							func() {
								prsngrdr.mltiargsrdr.InsertArgs(inclbuf.Reader())
								frases = ""
							}()
						} else if prsng.LookupTemplate != nil {
							if r, _ := prsng.LookupTemplate(frases); r != nil {
								func() {
									if prsngrdr.includes == nil {
										prsngrdr.includes = map[string]*iorw.Buffer{}
									}
									prsngrdr.includes[frases] = iorw.NewBuffer()
									prsngrdr.includes[frases].ReadFrom(r)
									prsngrdr.mltiargsrdr.InsertArgs(prsngrdr.includes[frases].Reader())
									frases = ""
								}()
							}
						} else if prsng.AtvActv != nil {
							if r, _ := prsng.AtvActv.AltLookupTemplate(frases); r != nil {
								func() {
									if prsngrdr.includes == nil {
										prsngrdr.includes = map[string]*iorw.Buffer{}
									}
									prsngrdr.includes[frases] = iorw.NewBuffer()
									prsngrdr.includes[frases].ReadFrom(r)
									prsngrdr.mltiargsrdr.InsertArgs(prsngrdr.includes[frases].Reader())
									frases = ""
								}()
							}
						}
					}
				}()
			}
		} else if frasebuf.HasPrefix("include-template") {
			func() {
				var tmpltrdr = frasebuf.Reader()
				defer tmpltrdr.Close()
				tmpltrdr.Seek(int64(len("include-template")), io.SeekStart)
				var tmplerr error = nil
				var tmplname = ""
				for tmplerr == nil {
					if tmplname, tmplerr = tmpltrdr.Readln(); tmplname != "" {
						tmplname = strings.TrimSpace(tmplname)

						if tmplname != "" {
							var tmpltbuf = iorw.NewBuffer()
							tmpltbuf.ReadRunesFrom(tmpltrdr)
							if tmpltbuf.Size() > 0 {
								if prsngrdr.tmplts == nil {
									prsngrdr.tmplts = map[string]*iorw.Buffer{}
								}
								if tmplbf, _ := prsngrdr.tmplts[tmplname]; tmplbf != nil {
									tmplbf.Clear()
									tmpltrdr.Close()
									tmpltrdr = tmpltbuf.Reader()
									tmplbf.ReadFrom(tmpltrdr)
									tmpltrdr.Close()
								} else {
									prsngrdr.tmplts[tmplname] = tmpltbuf
								}
							} else {
								tmpltbuf.Close()
							}
						}
						break
					} else {
						break
					}
				}
			}()
		} else if frasebuf.HasPrefix("tmpl:") {
			if bufr := frasebuf.Reader(); bufr != nil {
				func() {
					defer bufr.Close()
					var ln = len([]byte("tmpl:"))
					var frases = strings.TrimSpace(bufr.SubString(int64(ln), frasebuf.Size()))
					if frases != "" {
						if tmplbuf, _ := prsngrdr.tmplts[frases]; tmplbuf != nil {
							func() {
								prsngrdr.mltiargsrdr.InsertArgs(tmplbuf.Reader())
								frases = ""
							}()
						}
					}
				}()
			}
		}
	}
}

func (prsngrdr *ParsingReader) ReadRune() (r rune, size int, err error) {
	r, size, err = internalReadRune(prsngrdr, prsngrdr.mltiargsrdr)
	return
}

func internalReadRune(prsngrdr *ParsingReader, mltiargsrdrd *iorw.MultiArgsReader) (r rune, size int, err error) {
	if prsngrdr.bufoutwtr.Size() == 0 {
		for prsngrdr.bufoutwtr.Size() == 0 {
			prsngrdr.cchdr, prsngrdr.cchdsize, prsngrdr.errrns = mltiargsrdrd.ReadRune()
			if prsngrdr.cchdsize > 0 && prsngrdr.errrns == nil {
				prsngrdr.intparsechar(prsngrdr.cchdr)
			} else if prsngrdr.errrns != nil {
				if prsngrdr.cchdsize > 0 && prsngrdr.errrns == io.EOF {
					prsngrdr.errrns = nil
					prsngrdr.intparsechar(prsngrdr.cchdr)
				}
				if prsngrdr.errrns == io.EOF {
					break
				} else {
					break
				}
			}
		}
		prsngrdr.flushfin()
		if prsngrdr.bufoutwtr.Size() > 0 {
			prsngrdr.bufoutrdr.Close()
			prsngrdr.bufoutrdr = prsngrdr.bufoutwtr.Reader()
		} else {
			err = prsngrdr.errrns
			return
		}
	}
	r, size, err = prsngrdr.bufoutrdr.ReadRune()
	if err == io.EOF {
		prsngrdr.bufoutwtr.Clear()
		r, size, err = internalReadRune(prsngrdr, mltiargsrdrd)
	}
	return
}

func (prsngrdr *ParsingReader) Close() (err error) {
	if prsngrdr != nil {
		resetprsngrdr(prsngrdr)
		prsgnrdrpool.Put(prsngrdr)
		prsngrdr = nil
	}
	return
}

package parsing

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/evocert/kwe/iorw"
)

func Passiveout(prsng *Parsing, i int) {
	if prsng != nil {
		if psvl := len(prsng.psvmap); psvl > 0 && i >= 0 && i < psvl {
			psvcoors := prsng.psvmap[i]
			if psvcoors[1] > psvcoors[0] {
				rdr := prsng.Reader()
				rdr.Seek(psvcoors[0], 0)
				io.CopyN(prsng.wout, rdr, psvcoors[1]-psvcoors[0])
			}
		}
	}
}

func parsepsvrune(prsng *Parsing, rn rune) (err error) {
	prsng.flushCde()
	if prsng.hascde {
		prsng.hascde = false
	}
	err = parseelmpsvrrune(prsng, prsng.elmoffset, prsng.elmlbli, prsng.elmprvrns, rn)
	return
}

func parsepsvphrase(prsng *Parsing, psvsctn *psvsection, phrslbli []int, rn rune) (err error) {
	if phrslbli[1] == 0 && phrslbli[0] < len(phrslbl[0]) {
		if phrslbli[0] > 0 && phrslbl[0][phrslbli[0]-1] == psvsctn.phrsprvrn && phrslbl[0][phrslbli[0]] != rn {
			phrsi := phrslbli[0]
			phrslbli[0] = 0
			psvsctn.phrsprvrn = rune(0)
			psvsctn.CachedBuf().WriteRunes(phrslbl[0][:phrsi]...)
		}
		if phrslbl[0][phrslbli[0]] == rn {
			phrslbli[0]++
			if phrslbli[0] == len(phrslbl[0]) {
				psvsctn.phrsprvrn = rune(0)
			} else {
				psvsctn.phrsprvrn = rn
			}
		} else {
			if phrsi := phrslbli[0]; phrsi > 0 {
				phrslbli[0] = 0
				psvsctn.phrsprvrn = rune(0)
				psvsctn.CachedBuf().WriteRunes(phrslbl[0][:phrsi]...)
			}
			psvsctn.phrsprvrn = rn
			psvsctn.CachedBuf().WriteRune(rn)
		}
	} else if phrslbli[0] == len(phrslbl[0]) && phrslbli[1] < len(phrslbl[1]) {
		if phrslbl[1][phrslbli[1]] == rn {
			phrslbli[1]++
			if phrslbli[1] == len(phrslbl[1]) {
				var phrsfound = ""
				if psvsctn.phrsbuf != nil && psvsctn.phrsbuf.Size() > 0 {
					if phrsfound = psvsctn.phrsbuf.String(); phrsfound != "" {
						psvsctn.phrsbuf.Clear()
						phrslbli[1] = 0
						phrslbli[0] = 0
						psvsctn.phrsprvrn = rune(0)
						if phrscoord, phrsok := psvsctn.phrsmap[phrsfound]; phrsok && phrscoord[1] > phrscoord[0] {
							if phrsrdr := psvsctn.PhraseTemplateBuf().Reader(); phrsrdr != nil {
								func() {
									defer func() {
										phrsrdr.Close()
										phrsrdr = nil
									}()
									phrsrdr.Seek(phrscoord[0], io.SeekStart)
									phrsrdr.MaxRead = phrscoord[1] - phrscoord[0]
									err = ParsePrsng(psvsctn.prsng, false, nil, phrsrdr)
								}()
								return
							}
						} else {
							if phrsfound != "content" {
								psvsctn.CachedBuf().WriteRunes(phrslbl[0]...)
								psvsctn.CachedBuf().WriteRunes([]rune(phrsfound)...)
								psvsctn.CachedBuf().WriteRunes(phrslbl[1]...)
							}
						}
					}
				} else {
					psvsctn.CachedBuf().WriteRunes(phrslbl[0]...)
					psvsctn.CachedBuf().WriteRunes(phrslbl[1]...)
				}
				phrslbli[1] = 0
				phrslbli[0] = 0
				psvsctn.phrsprvrn = rune(0)

			}
		} else {
			if phrsi := phrslbli[1]; phrsi > 0 {
				psvsctn.CachedBuf().WriteRunes(phrslbl[0]...)
				if psvsctn.phrsbuf != nil && psvsctn.phrsbuf.Size() > 0 {
					psvsctn.CachedBuf().Print(psvsctn.phrsbuf.String())
					psvsctn.phrsbuf.Clear()
				}
				psvsctn.CachedBuf().WriteRunes(phrslbl[1]...)
				phrslbli[1] = 0
				phrslbli[0] = 0
				psvsctn.phrsprvrn = rune(0)
				return
			}
			if strings.TrimSpace(string(rn)) != "" {
				psvsctn.PhraseBuf().WriteRune(rn)
			} else {
				psvsctn.CachedBuf().WriteRunes(phrslbl[0]...)
				if psvsctn.phrsbuf != nil && psvsctn.phrsbuf.Size() > 0 {
					psvsctn.CachedBuf().Print(psvsctn.phrsbuf.String())
					psvsctn.phrsbuf.Clear()
				}
				psvsctn.CachedBuf().WriteRune(rn)
				phrslbli[1] = 0
				phrslbli[0] = 0
				psvsctn.phrsprvrn = rune(0)
			}
		}
	}
	return
}

func parseelmpsvrrune(prsng *Parsing, elmoffset int, elmlbli []int, elmprvrns []rune, rn rune) (err error) {
	if elmoffset == -1 {
		elmoffset = 0
		prsng.elmoffset = elmoffset
	} else {
		if rn == '/' {
			if elmoffset == 0 && elmlbli[1] == 0 && elmlbli[0] == 1 && elmprvrns[0] == '<' {
				elmoffset = 2
				prsng.elmoffset = elmoffset
				elmprvrns[0] = rune(0)
			} else if elmoffset == 0 && elmlbli[0] == len(elmlbl[elmoffset]) {
				elmoffset = 4
				elmlbli[0] = len(elmlbl[elmoffset])
				prsng.elmoffset = elmoffset
				elmprvrns[1] = rune(0)
			}
		}
	}
	if elmlbli[1] == 0 && elmlbli[0] < len(elmlbl[elmoffset]) {
		if elmlbli[0] > 0 && elmlbl[elmoffset][elmlbli[0]-1] == elmprvrns[0] && elmlbl[elmoffset][elmlbli[0]] != rn {
			elmri := elmlbli[0]
			elmlbli[0] = 0
			elmprvrns[0] = rune(0)
			prsng.writePsv(elmlbl[elmoffset][:elmri]...)
			if elmoffset > 0 {
				prsng.elmoffset = 0
				elmoffset = 0
			}
		}
		if elmlbl[elmoffset][elmlbli[0]] == rn {
			elmlbli[0]++
			if elmlbli[0] == len(elmlbl[elmoffset]) {
				if prsng.tmpbuf != nil {
					prsng.tmpbuf.Clear()
				}
				elmprvrns[0] = rune(0)
			} else {
				elmprvrns[0] = rn
			}
		} else {
			if elmlbli[0] > 0 {
				elmri := elmlbli[0]
				elmlbli[0] = 0
				elmprvrns[0] = rune(0)
				prsng.writePsv(elmlbl[elmoffset][:elmri]...)
				if elmoffset > 0 {
					prsng.elmoffset = 0
					elmoffset = 0
				}
			}
			elmprvrns[0] = rn
			prsng.writePsv(rn)
		}
	} else if elmlbli[0] == len(elmlbl[elmoffset]) && elmlbli[1] < len(elmlbl[elmoffset+1]) {
		if elmlbl[elmoffset+1][elmlbli[1]] == rn {
			elmlbli[1]++
			if elmlbli[1] == len(elmlbl[elmoffset+1]) {
				prsng.elmoffset = -1
				elmlbli[0] = 0
				elmlbli[1] = 0
				elmprvrns[0] = rune(0)
				elmprvrns[1] = rune(0)
				valid, elmTpe, psvsctn, verr := validElemParsing(prsng, elmoffset, prsng.crntpsvsctn)
				if verr != nil {
					err = verr
					return
				} else if valid && psvsctn != nil {
					if elmTpe == ElemEnd || elmTpe == ElemSingle {
						err = processPsvSection(psvsctn)
					}
				} else {
					if err = prsng.writePsv(elmlbl[elmoffset]...); err != nil {
						return
					}
					if prsng.tmpbuf != nil && prsng.tmpbuf.Size() > 0 {
						if tmprdr := prsng.tmpbuf.Reader(); tmprdr != nil {
							for {
								tmpr, tmps, tmperr := tmprdr.ReadRune()
								if tmps > 0 {
									prsng.writePsv(tmpr)
								}
								if tmperr != nil {
									tmprdr.Close()
									tmprdr = nil
									if tmperr != io.EOF {
										err = tmperr
										return
									}
									break
								}
							}
						}
					}
					if err = prsng.writePsv(elmlbl[elmoffset+1]...); err != nil {
						return
					}
				}
			}
		} else {
			if elmlbli[1] > 0 {
				elmrl := elmlbli[0]
				elmlbli[0] = 0
				elmprvrns[0] = rune(0)
				for _, tmprn := range elmlbl[elmoffset][:elmrl] {
					if err = prsng.tempbuf().WriteRune(tmprn); err != nil {
						return
					}
				}
			}
			elmprvrns[1] = rn
			err = prsng.tempbuf().WriteRune(rn)
		}
	}
	return
}

func validElemParsing(prsng *Parsing, elmoffset int, crntpsvsctn *psvsection) (valid bool, elmTpe elemtype, psvsctn *psvsection, err error) {
	if prsng.tmpbuf == nil || prsng.tmpbuf.Size() == 0 {
		return
	}
	valid = true
	elmTpe = ElemStart
	if elmoffset == 2 {
		elmTpe = ElemEnd
	} else if elmoffset == 4 {
		elmTpe = ElemSingle
	}
	if elmTpe == ElemStart || elmTpe == ElemSingle {
		if psvsctn = newPsvSection(prsng, elmTpe, prsng.tmpbuf, crntpsvsctn); psvsctn == nil {
			valid = false
		}
	} else {
		psvsctn = crntpsvsctn
	}
	return
}

type elemtype int

const (
	ElemNone elemtype = iota
	//ElemStart - elem start
	ElemStart
	//ElemEnd - elem end
	ElemEnd
	//ElemSingle - elem single
	ElemSingle
)

func (elmtpe elemtype) String() (s string) {
	if elmtpe == ElemEnd {
		s = "ELEM-END"
	} else if elmtpe == ElemSingle {
		s = "ELEM-SINGLE"
	} else if elmtpe == ElemStart {
		s = "ELEM-START"
	}
	return
}

type psvsection struct {
	prsng        *Parsing
	elmtpe       elemtype
	tmpbuf       *iorw.Buffer
	prvsctn      *psvsection
	nxtsctn      *psvsection
	chcdbf       *iorw.Buffer
	phrsbuf      *iorw.Buffer
	phrstmpltbuf *iorw.Buffer
	phrsmap      map[string][]int64
	canphrs      bool
	phrslbli     []int
	phrsprvrn    rune
	tmpltpath    string
	tmpstrti     int64
	tmpendi      int16
}

func removePsvSection(prsng *Parsing, psvsctn *psvsection) {
	if prsng != nil && psvsctn != nil && psvsctn.prsng == prsng {
		prvsctn := psvsctn.prvsctn
		nxtsctn := psvsctn.nxtsctn

		prsng.crntpsvsctn = prvsctn

		if prsng.headpsvsctn == psvsctn {
			if prsng.tailpsvsctn == psvsctn {
				prsng.headpsvsctn = nil
				prsng.tailpsvsctn = nil
			} else {
				prsng.headpsvsctn.nxtsctn = prvsctn
			}
		} else if prsng.tailpsvsctn == psvsctn {
			prsng.tailpsvsctn = prvsctn
		}
		if nxtsctn != nil {
			nxtsctn.prvsctn = prvsctn
		}
		if prvsctn != nil {
			prvsctn.nxtsctn = nxtsctn
		}
		prsng = nil
	}
}

func (psvsctn *psvsection) dispose() {
	if psvsctn != nil {
		if psvsctn.prsng != nil {
			removePsvSection(psvsctn.prsng, psvsctn)
			psvsctn.prsng = nil
		}
		if psvsctn.chcdbf != nil {
			psvsctn.chcdbf.Close()
			psvsctn.chcdbf = nil
		}
		if psvsctn.phrsmap != nil {
			for k := range psvsctn.phrsmap {
				psvsctn.phrsmap[k] = nil
				delete(psvsctn.phrsmap, k)
			}
			psvsctn.phrsmap = nil
		}
		if psvsctn.phrsbuf != nil {
			psvsctn.phrsbuf.Close()
			psvsctn.phrsbuf = nil
		}
		if psvsctn.tmpbuf != nil {
			psvsctn.tmpbuf.Close()
			psvsctn.tmpbuf = nil
		}
		psvsctn = nil
	}
}

func (psvsctn *psvsection) PhraseBuf() *iorw.Buffer {
	if psvsctn.phrsbuf == nil {
		psvsctn.phrsbuf = iorw.NewBuffer()
	}
	return psvsctn.phrsbuf
}

func (psvsctn *psvsection) PhraseTemplateBuf() *iorw.Buffer {
	if psvsctn.phrstmpltbuf == nil {
		psvsctn.phrstmpltbuf = iorw.NewBuffer()
	}
	return psvsctn.phrstmpltbuf
}

func (psvsctn *psvsection) path() (path string) {
	path = strings.Replace(psvsctn.tmpltpath, "|", "/", -1)
	prsngext := filepath.Ext(psvsctn.prsng.Prsvpth)
	prvsctnext := ""
	if strings.HasPrefix(path, ".") {
		if psvsctn.prvsctn != nil {
			prvsctnpth := psvsctn.prvsctn.path()
			prvsctnext = filepath.Ext(prvsctnpth)
			path = prvsctnpth[:strings.LastIndex(prvsctnpth, "/")+1] + path[1:]
		} else {
			prsngpth := psvsctn.prsng.Prsvpth
			path = prsngpth[:strings.LastIndex(prsngpth, "/")+1] + path[1:]
		}
	} else if !strings.HasPrefix(path, "/") {
		prsngpth := psvsctn.prsng.Prsvpth
		path = prsngpth[:strings.LastIndex(prsngpth, "/")+1] + path
	}
	if pthext := filepath.Ext(path); pthext == "" {
		if prvsctnext != "" {
			pthext = prvsctnext
		} else if prsngext != "" {
			pthext = prsngext
		}
		path = path + pthext
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return
}

func newPsvSection(prsng *Parsing, elmtpe elemtype, tmpbuf *iorw.Buffer, crntsctn *psvsection) (psvsctn *psvsection) {
	if tmpbuf != nil && tmpbuf.Size() > 0 {

		var tmpltpath = ""

		tmprdr := tmpbuf.Reader()
		var nxttmpbuf *iorw.Buffer = nil
		var rnferr error = nil
		tmpclbli := []int{0, 0}
		tmpclbl := [][]rune{[]rune("{@"), []rune("@}")}
		tmpclprnrn := rune(0)
		foundtmltnme := false
		for !foundtmltnme {
			rn, rns, rnerr := tmprdr.ReadRune()
			if rns > 0 {
				if strings.TrimSpace(string(rn)) == "" {
					tmprdr.Close()
					return
				}
				if tmpclbli[1] == 0 && tmpclbli[0] < len(tmpclbl[0]) {
					if tmpclbli[0] > 0 && tmpclbl[0][tmpclbli[0]-1] == tmpclprnrn && tmpclbl[0][tmpclbli[0]] != rn {
						tmprdr.Close()
						return
					}
					if tmpclbl[0][tmpclbli[0]] == rn {
						tmpclbli[0]++
						if tmpclbli[0] == len(tmpclbl[0]) {
							if tmpltpath == "" {
								tmprdr.Close()
								return
							}
							foundtmltnme = true
							nxttmpbuf = iorw.NewBuffer()
							nxttmpbuf.WriteRunes(prslbl[0]...)
							tmpclprnrn = rune(0)
						}
					} else {
						if tmpclbli[0] > 0 {
							tmprdr.Close()
							return
						}
						tmpclprnrn = rn
						tmpltpath += string(rn)
					}
				}
			}
			if rnerr != nil {
				if rnerr != io.EOF {
					rnferr = rnerr
				}
				break
			}
		}

		if rnferr == nil {
			if foundtmltnme {
				for {
					rn, rns, rnerr := tmprdr.ReadRune()
					if rns > 0 {
						if tmpclbli[1] == 0 && tmpclbli[0] < len(tmpclbl[0]) {
							if tmpclbli[0] > 0 && tmpclbl[0][tmpclbli[0]-1] == tmpclprnrn && tmpclbl[0][tmpclbli[0]] != rn {
								tmprdr.Close()
								nxttmpbuf.Close()
								return
							}
							if tmpclbl[0][tmpclbli[0]] == rn {
								tmpclbli[0]++
								if tmpclbli[0] == len(tmpclbl[0]) {
									nxttmpbuf.WriteRunes(prslbl[0]...)
									tmpclprnrn = rune(0)
								} else {
									tmpclprnrn = rn
								}
							} else {
								if tmpclbli[0] > 0 {
									tmprdr.Close()
									return
								}
								tmpclprnrn = rn
								tmpltpath += string(rn)
							}
						} else if tmpclbli[0] == len(tmpclbl[0]) && tmpclbli[1] < len(tmpclbl[1]) {
							if tmpclbl[1][tmpclbli[1]] == rn {
								tmpclbli[1]++
								if tmpclbli[1] == len(tmpclbl[1]) {
									nxttmpbuf.WriteRunes(prslbl[1]...)
									tmpclbli[0] = 0
									tmpclbli[1] = 0
									tmpclprnrn = rune(0)
								}
							} else {
								if tmpclbli[1] > 0 {
									tmprdr.Close()
									nxttmpbuf.Close()
									return
								}
								nxttmpbuf.WriteRune(rn)
							}
						}
					}
					if rnerr != nil {
						if rnerr != io.EOF || tmpclbli[0] == len(tmpclbl[0]) {
							tmprdr.Close()
							nxttmpbuf.Close()
							return
						}
						break
					}
				}
				tmprdr.Close()
				if nxttmpbuf != nil && nxttmpbuf.Size() > 0 {
					tmpbuf.Clear()
					tmpbuf.Print(nxttmpbuf)
				}
			} else {
				tmprdr.Close()
				tmpbuf.Clear()
			}
		} else {
			tmprdr.Close()
			return
		}

		psvsctn = &psvsection{prsng: prsng, elmtpe: elmtpe, prvsctn: crntsctn, tmpbuf: tmpbuf.Clone(), chcdbf: nil, tmpstrti: -1, tmpendi: -1,
			phrsmap: map[string][]int64{}, phrslbli: []int{0, 0}, phrsprvrn: rune(0)}
		if prsng.headpsvsctn == nil {
			prsng.headpsvsctn = psvsctn
		}
		if prsng.tailpsvsctn != nil {
			prsng.tailpsvsctn.nxtsctn = psvsctn
		}
		prsng.tailpsvsctn = psvsctn
		psvsctn.tmpltpath = tmpltpath
		prsng.crntpsvsctn = psvsctn
	}
	return
}

func (psvsctn *psvsection) Elemtype() elemtype {
	return psvsctn.elmtpe
}

func (psvsctn *psvsection) CachedBuf() *iorw.Buffer {
	if psvsctn.chcdbf == nil {
		psvsctn.chcdbf = iorw.NewBuffer()
	}
	return psvsctn.chcdbf
}

func processPsvSection(psvsctn *psvsection) (err error) {
	if psvsctn.nxtsctn == nil {
		prsng := psvsctn.prsng
		var rnrdr io.RuneReader = nil
		if tmpltpath := psvsctn.path(); tmpltpath != "" {
			var tmpltcoordsok bool = false
			if _, tmpltcoordsok = psvsctn.prsng.tmpltMap()[tmpltpath]; !tmpltcoordsok {
				if psvsctn.prsng.AtvActv != nil {
					lkprdr, lkperr := psvsctn.prsng.AtvActv.AltLookupTemplate(tmpltpath)
					if lkperr != nil {
						err = lkperr
					} else if lkprdr != nil {
						tmpltsi := psvsctn.prsng.tmpltBuf().Size()
						psvsctn.prsng.tmpltBuf().Print(lkprdr)
						if tmpltei := psvsctn.prsng.tmpltBuf().Size(); tmpltei > tmpltsi {
							psvsctn.prsng.tmpltMap()[tmpltpath] = []int64{tmpltsi, tmpltei}
							tmpltcoordsok = true
						}
					}
				}
			}
			if tmpltcoordsok {
				if tmplrdr, mxlen := psvsctn.prsng.tmpltrdr(tmpltpath); mxlen > 0 {
					rnrdr = tmplrdr
				}
			}
		}

		if elmtpe := psvsctn.elmtpe; elmtpe == ElemSingle || elmtpe == ElemStart {
			if elmtpe == ElemStart {
				if psvsctn.chcdbf != nil && psvsctn.chcdbf.Size() > 0 {
					if _, cntok := psvsctn.phrsmap["content"]; !cntok {
						stri := psvsctn.PhraseTemplateBuf().Size()
						if cntrdr := psvsctn.chcdbf.Reader(); cntrdr != nil {
							psvsctn.PhraseTemplateBuf().Print(cntrdr)
							cntrdr.Close()
							cntrdr = nil
							endi := psvsctn.PhraseTemplateBuf().Size()
							psvsctn.phrsmap["content"] = []int64{stri, endi}
						}
					}
					psvsctn.chcdbf.Clear()
				}
			}
			ParsePrsng(prsng, false, nil, "<@((...arguments)=>{@>")

			if rnrdr != nil {
				psvsctn.canphrs = true
				ParsePrsng(prsng, false, nil, rnrdr)
				psvsctn.canphrs = false
			}
			if psvsctn.tmpbuf != nil && psvsctn.tmpbuf.Size() > 0 {
				ParsePrsng(prsng, false, nil, "<@})(@>"+psvsctn.tmpbuf.String()+"<@);@>")
				psvsctn.tmpbuf.Clear()
			} else {
				ParsePrsng(prsng, false, nil, "<@})();@>")
			}

			decpsvcsection(psvsctn)
			if psvsctn.chcdbf != nil && psvsctn.chcdbf.Size() > 0 {
				rnrdr = psvsctn.chcdbf.Reader()
				psvsctn.canphrs = true
				ParsePrsng(prsng, false, nil, rnrdr)
				psvsctn.canphrs = false
				rnrdr = nil
			}

			/**/
		}
	} else {
		decpsvcsection(psvsctn)
	}
	psvsctn.dispose()
	psvsctn = nil
	return
}

func decpsvcsection(psvsctn *psvsection) {
	if psvsctn.prsng != nil {
		removePsvSection(psvsctn.prsng, psvsctn)
		psvsctn.prsng = nil
	}
}

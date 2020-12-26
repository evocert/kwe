package active

import (
	"bufio"
	"io"

	"github.com/evocert/kwe/iorw"
)

const (
	none int = iota
	gthan
	start
	end
	done
)

const (
	elemnone int = iota
	//ElemStart - elem start
	ElemStart
	//ElemEnd - elem end
	ElemEnd
	//ElemSingle - elem single
	ElemSingle
)

type passivectrl struct {
	prvpsvctrl   *passivectrl
	nxtpsvctrl   *passivectrl
	ctrlstage    int
	prsng        *parsing
	rawr         io.Reader
	bfrawr       *bufio.Reader
	tmpbuf       *iorw.Buffer
	cchdbuf      *iorw.Buffer
	prvrn        rune
	elmtype      int
	elmName      string
	lastElmType  int
	lastElemName string
	phrslbli     []int
	phrsprvr     rune
	phrsbuf      *iorw.Buffer
	cntntbuf     *iorw.Buffer
	prepped      bool
}

var phrslbl [][]rune = [][]rune{[]rune("{:"), []rune(":}")}

func newpsvctrl(prsng *parsing, prvpsvctrl *passivectrl) (psvctrl *passivectrl) {
	psvctrl = &passivectrl{prsng: prsng, prvrn: rune(0), prvpsvctrl: prvpsvctrl, phrslbli: []int{0, 0}, phrsprvr: rune(0)}
	if prvpsvctrl != nil {
		psvctrl.rawr = prvpsvctrl.rawr
		prvpsvctrl.rawr = nil
		if prsng.prvpsvctrls == nil {
			prsng.prvpsvctrls = map[*passivectrl]*passivectrl{}
		}
		prvpsvctrl.nxtpsvctrl = psvctrl
		prsng.prvpsvctrls[prvpsvctrl] = psvctrl
	}
	return
}

func (psvctrl *passivectrl) cachedbuf() *iorw.Buffer {
	if psvctrl.cchdbuf == nil {
		psvctrl.cchdbuf = iorw.NewBuffer()
	}
	return psvctrl.cchdbuf
}

func (psvctrl *passivectrl) phrasebuf() *iorw.Buffer {
	if psvctrl.phrsbuf == nil {
		psvctrl.phrsbuf = iorw.NewBuffer()
	}
	return psvctrl.phrsbuf
}

func (psvctrl *passivectrl) ReadRune() (r rune, size int, err error) {
	if psvctrl.rawr != nil {
		if psvctrl.bfrawr == nil {
			psvctrl.bfrawr = bufio.NewReader(psvctrl.rawr)
		}
		r, size, err = psvctrl.bfrawr.ReadRune()
	} else {
		err = io.EOF
	}
	return
}

func (psvctrl *passivectrl) clearcchdbuf() {
	if psvctrl.cchdbuf != nil {
		psvctrl.cchdbuf.Clear()
	}
}

func (psvctrl *passivectrl) reset() {
	if psvctrl.tmpbuf != nil {
		psvctrl.tmpbuf.Clear()
	}

	/*if psvctrl.tmpcde != nil {
		psvctrl.tmpcde.Clear()
	}*/
	if psvctrl.rawr != nil {
		psvctrl.rawr = nil
	}
	if psvctrl.bfrawr != nil {
		psvctrl.bfrawr = nil
	}
	if psvctrl.prepped {
		psvctrl.prepped = false
	}
	psvctrl.ctrlstage = none
	psvctrl.elmtype = elemnone
	psvctrl.prvrn = rune(0)
}

func (psvctrl *passivectrl) buf() *iorw.Buffer {
	if psvctrl.tmpbuf == nil {
		psvctrl.tmpbuf = iorw.NewBuffer()
	}
	return psvctrl.tmpbuf
}

func (psvctrl *passivectrl) bufsize() int64 {
	if psvctrl.tmpbuf == nil {
		return 0
	}
	return psvctrl.tmpbuf.Size()
}

func (psvctrl *passivectrl) close() {
	if psvctrl != nil {
		if psvctrl.prsng != nil {
			psvctrl.prsng = nil
		}
		if psvctrl.tmpbuf != nil {
			psvctrl.tmpbuf.Close()
			psvctrl.tmpbuf = nil
		}
		if psvctrl.cchdbuf != nil {
			psvctrl.cchdbuf.Close()
			psvctrl.cchdbuf = nil
		}
		/*if psvctrl.tmpcde != nil {
			psvctrl.tmpcde.Close()
			psvctrl.tmpcde = nil
		}*/
		if psvctrl.prvpsvctrl != nil {
			psvctrl.prvpsvctrl = nil
		}
		if psvctrl.nxtpsvctrl != nil {
			psvctrl.nxtpsvctrl = nil
		}
	}
}

func (psvctrl *passivectrl) validate() (valid bool) {
	valid = true
	return
}

func (psvctrl *passivectrl) outputrn(rn rune) (err error) {
	//err = psvctrl._outputrn(rn)
	err = parsepsvctrl(psvctrl, psvctrl.phrslbli, rn)
	return
}

func (psvctrl *passivectrl) _outputruns(p []rune) (err error) {
	if pl := len(p); pl > 0 {
		for _, rn := range p {
			if err = psvctrl._outputrn(rn); err != nil {
				break
			}
		}
	}
	return
}

func (psvctrl *passivectrl) _outputrn(rn rune) (err error) {
	psvctrl.prsng.psvr[psvctrl.prsng.psvri] = rn
	psvctrl.prsng.psvri++
	if psvctrl.prsng.psvri == len(psvctrl.prsng.psvr) {
		psvctrl.prsng.psvri = 0
		if psvctrl.lastElmType == ElemStart {
			err = psvctrl.cachedbuf().WriteRunes(psvctrl.prsng.psvr)
		} else {
			err = psvctrl.prsng.writePsv(psvctrl.prsng.psvr)
		}
	}
	return
}

func processPhrase(psvctrl *passivectrl, phrsbuf *iorw.Buffer) (err error) {
	if bufs := phrsbuf.String(); bufs != "" {
		if bufs == "content" {
			if psvctrl.cntntbuf != nil && psvctrl.cntntbuf.Size() > 0 {
				func() {
					if cntntr := psvctrl.cntntbuf.Reader(); cntntr != nil {
						defer cntntr.Close()
						err = parseprsngrunerdr(psvctrl.prsng, cntntr, false)
					}
				}()
			}
		}
	}
	return
}

func parsepsvctrl(psvctrl *passivectrl, phrslbli []int, pr rune) (err error) {
	if phrslbli[1] == 0 && phrslbli[0] < len(phrslbl[0]) {
		if phrslbli[0] > 0 && phrslbl[0][phrslbli[0]-1] == psvctrl.phrsprvr && phrslbl[0][phrslbli[0]] != pr {
			if phrsi := phrslbli[0]; phrsi > 0 {
				phrslbli[0] = 0
				err = psvctrl._outputruns(phrslbl[0][:phrsi])
			}
		}
		if phrslbl[0][phrslbli[0]] == pr {
			phrslbli[0]++
			if phrslbli[0] == len(phrslbl[0]) {

				psvctrl.phrsprvr = 0
			}
			psvctrl.phrsprvr = pr
		} else {
			if phrsi := phrslbli[0]; phrsi > 0 {
				phrslbli[0] = 0
				err = psvctrl._outputruns(phrslbl[0][:phrsi])
			}
			psvctrl.phrsprvr = pr
			err = psvctrl._outputrn(pr)
		}
	} else if phrslbli[0] == len(phrslbl[0]) && phrslbli[1] < len(phrslbl[1]) {
		if phrslbl[1][phrslbli[1]] == pr {
			phrslbli[1]++
			if phrslbli[1] == len(phrslbl[1]) {

				phrslbli[0] = 0
				phrslbli[1] = 0
				psvctrl.phrsprvr = 0
				if psvctrl.phrsbuf != nil {
					if psvctrl.phrsbuf.Size() > 0 {
						err = processPhrase(psvctrl, psvctrl.phrsbuf)
					}
					psvctrl.phrsbuf.Clear()
				}
			}
		} else {
			if phrsi := phrslbli[1]; phrsi > 0 {
				phrslbli[1] = 0
				err = psvctrl.phrasebuf().WriteRunes(phrslbl[1][:phrsi])
			}
			psvctrl.phrsprvr = pr
			err = psvctrl.phrasebuf().WriteRune(pr)
		}
	}
	return
}

func (psvctrl *passivectrl) flushout(rns ...rune) (err error) {
	psvctrl.ctrlstage = none
	if psvctrl.elmtype != elemnone {
		psvctrl.outputrn('<')
		if psvctrl.elmtype == ElemEnd {
			psvctrl.outputrn('/')
		}
	}
	if psvctrl.bufsize() > 0 {
		for _, rn := range []rune(psvctrl.buf().String()) {
			if err = psvctrl.outputrn(rn); err != nil {
				break
			}
		}
	}
	if psvctrl.elmtype == ElemSingle {
		psvctrl.outputrn('/')
	}
	if len(rns) > 0 {
		for _, rn := range rns {
			if err = psvctrl.outputrn(rn); err != nil {
				break
			}
		}
	}
	psvctrl.reset()
	return err
}

func (psvctrl *passivectrl) validrune(rn rune) bool {
	//  (A-Z) || (a-z) || (0-9) || ('.','/','|')
	return (rn >= 65 && rn <= 90) || (rn >= 97 && rn <= 122) || (rn >= 30 && rn <= 39) || (rn == '.' || rn == '/' || rn == '|')
}

func (psvctrl *passivectrl) resetphrase() {
	psvctrl.phrslbli[0] = 0
	psvctrl.phrslbli[1] = 0
	if psvctrl.phrsbuf != nil && psvctrl.phrsbuf.Size() > 0 {
		psvctrl.phrsbuf.Clear()
	}
}

func (psvctrl *passivectrl) processrn(rn rune) (err error) {
	if psvctrl.ctrlstage == none {
		if rn == '<' {
			psvctrl.resetphrase()
			psvctrl.ctrlstage = gthan
			psvctrl.prvrn = rn
			psvctrl.elmtype = ElemStart
		} else {
			err = psvctrl.outputrn(rn)
		}
	} else if psvctrl.ctrlstage == gthan {
		if rn == ':' && (psvctrl.prvrn == '/' || psvctrl.prvrn == '<') {
			psvctrl.ctrlstage = start
			if psvctrl.elmtype == ElemEnd {
				psvctrl.buf().Print(string('/'))
			}
			psvctrl.buf().Print(string(rn))
			psvctrl.prvrn = 0
			psvctrl.prvrn = rn
		} else if rn == '/' && psvctrl.prvrn == '<' {
			if psvctrl.elmtype == ElemStart {
				psvctrl.elmtype = ElemEnd
				psvctrl.prvrn = rn
			} else {
				psvctrl.flushout(rn)
				//Flush
			}
		} else {
			psvctrl.flushout(rn)
			//Flush
		}
	} else if psvctrl.ctrlstage == start {
		if rn == '>' {
			if psvctrl.prvrn == '/' {
				if psvctrl.elmtype == ElemStart {
					psvctrl.elmtype = ElemSingle
				}
			}
			if psvctrl.validate() {
				if err = parseValidity(psvctrl); err == nil {
					psvctrl.reset()
				}
			} else {
				psvctrl.flushout(rn)
			}
		} else if psvctrl.validrune(rn) {
			if rn == '/' {
				if psvctrl.elmtype == ElemEnd {
					err = psvctrl.flushout(rn)
				} else {
					psvctrl.prvrn = rn
				}
			} else {
				psvctrl.buf().Print(string(rn))
				psvctrl.prvrn = rn
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
	if prsng.psvctrl == nil {
		prsng.psvctrl = newpsvctrl(prsng, nil)
	}
	err = prsng.psvctrl.processrn(rn)
	return
}

func parseValidity(psvctrl *passivectrl) (err error) {
	var elmtype = psvctrl.elmtype
	var elmname = psvctrl.buf().String()
	var elmpath = elmname[1:]
	psvctrl.reset()
	if elmname != "" {
		if elmtype == ElemStart || elmtype == ElemSingle {
			if psvctrl.prsng.atv.LookupTemplate != nil {
				if rawr, rawrerr := psvctrl.prsng.atv.LookupTemplate(elmpath); rawr != nil && rawrerr == nil {
					psvctrl.rawr = rawr
					psvctrl.prsng.flushPsv()
					if elmtype == ElemSingle {
						err = parseprsngrunerdr(psvctrl.prsng, psvctrl, false)
					} else {
						if psvctrl.lastElmType == elemnone && psvctrl.lastElemName == "" {
							psvctrl.lastElmType = elmtype
							psvctrl.lastElemName = elmname
							newpsvctrl(psvctrl.prsng, psvctrl)
						}
					}
				} else if rawrerr != nil {
					err = rawrerr
				}
			}
			//}
		} else if elmtype == ElemEnd {
			if psvctrl.nxtpsvctrl != nil {
				if psvctrl.lastElemName != "" && psvctrl.lastElmType == ElemStart {
					psvctrl.prsng.flushPsv()
					psvctrl.lastElemName = ""
					psvctrl.lastElmType = elemnone
					if nxtpsvctrl := psvctrl.nxtpsvctrl; nxtpsvctrl != nil {
						if psvctrl.cchdbuf != nil && psvctrl.cchdbuf.Size() > 0 {
							if nxtpsvctrl.cntntbuf != nil {
								nxtpsvctrl.cntntbuf.Close()
								nxtpsvctrl.cntntbuf = nil
							}
							nxtpsvctrl.cntntbuf = psvctrl.cchdbuf
							psvctrl.cchdbuf = nil
						}
						psvctrl.prsng.psvctrl = psvctrl.nxtpsvctrl
						if err = parseprsngrunerdr(psvctrl.prsng, psvctrl.prsng.psvctrl, false); err == nil || err == io.EOF {
							psvctrl.prsng.psvctrl = nxtpsvctrl.prvpsvctrl
						}
						psvctrl.prsng.psvctrl.reset()
						nxtpsvctrl.close()
					}
				}
			}
		}
	}
	return
}

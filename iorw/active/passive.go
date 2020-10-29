package active

import (
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

type passivecontrol struct {
	ctrlstage int
	prsng     *parsing
	tmpbuf    *iorw.Buffer
	prvrn     rune
	elmtype   int
	tmpcde    *iorw.Buffer
}

func newpsvctrl(prsng *parsing) (psvctrl *passivecontrol) {
	psvctrl = &passivecontrol{prsng: prsng, prvrn: rune(0)}
	return
}

func (psvctrl *passivecontrol) reset() {
	if psvctrl.tmpbuf != nil {
		psvctrl.tmpbuf.Clear()
	}
	if psvctrl.tmpcde != nil {
		psvctrl.tmpcde.Clear()
	}
	psvctrl.ctrlstage = none
	psvctrl.elmtype = elemnone
	psvctrl.prvrn = rune(0)
}

func (psvctrl *passivecontrol) buf() *iorw.Buffer {
	if psvctrl.tmpbuf == nil {
		psvctrl.tmpbuf = iorw.NewBuffer()
	}
	return psvctrl.tmpbuf
}

func (psvctrl *passivecontrol) bufsize() int64 {
	if psvctrl.tmpbuf == nil {
		return 0
	}
	return psvctrl.tmpbuf.Size()
}

func (psvctrl *passivecontrol) close() {
	if psvctrl != nil {
		if psvctrl.prsng != nil {
			psvctrl.prsng = nil
		}
		if psvctrl.tmpbuf != nil {
			psvctrl.tmpbuf.Close()
			psvctrl.tmpbuf = nil
		}
		if psvctrl.tmpcde != nil {
			psvctrl.tmpcde.Close()
			psvctrl.tmpcde = nil
		}
	}
}

func (psvctrl *passivecontrol) validate() (valid bool) {
	valid = true
	return
}

func (psvctrl *passivecontrol) outputrn(rn rune) (err error) {
	psvctrl.prsng.psvr[psvctrl.prsng.psvri] = rn
	psvctrl.prsng.psvri++
	if psvctrl.prsng.psvri == len(psvctrl.prsng.psvr) {
		psvctrl.prsng.psvri = 0
		err = psvctrl.prsng.writePsv(psvctrl.prsng.psvr)
	}
	return
}

func (psvctrl *passivecontrol) flushout(rns ...rune) (err error) {
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

func (psvctrl *passivecontrol) validrune(rn rune) bool {
	//  (A-Z) || (a-z) || (0-9) || ('.','/','|')
	return (rn >= 65 && rn <= 90) || (rn >= 97 && rn <= 122) || (rn >= 30 && rn <= 39) || (rn == '.' || rn == '/' || rn == '|')
}

func (psvctrl *passivecontrol) processrn(rn rune) (err error) {
	if psvctrl.ctrlstage == none {
		if rn == '<' {
			psvctrl.ctrlstage = gthan
			psvctrl.prvrn = rn
			psvctrl.elmtype = ElemStart
		} else {
			err = psvctrl.outputrn(rn)
		}
	} else if psvctrl.ctrlstage == gthan {
		if rn == '#' && (psvctrl.prvrn == '/' || psvctrl.prvrn == '<') {
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
				if err = psvctrl.prsng.prepPsvValidity(psvctrl); err == nil {
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
		prsng.psvctrl = newpsvctrl(prsng)
	}
	//err = prsng.psvctrl.outputrn(rn)
	err = prsng.psvctrl.processrn(rn)
	return
}

func (prsng *parsing) prepPsvValidity(psvctrl *passivecontrol) (err error) {
	var elmtype = psvctrl.elmtype
	psvctrl.reset()
	if elmtype == ElemStart || elmtype == ElemSingle {

	} else if elmtype == ElemEnd {

	}
	return
}

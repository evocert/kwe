package active

import (
	"github.com/dop251/goja"
	"github.com/evocert/kwe/iorw"
)

//Active - struct
type Active struct {
	*iorw.Buffer
	atvrntmes []*atvruntime
}

//NewActive - instance
func NewActive() (atv *Active) {
	atv = &Active{Buffer: iorw.NewBuffer(), atvrntmes: []*atvruntime{}}
	return
}

//Close - refer to  io.Closer
func (atv *Active) Close() (err error) {
	if atv.atvrntmes != nil {
		for len(atv.atvrntmes) > 0 {
			atv.atvrntmes[len(atv.atvrntmes)-1].close()
			atv.atvrntmes[len(atv.atvrntmes)-1] = nil
			atv.atvrntmes = atv.atvrntmes[0 : len(atv.atvrntmes)-1]
		}
	}
	return
}

type parsing struct {
	*iorw.Buffer
	atv       *Active
	prslbl    [][]rune
	prslbli   []int
	prslblprv []rune
	prntprsng *parsing
	foundcde  bool
	hascde    bool
	cdemap    map[int][]int64
	psvmap    map[int][]int64
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
		prsng = nil
	}
}

func (prsng *parsing) topprsng() *parsing {
	if prsng.prntprsng == nil {
		return prsng
	}
	return prsng.prntprsng.topprsng()
}

func parsepsvrunes(prsng *parsing, p []rune) {

}

func parseprsng(prsng *parsing, prslbl [][]rune, prslbli []int, prslblprv []rune, pr rune) {
	if prslbli[1] == 0 && prslbl[0][prslbli[0]] != pr {
		if prslbli[0] > 0 && prslbl[0][prslbli[0]-1] == prslblprv[0] && prslbl[0][prslbli[0]] != pr {

		}
		if prslbl[0][prslbli[0]] == pr {
			prslbli[0]++
			if prslbli[0] == len(prslbl[0]) {

				prslblprv[0] = 0
			} else {
				prslblprv[0] = pr
			}
		} else {

		}
	} else if prslbli[0] == len(prslbl[0]) && prslbl[1][prslbli[1]] != pr {
		if prslbl[1][prslbli[1]] == pr {
			prslbli[1]++
			if prslbli[1] == len(prslbl[1]) {

			}
		} else {
			if prslbli[1] > 0 {

			}
		}
	}
}

func nextparsing(atv *Active, prntprsng *parsing) (prsng *parsing) {
	prsng = &parsing{prntprsng: prntprsng, atv: atv, prslbl: [][]rune{[]rune("<@"), []rune("@>")}, prslbli: []int{0, 0}, prslblprv: []rune{0, 0}}
	return
}

type atvruntime struct {
	*parsing
	atv *Active
	vm  *goja.Runtime
}

func (atvrntme *atvruntime) close() {
	if atvrntme != nil {
		if atvrntme.parsing != nil {
			atvrntme.parsing.close()
			atvrntme.parsing = nil
		}
		if atvrntme.atv != nil {
			atvrntme.atv = nil
		}
		if atvrntme.vm != nil {
			atvrntme.vm = nil
		}
	}
}

func newatvruntime(atv *Active) (atvrntme *atvruntime) {
	atvrntme = &atvruntime{atv: atv, parsing: nextparsing(atv, nil), vm: goja.New()}

	return
}

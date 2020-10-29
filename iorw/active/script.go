package active

import "strings"

func parseatvrune(prsng *parsing, rn rune) (err error) {
	if !prsng.hascde {
		prsng.flushPsv()
		if strings.TrimSpace(string(rn)) != "" {
			prsng.hascde = true
		}
	}
	if prsng.hascde {
		if !prsng.foundcde {
			prsng.foundcde = true
		}
		prsng.cder[prsng.cderi] = rn
		prsng.cderi++
		if prsng.cderi == len(prsng.cder) {
			prsng.cderi = 0
			err = prsng.writeCde(prsng.cder)
		}
	}
	return
}

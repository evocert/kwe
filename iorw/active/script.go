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
		if prsng.cdetxt == rune(0) {
			if (rn == '\'' || rn == '"') && prsng.prslblprv[1] != '\\' {
				prsng.cdetxt = rn
			}
		} else if prsng.cdetxt > rune(0) && prsng.cdetxt == rn {
			if (rn == '\'' || rn == '"') && prsng.prslblprv[1] != '\\' {
				prsng.cdetxt = rune(0)
			}
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

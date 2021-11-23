package active

import (
	"strings"
)

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
			if (rn == '\'' || rn == '"' || rn == '`') && prsng.prslblprv[1] != '\\' {
				prsng.cdetxt = rn
			}
		} else if prsng.cdetxt > rune(0) && prsng.cdetxt == rn {
			if (rn == '\'' || rn == '"' || rn == '`') && prsng.prslblprv[1] != '\\' {
				prsng.cdetxt = rune(0)
			}
		}
		if prsng.cdetxt == '`' {
			if (rn == '\'' || rn == '"' || rn == '\\' || rn == '\t' || rn == '\r' || rn == '\n') && prsng.prslblprv[1] != '\\' {
				if rn == '\t' {
					rn = 't'
				} else if rn == '\r' {
					rn = 'r'
				} else if rn == '\n' {
					rn = 'n'
				}
				for _, sbcdtxtr := range []rune{'\\', rn} {
					prsng.cder[prsng.cderi] = sbcdtxtr
					prsng.cderi++
					if prsng.cderi == len(prsng.cder) {
						prsng.cderi = 0
						err = prsng.writeCde(prsng.cder)
					}
				}
			} else {
				prsng.cder[prsng.cderi] = rn
				prsng.cderi++
				if prsng.cderi == len(prsng.cder) {
					prsng.cderi = 0
					err = prsng.writeCde(prsng.cder)
				}
			}
		} else {
			prsng.cder[prsng.cderi] = rn
			prsng.cderi++
			if prsng.cderi == len(prsng.cder) {
				prsng.cderi = 0
				err = prsng.writeCde(prsng.cder)
			}
		}
	}
	return
}

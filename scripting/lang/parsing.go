package lang

import (
	"bufio"
	"context"
	"io"

	"github.com/evocert/kwe/iorw"
)

type Parser struct {
}

func NewParser() (prsr *Parser) {
	prsr = &Parser{}

	return
}

func (prsr *Parser) Parse(lang string, a ...interface{}) (prgrm *Program, err error) {
	if len(a) > 0 {
		prgrm = &Program{}
		var pi, pw = io.Pipe()
		go func() {
			var pwerr error = nil
			defer func() {
				if pwerr != nil {
					pw.CloseWithError(pwerr)
				} else {
					pw.Close()
				}
			}()
			pwerr = iorw.Fprint(pw, a...)
		}()

		var ctx, ctxcnl = context.WithCancel(context.Background())
		go func() {
			var pierr error = nil
			defer func() {
				if pierr != nil {
					err = pierr
				}
				defer ctxcnl()
			}()
			var bufpi = bufio.NewReader(pi)
			var tmps = ""
			var lngkwrdmp = LangKeyWordMap(lang)
			for pierr == nil {
				r, size, rerr := bufpi.ReadRune()
				if size > 0 {
					if r > 0 {
						tmps += (string(r) + "")
						if tmpkwrd, _ := lngkwrdmp[tmps]; tmpkwrd != nil {
							if !tmpkwrd.FutureKeyword {

							}
						}
					}
				}
				if rerr != nil {
					if rerr != io.EOF {
						pierr = rerr
					} else {
						break
					}
				}
			}
		}()
		<-ctx.Done()
	}
	return
}

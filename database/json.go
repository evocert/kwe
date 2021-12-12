package database

import (
	"bufio"
	"encoding/json"
	"io"

	"context"

	"github.com/evocert/kwe/iorw"
)

//JSONReader - struct
type JSONReader struct {
	rdr   *Reader
	bfr   *bufio.Reader
	exctr *Executor
	pr    *io.PipeReader
	pw    *io.PipeWriter
	err   error
}

//NewJSONReader - over rdr*Reader or exctr*Executor
func NewJSONReader(rdr *Reader, exctr *Executor, err error) (jsnr *JSONReader) {
	jsnr = &JSONReader{rdr: rdr, exctr: exctr, err: err}
	return
}

//Read - refer to io.Reader
func (jsnr *JSONReader) Read(p []byte) (n int, err error) {
	if jsnr.pr == nil && jsnr.pw == nil {
		jsnr.pr, jsnr.pw = io.Pipe()
		ctx, ctxcancel := context.WithCancel(context.Background())
		go func(rdr *Reader, exctr *Executor) {
			defer func() {
				jsnr.pw.Close()
			}()
			enc := json.NewEncoder(jsnr.pw)
			ctxcancel()
			if rdr != nil {
				iorw.Fprint(jsnr.pw, "{\"columns\":[")
				for cn, c := range rdr.cls {
					iorw.Fprint(jsnr.pw, "{")
					t := rdr.cltpes[cn]
					var ctpm = map[string]interface{}{"title": c, "name": c, "dbtype": t.DatabaseType(), "type": t.Type().Name(), "length": t.Length(), "numeric": t.Numeric(), "scale": t.Scale(), "precision": t.Precision()}
					var ctpml = len(ctpm)
					for ctpmk, ctpmv := range ctpm {
						ctpml--
						enc.Encode(ctpmk)
						iorw.Fprint(jsnr.pw, ":")
						enc.Encode(ctpmv)
						if ctpml > 0 {
							iorw.Fprint(jsnr.pw, ",")
						}
					}
					iorw.Fprint(jsnr.pw, "}")
					if cn < len(rdr.cls)-1 {
						iorw.Fprint(jsnr.pw, ",")
					}
				}
				iorw.Fprint(jsnr.pw, "]")
				iorw.Fprint(jsnr.pw, ",\"data\":")
				iorw.Fprint(jsnr.pw, "[")
				if nxt, nxterr := rdr.Next(); nxt {
					var nxtdata []interface{} = rdr.Data()
					for {
						iorw.Fprint(jsnr.pw, "[")
						for nd, d := range nxtdata {
							enc.Encode(d)
							if nd < len(nxtdata)-1 {
								iorw.Fprint(jsnr.pw, ",")
							}
						}
						iorw.Fprint(jsnr.pw, "]")
						nxt, nxterr = rdr.Next()
						if nxt {
							iorw.Fprint(jsnr.pw, ",")
							nxtdata = rdr.Data()
						} else if !nxt || nxterr != nil {
							break
						}
					}
				}
				iorw.Fprint(jsnr.pw, "]")
				iorw.Fprint(jsnr.pw, "}")
			} else if exctr != nil {
				if jsnr.err == nil && exctr.lasterr == nil {
					iorw.Fprint(jsnr.pw, "{")
					iorw.Fprint(jsnr.pw, "}")
				} else if jsnr.err != nil || exctr.lasterr != nil {
					iorw.Fprint(jsnr.pw, "{\"error\":")
					if jsnr.err != nil {
						enc.Encode(jsnr.err.Error())
					} else if exctr.lasterr != nil {
						enc.Encode(exctr.lasterr.Error())
					}
					iorw.Fprint(jsnr.pw, "}")
				}
			} else {
				iorw.Fprint(jsnr.pw, "{\"error\":")
				if jsnr.err != nil {
					enc.Encode(jsnr.err.Error())
				} else {
					enc.Encode("empty")
				}
				iorw.Fprint(jsnr.pw, "}")
			}
		}(jsnr.rdr, jsnr.exctr)
		<-ctx.Done()
		ctx = nil
	}
	if jsnr.pr != nil {
		n, err = jsnr.pr.Read(p)
	}
	if n == 0 && err == nil {
		err = io.EOF
	}
	return
}

//ReadRune - refer to io.RuneReader
func (jsnr *JSONReader) ReadRune() (r rune, size int, err error) {
	if jsnr.bfr == nil {
		jsnr.bfr = bufio.NewReader(jsnr)
	}
	r, size, err = jsnr.bfr.ReadRune()
	return
}

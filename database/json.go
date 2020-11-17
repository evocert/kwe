package database

import (
	"bufio"
	"encoding/json"
	"io"
	"sync"

	"github.com/evocert/kwe/iorw"
)

//JSONReader - struct
type JSONReader struct {
	rdr   *Reader
	bfr   *bufio.Reader
	exctr *Executor
	pr    *io.PipeReader
	pw    *io.PipeWriter
}

//NewJSONReader - over rdr*Reader or exctr*Executor
func NewJSONReader(rdr *Reader, exctr *Executor) (jsnr *JSONReader) {
	jsnr = &JSONReader{rdr: rdr, exctr: exctr}
	return
}

//Read - refer to io.Reader
func (jsnr *JSONReader) Read(p []byte) (n int, err error) {
	if jsnr.pr == nil && jsnr.pw == nil {
		jsnr.pr, jsnr.pw = io.Pipe()
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func(rdr *Reader, exctr *Executor) {
			defer func() {
				jsnr.pw.Close()
			}()
			if rdr != nil {
				enc := json.NewEncoder(jsnr.pw)
				wg.Done()
				iorw.Fprint(jsnr.pw, "{columns:[")
				for cn, c := range rdr.cls {
					iorw.Fprint(jsnr.pw, "{")
					t := rdr.cltpes[cn]
					var ctpm = map[string]interface{}{"name": c, "dbtype": t.DatabaseType(), "type": t.Type().Name(), "length": t.Length(), "numeric": t.Numeric(), "scale": t.Scale(), "precision": t.Precision()}
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
				iorw.Fprint(jsnr.pw, "]}")
				iorw.Fprint(jsnr.pw, ",data:")
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
			} else {
				wg.Done()
			}
		}(jsnr.rdr, jsnr.exctr)
		wg.Wait()
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

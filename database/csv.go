package database

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/evocert/kwe/iorw"
)

//CSVReader -
type CSVReader struct {
	rdr        *Reader
	bfr        *bufio.Reader
	err        error
	Headers    bool
	ColDelim   string
	RowDelim   string
	IncludeEOF bool
	AltHeaders []string
	pr         *io.PipeReader
	pw         *io.PipeWriter
}

//NewCSVReader - over rdr*Reader
func NewCSVReader(rdr *Reader, err error, a ...interface{}) (csvr *CSVReader) {
	csvr = &CSVReader{rdr: rdr, err: err, Headers: true, ColDelim: ",", RowDelim: "\r\n", AltHeaders: nil}
	for len(a) > 0 {
		var d = a[0]
		if ctngs, ctngsok := d.(map[string]interface{}); ctngsok {
			for ctngk := range ctngs {
				ctngv := ctngs[ctngk]
				if strings.ToLower(ctngk) == "coldelim" {
					if cldelim, cldelimok := ctngv.(string); cldelimok {
						if cldelim != "" {
							csvr.ColDelim = cldelim
						}
					}
				} else if strings.ToLower(ctngk) == "headers" {
					if hdrs, hdrsok := ctngv.(bool); hdrsok {
						csvr.Headers = hdrs
					}
				} else if strings.ToLower(ctngk) == "altheaders" {
					if althdrs, althdrsok := ctngv.([]interface{}); althdrsok {
						althrsds := make([]string, len(althdrs))
						if len(althrsds) > 0 {
							for n := range althdrs {
								althrsds[n], _ = althdrs[n].(string)
							}
							csvr.AltHeaders = althrsds
						}
					}
				}
			}
		}
		a = a[1:]
	}
	return
}

//Read - refer to io.Reader
func (csvr *CSVReader) Read(p []byte) (n int, err error) {
	if csvr.pr == nil && csvr.pw == nil {
		csvr.pr, csvr.pw = io.Pipe()
		ctx, ctxcancel := context.WithCancel(context.Background())
		go func(rdr *Reader) {
			defer func() {
				csvr.pw.Close()
			}()
			//enc := json.NewEncoder(jsnr.pw)
			ctxcancel()
			if csvr.err == nil {
				if rdr != nil {
					var sval = func(v string) (s string) {
						if v != "" {
							if strings.Index(v, csvr.ColDelim) >= 0 || strings.Index(v, "\"") >= 0 {
								s = "\"" + strings.Replace(v, "\"", "\"\"", -1) + "\""
							} else {
								s = v
							}
						}
						return
					}
					var canPrintCols = false
					var colcount = 0
					var cls []string = nil
					if csvr.Headers && len(csvr.AltHeaders) == 0 {
						cls = rdr.cls[:]
						canPrintCols = true
					} else if len(csvr.AltHeaders) > 0 {
						cls = csvr.AltHeaders[:]
						canPrintCols = true
					}
					colcount = len(cls)
					var dta []interface{} = nil
					if nxt, nxterr := rdr.Next(); nxterr == nil {
						if nxt {
							dta = rdr.Data()
							if colcount == 0 && len(dta) > 0 {
								colcount = len(dta)
								cls = make([]string, colcount)
								for n := range cls {
									cls[n] = "Column" + fmt.Sprintf("%d", n)
								}
							}
						}
						if canPrintCols && len(cls) > 0 {
							for n := range cls {
								c := cls[n]
								iorw.Fprint(csvr.pw, sval(strings.ToUpper(c)))
								if n < colcount-1 {
									iorw.Fprint(csvr.pw, csvr.ColDelim)
								}
							}
						}
						var firstRow = true
						for {
							if nxt && len(dta) == colcount {
								if firstRow {
									firstRow = false
									if canPrintCols {
										iorw.Fprint(csvr.pw, csvr.RowDelim)
									}
								} else {
									iorw.Fprint(csvr.pw, csvr.RowDelim)
								}
								for n := range dta {
									d := dta[n]
									if s, sok := d.(string); sok {
										if fltval, nrerr := strconv.ParseFloat(s, 64); nrerr == nil {
											if tstintval := int64(fltval); float64(tstintval) == fltval {
												iorw.Fprint(csvr.pw, fmt.Sprintf("%d", tstintval))
											} else {
												iorw.Fprint(csvr.pw, fmt.Sprintf("%.0f", fltval))
											}
										} else if intval, nrerr := strconv.ParseInt(s, 10, 64); nrerr == nil {
											iorw.Fprint(csvr.pw, fmt.Sprintf("%d", intval))
										} else {
											if _, terr := time.Parse("2006-01-02T15:04:05", s); terr == nil {
												iorw.Fprint(csvr.pw, strings.Replace(s, "T", " ", -1))
											} else {
												iorw.Fprint(csvr.pw, sval(s))
											}
										}
									} else if t, tok := d.(time.Time); tok {
										iorw.Fprint(csvr.pw, t)
									} else {
										iorw.Fprint(csvr.pw, d)
									}
									if n < colcount-1 {
										iorw.Fprint(csvr.pw, csvr.ColDelim)
									}
								}
								nxt, nxterr = rdr.Next()
								if nxterr != nil {
									nxt = false
								} else if nxt {
									dta = rdr.Data()
								}
							} else {
								break
							}
						}
						if csvr.IncludeEOF {
							iorw.Fprint(csvr.pw, csvr.RowDelim)
						}
					}
				}
			}
		}(csvr.rdr)
		<-ctx.Done()
		ctx = nil
	}
	if csvr.pr != nil {
		n, err = csvr.pr.Read(p)
	}
	if n == 0 && err == nil {
		err = io.EOF
	}
	return
}

//ReadRune - refer to io.RuneReader
func (csvr *CSVReader) ReadRune() (r rune, size int, err error) {
	if csvr.bfr == nil {
		csvr.bfr = bufio.NewReader(csvr)
	}
	r, size, err = csvr.bfr.ReadRune()
	return
}

package chnls

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/evocert/kwe/iorw"
)

type requeststdio struct {
	wg     *sync.WaitGroup
	rqst   *Request
	inbuf  *iorw.Buffer
	outbuf *iorw.Buffer
}

func (rqststdio *requeststdio) captureRunes(p ...rune) (err error) {
	if len(p) > 0 {
		err = rqststdio.inbuf.WriteRunes(p...)
	}
	return
}

func (rqststdio *requeststdio) executeStdIO() {
	rqststdio.wg.Add(1)
	go func() {
		rns := make([]rune, 1024)
		rnsi := 0
		var rnserr error = nil
		if stdio, stdiook := rqststdio.rqst.rqstr.(*os.File); stdiook {
			rdr := bufio.NewScanner(stdio)
			firstScan := true
			for {
				if firstScan {
					firstScan = false
					fmt.Print("")
					rdr.Scan()
				} else {
					rdr.Scan()
				}
				if text := rdr.Text(); text != "" {
					if strings.HasPrefix(text, "!!js:") {
						text = text[len("!!js:"):]
						if rqststdio.inbuf.Size() > 0 {
							bfr := rqststdio.inbuf.Reader()
							if text != "" {
								if filepath.Ext(text) == "" {
									text = text + ".js"
								}
								rqststdio.rqst.MapResource(text, bfr)
								rqststdio.rqst.AddPath(text)
								rqststdio.rqst.processPaths(false)
							} else {
								rqststdio.rqst.copy(bfr, rqststdio.rqst, true)
								bfr.Close()
							}
							rqststdio.inbuf.Clear()
						}
					} else {
						for _, r := range rdr.Text() {
							rns[rnsi] = r
							rnsi++
							if rnsi == len(rns) {
								rnserr = rqststdio.captureRunes(rns[:rnsi]...)
								rnsi = 0
								if rnserr != nil {
									break
								}
							}
						}

						if rnsi > 0 {
							rnserr = rqststdio.captureRunes(rns[:rnsi]...)
							rnsi = 0
							if rnserr != nil {
								break
							}
						}
					}
				}
			}
		} else {
			rdr := bufio.NewReader(rqststdio.rqst.rqstr)
			for {
				r, s, rerr := rdr.ReadRune()
				if s > 0 {
					rns[rnsi] = r
					rnsi++
					if rnsi == len(rns) {
						rqststdio.captureRunes(rns[:rnsi]...)
						rnsi = 0
					}
				}
				if rerr != nil {
					if rerr == io.EOF {
						if rnsi > 0 {
							rqststdio.captureRunes(rns[:rnsi]...)
							rnsi = 0
						}
						time.Sleep(10)
					}
				} else {
					break
				}
			}
		}
	}()
	rqststdio.wg.Wait()
}

func (rqststdio *requeststdio) Print(a ...interface{}) {
	iorw.Fprint(rqststdio, a...)
}

func (rqststdio *requeststdio) Println(a ...interface{}) {
	iorw.Fprintln(rqststdio, a...)
}

func newrequeststdio(rqst *Request) (rqststdio *requeststdio) {
	rqststdio = &requeststdio{rqst: rqst, wg: &sync.WaitGroup{}, inbuf: iorw.NewBuffer(), outbuf: iorw.NewBuffer()}
	return
}

func (rqststdio *requeststdio) dispose() {
	if rqststdio != nil {
		if rqststdio.inbuf != nil {
			rqststdio.inbuf.Close()
			rqststdio.inbuf = nil
		}
		if rqststdio.outbuf != nil {
			rqststdio.outbuf.Close()
			rqststdio.outbuf = nil
		}
		rqststdio = nil
	}
}

func (rqststdio *requeststdio) Write(p []byte) (n int, err error) {
	if rqststdio.rqst != nil {
		n, err = rqststdio.rqst.Write(p)
	}
	return
}

package chnls

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/evocert/kwe/database"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/osprc"
	"github.com/evocert/kwe/ws"
)

type requeststdio struct {
	wg            *sync.WaitGroup
	rqst          *Request
	isDone        bool
	prvinr        rune
	hashbang      bool
	lsthshlnk     string
	lsthshlnkargs []string
	tmpbuf        *iorw.Buffer
	inbuf         *iorw.Buffer
	outbuf        *iorw.Buffer
	//
	cmd *osprc.Command
}

func (rqststdio *requeststdio) captureRunes(eof bool, p ...rune) (err error) {
	if len(p) > 0 {
		for _, r := range p {
			if err = rqststdio.captureRune(r); err != nil {
				break
			}
		}
	}
	return
}

func (rqststdio *requeststdio) captureRune(r rune) (err error) {
	if r == rune(10) {
		if rqststdio.hashbang {
			if s := rqststdio.tmpbuf.String(); s != "" {
				rqststdio.tmpbuf.Clear()
				sargs := strings.Split(s, " ")
				sargsi := 0
				for sargsi < len(sargs) {
					if strings.TrimSpace(sargs[sargsi]) == "" {
						sargs = append(sargs[:sargsi], sargs[sargsi+1:]...)
					} else {
						sargs[sargsi] = strings.TrimSpace(sargs[sargsi])
						sargsi++
					}
				}
				if sargs[0] == "#!commit" || sargs[0] == "#!close" || sargs[0] == "#!exit" {
					if sargs[0] == "#!exit" {
						if rqststdio.cmd != nil {
							rqststdio.rqst.copy(rqststdio.cmd, rqststdio.rqst, false, "")
							rqststdio.cmd.Close()
							rqststdio.cmd = nil
						}
					} else {
						if rqststdio.inbuf.Size() > 0 {
							rqststdio.tmpbuf.Clear()
						}
						if rqststdio.lsthshlnk == "#!js" {
							if rqststdio.inbuf.Size() > 0 {
								if bfr := rqststdio.inbuf.Reader(); bfr != nil {
									rqststdio.rqst.copy(bfr, rqststdio.rqst, true, "")
									bfr.Close()
								}
								rqststdio.inbuf.Clear()
							}
						} else if rqststdio.lsthshlnk == "#!dbms" {
							if rqststdio.inbuf.Size() > 0 {
								if bfr := rqststdio.inbuf.Reader(); bfr != nil {
									if dbmserr := database.GLOBALDBMS().InOut(bfr, rqststdio.rqst, rqststdio.rqst.Parameters()); dbmserr != nil {
										rqststdio.isDone = true
										err = dbmserr
									}
								}
								rqststdio.inbuf.Clear()
							}
						} else if sargs[0] == "#!commit" && rqststdio.cmd != nil {
							rqststdio.rqst.copy(rqststdio.cmd, rqststdio.rqst, false, "")
						}
						rqststdio.isDone = sargs[0] == "#!close"
					}
				} else {
					if rqststdio.lsthshlnk = sargs[0]; strings.HasPrefix(rqststdio.lsthshlnk, "#!js:") {
						if path := rqststdio.lsthshlnk[len("#!js:"):]; path != "" {
							rqststdio.lsthshlnk = ""
							rqststdio.rqst.AddPath(path)
							rqststdio.rqst.processPaths(false)
						} else {
							rqststdio.lsthshlnk = ""
						}
					} else if rqststdio.lsthshlnk != "" && strings.HasPrefix(rqststdio.lsthshlnk, "#!") && !(rqststdio.lsthshlnk == "#!js" || rqststdio.lsthshlnk == "#!close" || rqststdio.lsthshlnk == "#!exit" || rqststdio.lsthshlnk == "#!commit" || rqststdio.lsthshlnk == "#!dbms") {
						if cmd, cmderr := osprc.NewCommand(rqststdio.lsthshlnk[len("#!"):], rqststdio.lsthshlnkargs...); cmderr == nil {
							if rqststdio.cmd != nil {
								rqststdio.cmd.Close()
								rqststdio.cmd = nil
							}
							rqststdio.cmd = cmd
							rqststdio.rqst.copy(rqststdio.cmd, rqststdio.rqst, false, "")
							/*go func() {
								//cmdp := make([]byte, 1024)
								//cmdpn := 0
								//var cmderr error
								for rqststdio.cmd != nil && !rqststdio.isDone {
									//cmdpn, cmderr = cmd.Read(cmdp)
									//cmd.Print()
									//if cmdpn > 0 && (cmderr == nil || cmderr == io.EOF) {
									//	rqststdio.rqst.Write(cmdp[:cmdpn])
									//}
									rqststdio.Print(rqststdio.cmd)
								}
							}()*/
						} else {
							err = cmderr
							rqststdio.isDone = true
						}
					}
				}
			}
			rqststdio.hashbang = false
			rqststdio.prvinr = rune(0)
			return
		}
		if rqststdio.prvinr == rune(13) {
			rqststdio.inbuf.Print(string(rqststdio.prvinr), string(r))
			rqststdio.prvinr = 0
			return
		}
		rqststdio.inbuf.Print(string(r))
	} else {
		if r == '\r' {
			if rqststdio.prvinr == '\r' {
				if rqststdio.tmpbuf.Size() > 0 {
					rqststdio.inbuf.Print(rqststdio.tmpbuf)
					rqststdio.tmpbuf.Clear()
				}
				rqststdio.inbuf.Print(string(rqststdio.prvinr), string(r))
			}
		} else {
			if rqststdio.prvinr == '\r' {
				if rqststdio.tmpbuf.Size() > 0 {
					rqststdio.inbuf.Print(rqststdio.tmpbuf)
					rqststdio.tmpbuf.Clear()
				}
				rqststdio.inbuf.Print(string(rqststdio.prvinr), string(r))
			} else if r == '#' {
				if rqststdio.prvinr != rune(0) && rqststdio.prvinr != '\n' {
					rqststdio.inbuf.Print(string(r))
				}
			} else {
				if !rqststdio.hashbang {
					if r == '!' {
						if rqststdio.prvinr == '#' {
							if rqststdio.tmpbuf.Size() == 0 {
								rqststdio.hashbang = true
							}
							rqststdio.tmpbuf.Print(string(rqststdio.prvinr), string(r))
						} else {
							if rqststdio.cmd == nil {
								rqststdio.inbuf.Print(string(r))
							} else {
								rqststdio.cmd.Print(string(r))
							}
						}
					} else {
						if rqststdio.cmd == nil {
							rqststdio.inbuf.Print(string(r))
						} else {
							rqststdio.cmd.Print(string(r))
						}
					}
				} else {
					rqststdio.tmpbuf.Print(string(r))
				}
			}
		}
	}
	rqststdio.prvinr = r
	return
}

func (rqststdio *requeststdio) executeStdIO() (err error) {
	rqststdio.wg.Add(1)
	go func() {
		defer rqststdio.wg.Done()
		var rnserr error = nil
		var rdr io.RuneReader = nil
		if rqstr := rqststdio.rqst.rqstr; rqstr != nil {
			if stdio, stdiook := rqstr.(*os.File); stdiook {
				pr, pw := io.Pipe()
				wg := &sync.WaitGroup{}
				wg.Add(1)
				go func() {
					defer pw.Close()
					wg.Done()
					scnr := bufio.NewScanner(stdio)
					fmt.Print("")
					txt := ""
					for {
						if scnr.Scan() {
							if txt = scnr.Text(); txt != "" {
								txt += "\r\n"
								pw.Write([]byte(txt))
							}
						}
					}
				}()
				wg.Wait()
				rdr = bufio.NewReader(pr)

			} else if wsdio, wsiook := rqstr.(*ws.ReaderWriter); wsiook {
				rdr = wsdio
			}
			if rdr == nil {
				rdr = bufio.NewReader(rqstr)
			}
			if rqststdio.rqst.initPath != "" {

			} else {
				for {
					r, s, rerr := rdr.ReadRune()
					if s > 0 && (rerr == nil || rerr == io.EOF) {
						if rnserr = rqststdio.captureRunes(false, r); rnserr != nil {
							if rerr == nil || rerr == io.EOF {
								rerr = rnserr
							}
						}
					}
					if rqststdio.isDone || rerr != nil {
						if rerr == io.EOF {
							if rqststdio.isDone {
								break
							}
							time.Sleep(time.Nanosecond * 10)
						} else {
							err = rerr
							break
						}
					}
				}
			}
		}
	}()
	rqststdio.wg.Wait()
	return
}

func (rqststdio *requeststdio) Print(a ...interface{}) {
	iorw.Fprint(rqststdio, a...)
}

func (rqststdio *requeststdio) Println(a ...interface{}) {
	iorw.Fprintln(rqststdio, a...)
}

func newrequeststdio(rqst *Request) (rqststdio *requeststdio) {
	rqststdio = &requeststdio{cmd: nil, rqst: rqst, wg: &sync.WaitGroup{}, isDone: false, hashbang: false, prvinr: rune(0), tmpbuf: iorw.NewBuffer(), inbuf: iorw.NewBuffer(), outbuf: iorw.NewBuffer()}
	return
}

func (rqststdio *requeststdio) dispose() {
	if rqststdio != nil {
		if rqststdio.tmpbuf != nil {
			rqststdio.tmpbuf.Close()
			rqststdio.tmpbuf = nil
		}
		if rqststdio.inbuf != nil {
			rqststdio.inbuf.Close()
			rqststdio.inbuf = nil
		}
		if rqststdio.outbuf != nil {
			rqststdio.outbuf.Close()
			rqststdio.outbuf = nil
		}
		if rqststdio.cmd != nil {
			rqststdio.cmd.Close()
			rqststdio.cmd = nil
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

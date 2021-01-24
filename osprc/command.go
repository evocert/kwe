package osprc

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/evocert/kwe/iorw"
)

//Command - struct
type Command struct {
	excmd      *exec.Cmd
	OnClose    func(int)
	excmdprcid int
	ctx        context.Context
	ctxcancel  context.CancelFunc
	cmdin      io.WriteCloser
	bfr        *bufio.Reader
	cmdout     io.ReadCloser
	cmdoutp    chan []byte
	cmdouterr  chan error
	cmdtmpp    []byte
	cmdtmppi   int
	cmdtmppl   int
	stdinpark  []byte
	stdinparkl int
	stdinparki int
	cancmdout  bool
}

//NewCommand return cmd *Command instance or err error
func NewCommand(execpath string, execargs ...string) (cmd *Command, err error) {
	var ctx, ctxcancel = context.WithCancel(context.Background())
	excmd := exec.CommandContext(ctx, execpath, execargs...)
	if cmdout, cmdouterr := excmd.StdoutPipe(); cmdouterr == nil {
		if cmdin, cmdinerr := excmd.StdinPipe(); cmdinerr == nil {
			if err = excmd.Start(); err == nil {
				cmd = &Command{excmd: excmd, excmdprcid: -1, OnClose: nil, ctx: ctx, ctxcancel: ctxcancel, cmdin: cmdin, cancmdout: false, cmdtmpp: make([]byte, 1024), stdinparkl: 0, stdinparki: 0, stdinpark: make([]byte, 1024), cmdtmppi: 0, cmdtmppl: 0, cmdoutp: make(chan []byte, 1), cmdouterr: make(chan error, 1), cmdout: cmdout}
				cmd.excmdprcid = excmd.Process.Pid
				go func() {
					p := make([]byte, 1024)
					n := 0
					err := error(nil)
					for {
						n, err = cmd.cmdout.Read(p)
						var bts []byte = make([]byte, n)
						if n > 0 {
							copy(bts, p)
						}
						cmd.cmdoutp <- bts
						cmd.cmdouterr <- err
						if err != nil && err != io.EOF {
							break
						}
					}
				}()
			} else {
				cmdin = nil
				cmdout = nil
				ctxcancel()
			}
		} else {
			err = cmdinerr
			ctxcancel()
		}
	} else {
		err = cmdouterr
		ctxcancel()
	}
	return
}

//PrcID underlying os Process ID
func (cmd *Command) PrcID() int {
	if cmd != nil {
		return cmd.excmdprcid
	}
	return -1
}

//Print - similar to fmt.Fprint just direct on *Command
func (cmd *Command) Print(a ...interface{}) {
	if len(a) > 0 {
		iorw.Fprint(cmd, a...)
	}
}

//Println - similar to fmt.Fprint just direct on *Command
func (cmd *Command) Println(a ...interface{}) {
	if len(a) > 0 {
		iorw.Fprint(cmd, a...)
	}
	iorw.Fprint(cmd, "\n")
}

//Readln - read line from cmd and return s string or err error
func (cmd *Command) Readln() (s string, err error) {
	if cmd.bfr == nil {
		cmd.bfr = bufio.NewReader(cmd)
	}
	s, err = iorw.ReadLine(cmd.bfr)
	if err == io.EOF {
		err = nil
	}
	return
}

//Readlines - read lines []string from cmd or err error
func (cmd *Command) Readlines() (lines []string, err error) {
	s := ""
	for err == nil {
		if s, err = iorw.ReadLine(cmd.bfr); err == nil || err == io.EOF {
			if lines == nil {
				lines = []string{}
			}
			lines = append(lines, s)
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

//ReadAll read and return content as s string or err error
func (cmd *Command) ReadAll() (s string, err error) {
	if cmd.bfr == nil {
		cmd.bfr = bufio.NewReader(cmd)
	}
	s, err = iorw.ReaderToString(cmd.bfr)
	return
}

//Close - *Command
func (cmd *Command) Close() (err error) {
	if cmd != nil {
		if cmd.OnClose != nil {
			cmd.OnClose(cmd.excmdprcid)
			cmd.OnClose = nil
		}
		if cmd.ctxcancel != nil {
			cmd.ctxcancel()
			cmd.ctxcancel = nil
		}
		if cmd.ctx != nil {
			cmd.ctx = nil
		}
		if cmd.cmdin != nil {
			cmd.cmdin.Close()
			cmd.cmdin = nil
		}
		if cmd.cmdout != nil {
			cmd.cmdout.Close()
			cmd.cmdout = nil
		}
		if cmd.excmd != nil {
			cmd.excmd.Wait()
			if rlserr := cmd.excmd.Process.Release(); rlserr != nil {
				cmd.excmd.Process.Kill()
			}
			cmd.excmd = nil
		}
		if cmd.cmdoutp != nil {
			close(cmd.cmdoutp)
			cmd.cmdoutp = nil
		}
		if cmd.cmdouterr != nil {
			close(cmd.cmdouterr)
			cmd.cmdouterr = nil
		}
		if cmd.bfr != nil {
			cmd.bfr = nil
		}
		cmd = nil
	}
	return
}

//ReadRune - refer to io.RuneReader
func (cmd *Command) ReadRune() (r rune, size int, err error) {
	if cmd.bfr == nil {
		cmd.bfr = bufio.NewReader(cmd)
	}
	r, size, err = cmd.bfr.ReadRune()
	return
}

//Dir return executing command directory
func (cmd *Command) Dir() string {
	if cmd != nil && cmd.excmd != nil {
		if pth := strings.Replace(cmd.excmd.Path, "\\", "/", -1); pth != "" {
			if strings.Index(pth, "/") > -1 {
				return pth[:strings.LastIndex(pth, "/")+1]
			}
		}
	}
	return ""
}

//Read - refer to io.Reader
func (cmd *Command) Read(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		lststderr := error(nil)
		for n < pl {
			if cmd.cmdtmppl == 0 || (cmd.cmdtmppl > 0 && cmd.cmdtmppl == cmd.cmdtmppi) {
				if cmd.cmdtmppi > 0 {
					cmd.cmdtmppi = 0
				}
				if cmd.cmdtmppl > 0 {
					cmd.cmdtmppl = 0
				}
				canCapture := true
				for canCapture {
					if cmd.stdinparkl == 0 || (cmd.stdinparkl > 0 && cmd.stdinparki == cmd.stdinparkl) {
						if cmd.stdinparki > 0 {
							cmd.stdinparki = 0
						}
						if cmd.stdinparkl > 0 {
							cmd.stdinparkl = 0
						}
						select {
						case stdin, ok := <-cmd.cmdoutp:
							if !ok {
								canCapture = false
							} else {
								lststderr = <-cmd.cmdouterr
								if cmd.stdinparkl = len(stdin); cmd.stdinparkl > 0 {
									cmd.stdinparkl = copy(cmd.stdinpark[:cmd.stdinparkl], stdin)
									canCapture = true
								} else {
									canCapture = false
								}
							}
						case <-time.After(1 * time.Second):
							canCapture = false
						}
					}
					if canCapture {
						for {
							if cmdl := len(cmd.cmdtmpp); cmd.cmdtmppi < cmdl {
								if cl := (cmd.stdinparkl - cmd.stdinparki); cl <= (cmdl - cmd.cmdtmppi) {
									copy(cmd.cmdtmpp[cmd.cmdtmppi:cmd.cmdtmppi+cl], cmd.stdinpark[cmd.stdinparki:cmd.stdinparki+cl])
									cmd.cmdtmppl += cl
									cmd.cmdtmppi += cl
									cmd.stdinparki += cl
								} else if cl := (cmdl - cmd.cmdtmppi); cl < (cmd.stdinparkl - cmd.stdinparki) {
									copy(cmd.cmdtmpp[cmd.cmdtmppi:cmd.cmdtmppi+cl], cmd.stdinpark[cmd.stdinparki:cmd.stdinparki+cl])
									cmd.cmdtmppl += cl
									cmd.cmdtmppi += cl
									cmd.stdinparki += cl
								}
								if cmdl == cmd.cmdtmppi {
									canCapture = false
									break
								}
								if cmd.stdinparki == cmd.stdinparkl {
									break
								}
							}
						}
					}
				}
				if cmd.cmdtmppl == 0 {
					break
				} else {
					cmd.cmdtmppi = 0
				}
			}
			for n < pl && cmd.cmdtmppi < cmd.cmdtmppl {
				if cl := (cmd.cmdtmppl - cmd.cmdtmppi); cl <= (pl - n) {
					copy(p[n:n+cl], cmd.cmdtmpp[cmd.cmdtmppi:cmd.cmdtmppi+cl])
					cmd.cmdtmppi += cl
					n += cl
				} else if cl := (pl - n); cl < (cmd.cmdtmppl - cmd.cmdtmppi) {
					copy(p[n:n+cl], cmd.cmdtmpp[cmd.cmdtmppi:cmd.cmdtmppi+cl])
					cmd.cmdtmppi += cl
					n += cl
				}
			}
			if cmd.cmdtmppi == cmd.cmdtmppl {
				break
			}
		}
		//if !cmd.cancmdout {
		//	go func() {
		//n, err = cmd.cmdout.Read(p)
		//		cmd.cancmdout = false
		//		cmd.cmdoutdne <- true
		//	}()
		//}
		//select {
		//case <-cmd.cmdoutdne:
		//case <-time.After(500 * time.Millisecond):
		//}
		if lststderr != nil && lststderr != io.EOF {

		}
		if n == 0 && err == nil {
			err = io.EOF
		}
	}
	return
}

//Write - refer to io.Writer
func (cmd *Command) Write(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		n, err = cmd.cmdin.Write(p)
	}
	return
}

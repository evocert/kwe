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
	//cmdinbufw  *bufio.Writer
	cmdout     io.ReadCloser
	cmdoutbufr *bufio.Reader
	cmdtmpp    []byte
	cmdtmppi   int
	cmdtmppl   int
	stdinpark  []byte
	stdinerr   chan error
	stdinparkl int
	stdinparki int
	milseconds int64
}

//NewCommand return cmd *Command instance or err error
func NewCommand(execpath string, execargs ...string) (cmd *Command, err error) {
	var ctx, ctxcancel = context.WithCancel(context.Background())
	excmd := exec.CommandContext(ctx, execpath, execargs...)
	if cmdout, cmdouterr := excmd.StdoutPipe(); cmdouterr == nil {
		if cmdin, cmdinerr := excmd.StdinPipe(); cmdinerr == nil {
			if err = excmd.Start(); err == nil {
				cmd = &Command{excmd: excmd, excmdprcid: -1, milseconds: 100, OnClose: nil, ctx: ctx, ctxcancel: ctxcancel, cmdin: cmdin, cmdtmpp: make([]byte, 1024), stdinerr: make(chan error, 1), stdinparkl: 0, stdinparki: 0, stdinpark: make([]byte, 1024), cmdtmppi: 0, cmdtmppl: 0 /* cmdoutp: make(chan []byte, 1), cmdouterr: make(chan error, 1),*/, cmdout: cmdout}
				cmd.excmdprcid = excmd.Process.Pid
				cmd.cmdoutbufr = bufio.NewReader(cmd)

				go func() {
					running := true
					for running {
						var stdinerr error = nil
						cmd.stdinparkl, stdinerr = cmd.cmdout.Read(cmd.stdinpark)
						cmd.stdinerr <- stdinerr
						if stdinerr != nil && stdinerr != io.EOF {
							running = false
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

//SetReadTimeout set read timeout in milliseconds int64
func (cmd *Command) SetReadTimeout(milseconds int64) {
	if milseconds < 100 {
		milseconds = 100
	}
	cmd.milseconds = milseconds
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

//Seek - refer tio iorw.Reader empty implementation
func (cmd *Command) Seek(offset int64, whence int) (n int64, err error) {
	return
}

//Readln - read line from cmd and return s string or err error
func (cmd *Command) Readln() (s string, err error) {
	s, err = iorw.ReadLine(cmd)
	if err == io.EOF {
		err = nil
	}
	return
}

//Readlines - read lines []string from cmd or err error
func (cmd *Command) Readlines() (lines []string, err error) {
	s := ""
	for err == nil {
		if s, err = cmd.Readln(); err == nil || err == io.EOF {
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
	s, err = iorw.ReaderToString(cmd)
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
		if cmd.cmdoutbufr != nil {
			cmd.cmdoutbufr = nil
		}
		/*if cmd.cmdinbufw != nil {
			cmd.cmdinbufw = nil
		}*/
		cmd = nil
	}
	return
}

//ReadRune - refer to io.RuneReader
func (cmd *Command) ReadRune() (r rune, size int, err error) {
	r, size, err = cmd.cmdoutbufr.ReadRune()
	return
}

//Dir return executing command directory
func (cmd *Command) Dir() string {
	if cmd != nil && cmd.excmd != nil {
		if pth := strings.Replace(cmd.excmd.Path, "\\", "/", -1); pth != "" {
			if strings.Contains(pth, "/") {
				return pth[:strings.LastIndex(pth, "/")+1]
			}
		}
	}
	return ""
}

//Flush return error
//flush wrapping *bufio.Writer if buffered
func (cmd *Command) Flush() (err error) {
	/*if cmd.cmdinbufw.Buffered() > 0 {
		err = cmd.cmdinbufw.Flush()
	}*/
	return
}

func (cmd *Command) Reset() {
	cmd.ResetRead()
	cmd.ResetWrite()
}

func (cmd *Command) ResetRead() {
	cmd.cmdoutbufr.Reset(cmd)
	cmd.stdinerr = nil
	cmd.stdinparki = 0
	cmd.stdinparkl = 0
}

func (cmd *Command) ResetWrite() {
	//cmd.cmdinbufw.Reset(cmd)
}

//Read - refer to io.Reader
func (cmd *Command) Read(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		lststderr := error(nil)
		if err = cmd.Flush(); err == nil {
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

							tmelpsed := time.After(time.Duration(cmd.milseconds) * time.Millisecond)
							var stdinerr error = nil
							var stdok = false
							select {
							case <-tmelpsed:
								canCapture = false
							case stdinerr, stdok = <-cmd.stdinerr:
								if !stdok {
									canCapture = false
								}
								if stdinerr != nil {
									if stdinerr != io.EOF {
										err = stdinerr
										canCapture = false
									}
								}
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
		}
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

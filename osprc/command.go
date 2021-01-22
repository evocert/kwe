package osprc

import (
	"bufio"
	"context"
	"io"
	"os/exec"

	"github.com/evocert/kwe/iorw"
)

//Command - struct
type Command struct {
	excmd     *exec.Cmd
	ctx       context.Context
	ctxcancel context.CancelFunc
	cmdin     io.WriteCloser
	pin       chan []byte
	bfr       *bufio.Reader
	cmdout    io.ReadCloser
	pout      chan []byte
}

//NewCommand return cmd *Command instance or err error
func NewCommand(execpath string, execargs ...string) (cmd *Command, err error) {
	var ctx, ctxcancel = context.WithCancel(context.Background())
	excmd := exec.CommandContext(ctx, execpath, execargs...)

	if cmdout, cmdouterr := excmd.StdoutPipe(); cmdouterr == nil {
		if cmdin, cmdinerr := excmd.StdinPipe(); cmdinerr == nil {
			if err = excmd.Start(); err == nil {
				cmd = &Command{excmd: excmd, ctx: ctx, ctxcancel: ctxcancel, cmdin: cmdin, cmdout: cmdout}
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

//Read - refer to io.Reader
func (cmd *Command) Read(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		n, err = cmd.cmdout.Read(p)
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

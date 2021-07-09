package jsext

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/evocert/kwe/iorw"
)

func Register_jsext_consoleutils(vm *goja.Runtime) {
	//vm.SetFieldNameMapper(goja.TagFieldNameMapper("json",true))
	type Version struct {
		Major int `json:"major"`
		Minor int `json:"minor"`
		Bump  int `json:"bump"`
	}
	//log.SetFlags(log.LstdFlags | log.Lmicroseconds) //global???
	//todo: namespace everything kwe.fsutils.etcetcetc
	//first test for kwe then do set fsutils on kwe
	vm.Set("console", struct {
		Version Version              `json:"version"`
		Log     func(...interface{}) `json:"log"`
		Warn    func(...interface{}) `json:"warn"`
		Error   func(...interface{}) `json:"error"`
		Debug   func(...interface{}) `json:"debug"`
		Trace   func(...interface{}) `json:"trace"`
	}{
		Version: Version{
			Major: 0,
			Minor: 0,
			Bump:  1,
		},
		//todo: colors
		Log: func(msg ...interface{}) {
			//buf := iorw.NewBuffer()
			//buf.Print(msg...)
			//lgmsg := buf.String()
			//buf.Close()
			//buf = nil
			//log.Println("LOG:   ", lgmsg)
			//logger.Output(2, fmt.Sprintln("LOG:   "+lgmsg))
			if len(msg) > 0 {
				msg = append([]interface{}{"LOG   "}, msg...)
			}
			logger.Output(2, fmt.Sprintln(msg...))
			//rw.Println(msg)
		},
		Warn: func(msg ...interface{}) {
			//buf := iorw.NewBuffer()
			//buf.Print(msg...)
			//lgmsg := buf.String()
			//buf = nil
			//log.Println("WARN:  ", lgmsg)
			//logger.Output(2, fmt.Sprintln("WARN:   "+lgmsg))
			if len(msg) > 0 {
				msg = append([]interface{}{"WARN   "}, msg...)
			}
			logger.Output(2, fmt.Sprintln(msg...))
			rw.Println(msg)
		},
		Error: func(msg ...interface{}) {
			//buf := iorw.NewBuffer()
			//buf.Print(msg...)
			//lgmsg := buf.String()
			//buf = nil
			//log.Println("ERROR: ", lgmsg)
			//logger.Output(2, fmt.Sprintln("ERROR:   "+lgmsg))
			if len(msg) > 0 {
				msg = append([]interface{}{"ERROR   "}, msg...)
			}
			logger.Output(2, fmt.Sprintln(msg...))
			//rw.Println(msg)
		},
		Debug: func(msg ...interface{}) {
			//buf := iorw.NewBuffer()
			//buf.Print(msg...)
			//lgmsg := buf.String()
			//buf = nil
			//log.Println("DEBUG: ", lgmsg)
			//logger.Output(2, fmt.Sprintln("DEBUG:   "+lgmsg))
			if len(msg) > 0 {
				msg = append([]interface{}{"DEBUG   "}, msg...)
			}
			logger.Output(2, fmt.Sprintln(msg...))
			//rw.Println(msg)
		},
		Trace: func(msg ...interface{}) {
			//buf := iorw.NewBuffer()
			//buf.Print(msg...)
			//lgmsg := buf.String()
			//buf = nil
			//log.Println("TRACE: ", lgmsg)
			if len(msg) > 0 {
				msg = append([]interface{}{"TRACE   "}, msg...)
			}
			logger.Output(2, fmt.Sprintln(msg...))
			//rw.Println(msg)
		},
	})
}

var logger *log.Logger = nil

type readwrite struct {
	inbytes  []byte
	inbf     *iorw.Buffer
	outw     io.Writer
	lckinout *sync.RWMutex
	flshdr   time.Duration
}

func (rw *readwrite) Write(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		bi := 0
		bl := len(rw.inbytes)
		bn := 0
		for _, pb := range p {
			if pb == '\n' {
				func() {
					if bi > 0 {
						//rw.lckinout.RLock()
						//defer rw.lckinout.RUnlock()
						bn, err = rw.inbf.Write(rw.inbytes[:bi])
						if bn > 0 {
							n += bn
						}
						bi = 0
					}
					//rw.lckinout.Lock()
					//defer rw.lckinout.Unlock()
					if hdrln := rw.inbf.String(); hdrln != "" {
						rw.inbf.Clear()
						rlines <- (strings.TrimSpace(hdrln) + "\n")
					}
				}()
			} else {
				rw.inbytes[bi] = pb
				bi++
				if bi == bl {
					func() {
						//rw.lckinout.RLock()
						//defer rw.lckinout.RUnlock()
						bn, err = rw.inbf.Write(rw.inbytes[:bi])
						if bn > 0 {
							n += bn
						}
					}()
					bi = 0
				}
			}
		}

		if bi > 0 {
			func() {
				//rw.lckinout.RLock()
				//defer rw.lckinout.RUnlock()
				bn, err = rw.inbf.Write(rw.inbytes[:bi])
				if bn > 0 {
					n += bn
				}
			}()
		}
	}
	return
}

func (rw *readwrite) Print(a ...interface{}) {
	if rw != nil {
		iorw.Fprint(rw, a...)
	}
}

func (rw *readwrite) Println(a ...interface{}) {
	if rw != nil {
		iorw.Fprintln(rw, a...)
	}
}

var rlines chan string = make(chan string)

var rw *readwrite = nil

func init() {
	rw = &readwrite{inbf: iorw.NewBuffer(), lckinout: &sync.RWMutex{}, flshdr: time.Millisecond * 500, outw: os.Stderr, inbytes: make([]byte, 8192)}
	go func() {
		for {
			for lgln := range rlines {
				if lgln != "" {
					rw.outw.Write([]byte(lgln))
				}
			}
		}
	}()
	logger = log.New(rw, "", log.LstdFlags|log.Lmicroseconds)
}

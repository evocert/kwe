package ws

import (
	"bufio"
	"io"
	"strings"

	"github.com/evocert/kwe/iorw"
	"github.com/gorilla/websocket"
)

//ReaderWriter - struct
type ReaderWriter struct {
	ws       *websocket.Conn
	r        io.Reader
	rbuf     *bufio.Reader
	rerr     error
	w        io.WriteCloser
	wbuf     *bufio.Writer
	werr     error
	isText   bool
	isBinary bool
}

//NewReaderWriter - instance
func NewReaderWriter(ws *websocket.Conn) (wsrw *ReaderWriter) {
	wsrw = &ReaderWriter{ws: ws, isText: false, isBinary: false, rerr: nil, werr: nil}
	return
}

//ReadRune - refer to io.RuneReader
func (wsrw *ReaderWriter) ReadRune() (r rune, size int, err error) {
	if wsrw.rbuf == nil {
		wsrw.rbuf = bufio.NewReader(wsrw)
	}
	r, size, err = wsrw.rbuf.ReadRune()
	return
}

//WriteRune - refer to bufio.Writer - WriteRune
func (wsrw *ReaderWriter) WriteRune(r rune) (size int, err error) {
	if wsrw.wbuf == nil {
		wsrw.wbuf = bufio.NewWriter(wsrw)
	}
	size, err = wsrw.wbuf.WriteRune(r)
	return
}

//CanRead - can Read
func (wsrw *ReaderWriter) CanRead() bool {
	return wsrw.rerr == nil
}

//CanWrite - can Write
func (wsrw *ReaderWriter) CanWrite() bool {
	return wsrw.werr == nil
}

//Read - refer io.Reader
func (wsrw *ReaderWriter) Read(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		if wsrw.r == nil {
			if err = wsrw.Flush(); err == nil {
				if wsrw.CanRead() {
					var messageType int
					messageType, wsrw.r, wsrw.rerr = wsrw.ws.NextReader()
					wsrw.isText = messageType == websocket.TextMessage
					wsrw.isBinary = messageType == websocket.BinaryMessage
					if wsrw.rerr != nil {
						if wsrw.rerr != io.EOF {
							return 0, wsrw.rerr
						}
						return 0, io.EOF
					}
				}
			} else {
				return 0, io.EOF
			}
		}
		for n = 0; n < len(p); {
			var m int
			m, err = wsrw.r.Read(p[n:])
			n += m
			if err != nil {
				if err == io.EOF {
					// done
					wsrw.r = nil
					break
				} else {
					wsrw.rerr = err
				}
			}
			if err != nil {
				break
			}
		}

		if n == 0 && err == nil {
			err = io.EOF
		}
	}
	return
}

//Readln - read single line
func (wsrw *ReaderWriter) Readln() (s string, err error) {
	s = ""
	var rns = make([]rune, 1024)
	var rnsi = 0
	for {
		rn, size, rnerr := wsrw.ReadRune()
		if size > 0 {
			if rn == rune(10) {
				if rnsi > 0 {
					s += string(rns[:rnsi])
					rnsi = 0
				}
				break
			} else {
				rns[rnsi] = rn
				rnsi++
				if rnsi == len(rns) {
					s += string(rns[:rnsi])
					rnsi = 0
				}
			}
		}
		if rnerr != nil {
			if rnerr != io.EOF {
				err = rnerr
			}
			break
		}
	}
	if s == "" && rnsi > 0 {
		s += string(rns[:rnsi])
		rnsi = 0
	}
	if s != "" {
		s = strings.TrimSpace(s)
	}
	return
}

//Readlines - return lines []string slice
func (wsrw *ReaderWriter) Readlines() (lines []string, err error) {
	var line = ""
	for {
		if line, err = wsrw.Readln(); line != "" && (err == nil || err == io.EOF) {
			if lines == nil {
				lines = []string{}
			}
			lines = append(lines, line)
		}
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
	}
	return
}

//ReadAll - return all read content as string
func (wsrw *ReaderWriter) ReadAll() (s string, err error) {
	s = ""
	var rns = make([]rune, 1024)
	var rnsi = 0
	for {
		rn, size, rnerr := wsrw.ReadRune()
		if size > 0 {
			rns[rnsi] = rn
			rnsi++
			if rnsi == len(rns) {
				s += string(rns[:rnsi])
				rnsi = 0
			}
		}
		if rnerr != nil {
			if rnerr != io.EOF {
				err = rnerr
			}
			break
		}
	}
	if rnsi > 0 {
		s += string(rns[:rnsi])
		rnsi = 0
	}
	return
}

func (wsrw *ReaderWriter) socketIOType() int {
	if wsrw.isText {
		return websocket.TextMessage
	} else if wsrw.isBinary {
		return websocket.BinaryMessage
	}
	return websocket.TextMessage
}

//Flush - flush invoke done onmessage
func (wsrw *ReaderWriter) Flush() (err error) {
	if wsrw.wbuf != nil {
		if err = wsrw.wbuf.Flush(); err != nil {
			return
		}
	}
	if wsrw.w != nil {
		err = wsrw.w.Close()
		wsrw.w = nil
	}
	return
}

//Print - refer to fmt.Fprint
func (wsrw *ReaderWriter) Print(a ...interface{}) {
	iorw.Fprint(wsrw, a...)
	wsrw.Flush()
}

//Println - refer to fmt.Fprintln
func (wsrw *ReaderWriter) Println(a ...interface{}) {
	iorw.Fprintln(wsrw, a...)
	wsrw.Flush()
}

//Write - refer io.Writer
func (wsrw *ReaderWriter) Write(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		if wsrw.w == nil && wsrw.CanWrite() {
			wsrw.w, wsrw.werr = wsrw.ws.NextWriter(wsrw.socketIOType())
			if wsrw.werr != nil {
				err = wsrw.werr
				return 0, err
			}
		}
		for n = 0; n < len(p); {
			var m int
			m, err = wsrw.w.Write(p[n : n+(len(p)-n)])
			n += m
			if err != nil {
				break
			}
		}
		if n == 0 && err == nil {
			err = io.EOF
		}
	}
	return
}

//Close - refer io.Closer
func (wsrw *ReaderWriter) Close() (err error) {
	if wsrw != nil {
		if wsrw.r != nil {
			wsrw.r = nil
		}
		if wsrw.w != nil {
			wsrw.w.Close()
			wsrw.w = nil
		}
		if wsrw.rbuf != nil {
			wsrw.rbuf = nil
		}
		if wsrw.wbuf != nil {
			wsrw.wbuf = nil
		}
		if wsrw.ws != nil {
			err = wsrw.ws.Close()
			wsrw.ws = nil
		}
		wsrw = nil
	}
	return
}

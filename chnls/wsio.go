package chnls

import (
	"bufio"
	"io"

	"github.com/evocert/kwe/iorw"
	"github.com/gorilla/websocket"
)

//WsReaderWriter - struct
type WsReaderWriter struct {
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

//NewWsReaderWriter - instance
func NewWsReaderWriter(ws *websocket.Conn) (wsrw *WsReaderWriter) {
	wsrw = &WsReaderWriter{ws: ws, isText: false, isBinary: false, rerr: nil, werr: nil}
	return
}

//ReadRune - refer to io.RuneReader
func (wsrw *WsReaderWriter) ReadRune() (r rune, size int, err error) {
	if wsrw.rbuf == nil {
		wsrw.rbuf = bufio.NewReader(wsrw)
	}
	r, size, err = wsrw.rbuf.ReadRune()
	return
}

//WriteRune - refer to bufio.Writer - WriteRune
func (wsrw *WsReaderWriter) WriteRune(r rune) (size int, err error) {
	if wsrw.wbuf == nil {
		wsrw.wbuf = bufio.NewWriter(wsrw)
	}
	size, err = wsrw.wbuf.WriteRune(r)
	return
}

//CanRead - can Read
func (wsrw *WsReaderWriter) CanRead() bool {
	return wsrw.rerr == nil
}

//CanWrite - can Write
func (wsrw *WsReaderWriter) CanWrite() bool {
	return wsrw.werr == nil
}

//Read - refer io.Reader
func (wsrw *WsReaderWriter) Read(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		if wsrw.r == nil {
			var messageType int
			messageType, wsrw.r, wsrw.rerr = wsrw.ws.NextReader()
			wsrw.isText = messageType == websocket.TextMessage
			wsrw.isBinary = messageType == websocket.BinaryMessage
			//_, wsrw.r, wsrw.rerr = wsrw.ws.NextReader()
			if wsrw.rerr != nil {
				return 0, io.EOF
			}
		}
		for n = 0; n < len(p); {
			var m int
			m, err = wsrw.r.Read(p[n:])
			n += m
			if err == io.EOF {
				// done
				wsrw.r = nil
				break
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

func (wsrw *WsReaderWriter) socketIOType() int {
	if wsrw.isText {
		return websocket.TextMessage
	} else if wsrw.isBinary {
		return websocket.BinaryMessage
	}
	return websocket.TextMessage
}

//Flush - flush invoke done onmessage
func (wsrw *WsReaderWriter) Flush() (err error) {
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
func (wsrw *WsReaderWriter) Print(a ...interface{}) {
	iorw.Fprint(wsrw, a...)
}

//Println - refer to fmt.Fprintln
func (wsrw *WsReaderWriter) Println(a ...interface{}) {
	iorw.Fprintln(wsrw, a...)
}

//Write - refer io.Writer
func (wsrw *WsReaderWriter) Write(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		if wsrw.w == nil {
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
func (wsrw *WsReaderWriter) Close() (err error) {
	if wsrw != nil {
		if wsrw.ws != nil {
			err = wsrw.ws.Close()
		}
		wsrw.ws = nil
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
		wsrw = nil
	}
	return
}

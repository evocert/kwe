package chnls

import (
	"io"

	"github.com/gorilla/websocket"
)

//WsReaderWriter - struct
type WsReaderWriter struct {
	ws       *websocket.Conn
	r        io.Reader
	w        io.Writer
	isText   bool
	isBinary bool
}

//NewWsReaderWriter - instance
func NewWsReaderWriter(ws *websocket.Conn) (wsrw *WsReaderWriter) {
	wsrw = &WsReaderWriter{ws: ws, isText: false, isBinary: false}
	return
}

//Read - refer io.Reader
func (wsrw *WsReaderWriter) Read(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {

		if wsrw.r == nil {
			var messageType int
			messageType, wsrw.r, err = wsrw.ws.NextReader()
			wsrw.isText = messageType == websocket.TextMessage
			wsrw.isBinary = messageType == websocket.BinaryMessage
			_, wsrw.r, err = wsrw.ws.NextReader()
			if err != nil {
				return 0, err
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

//Write - refer io.Writer
func (wsrw *WsReaderWriter) Write(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		if wsrw.w == nil {
			wsrw.w, err = wsrw.ws.NextWriter(wsrw.socketIOType())
			if err != nil {
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
			wsrw.w = nil
		}
		wsrw = nil
	}
	return
}

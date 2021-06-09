package caching

import (
	"context"
	"io"

	"github.com/evocert/kwe/iorw"
)

type MapAPI interface {
	Keys(...interface{}) []interface{}
	Values(...interface{}) []interface{}
	Put(interface{}, ...interface{})
	Remove(...interface{})
	Find(...interface{}) interface{}
	Fprint(io.Writer, ...interface{}) error
	String(...interface{}) string
	Focus(...interface{}) bool
	Reset(...interface{})
	Clear(...interface{})
	Close(...interface{})
	//Array
	Push(interface{}, ...interface{}) int
	Pop(interface{}, ...interface{}) interface{}
	Shift(interface{}, ...interface{}) int
	Unshift(interface{}, ...interface{}) interface{}
	At(interface{}, ...interface{}) interface{}
	FocusAt(interface{}, ...interface{}) bool
}

func MapReader(mapapi MapAPI, ks ...interface{}) (rdr *iorw.EOFCloseSeekReader) {
	pi, pw := io.Pipe()
	cntx, cntxcancel := context.WithCancel(context.Background())
	go func() {
		defer func() {
			pw.Close()
		}()
		cntxcancel()
		if mapapi != nil {
			mapapi.Fprint(pw, ks...)
		} else {
			iorw.Fprint(pw, "{}")
		}
	}()
	<-cntx.Done()
	rdr = iorw.NewEOFCloseSeekReader(pi)
	return
}

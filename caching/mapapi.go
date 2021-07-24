package caching

import (
	"context"
	"io"

	"github.com/evocert/kwe/iorw"
)

type MapAPI interface {
	Keys(...interface{}) []interface{}
	Values(...interface{}) []interface{}
	IsMap(...interface{}) bool
	Exists(...interface{}) bool
	Find(...interface{}) interface{}
	Put(interface{}, ...interface{})
	Remove(...interface{})
	Fprint(io.Writer, ...interface{}) error
	String(...interface{}) string
	Focus(...interface{}) bool
	Reset(...interface{}) bool
	Clear(...interface{}) bool
	Close(...interface{}) bool
	//Array
	IsMapAt(interface{}, ...interface{}) bool
	ExistsAt(interface{}, ...interface{}) bool
	Push(interface{}, ...interface{}) int
	Pop(interface{}, ...interface{}) interface{}
	Shift(interface{}, ...interface{}) int
	Unshift(interface{}, ...interface{}) interface{}
	At(interface{}, ...interface{}) interface{}
	FocusAt(interface{}, ...interface{}) bool
	ClearAt(interface{}, ...interface{}) bool
	CloseAt(interface{}, ...interface{}) bool
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

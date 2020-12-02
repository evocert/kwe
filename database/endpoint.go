package database

import "io"

//EndPoint - struct
type EndPoint struct {
	datasource string
	args       []interface{}
}

func newEndPoint(datasource string, a ...interface{}) (endpnt *EndPoint) {
	endpnt = &EndPoint{datasource: datasource, args: a}
	return
}

func (endpnt *EndPoint) query(out io.Writer, iorags ...interface{}) (err error) {
	return
}

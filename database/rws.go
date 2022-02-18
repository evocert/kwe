package database

import (
	"database/sql"
	"io"
)

type RWSReader struct {
	rdr       io.RuneReader
	lsterr    error
	strmstngs map[string]interface{}
	sqlrws    *sql.Rows
	coltypes  []*ColumnType
	cls       []string
}

func newRWSReader(sqlrws *sql.Rows, strmstngs map[string]interface{}) (rwsrrdr *RWSReader, err error) {
	if len(strmstngs) > 0 {
		rwsrrdr = &RWSReader{strmstngs: strmstngs}
	}
	return
}

func (rwsrdr *RWSReader) Close() (err error) {
	if rwsrdr != nil {
		if rwsrdr.sqlrws != nil {
			err = rwsrdr.sqlrws.Close()
			rwsrdr.sqlrws = nil
		}
		if rwsrdr.strmstngs != nil {
			for strmk := range rwsrdr.strmstngs {
				delete(rwsrdr.strmstngs, strmk)
			}
			rwsrdr.strmstngs = nil
		}
		if rwsrdr.coltypes != nil {
			rwsrdr.coltypes = nil
		}
		rwsrdr = nil
	}
	return
}

func (rwsrdr *RWSReader) Err() (err error) {
	if rwsrdr != nil {
		if rwsrdr.sqlrws != nil && len(rwsrdr.strmstngs) == 0 && rwsrdr.rdr == nil {
			rwsrdr.lsterr = rwsrdr.sqlrws.Err()
		}
		err = rwsrdr.lsterr
	}
	return
}

func (rwsrdr *RWSReader) Next() (nxt bool) {
	if rwsrdr != nil {
		if rwsrdr.sqlrws != nil && len(rwsrdr.strmstngs) == 0 && rwsrdr.rdr == nil {
			nxt = rwsrdr.sqlrws.Next()
		}
	}
	return
}

func (rwsrdr *RWSReader) Scan(dest ...interface{}) (err error) {
	if rwsrdr != nil {
		if rwsrdr.sqlrws != nil && len(rwsrdr.strmstngs) == 0 && rwsrdr.rdr == nil {
			err = rwsrdr.sqlrws.Scan(dest...)
		}
	}
	return
}

func (rwsrdr *RWSReader) ColumnTypes() (coltypes []*ColumnType, err error) {
	if rwsrdr != nil {
		prepRWSColumns(rwsrdr)
	}
	return
}

func (rwsrdr *RWSReader) Columns() (cls []string, err error) {
	if rwsrdr != nil {
		prepRWSColumns(rwsrdr)
	}
	return
}

func prepRWSColumns(rwsrdr *RWSReader) {
	if len(rwsrdr.cls) == 0 {
		if rwsrdr.sqlrws != nil && len(rwsrdr.strmstngs) == 0 && rwsrdr.rdr == nil {
			cltps, _ := rwsrdr.sqlrws.ColumnTypes()
			cls, _ := rwsrdr.sqlrws.Columns()
			if len(rwsrdr.cls) == 0 && len(cls) > 0 {
				rwsrdr.cls = cls[:]
				rwsrdr.coltypes = columnTypes(cltps, cls)
			}
		}
	}
}

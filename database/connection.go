package database

import (
	"bufio"
	"database/sql"
	"io"
	"runtime"
	"strings"
)

//Connection - struct
type Connection struct {
	dbms                       *DBMS
	driverName, dataSourceName string
	db                         *sql.DB
	dbinvoker                  func(string) (*sql.DB, error)
}

func runeReaderToString(rnr io.RuneReader) (s string) {
	return
}

func queryToStatement(cn *Connection, query interface{}) (stmnt string) {
	var rnrr io.RuneReader = nil
	if qrys, qrysok := query.(string); qrysok && qrys != "" {
		rnrr = bufio.NewReader(strings.NewReader(qrys))
	} else if qryrnr, qryrnrok := query.(io.RuneReader); qryrnrok {
		rnrr = qryrnr
	} else if qryr, qryrok := query.(io.Reader); qryrok {
		rnrr = bufio.NewReader(qryr)
	}

	stmnt = ""

	var rns = make([]rune, 1024)
	var rnsi = 0
	for rnrr != nil {
		r, s, e := rnrr.ReadRune()
		if s > 0 {
			rns[rnsi] = r
			rnsi++
			if rnsi == len(rns) {
				stmnt += string(rns)
			}
		}
		if e != nil {
			break
		}
	}
	if rnsi > 0 {
		stmnt += string(rns[:rnsi])
	}
	return
}

func (cn *Connection) query(query interface{}, noreader bool, prms ...interface{}) (reader *Reader, exctr *Executor, err error) {
	var stmnt = queryToStatement(cn, query)
	if cn.db == nil {
		if cn.dbinvoker == nil {
			if dbinvoker, hasdbinvoker := cn.dbms.driverDbInvoker(cn.driverName); hasdbinvoker {
				cn.dbinvoker = dbinvoker
			}
		}
		if cn.db, err = cn.dbinvoker(cn.dataSourceName); err == nil && cn.db != nil {
			cn.db.SetMaxIdleConns(runtime.NumCPU() * 4)
		}
	}
	if cn.db != nil {
		if stmnt != "" {
			exctr = newExecutor(cn, cn.db, stmnt)
			if noreader {
				exctr.execute()
			} else {
				reader = newReader(exctr)
				reader.execute()
			}
		}
	}
	return
}

//NewConnection - dbms,driver name and datasource name (cn-string)
func NewConnection(dbms *DBMS, driverName, dataSourceName string) (cn *Connection) {
	cn = &Connection{dbms: dbms, driverName: driverName, dataSourceName: dataSourceName}
	return
}

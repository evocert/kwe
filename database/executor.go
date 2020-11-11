package database

import (
	"database/sql"

	"github.com/evocert/kwe/parameters"
)

//Executor - struct
type Executor struct {
	stmnt        string
	db           *sql.DB
	cn           *Connection
	stmt         *sql.Stmt
	lasterr      error
	lastInsertID int64
	rowsAffected int64
	mappedArgs   map[string]interface{}
	qryArgs      []interface{}
	OnSuccess    interface{}
	OnError      interface{}
	OnFinalize   interface{}
}

func newExecutor(cn *Connection, db *sql.DB, query interface{}, onsuccess, onerror, onfinalize interface{}, args ...interface{}) (exctr *Executor) {
	var argsn = 0
	for argsn < len(args) {
		var d = args[argsn]
		if _, dok := d.(*parameters.Parameters); dok {
			argsn++
		} else if _, dok := d.(map[string]interface{}); dok {
			argsn++
		} else {
			if d != nil {
				if onsuccess == nil {
					onsuccess = d
				} else if onerror == nil {
					onerror = d
				} else if onfinalize == nil {
					onfinalize = d
				}
				if argsn == len(args) {
					args = append(args[:argsn])
				} else if argsn < len(args) {
					args = append(args[:argsn], args[argsn+1:]...)
				}
			}
		}
	}
	exctr = &Executor{db: db, cn: cn, OnSuccess: onsuccess, OnError: onerror, OnFinalize: onfinalize}
	exctr.stmnt, exctr.qryArgs, exctr.mappedArgs = queryToStatement(exctr, query, args...)
	return
}

func (exctr *Executor) execute(forrows ...bool) (rws *sql.Rows, cltpes []*sql.ColumnType, cls []string) {
	if exctr.stmt, exctr.lasterr = exctr.db.Prepare(exctr.stmnt); exctr.lasterr == nil && exctr.stmt != nil {
		exctr.lastInsertID = -1
		exctr.rowsAffected = -1
		if len(forrows) >= 1 && forrows[0] {
			if rws, exctr.lasterr = exctr.stmt.Query(exctr.qryArgs...); rws != nil && exctr.lasterr == nil {
				invokeSuccess(exctr.OnSuccess)
				cltpes, _ = rws.ColumnTypes()
				cls, _ = rws.Columns()
			} else if exctr.lasterr != nil {
				invokeError(exctr.lasterr, exctr.OnError)
			}
		} else {
			if rslt, rslterr := exctr.stmt.Exec(exctr.qryArgs...); rslterr == nil {
				if exctr.lastInsertID, rslterr = rslt.LastInsertId(); rslterr != nil {
					exctr.lastInsertID = -1
				}
				if exctr.rowsAffected, rslterr = rslt.RowsAffected(); rslterr != nil {
					exctr.rowsAffected = -1
				}
				invokeSuccess(exctr.OnSuccess)

			} else {
				exctr.lasterr = rslterr
				invokeError(exctr.lasterr, exctr.OnSuccess)
			}
		}
	}
	return
}

//Close - Executor
func (exctr *Executor) Close() (err error) {
	if exctr != nil {
		if exctr.db != nil {
			exctr.db = nil
		}
		if exctr.cn != nil {
			exctr.cn = nil
		}
		if exctr.stmt != nil {
			err = exctr.stmt.Close()
			exctr.stmt = nil
		}
		if exctr.lasterr != nil {
			exctr.lasterr = nil
		}
		invokeFinalize(exctr.OnFinalize)
		if exctr.OnSuccess != nil {
			exctr.OnSuccess = nil
		}
		if exctr.OnError != nil {
			exctr.OnError = nil
		}
		if exctr.OnFinalize != nil {
			exctr.OnFinalize = nil
		}
		exctr = nil
	}
	return
}

//Err - return last Error
func (exctr *Executor) Err() (err error) {
	if exctr != nil {
		err = exctr.lasterr
	}
	return
}

package database

import (
	"database/sql"

	"github.com/evocert/kwe/iorw/active"
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
	argNames     []string
	qryArgs      []interface{}
	OnSuccess    interface{}
	OnError      interface{}
	OnFinalize   interface{}
	script       active.Runtime
	canRepeat    bool
}

func newExecutor(cn *Connection, db *sql.DB, query interface{}, canRepeat bool, script active.Runtime, onsuccess, onerror, onfinalize interface{}, args ...interface{}) (exctr *Executor) {
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
	exctr = &Executor{db: db, cn: cn, script: script, canRepeat: canRepeat, OnSuccess: onsuccess, OnError: onerror, OnFinalize: onfinalize}
	exctr.stmnt, exctr.argNames, exctr.mappedArgs = queryToStatement(exctr, query, args...)
	return
}

func (exctr *Executor) execute(forrows ...bool) (rws *sql.Rows, cltpes []*sql.ColumnType, cls []string) {
	if exctr.stmt, exctr.lasterr = exctr.db.Prepare(exctr.stmnt); exctr.lasterr == nil && exctr.stmt != nil {
		exctr.lastInsertID = -1
		exctr.rowsAffected = -1
		if exctr.canRepeat && len(exctr.argNames) > 0 {
			for agrn, argnme := range exctr.argNames {
				if prmv, prmvok := exctr.mappedArgs[argnme]; prmvok {
					parseParam(exctr, prmv, agrn)
				} else {
					parseParam(exctr, nil, agrn)
				}
			}
		}
		if len(forrows) >= 1 && forrows[0] {
			if rws, exctr.lasterr = exctr.stmt.Query(exctr.qryArgs...); rws != nil && exctr.lasterr == nil {
				cltpes, _ = rws.ColumnTypes()
				cls, _ = rws.Columns()
			} else if exctr.lasterr != nil {
				invokeError(exctr.script, exctr.lasterr, exctr.OnError)
			}
		} else {
			if rslt, rslterr := exctr.stmt.Exec(exctr.qryArgs...); rslterr == nil {
				if exctr.lastInsertID, rslterr = rslt.LastInsertId(); rslterr != nil {
					exctr.lastInsertID = -1
				}
				if exctr.rowsAffected, rslterr = rslt.RowsAffected(); rslterr != nil {
					exctr.rowsAffected = -1
				}
				invokeSuccess(exctr.script, exctr.OnSuccess, exctr)

			} else {
				exctr.lasterr = rslterr
				invokeError(exctr.script, exctr.lasterr, exctr.OnError)
			}
		}
	}
	return
}

//Repeat - repeat last query by repopulating parameters but dont regenerate last statement
func (exctr *Executor) Repeat(a ...interface{}) (err error) {
	exctr.execute()
	err = exctr.lasterr
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
		invokeFinalize(exctr.script, exctr.OnFinalize)
		if exctr.script != nil {
			exctr.script = nil
		}
		if exctr.OnSuccess != nil {
			exctr.OnSuccess = nil
		}
		if exctr.OnError != nil {
			exctr.OnError = nil
		}
		if exctr.OnFinalize != nil {
			exctr.OnFinalize = nil
		}
		if exctr.argNames != nil {
			exctr.argNames = nil
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

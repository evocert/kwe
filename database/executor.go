package database

import "database/sql"

//Executor - struct
type Executor struct {
	stmnt        string
	db           *sql.DB
	cn           *Connection
	stmt         *sql.Stmt
	lasterr      error
	lastInsertID int64
	rowsAffected int64
}

func newExecutor(cn *Connection, db *sql.DB, stmnt string) (exctr *Executor) {
	exctr = &Executor{stmnt: stmnt, db: db, cn: cn}
	return
}

func (exctr *Executor) execute(forrows ...bool) (rws *sql.Rows, cltpes []*sql.ColumnType, cls []string) {
	if exctr.stmt, exctr.lasterr = exctr.db.Prepare(exctr.stmnt); exctr.lasterr == nil && exctr.stmt != nil {
		exctr.lastInsertID = -1
		exctr.rowsAffected = 0
		if len(forrows) >= 1 && forrows[0] {
			if rws, exctr.lasterr = exctr.stmt.Query(); rws != nil && exctr.lasterr == nil {
				cltpes, _ = rws.ColumnTypes()
				cls, _ = rws.Columns()
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
		exctr = nil
	}
	return
}

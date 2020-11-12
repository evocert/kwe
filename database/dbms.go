package database

import (
	"database/sql"
	"encoding/json"
	"io"

	"github.com/evocert/kwe/iorw/active"
)

//DBMS - struct
type DBMS struct {
	cnctns  map[string]*Connection
	drivers map[string]func(string, ...interface{}) (*sql.DB, error)
}

//RegisterConnection - alias, driverName, dataSourceName
func (dbms *DBMS) RegisterConnection(alias string, driver string, datasource string) (registered bool) {
	if alias != "" && driver != "" && datasource != "" {
		if _, drvinvok := dbms.drivers[driver]; drvinvok {
			if cn, cnok := dbms.cnctns[alias]; cnok {
				if cn.driverName != driver {
					cn.driverName = driver
					cn.dataSourceName = datasource
					registered = true
				}
			} else if cn := NewConnection(dbms, driver, datasource); cn != nil {
				dbms.cnctns[alias] = cn
				registered = true
			}
		}
	}
	return

}

//RegisterDriver - register driver name for invokable db call
func (dbms *DBMS) RegisterDriver(driver string, invokedbcall func(string, ...interface{}) (*sql.DB, error)) {
	if driver != "" && invokedbcall != nil {
		dbms.drivers[driver] = invokedbcall
	}
}

func (dbms *DBMS) aliasExists(alias string) (exists bool, dbcn *Connection) {
	if alias != "" && len(dbms.cnctns) > 0 {
		dbcn, exists = dbms.cnctns[alias]
	}
	return
}

func (dbms *DBMS) driverDbInvoker(driver string) (dbinvoker func(string, ...interface{}) (*sql.DB, error), hasdbinvoker bool) {
	if driver != "" && len(dbms.drivers) > 0 {
		dbinvoker, hasdbinvoker = dbms.drivers[driver]
	}
	return
}

//QuerySettings - map[string]interface{} settings wrapper for Query
// settings :
// alias -  cn alias
// query -  statement
// args - [] slice of arguments
// success - func(r) event when ready
// error - func(error) event when encountering an error
// finalize - func() final wrapup event
// repeatable - true keep underlying stmnt open and allows for repeating query
// script - script handle
func (dbms *DBMS) QuerySettings(a interface{}) (reader *Reader) {
	if a != nil {
		if sttngs, sttngsok := a.(map[string]interface{}); sttngsok {
			var alias = ""
			var query interface{} = nil
			var prms = []interface{}{}
			var onsuccess interface{} = nil
			var onerror interface{} = nil
			var onfinalize interface{} = nil
			var script active.Runtime = nil
			var canRepeat = false
			var stngok = false
			var args []interface{} = nil
			for stngk, stngv := range sttngs {
				if stngk == "alias" {
					alias, _ = stngv.(string)
				} else if stngk == "query" {
					query = stngv
				} else if stngk == "repeatable" {
					if canRepeat, stngok = stngv.(bool); stngok {
						if canRepeat {
							prms = append(prms, stngv)
						}
					}
				} else if stngk == "success" {
					onsuccess = stngv
				} else if stngk == "error" {
					onerror = stngv
				} else if stngk == "finalize" {
					onfinalize = stngv
				} else if stngk == "script" {
					if script, stngok = stngv.(active.Runtime); stngok {
						prms = append(prms, script)
					}
				} else if stngk == "prms" || stngk == "args" {
					if args, stngok = stngv.([]interface{}); stngok && len(args) > 0 {
						prms = append(prms, args...)
					}
				}
			}
			if exists, dbcn := dbms.aliasExists(alias); exists {
				var err error = nil
				reader, _, err = dbcn.query(query, false, onsuccess, onerror, onfinalize, prms...)
				if err != nil && reader == nil {

				}
			}
		}
	}
	return
}

//Query - query database by alias - return Reader for underlying dataset
func (dbms *DBMS) Query(alias string, query interface{}, prms ...interface{}) (reader *Reader) {
	if exists, dbcn := dbms.aliasExists(alias); exists {
		var err error = nil
		reader, _, err = dbcn.query(query, false, nil, nil, nil, prms...)
		if err != nil && reader == nil {

		}
	}
	return
}

//InOut - OO{ in io.Reader -> out io.Writer } loop till no input
func (dbms *DBMS) InOut(in io.Reader, out io.Writer, outformat string) (err error) {
	var decoder *json.Decoder = nil
	if in != nil {
		decoder = json.NewDecoder(in)
		decoder.Token()
	}

	var data map[string]interface{} = nil
	for decoder != nil {
		if data == nil {
			data = map[string]interface{}{}
		}

	}
	return
}

//ExecuteSettings - map[string]interface{} settings wrapper for Execute
// settings :
// alias -  cn alias
// query -  statement
// args - [] slice of arguments
// success - func(r) event when ready
// error - func(error) event when encountering an error
// finalize - func() final wrapup event
// repeatable - true keep underlying stmnt open and allows for repeating query
// script - script handle
func (dbms *DBMS) ExecuteSettings(a interface{}) (exctr *Executor) {
	if sttngs, sttngsok := a.(map[string]interface{}); sttngsok {
		var alias = ""
		var query interface{} = nil
		var prms = []interface{}{}
		var onsuccess interface{} = nil
		var onerror interface{} = nil
		var onfinalize interface{} = nil
		var script active.Runtime = nil
		var canRepeat = false
		var stngok = false
		var args []interface{} = nil
		for stngk, stngv := range sttngs {
			if stngk == "alias" {
				alias, _ = stngv.(string)
			} else if stngk == "query" {
				query = stngv
			} else if stngk == "repeatable" {
				if canRepeat, stngok = stngv.(bool); stngok {
					if canRepeat {
						prms = append(prms, stngv)
					}
				}
			} else if stngk == "success" {
				onsuccess = stngv
			} else if stngk == "error" {
				onerror = stngv
			} else if stngk == "finalize" {
				onfinalize = stngv
			} else if stngk == "script" {
				if script, stngok = stngv.(active.Runtime); stngok {
					prms = append(prms, script)
				}
			} else if stngk == "prms" || stngk == "args" {
				if args, stngok = stngv.([]interface{}); stngok && len(args) > 0 {
					prms = append(prms, args...)
				}
			}
		}
		if exists, dbcn := dbms.aliasExists(alias); exists {
			var err error = nil
			if _, exctr, err = dbcn.query(query, true, onsuccess, onerror, onfinalize, prms...); err != nil {

			}
		}
	}
	return
}

//Execute - query database by alias - no result actions
func (dbms *DBMS) Execute(alias string, query interface{}, prms ...interface{}) (exctr *Executor) {
	if exists, dbcn := dbms.aliasExists(alias); exists {
		var err error = nil
		if _, exctr, err = dbcn.query(query, true, nil, nil, nil, prms...); err != nil {

		}
	}
	return
}

//NewDBMS - instance
func NewDBMS() (dbms *DBMS) {
	dbms = &DBMS{cnctns: map[string]*Connection{}, drivers: map[string]func(string, ...interface{}) (*sql.DB, error){}}

	return
}

var glbdbms *DBMS

//GLOBALDBMS - Global DBMS instance
func GLOBALDBMS() *DBMS {
	return glbdbms
}

func init() {
	if glbdbms == nil {
		glbdbms = NewDBMS()
	}
}

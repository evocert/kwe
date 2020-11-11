package database

import (
	"database/sql"
	"encoding/json"
	"io"
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

package database

import (
	"database/sql"
)

//DBMS - struct
type DBMS struct {
	cnctns  map[string]*Connection
	drivers map[string]func(string) (*sql.DB, error)
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
func (dbms *DBMS) RegisterDriver(driver string, invokedbcall func(string) (*sql.DB, error)) {
	if driver != "" && invokedbcall != nil {
		dbms.drivers[driver] = invokedbcall
	}
}

//NewDBMS - instance
func NewDBMS() (dbms *DBMS) {
	dbms = &DBMS{cnctns: map[string]*Connection{}, drivers: map[string]func(string) (*sql.DB, error){}}

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

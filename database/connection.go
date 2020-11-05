package database

//Connection - struct
type Connection struct {
	dbms                       *DBMS
	driverName, dataSourceName string
}

//NewConnection - dbms,driver name and datasource name (cn-string)
func NewConnection(dbms *DBMS, driverName, dataSourceName string) (cn *Connection) {
	cn = &Connection{dbms: dbms, driverName: driverName, dataSourceName: dataSourceName}
	return
}

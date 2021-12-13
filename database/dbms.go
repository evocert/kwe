package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/parameters"
)

type DBMSAPI interface {
	Connections() (cns []string)
	UnregisterConnection(alias string) (unregistered bool)
	RegisterConnection(alias string, driver string, datasource string, a ...interface{}) (registered bool)
	Exists(alias string) (exists bool)
	Info(alias ...string) map[string]interface{}
	Query(a interface{}, qryargs ...interface{}) (reader *Reader)
	Execute(a interface{}, excargs ...interface{}) (exctr *Executor)
	InOut(in interface{}, out io.Writer, ioargs ...interface{}) (err error)
	DriverName(alias string) string
}

type ActiveDBMS struct {
	dbms     *DBMS
	atvrntme active.Runtime
	prmsfnc  func() parameters.ParametersAPI
}

func (atvdbms *ActiveDBMS) Connections() (cns []string) {
	if atvdbms != nil && atvdbms.dbms != nil {
		cns = atvdbms.dbms.Connections()
	}
	return
}

func (atvdbms *ActiveDBMS) Info(alias ...string) (info map[string]interface{}) {
	if atvdbms != nil && atvdbms.dbms != nil {
		info = atvdbms.dbms.Info(alias...)
	}
	return
}

func (atvdbms *ActiveDBMS) DriverName(alias string) (driver string) {
	if atvdbms != nil && atvdbms.dbms != nil {
		driver = atvdbms.dbms.DriverName(alias)
	}
	return
}

func (atvdbms *ActiveDBMS) Dispose() {
	if atvdbms != nil {
		if atvdbms.atvrntme != nil {
			atvdbms.atvrntme = nil
		}
		if atvdbms.dbms != nil {
			atvdbms.dbms = nil
		}
		if atvdbms.prmsfnc != nil {
			atvdbms.prmsfnc = nil
		}
		atvdbms = nil
	}
}

func (atvdbms *ActiveDBMS) UnregisterConnection(alias string) (unregistered bool) {
	if atvdbms != nil && atvdbms.dbms != nil {
		unregistered = atvdbms.dbms.UnregisterConnection(alias)
	}
	return
}

func (atvdbms *ActiveDBMS) RegisterConnection(alias string, driver string, datasource string, a ...interface{}) (registered bool) {
	if atvdbms != nil && atvdbms.dbms != nil {
		if atvdbms.atvrntme != nil {
			a = append([]interface{}{atvdbms.atvrntme}, a...)
		}
		registered = atvdbms.dbms.RegisterConnection(alias, driver, datasource, a...)
	}
	return
}

func (atvdbms *ActiveDBMS) Exists(alias string) (exists bool) {
	if atvdbms != nil && atvdbms.dbms != nil {
		exists, _ = atvdbms.dbms.Exists(alias)
	}
	return
}

func (atvdbms *ActiveDBMS) Query(a interface{}, qryargs ...interface{}) (reader *Reader) {
	if atvdbms != nil && atvdbms.dbms != nil {
		if atvdbms.atvrntme != nil {
			qryargs = append([]interface{}{atvdbms.atvrntme}, qryargs...)
		}
		if atvdbms.prmsfnc != nil {
			if prms := atvdbms.prmsfnc(); prms != nil {
				qryargs = append([]interface{}{prms}, qryargs...)
			}
		}
		reader = atvdbms.dbms.Query(a, qryargs...)
	}
	return
}

func (atvdbms *ActiveDBMS) Execute(a interface{}, excargs ...interface{}) (exctr *Executor) {
	if atvdbms != nil && atvdbms.dbms != nil {
		if atvdbms.atvrntme != nil {
			excargs = append([]interface{}{atvdbms.atvrntme}, excargs...)
		}
		exctr = atvdbms.dbms.Execute(a, excargs...)
	}
	return
}

func (atvdbms *ActiveDBMS) InOut(in interface{}, out io.Writer, ioargs ...interface{}) (err error) {
	if atvdbms != nil && atvdbms.dbms != nil {
		if atvdbms.atvrntme != nil {
			ioargs = append([]interface{}{atvdbms.atvrntme}, ioargs...)
		}
		err = atvdbms.dbms.InOut(in, out, ioargs...)
	}
	return
}

func newActiveDBMS(dbms *DBMS, rntme active.Runtime, prmsfnc func() parameters.ParametersAPI) (atvdbms *ActiveDBMS) {
	if dbms != nil && rntme != nil {
		atvdbms = &ActiveDBMS{dbms: dbms, atvrntme: rntme, prmsfnc: prmsfnc}
	}
	return
}

//DBMS - struct
type DBMS struct {
	cnctns  map[string]*Connection
	drivers map[string]func(string, ...interface{}) (*sql.DB, error)
}

//ActiveDBMS return registered connections
func (dbms *DBMS) ActiveDBMS(rntme active.Runtime, prmsfnc func() parameters.ParametersAPI) (atvdbms *ActiveDBMS) {
	return newActiveDBMS(dbms, rntme, prmsfnc)
}

//Connection return registered connections
func (dbms *DBMS) Connection(alias string) (cn *Connection) {
	if alias = strings.TrimSpace(alias); alias != "" {
		cn = dbms.cnctns[alias]
	}
	return
}

//Connections return list of registered connection aliases
func (dbms *DBMS) Connections() (cns []string) {
	if cnsl := len(dbms.cnctns); cnsl > 0 {
		cns = make([]string, cnsl)
		cnsi := 0
		for cnsk := range dbms.cnctns {
			cns[cnsi] = cnsk
			cnsi++
			if cnsi == cnsl {
				break
			}
		}
	}
	return
}

//Info return status info of all aliases or alias... provided
func (dbms *DBMS) Info(alias ...string) (info map[string]interface{}) {
	if dbms != nil {
		if cnsl := len(dbms.cnctns); cnsl > 0 {
			info = map[string]interface{}{}
			if aliasln := len(alias); aliasln > 0 {
				for _, cnsk := range alias {
					if cnsk != "" {
						if cn := dbms.cnctns[cnsk]; cn != nil {
							info[cnsk] = cn.Info()
						}
					}
				}
			} else {
				for cnsk := range dbms.cnctns {
					if cn := dbms.cnctns[cnsk]; cn != nil {
						info[cnsk] = cn.Info()
					}
				}
			}
		}
	}
	return
}

func (dbms *DBMS) DriverName(alias string) (driver string) {
	if dbms != nil {
		if cn, cnok := dbms.cnctns[alias]; cnok {
			driver = cn.driverName
		}
	}
	return
}

//UnegisterConnection - alias
func (dbms *DBMS) UnregisterConnection(alias string) (unregistered bool) {
	if alias != "" {
		if cn, cnok := dbms.cnctns[alias]; cnok {
			cn.Dispose()
			dbms.cnctns[alias] = nil
			delete(dbms.cnctns, alias)
			unregistered = true
		}
	}
	return
}

//RegisterConnection - alias, driverName, dataSourceName
func (dbms *DBMS) RegisterConnection(alias string, driver string, datasource string, a ...interface{}) (registered bool) {
	if alias != "" && driver != "" && datasource != "" {
		if strings.HasPrefix(datasource, "http://") || strings.HasPrefix(datasource, "https://") || strings.HasPrefix(datasource, "ws://") || strings.HasPrefix(datasource, "wss://") {
			if cn, cnok := dbms.cnctns[alias]; cnok {
				if cn.driverName != driver {
					cn.driverName = driver
					cn.dataSourceName = datasource
					calibrateConnection(cn, a...)
					registered = true
				}
			} else if cn := NewConnection(dbms, driver, datasource); cn != nil {
				dbms.cnctns[alias] = cn
				calibrateConnection(cn, a...)
				registered = true
			}
		} else if _, drvinvok := dbms.drivers[driver]; drvinvok {
			if cn, cnok := dbms.cnctns[alias]; cnok {
				if cn.driverName != driver {
					cn.driverName = driver
					cn.dataSourceName = datasource
					registered = true
					calibrateConnection(cn, a...)
				}
			} else if cn := NewConnection(dbms, driver, datasource); cn != nil {
				dbms.cnctns[alias] = cn
				calibrateConnection(cn, a...)
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

//Exists - alias exist <= exist[true], dbcn[*Connection]
func (dbms *DBMS) Exists(alias string) (exists bool, dbcn *Connection) {
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

//Query - map[string]interface{} settings wrapper for Query
// settings :
// alias -  cn alias
// query -  statement
// args - [] slice of arguments
// success - func(r) event when ready
// error - func(error) event when encountering an error
// finalize - func() final wrapup event
// repeatable - true keep underlying stmnt open and allows for repeating query
// script - script handle
func (dbms *DBMS) Query(a interface{}, qryargs ...interface{}) (reader *Reader) {
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
			var execargs []map[string]interface{} = nil
			var argsmap map[string]interface{} = nil
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
				} else if stngk == "exec" {
					if execargsv, _ := stngv.([]interface{}); len(execargsv) > 0 {
						for _, execstngv := range execargsv {
							if execstngv != nil {
								if execmpv, _ := execstngv.(map[string]interface{}); execmpv != nil && len(execmpv) > 0 {
									if execargs == nil {
										execargs = []map[string]interface{}{}
									}
									execargs = append(execargs, execmpv)
								}
							}
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
					} else if argsmap, stngok = stngv.(map[string]interface{}); stngok && len(argsmap) > 0 {
						prms = append(prms, argsmap)
					}
				}
			}
			if len(qryargs) > 0 {
				prms = append(prms, qryargs...)
			}
			if exists, dbcn := dbms.Exists(alias); exists {
				var err error = nil
				reader, _, err = internquery(dbcn, query, false, execargs, onsuccess, onerror, onfinalize, prms...)
				if err != nil && reader == nil {

				}
			}
		}
	}
	return
}

func (dbms *DBMS) QueryJSON(query interface{}, prms ...interface{}) (reader *Reader) {
	var err error = nil
	reader, _, err = internquery(nil, query, false, nil, nil, nil, nil, prms...)
	if err != nil && reader == nil {

	}
	return
}

/*//Query - query database by alias - return Reader for underlying dataset
func (dbms *DBMS) Query(alias string, query interface{}, prms ...interface{}) (reader *Reader) {
	if exists, dbcn := dbms.AliasExists(alias); exists {
		var err error = nil
		reader, _, err = internquery(dbcn, query, false, nil, nil, nil, prms...)
		if err != nil && reader == nil {

		}
	}
	return
}*/

func (dbms *DBMS) inReaderOut(ri io.Reader, out io.Writer, ioargs ...interface{}) (hasoutput bool, err error) {
	if ri != nil {
		func() {
			var buff = iorw.NewBuffer()
			defer buff.Close()
			buffl, bufferr := io.Copy(buff, ri)
			if bufferr == nil || bufferr == io.EOF {
				if buffl > 0 {
					func() {
						var buffr = buff.Reader()
						defer func() {
							buffr.Close()
						}()
						d := json.NewDecoder(buffr)
						rqstmp := map[string]interface{}{}
						if jsnerr := d.Decode(&rqstmp); jsnerr == nil {
							if len(rqstmp) > 0 {
								hasoutput, err = dbms.inMapOut(rqstmp, out, ioargs...)
							}
						} else {
							err = jsnerr
						}
					}()
				}
			}
		}()
	}
	return
}

func (dbms *DBMS) inMapOut(mpin map[string]interface{}, out io.Writer, ioargs ...interface{}) (hasoutput bool, err error) {
	if mpl := len(mpin); mpl > 0 {
		if out != nil {
			hasoutput = true
			iorw.Fprint(out, "{")
		}
		var dfltalias string = ""
		var dfltcn, crntcn *Connection = nil, nil
		if aliasv, aliasok := mpin["alias"]; aliasok {
			if dfltalias == "" {
				if aliass, aliassok := aliasv.(string); aliassok {
					dfltalias = aliass
					if dfcnok, dfcn := dbms.Exists(aliass); dfcnok {
						dfltcn = dfcn
					}
				}
			}
			mpl--
			delete(mpin, "alias")
		}
		for mk, mv := range mpin {
			mpl--
			if out != nil {
				hasoutput = true
				iorw.Fprint(out, "\""+mk+"\":")
			}
			if mvp, mvpok := mv.(map[string]interface{}); mvpok {
				crntcn = nil
				if dalias, daliasok := mvp["alias"]; daliasok && dalias != nil {
					if salias, saliasok := dalias.(string); saliasok && salias != "" {
						if cnok, cn := dbms.Exists(salias); cnok {
							crntcn = cn
						}
					}
					if crntcn == nil {
						if out != nil {
							hasoutput = true
							jsnrdr := NewJSONReader(nil, nil, fmt.Errorf("alias does not exist"))
							io.Copy(out, jsnrdr)
							jsnrdr = nil
						}
					}
				} else {
					if crntcn == nil && dfltcn != nil {
						crntcn = dfltcn
					} else {
						if dfltalias != "" {
							if out != nil {
								hasoutput = true
								jsnrdr := NewJSONReader(nil, nil, fmt.Errorf("default alias does not exist"))
								io.Copy(out, jsnrdr)
								jsnrdr = nil
							}
						} else {
							if out != nil {
								hasoutput = true
								jsnrdr := NewJSONReader(nil, nil, fmt.Errorf("no alias"))
								io.Copy(out, jsnrdr)
								jsnrdr = nil
							}
						}
					}
				}

				if crntcn != nil {
					if cmd, cmdok := mvp["execute"]; cmdok {
						delete(mvp, "execute")
						exctr, exctrerr := crntcn.GblExecute(cmd, mvp, ioargs)
						if out != nil {
							hasoutput = true
							jsnrdr := NewJSONReader(nil, exctr, exctrerr)
							io.Copy(out, jsnrdr)
							jsnrdr = nil
						}
					} else if cmd, cmdok := mvp["query"]; cmdok {
						delete(mvp, "query")
						rdr, rdrerr := crntcn.GblQuery(cmd, mvp, ioargs)
						if out != nil {
							hasoutput = true
							jsnrdr := NewJSONReader(rdr, nil, rdrerr)
							io.Copy(out, jsnrdr)
							jsnrdr = nil
						}
					} else {
						if out != nil {
							hasoutput = true
							jsnrdr := NewJSONReader(nil, nil, fmt.Errorf("no request"))
							io.Copy(out, jsnrdr)
							jsnrdr = nil
						}
					}
				}
			} else {
				if out != nil {
					hasoutput = true
					jsnrdr := NewJSONReader(nil, nil, fmt.Errorf("invalid request"))
					io.Copy(out, jsnrdr)
					jsnrdr = nil
				}
			}
			if mpl >= 1 {
				if out != nil {
					hasoutput = true
					iorw.Fprint(out, ",")
				}
			}
		}
		if out != nil {
			hasoutput = true
			iorw.Fprint(out, "}")
		}
	}
	return
}

//InOutS - OO{ in io.Reader -> out string } loop till no input
func (dbms *DBMS) InOutS(in interface{}, ioargs ...interface{}) (out string, err error) {
	var buff = iorw.NewBuffer()
	defer buff.Close()
	err = dbms.InOut(in, buff)
	out = buff.String()
	return
}

//InOut - OO{ in io.Reader -> out io.Writer } loop till no input
func (dbms *DBMS) InOut(in interface{}, out io.Writer, ioargs ...interface{}) (err error) {
	if in != nil {
		var hasoutput = false
		if mp, mpok := in.(map[string]interface{}); mpok {
			hasoutput, err = dbms.inMapOut(mp, out, ioargs...)
		} else if ri, riok := in.(io.Reader); riok && ri != nil {
			hasoutput, err = dbms.inReaderOut(ri, out, ioargs...)
		} else if si, siok := in.(string); siok && si != "" {
			hasoutput, err = dbms.inReaderOut(strings.NewReader(si), out, ioargs...)
		}
		if !hasoutput {
			if out != nil {
				if err != nil {
					iorw.Fprint(out, "{\"error\":\""+err.Error()+"\"}")
				} else {
					iorw.Fprint(out, "{}")
				}
			}
		}
	} else {
		if out != nil {
			iorw.Fprint(out, "{}")
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
func (dbms *DBMS) Execute(a interface{}, excargs ...interface{}) (exctr *Executor) {
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
		var argsmap map[string]interface{} = nil
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
				} else if argsmap, stngok = stngv.(map[string]interface{}); stngok && len(argsmap) > 0 {
					prms = append(prms, argsmap)
				}
			}
		}
		if len(excargs) > 0 {
			prms = append(prms, excargs...)
		}
		if exists, dbcn := dbms.Exists(alias); exists {
			var err error = nil
			if _, exctr, err = internquery(dbcn, query, true, nil, onsuccess, onerror, onfinalize, prms...); err != nil {

			}
		}
	}
	return
}

/*//Execute - query database by alias - no result actions
func (dbms *DBMS) Execute(alias string, query interface{}, prms ...interface{}) (exctr *Executor) {
	if exists, dbcn := dbms.AliasExists(alias); exists {
		var err error = nil
		if _, exctr, err = internquery(dbcn, query, true, nil, nil, nil, prms...); err != nil {

		}
	}
	return
}*/

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

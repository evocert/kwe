package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
)

//DBMS - struct
type DBMS struct {
	cnctns  map[string]*Connection
	drivers map[string]func(string, ...interface{}) (*sql.DB, error)
}

//RegisterConnection - alias, driverName, dataSourceName
func (dbms *DBMS) RegisterConnection(alias string, driver string, datasource string, a ...interface{}) (registered bool) {
	if alias != "" && driver != "" && datasource != "" {
		if strings.HasPrefix(datasource, "http://") || strings.HasPrefix(datasource, "https://") || strings.HasPrefix(datasource, "ws://") || strings.HasPrefix(datasource, "wss://") {
			if cn, cnok := dbms.cnctns[alias]; cnok {
				if cn.driverName != driver {
					cn.driverName = driver
					cn.dataSourceName = datasource
					registered = true
					cn.endpnt = newEndPoint(cn.dataSourceName, a...)
				}
			} else if cn := NewConnection(dbms, driver, datasource); cn != nil {
				cn.endpnt = newEndPoint(cn.dataSourceName, a...)
				dbms.cnctns[alias] = cn
				registered = true
			}
		} else if _, drvinvok := dbms.drivers[driver]; drvinvok {
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

//AliasExists - alias exist <= exist[true], dbcn[*Connection]
func (dbms *DBMS) AliasExists(alias string) (exists bool, dbcn *Connection) {
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
			if exists, dbcn := dbms.AliasExists(alias); exists {
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
	if exists, dbcn := dbms.AliasExists(alias); exists {
		var err error = nil
		reader, _, err = dbcn.query(query, false, nil, nil, nil, prms...)
		if err != nil && reader == nil {

		}
	}
	return
}

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
					if dfcnok, dfcn := dbms.AliasExists(aliass); dfcnok {
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
						if cnok, cn := dbms.AliasExists(salias); cnok {
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
		if exists, dbcn := dbms.AliasExists(alias); exists {
			var err error = nil
			if _, exctr, err = dbcn.query(query, true, onsuccess, onerror, onfinalize, prms...); err != nil {

			}
		}
	}
	return
}

//Execute - query database by alias - no result actions
func (dbms *DBMS) Execute(alias string, query interface{}, prms ...interface{}) (exctr *Executor) {
	if exists, dbcn := dbms.AliasExists(alias); exists {
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

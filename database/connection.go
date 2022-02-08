package database

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/parameters"
)

//Connection - struct
type Connection struct {
	dbms                       *DBMS
	driverName, dataSourceName string
	dbi                        interface{}
	db                         *sql.DB
	args                       []interface{}
	dbinvoker                  func(string, ...interface{}) (*sql.DB, error)
	lastmaxidecons             int
	maxidlecons                int
	lastmaxopencons            int
	maxopencons                int
}

func (cn *Connection) Info() (info map[string]interface{}) {
	if cn != nil {
		info = map[string]interface{}{}
		info["driver"] = cn.Driver()
		info["connected"] = cn.IsConnected()
		info["datasource"] = cn.dataSourceName
		if db := cn.db; db != nil {
			var dbstats = db.Stats()
			info["stats"] = dbstats
		}
	}
	return
}

func (cn *Connection) SetMaxIdleConns(idlcons int) {
	if cn != nil {
		cn.maxidlecons = idlcons
		if cn.db != nil {
			if cn.lastmaxidecons != cn.maxidlecons {
				cn.db.SetMaxIdleConns(cn.maxidlecons)
				cn.lastmaxidecons = cn.maxidlecons
			}
		}
	}
}

func (cn *Connection) SetMaxOpenConns(opencons int) {
	if cn != nil {
		cn.maxopencons = opencons
		if cn.db != nil {
			if cn.lastmaxidecons != cn.maxopencons {
				cn.db.SetMaxOpenConns(cn.maxopencons)
				cn.lastmaxidecons = cn.maxopencons
			}
		}
	}
}

func (cn *Connection) IsRemote() bool {
	return (strings.HasPrefix(cn.dataSourceName, "http://") || strings.HasPrefix(cn.dataSourceName, "https://") || strings.HasPrefix(cn.dataSourceName, "ws://") || strings.HasPrefix(cn.dataSourceName, "wss://"))
}

func (cn *Connection) Dispose() (err error) {
	if cn != nil {
		if cn.db != nil {
			err = cn.db.Close()
			cn.db = nil
		}
		if cn.args != nil {
			cn.args = nil
		}
		if cn.dbi != nil {
			cn.dbi = nil
		}
		if cn.dbinvoker != nil {
			cn.dbinvoker = nil
		}
		if cn.dbms != nil {
			cn.dbms = nil
		}
		cn = nil
	}
	return
}

func runeReaderToString(rnr io.RuneReader) (s string) {
	return
}

func (cn *Connection) Ping() (status map[string]interface{}) {
	status = map[string]interface{}{}
	if cn.db != nil {
		pngerr := cn.db.Ping()
		if pngerr == nil {
			status["status"] = "Ok"
		} else {
			status["status"] = "Failed"
			status["error"] = pngerr.Error()
		}
	}
	return
}

func (cn *Connection) DataSource() (datasource string) {
	if cn != nil {
		datasource = cn.dataSourceName
	}
	return
}

func (cn *Connection) Driver() (driver string) {
	if cn != nil {
		driver = cn.driverName
	}
	return
}

func (cn *Connection) IsConnected() (connected bool) {
	if cn.db != nil {
		pngerr := cn.db.Ping()
		connected = pngerr == nil
	}
	return
}

func parseParam(exctr *Executor, prmval interface{}, argi int) (s string) {
	if exctr.cn.driverName == "sqlserver" {
		if argi == -1 {
			s = ("@p" + fmt.Sprintf("%d", len(exctr.qryArgs)))
			prmnme := "p" + fmt.Sprintf("%d", len(exctr.qryArgs))
			exctr.qryArgs = append(exctr.qryArgs, sql.Named(prmnme, prmval))
		} else {
			prmnme := "p" + fmt.Sprintf("%d", argi)
			exctr.qryArgs[argi] = sql.Named(prmnme, prmval)
		}
	} else if exctr.cn.driverName == "postgres" || exctr.cn.driverName == "sqlite3" || exctr.cn.driverName == "kwesqlite" {
		if argi == -1 {
			s = ("$" + fmt.Sprintf("%d", len(exctr.qryArgs)+1))
		}
		if argvs, argvsok := prmval.(string); argvsok {
			argvs = strings.TrimSpace(argvs)
			/*if fltval, nrerr := strconv.ParseFloat(argvs, 64); nrerr == nil {
				if tstintval := int64(fltval); float64(tstintval) == fltval {
					if argi == -1 {
						exctr.qryArgs = append(exctr.qryArgs, tstintval)
					} else {
						exctr.qryArgs[argi] = tstintval
					}
				} else {
					if argi == -1 {
						exctr.qryArgs = append(exctr.qryArgs, fltval)
					} else {
						exctr.qryArgs[argi] = fltval
					}
				}
			} else if intval, nrerr := strconv.ParseInt(argvs, 10, 64); nrerr == nil {
				if argi == -1 {
					exctr.qryArgs = append(exctr.qryArgs, intval)
				} else {
					exctr.qryArgs[argi] = intval
				}
			} else {
				if argi == -1 {
					exctr.qryArgs = append(exctr.qryArgs, argvs)
				} else {
					exctr.qryArgs[argi] = argvs
				}
			}*/
			if argi == -1 {
				exctr.qryArgs = append(exctr.qryArgs, argvs)
			} else {
				exctr.qryArgs[argi] = argvs
			}
		} else if argvb, argvsok := prmval.(bool); argvsok {
			if argi == -1 {
				exctr.qryArgs = append(exctr.qryArgs, argvb)
			} else {
				exctr.qryArgs[argi] = argvb
			}
		} else {
			if argi == -1 {
				exctr.qryArgs = append(exctr.qryArgs, prmval)
			} else {
				exctr.qryArgs[argi] = prmval
			}
		}
	} else if exctr.cn.driverName == "oracle" {
		if argi == -1 {
			s = (":" + fmt.Sprintf("%d", len(exctr.qryArgs)+1))
		}
		if argvs, argvsok := prmval.(string); argvsok {
			argvs = strings.TrimSpace(argvs)
			/*if fltval, nrerr := strconv.ParseFloat(argvs, 64); nrerr == nil {
				if tstintval := int64(fltval); float64(tstintval) == fltval {
					if argi == -1 {
						exctr.qryArgs = append(exctr.qryArgs, tstintval)
					} else {
						exctr.qryArgs[argi] = tstintval
					}
				} else {
					if argi == -1 {
						exctr.qryArgs = append(exctr.qryArgs, fltval)
					} else {
						exctr.qryArgs[argi] = fltval
					}
				}
			} else if intval, nrerr := strconv.ParseInt(argvs, 10, 64); nrerr == nil {
				if argi == -1 {
					exctr.qryArgs = append(exctr.qryArgs, intval)
				} else {
					exctr.qryArgs[argi] = intval
				}
			} else {
				if argi == -1 {
					exctr.qryArgs = append(exctr.qryArgs, argvs)
				} else {
					exctr.qryArgs[argi] = argvs
				}
			}*/
			if argi == -1 {
				exctr.qryArgs = append(exctr.qryArgs, argvs)
			} else {
				exctr.qryArgs[argi] = argvs
			}
		} else if argvb, argvsok := prmval.(bool); argvsok {
			if argi == -1 {
				exctr.qryArgs = append(exctr.qryArgs, argvb)
			} else {
				exctr.qryArgs[argi] = argvb
			}
		} else {
			if argi == -1 {
				exctr.qryArgs = append(exctr.qryArgs, prmval)
			} else {
				exctr.qryArgs[argi] = prmval
			}
		}
	} else {
		if argi == -1 {
			exctr.qryArgs = append(exctr.qryArgs, prmval)
			s = "?"
		} else {
			exctr.qryArgs[argi] = prmval
		}
	}
	return
}

func queryToStatement(exctr *Executor, query interface{}, args ...interface{}) (stmnt string, validNames []string, mappedVals map[string]interface{}) {
	var rnrr io.RuneReader = nil
	var sqlbuf *iorw.Buffer = nil
	if qrys, qrysok := query.(string); qrysok && qrys != "" {
		rnrr = bufio.NewReader(strings.NewReader(qrys))
	} else if qryrnr, qryrnrok := query.(io.RuneReader); qryrnrok {
		rnrr = qryrnr
	} else if qryr, qryrok := query.(io.Reader); qryrok {
		rnrr = bufio.NewReader(qryr)
	}

	mappedVals = map[string]interface{}{}
	var foundTxt = false

	var rdr *Reader = nil
	for len(args) > 0 {
		if pargs, ispargs := args[0].(*parameters.Parameters); ispargs && pargs != nil {
			for _, skey := range pargs.StandardKeys() {
				mappedVals[skey] = strings.Join(pargs.Parameter(skey), "")
			}
		} else if rdrargs, isrdrargs := args[0].(*Reader); isrdrargs && rdrargs != nil {
			rdr = rdrargs

			if cols := rdr.Columns(); len(cols) > 0 {
				data := rdr.Data()
				if len(data) == len(cols) {
					for cn := range cols {
						mappedVals[cols[cn]] = data[cn]
					}
				}
				data = nil
				cols = nil
			}
		} else if pmargs, ispmargs := args[0].(map[string]interface{}); ispmargs {
			for pmk := range pmargs {
				if mpv, mpvok := pmargs[pmk].(map[string]interface{}); mpvok && mpv != nil && len(mpv) > 0 {

				} else {
					mappedVals[pmk] = pmargs[pmk]
				}
			}
		}
		args = args[1:]
	}
	if len(exctr.qryArgs) == 0 {
		exctr.qryArgs = []interface{}{}
	}

	stmnt = ""

	var prvr = rune(0)
	var prmslbl = [][]rune{[]rune("@@"), []rune("@@")}
	var prmslbli = []int{0, 0}

	var appr = func(r rune) {
		if sqlbuf == nil {
			sqlbuf = iorw.NewBuffer()
		}
		sqlbuf.Print(string(r))
	}

	var apprs = func(p []rune) {
		if pl := len(p); pl > 0 {
			if sqlbuf == nil {
				sqlbuf = iorw.NewBuffer()
			}
			sqlbuf.Print(string(p))
		}
	}

	if len(mappedVals) == 0 {
		stmnt, _ = iorw.ReaderToString(rnrr)
	} else {
		var psblprmnme = make([]rune, 8192)
		var psblprmnmei = 0
		iorw.ReadRunesEOFFunc(rnrr, func(r rune) error {
			if foundTxt {
				appr(r)
				if r == '\'' {
					foundTxt = false
					prvr = rune(0)
				} else {
					prvr = r
				}
			} else {
				if prmslbli[1] == 0 && prmslbli[0] < len(prmslbl[0]) {
					if prmslbli[0] > 0 && prmslbl[0][prmslbli[0]-1] == prvr && prmslbl[0][prmslbli[0]] != r {
						if prmsl := prmslbli[0]; prmsl > 0 {
							prmslbli[0] = 0
							apprs(prmslbl[0][:prmsl])
						}
					}
					if prmslbl[0][prmslbli[0]] == r {
						prmslbli[0]++
						if prmslbli[0] == len(prmslbl[0]) {

							prvr = rune(0)
						} else {
							prvr = r
						}
					} else {
						if prmsl := prmslbli[0]; prmsl > 0 {
							prmslbli[0] = 0
							apprs(prmslbl[0][:prmsl])
						}
						appr(r)
						if r == '\'' {
							foundTxt = true
							prvr = rune(0)
						} else {
							prvr = r
						}
					}
				} else if prmslbli[0] == len(prmslbl[0]) && prmslbli[1] < len(prmslbl[1]) {
					if prmslbl[1][prmslbli[1]] == r {
						prmslbli[1]++
						if prmslbli[1] == len(prmslbl[1]) {
							if psblprmnmei > 0 {
								if psbprmnme := string(psblprmnme[:psblprmnmei]); psbprmnme != "" {
									fndprm := false
									if !exctr.isRemote() {
										for mpvk := range mappedVals {
											mpv := mappedVals[mpvk]
											if fndprm = strings.ToUpper(psbprmnme) == strings.ToUpper(mpvk); fndprm {
												if validNames == nil {
													validNames = []string{}
												}
												validNames = append(validNames, mpvk)
												apprs([]rune(parseParam(exctr, mpv, -1)))
												break
											}
										}
									}
									if !fndprm {
										apprs(prmslbl[0])
										apprs(psblprmnme[:psblprmnmei])
										apprs(prmslbl[1])
									}
								} else {
									apprs(prmslbl[0])
									apprs(prmslbl[1])
								}
								psblprmnmei = 0
							} else {
								apprs(prmslbl[0])
								apprs(prmslbl[1])
							}
							prmslbli[1] = 0
							prvr = rune(0)
							prmslbli[0] = 0
						}
					} else {
						if prmsl := prmslbli[1]; prmsl > 0 {
							//Invalid End Parameter
							prmslbli[1] = 0
							prvr = rune(0)
							prmslbli[0] = 0
							apprs(prmslbl[0])
							if psblprmnmei > 0 {
								apprs(psblprmnme[:psblprmnmei])
								psblprmnmei = 0
							}
							apprs(prmslbl[1][:prmsl])
						} else {
							psblprmnme[psblprmnmei] = r
							psblprmnmei++
							prvr = r
							if psblprmnmei == len(psblprmnme) {
								//Invalid Parameter Length
								prmslbli[1] = 0
								prvr = rune(0)
								prmslbli[0] = 0
								apprs(prmslbl[0])
								if psblprmnmei > 0 {
									apprs(psblprmnme[:psblprmnmei])
									psblprmnmei = 0
								}
							}
						}
					}
				}
			}
			return nil
		})

		if sqlbuf != nil {
			if sqlbuf.Size() > 0 {
				stmnt = sqlbuf.String()
			}
			sqlbuf.Close()
			sqlbuf = nil
		} else {
			stmnt = ""
		}
	}
	return
}

//GblExecute - public for query()*Executor
func (cn *Connection) GblExecute(query interface{}, prms ...interface{}) (exctr *Executor, err error) {
	if _, exctr, err = internquery(cn, query, true, nil, nil, nil, nil, prms...); err != nil {

	}
	return
}

//GblQuery - public for query() *Reader
func (cn *Connection) GblQuery(query interface{}, prms ...interface{}) (reader *Reader, err error) {
	reader, _, err = internquery(cn, query, false, nil, nil, nil, nil, prms...)
	if err != nil && reader == nil {

	}
	return
}

func (cn *Connection) inReaderOut(rin io.Reader, out io.Writer, ioargs ...interface{}) (hasoutput bool, err error) {
	if rin != nil {
		func() {
			var buff = iorw.NewBuffer()
			defer buff.Close()
			buff.Print(rin)
			buffl := buff.Size()

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
							hasoutput, err = cn.inMapOut(rqstmp, out, ioargs...)
						}
					} else {
						err = jsnerr
					}
				}()
			}
			//}
		}()
	}
	return
}

func (cn *Connection) inMapOut(mpin map[string]interface{}, out io.Writer, ioargs ...interface{}) (hasoutput bool, err error) {
	if mpin != nil {
		if mpl := len(mpin); mpl > 0 {
			if out != nil {
				hasoutput = true
				iorw.Fprint(out, "{")
			}
			for mk := range mpin {
				mv := mpin[mk]
				mpl--
				if out != nil {
					hasoutput = true
					iorw.Fprint(out, "\""+mk+"\":")
				}
				if mvp, mvpok := mv.(map[string]interface{}); mvpok {
					if cmd, cmdok := mvp["execute"]; cmdok {
						delete(mvp, "execute")
						exctr, exctrerr := cn.GblExecute(cmd, mvp, ioargs)
						if out != nil {
							hasoutput = true
							jsnrdr := NewJSONReader(nil, exctr, exctrerr)
							io.Copy(out, jsnrdr)
							jsnrdr = nil
						}
					} else if cmd, cmdok := mvp["query"]; cmdok {
						delete(mvp, "query")
						rdr, rdrerr := cn.GblQuery(cmd, mvp, ioargs)
						if out != nil {
							hasoutput = true
							jsnrdr := NewJSONReader(rdr, nil, rdrerr)
							io.Copy(out, jsnrdr)
							jsnrdr = nil
						}
						if rdr != nil {
							rdr.Close()
						}
					} else {
						if out != nil {
							hasoutput = true
							jsnrdr := NewJSONReader(nil, nil, fmt.Errorf("no request"))
							io.Copy(out, jsnrdr)
							jsnrdr = nil
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
	}
	return
}

//InOut - OO{ in interface{} -> out io.Writer } loop till no input
func (cn *Connection) InOut(in interface{}, out io.Writer, ioargs ...interface{}) {
	if in != nil {
		var hasoutput = false
		var ioerr error = nil
		if mp, mpok := in.(map[string]interface{}); mpok {
			hasoutput, ioerr = cn.inMapOut(mp, out, ioargs...)
		} else if mr, mrok := in.(io.Reader); mrok && mr != nil {
			hasoutput, ioerr = cn.inReaderOut(mr, out, ioargs...)
		} else if si, siok := in.(string); siok && si != "" {
			hasoutput, ioerr = cn.inReaderOut(strings.NewReader(si), out, ioargs...)
		}
		if !hasoutput {
			if out != nil {
				if ioerr != nil {
					iorw.Fprint(out, "{\"error\":\""+ioerr.Error()+"\"}")
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
}

func internquery(cn *Connection, query interface{}, noreader bool, execargs []map[string]interface{}, onsuccess, onerror, onfinalize interface{}, args ...interface{}) (reader *Reader, exctr *Executor, err error) {
	var argsn = 0
	var script active.Runtime = nil
	var canRepeat = false
	for argsn < len(args) {
		var d = args[argsn]
		if _, dok := d.(*parameters.Parameters); dok {
			argsn++
		} else if _, dok := d.(*Reader); dok {
			argsn++
		} else if _, dok := d.(map[string]interface{}); dok {
			argsn++
		} else if _, dok := d.(bool); dok {
			argsn++
		} else if dbool, dok := d.(bool); dok {
			canRepeat = dbool
			if argsn == len(args) {
				args = append(args[:argsn])
			} else if argsn < len(args) {
				args = append(args[:argsn], args[argsn+1:]...)
			}
		} else {
			if d != nil {
				if scrpt, scrptok := d.(active.Runtime); scrptok {
					if script == nil && scrpt != nil {
						script = scrpt
					}
				} else if onsuccess == nil {
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

	if cn.db == nil && !cn.IsRemote() {
		if cn.dbinvoker == nil {
			if dbinvoker, hasdbinvoker := cn.dbms.driverDbInvoker(cn.driverName); hasdbinvoker {
				cn.dbinvoker = dbinvoker
			}
		}
		if cn.dbinvoker != nil {
			if cn.dbi, err = cn.dbinvoker(cn.dataSourceName); err == nil && cn.dbi != nil {
				cn.db, _ = cn.dbi.(*sql.DB)
			}
			if err != nil && onerror != nil {
				invokeError(script, err, onerror)
			}
		}
	}
	if cn.db != nil {
		if cn.lastmaxidecons != cn.maxidlecons {
			cn.db.SetMaxIdleConns(cn.maxidlecons)
			cn.lastmaxidecons = cn.maxidlecons
		}
		if cn.lastmaxopencons != cn.maxopencons {
			cn.db.SetMaxOpenConns(cn.maxopencons)
			cn.lastmaxopencons = cn.maxopencons
		}

		if err = cn.db.Ping(); err != nil {
			cn.db.Close()
			cn.db = nil
			if cn.dbinvoker == nil {
				if dbinvoker, hasdbinvoker := cn.dbms.driverDbInvoker(cn.driverName); hasdbinvoker {
					cn.dbinvoker = dbinvoker
				}
			}
			if cn.dbi, err = cn.dbinvoker(cn.dataSourceName); err == nil && cn.dbi != nil {
				if cn.db, _ = cn.dbi.(*sql.DB); cn.db != nil {
					cn.db.Close()
				}
				if cn.dbi, err = cn.dbinvoker(cn.dataSourceName); err == nil && cn.dbi != nil {
					cn.db, _ = cn.dbi.(*sql.DB)
				}
			}
			if err != nil && onerror != nil {
				invokeError(script, err, onerror)
			}
			if err == nil {
				if err = cn.db.Ping(); err != nil {
					invokeError(script, err, onerror)
				}
			}
		}
		if err == nil {
			if query != nil {
				exctr = newExecutor(cn, cn.db, query, canRepeat, script, onsuccess, onerror, onfinalize, args...)
				if noreader {
					exctr.execute(false)
					if err = exctr.lasterr; err != nil {
						invokeError(exctr.script, err, onerror)
						exctr.Close()
						exctr = nil
					}
				} else {
					reader = newReader(exctr)
					reader.execute()
					if err = reader.lasterr; err != nil {
						invokeError(reader.script, err, onerror)
						reader.Close()
						reader = nil
					}
				}
			}
		}
	} else if cn.IsRemote() {
		if query != nil {
			exctr = newExecutor(cn, cn.db, query, canRepeat, script, onsuccess, onerror, onfinalize, args...)
			if noreader {
				exctr.execute(false)
				if err = exctr.lasterr; err != nil {
					invokeError(exctr.script, err, onerror)
					exctr.Close()
					exctr = nil
				}
			} else {
				reader = newReader(exctr)
				if len(execargs) > 0 {
					for execmpi := range execargs {
						if execmp := execargs[execmpi]; len(execmp) > 0 {

						}
					}
				}
				reader.execute()
				if err = reader.lasterr; err != nil {
					invokeError(reader.script, err, onerror)
					reader.Close()
					reader = nil
				}
			}
		}
	}
	return
}

func invokeError(script active.Runtime, err error, onerror interface{}) {
	if onerror != nil {
		if fncerror, fncerrorok := onerror.(func(error)); fncerrorok {
			fncerror(err)
		} else if script != nil {
			script.InvokeFunction(onerror, err)
		}
	}
}

func invokeSuccess(script active.Runtime, onsuccess interface{}, args ...interface{}) {
	if onsuccess != nil {
		if fncsuccess, fncsuccessok := onsuccess.(func()); fncsuccessok {
			fncsuccess()
		} else if script != nil {
			script.InvokeFunction(onsuccess, args...)
		}
	}
}

func invokeFinalize(script active.Runtime, onfinalize interface{}) {
	if onfinalize != nil {
		if fncfinalize, fncfinalizeok := onfinalize.(func()); fncfinalizeok {
			fncfinalize()
		} else if script != nil {

			script.InvokeFunction(onfinalize)
		}
	}
}

//NewConnection - dbms,driver name and datasource name (cn-string)
func NewConnection(dbms *DBMS, driverName, dataSourceName string) (cn *Connection) {
	cn = &Connection{dbms: dbms, driverName: driverName, dataSourceName: dataSourceName, lastmaxopencons: -1, lastmaxidecons: -1, maxopencons: -1, maxidlecons: -1}
	return
}

func calibrateConnection(cn *Connection, a ...interface{}) {
	if cn != nil && len(a) > 0 {
		var idlecons = cn.maxidlecons
		var opencons = cn.maxopencons
		for di := range a {
			if d := a[di]; d != nil {
				if mpcnsttngs, _ := d.(map[string]interface{}); mpcnsttngs != nil && len(mpcnsttngs) > 0 {
					for k := range mpcnsttngs {
						v := mpcnsttngs[k]
						if k == "max-idle-cons" {
							if vidlecons, _ := v.(int); vidlecons != idlecons {
								idlecons = vidlecons
							}
						} else if k == "max-open-cons" {
							if vopencons, _ := v.(int); vopencons != opencons {
								opencons = vopencons
							}
						}
					}
				}
			}
		}
		if idlecons != cn.maxidlecons {
			cn.SetMaxIdleConns(idlecons)
		}
		if opencons != cn.maxopencons {
			cn.SetMaxOpenConns(opencons)
		}
	}
}

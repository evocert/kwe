package database

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/parameters"
)

//Connection - struct
type Connection struct {
	dbms                       *DBMS
	driverName, dataSourceName string
	db                         *sql.DB
	dbinvoker                  func(string, ...interface{}) (*sql.DB, error)
}

func runeReaderToString(rnr io.RuneReader) (s string) {
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
	} else if exctr.cn.driverName == "postgres" {
		if argi == -1 {
			s = ("$" + fmt.Sprintf("%d", len(exctr.qryArgs)+1))
		}
		if argvs, argvsok := prmval.(string); argvsok {
			if fltval, nrerr := strconv.ParseFloat(argvs, 64); nrerr == nil {
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
	if qrys, qrysok := query.(string); qrysok && qrys != "" {
		rnrr = bufio.NewReader(strings.NewReader(qrys))
	} else if qryrnr, qryrnrok := query.(io.RuneReader); qryrnrok {
		rnrr = qryrnr
	} else if qryr, qryrok := query.(io.Reader); qryrok {
		rnrr = bufio.NewReader(qryr)
	}

	mappedVals = map[string]interface{}{}
	var foundTxt = false
	if len(args) == 1 {

		if pargs, ispargs := args[0].(*parameters.Parameters); ispargs {
			for _, skey := range pargs.StandardKeys() {
				mappedVals[skey] = strings.Join(pargs.Parameter(skey), "")
			}
		} else if pmargs, ispmargs := args[0].(map[string]interface{}); ispmargs {
			for pmk, pmv := range pmargs {
				if mpv, mpvok := pmv.(map[string]interface{}); mpvok && mpv != nil && len(mpv) > 0 {

				} else {
					mappedVals[pmk] = pmv
				}
			}
		}
	}

	if len(exctr.qryArgs) == 0 {
		exctr.qryArgs = []interface{}{}
	}

	stmnt = ""

	var rns = make([]rune, 1024)
	var rnsi = 0
	var prvr = rune(0)
	var prmslbl = [][]rune{[]rune("@@"), []rune("@@")}
	var prmslbli = []int{0, 0}

	var appr = func(r rune) {
		rns[rnsi] = r
		rnsi++
		if rnsi == len(rns) {
			stmnt += string(rns)
			rnsi = 0
		}
	}

	var apprs = func(p []rune) {
		if pl := len(p); pl > 0 {
			pi := 0
			for pi < pl {
				if l := (len(rns) - rnsi); (pl - pi) >= l {
					copy(rns[rnsi:rnsi+l], p[pi:pi+l])
					rnsi += l
					pi += l
				} else if l := (pl - pi); l < (len(rns) - rnsi) {
					copy(rns[rnsi:rnsi+l], p[pi:pi+l])
					rnsi += l
					pi += l
				}
				if rnsi == len(rns) {
					stmnt += string(rns)
					rnsi = 0
				}
			}
		}
	}

	var psblprmnme = make([]rune, 8192)
	var psblprmnmei = 0
	for rnrr != nil {
		r, s, e := rnrr.ReadRune()
		if s > 0 {
			if len(mappedVals) == 0 {
				appr(r)
			} else {
				if foundTxt {
					appr(r)
					if r == '\'' {
						if prvr == r {
							foundTxt = false
							prvr = rune(0)
						} else {
							prvr = r
						}
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
										for mpvk, mpv := range mappedVals {
											if fndprm = strings.ToUpper(psbprmnme) == strings.ToUpper(mpvk); fndprm {
												if validNames == nil {
													validNames = []string{}
												}
												validNames = append(validNames, mpvk)
												apprs([]rune(parseParam(exctr, mpv, -1)))
												break
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
			}
		}
		if e != nil {
			break
		}
	}
	if rnsi > 0 {
		stmnt += string(rns[:rnsi])
	}
	return
}

//GblExecute - public for query()*Executor
func (cn *Connection) GblExecute(query interface{}, prms ...interface{}) (exctr *Executor, err error) {
	if _, exctr, err = cn.query(query, true, nil, nil, nil, prms...); err != nil {

	}
	return
}

//GblQuery - public for query() *Reader
func (cn *Connection) GblQuery(query interface{}, prms ...interface{}) (reader *Reader, err error) {
	reader, _, err = cn.query(query, false, nil, nil, nil, prms...)
	if err != nil && reader == nil {

	}
	return
}

func (cn *Connection) inMapOut(mpin map[string]interface{}, out io.Writer, ioargs ...interface{}) (hasoutput bool) {
	if mpin != nil {
		if mpl := len(mpin); mpl > 0 {
			if out != nil {
				hasoutput = true
				iorw.Fprint(out, "{")
			}
			for mk, mv := range mpin {
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
				if mpl > 0 {
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
		if mr, mrok := in.(io.Reader); mrok && mr != nil {
			var buff = iorw.NewBuffer()
			func() {
				defer buff.Close()
				if buffl, bufferr := io.Copy(buff, mr); bufferr == nil {
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
									hasoutput = cn.inMapOut(rqstmp, out, ioargs...)
								}
							} else {

							}
						}()
					}
				}
			}()
		}
		if !hasoutput {
			if out != nil {
				iorw.Fprint(out, "{}")
			}
		}
	} else {
		if out != nil {
			iorw.Fprint(out, "{}")
		}
	}
}

func (cn *Connection) query(query interface{}, noreader bool, onsuccess, onerror, onfinalize interface{}, args ...interface{}) (reader *Reader, exctr *Executor, err error) {
	var argsn = 0
	var script active.Runtime = nil
	var canRepeat = false
	for argsn < len(args) {
		var d = args[argsn]
		if _, dok := d.(*parameters.Parameters); dok {
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

	if cn.db == nil {
		if cn.dbinvoker == nil {
			if dbinvoker, hasdbinvoker := cn.dbms.driverDbInvoker(cn.driverName); hasdbinvoker {
				cn.dbinvoker = dbinvoker
			}
		}
		if cn.db, err = cn.dbinvoker(cn.dataSourceName); err == nil && cn.db != nil {
			//cn.db.SetMaxIdleConns(runtime.NumCPU() * 4)
		}
		if err != nil && onerror != nil {
			invokeError(script, err, onerror)
		}
	}
	if cn.db != nil {
		if err = cn.db.Ping(); err != nil {
			cn.db.Close()
			cn.db = nil
			if cn.dbinvoker == nil {
				if dbinvoker, hasdbinvoker := cn.dbms.driverDbInvoker(cn.driverName); hasdbinvoker {
					cn.dbinvoker = dbinvoker
				}
			}
			if cn.db, err = cn.dbinvoker(cn.dataSourceName); err == nil && cn.db != nil {
				cn.db.Close()
				cn.db, err = cn.dbinvoker(cn.dataSourceName)
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
	cn = &Connection{dbms: dbms, driverName: driverName, dataSourceName: dataSourceName}
	return
}

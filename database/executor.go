package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/parameters"
	"github.com/evocert/kwe/web"
)

//Executor - struct
type Executor struct {
	orgstmnt     string
	Stmnt        string
	jsndcdr      *json.Decoder
	lastdlm      string
	tknlvl       int
	db           *sql.DB
	cn           *Connection
	stmt         *sql.Stmt
	strmqrystngs map[string]interface{}
	lasterr      error
	lastInsertID int64
	rowsAffected int64
	mappedArgs   map[string]interface{}
	argNames     []string
	qryArgs      []interface{}
	OnSuccess    interface{}
	OnError      interface{}
	OnFinalize   interface{}
	OnClose      func(*Executor)
	script       active.Runtime
	canRepeat    bool
}

func newExecutor(cn *Connection, db *sql.DB, query interface{}, strmqrystngs map[string]interface{}, canRepeat bool, script active.Runtime, onsuccess, onerror, onfinalize interface{}, args ...interface{}) (exctr *Executor) {
	var argsn = 0
	for argsn < len(args) {
		var d = args[argsn]
		if _, dok := d.(*parameters.Parameters); dok {
			argsn++
		} else if _, dok := d.(*Reader); dok {
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
					args = args[:argsn]
				} else if argsn < len(args) {
					args = append(args[:argsn], args[argsn+1:]...)
				}
			}
		}
	}
	exctr = &Executor{db: db, cn: cn, strmqrystngs: strmqrystngs, script: script, canRepeat: canRepeat, OnSuccess: onsuccess, OnError: onerror, OnFinalize: onfinalize}
	exctr.Stmnt, exctr.argNames, exctr.mappedArgs = queryToStatement(exctr, query, args...)
	return
}

func getTypeByName(tpmn string) (t reflect.Type) {
	if tpmn == "bool" {
		t = reflect.TypeOf(false)
	} else if tpmn == "int" {
		t = reflect.TypeOf(int(0))
	} else if tpmn == "int8" {
		t = reflect.TypeOf(int8(0))
	} else if tpmn == "int16" {
		t = reflect.TypeOf(int16(0))
	} else if tpmn == "int32" {
		t = reflect.TypeOf(int32(0))
	} else if tpmn == "int64" {
		t = reflect.TypeOf(int64(0))
	} else if tpmn == "uint" {
		t = reflect.TypeOf(uint(0))
	} else if tpmn == "uint8" {
		t = reflect.TypeOf(uint8(0))
	} else if tpmn == "uint16" {
		t = reflect.TypeOf(uint16(0))
	} else if tpmn == "uint32" {
		t = reflect.TypeOf(uint32(0))
	} else if tpmn == "uint64" {
		t = reflect.TypeOf(uint64(0))
	} else if tpmn == "float32" {
		t = reflect.TypeOf(float32(0))
	} else if tpmn == "float64" {
		t = reflect.TypeOf(float64(0))
	} else if tpmn == "complex64" {
		t = reflect.TypeOf(complex64(0))
	} else if tpmn == "complex128" {
		t = reflect.TypeOf(complex128(0))
	} else if tpmn == "Time" {
		t = reflect.TypeOf(time.Now())
	} else {
		t = reflect.TypeOf("")
	}
	return
}

func (exctr *Executor) isRemote() bool {
	return exctr.cn != nil && exctr.cn.IsRemote()
}

//ExecError - struct
type ExecError struct {
	err   error
	stmnt string
}

//Statement return statement executed that caused error
func (execerr *ExecError) Statement() string {
	if execerr != nil {
		return execerr.stmnt
	}
	return ""
}

func (execerr *ExecError) Error() string {
	if execerr != nil && execerr.err != nil {
		return execerr.err.Error()
	}
	return ""
}

func newExecErr(err error, stmnt string) (execerr *ExecError) {
	execerr = &ExecError{err: err, stmnt: stmnt}
	return
}

func (exctr *Executor) execute(forrows ...bool) (rws RWSAPI, cltpes []*ColumnType, cls []string) {
	if exctr.isRemote() {
		pi, po := io.Pipe()
		ctx, ctxcancel := context.WithCancel(context.Background())
		go func() {
			defer func() {
				po.Close()
			}()
			ctxcancel()
			exctr.webquery(len(forrows) == 1 && forrows[0], po)
		}()
		<-ctx.Done()
		ctx = nil
		exctr.jsndcdr = json.NewDecoder(pi)
		exctr.tknlvl = 0
		for {
			tkn, tknerr := exctr.jsndcdr.Token()
			if tknerr != nil {
				if tknerr != io.EOF {
					exctr.lasterr = tknerr
				}
				break
			} else {
				if dlm, dlmok := tkn.(json.Delim); dlmok {
					if exctr.lastdlm = dlm.String(); exctr.lastdlm == "{" {
						exctr.tknlvl++
					} else if exctr.lastdlm == "}" {
						exctr.tknlvl--
					}
				} else {
					s, _ := tkn.(string)
					if exctr.tknlvl == 1 && s != "" {

					} else {
						if exctr.tknlvl == 2 && s != "" {
							if s == "columns" {
								clsarr := []interface{}{}
								if tknerr = exctr.jsndcdr.Decode(&clsarr); tknerr != nil {
									exctr.lasterr = tknerr
									break
								}
								if l := len(clsarr); l > 0 {
									cls = make([]string, l)
									cltpes = make([]*ColumnType, l)
									for cn := range clsarr {
										cltp := &ColumnType{}
										if c := clsarr[cn]; c != nil {
											cmp, _ := c.(map[string]interface{})
											for ck := range cmp {
												cv := cmp[ck]
												if ck == "name" {
													cls[cn], _ = cv.(string)
													cltp.name = cls[cn]
												} else if ck == "length" {
													if flt, fltok := cv.(float64); fltok {
														cltp.length = int64(flt)
													}
												} else if ck == "dbtype" {
													cltp.databaseType, _ = cv.(string)
												} else if ck == "numeric" {
													cltp.hasPrecisionScale, _ = cv.(bool)
												} else if ck == "scale" {
													if flt, fltok := cv.(float64); fltok {
														cltp.scale = int64(flt)
													}
												} else if ck == "precision" {
													if flt, fltok := cv.(float64); fltok {
														cltp.precision = int64(flt)
													}
												} else if ck == "type" {
													if tpnm, _ := cv.(string); tpnm != "" {
														cltp.scanType = getTypeByName(tpnm)
													}
												}
											}
										}
										cltpes[cn] = cltp
									}
								}
								break
							} else if s == "error" {
								if tkn, tknerr = exctr.jsndcdr.Token(); tknerr != nil {
									exctr.lasterr = tknerr
									break
								}
								if serr, serrok := tkn.(string); serrok {
									exctr.lasterr = fmt.Errorf("%v", serr)
								} else {
									exctr.lasterr = fmt.Errorf("%v", "unknown error")
								}
								break
							} else if s == "data" {
								break
							}
						} else {
							break
						}
					}
				}
			}
		}
	} else {
		if exctr.stmt == nil {
			if len(exctr.strmqrystngs) == 0 {
				if exctr.stmt, exctr.lasterr = exctr.db.Prepare(exctr.Stmnt); exctr.lasterr != nil {
					exctr.lasterr = newExecErr(exctr.lasterr, exctr.Stmnt)
				}
			}
		}
		if exctr.lasterr == nil {

			exctr.lastInsertID = -1
			exctr.rowsAffected = -1
			if exctr.stmt != nil {
				if exctr.canRepeat && len(exctr.argNames) > 0 {
					for argn := range exctr.argNames {
						if argnme := exctr.argNames[argn]; argnme != "" {
							if prmv, prmvok := exctr.mappedArgs[argnme]; prmvok {
								parseParam(exctr, prmv, argn)
							} else {
								parseParam(exctr, nil, argn)
							}
						}
					}
				}
			}
			if len(forrows) >= 1 && forrows[0] {
				if rws, exctr.lasterr = func() (rswapi RWSAPI, rwserr error) {
					if exctr.stmt != nil && exctr.strmqrystngs == nil {
						var rws *sql.Rows = nil
						if rws, rwserr = exctr.stmt.Query(exctr.qryArgs...); rws != nil && rwserr == nil {
							rswapi, rwserr = newRWSReader(rws, nil)
						}
					} else {
						rswapi, rwserr = newRWSReader(nil, exctr.strmqrystngs)
					}
					return
				}(); rws != nil && exctr.lasterr == nil {
					cltps, _ := rws.ColumnTypes()
					cls, _ = rws.Columns()
					if len(cls) > 0 {
						clsdstnc := map[string]int{}
						clsdstncorg := map[string]int{}
						cltpes = cltps[:] //columnTypes(cltps, cls)
						for cn := range cls {
							c := cls[cn]
							if c != "" {
								if ci, ciok := clsdstnc[c]; ciok {
									if orgcn, orgok := clsdstncorg[c]; orgok && cls[orgcn] == c {
										cls[orgcn] = fmt.Sprintf("%s%d", c, 0)
									}
									clsdstnc[c]++
									c = fmt.Sprintf("%s%d", c, ci+1)
								} else {
									if _, orgok := clsdstncorg[c]; !orgok {
										clsdstncorg[c] = cn
									}
									clsdstnc[c] = 0
								}
							}
							cls[cn] = c

						}
					}
				} else if exctr.lasterr != nil {
					exctr.lasterr = newExecErr(exctr.lasterr, exctr.Stmnt)
					invokeError(exctr.script, exctr.lasterr, exctr.OnError)
				}
			} else {
				if rslt, rslterr := exctr.stmt.Exec(exctr.qryArgs...); rslterr == nil {
					if exctr.cn.driverName == "sqlserver" {
						exctr.lastInsertID = -1
					} else if exctr.lastInsertID, rslterr = rslt.LastInsertId(); rslterr != nil {
						exctr.lastInsertID = -1
					}
					if exctr.rowsAffected, rslterr = rslt.RowsAffected(); rslterr != nil {
						exctr.rowsAffected = -1
					}
					invokeSuccess(exctr.script, exctr.OnSuccess, exctr)
				} else {
					exctr.lasterr = rslterr
					exctr.lasterr = newExecErr(exctr.lasterr, exctr.Stmnt)
					invokeError(exctr.script, exctr.lasterr, exctr.OnError)
				}
			}
		}
	}
	return
}

func (exctr *Executor) webquery(forrows bool, out io.Writer, iorags ...interface{}) (err error) {
	pi, pw := io.Pipe()
	ctx, ctxcancel := context.WithCancel(context.Background())
	func() {
		defer func() {
			pi.Close()
			pi = nil
		}()
		go func() {
			var pwerr error = nil
			defer func() {
				if pwerr != nil {
					pw.CloseWithError(pwerr)
				} else {
					pw.Close()
				}
			}()
			ctxcancel()
			encw := json.NewEncoder(pw)
			rqstmpstngs := map[string]interface{}{}
			if len(exctr.mappedArgs) > 0 {
				for kmp := range exctr.mappedArgs {
					vmp := exctr.mappedArgs[kmp]
					rqstmpstngs[kmp] = vmp
				}
			}
			if forrows {
				rqstmpstngs["query"] = exctr.Stmnt
			} else {
				rqstmpstngs["execute"] = exctr.Stmnt
			}

			rqstmp := map[string]interface{}{"1": rqstmpstngs}
			encw.Encode(&rqstmp)
			encw = nil
			rqstmp = nil
		}()
		<-ctx.Done()
		ctx = nil
		datasource := exctr.cn.dataSourceName
		if strings.HasPrefix(datasource, "http://") || strings.HasPrefix(datasource, "https://") {
			func() {
				/*var rspheaders = map[string]string{}*/
				var rqstheaders = map[string]string{}
				rqstheaders["Content-Type"] = "application/json"
				args := []interface{}{rqstheaders}
				if len(exctr.cn.args) > 0 {
					exctr.cn.args = append(exctr.cn.args, pi)
					args := append(args, exctr.cn.args...)
					if rspr, rsprerr := web.DefaultClient.Send(datasource, args...); rsprerr == nil {
						if rspr != nil {
							err = iorw.Fprint(out, rspr)
						}
					} else {
						err = rsprerr
					}
				} else {
					args = append(args, pi)
					if rspr, rsprerr := web.DefaultClient.Send(datasource, args...); rsprerr == nil {
						if rspr != nil {
							err = iorw.Fprint(out, rspr)
						}
					} else if rsprerr != nil {
						err = rsprerr
					}
				}
				rqstheaders = nil
				args = nil
			}()
		} else if strings.HasPrefix(datasource, "ws://") || strings.HasPrefix(datasource, "wss://") {
			func() {
				var rqstheaders = map[string]string{}
				//rqstheaders["Content-Type"] = "application/json"
				args := []interface{}{rqstheaders}
				if len(exctr.cn.args) > 0 {
					//exctr.cn.args = append(exctr.cn.args, pi)
					args := append(args, exctr.cn.args...)
					if rsprw, rsprerr := web.DefaultClient.SendReceive(datasource, args...); rsprerr == nil {
						if pi != nil {
							err = rsprw.Print(pi)
						}
						if rsprw != nil {
							err = iorw.Fprint(out, rsprw)
						}
					} else {
						err = rsprerr
					}
				} else {
					//args = append(args, pi)
					if rsprw, rsprerr := web.DefaultClient.SendReceive(datasource, args...); rsprerr == nil {
						if pi != nil {
							err = rsprw.Print(pi)
						}
						if rsprw != nil {
							err = iorw.Fprint(out, rsprw)
						}
					} else if rsprerr != nil {
						err = rsprerr
					}
				}
				rqstheaders = nil
				args = nil
			}()
		}
	}()
	return
}

//Repeat - repeat last query by repopulating parameters but dont regenerate last statement
func (exctr *Executor) Repeat(args ...interface{}) (err error) {
	if len(args) == 1 {
		if pargs, ispargs := args[0].(*parameters.Parameters); ispargs {
			for _, skey := range pargs.StandardKeys() {
				for argn := range exctr.argNames {
					argnme := exctr.argNames[argn]
					if strings.EqualFold(skey, argnme) {
						exctr.mappedArgs[argnme] = strings.Join(pargs.Parameter(skey), "")
						break
					}
				}
			}
		} else if pmargs, ispmargs := args[0].(map[string]interface{}); ispmargs {
			for pmk := range pmargs {
				pmv := pmargs[pmk]
				if mpv, mpvok := pmv.(map[string]interface{}); mpvok && mpv != nil {

				} else {
					for argn := range exctr.argNames {
						argnme := exctr.argNames[argn]
						//strings.ToLower(pmk) == strings.ToLower(argnme)
						if strings.EqualFold(pmk, argnme) {
							exctr.mappedArgs[argnme] = pmv
							break
						}
					}
				}
			}
		}
	}
	if !exctr.canRepeat {
		exctr.canRepeat = true
	}
	exctr.execute()
	if err = exctr.lasterr; err != nil {
		invokeError(exctr.script, err, exctr.OnError)
	}
	return
}

//ToJSON write *Executor out to json
func (exctr *Executor) ToJSON(w io.Writer) (err error) {
	if w != nil {
		if jsnrdr := exctr.JSONReader(); jsnrdr != nil {
			func() {
				defer func() { jsnrdr = nil }()
				if _, err = io.Copy(w, jsnrdr); err != nil && err == io.EOF {
					err = nil
				}
			}()
		}
	}
	return
}

//JSONReader return *JSONReader
func (exctr *Executor) JSONReader() (jsnrdr *JSONReader) {
	jsnrdr = NewJSONReader(nil, exctr, nil)
	return
}

//JSON execute *Executor and return json as string
func (exctr *Executor) JSON() (s string, err error) {
	bufr := iorw.NewBuffer()
	func() {
		defer bufr.Close()
		if err = exctr.ToJSON(bufr); err == nil {
			s = bufr.String()
		}
	}()
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
		if exctr.db != nil {
			exctr.db = nil
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
		if exctr.jsndcdr != nil {
			exctr.jsndcdr = nil
		}
		if exctr.mappedArgs != nil {
			exctr.mappedArgs = nil
			for mk := range exctr.mappedArgs {
				exctr.mappedArgs[mk] = nil
				delete(exctr.mappedArgs, mk)
			}
		}
		if exctr.qryArgs != nil {
			for len(exctr.qryArgs) > 0 {
				exctr.qryArgs[0] = nil
				exctr.qryArgs = exctr.qryArgs[1:]
			}
			exctr.qryArgs = nil
		}
		if exctr.OnClose != nil {
			exctr.OnClose(exctr)
			exctr.OnClose = nil
		}
		if exctr.strmqrystngs != nil {
			if len(exctr.strmqrystngs) > 0 {
				for strnk := range exctr.strmqrystngs {
					delete(exctr.strmqrystngs, strnk)
				}
			}
			exctr.strmqrystngs = nil
		}
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

package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/parameters"
)

//Executor - struct
type Executor struct {
	orgstmnt     string
	stmnt        string
	endpnt       *EndPoint
	jsndcdr      *json.Decoder
	lastdlm      string
	tknlvl       int
	db           *sql.DB
	cn           *Connection
	stmt         *sql.Stmt
	lasterr      error
	lastInsertID int64
	rowsAffected int64
	mappedArgs   map[string]interface{}
	argNames     []string
	qryArgs      []interface{}
	OnSuccess    interface{}
	OnError      interface{}
	OnFinalize   interface{}
	script       active.Runtime
	canRepeat    bool
}

func newExecutor(cn *Connection, db *sql.DB, endpnt *EndPoint, query interface{}, canRepeat bool, script active.Runtime, onsuccess, onerror, onfinalize interface{}, args ...interface{}) (exctr *Executor) {
	var argsn = 0
	for argsn < len(args) {
		var d = args[argsn]
		if _, dok := d.(*parameters.Parameters); dok {
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
					args = append(args[:argsn])
				} else if argsn < len(args) {
					args = append(args[:argsn], args[argsn+1:]...)
				}
			}
		}
	}
	exctr = &Executor{endpnt: endpnt, db: db, cn: cn, script: script, canRepeat: canRepeat, OnSuccess: onsuccess, OnError: onerror, OnFinalize: onfinalize}
	exctr.stmnt, exctr.argNames, exctr.mappedArgs = queryToStatement(exctr, query, args...)
	return
}

func getTypeByName(tpmn string) (t reflect.Type) {
	fmt.Println(tpmn)
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

func (exctr *Executor) execute(forrows ...bool) (rws *sql.Rows, cltpes []*ColumnType, cls []string) {
	if exctr.endpnt == nil {
		if exctr.stmt == nil {
			exctr.stmt, exctr.lasterr = exctr.db.Prepare(exctr.stmnt)
		}
		if exctr.lasterr == nil && exctr.stmt != nil {
			exctr.lastInsertID = -1
			exctr.rowsAffected = -1
			if exctr.canRepeat && len(exctr.argNames) > 0 {
				for agrn, argnme := range exctr.argNames {
					if prmv, prmvok := exctr.mappedArgs[argnme]; prmvok {
						parseParam(exctr, prmv, agrn)
					} else {
						parseParam(exctr, nil, agrn)
					}
				}
			}

			if len(forrows) >= 1 && forrows[0] {
				if rws, exctr.lasterr = exctr.stmt.Query(exctr.qryArgs...); rws != nil && exctr.lasterr == nil {
					cltps, _ := rws.ColumnTypes()
					cls, _ = rws.Columns()
					if len(cls) > 0 {
						clsdstnc := map[string]int{}
						clsdstncorg := map[string]int{}
						cltpes = columnTypes(cltps, cls)
						for cn, c := range cls {
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
							cls[cn] = c
						}
					}
				} else if exctr.lasterr != nil {
					invokeError(exctr.script, exctr.lasterr, exctr.OnError)
				}
			} else {
				if rslt, rslterr := exctr.stmt.Exec(exctr.qryArgs...); rslterr == nil {
					if exctr.lastInsertID, rslterr = rslt.LastInsertId(); rslterr != nil {
						exctr.lastInsertID = -1
					}
					if exctr.rowsAffected, rslterr = rslt.RowsAffected(); rslterr != nil {
						exctr.rowsAffected = -1
					}
					invokeSuccess(exctr.script, exctr.OnSuccess, exctr)
				} else {
					exctr.lasterr = rslterr
					invokeError(exctr.script, exctr.lasterr, exctr.OnError)
				}
			}
		}
	} else {
		pi, po := io.Pipe()
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer func() {
				po.Close()
			}()
			wg.Done()
			exctr.endpnt.query(exctr, len(forrows) == 1 && forrows[0], po)
		}()
		wg.Wait()
		exctr.jsndcdr = json.NewDecoder(pi)
		exctr.tknlvl = 0
		for {
			tkn, tknerr := exctr.jsndcdr.Token()
			if tknerr != nil {
				exctr.lasterr = tknerr
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
									for cn, c := range clsarr {
										cltp := &ColumnType{}
										if c != nil {
											cmp, _ := c.(map[string]interface{})
											for ck, cv := range cmp {
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
	}
	return
}

//Repeat - repeat last query by repopulating parameters but dont regenerate last statement
func (exctr *Executor) Repeat(args ...interface{}) (err error) {
	if len(args) == 1 {
		if pargs, ispargs := args[0].(*parameters.Parameters); ispargs {
			for _, skey := range pargs.StandardKeys() {
				for _, argnme := range exctr.argNames {
					if strings.ToLower(skey) == strings.ToLower(argnme) {
						exctr.mappedArgs[argnme] = strings.Join(pargs.Parameter(skey), "")
						break
					}
				}
			}
		} else if pmargs, ispmargs := args[0].(map[string]interface{}); ispmargs {
			for pmk, pmv := range pmargs {
				if mpv, mpvok := pmv.(map[string]interface{}); mpvok && mpv != nil {

				} else {
					for _, argnme := range exctr.argNames {
						if strings.ToLower(pmk) == strings.ToLower(argnme) {
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

//Close - Executor
func (exctr *Executor) Close() (err error) {
	if exctr != nil {
		if exctr.db != nil {
			exctr.db = nil
		}
		if exctr.cn != nil {
			exctr.cn = nil
		}
		if exctr.endpnt != nil {
			exctr.endpnt = nil
		}
		if exctr.stmt != nil {
			err = exctr.stmt.Close()
			exctr.stmt = nil
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

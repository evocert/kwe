package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/dop251/goja"
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

func queryToStatement(exctr *Executor, query interface{}, args ...interface{}) (stmnt string, qryArgs []interface{}, mappedVals map[string]interface{}) {
	var rnrr io.RuneReader = nil
	if qrys, qrysok := query.(string); qrysok && qrys != "" {
		rnrr = bufio.NewReader(strings.NewReader(qrys))
	} else if qryrnr, qryrnrok := query.(io.RuneReader); qryrnrok {
		rnrr = qryrnr
	} else if qryr, qryrok := query.(io.Reader); qryrok {
		rnrr = bufio.NewReader(qryr)
	}

	//var validNames = []string{}
	mappedVals = map[string]interface{}{}
	var foundTxt = false
	if len(args) == 1 {

		if pargs, ispargs := args[0].(*parameters.Parameters); ispargs {
			//ignoreCase = true
			for _, skey := range pargs.StandardKeys() {
				//validNames = append(validNames, skey)
				mappedVals[skey] = strings.Join(pargs.Parameter(skey), "")
			}
		} else if pmargs, ispmargs := args[0].(map[string]interface{}); ispmargs {
			for pmk, pmv := range pmargs {
				if mpv, mpvok := pmv.(map[string]interface{}); mpvok && mpv != nil {

				} else {
					//validNames = append(validNames, pmk)
					mappedVals[pmk] = pmv
				}
			}
		}
	}

	qryArgs = []interface{}{}

	var parseParam = func(prmval interface{}) (s string) {
		if exctr.cn.driverName == "sqlserver" {
			s = ("@p" + fmt.Sprintf("%d", len(qryArgs)))
			qryArgs = append(qryArgs, sql.Named("p"+fmt.Sprintf("%d", len(qryArgs)), prmval))
		} else if exctr.cn.driverName == "postgres" {
			//qry += ("$S" + fmt.Sprintf("%d", len(txtArgs)))
			/*argv := prmval
			if argvs, argvsok := argv.(string); argvsok {
				qry += "CONVERT_FROM(DECODE('" + base64.URLEncoding.EncodeToString([]byte(argvs)) + "', 'BASE64'), 'UTF-8')"
			} else {
				qry += fmt.Sprint(argv)
			}*/

			prmname := "$" + fmt.Sprintf("%d", len(qryArgs))
			qryArgs = append(qryArgs, prmval)
			s = (prmname)
		} else {
			qryArgs = append(qryArgs, prmval)
			s = "?"
		}
		return
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
									if mpv, mpvok := mappedVals[string(psblprmnme[:psblprmnmei])]; mpvok {
										apprs([]rune(parseParam(mpv)))
									} else {
										apprs(prmslbl[0])
										apprs(psblprmnme[:psblprmnmei])
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

func (cn *Connection) query(query interface{}, noreader bool, onsuccess, onerror, onfinalize interface{}, args ...interface{}) (reader *Reader, exctr *Executor, err error) {
	var argsn = 0
	var script active.Runtime = nil
	for argsn < len(args) {
		var d = args[argsn]
		if _, dok := d.(*parameters.Parameters); dok {
			argsn++
		} else if _, dok := d.(map[string]interface{}); dok {
			argsn++
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
			cn.db.SetMaxIdleConns(runtime.NumCPU() * 4)
		}
		if err != nil && onerror != nil {
			invokeError(err, onerror)
		}
	}
	if cn.db != nil {
		if query != nil {
			exctr = newExecutor(cn, cn.db, query, onsuccess, onerror, onfinalize, args...)
			if noreader {
				exctr.execute()
			} else {
				reader = newReader(exctr)
				reader.execute()
			}
		}
	}
	return
}

func invokeError(err error, onerror interface{}) {
	if onerror != nil {
		if fncerror, fncerrorok := onerror.(func(error)); fncerrorok {
			fncerror(err)
		} else if atverror, atverrorok := onerror.(func(goja.FunctionCall) goja.Value); atverrorok {
			var fnccall = goja.FunctionCall{}
			if atvval, atvvalok := onerror.(*goja.Callable); atvvalok {
				if atvval != nil {

				}
			}
			atverror(fnccall)
		}
	}
}

func invokeSuccess(onsuccess interface{}) {
	if onsuccess != nil {
		if fncsuccess, fncsuccessok := onsuccess.(func()); fncsuccessok {
			fncsuccess()
		} else if atvsuccess, atvsuccessok := onsuccess.(func(goja.FunctionCall) goja.Value); atvsuccessok {
			var fnccall = goja.FunctionCall{}
			atvsuccess(fnccall)
		}
	}
}

func invokeFinalize(onfinalize interface{}) {
	if onfinalize != nil {
		if fncfinalize, fncfinalizeok := onfinalize.(func()); fncfinalizeok {
			fncfinalize()
		} else if atvfinalize, atvfinalizeok := onfinalize.(func(goja.FunctionCall) goja.Value); atvfinalizeok {

			atvfinalize(goja.FunctionCall{})
		}
	}
}

//NewConnection - dbms,driver name and datasource name (cn-string)
func NewConnection(dbms *DBMS, driverName, dataSourceName string) (cn *Connection) {
	cn = &Connection{dbms: dbms, driverName: driverName, dataSourceName: dataSourceName}
	return
}

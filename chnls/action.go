package chnls

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"github.com/evocert/kwe/database"
	"github.com/evocert/kwe/mimes"
	"github.com/evocert/kwe/scheduling"
)

//Action - struct
type Action struct {
	rqst    *Request
	rspath  string
	sttngs  map[string]interface{}
	prvactn *Action
}

func newAction(rqst *Request, rspath string) (actn *Action) {
	actn = &Action{rqst: rqst, rspath: rspath, prvactn: rqst.lstexctngactng}
	rqst.lstexctngactng = actn
	return
}

//Path of resource action thats currently processing
func (actn *Action) Path() string {
	if actn != nil {
		return actn.rspath
	}
	return ""
}

func executeAction(actn *Action, w ...io.Writer) (err error) {
	altw := func() io.Writer {
		if len(w) == 1 && w[0] != nil {
			return w[0]
		}
		return nil
	}

	defer func() {
		if actn != nil {
			actn.Close()
			actn = nil
		}
	}()
	var rspath = actn.rspath // actn.rsngpth.Path
	var rspathext = filepath.Ext(rspath)
	var isTextRequest = false
	var dbmsaliases map[string]*database.Connection = nil
	var schdlsaliases map[string]*scheduling.Schedule = nil
	if strings.HasPrefix(rspath, "/scheduling/") || strings.HasPrefix(rspath, "/scheduling-") {
		func() {
			defer func() {
				if schdlsaliases != nil {
					if aliasesl := len(schdlsaliases); aliasesl > 0 {
						aliasks := make([]string, aliasesl)
						aliasksi := 0
						for aliask := range schdlsaliases {
							aliasks[aliasksi] = aliask
							aliasksi++
						}
						for _, aliask := range aliasks {
							schdlsaliases[aliask] = nil
							delete(schdlsaliases, aliask)
						}
					}
					schdlsaliases = nil
				}
			}()
			var dbmspath = rspath
			var alias = "all"
			if strings.HasPrefix(dbmspath, "/scheduling/") {
				dbmspath = dbmspath[len("/scheduling/")-1:]
			} else if strings.HasPrefix(dbmspath, "/scheduling-") {
				dbmspath = dbmspath[len("/scheduling-"):]
				if strings.Index(dbmspath, "/") > 0 {
					alias = dbmspath[:strings.Index(dbmspath, "/")]
					dbmspath = dbmspath[strings.Index(dbmspath, "/"):]
				}
			}
			if alias != "" {
				if alias == "all" {
					if actn.rqst.Parameters().ContainsParameter("scheduling-alias") {
						if aliassesfound := actn.rqst.Parameters().Parameter("scheduling-alias"); len(aliassesfound) > 0 {
							for _, kalias := range aliassesfound {
								if exists, schdl := scheduling.GLOBALSCHEDULES().ScheduleExists(kalias); exists {
									if schdlsaliases == nil {
										schdlsaliases = map[string]*scheduling.Schedule{}
									}
									schdlsaliases[kalias] = schdl
								}
							}
						}
					} else {
						if rspathext == ".json" {
							actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(".json", "text/plain")
							var jsnr io.Reader = nil
							if actn.rqst.Parameters().ContainsParameter("scheduling:json") {
								jsnr = strings.NewReader(strings.Join(actn.rqst.Parameters().Parameter("scheduling:json"), ""))
							} else {
								jsnr = actn.rqst.RequestBody()
							}
							if jsnr != nil {
								scheduling.GLOBALSCHEDULES().InOut(jsnr, actn.rqst, actn.rqst.Parameters())
							}
						}
					}
				} else {
					if exists, schdl := scheduling.GLOBALSCHEDULES().ScheduleExists(alias); exists {
						if schdlsaliases == nil {
							schdlsaliases = map[string]*scheduling.Schedule{}
						}
						schdlsaliases[alias] = schdl
					}
				}
				if len(schdlsaliases) > 0 {
					for kalias, schdl := range schdlsaliases {
						if actn.rqst.Parameters().ContainsParameter(kalias + ":json") {
							jsnval := strings.Join(actn.rqst.Parameters().Parameter(kalias+":json"), "")
							actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(".json", "text/plain")
							schdl.InOut(strings.NewReader(jsnval), actn.rqst, actn.rqst.Parameters())
						} /*else if actn.rqst.Parameters().ContainsParameter(kalias + ":query") {
							dbrdr, dbrdrerr := dbcn.GblQuery(strings.Join(actn.rqst.Parameters().Parameter(kalias+":query"), ""), actn.rqst.Parameters())
							if rspathext != "" {
								if rspathext == ".json" {
									actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
									actn.rqst.copy(io.MultiReader(database.NewJSONReader(dbrdr, nil, dbrdrerr)), nil, false)
								} else if rspathext == ".js" {
									var script = true
									if actn.rqst.Parameters().ContainsParameter(kalias+":script") && strings.Join(actn.rqst.Parameters().Parameter(kalias+":script"), "") == "false" {
										script = false
									}
									if actn.rqst.Parameters().ContainsParameter(kalias + ":jscall") {
										if jscall := strings.Join(actn.rqst.Parameters().Parameter(kalias+":jscall"), ""); jscall != "" {
											if actn.rqst.mimetype == "" {
												actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
											}
											if script {
												actn.rqst.copy(io.MultiReader(strings.NewReader("script||"), strings.NewReader(jscall+"("), database.NewJSONReader(dbrdr, nil, dbrdrerr), strings.NewReader(");"), strings.NewReader("||script")), nil, false)
											} else {
												actn.rqst.copy(io.MultiReader(strings.NewReader(jscall+"("), database.NewJSONReader(dbrdr, nil, dbrdrerr), strings.NewReader(");")), nil, false)
											}
										}
									} else if actn.rqst.Parameters().ContainsParameter(kalias + ":jsvar") {
										if jsvar := strings.Join(actn.rqst.Parameters().Parameter(kalias+":jsvar"), ""); jsvar != "" {
											if actn.rqst.mimetype == "" {
												actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
											}
											if script {
												actn.rqst.copy(io.MultiReader(strings.NewReader("script||"), strings.NewReader(jsvar+"="), database.NewJSONReader(dbrdr, nil, dbrdrerr), strings.NewReader(";"), strings.NewReader("||script")), nil, false)
											} else {
												actn.rqst.copy(io.MultiReader(strings.NewReader(jsvar+"="), database.NewJSONReader(dbrdr, nil, dbrdrerr), strings.NewReader(";")), nil, false)
											}
										}
									}
								} else if rspathext == ".csv" {
									var csvsttngs = map[string]interface{}{}

									if actn.rqst.Parameters().ContainsParameter(kalias + ":csv") {
										var csvsttngsval = strings.Join(actn.rqst.Parameters().Parameter(kalias+":csv"), "")
										var decsttngs = json.NewDecoder(strings.NewReader(csvsttngsval))
										if decerr := decsttngs.Decode(&csvsttngs); decerr == nil {

										}
									}
									actn.rqst.copy(io.MultiReader(database.NewCSVReader(dbrdr, dbrdrerr, csvsttngs)), nil, false)
								}
							} else {

							}
							if dbrdr != nil {
								dbrdr.Close()
							}
						} else if actn.rqst.Parameters().ContainsParameter(kalias + ":execute") {
							exctr, exctrerr := dbcn.GblExecute(strings.Join(actn.rqst.Parameters().Parameter(kalias+":execute"), ""), actn.rqst.Parameters())
							if rspathext == "" {
								rspathext = ".json"
							}
							if rspathext != "" {
								if rspathext == ".json" {
									actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
									actn.rqst.copy(io.MultiReader(database.NewJSONReader(nil, exctr, exctrerr)), nil, false)
								} else if rspathext == ".js" {
									var script = true
									if actn.rqst.Parameters().ContainsParameter(kalias+":script") && strings.Join(actn.rqst.Parameters().Parameter(kalias+":script"), "") == "false" {
										script = false
									}
									if actn.rqst.Parameters().ContainsParameter(kalias + ":jscall") {
										if jscall := strings.Join(actn.rqst.Parameters().Parameter(kalias+":jscall"), ""); jscall != "" {
											if actn.rqst.mimetype == "" {
												actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
											}
											if script {
												actn.rqst.copy(io.MultiReader(strings.NewReader("script||"), strings.NewReader(jscall+"("), database.NewJSONReader(nil, exctr, exctrerr), strings.NewReader(");"), strings.NewReader("||script")), nil, false)
											} else {
												actn.rqst.copy(io.MultiReader(strings.NewReader(jscall+"("), database.NewJSONReader(nil, exctr, exctrerr), strings.NewReader(");")), nil, false)
											}
										}
									} else if actn.rqst.Parameters().ContainsParameter(kalias + ":jsvar") {
										if jsvar := strings.Join(actn.rqst.Parameters().Parameter(kalias+":jsvar"), ""); jsvar != "" {
											if actn.rqst.mimetype == "" {
												actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
											}
											if script {
												actn.rqst.copy(io.MultiReader(strings.NewReader("script||"), strings.NewReader(jsvar+"="), database.NewJSONReader(nil, exctr, exctrerr), strings.NewReader(";"), strings.NewReader("||script")), nil, false)
											} else {
												actn.rqst.copy(io.MultiReader(strings.NewReader(jsvar+"="), database.NewJSONReader(nil, exctr, exctrerr), strings.NewReader(";")), nil, false)
											}
										}
									}
								}
							} else {

							}
							if exctr != nil {
								exctr.Close()
							}
						} else {
							if rspathext == "" {
								rspathext = ".json"
							}
							if rspathext == ".json" {
								actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
								schdl.InOut(actn.rqst.RequestBody(), actn.rqst, actn.rqst.Parameters())
							}
						}*/
					}
				}
			}
		}()
	} else if strings.HasPrefix(rspath, "/dbms/") || strings.HasPrefix(rspath, "/dbms-") {
		func() {
			defer func() {
				if dbmsaliases != nil {
					if aliasesl := len(dbmsaliases); aliasesl > 0 {
						aliasks := make([]string, aliasesl)
						aliasksi := 0
						for aliask := range dbmsaliases {
							aliasks[aliasksi] = aliask
							aliasksi++
						}
						for _, aliask := range aliasks {
							dbmsaliases[aliask] = nil
							delete(dbmsaliases, aliask)
						}
					}
					dbmsaliases = nil
				}
			}()
			var dbmspath = rspath
			var alias = "all"
			if strings.HasPrefix(dbmspath, "/dbms/") {
				dbmspath = dbmspath[len("/dbms/")-1:]
			} else if strings.HasPrefix(dbmspath, "/dbms-") {
				dbmspath = dbmspath[len("/dbms-"):]
				if strings.Index(dbmspath, "/") > 0 {
					alias = dbmspath[:strings.Index(dbmspath, "/")]
					dbmspath = dbmspath[strings.Index(dbmspath, "/"):]
				}
			}
			if alias != "" {
				if alias == "all" {
					if actn.rqst.Parameters().ContainsParameter("dbms-alias") {
						if aliassesfound := actn.rqst.Parameters().Parameter("dbms-alias"); len(aliassesfound) > 0 {
							for _, kalias := range aliassesfound {
								if exists, dbcn := database.GLOBALDBMS().AliasExists(kalias); exists {
									if dbmsaliases == nil {
										dbmsaliases = map[string]*database.Connection{}
									}
									dbmsaliases[kalias] = dbcn
								}
							}
						}
					} else {
						if rspathext == ".json" {
							actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(".json", "text/plain")
							var jsnr io.Reader = nil
							if actn.rqst.Parameters().ContainsParameter("dbms:json") {
								jsnr = strings.NewReader(strings.Join(actn.rqst.Parameters().Parameter("dbms:json"), ""))
							} else {
								jsnr = actn.rqst.RequestBody()
							}
							if jsnr != nil {
								database.GLOBALDBMS().InOut(jsnr, actn.rqst, actn.rqst.Parameters())
							}
						}
					}
				} else {
					if exists, dbcn := database.GLOBALDBMS().AliasExists(alias); exists {
						if dbmsaliases == nil {
							dbmsaliases = map[string]*database.Connection{}
						}
						dbmsaliases[alias] = dbcn
					}
				}
				if len(dbmsaliases) > 0 {
					for kalias, dbcn := range dbmsaliases {
						if actn.rqst.Parameters().ContainsParameter(kalias + ":json") {
							jsnval := strings.Join(actn.rqst.Parameters().Parameter(kalias+":json"), "")
							actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(".json", "text/plain")
							dbcn.InOut(strings.NewReader(jsnval), actn.rqst, actn.rqst.Parameters())
						} else if actn.rqst.Parameters().ContainsParameter(kalias + ":query") {
							dbrdr, dbrdrerr := dbcn.GblQuery(strings.Join(actn.rqst.Parameters().Parameter(kalias+":query"), ""), actn.rqst.Parameters())
							if rspathext != "" {
								if rspathext == ".json" {
									actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
									actn.rqst.copy(io.MultiReader(database.NewJSONReader(dbrdr, nil, dbrdrerr)), nil, false)
								} else if rspathext == ".js" {
									var script = true
									if actn.rqst.Parameters().ContainsParameter(kalias+":script") && strings.Join(actn.rqst.Parameters().Parameter(kalias+":script"), "") == "false" {
										script = false
									}
									if actn.rqst.Parameters().ContainsParameter(kalias + ":jscall") {
										if jscall := strings.Join(actn.rqst.Parameters().Parameter(kalias+":jscall"), ""); jscall != "" {
											if actn.rqst.mimetype == "" {
												actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
											}
											if script {
												actn.rqst.copy(io.MultiReader(strings.NewReader("script||"), strings.NewReader(jscall+"("), database.NewJSONReader(dbrdr, nil, dbrdrerr), strings.NewReader(");"), strings.NewReader("||script")), nil, false)
											} else {
												actn.rqst.copy(io.MultiReader(strings.NewReader(jscall+"("), database.NewJSONReader(dbrdr, nil, dbrdrerr), strings.NewReader(");")), nil, false)
											}
										}
									} else if actn.rqst.Parameters().ContainsParameter(kalias + ":jsvar") {
										if jsvar := strings.Join(actn.rqst.Parameters().Parameter(kalias+":jsvar"), ""); jsvar != "" {
											if actn.rqst.mimetype == "" {
												actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
											}
											if script {
												actn.rqst.copy(io.MultiReader(strings.NewReader("script||"), strings.NewReader(jsvar+"="), database.NewJSONReader(dbrdr, nil, dbrdrerr), strings.NewReader(";"), strings.NewReader("||script")), nil, false)
											} else {
												actn.rqst.copy(io.MultiReader(strings.NewReader(jsvar+"="), database.NewJSONReader(dbrdr, nil, dbrdrerr), strings.NewReader(";")), nil, false)
											}
										}
									}
								} else if rspathext == ".csv" {
									var csvsttngs = map[string]interface{}{}

									if actn.rqst.Parameters().ContainsParameter(kalias + ":csv") {
										var csvsttngsval = strings.Join(actn.rqst.Parameters().Parameter(kalias+":csv"), "")
										var decsttngs = json.NewDecoder(strings.NewReader(csvsttngsval))
										if decerr := decsttngs.Decode(&csvsttngs); decerr == nil {

										}
									}
									actn.rqst.copy(io.MultiReader(database.NewCSVReader(dbrdr, dbrdrerr, csvsttngs)), nil, false)
								}
							} else {

							}
							if dbrdr != nil {
								dbrdr.Close()
							}
						} else if actn.rqst.Parameters().ContainsParameter(kalias + ":execute") {
							exctr, exctrerr := dbcn.GblExecute(strings.Join(actn.rqst.Parameters().Parameter(kalias+":execute"), ""), actn.rqst.Parameters())
							if rspathext == "" {
								rspathext = ".json"
							}
							if rspathext != "" {
								if rspathext == ".json" {
									actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
									actn.rqst.copy(io.MultiReader(database.NewJSONReader(nil, exctr, exctrerr)), nil, false)
								} else if rspathext == ".js" {
									var script = true
									if actn.rqst.Parameters().ContainsParameter(kalias+":script") && strings.Join(actn.rqst.Parameters().Parameter(kalias+":script"), "") == "false" {
										script = false
									}
									if actn.rqst.Parameters().ContainsParameter(kalias + ":jscall") {
										if jscall := strings.Join(actn.rqst.Parameters().Parameter(kalias+":jscall"), ""); jscall != "" {
											if actn.rqst.mimetype == "" {
												actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
											}
											if script {
												actn.rqst.copy(io.MultiReader(strings.NewReader("script||"), strings.NewReader(jscall+"("), database.NewJSONReader(nil, exctr, exctrerr), strings.NewReader(");"), strings.NewReader("||script")), nil, false)
											} else {
												actn.rqst.copy(io.MultiReader(strings.NewReader(jscall+"("), database.NewJSONReader(nil, exctr, exctrerr), strings.NewReader(");")), nil, false)
											}
										}
									} else if actn.rqst.Parameters().ContainsParameter(kalias + ":jsvar") {
										if jsvar := strings.Join(actn.rqst.Parameters().Parameter(kalias+":jsvar"), ""); jsvar != "" {
											if actn.rqst.mimetype == "" {
												actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
											}
											if script {
												actn.rqst.copy(io.MultiReader(strings.NewReader("script||"), strings.NewReader(jsvar+"="), database.NewJSONReader(nil, exctr, exctrerr), strings.NewReader(";"), strings.NewReader("||script")), nil, false)
											} else {
												actn.rqst.copy(io.MultiReader(strings.NewReader(jsvar+"="), database.NewJSONReader(nil, exctr, exctrerr), strings.NewReader(";")), nil, false)
											}
										}
									}
								}
							} else {

							}
							if exctr != nil {
								exctr.Close()
							}
						} else {
							if rspathext == "" {
								rspathext = ".json"
							}
							if rspathext == ".json" {
								actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
								dbcn.InOut(actn.rqst.RequestBody(), actn.rqst, actn.rqst.Parameters())
							}
						}
					}
				}
			}
		}()
	} else {
		if curactnhndlr := actn.ActionHandler(); curactnhndlr == nil {
			if actn.rqst.isFirstRequest {
				actn.rqst.isFirstRequest = false
				if actn.rqst.mimetype == "" {
					actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
				}
				if rspath != "" {
					if strings.LastIndex(rspath, ".") == -1 {
						if !strings.HasSuffix(rspath, "/") {
							rspath = rspath + "/"
						}
						rspath = rspath + "index.html"
						actn.rspath = rspath
						actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
						if curactnhndlr = actn.ActionHandler(); curactnhndlr == nil {
							actn.rqst.mimetype = "text/plain"
							isTextRequest = false
						} else {
							func() {
								defer func() {
									curactnhndlr.Close()
									curactnhndlr = nil
								}()
								if isTextRequest {
									isTextRequest = false
									actn.rqst.copy(curactnhndlr, altw(), true, actn.rspath) // actn.rsngpth.Path)
								} else {
									actn.rqst.copy(curactnhndlr, nil, false, actn.rspath) // actn.rsngpth.Path)
								}
							}()
						}
					}
				}
			}
		} else if curactnhndlr != nil {
			func() {
				defer func() {
					curactnhndlr.Close()
					curactnhndlr = nil
				}()
				if actn.rqst.isFirstRequest {
					if actn.rqst.mimetype == "" {
						actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
					} else {
						_, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
					}
					actn.rqst.isFirstRequest = false
				} else {
					_, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
				}
				if isTextRequest {
					isTextRequest = false
					actn.rqst.copy(curactnhndlr, altw(), true, actn.rspath) // actn.rsngpth.Path)

				} else {
					actn.rqst.copy(curactnhndlr, nil, false, actn.rspath) // actn.rsngpth.Path)
				}
			}()
		}
	}
	return
}

//ActionHandler - handle individual action io
func (actn *Action) ActionHandler() (actnhndl *ActionHandler) {
	actnhndl = NewActionHandler(actn)
	return
}

//Close - action
func (actn *Action) Close() (err error) {
	if actn != nil {
		if actn.rqst != nil {
			actn.rqst.detachAction(actn)
			actn.rqst = nil
		}
		if actn.prvactn != nil {
			actn.prvactn = nil
		}
		/*if actn.rsngpth != nil {
			actn.rsngpth.Close()
			actn.rsngpth = nil
		}*/
		actn = nil
	}
	return
}

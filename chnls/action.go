package chnls

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"github.com/evocert/kwe/database"
	"github.com/evocert/kwe/mimes"
	"github.com/evocert/kwe/resources"
)

//Action - struct
type Action struct {
	rqst    *Request
	rsngpth *resources.ResourcingPath
	sttngs  map[string]interface{}
	prvactn *Action
}

func newAction(rqst *Request, rsngpth *resources.ResourcingPath) (actn *Action) {
	actn = &Action{rqst: rqst, rsngpth: rsngpth, prvactn: rqst.lstexctngactng}
	rqst.lstexctngactng = actn
	return
}

//Path of resource action is currently processing
func (actn *Action) Path() string {
	if actn != nil {
		if actn.rsngpth != nil {
			return actn.rsngpth.Path
		}
	}
	return ""
}

func executeAction(actn *Action) (err error) {
	defer func() {
		if actn != nil {
			actn.Close()
			actn = nil
		}
	}()
	var rspath = actn.rsngpth.Path
	var rspathext = filepath.Ext(rspath)
	var isTextRequest = false
	var aliases map[string]*database.Connection = nil
	if strings.HasPrefix(rspath, "/dbms/") || strings.HasPrefix(rspath, "/dbms-") {
		func() {
			defer func() {
				if aliases != nil {
					if aliasesl := len(aliases); aliasesl > 0 {
						aliasks := make([]string, aliasesl)
						aliasksi := 0
						for aliask := range aliases {
							aliasks[aliasksi] = aliask
							aliasksi++
						}
						for _, aliask := range aliasks {
							aliases[aliask] = nil
							delete(aliases, aliask)
						}
					}
					aliases = nil
				}
				if actn != nil {
					if rspth := actn.rsngpth.Path; rspth != "" {
						if _, ok := actn.rqst.rsngpthsref[rspth]; ok {
							actn.rqst.rsngpthsref[rspth] = nil
							delete(actn.rqst.rsngpthsref, rspth)
						}
					}
					actn.Close()
					actn = nil
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
									if aliases == nil {
										aliases = map[string]*database.Connection{}
									}
									aliases[kalias] = dbcn
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
						if aliases == nil {
							aliases = map[string]*database.Connection{}
						}
						aliases[alias] = dbcn
					}
				}
				if len(aliases) > 0 {
					for kalias, dbcn := range aliases {
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
			if rspth := actn.rsngpth.Path; rspth != "" {
				if _, ok := actn.rqst.rsngpthsref[rspth]; ok {
					actn.rqst.rsngpthsref[rspth] = nil
					delete(actn.rqst.rsngpthsref, rspth)
				}
			}
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
						actn.rsngpth.Path = rspath
						actn.rsngpth.LookupPath = actn.rsngpth.Path
						actn.rqst.mimetype, isTextRequest = mimes.FindMimeType(rspath, "text/plain")
						if curactnhndlr = actn.ActionHandler(); curactnhndlr == nil {
							actn.rqst.mimetype = "text/plain"
							isTextRequest = false
						} else {
							actn.rqst.rsngpthsref[actn.rsngpth.Path] = actn.rsngpth
							if isTextRequest && actn.rsngpth.Path != actn.rsngpth.LookupPath {
								isTextRequest = false
							}
							if isTextRequest {
								isTextRequest = false
								actn.rqst.copy(curactnhndlr, nil, true)
							} else {
								actn.rqst.copy(curactnhndlr, nil, false)
							}
							curactnhndlr.Close()
							curactnhndlr = nil
						}
					} else {
						actn.Close()
					}
				} else {
					actn.Close()
				}
			} else {
				actn.Close()
			}
			actn = nil
		} else if curactnhndlr != nil {
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
			actn.rqst.rsngpthsref[actn.rsngpth.Path] = actn.rsngpth
			if isTextRequest && actn.rsngpth.Path != actn.rsngpth.LookupPath {
				isTextRequest = false
			}
			if isTextRequest {
				isTextRequest = false
				actn.rqst.copy(curactnhndlr, nil, true)
			} else {
				actn.rqst.copy(curactnhndlr, nil, false)
			}
			if curactnhndlr != nil {
				curactnhndlr.Close()
				curactnhndlr = nil
			}
			actn.Close()
			actn = nil
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
		if actn.rsngpth != nil {
			actn.rsngpth.Close()
			actn.rsngpth = nil
		}
		actn = nil
	}
	return
}

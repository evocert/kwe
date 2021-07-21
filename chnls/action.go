package chnls

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"github.com/evocert/kwe/database"
	"github.com/evocert/kwe/mimes"
)

//Action - struct
type Action struct {
	rqst    *Request
	rspath  string
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

func executeAction(actn *Action) (err error) {
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
	if strings.HasPrefix(rspath, "/dbms/") || strings.HasPrefix(rspath, "/dbms-") {
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
			var mimetype = ""
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
							mimetype, isTextRequest = mimes.FindMimeType(".json", "text/plain")
							if actn.rqst.rqstrw.Response().ContentType() == "" {
								actn.rqst.rqstrw.Response().SetContentType(mimetype)
							}
							var jsnr io.Reader = nil
							if actn.rqst.Parameters().ContainsParameter("dbms:json") {
								jsnr = strings.NewReader(strings.Join(actn.rqst.Parameters().Parameter("dbms:json"), ""))
							} else {
								jsnr = actn.rqst.RequestBody()
							}
							if jsnr != nil {
								database.GLOBALDBMS().InOut(jsnr, actn.rqst.rqstrw.Response(), actn.rqst.Parameters())
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
							mimetype, isTextRequest = mimes.FindMimeType(".json", "text/plain")
							if actn.rqst.rqstrw.Response().ContentType() == "" {
								actn.rqst.rqstrw.Response().SetContentType(mimetype)
							}
							dbcn.InOut(strings.NewReader(jsnval), actn.rqst.rqstrw.Response(), actn.rqst.Parameters())
						} else if actn.rqst.Parameters().ContainsParameter(kalias + ":query") {
							dbrdr, dbrdrerr := dbcn.GblQuery(strings.Join(actn.rqst.Parameters().Parameter(kalias+":query"), ""), actn.rqst.Parameters())
							if rspathext != "" {
								if rspathext == ".json" {
									mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
									if actn.rqst.rqstrw.Response().ContentType() == "" {
										actn.rqst.rqstrw.Response().SetContentType(mimetype)
									}
									actn.rqst.copy(io.MultiReader(database.NewJSONReader(dbrdr, nil, dbrdrerr)), nil, false, false, "")
								} else if rspathext == ".js" {
									var script = true
									if actn.rqst.Parameters().ContainsParameter(kalias+":script") && strings.Join(actn.rqst.Parameters().Parameter(kalias+":script"), "") == "false" {
										script = false
									}
									if actn.rqst.Parameters().ContainsParameter(kalias + ":jscall") {
										if jscall := strings.Join(actn.rqst.Parameters().Parameter(kalias+":jscall"), ""); jscall != "" {
											if mimetype == "" {
												mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
											}
											if actn.rqst.rqstrw.Response().ContentType() == "" {
												actn.rqst.rqstrw.Response().SetContentType(mimetype)
											}
											if script {
												actn.rqst.copy(io.MultiReader(strings.NewReader("script||"), strings.NewReader(jscall+"("), database.NewJSONReader(dbrdr, nil, dbrdrerr), strings.NewReader(");"), strings.NewReader("||script")), nil, false, false, "")
											} else {
												actn.rqst.copy(io.MultiReader(strings.NewReader(jscall+"("), database.NewJSONReader(dbrdr, nil, dbrdrerr), strings.NewReader(");")), nil, false, false, "")
											}
										}
									} else if actn.rqst.Parameters().ContainsParameter(kalias + ":jsvar") {
										if jsvar := strings.Join(actn.rqst.Parameters().Parameter(kalias+":jsvar"), ""); jsvar != "" {
											if mimetype == "" {
												mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
											}
											if actn.rqst.rqstrw.Response().ContentType() == "" {
												actn.rqst.rqstrw.Response().SetContentType(mimetype)
											}
											if script {
												actn.rqst.copy(io.MultiReader(strings.NewReader("script||"), strings.NewReader(jsvar+"="), database.NewJSONReader(dbrdr, nil, dbrdrerr), strings.NewReader(";"), strings.NewReader("||script")), nil, false, false, "")
											} else {
												actn.rqst.copy(io.MultiReader(strings.NewReader(jsvar+"="), database.NewJSONReader(dbrdr, nil, dbrdrerr), strings.NewReader(";")), nil, false, false, "")
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
									actn.rqst.copy(io.MultiReader(database.NewCSVReader(dbrdr, dbrdrerr, csvsttngs)), nil, false, false, "")
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
									mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
									if actn.rqst.rqstrw.Response().ContentType() == "" {
										actn.rqst.rqstrw.Response().SetContentType(mimetype)
									}
									actn.rqst.copy(io.MultiReader(database.NewJSONReader(nil, exctr, exctrerr)), nil, false, false, "")
								} else if rspathext == ".js" {
									var script = true
									if actn.rqst.Parameters().ContainsParameter(kalias+":script") && strings.Join(actn.rqst.Parameters().Parameter(kalias+":script"), "") == "false" {
										script = false
									}
									if actn.rqst.Parameters().ContainsParameter(kalias + ":jscall") {
										if jscall := strings.Join(actn.rqst.Parameters().Parameter(kalias+":jscall"), ""); jscall != "" {
											if mimetype == "" {
												mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
											}
											if actn.rqst.rqstrw.Response().ContentType() == "" {
												actn.rqst.rqstrw.Response().SetContentType(mimetype)
											}
											if script {
												actn.rqst.copy(io.MultiReader(strings.NewReader("script||"), strings.NewReader(jscall+"("), database.NewJSONReader(nil, exctr, exctrerr), strings.NewReader(");"), strings.NewReader("||script")), nil, false, false, "")
											} else {
												actn.rqst.copy(io.MultiReader(strings.NewReader(jscall+"("), database.NewJSONReader(nil, exctr, exctrerr), strings.NewReader(");")), nil, false, false, "")
											}
										}
									} else if actn.rqst.Parameters().ContainsParameter(kalias + ":jsvar") {
										if jsvar := strings.Join(actn.rqst.Parameters().Parameter(kalias+":jsvar"), ""); jsvar != "" {
											if mimetype == "" {
												mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
											}
											if actn.rqst.rqstrw.Response().ContentType() == "" {
												actn.rqst.rqstrw.Response().SetContentType(mimetype)
											}
											if script {
												actn.rqst.copy(io.MultiReader(strings.NewReader("script||"), strings.NewReader(jsvar+"="), database.NewJSONReader(nil, exctr, exctrerr), strings.NewReader(";"), strings.NewReader("||script")), nil, false, false, "")
											} else {
												actn.rqst.copy(io.MultiReader(strings.NewReader(jsvar+"="), database.NewJSONReader(nil, exctr, exctrerr), strings.NewReader(";")), nil, false, false, "")
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
								mimetype, isTextRequest = mimes.FindMimeType(rspathext, "text/plain")
								if actn.rqst.rqstrw.Response().ContentType() == "" {
									actn.rqst.rqstrw.Response().SetContentType(mimetype)
								}
								dbcn.InOut(actn.rqst, actn.rqst.rqstrw.Response(), actn.rqst.Parameters())
							}
						}
					}
				}
			}
		}()
	} else {
		if curactnhndlr := actn.ActionHandler(); curactnhndlr != nil {
			func() {
				defer func() {
					curactnhndlr.Close()
					curactnhndlr = nil
				}()
				if actn.rqst.isFirstRequest {
					if actn.rqst.rqstrw.Response().ContentType() == "" {
						if curactnhndlr.raw {
							actn.rqst.rqstrw.Response().SetContentType("text/plain")
						} else {
							mimetype := ""
							mimetype, isTextRequest = mimes.FindMimeType(actn.rspath, "text/plain")
							actn.rqst.rqstrw.Response().SetContentType(mimetype)
						}
					} else {
						_, isTextRequest = mimes.FindMimeType(actn.rspath, "text/plain")
					}
					actn.rqst.isFirstRequest = false
				} else {
					_, isTextRequest = mimes.FindMimeType(actn.rspath, "text/plain")
				}
				if isTextRequest && !curactnhndlr.raw {
					isTextRequest = false
					actn.rqst.copy(curactnhndlr, nil, true, curactnhndlr.active && !curactnhndlr.raw, actn.rspath)
				} else {
					actn.rqst.copy(curactnhndlr, nil, false, false, actn.rspath)
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

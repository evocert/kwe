package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/dop251/goja"
	_ "github.com/evocert/kwe/alertify"
	_ "github.com/evocert/kwe/bootstrap"
	"github.com/evocert/kwe/caching"
	"github.com/evocert/kwe/database"
	_ "github.com/evocert/kwe/datepicker"
	"github.com/evocert/kwe/enumeration"
	"github.com/evocert/kwe/env"
	_ "github.com/evocert/kwe/fonts/material"
	_ "github.com/evocert/kwe/fonts/robotov27latin"
	"github.com/evocert/kwe/fsutils"
	_ "github.com/evocert/kwe/goldenlayout"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
	_ "github.com/evocert/kwe/jspanel"
	"github.com/evocert/kwe/listen"
	"github.com/evocert/kwe/mimes"
	"github.com/evocert/kwe/mqtt"
	"github.com/evocert/kwe/osprc"
	"github.com/evocert/kwe/requirejs"
	scheduling "github.com/evocert/kwe/scheduling/ext"
	"github.com/evocert/kwe/service"
	_ "github.com/evocert/kwe/sip"
	_ "github.com/evocert/kwe/typescript"
	"github.com/evocert/kwe/web"

	_ "github.com/evocert/kwe/database/mysql"
	//_ "github.com/evocert/kwe/database/ora"
	_ "github.com/evocert/kwe/database/postgres"
	_ "github.com/evocert/kwe/database/sqlserver"

	"github.com/evocert/kwe/requesting"
	"github.com/evocert/kwe/resources"
	_ "github.com/evocert/kwe/webactions"
)

type exepath struct {
	path string
	args []interface{}
}

func (expth *exepath) Path() string {
	if expth != nil {
		if strings.LastIndex(expth.path, "/") < strings.Index(expth.path, "?") {
			return expth.path[:strings.Index(expth.path, "?")]
		} else {
			return expth.path
		}
	}
	return ""
}

func (expth *exepath) PathRoot() (pathroot string) {
	if expth != nil {
		if strings.LastIndex(expth.path, "/") > -1 {
			pathroot = expth.path[0 : strings.LastIndex(expth.path, "/")+1]
		} else {
			pathroot = "/"
		}
	}
	return
}

func (expth *exepath) Args() (args []interface{}) {
	if expth != nil {
		args = expth.args
	}
	return
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU() * 100)
	var serveRequest func(
		rqst requesting.RequestAPI,
		mqttmsg mqtt.Message,
		mqttevent mqtt.MqttEvent,
		atv *active.Active,
		schdl scheduling.ScheduleAPI,
		a ...interface{}) (err error) = nil
	lstnr := listen.NewListener()
	var glblutilsfs = fsutils.NewFSUtils()
	var glbldbms = database.GLOBALDBMS
	var glblrsfs = resources.GLOBALRSNG().FS
	glblrsfs().MKDIR("/require/js", "")
	glblrsfs().SET("/require/js/require.js", requirejs.RequireJS())
	var glblchng = caching.GLOBALMAPHANDLER

	var prepScheduleActionArgs func(scheduling.ScheduleAPI, ...interface{}) ([]interface{}, error) = nil

	var glblschdlng = func() *scheduling.Schedules {
		return scheduling.GLOBALSCHEDULES(func(schdl scheduling.ScheduleAPI, a ...interface{}) (nargs []interface{}, err error) {
			if prepScheduleActionArgs != nil {
				nargs, err = prepScheduleActionArgs(schdl, a...)
			}
			return
		}, func(rqst requesting.RequestAPI, atv *active.Active, schdl scheduling.ScheduleAPI, a ...interface{}) error {
			return serveRequest(rqst, nil, nil, atv, nil, a...)
		})
	}
	var glblenv = env.Env()
	active.LoadGlobalModule("kwe.js", sysjsTemplate("kwe",
		map[string]interface{}{
			"env":         "_env",
			"command":     "_command",
			"path":        "_path",
			"in":          "_in",
			"dbms":        "_dbms",
			"caching":     "_caching",
			"out":         "_out",
			"fs":          "_fs",
			"fsutils":     "_fsutils",
			"scheduling":  "_scheduling",
			"schdl":       "_schdl",
			"mqtting":     "_mqtting",
			"mqttmsg":     "_mqttmsg",
			"mqttevent":   "_mqttevent",
			"extMimetype": "_extMimetype",
			"addPath":     "_addpath",
			"send":        "_send",
			"sendEval":    "_sendeval",
			"sendreceive": "_sendreceive",
			"listen":      "_listen"}))
	var glblmqttng = mqtt.NewMQTTManager(mqtt.MqttEventing(func(event mqtt.MqttEvent) {
		if rqst := requesting.NewRequest(nil, event.EventPath()); rqst != nil {
			defer rqst.Close()
			if lstnr.ServeRequest != nil {
				lstnr.ServeRequest(rqst, event)
			}
		}
	}), mqtt.MqttMessaging(func(message mqtt.Message) {
		if rqst := requesting.NewRequest(nil, message.TopicPath()); rqst != nil {
			defer rqst.Close()
			if lstnr.ServeRequest != nil {
				lstnr.ServeRequest(rqst, message)
			}
		}
	}))

	var executeSchdlDbms = func(schdl scheduling.ScheduleAPI, a ...interface{}) (err error) {
		if len(a) > 0 {
			err = database.GLOBALDBMS().InOut(a[0], nil, a[1:]...)
		}
		return
	}

	var executeSchdlCommand = func(schdl scheduling.ScheduleAPI, a ...interface{}) (err error) {
		if len(a) > 0 {

		}
		return
	}

	var executeSchdlRequest = func(schdl scheduling.ScheduleAPI, a ...interface{}) (err error) {
		if serveRequest != nil {
			rqst := requesting.NewRequest(nil, a...)
			defer rqst.Close()
			serveRequest(rqst, nil, nil, schdl.Active(), schdl)
		}
		return
	}

	var executeScheduleAction = func(a ...interface{}) (err error) {
		scdl, _ := a[0].(scheduling.ScheduleAPI)
		a = a[1:]
		for len(a) > 1 {
			if cmd, cmdok := a[0].(string); cmdok && cmd != "" {
				a = a[1:]
				var cmdfnctoexec func(scheduling.ScheduleAPI, ...interface{}) error = nil
				if cmd == "dbms" {
					cmdfnctoexec = executeSchdlDbms
				} else if cmd == "request" {
					cmdfnctoexec = executeSchdlRequest
				} else if cmd == "command" {
					cmdfnctoexec = executeSchdlCommand
				} else if cmd == "script" {
					//cmdfnctoexec = rqst.executeSchdlScript
				}
				if cmdfnctoexec != nil {
					if cmdmap, cmdmapok := a[0].(map[string]interface{}); cmdmapok && len(cmdmap) > 0 {
						cmdfnctoexec(scdl, cmdmap)
					} else if cmdargs, cmdargsok := a[0].([]interface{}); cmdargsok && len(cmdargs) > 0 {
						cmdfnctoexec(scdl, cmdargs)
					} else if cmdarg := a[0]; cmdarg == nil || cmdarg != nil {
						cmdfnctoexec(scdl, cmdarg)
					} else {
						break
					}
					a = a[1:]
				} else {
					break
				}
			} else {
				break
			}
		}

		return
	}

	prepScheduleActionArgs = func(scdl scheduling.ScheduleAPI, a ...interface{}) (preppedargs []interface{}, err error) {
		if al := len(a); al > 0 {
			ai := 0
			var schdlatv = func() (atv *active.Active) {
				atv = scdl.Active()
				if atv.ObjectMapRef == nil {
					if serveRequest != nil {
						rqst := requesting.NewRequest(nil, "dummy.js")
						defer rqst.Close()
						serveRequest(rqst, nil, nil, atv, scdl)
					}
				}
				return
			}
			regatvfnc := func(atvfnc func(goja.FunctionCall) goja.Value) bool {
				if atvfnc != nil {
					var prppdatvfnc scheduling.FuncArgsErrHandle = nil
					prppdatvfnc = func(args ...interface{}) (rserr error) {
						atv := schdlatv()
						if rslt := atv.InvokeFunction(atvfnc, args...); rslt != nil {
							if dne, dneok := rslt.(bool); dneok && dne {
								rserr = fmt.Errorf("DONE")
							}
						}
						return
					}
					a[ai] = nil
					a[ai] = prppdatvfnc
					return true
				}
				return false
			}
			for ai < al {
				d := a[ai]
				if sfnc, sfncok := d.(string); sfncok {
					if sfnc != "" {
						if strings.HasPrefix(sfnc, "function(") {
							schdlatv().InvokeVM(func(vm *goja.Runtime) (vmerr error) {
								atvfncval, _ := vm.RunString("(" + sfnc + ")")
								var atvfncref func(goja.FunctionCall) goja.Value = nil
								vm.ExportTo(atvfncval, &atvfncref)
								if !regatvfnc(atvfncref) {

								}
								return
							})
						}
					}
				} else if rqstactnmap, rqstactnmapok := d.(map[string]interface{}); rqstactnmapok {
					ignore := false
					for rqstmk, rqstmv := range rqstactnmap {
						if rqstmk != "" && strings.Contains("|request|dbms|command|script|", "|"+rqstmk+"|") {
							if !ignore {
								ignore = true
							}
							a[ai] = scheduling.FuncArgsErrHandle(executeScheduleAction)
							tmpa := append([]interface{}{scdl}, a[ai+1:])
							a = append(append(a[:ai+1], []interface{}{rqstmk, rqstmv}), tmpa...)
							al = len(a)
							ai++
							ai++
						}
					}
					if ignore {

						continue
					}
				}
				ai++
			}
			preppedargs = a[:]
		}
		return
	}

	serveRequest = func(
		rqst requesting.RequestAPI,
		mqttmsg mqtt.Message,
		mqttevent mqtt.MqttEvent,
		atv *active.Active,
		schdl scheduling.ScheduleAPI,
		a ...interface{}) (err error) {
		defer func() { runtime.GC() }()
		var exitingatv = atv != nil

		var invokeCommand func(execpath string, execargs ...string) (cmd *osprc.Command, err error) = nil

		var cmnds map[int]*osprc.Command = nil

		var closecmd = func(prcid int) {
			if cmnds != nil {
				if cmdf := cmnds[prcid]; cmdf != nil {
					cmdf.OnClose = nil
					delete(cmnds, prcid)
				}
			}
		}

		invokeCommand = func(execpath string, execargs ...string) (cmd *osprc.Command, err error) {
			cmd, err = osprc.NewCommand(execpath, execargs...)

			if err == nil && cmd != nil {
				cmd.OnClose = closecmd
				if cmnds == nil {
					cmnds = map[int]*osprc.Command{}
				}
				cmnds[cmd.PrcID()] = cmd
			}
			return
		}

		defer func() {
			if len(cmnds) > 0 {
				prcsids := make([]int, len(cmnds))
				prcsidsi := 0
				for prcid := range cmnds {
					prcsids[prcsidsi] = prcid
					prcsidsi++
				}

				for _, prcid := range prcsids {
					cmnds[prcid].Close()
				}
				cmnds = nil
			}
		}()

		if len(a) > 0 {
			for _, d := range a {
				if mqttmsgd, _ := d.(mqtt.Message); mqttmsgd != nil && mqttmsg == nil {
					mqttmsg = mqttmsgd
				} else if mqtteventd, _ := d.(mqtt.MqttEvent); mqtteventd != nil && mqttevent == nil {
					mqttevent = mqtteventd
				} else if atvd, _ := d.(*active.Active); atvd != nil {
					if atv == nil {
						atv = atvd
						exitingatv = atv != nil
					}
				}
			}
		}

		rspns := rqst.Response()

		var rqstdpaths *enumeration.List = enumeration.NewList()
		var addNextPath = func(nxtpth ...string) {
			if nxtpthl := len(nxtpth); nxtpthl > 0 {
				nxtpthi := 0
				var nxtToAdd []string = nil
				for nxtpthi < nxtpthl {
					if nxtp := strings.TrimSpace(nxtpth[nxtpthi]); nxtp != "" {
						for _, nxp := range strings.Split(nxtp, "|") {
							if nxp != "" {
								if nxtToAdd == nil {
									nxtToAdd = []string{}
								}
								nxtToAdd = append(nxtToAdd, nxp)
							}
						}
					}
					nxtpthi++
				}
				if len(nxtToAdd) > 0 {
					for _, nxttadd := range nxtToAdd {
						rqstdpaths.InsertAfter(nil, nil, rqstdpaths.CurrentDoing(), &exepath{path: nxttadd})
					}
					nxtToAdd = nil
				}
			}
		}

		if ppath := rqst.Path(); ppath != "" {
			addNextPath(ppath)
			func() {
				defer func() {
					if atv != nil && !exitingatv {
						atv.Close()
					}
				}()
				var crntexpths *enumeration.List = nil
				defer func() {
					if crntexpths != nil {
						crntexpths.Dispose(nil, nil)
						crntexpths = nil
					}
				}()
				var processPath = func(expth *exepath) (err error) {
					if expth == nil {
						return
					}
					if crntexpths == nil {
						crntexpths = enumeration.NewList()
					}
					crntexpths.Push(nil, nil, expth)
					var cnrtextpthsnd = crntexpths.Tail()
					defer func() {
						if crntexpths != nil && cnrtextpthsnd != nil {
							cnrtextpthsnd.Dispose(nil, nil)
						}
					}()
					var path = expth.Path()
					var convertactive bool = false
					var israw = false
					var mimetype, isactive = mimes.FindMimeType(path, "text/plain")
					if israw = strings.Contains(path, "/raw:"); israw {
						path = strings.Replace(path, "/raw:", "/", 1)
						isactive = !israw
					}
					if convertactive = strings.Contains(path, "/active:"); convertactive {
						path = strings.Replace(path, "/active:", "/", 1)
						if !israw && !isactive {
							isactive = convertactive
						}
					}
					if strings.HasSuffix(path, "/typescript.js") {
						isactive = false
					}
					if rspns != nil {
						rspns.SetHeader("Content-Type", mimetype)
					}
					var rs io.Reader = nil
					if rs = glblrsfs().CAT(path); rs == nil && (strings.LastIndex(path, ".") == -1 || strings.LastIndex(path, "/") > strings.LastIndex(path, ".")) {
						if !strings.HasSuffix(path, "/") {
							path += "/"
						}
						for _, pth := range strings.Split("html,xml,svg,js,json,css", ",") {
							if rs = glblrsfs().CAT(path + "index" + "." + pth); rs == nil {
								if rs = glblrsfs().CAT(path + "main" + "." + pth); rs == nil {
									continue
								} else {
									mimetype, isactive = mimes.FindMimeType(path+"main"+"."+pth, "text/plain")
									if rspns != nil {
										rspns.SetHeader("Content-Type", mimetype)
									}
									break
								}
							} else {
								mimetype, isactive = mimes.FindMimeType(path+"index"+"."+pth, "text/plain")
								if rspns != nil {
									rspns.SetHeader("Content-Type", mimetype)
								}
								break
							}
						}
					}
					if rs == nil && path == "dummy.js" {
						rs = iorw.NewEOFCloseSeekReader(strings.NewReader("<@ /**/ @>"))
					}
					if rs != nil {
						if isactive {
							if atv == nil {
								atv = active.NewActive()
							}
							if atv.LookupTemplate == nil {
								atv.LookupTemplate = func(lkppath string, a ...interface{}) (lkpr io.Reader, lkperr error) {
									if glblrsfs != nil {
										if lkppath != "" && (strings.HasSuffix(lkppath, ".js") || strings.HasSuffix(lkppath, ".html") || strings.HasSuffix(lkppath, ".xml") || strings.HasSuffix(lkppath, ".svg")) {
											if !strings.HasPrefix(lkppath, "/") {
												if crntexpths != nil && crntexpths.Length() > 0 {
													if val := crntexpths.Tail().Value(); val != nil {
														if dngexpth, _ := val.(*exepath); dngexpth != nil {
															lkppath = dngexpth.PathRoot() + lkppath
														} else {
															lkppath = expth.PathRoot() + lkppath
														}
													} else {
														lkppath = expth.PathRoot() + lkppath
													}
												} else {
													lkppath = expth.PathRoot() + lkppath
												}
											}
											if lkpr = glblrsfs().CAT(lkppath); lkpr == nil {
												lkpr = glblutilsfs.CAT(lkppath)
											}
										}
									}
									return
								}
							}
							if atv.ObjectMapRef == nil {
								atv.ObjectMapRef = func() (objrf map[string]interface{}) {
									var objref = map[string]interface{}{}
									objref["_command"] = invokeCommand
									objref["_path"] = func() *exepath {
										if crntexpths != nil && crntexpths.Length() > 0 {
											if val := crntexpths.Tail().Value(); val != nil {
												if dngexpth, _ := val.(*exepath); dngexpth != nil {
													return dngexpth
												}
											} else {
												return expth
											}
										}
										return expth
									}
									objref["_env"] = glblenv
									objref["_in"] = rqst
									objref["_dbms"] = glbldbms().ActiveDBMS(atv, rqst.Parameters())
									objref["_caching"] = glblchng().ActiveHandler(atv, rqst.Parameters())
									objref["_out"] = rspns
									objref["_fs"] = glblrsfs()
									objref["_fsutils"] = glblutilsfs
									objref["_scheduling"] = glblschdlng().ActiveSCHEDULING(atv, rqst.Parameters())
									objref["_schdl"] = schdl
									objref["_mqtting"] = glblmqttng
									objref["_mqttmsg"] = mqttmsg
									objref["_mqttevent"] = mqttevent
									objref["_extMimetype"] = mimes.ExtMimeType
									objref["_addpath"] = addNextPath
									objref["_send"] = func(rqstpath string, a ...interface{}) (rdr iorw.Reader, err error) {
										return send(atv, glblrsfs(), rqst, rqstpath, false, a...)
									}
									objref["_sendeval"] = func(rqstpath string, a ...interface{}) (rdr iorw.Reader, err error) {
										return send(atv, glblrsfs(), rqst, rqstpath, true, a...)
									}
									objref["_sendreceive"] = func(rqstpath string, a ...interface{}) (rdr iorw.PrinterReader, err error) {
										return sendreceive(atv, glblrsfs(), rqst, rqstpath, a...)
									}
									objref["_listen"] = func(addr ...string) {
										lstnr.Listen("tcp", addr...)
									}
									objrf = objref
									return
								}
							}
							func() {
								var evalerr error = nil
								evalerr = atv.Eval(rspns, rqst, path, convertactive, rs)
								if evalerr != nil {
									if rspns != nil {
										rspns.SetHeader("Content-Type", "application/javascript")
										rspns.SetStatus(500)
										rspns.Print(evalerr)
									} else {
										println(evalerr.Error())
									}
								}
							}()
						} else if rspns != nil {
							rspns.Print(rs)
						}
					}
					return
				}

				rqstdpaths.Do(
					func(nde *enumeration.Node, val interface{}) (donepath bool, doneerr error) {
						if err == nil {
							donepath = true
							var expath, _ = val.(*exepath)
							defer func() {
								if expath != nil {
									expath.args = nil
									expath = nil
								}
							}()
							if doneerr = processPath(expath); doneerr != nil && err == nil {
								err = doneerr
							}
							nde.Set(nil)
						} else {
							donepath = true
						}
						return
					}, nil, nil, nil)
			}()
		}
		return
	}

	lstnr.ServeRequest = func(rqst requesting.RequestAPI, a ...interface{}) error {
		return serveRequest(rqst, nil, nil, nil, nil, a...)
	}
	service.ServeRequest = func(rqst requesting.RequestAPI, a ...interface{}) error {
		return serveRequest(rqst, nil, nil, nil, nil, a...)
	}
	service.RunService(os.Args...)
}

func send(atv *active.Active, fs *fsutils.FSUtils, rqst requesting.RequestAPI, rqstpath string, andeval bool, a ...interface{}) (rdr iorw.Reader, err error) {
	convertactive := false
	if convertactive = strings.Contains(rqstpath, "/active:"); convertactive {
		rqstpath = strings.Replace(rqstpath, "/active:", "/", 1)
	}
	if strings.HasPrefix(rqstpath, "http://") || strings.HasPrefix(rqstpath, "https://") {
		a = append([]interface{}{atv}, a...)
		rdr, err = web.DefaultClient.Send(rqstpath, a...)
	} else {
		rdr = iorw.NewEOFCloseSeekReader(fs.CAT(rqstpath))
	}
	if rdr != nil && andeval && err == nil {
		ctx, ctxcancel := context.WithCancel(context.Background())
		pr, pw := io.Pipe()
		go func() {
			ctxcancel()
			var pwerr error = nil
			pwerr = atv.Eval(pw, nil, rqstpath, convertactive, rdr)
			defer func() {
				if pwerr == nil {
					pw.Close()
				} else {
					pw.CloseWithError(pwerr)
				}
			}()
		}()
		<-ctx.Done()
		rdr = iorw.NewEOFCloseSeekReader(pr)
	}
	return
}

func sendreceive(atv *active.Active, fs *fsutils.FSUtils, rqst requesting.RequestAPI, rqstpath string, a ...interface{}) (rdrwtr iorw.PrinterReader, err error) {
	if strings.HasPrefix(rqstpath, "ws://") || strings.HasPrefix(rqstpath, "wss://") {
		a = append([]interface{}{atv}, a...)
		rdrwtr, err = web.DefaultClient.SendReceive(rqstpath, a...)
	}
	return
}

func sysjsTemplate(nmspace string, objmap ...map[string]interface{}) (sysjscode string) {
	var objmptolistcode = func() (cde string) {
		if len(objmap) > 0 && objmap[0] != nil {
			for objk, objv := range objmap[0] {
				if objv != nil && objk != "" {
					if objs, _ := objv.(string); objs != "" {
						cde += `if (` + objs + `!==undefined && ` + objs + `!==null) {obj` + nmspace + `["` + objk + `"]=` + objs + `;}`
					}
				}
			}
		}
		return
	}
	sysjscode = "function " + nmspace + `(){
	var obj` + nmspace + `={};
	obj` + nmspace + `.methods = (obj) => {
		let properties = new Set()
		let currentObj = obj
		Object.entries(currentObj).forEach((key)=>{
			key=(key=(key+"")).indexOf(",")>0?key.substring(0,key.indexOf(',')):key;
			if (typeof currentObj[key] === 'function') {
				var item=key;
				properties.add(item);
			}
		});
		if (properties.size===0) {
			do {
				Object.getOwnPropertyNames(currentObj).map(item => properties.add(item))
			} while ((currentObj = Object.getPrototypeOf(currentObj)))
		}
		return [...properties.keys()].filter(item => typeof obj[item] === 'function')
	}

	obj` + nmspace + `.fields = (obj) => {
		let properties = new Set()
		let currentObj = obj
		Object.entries(currentObj).forEach((key)=>{
			key=(key=(key+"")).indexOf(",")>0?key.substring(0,key.indexOf(',')):key;
			if (typeof currentObj[key] !== 'function') {
				var item=key;
				properties.add(item);
			}
		});
		if (properties.size===0) {
			do {
				Object.getOwnPropertyNames(currentObj).map(item => properties.add(item))
			} while ((currentObj = Object.getPrototypeOf(currentObj)))
		}
		return [...properties.keys()].filter(item => item!=='__proto__' && typeof obj[item] !== 'function')
	}
	` + objmptolistcode() + `
	return obj` + nmspace + `;
}

if (typeof this.` + nmspace + `==="function") {
	this.` + nmspace + `=this.` + nmspace + `();
}
`
	return
}

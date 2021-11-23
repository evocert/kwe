package channeling

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/evocert/kwe/api"
	"github.com/evocert/kwe/caching"
	"github.com/evocert/kwe/database"
	"github.com/evocert/kwe/enumeration"
	"github.com/evocert/kwe/env"
	"github.com/evocert/kwe/fsutils"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/listen"
	"github.com/evocert/kwe/mimes"
	"github.com/evocert/kwe/mqtt"
	"github.com/evocert/kwe/osprc"
	"github.com/evocert/kwe/parameters"
	"github.com/evocert/kwe/requesting"
	"github.com/evocert/kwe/requirejs"
	"github.com/evocert/kwe/resources"
	"github.com/evocert/kwe/web"
)

func NewSession(a ...interface{}) (session api.SessionAPI) {
	var mqttmsg mqtt.Message
	var mqttevent mqtt.MqttEvent
	var mqttmngr mqtt.MQTTManagerAPI
	var atvschdlngmngr api.ActiveSchedulesAPI = nil
	var schdlngmngr api.SchedulesAPI = nil
	var schdl api.ScheduleAPI = nil
	var lstnr listen.ListenerAPI = nil
	//var atv *active.Active = nil
	if len(a) > 0 {
		for _, d := range a {
			if d != nil {
				if mqtteventd, _ := d.(mqtt.MqttEvent); mqtteventd != nil {
					if mqttevent == nil {
						mqttevent = mqtteventd
					}
				} else if mqttmsgd, _ := d.(mqtt.Message); mqttmsgd != nil {
					if mqttmsg == nil {
						mqttmsg = mqttmsgd
					}
				} else if mqttmngrd, _ := d.(mqtt.MQTTManagerAPI); mqttmngrd != nil {
					if mqttmngr == nil {
						mqttmngr = mqttmngrd
					}
				} else if atvschdlngmngrd, _ := d.(api.ActiveSchedulesAPI); atvschdlngmngrd != nil {
					if atvschdlngmngr == nil {
						atvschdlngmngr = atvschdlngmngrd
					}
				} else if schdlngmngrd, _ := d.(api.SchedulesAPI); schdlngmngrd != nil {
					if schdlngmngr == nil {
						schdlngmngr = schdlngmngrd
					}
				} else if schdld, _ := d.(api.ScheduleAPI); schdld != nil {
					if schdl == nil {
						schdl = schdld
					}
				} else if lstnrd, _ := d.(listen.ListenerAPI); lstnrd != nil {
					if lstnr == nil {
						lstnr = lstnrd
					}
				}
			}
		}
	}

	var rsmngr = resources.GLOBALRSNG()

	if mqttmngr == nil {
		if mqttevent != nil {
			mqttmngr = mqttevent.MqttManager()
		} else if mqttmsg != nil {
			mqttmngr = mqttmsg.Manager()
		}
	}
	var ssn = &Session{atv: active.NewActive(), rsmngr: rsmngr, ssnrsmngr: resources.NewResourcingManager()}
	ssn.atvdbms = database.GLOBALDBMS().ActiveDBMS(ssn.atv)
	ssn.chnghndlr = caching.GLOBALMAPHANDLER().ActiveHandler(ssn.atv)
	ssn.mqttevent = mqttevent
	ssn.mqttmsg = mqttmsg
	ssn.mqttmngr = mqttmngr
	ssn.lstnr = lstnr
	if schdl != nil && schdlngmngr == nil {
		schdlngmngr = schdl.Schedules()
	}
	if schdlngmngr != nil && atvschdlngmngr == nil {
		atvschdlngmngr = schdlngmngr.ActiveSCHEDULING(ssn.atv)
	}
	ssn.schdlngmngr = atvschdlngmngr
	ssn.schdl = schdl
	session = ssn
	return
}

type Session struct {
	lstnr       listen.ListenerAPI
	mqttmsg     mqtt.Message
	mqttevent   mqtt.MqttEvent
	mqttmngr    mqtt.MQTTManagerAPI
	schdlngmngr api.ActiveSchedulesAPI
	schdl       api.ScheduleAPI
	atv         *active.Active
	atvdbms     *database.ActiveDBMS
	rqst        requesting.RequestAPI
	rsmngr      *resources.ResourcingManager
	ssnrsmngr   *resources.ResourcingManager
	chnghndlr   *caching.ActiveHandler
	addNextPath func(nxtpth ...string)
	pathfunc    func() *exepath
	cmnds       map[int]*osprc.Command
}

func (ssn *Session) closecmd(prcid int) {
	if ssn != nil && ssn.cmnds != nil {
		if cmdf := ssn.cmnds[prcid]; cmdf != nil {
			cmdf.OnClose = nil
			delete(ssn.cmnds, prcid)
		}
	}
}

func (ssn *Session) Command(execpath string, execargs ...string) (cmd *osprc.Command, err error) {
	if ssn != nil {
		cmd, err = osprc.NewCommand(execpath, execargs...)
		if err == nil && cmd != nil {
			cmd.OnClose = ssn.closecmd
			if ssn.cmnds == nil {
				ssn.cmnds = make(map[int]*osprc.Command)
			}
			ssn.cmnds[cmd.PrcID()] = cmd
		}
	}
	return
}

func (ssn *Session) Active(a ...interface{}) (atv *active.Active) {
	if ssn != nil {
		if ssn.atv == nil {
			ssn.atv = active.NewActive()
		}
		atv = ssn.atv
	}
	return
}

func (ssn *Session) Scheduling() (schdlsmngr api.ActiveSchedulesAPI) {
	if ssn != nil && ssn.schdlngmngr != nil {
		schdlsmngr = ssn.schdlngmngr
	}
	return
}

func (ssn *Session) Schedule() (schdl api.ScheduleAPI) {
	if ssn != nil && ssn.schdl != nil {
		schdl = ssn.schdl
	}
	return
}

func (ssn *Session) MQTTManager() (mqttmngr mqtt.MQTTManagerAPI) {
	if ssn != nil {
		mqttmngr = ssn.mqttmngr
	}
	return
}

func (ssn *Session) MQTTEvent() (mqttevent mqtt.MqttEvent) {
	if ssn != nil {
		mqttevent = ssn.mqttevent
	}
	return
}

func (ssn *Session) MQTTMessage() (mqttmsg mqtt.Message) {
	if ssn != nil {
		mqttmsg = ssn.mqttmsg
	}
	return
}

func (ssn *Session) Path() (expth api.PathAPI) {
	if ssn != nil && ssn.pathfunc != nil {
		expth = ssn.pathfunc()
	}
	return
}

func (ssn *Session) Env() (env env.EnvAPI) {
	if ssn != nil {
		env = glblenv
	}
	return
}

func (ssn *Session) Listen(network string, addr ...string) (err error) {
	if ssn != nil && ssn.lstnr != nil {
		err = ssn.lstnr.Listen(network, addr...)
	}
	return
}

func (ssn *Session) Send(rqstpath string, a ...interface{}) (rdr iorw.Reader, err error) {
	if ssn != nil {
		rdr, err = internSend(ssn.atv, ssn.FS(), ssn.SessionFS(), ssn.FSUTILS(), ssn.rqst, rqstpath, false, a...)
	}
	return
}

func (ssn *Session) SendRecieve(rqstpath string, a ...interface{}) (rdrwtr iorw.PrinterReader, err error) {
	if ssn != nil {
		rdrwtr, err = internSendreceive(ssn.atv, ssn.FS(), ssn.SessionFS(), ssn.FSUTILS(), ssn.rqst, rqstpath, a...)
	}
	return
}

func (ssn *Session) SessionSend(rqstpath string, a ...interface{}) (rdr iorw.Reader, err error) {
	if ssn != nil {
		rdr, err = internSend(ssn.atv, ssn.FS(), ssn.SessionFS(), ssn.FSUTILS(), ssn.rqst, rqstpath, false, a...)
	}
	return
}

func (ssn *Session) SessionSendRecieve(rqstpath string, a ...interface{}) (rdrwtr iorw.PrinterReader, err error) {
	if ssn != nil {
		rdrwtr, err = internSendreceive(ssn.atv, ssn.FS(), ssn.SessionFS(), ssn.FSUTILS(), ssn.rqst, rqstpath, a...)
	}
	return
}

func (ssn *Session) Caching() (ccngapi caching.MapAPI) {
	if ssn != nil {
		ccngapi = ssn.chnghndlr
	}
	return
}

func internSend(atv *active.Active, fs *fsutils.FSUtils, fssession *fsutils.FSUtils, fsutls fsutils.FSUtils, rqst requesting.RequestAPI, rqstpath string, andeval bool, a ...interface{}) (rdr iorw.Reader, err error) {
	convertactive := false
	if convertactive = strings.Contains(rqstpath, "/active:"); convertactive {
		rqstpath = strings.Replace(rqstpath, "/active:", "/", 1)
	}
	if strings.HasPrefix(rqstpath, "http://") || strings.HasPrefix(rqstpath, "https://") {
		a = append([]interface{}{atv}, a...)
		rdr, err = web.DefaultClient.Send(rqstpath, a...)
	} else {
		rdr = iorw.NewEOFCloseSeekReader(func() (r io.Reader) {
			if r = fs.CAT(rqstpath); r == nil {
				if r = fsutls.CAT(rqstpath); r == nil {
					r = fssession.CAT(rqstpath)
				}
			}
			return r
		}())
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

func internSendreceive(atv *active.Active, fs *fsutils.FSUtils, fssession *fsutils.FSUtils, fsutls fsutils.FSUtils, rqst requesting.RequestAPI, rqstpath string, a ...interface{}) (rdrwtr iorw.PrinterReader, err error) {
	if strings.HasPrefix(rqstpath, "ws://") || strings.HasPrefix(rqstpath, "wss://") {
		a = append([]interface{}{atv}, a...)
		rdrwtr, err = web.DefaultClient.SendReceive(rqstpath, a...)
	}
	return
}

func (ssn *Session) DBMS() (atvdbms database.DBMSAPI) {
	if ssn != nil {
		atvdbms = ssn.atvdbms
	}
	return
}

func (ssn *Session) Close() (err error) {
	if ssn != nil {
		if ssn.atv != nil {
			ssn.atv.Close()
			ssn.atv = nil
		}
		if ssn.rqst != nil {
			ssn.rqst = nil
		}
		if ssn.atvdbms != nil {
			ssn.atvdbms.Dispose()
			ssn.atvdbms = nil
		}
		if ssn.rqst != nil {
			ssn.rqst = nil
		}
		if ssn.lstnr != nil {
			ssn.lstnr = nil
		}
		if ssn.rsmngr != nil {
			ssn.rsmngr = nil
		}
		if ssn.ssnrsmngr != nil {
			ssn.ssnrsmngr.Close()
			ssn.ssnrsmngr = nil
		}
		if ssn.mqttevent != nil {
			ssn.mqttevent = nil
		}
		if ssn.mqttmsg != nil {
			ssn.mqttmsg = nil
		}
		if ssn.mqttmngr != nil {
			ssn.mqttmngr = nil
		}
		if ssn.schdlngmngr != nil {
			ssn.schdlngmngr.Dispose()
			ssn.schdlngmngr = nil
		}
		if ssn.schdl != nil {
			ssn.schdl = nil
		}
		if ssn.cmnds != nil {
			func() {
				if len(ssn.cmnds) > 0 {
					prcsids := make([]int, len(ssn.cmnds))
					prcsidsi := 0
					for prcid := range ssn.cmnds {
						prcsids[prcsidsi] = prcid
						prcsidsi++
					}

					for _, prcid := range prcsids {
						ssn.cmnds[prcid].Close()
					}
					ssn.cmnds = nil
				}
			}()
		}
		ssn = nil
	}
	return
}

func (ssn *Session) In() (rqst requesting.RequestAPI) {
	if ssn != nil && ssn.rqst != nil {
		rqst = ssn.rqst
	}
	return
}

func (ssn *Session) Out() (rspns requesting.ResponseAPI) {
	if ssn != nil && ssn.rqst != nil {
		rspns = ssn.rqst.Response()
	}
	return
}

func (ssn *Session) FS() (fs *fsutils.FSUtils) {
	if ssn != nil && ssn.rsmngr != nil {
		fs = ssn.rsmngr.FS()
	}
	return
}

func (ssn *Session) SessionFS() (fs *fsutils.FSUtils) {
	if ssn != nil && ssn.ssnrsmngr != nil {
		fs = ssn.ssnrsmngr.FS()
	}
	return
}

func (ssn *Session) FSUTILS() (fs fsutils.FSUtils) {
	if ssn != nil {
		fs = fslcl
	}
	return
}

func (ssn *Session) Parameters() (prms parameters.ParametersAPI) {
	if ssn != nil && ssn.rqst != nil {
		prms = ssn.rqst.Parameters()
	}
	return
}

func (ssn *Session) AddPath(nxtpth ...string) {
	if ssn != nil && ssn.addNextPath != nil {
		ssn.addNextPath(nxtpth...)
	}
}

func (ssn *Session) Execute(a ...interface{}) (err error) {
	if ssn != nil {
		var ai = 0
		var nxtrqst requesting.RequestAPI = nil
		if len(a) > 0 {
			for _, d := range a {
				if d != nil {
					if rqstd, _ := d.(requesting.RequestAPI); rqstd != nil {
						if nxtrqst == nil {
							nxtrqst = rqstd
						}
						continue
					}
				}
				ai++
			}
		}
		var prvrqst requesting.RequestAPI = ssn.rqst
		defer func() {
			if nxtrqst != nil {
				ssn.rqst = prvrqst
				nxtrqst = nil
			}
		}()
		if nxtrqst != nil {
			ssn.rqst = nxtrqst
		}
		if rqst := ssn.In(); rqst != nil {
			rspns := rqst.Response()
			prtclrangetype := rqst.RangeType()
			prtclrangeoffset := rqst.RangeOffset()

			var rqstdpaths *enumeration.List = enumeration.NewList()
			ssn.addNextPath = func(nxtpth ...string) {
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
			defer func() { ssn.addNextPath = nil }()

			if ppath := rqst.Path(); ppath != "" {
				ssn.addNextPath(ppath)
				func() {
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
						ssn.pathfunc = func() *exepath {
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

						defer func() {
							ssn.pathfunc = nil
						}()
						var path = expth.Path()
						var pathext = filepath.Ext(path)
						var convertactive bool = false
						var israw = false
						var mimetype, isactive, ismedia = mimes.FindMimeType(path, "text/plain")
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
						if rs = ssn.FS().CAT(path); rs == nil && (strings.LastIndex(path, ".") == -1 || strings.LastIndex(path, "/") > strings.LastIndex(path, ".")) {
							if !strings.HasSuffix(path, "/") {
								path += "/"
							}
							for _, pth := range strings.Split("html,xml,svg,js,json,css", ",") {
								if rs = ssn.FS().CAT(path + "index" + "." + pth); rs == nil {
									if rs = ssn.FS().CAT(path + "main" + "." + pth); rs == nil {
										continue
									} else {
										pathext = "." + pth
										mimetype, isactive, ismedia = mimes.FindMimeType(path+"main"+"."+pth, "text/plain")
										if rspns != nil {
											rspns.SetHeader("Content-Type", mimetype)
										}
										break
									}
								} else {
									pathext = "." + pth
									mimetype, isactive, ismedia = mimes.FindMimeType(path+"index"+"."+pth, "text/plain")
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
								if ssn.atv.LookupTemplate == nil {
									ssn.atv.LookupTemplate = func(lkppath string, a ...interface{}) (lkpr io.Reader, lkperr error) {
										if lkppath != "" && strings.LastIndex(lkppath, ".") == -1 {
											if crntexpths != nil && crntexpths.Length() > 0 {
												if val := crntexpths.Tail().Value(); val != nil {
													if dngexpth, _ := val.(*exepath); dngexpth != nil {
														if dngext := dngexpth.Ext(); dngext != "" {
															lkppath = lkppath + dngext
														} else {
															lkppath = lkppath + pathext
														}

													} else {
														if dngext := expth.Ext(); dngext != "" {
															lkppath = lkppath + dngext
														} else {
															lkppath = lkppath + pathext
														}
													}
												} else {
													if dngext := expth.Ext(); dngext != "" {
														lkppath = lkppath + dngext
													} else {
														lkppath = lkppath + pathext
													}
												}
											} else {
												if dngext := expth.Ext(); dngext != "" {
													lkppath = lkppath + dngext
												} else {
													lkppath = lkppath + pathext
												}
											}
										}
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
											if lkpr = ssn.FS().CAT(lkppath); lkpr == nil {
												lkpr = ssn.FSUTILS().CAT(lkppath)
											}
										}
										return
									}
								}
								if ssn.atv.ObjectMapRef == nil {
									ssn.atv.ObjectMapRef = func() (objrf map[string]interface{}) {
										var objref = map[string]interface{}{}
										objref["kwe"] = ssn
										objrf = objref
										return
									}
								}
								func() {
									var evalerr error = nil
									evalerr = ssn.atv.Eval(rspns, rqst, path, convertactive, rs)
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
								if ismedia {
									if eofrs, _ := rs.(*iorw.EOFCloseSeekReader); eofrs != nil {
										eofrs.Seek(prtclrangeoffset, 0)
										if rssize := eofrs.Size(); rssize > 0 {
											if prtclrangetype == "bytes" && prtclrangeoffset > -1 {
												maxoffset := int64(0)
												maxlen := int64(0)
												if maxoffset = prtclrangeoffset + (rssize - prtclrangeoffset); maxoffset > 0 {
													maxlen = maxoffset - prtclrangeoffset
													maxoffset--
												}

												if maxoffset < prtclrangeoffset {
													maxoffset = prtclrangeoffset
													maxlen = 0
												}

												if maxlen > 1024*1024 {
													maxlen = 1024 * 1024
													maxoffset = prtclrangeoffset + (maxlen - 1)
												}
												contentrange := fmt.Sprintf("%s %d-%d/%d", rqst.RangeType(), prtclrangeoffset, maxoffset, rssize)
												rspns.SetHeader("Content-Range", contentrange)
												rspns.SetHeader("Content-Length", fmt.Sprintf("%d", maxlen))
												eofrs.MaxRead = maxlen
											} else {
												rspns.SetHeader("Content-Length", fmt.Sprintf("%d", rssize))
												eofrs.MaxRead = rssize
											}
										}
										rspns.Print(rs)
									} else {
										rspns.Print(rs)
									}
									prtclrangeoffset = -1
								} else {
									rspns.Print(rs)
								}
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

		}
	}
	return
}

var fslcl fsutils.FSUtils

var glblenv = env.Env()

func init() {
	fslcl = fsutils.NewFSUtils()
	active.LoadGlobalModule("kwesession.js", `function kwefields(obj){
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
	function kwemethods(obj){
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
	}`)
	active.LoadGlobalModule("kwe.js", sysjsTemplate("kwe",
		map[string]interface{}{
			"listen":      "_listen",
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
			"sendreceive": "_sendreceive"}))
	var rsmngr = resources.GLOBALRSNG()
	rsmngr.FS().MKDIR("/require/js", "")
	rsmngr.FS().SET("/require/js/require.js", requirejs.RequireJS())
}

type exepath struct {
	path string
	args []interface{}
}

func (expth *exepath) Ext() (ext string) {
	if expth != nil {
		if pth := expth.Path(); pth != "" {
			ext = filepath.Ext(pth)
		}
	}
	return
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

func ExecuteSession(ssn api.SessionAPI, a ...interface{}) (err error) {
	var closessn = ssn == nil
	if closessn {
		ssn = InvokeSession(a...)
	}
	defer func() {
		if ssn != nil {
			if closessn {
				ssn.Close()
			}
			ssn = nil
		}
	}()
	err = ssn.Execute(a...)
	return
}

func InvokeSession(a ...interface{}) (ssn api.SessionAPI) {
	ssn = NewSession(a...)
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

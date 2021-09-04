package main

import (
	"context"
	"io"
	"os"
	"strings"

	_ "github.com/evocert/kwe/alertify"
	_ "github.com/evocert/kwe/bootstrap"
	"github.com/evocert/kwe/caching"
	"github.com/evocert/kwe/database"
	_ "github.com/evocert/kwe/datepicker"
	_ "github.com/evocert/kwe/fonts/material"
	_ "github.com/evocert/kwe/fonts/robotov27latin"
	"github.com/evocert/kwe/fsutils"
	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/iorw/active"
	"github.com/evocert/kwe/listen"
	"github.com/evocert/kwe/mimes"
	"github.com/evocert/kwe/mqtt"
	scheduling "github.com/evocert/kwe/scheduling/ext"
	"github.com/evocert/kwe/service"
	_ "github.com/evocert/kwe/typescript"
	"github.com/evocert/kwe/web"

	_ "github.com/evocert/kwe/database/mysql"
	_ "github.com/evocert/kwe/database/ora"
	_ "github.com/evocert/kwe/database/postgres"
	_ "github.com/evocert/kwe/database/sqlserver"
	"github.com/evocert/kwe/requesting"
	"github.com/evocert/kwe/resources"
	_ "github.com/evocert/kwe/webactions"
)

func main() {
	lstnr := listen.NewListener()
	var glblutilsfs = fsutils.NewFSUtils()
	var glbldbms = database.GLOBALDBMS
	var glblrsfs = resources.GLOBALRSNG().FS
	var glblchng = caching.GLOBALMAPHANDLER
	var glblschdlng = scheduling.GLOBALSCHEDULES
	active.LoadGlobalModule("kwe.js", sysjsTemplate("kwe",
		map[string]interface{}{
			"in":          "_in",
			"dbms":        "_dbms",
			"caching":     "_caching",
			"out":         "_out",
			"fs":          "_fs",
			"fsutils":     "_fsutils",
			"scheduling":  "_scheduling",
			"mqtting":     "_mqtting",
			"mqttmsg":     "_mqttmsg",
			"mqttevent":   "_mqttevent",
			"extMimetype": "_extMimetype",
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

	lstnr.ServeRequest = func(rqst requesting.RequestAPI, a ...interface{}) (err error) {
		var mqttmsg mqtt.Message = nil
		var mqttevent mqtt.MqttEvent = nil
		if len(a) > 0 {
			for _, d := range a {
				if mqttmsgd, _ := d.(mqtt.Message); mqttmsgd != nil && mqttmsg == nil {
					mqttmsg = mqttmsgd
				} else if mqtteventd, _ := d.(mqtt.MqttEvent); mqtteventd != nil && mqttevent == nil {
					mqttevent = mqtteventd
				}
			}
		}
		rspns := rqst.Response()
		if path := rqst.Path(); path != "" {
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
			if rs != nil {
				if isactive {
					var atv = active.NewActive()
					if atv.ObjectMapRef == nil {
						atv.ObjectMapRef = func() (objref map[string]interface{}) {
							objref = map[string]interface{}{}
							objref["_in"] = rqst
							objref["_dbms"] = glbldbms().ActiveDBMS(atv, rqst.Parameters())
							objref["_caching"] = glblchng().ActiveHandler(atv, rqst.Parameters())
							objref["_out"] = rspns
							objref["_fs"] = glblrsfs()
							objref["_fsutils"] = glblutilsfs
							objref["_scheduling"] = glblschdlng().ActiveSCHEDULING(atv, rqst.Parameters())
							objref["_mqtting"] = glblmqttng
							objref["_mqttmsg"] = mqttmsg
							objref["_mqttevent"] = mqttevent
							objref["_extMimetype"] = mimes.ExtMimeType
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
							return
						}
					}
					func() {
						defer atv.Close()
						var evalerr error = nil
						if convertactive {
							evalerr = atv.Eval(rspns, rqst, path, "<@", "\r\n", rs, "@>")
						} else {
							evalerr = atv.Eval(rspns, rqst, path, rs)
						}
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
		} else {

		}
		return
	}

	service.ServeRequest = lstnr.ServeRequest
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
			if convertactive {
				pwerr = atv.Eval(pw, nil, rqstpath, "<@", "\r\n", rdr, "@>")
			} else {
				pwerr = atv.Eval(pw, nil, rqstpath, rdr)
			}
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
						cde += `obj` + nmspace + `["` + objk + `"]=` + objs + `;`
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

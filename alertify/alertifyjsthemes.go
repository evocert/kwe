package alertify

import (
	"io"
	"strings"

	"github.com/evocert/kwe/resources"
)

const alertifythemesdefaultcss string = `/**
* alertifyjs 1.13.1 http://alertifyjs.com
* AlertifyJS is a javascript framework for developing pretty browser dialogs and notifications.
* Copyright 2019 Mohammad Younes <Mohammad@alertifyjs.com> (http://alertifyjs.com) 
* Licensed under GPL 3 <https://opensource.org/licenses/gpl-3.0>*/
.alertify .ajs-dialog{background-color:#fff;-webkit-box-shadow:0 15px 20px 0 rgba(0,0,0,.25);box-shadow:0 15px 20px 0 rgba(0,0,0,.25);border-radius:2px}.alertify .ajs-header{color:#000;font-weight:700;background:#fafafa;border-bottom:#eee 1px solid;border-radius:2px 2px 0 0}.alertify .ajs-body{color:#000}.alertify .ajs-body .ajs-content .ajs-input{display:block;width:100%;padding:8px;margin:4px;border-radius:2px;border:1px solid #ccc}.alertify .ajs-body .ajs-content p{margin:0}.alertify .ajs-footer{background:#fbfbfb;border-top:#eee 1px solid;border-radius:0 0 2px 2px}.alertify .ajs-footer .ajs-buttons .ajs-button{background-color:transparent;color:#000;border:0;font-size:14px;font-weight:700;text-transform:uppercase}.alertify .ajs-footer .ajs-buttons .ajs-button.ajs-ok{color:#3593d2}.alertify-notifier .ajs-message{background:rgba(255,255,255,.95);color:#000;text-align:center;border:solid 1px #ddd;border-radius:2px}.alertify-notifier .ajs-message.ajs-success{color:#fff;background:rgba(91,189,114,.95);text-shadow:-1px -1px 0 rgba(0,0,0,.5)}.alertify-notifier .ajs-message.ajs-error{color:#fff;background:rgba(217,92,92,.95);text-shadow:-1px -1px 0 rgba(0,0,0,.5)}.alertify-notifier .ajs-message.ajs-warning{background:rgba(252,248,215,.95);border-color:#999}`

func AlertifyThemesDefaultCSS() io.Reader {
	return strings.NewReader(alertifythemesdefaultcss)
}

func init() {
	gblrs := resources.GLOBALRSNG()
	gblrs.FS().MKDIR("/alertify/css/themes", "")
	gblrs.FS().SET("/alertify/css/themes/default.css", AlertifyThemesDefaultCSS())
}

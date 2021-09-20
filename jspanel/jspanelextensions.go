package jspanel

import (
	"io"
	"strings"
)

func JSPanelModalJS() io.Reader {
	return strings.NewReader(jspanelmodaljs)
}

const jspanelmodaljs string = `"use strict";jsPanel.modal||(jsPanel.modal={version:"1.2.5",date:"2020-04-26 23:23",defaults:{closeOnEscape:!0,closeOnBackdrop:!0,dragit:!1,headerControls:"closeonly",resizeit:!1,syncMargins:!1},addBackdrop:function(e){var n=document.getElementsByClassName("jsPanel-modal-backdrop").length,a=document.createElement("div");return a.id="jsPanel-modal-backdrop-"+e,0===n?a.className="jsPanel-modal-backdrop":n>0&&(a.className="jsPanel-modal-backdrop jsPanel-modal-backdrop-multi"),a.style.zIndex=this.ziModal.next(),a},removeBackdrop:function(e){var n=document.getElementById("jsPanel-modal-backdrop-".concat(e));n.classList.add("jsPanel-modal-backdrop-out");var a=1e3*parseFloat(getComputedStyle(n).animationDuration);window.setTimeout(function(){document.body.removeChild(n)},a)},create:function(){var e=arguments.length>0&&void 0!==arguments[0]?arguments[0]:{};e.paneltype="modal",e.id?"function"==typeof e.id&&(e.id=e.id()):e.id="jsPanel-".concat(jsPanel.idCounter+=1);var n=e,a=this.addBackdrop(n.id);return e.config&&delete(n=Object.assign({},e.config,e)).config,n=Object.assign({},this.defaults,n,{container:"window"}),document.body.append(a),jsPanel.create(n,function(e){e.style.zIndex=jsPanel.modal.ziModal.next(),e.header.style.cursor="default",e.footer.style.cursor="default",n.closeOnBackdrop&&jsPanel.pointerup.forEach(function(a){document.getElementById("jsPanel-modal-backdrop-".concat(n.id)).addEventListener(a,function(){e.close(null,!0)})}),e.options.onclosed.unshift(function(){return jsPanel.modal.removeBackdrop(n.id),!0})})}},jsPanel.modal.ziModal=function(){var e=jsPanel.ziBase+1e4;return{next:function(){return e++}}}()),"undefined"!=typeof module&&(module.exports=jsPanel);`

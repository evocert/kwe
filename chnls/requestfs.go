package chnls

import (
	"io"
	"strings"

	"github.com/evocert/kwe/fsutils"
	"github.com/evocert/kwe/resources"
)

//FS return fsutils.FSUtils implementation for *Request
func (rqst *Request) FS() *fsutils.FSUtils {
	if rqst.fsutils == nil {
		rqst.fsutils = &fsutils.FSUtils{
			FIND: func(path ...string) (finfos []fsutils.FileInfo) {
				finfos = rqst.fsfind(rqst.rscngmnger(), path...)
				return
			}, LS: func(path ...string) (finfos []fsutils.FileInfo) {
				finfos = rqst.fsls(rqst.rscngmnger(), path...)
				return
			}, MKDIR: func(path ...interface{}) bool {
				return rqst.fsmkdir(rqst.rscngmnger(), path...)
			}, MKDIRALL: func(path ...interface{}) bool {
				return rqst.fsmkdirall(rqst.rscngmnger(), path...)
			}, RM: func(path string) bool {
				return rqst.fsrm(rqst.rscngmnger(), path)
			}, MV: func(path string, destpath string) bool {
				return rqst.fsmv(rqst.rscngmnger(), path, destpath)
			}, TOUCH: func(path string) bool {
				return rqst.fstouch(rqst.rqstrsngmngr, path)
			}, PIPE: func(path string) io.Reader {
				return rqst.fspipe(rqst.rscngmnger(), path)
			}, PIPES: func(path string) string {
				return rqst.fspipes(rqst.rscngmnger(), path)
			}, CAT: func(path string) io.Reader {
				return rqst.fscat(rqst.rscngmnger(), path)
			}, CATS: func(path string) string {
				return rqst.fscats(rqst.rscngmnger(), path)
			}, SET: func(path string, a ...interface{}) bool {
				return rqst.fsset(rqst.rscngmnger(), path, a...)
			}, APPEND: func(path string, a ...interface{}) bool {
				return rqst.fsappend(rqst.rscngmnger(), path, a...)
			},
		}
	}
	return rqst.fsutils
}

func (rqst *Request) rscngmnger() *resources.ResourcingManager {
	if rqst != nil {
		if rqst.rqstrsngmngr == nil {
			rqst.rqstrsngmngr = resources.NewResourcingManager()
		}
		return rqst.rqstrsngmngr
	}
	return nil
}

func (rqst *Request) fsfind(rsngmngr *resources.ResourcingManager, path ...string) (finfos []fsutils.FileInfo) {
	finfos = rsngmngr.FS().FIND(path...)
	return
}

func (rqst *Request) fsls(rsngmngr *resources.ResourcingManager, path ...string) (finfos []fsutils.FileInfo) {
	finfos = rsngmngr.FS().LS(path...)
	return
}

func (rsqt *Request) fsmkdir(rsngmngr *resources.ResourcingManager, path ...interface{}) bool {
	if pthl := len(path); pthl > 0 {
		var pth1 = ""
		var pth2 = ""
		if pthl > 1 {
			pth2, _ = path[1].(string)
			pth2 = strings.TrimSpace(pth2)
			path[1] = pth2
		}
		if pth1, _ = path[0].(string); pth1 != "" && !strings.HasPrefix(pth1, "/") {
			pth1 = "/" + pth1
		}
		path[0] = pth1
		if !rsngmngr.FS().MKDIR(path...) {
			if pthl == 1 && pth1 != "" {
				rsngmngr.RegisterEndpoint(pth1, "")
				return true
			} else if pthl == 2 {
				rsngmngr.RegisterEndpoint(pth1, pth2)
				return true
			} else if pthl > 2 {
				rsngmngr.RegisterEndpoint(pth1, pth2, path[2:]...)
				return true
			}
		}
	}
	return false
}

func (rsqt *Request) fsmkdirall(rsngmngr *resources.ResourcingManager, path ...interface{}) bool {
	if pthl := len(path); pthl > 0 {
		var pth1 = ""
		var pth2 = ""
		if pthl > 1 {
			pth2, _ = path[1].(string)
			pth2 = strings.TrimSpace(pth2)
			path[1] = pth2
		}
		if pth1, _ = path[0].(string); pth1 != "" && !strings.HasPrefix(pth1, "/") {
			path[0] = "/" + pth1
		}
		if !rsngmngr.FS().MKDIRALL(path...) {
			if pthl == 1 && pth1 != "" {
				rsngmngr.RegisterEndpoint(pth1, "")
				return true
			} else if pthl == 2 {
				rsngmngr.RegisterEndpoint(pth1, pth2)
				return true
			} else if pthl > 2 {
				rsngmngr.RegisterEndpoint(pth1, pth2, path[2:]...)
				return true
			}
		}
	}
	return true
}

func (rsqt *Request) fsrm(rsngmngr *resources.ResourcingManager, path string) bool {
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return rsngmngr.FS().RM(path)
}

func (rsqt *Request) fstouch(rsngmngr *resources.ResourcingManager, path string) bool {
	return rsngmngr.FS().TOUCH(path)
}

func (rsqt *Request) fsset(rsngmngr *resources.ResourcingManager, path string, a ...interface{}) bool {
	return rsngmngr.FS().SET(path, a...)
}

func (rsqt *Request) fsappend(rsngmngr *resources.ResourcingManager, path string, a ...interface{}) bool {
	return rsngmngr.FS().APPEND(path, a...)
}

func (rsqt *Request) fsmv(rsngmngr *resources.ResourcingManager, path string, destpath string) bool {
	return rsngmngr.FS().MV(path, destpath)
}

func (rsqt *Request) fscat(rsngmngr *resources.ResourcingManager, path string) io.Reader {
	return rsngmngr.FS().CAT(path)
}

func (rsqt *Request) fscats(rsngmngr *resources.ResourcingManager, path string) string {
	return rsngmngr.FS().CATS(path)
}

func (rsqt *Request) fspipe(rsngmngr *resources.ResourcingManager, path string) io.Reader {
	return rsngmngr.FS().PIPE(path)
}

func (rsqt *Request) fspipes(rsngmngr *resources.ResourcingManager, path string) string {
	return rsngmngr.FS().PIPES(path)
}

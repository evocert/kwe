package fsutils

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/evocert/kwe/iorw"
)

//LS List dir content
func LS(path string, altpath ...string) (finfos []FileInfo, err error) {
	path = strings.Replace(path, "\\", "/", -1)
	var altpth = ""
	if len(altpath) == 1 && altpath[0] != "" {
		altpth = altpath[0]
	}
	if fi, fierr := os.Stat(path); fierr == nil {
		if fi.IsDir() {
			if !strings.HasSuffix(path, "/") {
				path += "/"
			}
			if fifis, fifpath, fifaltpath, fifiserr := internalFind(fi, path, altpth); fifiserr == nil {
				if !strings.HasSuffix(fifpath, "/") {
					fifpath += "/"
				}
				if !strings.HasSuffix(fifaltpath, "/") {
					fifaltpath += "/"
				}
				if fifaltpath != "" {
					fifpath = fifaltpath
				}
				for _, fifi := range fifis {
					if finfos == nil {
						finfos = []FileInfo{}
					}
					finfos = append(finfos, newFileInfo(fifi.Name(), path+fifi.Name(), fifpath+fifi.Name(), fifi.Size(), fifi.Mode(), fifi.ModTime()))
				}
			} else {

			}
		} else {
			finfos = []FileInfo{newFileInfo(fi.Name(), path+fi.Name(), path+fi.Name(), fi.Size(), fi.Mode(), fi.ModTime())}
		}
	} else {
		err = fierr
	}
	return
}

func internalFind(fi os.FileInfo, rootpath string, altrootpath string) (finfos []os.FileInfo, fipath string, fialtpath string, err error) {
	if strings.HasSuffix(rootpath, fi.Name()) {
		rootpath = rootpath[:len(rootpath)-len(fi.Name())]
	}
	rootpath = strings.Replace(rootpath, "\\", "/", -1)
	if !strings.HasSuffix(rootpath, "/") {
		rootpath += "/"
	}

	altrootpath = strings.Replace(altrootpath, "\\", "/", -1)
	if !strings.HasSuffix(altrootpath, "/") {
		altrootpath += "/"
	}
	if fi.IsDir() {
		if f, ferr := os.Open(rootpath + fi.Name()); ferr == nil {
			if fis, fiserr := f.Readdir(0); fiserr == nil && len(fis) > 0 {
				finfos = fis[:]
			}
			rootpath = rootpath + fi.Name()
			f.Close()
		}
	} else {
		finfos = []os.FileInfo{fi}
	}
	fipath = rootpath
	if altrootpath != "" {
		fialtpath = altrootpath
	}
	return
}

// A FileInfo describes a file
type FileInfo interface {
	Name() string         // base name of the file
	Path() string         // relative path of the file
	AbsolutePath() string // absolute path of the file
	Size() int64          // length in bytes for regular files; system-dependent for others
	Mode() os.FileMode    // file mode bits
	ModTime() time.Time   // modification time
	IsDir() bool          // abbreviation for Mode().IsDir()
	JSON() string         //json representation as a string
}

type fileInfo struct {
	name         string
	path         string
	absolutepath string
	size         int64
	mode         os.FileMode
	modtime      time.Time
}

func newFileInfo(name string,
	path string,
	absolutepath string,
	size int64,
	mode os.FileMode,
	modtime time.Time) (finfo *fileInfo) {
	finfo = &fileInfo{name: name, path: path, absolutepath: absolutepath, size: size, mode: mode, modtime: modtime}
	return
}

func (finfo *fileInfo) Name() string {
	return finfo.name
}

func (finfo *fileInfo) Path() string {
	return finfo.path
}

func (finfo *fileInfo) AbsolutePath() string {
	return finfo.absolutepath
}

func (finfo *fileInfo) Size() int64 {
	return finfo.size
}

func (finfo *fileInfo) Mode() os.FileMode {
	return finfo.mode
}

func (finfo *fileInfo) ModTime() time.Time {
	return finfo.modtime
}

func (finfo *fileInfo) IsDir() bool {
	return finfo != nil && finfo.mode.IsDir()
}

func (finfo *fileInfo) JSON() (s string) {
	buf := iorw.NewBuffer()
	enc := json.NewEncoder(buf)
	enc.Encode(map[string]interface{}{"Name": finfo.name, "Path": finfo.path, "Dir": finfo.IsDir(), "Modified": finfo.modtime, "Size": finfo.size})
	s = buf.String()
	buf.Close()
	buf = nil
	if s != "" {
		s = strings.TrimSpace(s)
	}
	return
}

//FIND list recursive dir content
func FIND(path string, altpath ...string) (finfos []FileInfo, err error) {
	var nxtfisfunc func(fi os.FileInfo, fipath string, fialtpath string) = nil
	var altpth = ""
	if len(altpath) == 1 && altpath[0] != "" {
		altpth = strings.Replace(altpath[0], "\\", "/", -1)
		if altpth != "/" && !strings.HasSuffix(altpth, "/") {
			altpth += "/"
		}
	}

	fisfunc := func(fi os.FileInfo, fipath string, fialtpath string) {
		if finfos == nil {
			finfos = []FileInfo{}
		}
		if strings.HasSuffix(fipath, fi.Name()) {
			fipath = fipath[:len(fipath)-len(fi.Name())]
		}
		fipath = strings.Replace(fipath, "\\", "/", -1)
		if !strings.HasSuffix(fipath, "/") {
			fipath += "/"
		}
		finfos = append(
			finfos,
			newFileInfo(fi.Name(), fipath+fi.Name(), fipath+fi.Name(), fi.Size(), fi.Mode(), fi.ModTime()),
		)
		if fifis, fifpath, fifaltpath, fifiserr := internalFind(fi, fipath, fialtpath); fifiserr == nil {
			if !strings.HasSuffix(fifpath, "/") {
				fifpath += "/"
			}
			if !strings.HasSuffix(fifaltpath, "/") {
				fifaltpath += "/"
			}
			for _, fifi := range fifis {
				if finfos == nil {
					finfos = []FileInfo{}
				}
				if fifi.IsDir() {
					nxtfisfunc(fifi, fifpath+fifi.Name(), fifaltpath+fifi.Name())
				} else {
					if fifaltpath != "" {
						finfos = append(finfos, newFileInfo(fifi.Name(), fifaltpath+fifi.Name(), fifaltpath+fifi.Name(), fifi.Size(), fifi.Mode(), fifi.ModTime()))
					} else {
						finfos = append(finfos, newFileInfo(fifi.Name(), fipath+fifi.Name(), fifpath+fifi.Name(), fifi.Size(), fifi.Mode(), fifi.ModTime()))
					}
				}
			}
		}
	}
	nxtfisfunc = fisfunc
	if fi, fierr := os.Stat(path); fierr == nil {
		fisfunc(fi, path, altpth)
	}
	return
}

//MKDIR make directory
func MKDIR(path string) error {
	return os.Mkdir(path, os.ModeDir)
}

//MKDIRALL make directory with all necessary parents
func MKDIRALL(path string) error {
	return os.MkdirAll(path, os.ModeDir)
}

//RM Remove file or directory
func RM(path string) (err error) {
	err = os.Remove(path)
	return
}

//MV Move file or directory
func MV(path string, destpath string) (err error) {
	err = os.Rename(path, destpath)
	return
}

//TOUCH Create an empty file if the file doesn’t already exist or
// if the file already exists then update the modified time of the file
func TOUCH(path string) (err error) {
	statf, staterr := os.Stat(path)
	if os.IsNotExist(staterr) {
		if file, ferr := os.Create(path); ferr == nil {
			defer file.Close()
		} else {
			err = ferr
		}
	} else if !statf.IsDir() {
		currentTime := time.Now().Local()
		err = os.Chtimes(path, currentTime, currentTime)
	}
	return
}

//FINFOPATHSJSON []FileInfo to JSON array
func FINFOPATHSJSON(a ...FileInfo) (s string) {
	s = "["
	for {
		if al := len(a); al > 0 {
			s += a[0].JSON()
			a = a[1:]
			if al > 0 {
				s += ","
			}
		} else {
			break
		}
	}
	s += "]"
	return
}

//FSUtils struct
type FSUtils struct {
	LS             func(path string) (finfos []FileInfo)   `json:"ls"`
	FIND           func(path string) (finfos []FileInfo)   `json:"find"`
	MKDIR          func(path string) bool                  `json:"mkdir"`
	MKDIRALL       func(path string) bool                  `json:"mkdirall"`
	RM             func(path string) bool                  `json:"rm"`
	MV             func(path string, destpath string) bool `json:"mv"`
	TOUCH          func(path string) bool                  `json:"touch"`
	FINFOPATHSJSON func(a ...FileInfo) (s string)          `json:"finfopathsjson"`
}

//NewFSUtils return instance of FSUtils
func NewFSUtils() (fsutlsstrct FSUtils) {
	fsutlsstrct = FSUtils{
		FIND: func(path string) (finfos []FileInfo) {
			finfos, _ = FIND(path)
			return
		},
		LS: func(path string) (finfos []FileInfo) {
			finfos, _ = LS(path)
			return
		},
		MKDIR: func(path string) bool {
			if err := MKDIR(path); err == nil {
				return true
			}
			return false
		},
		MKDIRALL: func(path string) bool {
			if err := MKDIRALL(path); err == nil {
				return true
			}
			return false
		},
		MV: func(path string, destpath string) bool {
			if err := MV(path, destpath); err == nil {
				return true
			}
			return false
		},
		RM: func(path string) bool {
			if err := RM(path); err == nil {
				return true
			}
			return false
		},
		TOUCH: func(path string) bool {
			if err := TOUCH(path); err == nil {
				return true
			}
			return false
		},
		FINFOPATHSJSON: func(a ...FileInfo) (s string) {
			s = FINFOPATHSJSON(a...)
			return
		}}
	return
}

package fsutils

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/evocert/kwe/iorw"
)

//LS List dir content
func LS(path string, altpath ...string) (finfos []FileInfo, err error) {
	path = strings.Replace(path, "\\", "/", -1)
	var altpth = ""
	if len(altpath) == 1 && altpath[0] != "" {
		altpth = strings.Replace(altpath[0], "\\", "/", -1)
	}

	if fi, fierr := os.Stat(path); fierr == nil {
		if fi.IsDir() {
			if fifis, fifpath, fifaltpath, fifiserr := internalFind(fi, path, altpth); fifiserr == nil {
				if !strings.HasSuffix(fifpath, "/") {
					fifpath += "/"
				}
				if fifaltpath != "" && !strings.HasSuffix(fifaltpath, "/") {
					fifaltpath += "/"
				}

				for fifin := range fifis {
					fifi := fifis[fifin]
					if finfos == nil {
						finfos = []FileInfo{}
					}
					if fifaltpath != "" {
						finfos = append(finfos, newFileInfo(fifi.Name(), fifaltpath+fifi.Name(), fifpath+fifi.Name(), fifi.Size(), fifi.Mode(), fifi.ModTime()))
					} else {
						finfos = append(finfos, newFileInfo(fifi.Name(), fifpath+fifi.Name(), fifpath+fifi.Name(), fifi.Size(), fifi.Mode(), fifi.ModTime()))
					}
				}
			}
		} else {
			fname := fi.Name()
			if strings.HasSuffix(path, fi.Name()) {
				path = path[:len(path)-len(fi.Name())]
			}
			if altpth != "" {
				if !strings.HasSuffix(altpth, fi.Name()) {
					if strings.LastIndex(altpth, ".") > strings.LastIndex(altpth, "/") {
						if strings.LastIndex(altpth, "/") > -1 {
							fname = altpth[strings.LastIndex(altpth, "/")+1:]
						} else {
							fname = altpth
						}
					} else {
						if !strings.HasSuffix(altpth, "/") {
							altpth += "/"
						}
						altpth += fi.Name()
					}
				}
				finfos = []FileInfo{newFileInfo(fname, altpth, path+fi.Name(), fi.Size(), fi.Mode(), fi.ModTime())}
			} else {
				finfos = []FileInfo{newFileInfo(fi.Name(), path+fi.Name(), path+fi.Name(), fi.Size(), fi.Mode(), fi.ModTime())}
			}
		}
	} else {
		var tmppath = ""
		var tmppaths = strings.Split(path, "/")
		for pn := range tmppaths {
			ps := tmppaths[pn]
			if tmpl := len(tmppaths); pn < tmpl {
				if fi, fierr := os.Stat(tmppath + ps + ".zip"); fierr == nil && !fi.IsDir() {
					var testpath = strings.Join(tmppaths[pn+1:tmpl], "/")
					var remainingpath = strings.Join(tmppaths[:pn+1], "/")
					if remainingpath != "" {
						if r, zerr := zip.OpenReader(tmppath + ps + ".zip"); zerr == nil {
							dirsfound := map[string]bool{}
							fname := ""
							for _, f := range r.File {
								if fname = f.Name; strings.HasPrefix(fname, testpath) {
									if fname[len(testpath):] == "" {
										if fname == testpath {
											fname = testpath
										}
									} else {
										fname = fname[len(testpath):]
									}

									if fname != "" {
										if fname[0:0] == "/" {
											fname = fname[1:]
										}
										if fname != "" {
											if fname != testpath {
												if strings.Contains(fname, "/") {
													fname = fname[:strings.Index(fname, "/")]
													if df, dfok := dirsfound[fname]; dfok {
														if !df {
															dirsfound[fname] = true
															if len(testpath) > 0 && !strings.HasSuffix(testpath, "/") {
																fmt.Println(testpath + "/" + fname)
															} else {
																fmt.Println(testpath + fname)
															}
														} else {
															fname = ""
														}
													} else {
														dirsfound[fname] = true
														if len(testpath) > 0 && !strings.HasSuffix(testpath, "/") {
															fname = testpath + "/" + fname

														} else {
															fname = testpath + fname
														}
													}
												} else {
													if len(testpath) > 0 && !strings.HasSuffix(testpath, "/") {
														fname = testpath + "/" + fname
													} else {
														fname = testpath + fname
													}
												}
											}
											if fname != "" {
												fifi := f.FileInfo()
												if finfos == nil {
													finfos = []FileInfo{}
												}
												if altpth != "" {
													if strings.LastIndex(altpth, ".") == -1 && !strings.HasSuffix(altpth, "/") {
														altpth += "/"
													}
													if fname == testpath {
														if !strings.HasSuffix(testpath, "/") {
															finfos = append(finfos, newFileInfo(fifi.Name(), altpth+fifi.Name(), remainingpath+".zip/"+testpath+"/"+fifi.Name(), fifi.Size(), fifi.Mode(), fifi.ModTime()))
														} else {
															finfos = append(finfos, newFileInfo(fifi.Name(), altpth+fifi.Name(), remainingpath+".zip/"+testpath+fifi.Name(), fifi.Size(), fifi.Mode(), fifi.ModTime()))
														}
													} else {
														if !strings.HasSuffix(testpath, "/") {
															finfos = append(finfos, newFileInfo(fifi.Name(), altpth+fifi.Name(), remainingpath+".zip/"+testpath+"/"+fifi.Name(), fifi.Size(), fifi.Mode(), fifi.ModTime()))
														} else {
															finfos = append(finfos, newFileInfo(fifi.Name(), altpth+fifi.Name(), remainingpath+".zip/"+testpath+fifi.Name(), fifi.Size(), fifi.Mode(), fifi.ModTime()))
														}
													}
												} else {
													if fname == testpath {
														finfos = append(finfos, newFileInfo(fifi.Name(), remainingpath+"/"+testpath, remainingpath+".zip/"+testpath, fifi.Size(), fifi.Mode(), fifi.ModTime()))
													} else {
														if !strings.HasSuffix(testpath, "/") {
															finfos = append(finfos, newFileInfo(fifi.Name(), remainingpath+"/"+testpath+"/"+fifi.Name(), remainingpath+".zip/"+testpath+"/"+fifi.Name(), fifi.Size(), fifi.Mode(), fifi.ModTime()))
														} else {
															finfos = append(finfos, newFileInfo(fifi.Name(), remainingpath+"/"+testpath+"/"+fifi.Name(), remainingpath+".zip/"+testpath+fifi.Name(), fifi.Size(), fifi.Mode(), fifi.ModTime()))
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
					break
				} else {
					tmppath = tmppath + ps + "/"
				}
			} else {
				break
			}
		}
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
	if altrootpath != "" && !strings.HasSuffix(altrootpath, "/") {
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
	enc.Encode(map[string]interface{}{"Name": finfo.name, "Path": finfo.path, "Absolute-Path": finfo.absolutepath, "Dir": finfo.IsDir(), "Modified": finfo.modtime, "Size": finfo.size})
	s = buf.String()
	buf.Close()
	buf = nil
	if s != "" {
		s = strings.TrimSpace(s)
	}
	return
}

//ABS return absolute path from relative path
func ABS(path string) (abspath string, err error) {
	if path != "" {
		path = strings.Replace(path, "\\", "/", -1)
	}
	if abspath, err = filepath.Abs(path); abspath != "" {
		abspath = strings.Replace(abspath, "\\", "/", -1)
	}
	return
}

//FIND list recursive dir content
func FIND(path string, altpath ...string) (finfos []FileInfo, err error) {
	var nxtfisfunc func(fi os.FileInfo, fipath string, fialtpath string) = nil
	var altpth = ""
	if len(altpath) == 1 && altpath[0] != "" {
		altpth = strings.Replace(altpath[0], "\\", "/", -1)
	}

	fisfunc := func(fi os.FileInfo, fipath string, fialtpath string) {
		if finfos == nil {
			finfos = []FileInfo{}
		}
		if strings.HasSuffix(fipath, fi.Name()) {
			fipath = fipath[:len(fipath)-len(fi.Name())]
		}
		fipath = strings.Replace(fipath, "\\", "/", -1)
		if fi.IsDir() {
			if !strings.HasSuffix(fipath, "/") {
				fipath += "/"
			}
			if fialtpath != "" {
				if fialtpath != "/" && !strings.HasSuffix(fialtpath, "/") {
					fialtpath += "/"
				}
				finfos = append(
					finfos,
					newFileInfo(fi.Name(), fialtpath, fipath+fi.Name(), fi.Size(), fi.Mode(), fi.ModTime()),
				)
			} else {
				finfos = append(
					finfos,
					newFileInfo(fi.Name(), fipath+fi.Name(), fipath+fi.Name(), fi.Size(), fi.Mode(), fi.ModTime()),
				)
			}
			if fifis, fifpath, fifaltpath, fifiserr := internalFind(fi, fipath, fialtpath); fifiserr == nil {
				if !strings.HasSuffix(fifpath, "/") {
					fifpath += "/"
				}
				if fifaltpath != "" && !strings.HasSuffix(fifaltpath, "/") {
					fifaltpath += "/"
				}
				for fifin := range fifis {
					fifi := fifis[fifin]
					if finfos == nil {
						finfos = []FileInfo{}
					}
					if fifi.IsDir() {
						if fifaltpath != "" {
							nxtfisfunc(fifi, fifpath+fifi.Name(), fifaltpath+fifi.Name())
						} else {
							nxtfisfunc(fifi, fifpath+fifi.Name(), "")
						}
					} else {
						if fifaltpath != "" {
							finfos = append(finfos, newFileInfo(fifi.Name(), fifaltpath+fifi.Name(), fifpath+fifi.Name(), fifi.Size(), fifi.Mode(), fifi.ModTime()))
						} else {
							finfos = append(finfos, newFileInfo(fifi.Name(), fifpath+fifi.Name(), fifpath+fifi.Name(), fifi.Size(), fifi.Mode(), fifi.ModTime()))
						}
					}
				}
			}
		} else {
			fname := fi.Name()
			if strings.HasSuffix(fipath, fi.Name()) {
				fipath = path[:len(fipath)-len(fi.Name())]
			}
			if fialtpath != "" {
				if !strings.HasSuffix(fialtpath, fi.Name()) {
					if strings.LastIndex(fialtpath, ".") > strings.LastIndex(fialtpath, "/") {
						if strings.LastIndex(fialtpath, "/") > -1 {
							fname = altpth[strings.LastIndex(fialtpath, "/")+1:]
						} else {
							fname = fialtpath
						}
					} else {
						if !strings.HasSuffix(fialtpath, "/") {
							fialtpath += "/"
						}
						fialtpath += fi.Name()
					}
				}
				finfos = []FileInfo{newFileInfo(fname, fialtpath, fipath+fi.Name(), fi.Size(), fi.Mode(), fi.ModTime())}
			} else {
				finfos = []FileInfo{newFileInfo(fi.Name(), fipath+fi.Name(), fipath+fi.Name(), fi.Size(), fi.Mode(), fi.ModTime())}
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

//RM Remove file or directory recursive
func RM(path string) (err error) {
	err = os.RemoveAll(path)
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

//CAT return file content if file exists else empty string
func CAT(path string) (r io.Reader, err error) {
	if statf, staterr := os.Stat(path); staterr != nil {
		err = staterr
	} else if !statf.IsDir() {
		if statf.Size() > 0 {
			if f, ferr := os.Open(path); ferr == nil {
				pr, pw := io.Pipe()
				ctx, ctxcancel := context.WithCancel(context.Background())
				go func() {
					var pwerr error = nil
					defer func() {
						f.Close()
						if pwerr == nil {
							pw.Close()
						} else {
							pw.CloseWithError(pwerr)
						}
					}()
					ctxcancel()
					if _, pwerr = io.Copy(pw, f); pwerr != nil {
						if pwerr == io.EOF {
							pwerr = nil
						}
					}
				}()
				<-ctx.Done()
				ctx = nil
				r = iorw.NewEOFCloseSeekReader(pr)
			} else {
				err = ferr
			}
		}
	}
	return
}

//MULTICAT return file(s) content if file(s) exists else empty string
func MULTICAT(path ...string) (r io.Reader, err error) {
	if pl := len(path); pl > 0 {
		var rdrs = []io.Reader{}
		for pthn := range path {
			if statf, staterr := os.Stat(path[pthn]); staterr != nil {
				err = staterr
			} else if !statf.IsDir() {
				if statf.Size() > 0 {
					if f, ferr := os.Open(path[pthn]); ferr == nil {
						pr, pw := io.Pipe()
						ctx, ctxcancel := context.WithCancel(context.Background())
						go func() {
							var pwerr error = nil
							defer func() {
								f.Close()
								if pwerr == nil {
									pw.Close()
								} else {
									pw.CloseWithError(pwerr)
								}
							}()
							ctxcancel()
							if _, pwerr = io.Copy(pw, f); pwerr != nil {
								if pwerr == io.EOF {
									pwerr = nil
								}
							}
						}()
						<-ctx.Done()
						ctx = nil
						rdrs = append(rdrs, pr)
					} else {
						err = ferr
					}
				}
			}
		}
		if len(rdrs) > 0 {
			r = iorw.NewMultiEOFCloseSeekReader(rdrs...)
		}
	}
	return
}

//CATS return file content if file exists else empty string
func CATS(path string) (cntnt string, err error) {
	var r io.Reader = nil
	if r, err = CAT(path); err == nil {
		if r != nil {
			var rc io.Closer = nil
			rc, _ = r.(io.Closer)
			func() {
				defer func() {
					if rc != nil {
						rc.Close()
						rc = nil
					}
					r = nil
				}()
				cntnt, err = iorw.ReaderToString(r)
			}()
		}
	}
	return
}

//MULTICATS return file(s) content if file(s) exists else empty string
func MULTICATS(path ...string) (cntnt string, err error) {
	if len(path) > 0 {
		var s = ""
		for pthn := range path {
			var r io.Reader = nil
			if r, err = CAT(path[pthn]); err == nil {
				if r != nil {
					var rc io.Closer = nil
					rc, _ = r.(io.Closer)
					func() {
						defer func() {
							if rc != nil {
								rc.Close()
								rc = nil
							}
							r = nil
						}()
						if s, err = iorw.ReaderToString(r); s != "" {
							cntnt += s
						}
					}()
					if err != nil {
						break
					}
				}
			}
		}
	}
	return
}

//PIPE return file content if file exists else empty string
func PIPE(path string) (r io.Reader, err error) {
	if statf, staterr := os.Stat(path); staterr != nil {
		err = staterr
	} else if !statf.IsDir() {
		if statf.Size() > 0 {
			if f, ferr := os.Open(path); ferr == nil {
				pr, pw := io.Pipe()
				ctx, ctxcancel := context.WithCancel(context.Background())
				go func() {
					var pwerr error = nil
					defer func() {
						f.Close()
						if pwerr == nil {
							pw.Close()
						} else {
							pw.CloseWithError(pwerr)
						}
					}()
					ctxcancel()
					if _, pwerr = io.Copy(pw, f); pwerr != nil {
						if pwerr == io.EOF {
							pwerr = nil
						}
					}
				}()
				<-ctx.Done()
				ctx = nil
				r = iorw.NewEOFCloseSeekReader(pr, false)
			} else {
				err = ferr
			}
		}
	}

	return
}

//PIPES return file content if file exists else empty string
func PIPES(path string) (cntnt string, err error) {
	var r io.Reader = nil
	if r, err = PIPE(path); err == nil {
		if r != nil {
			var rc io.Closer = nil
			rc, _ = r.(io.Closer)
			func() {
				defer func() {
					if rc != nil {
						rc.Close()
						rc = nil
					}
					r = nil
				}()
				cntnt, err = iorw.ReaderToString(r)
			}()
		}
	}
	return
}

//SET if file exists replace content else create file and append content
func SET(path string, a ...interface{}) (err error) {
	if f, ferr := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); ferr == nil {
		func() {
			defer f.Close()
			iorw.Fprint(f, a...)
		}()
	} else {
		err = ferr
	}
	return
}

//APPEND if file exists append content else create file and append content
func APPEND(path string, a ...interface{}) (err error) {
	if f, ferr := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); ferr == nil {
		func() {
			defer f.Close()
			iorw.Fprint(f, a...)
		}()
	} else {
		err = ferr
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
			if al > 1 {
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
	ABS            func(path string) string                                                                                     `json:"abs"`
	LS             func(path ...string) (finfos []FileInfo)                                                                     `json:"ls"`
	FIND           func(path ...string) (finfos []FileInfo)                                                                     `json:"find"`
	MKDIR          func(path ...interface{}) bool                                                                               `json:"mkdir"`
	MKDIRALL       func(path ...interface{}) bool                                                                               `json:"mkdirall"`
	RM             func(path string) bool                                                                                       `json:"rm"`
	MV             func(path string, destpath string) bool                                                                      `json:"mv"`
	TOUCH          func(path string) bool                                                                                       `json:"touch"`
	FINFOPATHSJSON func(a ...FileInfo) (s string)                                                                               `json:"finfopathsjson"`
	PIPE           func(path string) (r io.Reader)                                                                              `json:"pipe"`
	PIPES          func(path string) (s string)                                                                                 `json:"pipes"`
	CAT            func(path string) (r io.Reader)                                                                              `json:"cat"`
	MULTICAT       func(path ...string) (r io.Reader)                                                                           `json:"multicat"`
	CATS           func(path string) (s string)                                                                                 `json:"cats"`
	MULTICATS      func(path ...string) (s string)                                                                              `json:"multicats"`
	SET            func(path string, a ...interface{}) bool                                                                     `json:"set"`
	APPEND         func(path string, a ...interface{}) bool                                                                     `json:"append"`
	DUMMYFINFO     func(name string, path string, absolutepath string, size int64, mod os.FileMode, modtime time.Time) FileInfo `json:"dummyfino"`
}

//NewFSUtils return instance of FSUtils
func NewFSUtils() (fsutlsstrct FSUtils) {
	fsutlsstrct = FSUtils{
		ABS: func(path string) (abspath string) {
			abspath, _ = ABS(path)
			return
		},
		FIND: func(path ...string) (finfos []FileInfo) {
			if len(path) == 1 {
				finfos, _ = FIND(path[0])
			} else if len(path) == 2 {
				finfos, _ = FIND(path[0], path[1])
			}
			return
		},
		LS: func(path ...string) (finfos []FileInfo) {
			if len(path) == 1 {
				finfos, _ = LS(path[0])
			} else if len(path) == 2 {
				finfos, _ = LS(path[0], path[1])
			}
			return
		},
		MKDIR: func(path ...interface{}) bool {
			if len(path) == 1 {
				if pth, _ := path[0].(string); pth != "" {
					if err := MKDIR(pth); err == nil {
						return true
					}
				}
			}
			return false
		},
		MKDIRALL: func(path ...interface{}) bool {
			if len(path) == 0 {
				if pth, _ := path[0].(string); pth != "" {
					if err := MKDIRALL(pth); err == nil {
						return true
					}
				}
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
		PIPE: func(path string) (r io.Reader) {
			if catr, err := PIPE(path); err == nil {
				r = catr
			}
			return
		}, PIPES: func(path string) (s string) {
			if cats, err := PIPES(path); err == nil {
				s = cats
			}
			return
		},
		CAT: func(path string) (r io.Reader) {
			if catr, err := CAT(path); err == nil {
				r = catr
			}
			return
		}, CATS: func(path string) (s string) {
			if cats, err := CATS(path); err == nil {
				s = cats
			}
			return
		},
		SET: func(path string, a ...interface{}) bool {
			if err := SET(path, a...); err == nil {
				return true
			}
			return false
		},
		APPEND: func(path string, a ...interface{}) bool {
			if err := APPEND(path, a...); err == nil {
				return true
			}
			return false
		},
		FINFOPATHSJSON: func(a ...FileInfo) (s string) {
			s = FINFOPATHSJSON(a...)
			return
		},
		DUMMYFINFO: func(name string, path string, absolutepath string, size int64, mod os.FileMode, modtime time.Time) (finfo FileInfo) {
			finfo = newFileInfo(name, path, absolutepath, size, mod, modtime)
			return
		}}
	return
}

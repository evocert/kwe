package parameters

import (
	"bufio"
	"io"
	"mime/multipart"
	http "net/http"
	url "net/url"
	"strings"
)

//Parameters -> structure containing parameters
type Parameters struct {
	standard  map[string][]string
	filesdata map[string][]interface{}
}

var emptyParmVal = []string{}
var emptyParamFile = []interface{}{}

//StandardKeys - list of standard parameters names (keys)
func (params *Parameters) StandardKeys() (keys []string) {
	if len(params.standard) > 0 {
		if keys == nil {
			keys = make([]string, len(params.standard))
		}
		ki := 0
		for k := range params.standard {
			keys[ki] = k
			ki++
		}
	}
	return keys
}

//FileKeys - list of file parameters names (keys)
func (params *Parameters) FileKeys() (keys []string) {
	if len(params.filesdata) > 0 {
		if keys == nil {
			keys = make([]string, len(params.filesdata))
		}
		ki := 0
		for k := range params.filesdata {
			keys[ki] = k
			ki++
		}
	}
	return keys
}

//SetParameter -> set or append parameter value
//pname : name
//pvalue : value of strings to add
//clear : clear existing value of parameter
func (params *Parameters) SetParameter(pname string, clear bool, pvalue ...string) {
	if pname = strings.ToUpper(strings.TrimSpace(pname)); pname == "" {
		return
	}
	if params.standard == nil {
		params.standard = make(map[string][]string)
	}
	if val, ok := params.standard[pname]; ok {
		if clear {
			val = nil
			params.standard[pname] = nil
			val = []string{}
		}
		if len(pvalue) > 0 {
			val = append(val, pvalue...)
		}
		params.standard[pname] = val
	} else {
		if len(pvalue) > 0 {
			params.standard[pname] = pvalue[:]
		} else {
			params.standard[pname] = []string{}
		}
	}
}

//ContainsParameter -> check if parameter exist
//pname : name
func (params *Parameters) ContainsParameter(pname string) bool {
	if pname = strings.ToUpper(strings.TrimSpace(pname)); pname == "" {
		return false
	}
	if params.standard == nil {
		return false
	}
	_, ok := params.standard[pname]
	return ok
}

//RemoveParameter  -> remove parameter and return any slice of string value
func (params *Parameters) RemoveParameter(pname string) (value []string) {
	if pname = strings.ToUpper(strings.TrimSpace(pname)); pname == "" {
		return
	}
	if params.standard == nil {
		return
	}
	if _, ok := params.standard[pname]; ok {
		value = params.standard[pname][:]
		params.standard[pname] = nil
		delete(params.standard, pname)
	}
	return
}

//SetFileParameter -> set or append file parameter value
//pname : name
//pfile : value of interface to add either FileHeader from mime/multipart or any io.Reader implementation
//clear : clear existing value of parameter
func (params *Parameters) SetFileParameter(pname string, clear bool, pfile ...interface{}) {
	if pname = strings.ToUpper(strings.TrimSpace(pname)); pname == "" {
		return
	}
	if params.filesdata == nil {
		params.filesdata = make(map[string][]interface{})
	}
	if val, ok := params.filesdata[pname]; ok {
		if clear {
			val = nil
			params.filesdata[pname] = nil
			val = []interface{}{}
		}
		if len(pfile) > 0 {
			for _, pf := range pfile {
				/*if fheader, fheaderok := pf.(multipart.FileHeader); fheaderok {
					fheader.
					if fv, fverr := fheader.Open(); fverr == nil {
						if rval, rvalok := fv.(io.Reader); rvalok {
							val = append(val, rval)
						}
					}
				} else {*/
				val = append(val, pf)
				//}
			}
		}
		params.filesdata[pname] = val
	} else {
		if len(pfile) > 0 {
			val = []interface{}{}
			for _, pf := range pfile {
				if fheader, fheaderok := pf.(multipart.FileHeader); fheaderok {
					if fv, fverr := fheader.Open(); fverr == nil {
						if rval, rvalok := fv.(io.Reader); rvalok {
							val = append(val, rval)
						}
					}
				} else {
					val = append(val, pf)
				}
			}
			params.filesdata[pname] = val
		} else {
			params.filesdata[pname] = []interface{}{}
		}
	}
}

//ContainsFileParameter -> check if file parameter exist
//pname : name
func (params *Parameters) ContainsFileParameter(pname string) bool {
	if pname = strings.ToUpper(strings.TrimSpace(pname)); pname == "" {
		return false
	}
	if params.filesdata == nil {
		return false
	}
	_, ok := params.filesdata[pname]
	return ok
}

//Parameter - return a specific parameter values
func (params *Parameters) Parameter(pname string, index ...int) []string {
	if pname = strings.ToUpper(strings.TrimSpace(pname)); pname != "" {
		if params.standard != nil {
			if stdv, ok := params.standard[pname]; ok {
				if stdl := len(stdv); stdl > 0 {
					if il := len(index); il > 0 {
						idx := []int{}
						for _, id := range index {
							if id >= 0 && id < stdl {
								idx = append(idx, id)
							}
						}
						if len(idx) > 0 {
							stdvls := make([]string, len(idx))
							for in, id := range idx {
								stdvls[in] = stdv[id]
							}
							return stdvls
						}
					} else {
						return stdv
					}
				}
			}
		}
	}
	return emptyParmVal
}

//StringParameter return parameter as string concatenated with sep
func (params *Parameters) StringParameter(pname string, sep string, index ...int) (s string) {
	if pval := params.Parameter(pname, index...); len(pval) > 0 {
		return strings.Join(pval, sep)
	}
	if pval := params.FileReader(pname, index...); len(pval) > 0 {

		var rnrtos = func(br *bufio.Reader) (bs string, err error) {
			rns := make([]rune, 1024)
			rnsi := 0
			if br != nil {
				for {
					rn, size, rnerr := br.ReadRune()
					if size > 0 {
						rns[rnsi] = rn
						rnsi++
						if rnsi == len(rns) {
							bs += string(rns[:rnsi])
							rnsi = 0
						}
					}
					if rnerr != nil {
						if rnerr != io.EOF {
							err = rnerr
						}
						break
					}
				}
			}
			if rnsi > 0 {
				bs += string(rns[:rnsi])
				rnsi = 0
			}
			return
		}
		if sep == "" {
			s, _ = rnrtos(bufio.NewReader(io.MultiReader(pval...)))
			return
		}
		var bfr *bufio.Reader = nil
		for rn, r := range pval {
			if bfr == nil {
				bfr = bufio.NewReader(r)
			} else {
				bfr.Reset(r)
			}
			if bs, bserr := rnrtos(bfr); bserr == nil {
				s += bs
				if rn < len(pval)-1 {
					s += sep
				}
			} else if bserr != nil {
				break
			}
		}
	}
	return ""
}

//FileReader return file parameter - array of io.Reader
func (params *Parameters) FileReader(pname string, index ...int) (rdrs []io.Reader) {
	if flsv := params.FileParameter(pname, index...); len(flsv) > 0 {
		rdrs = make([]io.Reader, len(flsv))
		for nfls, fls := range flsv {
			if fhead, fheadok := fls.(multipart.FileHeader); fheadok {
				rdrs[nfls], _ = fhead.Open()
			} else if fr, frok := fls.(io.Reader); frok {
				rdrs[nfls] = fr
			}
		}
	}
	return
}

//FileName return file parameter name - array of string
func (params *Parameters) FileName(pname string, index ...int) (nmes []string) {
	if flsv := params.FileParameter(pname, index...); len(flsv) > 0 {
		nmes = make([]string, len(flsv))
		for nfls, fls := range flsv {
			if fhead, fheadok := fls.(multipart.FileHeader); fheadok {
				nmes[nfls] = fhead.Filename
			}
		}
	}
	return
}

//FileSize return file parameter size - array of int64)
func (params *Parameters) FileSize(pname string, index ...int) (sizes []int64) {
	if flsv := params.FileParameter(pname, index...); len(flsv) > 0 {
		sizes = make([]int64, len(flsv))
		for nfls, fls := range flsv {
			if fhead, fheadok := fls.(multipart.FileHeader); fheadok {
				sizes[nfls] = fhead.Size
			}
		}
	}
	return
}

//FileParameter return file paramater - array of file
func (params *Parameters) FileParameter(pname string, index ...int) []interface{} {
	if pname = strings.ToUpper(strings.TrimSpace(pname)); pname != "" {
		if params.filesdata != nil {
			if flsv, ok := params.filesdata[pname]; ok {
				if flsl := len(flsv); flsl > 0 {
					if il := len(index); il > 0 {
						idx := []int{}
						for _, id := range index {
							if id >= 0 && id < flsl {
								idx = append(idx, id)
							}
						}
						if len(idx) > 0 {
							flsvls := make([]interface{}, len(idx))
							for in, id := range idx {
								flsvls[in] = flsv[id]
							}
							return flsvls
						}
					} else {
						return flsv
					}
				}
			}
		}
	}
	return emptyParamFile
}

//CleanupParameters function that can be called to assist in cleaning up instance of Parameter container
func (params *Parameters) CleanupParameters() {
	if params.standard != nil {
		for pname := range params.standard {
			params.standard[pname] = nil
			delete(params.standard, pname)
		}
		params.standard = nil
	}
	if params.filesdata != nil {
		for pname := range params.filesdata {
			params.filesdata[pname] = nil
			delete(params.filesdata, pname)
		}
		params.filesdata = nil
	}
}

//NewParameters return new instance of Paramaters container
func NewParameters() *Parameters {
	return &Parameters{}
}

//LoadParametersFromRawURL - populate paramaters just from raw url
func LoadParametersFromRawURL(params *Parameters, rawURL string) {
	if rawURL != "" {
		if urlvals, e := url.ParseQuery(rawURL); e == nil {
			if urlvals != nil {
				for pname, pvalue := range urlvals {
					params.SetParameter(pname, false, pvalue...)
				}
			}
		}
	}
}

//LoadParametersFromHTTPRequest - Load Parameters from http.Request
func LoadParametersFromHTTPRequest(params *Parameters, r *http.Request) {
	if r.URL != nil {
		LoadParametersFromRawURL(params, r.URL.RawQuery)
		r.URL.RawQuery = ""
	}
	if err := r.ParseMultipartForm(0); err == nil {
		if r.MultipartForm != nil {
			for pname, pvalue := range r.MultipartForm.Value {
				params.SetParameter(pname, false, pvalue...)
			}
			for pname, pfile := range r.MultipartForm.File {
				if len(pfile) > 0 {
					pfilei := []interface{}{}
					for pf := range pfile {
						pfilei = append(pfilei, pf)
					}
					params.SetFileParameter(pname, false, pfilei...)
					pfilei = nil
				}
			}
		} else if r.Form != nil {
			for pname, pvalue := range r.Form {
				params.SetParameter(pname, false, pvalue...)
			}
		}
	} else if err := r.ParseForm(); err == nil {
		if r.Form != nil {
			for pname, pvalue := range r.Form {
				params.SetParameter(pname, false, pvalue...)
			}
		}
	}
}

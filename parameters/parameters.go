package parameters

import (
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
func (params *Parameters) Parameter(pname string) []string {
	if pname = strings.ToUpper(strings.TrimSpace(pname)); pname != "" {
		if params.standard != nil {
			if _, ok := params.standard[pname]; ok {
				return params.standard[pname]
			}
		}
	}
	return emptyParmVal
}

//StringParameter return parameter as string concatenated with sep
func (params *Parameters) StringParameter(pname string, sep string) string {
	if pval := params.Parameter(pname); len(pval) > 0 {
		return strings.Join(pval, sep)
	}
	return ""
}

//FileParameter return file paramater - array of file
func (params *Parameters) FileParameter(pname string) []interface{} {
	if pname = strings.ToUpper(strings.TrimSpace(pname)); pname != "" {
		if params.filesdata != nil {
			if _, ok := params.filesdata[pname]; ok {
				return params.filesdata[pname]
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

package parameters

import (
	"bufio"
	"io"
	"mime/multipart"
	http "net/http"
	url "net/url"
	"strings"
)

type ParametersAPI interface {
	StandardKeys() []string
	FileKeys() []string
	SetParameter(string, bool, ...string)
	AppendPhrase(...string)
	Phrases() []string
	ContainsPhrase(...string) bool
	ContainsParameter(string) bool
	RemoveParameter(string) []string
	SetFileParameter(string, bool, ...interface{})
	ContainsFileParameter(string) bool
	Parameter(string, ...int) []string
	StringParameter(string, string, ...int) string
	FileReader(string, ...int) []io.Reader
	FileName(string, ...int) []string
	FileSize(string, ...int) []int64
	FileParameter(string, ...int) []interface{}
	CleanupParameters()
}

//Parameters -> structure containing parameters
type Parameters struct {
	phrases   []string
	standard  map[string][]string
	filesdata map[string][]interface{}
}

var emptyParmVal = []string{}
var emptyParamFile = []interface{}{}

func (params *Parameters) AppendPhrase(phrases ...string) {
	if params != nil {
		if len(phrases) > 0 {
			for phrsn := range phrases {
				if phrs := phrases[phrsn]; phrs != "" {
					if params.phrases == nil {
						params.phrases = []string{}
					}
					params.phrases = append(params.phrases, phrs)
				}
			}
		}
	}
}

func (params *Parameters) Phrases() (phrases []string) {
	if params != nil && len(params.phrases) > 0 {
		phrases = params.phrases[:]
	}
	return
}

func (params *Parameters) ContainsPhrase(...string) (exists bool) {

	return
}

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
			for pf := range pfile {
				/*if fheader, fheaderok := pf.(multipart.FileHeader); fheaderok {
					fheader.
					if fv, fverr := fheader.Open(); fverr == nil {
						if rval, rvalok := fv.(io.Reader); rvalok {
							val = append(val, rval)
						}
					}
				} else {*/
				val = append(val, pfile[pf])
				//}
			}
		}
		params.filesdata[pname] = val
	} else {
		if len(pfile) > 0 {
			val = []interface{}{}
			for pf := range pfile {
				if fheader, fheaderok := pfile[pf].(multipart.FileHeader); fheaderok {
					if fv, fverr := fheader.Open(); fverr == nil {
						if rval, rvalok := fv.(io.Reader); rvalok {
							val = append(val, rval)
						}
					}
				} else {
					val = append(val, pfile[pf])
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
						for idn := range index {
							if id := index[idn]; id >= 0 && id < stdl {
								idx = append(idx, id)
							}
						}
						if len(idx) > 0 {
							stdvls := make([]string, len(idx))
							for in := range idx {
								stdvls[in] = stdv[idx[in]]
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
		var bfr *bufio.Reader = nil
		for rn := range pval {
			if r := pval[rn]; r != nil {
				if bfr == nil {
					bfr = bufio.NewReader(r)
				} else {
					bfr.Reset(r)
				}
				if bfr != nil {
					if bs, bserr := rnrtos(bfr); bserr == nil {
						s += bs
					} else if bserr != nil {
						break
					}
				}
			}
			if rn < len(pval)-1 {
				s += sep
			}
		}
	}
	return
}

//FileReader return file parameter - array of io.Reader
func (params *Parameters) FileReader(pname string, index ...int) (rdrs []io.Reader) {
	if flsv := params.FileParameter(pname, index...); len(flsv) > 0 {
		rdrs = make([]io.Reader, len(flsv))
		for nfls, fls := range flsv {
			if fhead, fheadok := fls.(*multipart.FileHeader); fheadok {
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
		for nfls := range flsv {
			if fhead, fheadok := flsv[nfls].(*multipart.FileHeader); fheadok {
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
	if params != nil {
		if phrsl := len(params.phrases); phrsl > 0 {
			for phrsl > 0 {
				params.phrases[0] = ""
				params.phrases = params.phrases[1:]
				phrsl--
			}
			params.phrases = nil
		}
	}
}

//NewParameters return new instance of Paramaters container
func NewParameters() *Parameters {
	return &Parameters{}
}

//LoadParametersFromRawURL - populate paramaters just from raw url
func LoadParametersFromRawURL(params ParametersAPI, rawURL string) {
	if params != nil && rawURL != "" {
		if rawURL != "" {
			var phrases = []string{}
			var rawUrls = strings.Split(rawURL, "&")
			rawURL = ""
			for _, rwurl := range rawUrls {
				if rwurl != "" {
					if strings.Contains(rwurl, "=") {
						rawURL += rwurl + "&"
					} else {
						phrases = append(phrases, rwurl)
					}
				}
			}
			if len(rawURL) > 1 && strings.HasSuffix(rawURL, "&") {
				rawURL = rawURL[:len(rawURL)-1]
			}
			if urlvals, e := url.ParseQuery(rawURL); e == nil {
				if len(urlvals) > 0 {
					for pname, pvalue := range urlvals {
						params.SetParameter(pname, false, pvalue...)
					}
				}
			}
			if len(phrases) > 0 {
				params.AppendPhrase(phrases...)
			}
		}
	}
}

//LoadParametersFromUrlValues - Load Parameters from url.Values
func LoadParametersFromUrlValues(params ParametersAPI, urlvalues url.Values) (err error) {
	if params != nil && urlvalues != nil {
		for pname, pvalue := range urlvalues {
			params.SetParameter(pname, false, pvalue...)
		}
	}
	return
}

//LoadParametersFromMultipartForm - Load Parameters from *multipart.Form
func LoadParametersFromMultipartForm(params ParametersAPI, mpartform *multipart.Form) (err error) {
	if params != nil && mpartform != nil {
		for pname, pvalue := range mpartform.Value {
			params.SetParameter(pname, false, pvalue...)
		}
		for pname, pfile := range mpartform.File {
			if len(pfile) > 0 {
				pfilei := []interface{}{}
				for _, pf := range pfile {
					pfilei = append(pfilei, pf)
				}
				params.SetFileParameter(pname, false, pfilei...)
				pfilei = nil
			}
		}
	}
	return
}

//LoadParametersFromHTTPRequest - Load Parameters from http.Request
func LoadParametersFromHTTPRequest(params ParametersAPI, r *http.Request) {
	if params != nil {
		if r.URL != nil {
			LoadParametersFromRawURL(params, r.URL.RawQuery)
			r.URL.RawQuery = ""
		}
		if err := r.ParseMultipartForm(0); err == nil {
			if r.MultipartForm != nil {
				LoadParametersFromMultipartForm(params, r.MultipartForm)
			} else if r.Form != nil {
				LoadParametersFromUrlValues(params, r.Form)
			}
		} else if err := r.ParseForm(); err == nil {
			LoadParametersFromUrlValues(params, r.Form)
		}
	}
}

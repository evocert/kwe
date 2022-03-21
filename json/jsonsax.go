package json

import (
	"encoding/json"
	jsn "encoding/json"
	"io"
	"strings"

	"github.com/evocert/kwe/iorw"
)

type JsonSax struct {
	r             io.Reader
	jsndcdr       *jsn.Decoder
	lstmnctpe     jsn.Delim
	crntky        string
	Level         int
	LevelKeys     map[int]string
	LevelType     map[int]rune
	Object        interface{}
	CallFunc      func(interface{}, ...interface{}) interface{}
	errfunc       interface{}
	Error         func(jsnsx *JsonSax, lasterr error)
	eoffunc       interface{}
	Eof           func(jsnsx *JsonSax)
	appendarrfunc interface{}
	AppendArr     func(jsnsx *JsonSax, val interface{}, vtpe rune)
	setkeyvalfunc interface{}
	SetKeyVal     func(jsnsx *JsonSax, k string, val interface{}, vtpe rune)
	startobjfunc  interface{}
	StartObj      func(jsnsx *JsonSax, k string)
	endobjfunc    interface{}
	EndObj        func(jsnsx *JsonSax) bool
	startarrfunc  interface{}
	StartArr      func(jsnsx *JsonSax, k string)
	endarrfunc    interface{}
	EndArr        func(jsnsx *JsonSax) bool
	OnClose       func(jsnsx *JsonSax)
}

func NewJsonSAX(a ...interface{}) (jsnsx *JsonSax) {
	var errfunc interface{}
	var eoffunc interface{}
	var appendarrfunc interface{}
	var setkeyvalfunc interface{}
	var startobjfunc interface{}
	var endobjfunc interface{}
	var startarrfunc interface{}
	var endarrfunc interface{}
	if al := len(a); al > 0 {
		ai := 0
		for ai < al {
			if d := a[ai]; d != nil {
				if mp, _ := d.(map[string]interface{}); len(mp) > 0 {
					for mk, mv := range mp {
						if strings.EqualFold(mk, "error") {
							errfunc = mv
						} else if strings.EqualFold(mk, "eof") {
							eoffunc = mv
						} else if strings.EqualFold(mk, "appendarr") {
							appendarrfunc = mv
						} else if strings.EqualFold(mk, "setkeyval") {
							setkeyvalfunc = mv
						} else if strings.EqualFold(mk, "startobj") {
							startobjfunc = mv
						} else if strings.EqualFold(mk, "endobj") {
							endobjfunc = mv
						} else if strings.EqualFold(mk, "startarr") {
							startarrfunc = mv
						} else if strings.EqualFold(mk, "endarr") {
							endarrfunc = mv
						}
					}
					al--
					a = a[1:]
					continue
				}
			}
			ai++
		}
	}
	r := iorw.NewMultiArgsReader(a...)

	jsnsx = &JsonSax{r: r, LevelKeys: map[int]string{}, LevelType: map[int]rune{},
		errfunc:       errfunc,
		eoffunc:       eoffunc,
		appendarrfunc: appendarrfunc,
		setkeyvalfunc: setkeyvalfunc,
		startobjfunc:  startobjfunc,
		endobjfunc:    endobjfunc,
		startarrfunc:  startarrfunc,
		endarrfunc:    endarrfunc}
	jsnsx.jsndcdr = jsn.NewDecoder(jsnsx)
	return
}

func (jsnsx *JsonSax) Read(p []byte) (n int, err error) {
	n, err = jsnsx.r.Read(p)
	return
}

func (jsnsx *JsonSax) Close() {
	if jsnsx != nil {
		if jsnsx.OnClose != nil {
			jsnsx.OnClose(jsnsx)
		}
		if jsnsx.Object != nil {
			jsnsx.Object = nil
		}
		if jsnsx.jsndcdr != nil {
			jsnsx.jsndcdr = nil
		}
		if jsnsx.r != nil {
			jsnsx.r = nil
		}
		if jsnsx.setkeyvalfunc != nil {
			jsnsx.setkeyvalfunc = nil
		}
		if jsnsx.SetKeyVal != nil {
			jsnsx.SetKeyVal = nil
		}
		if jsnsx.appendarrfunc != nil {
			jsnsx.appendarrfunc = nil
		}
		if jsnsx.AppendArr != nil {
			jsnsx.AppendArr = nil
		}
		if jsnsx.endarrfunc != nil {
			jsnsx.endarrfunc = nil
		}
		if jsnsx.EndArr != nil {
			jsnsx.EndArr = nil
		}
		if jsnsx.endobjfunc != nil {
			jsnsx.endobjfunc = nil
		}
		if jsnsx.EndObj != nil {
			jsnsx.EndObj = nil
		}
		if jsnsx.startarrfunc != nil {
			jsnsx.startarrfunc = nil
		}
		if jsnsx.StartArr != nil {
			jsnsx.StartArr = nil
		}
		if jsnsx.startobjfunc != nil {
			jsnsx.startobjfunc = nil
		}
		if jsnsx.StartObj != nil {
			jsnsx.StartObj = nil
		}
		if jsnsx.errfunc != nil {
			jsnsx.errfunc = nil
		}
		if jsnsx.Error != nil {
			jsnsx.Error = nil
		}
		if jsnsx.eoffunc != nil {
			jsnsx.eoffunc = nil
		}
		if jsnsx.Eof != nil {
			jsnsx.Eof = nil
		}
		if jsnsx.LevelKeys != nil {
			for k := range jsnsx.LevelKeys {
				delete(jsnsx.LevelKeys, k)
			}
			jsnsx.LevelKeys = nil
		}
		if jsnsx.LevelType != nil {
			for k := range jsnsx.LevelType {
				delete(jsnsx.LevelType, k)
			}
			jsnsx.LevelType = nil
		}
		jsnsx = nil
	}
}

func (jsnsx *JsonSax) Next() (canContinue bool, err error) {
	canContinue = true
	if tkn, tknerr := jsnsx.jsndcdr.Token(); tknerr == nil {
		if dlm, _ := tkn.(jsn.Delim); rune(dlm) == '{' {
			jsnsx.lstmnctpe = dlm
			jsnsx.Level++
			if jsnsx.crntky != "" {
				jsnsx.LevelKeys[jsnsx.Level] = jsnsx.crntky
				jsnsx.crntky = ""
			}
			jsnsx.LevelType[jsnsx.Level] = 'O'
			if jsnsx.CallFunc != nil && jsnsx.startobjfunc != nil {
				jsnsx.CallFunc(jsnsx.startobjfunc, jsnsx, jsnsx.LevelKeys[jsnsx.Level])
			}
			if jsnsx.StartObj != nil {
				jsnsx.StartObj(jsnsx, jsnsx.LevelKeys[jsnsx.Level])
			}
		} else if dlm, _ := tkn.(json.Delim); rune(dlm) == '}' {
			jsnsx.crntky = ""
			if jsnsx.CallFunc != nil && jsnsx.endobjfunc != nil {
				if cncntnue, _ := jsnsx.CallFunc(jsnsx.endobjfunc, jsnsx).(bool); cncntnue {
					canContinue = !cncntnue
				}
			}
			if jsnsx.EndObj != nil {
				canContinue = !jsnsx.EndObj(jsnsx)
			}
			delete(jsnsx.LevelKeys, jsnsx.Level)
			delete(jsnsx.LevelType, jsnsx.Level)
			jsnsx.Level--
		} else if dlm, _ := tkn.(json.Delim); rune(dlm) == '[' {
			jsnsx.lstmnctpe = dlm
			jsnsx.Level++
			if jsnsx.crntky != "" {
				jsnsx.LevelKeys[jsnsx.Level] = jsnsx.crntky
				jsnsx.crntky = ""
			}
			jsnsx.LevelType[jsnsx.Level] = 'A'
			if jsnsx.CallFunc != nil && jsnsx.startarrfunc != nil {
				jsnsx.CallFunc(jsnsx.startarrfunc, jsnsx, jsnsx.LevelKeys[jsnsx.Level])
			}
			if jsnsx.StartArr != nil {
				jsnsx.StartArr(jsnsx, jsnsx.LevelKeys[jsnsx.Level])
			}
		} else if dlm, _ := tkn.(json.Delim); rune(dlm) == ']' {
			if jsnsx.CallFunc != nil && jsnsx.endarrfunc != nil {
				if cncntnue, _ := jsnsx.CallFunc(jsnsx.endarrfunc, jsnsx).(bool); cncntnue {
					canContinue = !cncntnue
				}
			}
			if jsnsx.EndArr != nil {
				canContinue = !jsnsx.EndArr(jsnsx)
			}
			delete(jsnsx.LevelKeys, jsnsx.Level)
			delete(jsnsx.LevelType, jsnsx.Level)
			jsnsx.Level--
		} else if dlm, _ := tkn.(json.Delim); rune(dlm) == rune(0) && tkn == nil {
			if jsnsx.lstmnctpe == '{' {
				if jsnsx.CallFunc != nil && jsnsx.setkeyvalfunc != nil {
					jsnsx.CallFunc(jsnsx.setkeyvalfunc, jsnsx, jsnsx.crntky, nil, 'N')
				}
				if jsnsx.SetKeyVal != nil {
					jsnsx.SetKeyVal(jsnsx, jsnsx.crntky, nil, 'N')
				}
				jsnsx.crntky = ""
			} else {
				if jsnsx.CallFunc != nil && jsnsx.appendarrfunc != nil {
					jsnsx.CallFunc(jsnsx.appendarrfunc, jsnsx, nil, 'N')
				}
				if jsnsx.AppendArr != nil {
					jsnsx.AppendArr(jsnsx, nil, 'N')
				}
			}
		} else if s, sk := tkn.(string); sk {
			if jsnsx.lstmnctpe == '{' {
				if jsnsx.crntky == "" {
					jsnsx.crntky = s
				} else {
					if jsnsx.CallFunc != nil && jsnsx.setkeyvalfunc != nil {
						jsnsx.CallFunc(jsnsx.setkeyvalfunc, jsnsx, jsnsx.crntky, s, 'S')
					}
					if jsnsx.SetKeyVal != nil {
						jsnsx.SetKeyVal(jsnsx, jsnsx.crntky, s, 'S')
					}
					jsnsx.crntky = ""
				}
			} else if jsnsx.lstmnctpe == '[' {
				if jsnsx.CallFunc != nil && jsnsx.appendarrfunc != nil {
					jsnsx.CallFunc(jsnsx.appendarrfunc, jsnsx, s, 'S')
				}
				if jsnsx.AppendArr != nil {
					jsnsx.AppendArr(jsnsx, s, 'S')
				}
			}
		} else if b, bk := tkn.(bool); bk {
			if jsnsx.lstmnctpe == '{' {
				if jsnsx.CallFunc != nil && jsnsx.setkeyvalfunc != nil {
					jsnsx.CallFunc(jsnsx.setkeyvalfunc, jsnsx, jsnsx.crntky, b, 'B')
				}
				if jsnsx.SetKeyVal != nil {
					jsnsx.SetKeyVal(jsnsx, jsnsx.crntky, b, 'B')
				}
				jsnsx.crntky = ""
			} else {
				if jsnsx.CallFunc != nil && jsnsx.appendarrfunc != nil {
					jsnsx.CallFunc(jsnsx.appendarrfunc, jsnsx, b, 'B')
				}
				if jsnsx.AppendArr != nil {
					jsnsx.AppendArr(jsnsx, b, 'B')
				}
			}
		} else if f, fk := tkn.(float64); fk {
			if jsnsx.lstmnctpe == '{' {
				if jsnsx.CallFunc != nil && jsnsx.setkeyvalfunc != nil {
					jsnsx.CallFunc(jsnsx.setkeyvalfunc, jsnsx, jsnsx.crntky, f, 'F')
				}
				if jsnsx.SetKeyVal != nil {
					jsnsx.SetKeyVal(jsnsx, jsnsx.crntky, f, 'F')
				}
				jsnsx.crntky = ""
			} else {
				if jsnsx.CallFunc != nil && jsnsx.appendarrfunc != nil {
					jsnsx.CallFunc(jsnsx.appendarrfunc, jsnsx, f, 'F')
				}
				if jsnsx.AppendArr != nil {
					jsnsx.AppendArr(jsnsx, f, 'F')
				}
			}
		}
	} else {
		canContinue = false
		if tknerr == io.EOF {
			if jsnsx.CallFunc != nil && jsnsx.eoffunc != nil {
				jsnsx.CallFunc(jsnsx.eoffunc, jsnsx)
				if jsnsx.Eof != nil {
					jsnsx.Eof(jsnsx)
				}
				tknerr = nil
			} else if jsnsx.Eof != nil {
				jsnsx.Eof(jsnsx)
				tknerr = nil
			}
		}
		err = tknerr
	}
	return canContinue, err
}

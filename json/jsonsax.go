package json

import (
	"encoding/json"
	jsn "encoding/json"
	"io"

	"github.com/evocert/kwe/iorw"
)

type JsonSax struct {
	r         io.Reader
	jsndcdr   *jsn.Decoder
	lstmnctpe jsn.Delim
	crntky    string
	Level     int
	LevelKeys map[int]string
	LevelType map[int]rune
	AppendArr func(jsnsx *JsonSax, val interface{}, vtpe rune)
	SetKeyVal func(jsnsx *JsonSax, k string, val interface{}, vtpe rune)
	StartObj  func(jsnsx *JsonSax, k string)
	EndObj    func(jsnsx *JsonSax) bool
	StartArr  func(jsnsx *JsonSax, k string)
	EndArr    func(jsnsx *JsonSax) bool
}

func NewJsonSAX(a ...interface{}) (jsnsx *JsonSax) {
	r := iorw.NewMultiArgsReader(a...)
	jsndcdr := jsn.NewDecoder(r)
	jsnsx = &JsonSax{r: r, jsndcdr: jsndcdr, LevelKeys: map[int]string{}, LevelType: map[int]rune{}}
	return
}

func (jsnsx *JsonSax) Close() {
	if jsnsx != nil {
		if jsnsx.jsndcdr != nil {
			jsnsx.jsndcdr = nil
		}
		if jsnsx.r != nil {
			jsnsx.r = nil
		}
		if jsnsx.SetKeyVal != nil {
			jsnsx.SetKeyVal = nil
		}
		if jsnsx.AppendArr != nil {
			jsnsx.AppendArr = nil
		}
		if jsnsx.EndArr != nil {
			jsnsx.EndArr = nil
		}
		if jsnsx.EndObj != nil {
			jsnsx.EndObj = nil
		}
		if jsnsx.StartArr != nil {
			jsnsx.StartArr = nil
		}
		if jsnsx.StartObj != nil {
			jsnsx.StartObj = nil
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

func (jsnsx *JsonSax) ParseNext() (canContinue bool, err error) {
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
			if jsnsx.StartObj != nil {
				jsnsx.StartObj(jsnsx, jsnsx.LevelKeys[jsnsx.Level])
			}
		} else if dlm, _ := tkn.(json.Delim); rune(dlm) == '}' {
			jsnsx.crntky = ""
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
			if jsnsx.StartArr != nil {
				jsnsx.StartArr(jsnsx, jsnsx.LevelKeys[jsnsx.Level])
			}
		} else if dlm, _ := tkn.(json.Delim); rune(dlm) == ']' {
			if jsnsx.EndArr != nil {
				canContinue = !jsnsx.EndArr(jsnsx)
			}
			delete(jsnsx.LevelKeys, jsnsx.Level)
			delete(jsnsx.LevelType, jsnsx.Level)
			jsnsx.Level--
		} else if s, sk := tkn.(string); sk {
			if jsnsx.lstmnctpe == '{' {
				if jsnsx.crntky == "" {
					jsnsx.crntky = s
				} else {
					if jsnsx.SetKeyVal != nil {
						jsnsx.SetKeyVal(jsnsx, jsnsx.crntky, s, 'S')
					}
					jsnsx.crntky = ""
				}
			} else if jsnsx.lstmnctpe == '[' {
				if jsnsx.AppendArr != nil {
					jsnsx.AppendArr(jsnsx, s, 'S')
				}
			}
		} else if b, bk := tkn.(bool); bk {
			if jsnsx.lstmnctpe == '{' {
				if jsnsx.SetKeyVal != nil {
					jsnsx.SetKeyVal(jsnsx, jsnsx.crntky, b, 'B')
				}
				jsnsx.crntky = ""
			} else {
				if jsnsx.AppendArr != nil {
					jsnsx.AppendArr(jsnsx, b, 'B')
				}
			}
		} else if f, fk := tkn.(float64); fk {
			if jsnsx.lstmnctpe == '{' {
				if jsnsx.SetKeyVal != nil {
					jsnsx.SetKeyVal(jsnsx, jsnsx.crntky, f, 'F')
				}
				jsnsx.crntky = ""
			} else {
				if jsnsx.AppendArr != nil {
					jsnsx.AppendArr(jsnsx, f, 'F')
				}
			}
		}
	} else {
		canContinue = false
		err = tknerr
	}
	return canContinue, err
}

package jsonrw

import (
	"encoding/json"
	"io"
)

//ReaderParser - struct
type ReaderParser struct {
	rdr      io.Reader
	dcdr     *json.Decoder
	Delim    json.Delim
	IsDelim  bool
	String   string
	IsString bool
	IsNumber bool
	Float    float64
	Int      int64
	IsFloat  bool
	IsBool   bool
	Bool     bool
	IsNull   bool

	PrevToken json.Token
	Token     json.Token
	Depth     int
}

//NewReaderParser - JSON Reader Parser
func NewReaderParser(r io.Reader) (jsnr *ReaderParser) {
	jsnr = &ReaderParser{rdr: r, Depth: 0, dcdr: json.NewDecoder(r), PrevToken: nil, Token: nil}
	return
}

//PopulateMap - populate map[string]interface{}
func (jsnr *ReaderParser) PopulateMap(mptopop map[string]interface{}) (err error) {
	if mptopop != nil {
		err = jsnr.dcdr.Decode(&mptopop)
	}
	return
}

//PopulateArray - populate []interface{}
func (jsnr *ReaderParser) PopulateArray(arr []interface{}) (err error) {
	if arr != nil {
		err = jsnr.dcdr.Decode(&arr)
	}
	return
}

func (jsnr *ReaderParser) setValueState() {
	if jsnr.Token == nil {
		jsnr.IsNull = true
	} else {
		if dlm, dlmok := jsnr.Token.(json.Delim); dlmok {
			jsnr.Delim = dlm
			jsnr.IsDelim = true
		} else if nr, nrok := jsnr.Token.(json.Number); nrok {
			if nr != "" {
				jsnr.IsNumber = true
				if fltv, flterr := nr.Float64(); flterr == nil {
					if intv, interr := nr.Int64(); interr == nil {
						if float64(intv) == fltv {
							jsnr.Int = intv
						} else {
							jsnr.Float = fltv
							jsnr.IsFloat = true
						}
					} else {
						jsnr.Float = fltv
						jsnr.IsFloat = true
					}
				}
			}
		} else if str, strok := jsnr.Token.(string); strok {
			jsnr.IsString = true
			jsnr.String = str
		} else if bl, blok := jsnr.Token.(bool); blok {
			jsnr.IsBool = true
			jsnr.Bool = bl
		}
	}
}

//More - wrap arround json.Decoder.More
func (jsnr *ReaderParser) More(jsonevent ...func(*ReaderParser, error) error) (more bool, err error) {
	if jsnr != nil {
		if jsnr.dcdr != nil {
			var jseventerr error = nil
			jsnr.IsDelim, jsnr.IsBool, jsnr.IsString, jsnr.IsNumber, jsnr.IsFloat, jsnr.IsNull = false, false, false, false, false, false
			jsnr.String, jsnr.Bool, jsnr.Float, jsnr.Int = "", false, 0, 0
			if jsnr.Token == nil {
				jsnr.Token, err = jsnr.dcdr.Token()
				if err == nil {
					jsnr.setValueState()
				}
				if len(jsonevent) == 1 && jsonevent[0] != nil {
					if jseventerr = jsonevent[0](jsnr, err); jseventerr != nil {
						err = jseventerr
					}
				}
				if err == nil {
					if dlm, dlmok := jsnr.Token.(json.Delim); dlmok {
						jsnr.Delim = dlm
						jsnr.Depth++
					}
					more = true
				}
				jsnr.PrevToken = jsnr.Token

			} else if jsnr.dcdr.More() {
				jsnr.Token, err = jsnr.dcdr.Token()
				if err == nil {
					jsnr.setValueState()
				}
				if len(jsonevent) == 1 && jsonevent[0] != nil {
					if jseventerr = jsonevent[0](jsnr, err); jseventerr != nil {
						err = jseventerr
					}
				}
				if err == nil {
					if _, dlmok := jsnr.Token.(json.Delim); dlmok {
						jsnr.Depth++
					}
					more = true
				}
				jsnr.PrevToken = jsnr.Token
			} else {
				jsnr.Token, err = jsnr.dcdr.Token()
				if err == nil {
					if dlm, dlmok := jsnr.Token.(json.Delim); dlmok {
						jsnr.Delim = dlm
					}
				}
				jsnr.Depth--
				if len(jsonevent) == 1 && jsonevent[0] != nil {
					if jseventerr = jsonevent[0](jsnr, err); jseventerr != nil {
						err = jseventerr
					}
				}
				jsnr.PrevToken = jsnr.Token
				if err == nil {
					more = true
				} else {
					more = false
				}
			}
		}
	}
	return
}

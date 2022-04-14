package xml

import (
	"encoding/xml"
	"io"
	"strings"

	"github.com/evocert/kwe/iorw"
)

type Attr struct {
	Space string
	Name  string
	Value string
}

type XmlSax struct {
	r              io.Reader
	Object         interface{}
	Level          int
	lastStartLevel int
	lastEndLevel   int
	LevelNames     map[int][]string
	LevelAttrs     map[int][]Attr
	xmldcdr        *xml.Decoder
	CallFunc       func(interface{}, ...interface{}) interface{}
	errfunc        interface{}
	Error          func(xmlsx *XmlSax, lasterr error)
	startelemfunc  interface{}
	StartElement   func(xmlsx *XmlSax, space string, name string, attrs ...Attr) (done bool)
	elemdatafunc   interface{}
	ElemData       func(xmlsx *XmlSax, data []byte)
	endelemfunc    interface{}
	EndElement     func(xmlsx *XmlSax, space string, name string) (done bool)
	closefunc      interface{}
	OnClose        func(xmlsx *XmlSax)
	eoffunc        interface{}
	Eof            func(xmlsx *XmlSax)
}

func NewXmlSAX(a ...interface{}) (xmlsx *XmlSax) {
	var errfunc interface{}
	var eoffunc interface{}
	var startelemfunc interface{}
	var endelemfunc interface{}
	var elemdatafunc interface{}
	var closefunc interface{}
	if al := len(a); al > 0 {
		ai := 0
		for ai < al {
			if d := a[ai]; d != nil {
				if mp, _ := d.(map[string]interface{}); len(mp) > 0 {
					for mk, mv := range mp {
						if strings.EqualFold(mk, "error") {
							errfunc = mv
						} else if strings.EqualFold(mk, "close") {
							closefunc = mv
						} else if strings.EqualFold(mk, "eof") {
							eoffunc = mv
						} else if strings.EqualFold(mk, "startelem") {
							startelemfunc = mv
						} else if strings.EqualFold(mk, "elemdata") {
							elemdatafunc = mv
						} else if strings.EqualFold(mk, "endelem") {
							endelemfunc = mv
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
	xmldcdr := xml.NewDecoder(r)
	xmlsx = &XmlSax{r: r, xmldcdr: xmldcdr, LevelNames: map[int][]string{}, LevelAttrs: map[int][]Attr{}, closefunc: closefunc,
		startelemfunc: startelemfunc, endelemfunc: endelemfunc, elemdatafunc: elemdatafunc, errfunc: errfunc, eoffunc: eoffunc}
	return
}

func (xmlsx *XmlSax) Close() (err error) {
	if xmlsx != nil {
		if xmlsx.closefunc != nil {
			if xmlsx.CallFunc != nil {
				xmlsx.CallFunc(xmlsx.closefunc, xmlsx)
			}
			xmlsx.closefunc = nil
		}
		if xmlsx.OnClose != nil {
			xmlsx.OnClose(xmlsx)
			xmlsx.OnClose = nil
		}
		if xmlsx.Object != nil {
			xmlsx.Object = nil
		}
		if xmlsx.xmldcdr != nil {
			xmlsx.xmldcdr = nil
		}
		if xmlsx.Eof != nil {
			xmlsx.Eof = nil
		}
		if xmlsx.eoffunc != nil {
			xmlsx.eoffunc = nil
		}
		if xmlsx.Error != nil {
			xmlsx.Error = nil
		}
		if xmlsx.errfunc != nil {
			xmlsx.errfunc = nil
		}
		if xmlsx.ElemData != nil {
			xmlsx.ElemData = nil
		}
		if xmlsx.elemdatafunc != nil {
			xmlsx.elemdatafunc = nil
		}
		if xmlsx.EndElement != nil {
			xmlsx.EndElement = nil
		}
		if xmlsx.endelemfunc != nil {
			xmlsx.endelemfunc = nil
		}
		if xmlsx.StartElement != nil {
			xmlsx.StartElement = nil
		}
		if xmlsx.startelemfunc != nil {
			xmlsx.startelemfunc = nil
		}
		if xmlsx.LevelNames != nil {
			for k := range xmlsx.LevelNames {
				xmlsx.LevelNames[k] = nil
				delete(xmlsx.LevelNames, k)
			}
			xmlsx.LevelNames = nil
		}
		if xmlsx.LevelAttrs != nil {
			for k := range xmlsx.LevelAttrs {
				xmlsx.LevelAttrs[k] = nil
				delete(xmlsx.LevelAttrs, k)
			}
			xmlsx.LevelAttrs = nil
		}
		xmlsx = nil
	}
	return
}

func (xmlsx *XmlSax) Next() (canContinue bool, err error) {
	canContinue = true
	if tkn, tknerr := xmlsx.xmldcdr.Token(); tknerr == nil {
		switch se := tkn.(type) {
		case xml.StartElement:
			xmlsx.Level++
			xmlsx.LevelNames[xmlsx.Level] = []string{se.Name.Space, se.Name.Local}
			if attrl := len(se.Attr); attrl > 0 {
				var attrs = make([]Attr, attrl)
				for attrn, attr := range se.Attr {
					attrs[attrn] = Attr{Name: attr.Name.Local, Space: attr.Name.Space, Value: attr.Value}
				}
				xmlsx.LevelAttrs[xmlsx.Level] = attrs
			}
			if xmlsx.StartElement != nil {
				canContinue = !xmlsx.StartElement(xmlsx, xmlsx.LevelNames[xmlsx.Level][0], xmlsx.LevelNames[xmlsx.Level][1], xmlsx.LevelAttrs[xmlsx.Level]...)
			}
		case xml.EndElement:
			if xmlsx.EndElement != nil {
				canContinue = !xmlsx.EndElement(xmlsx, se.Name.Space, se.Name.Local)
			}
			if _, lvlnmesok := xmlsx.LevelNames[xmlsx.Level]; lvlnmesok {
				xmlsx.LevelNames[xmlsx.Level] = nil
				delete(xmlsx.LevelNames, xmlsx.Level)
			}
			if _, lvlattrsok := xmlsx.LevelAttrs[xmlsx.Level]; lvlattrsok {
				xmlsx.LevelAttrs[xmlsx.Level] = nil
				delete(xmlsx.LevelAttrs, xmlsx.Level)
			}
			xmlsx.Level--
		case xml.CharData:
			if len(se) > 0 {
				if xmlsx.ElemData != nil {
					xmlsx.ElemData(xmlsx, se)
				}
			}
		}
	} else {
		canContinue = false
		if tknerr == io.EOF {
			if xmlsx.CallFunc != nil && xmlsx.eoffunc != nil {
				xmlsx.CallFunc(xmlsx.eoffunc, xmlsx)
				if xmlsx.Eof != nil {
					xmlsx.Eof(xmlsx)
				}
				tknerr = nil
			} else if xmlsx.Eof != nil {
				xmlsx.Eof(xmlsx)
				tknerr = nil
			}
		}
		err = tknerr
	}
	return
}

func XmlJsonToString(level int, a ...interface{}) (s string, err error) {
	func() {
		buff := iorw.NewBuffer()
		defer buff.Close()
		if err = WriteXmlJson(buff, level, a...); err == nil {
			s = buff.String()
		}
	}()
	return
}

func WriteXmlJson(w io.Writer, level int, a ...interface{}) (err error) {
	if len(a) > 0 {
		a = append([]interface{}{`<?xml version="1.0" encoding="UTF-8"?>`}, a...)
		xml := iorw.NewMultiArgsReader(a...)
		if json, jsonerr := Convert(xml); jsonerr == nil {
			io.Copy(w, json)
		} else {
			err = jsonerr
		}
	}
	/*func() {
		dec := NewXmlSAX(a...)
		buf := iorw.NewBuffer()
		defer buf.Close()
		enc := json.NewEncoder(buf)
		enc.SetIndent("", "")
		enc.SetEscapeHTML(false)
		if al := len(a); al > 0 {
			dec.StartElement = func(xmlsx *XmlSax, space, name string, attrs ...Attr) (done bool) {
				if xmlsx.Level >= level {
					if xmlsx.lastStartLevel != xmlsx.Level {
						iorw.Fprint(w, "{")
					} else if xmlsx.lastStartLevel == xmlsx.Level {
						iorw.Fprint(w, ",")
					}
					if attrsl := len(attrs); attrsl > 0 {
						iorw.Fprint(w, "{")
						for attn, atr := range attrs {
							iorw.Fprint(w, "@")
							enc.Encode(atr.Name)
							iorw.Fprint(w, strings.TrimSpace(buf.String()))
							buf.Clear()
							iorw.Fprint(w, ":")
							enc.Encode(atr.Value)
							iorw.Fprint(w, strings.TrimSpace(buf.String()))
							buf.Clear()
							if attn < attrsl {
								iorw.Fprint(w, ",")
							}
						}
					}
					if space != "" {
						name = space + ":" + name
					}
					enc.Encode(name)
					iorw.Fprint(w, strings.TrimSpace(buf.String()))
					buf.Clear()
					iorw.Fprint(w, ":")
				}
				return
			}

			dec.ElemData = func(xmlsx *XmlSax, data []byte) {
				if xmlsx.Level > 0 {
					enc.Encode(string(data))
					iorw.Fprint(w, strings.TrimSpace(buf.String()))
					buf.Clear()
				}
			}

			dec.EndElement = func(xmlsx *XmlSax, space, name string) (done bool) {
				if xmlsx.lastEndLevel != xmlsx.lastStartLevel {
					if xmlsx.Level >= level-1 {
						iorw.Fprint(w, "}")
					}
				}
				return
			}

			for err == nil {
				if next, errnext := dec.Next(); (!next && errnext == nil) || errnext != nil {
					if errnext != io.EOF {
						err = errnext
					}
					break
				}
			}
		} else {
			enc.Encode(nil)
		}
	}()*/
	return
}

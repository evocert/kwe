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
	r             io.Reader
	Level         int
	LevelNames    map[int][]string
	LevelAttrs    map[int][]Attr
	xmldcdr       *xml.Decoder
	CallFunc      func(interface{}, ...interface{}) interface{}
	errfunc       interface{}
	Error         func(xmlsn *XmlSax, lasterr error)
	startelemfunc interface{}
	StartElement  func(xmlsn *XmlSax, space string, name string, attrs ...[]Attr) (done bool)
	elemdatafunc  interface{}
	ElemData      func(xmlsn *XmlSax, data []byte)
	endelemfunc   interface{}
	EndElement    func(xmlsn *XmlSax, space string, name string) (done bool)
	OnClose       func(xmlsn *XmlSax)
	eoffunc       interface{}
	Eof           func(xmlsn *XmlSax)
}

func NewXmlSAX(a ...interface{}) (xmlsx *XmlSax) {
	var errfunc interface{}
	var eoffunc interface{}
	var startelemfunc interface{}
	var endelemfunc interface{}
	var elemdatafunc interface{}
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
	xmlsx = &XmlSax{r: r, xmldcdr: xmldcdr, LevelNames: map[int][]string{}, LevelAttrs: map[int][]Attr{},
		startelemfunc: startelemfunc, endelemfunc: endelemfunc, elemdatafunc: elemdatafunc, errfunc: errfunc, eoffunc: eoffunc}
	return
}

func (xmlsx *XmlSax) Close() (err error) {
	if xmlsx != nil {
		if xmlsx.OnClose != nil {
			xmlsx.OnClose(xmlsx)
			xmlsx.OnClose = nil
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
				canContinue = !xmlsx.StartElement(xmlsx, xmlsx.LevelNames[xmlsx.Level][0], xmlsx.LevelNames[xmlsx.Level][1], xmlsx.LevelAttrs[xmlsx.Level])
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

package xml

import (
	"encoding/xml"
	"io"

	"github.com/evocert/kwe/iorw"
)

type XmlSax struct {
	r            io.Reader
	Level        int
	LevelNames   map[int][]string
	LevelAttrs   map[int][][]string
	xmldcdr      *xml.Decoder
	StartElement func(xmlsn *XmlSax, space string, name string, attrs ...[][]string) (done bool)
	ElemData     func(xmlsn *XmlSax, data []byte)
	EndElement   func(xmlsn *XmlSax, space string, name string) (done bool)
}

func NewXmlSAX(a ...interface{}) (xmlsx *XmlSax) {
	r := iorw.NewMultiArgsReader(a...)
	xmldcdr := xml.NewDecoder(r)
	xmlsx = &XmlSax{r: r, xmldcdr: xmldcdr, LevelNames: map[int][]string{}, LevelAttrs: map[int][][]string{}}
	return
}

func (xmlsx *XmlSax) Close() (err error) {
	if xmlsx != nil {
		if xmlsx.xmldcdr != nil {
			xmlsx.xmldcdr = nil
		}
		if xmlsx.ElemData != nil {
			xmlsx.ElemData = nil
		}
		if xmlsx.EndElement != nil {
			xmlsx.EndElement = nil
		}
		if xmlsx.StartElement != nil {
			xmlsx.StartElement = nil
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

func (xmlsx *XmlSax) ParseNext() (canContinue bool, err error) {
	canContinue = true
	if tkn, tknerr := xmlsx.xmldcdr.Token(); tknerr == nil {
		switch se := tkn.(type) {
		case xml.StartElement:
			xmlsx.Level++
			xmlsx.LevelNames[xmlsx.Level] = []string{se.Name.Space, se.Name.Local}
			if attrl := len(se.Attr); attrl > 0 {
				var attrs = make([][]string, attrl)
				for attrn, attr := range se.Attr {
					attrs[attrn] = []string{attr.Name.Space, attr.Name.Local, attr.Value}
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
		if tknerr != nil {
			canContinue = false
			err = tknerr
		}
	}
	return
}

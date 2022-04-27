package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/json"
	"github.com/evocert/kwe/xml"
)

type RWSReader struct {
	rdr       io.RuneReader
	lsterr    error
	strmtype  string
	strmstngs map[string]interface{}
	sqlrws    *sql.Rows
	coltypes  []*ColumnType
	cls       []string
	firstdata bool
	eof       bool
	data      []interface{}
	jsnsx     *json.JsonSax
	xmlsx     *xml.XmlSax
}

func newRWSReader(sqlrws *sql.Rows, strmstngs map[string]interface{}) (rwsrrdr *RWSReader, err error) {
	if len(strmstngs) > 0 {
		var strmtype = ""
		var stngs map[string]interface{}
		var rdr io.RuneReader = nil
		if strmtype, _ = strmstngs["stream-type"].(string); strmtype != "" {
			delete(strmstngs, "stream-type")
		}
		for strmk, strmv := range strmstngs {
			if strmv != nil {
				if strmtype == "csv" || strmtype == "json" || strmtype == "xml" {
					if strmtype == "csv" && strings.Contains("row-delim,col-delim,headers", strmk) {
						if stngs == nil {
							stngs = map[string]interface{}{}
						}
						stngs[strmk] = strmv
					} else if strmk == "data" {
						if args, _ := strmv.([]interface{}); len(args) > 0 {
							rdr = iorw.NewMultiArgsReader(args...)
						} else {
							rdr = iorw.NewMultiArgsReader(strmv)
						}
					}
				}
			}
		}
		if rdr != nil && ((strmtype == "csv" && len(stngs) > 0) || strmtype == "json" || strmtype == "xml") {
			rwsrrdr = &RWSReader{strmstngs: stngs, rdr: rdr, strmtype: strmtype}
		} else {
			err = fmt.Errorf("%s", "Unsupported Data Stream Configutaion")
		}
	} else if sqlrws != nil {
		rwsrrdr = &RWSReader{sqlrws: sqlrws}
	}
	return
}

func (rwsrdr *RWSReader) indexOfColumn(col string) (i int) {
	for p, v := range rwsrdr.cls {
		if v == col {
			return p
		}
	}
	return -1
}

func (rwsrdr *RWSReader) Close() (err error) {
	if rwsrdr != nil {
		if rwsrdr.sqlrws != nil {
			err = rwsrdr.sqlrws.Close()
			rwsrdr.sqlrws = nil
		}
		if rwsrdr.strmstngs != nil {
			for strmk := range rwsrdr.strmstngs {
				delete(rwsrdr.strmstngs, strmk)
			}
			rwsrdr.strmstngs = nil
		}
		if rwsrdr.coltypes != nil {
			rwsrdr.coltypes = nil
		}
		if rwsrdr.cls != nil {
			rwsrdr.cls = nil
		}
		if rwsrdr.data != nil {
			rwsrdr.data = nil
		}
		if rwsrdr.strmtype != "" {
			rwsrdr.strmtype = ""
		}
		if rwsrdr.jsnsx != nil {
			rwsrdr.jsnsx.Close()
			rwsrdr.jsnsx = nil
		}
		if rwsrdr.xmlsx != nil {
			rwsrdr.xmlsx.Close()
			rwsrdr.xmlsx = nil
		}
		rwsrdr = nil
	}
	return
}

func (rwsrdr *RWSReader) Err() (err error) {
	if rwsrdr != nil {
		if rwsrdr.sqlrws != nil && len(rwsrdr.strmstngs) == 0 && rwsrdr.rdr == nil {
			rwsrdr.lsterr = rwsrdr.sqlrws.Err()
		}
		err = rwsrdr.lsterr
	}
	return
}

func (rwsrdr *RWSReader) Next() (nxt bool) {
	if rwsrdr != nil {
		if rwsrdr.sqlrws != nil && len(rwsrdr.strmstngs) == 0 && rwsrdr.rdr == nil {
			nxt = rwsrdr.sqlrws.Next()
		} else if rwsrdr.rdr != nil /* && len(rwsrdr.strmstngs) > 0*/ && rwsrdr.strmtype != "" {
			if rwsrdr.firstdata {
				rwsrdr.firstdata = false
			} else {
				if rwsrdr.strmtype == "csv" && len(rwsrdr.data) > 0 {
					rwsrdr.data = nil
				}
				if len(rwsrdr.cls) == 0 {
					prepRWSColumns(rwsrdr)
				} else {
					rwsrdr.lsterr = populateRWSStreamData(rwsrdr)
				}
			}
			if rwsrdr.lsterr == nil {
				if nxt = !rwsrdr.eof; nxt {
					nxt = len(rwsrdr.data) > 0
				}
			} else {
				if rwsrdr.lsterr == io.EOF {
					rwsrdr.lsterr = nil
				}
				nxt = false
			}
		}
	}
	return
}

func (rwsrdr *RWSReader) Scan(dest ...interface{}) (err error) {
	if rwsrdr != nil {
		if rwsrdr.sqlrws != nil && len(rwsrdr.strmstngs) == 0 && rwsrdr.rdr == nil {
			err = rwsrdr.sqlrws.Scan(dest...)
		} else if rwsrdr.rdr != nil /*&& len(rwsrdr.strmstngs) > 0*/ && rwsrdr.strmtype != "" {
			if len(rwsrdr.data) > 0 && len(rwsrdr.cls) == len(rwsrdr.data) {
				if len(dest) == len(rwsrdr.data) {
					for destn, dta := range rwsrdr.data {
						switch d := dest[destn].(type) {
						case *interface{}:
							if d == nil {
								return errors.New("destination pointer is nil")
							}
							*d = dta
						}
					}
				}
			}
		}
	}
	return
}

func (rwsrdr *RWSReader) ColumnTypes() (coltypes []*ColumnType, err error) {
	if rwsrdr != nil {
		if err = prepRWSColumns(rwsrdr); err == nil {
			coltypes = rwsrdr.coltypes[:]
		} else {
			rwsrdr.lsterr = err
		}
	}
	return
}

func (rwsrdr *RWSReader) Columns() (cls []string, err error) {
	if rwsrdr != nil {
		if err = prepRWSColumns(rwsrdr); err == nil {
			cls = rwsrdr.cls[:]
		} else {
			rwsrdr.lsterr = err
		}
	}
	return
}

func prepRWSColumns(rwsrdr *RWSReader) (err error) {
	if len(rwsrdr.cls) == 0 {
		if rwsrdr.sqlrws != nil && len(rwsrdr.strmstngs) == 0 && rwsrdr.rdr == nil {
			if cltps, cltpserr := rwsrdr.sqlrws.ColumnTypes(); cltpserr == nil {
				if cls, clserr := rwsrdr.sqlrws.Columns(); clserr == nil {
					if len(rwsrdr.cls) == 0 && len(cls) > 0 {
						rwsrdr.cls = cls[:]
						rwsrdr.coltypes = columnTypes(cltps, cls)
					}
				} else {
					err = clserr
				}
			} else {
				err = cltpserr
			}
		} else if rwsrdr.sqlrws == nil && ((rwsrdr.strmtype == "csv" && len(rwsrdr.strmstngs) > 0) || rwsrdr.strmtype == "json" || rwsrdr.strmtype == "xml") && rwsrdr.rdr != nil {
			err = populateRWSStreamData(rwsrdr)
		}
	}
	return
}

func populateRWSStreamData(rwsrdr *RWSReader) (err error) {
	if rwsrdr.strmtype == "csv" {
		var headers, _ = rwsrdr.strmstngs["headers"].(bool)
		var coldelim, _ = rwsrdr.strmstngs["col-delim"].(string)
		var rowdelim, _ = rwsrdr.strmstngs["row-delim"].(string)
		err = parseCSVRWS(rwsrdr, headers, []rune(coldelim)[:], []rune(rowdelim)[:], rwsrdr.rdr, len(rwsrdr.cls) == 0)
	} else if rwsrdr.strmtype == "json" {
		err = parseJSONRWS(rwsrdr, rwsrdr.rdr, len(rwsrdr.cls) == 0)
	} else if rwsrdr.strmtype == "xml" {
		err = parseXMLRWS(rwsrdr, rwsrdr.rdr, len(rwsrdr.cls) == 0)
	}
	return
}

func parseXMLRWS(rwsrdr *RWSReader, rdr io.RuneReader, readcols bool) (err error) {

	if rwsrdr.xmlsx == nil {
		if r, _ := rdr.(io.Reader); r != nil {
			rwsrdr.xmlsx = xml.NewXmlSAX(io.Reader(r))
			rwsrdr.xmlsx.Eof = func(xmlsn *xml.XmlSax) {
				xmlsn.Close()
				rwsrdr.eof = true
			}
		}
	}
	if rwsrdr.xmlsx != nil {
		if len(rwsrdr.cls) == 0 {
			rwsrdr.xmlsx.StartElement = nil

			rwsrdr.xmlsx.ElemData = func(xmlsn *xml.XmlSax, data []byte) {
				if xmlsn.Level == 3 {
					if cl, dl := len(rwsrdr.cls), len(rwsrdr.data); (cl == 0 && dl == 0) || dl == cl {
						rwsrdr.data = append(rwsrdr.data, string(data))
					} else if cl != dl && dl >= cl {
						rwsrdr.data[cl] = string(data)
					}
					rwsrdr.cls = append(rwsrdr.cls, xmlsn.LevelNames[xmlsn.Level][1])

					rwsrdr.coltypes = append(rwsrdr.coltypes, &ColumnType{
						name:              xmlsn.LevelNames[xmlsn.Level][1],
						hasNullable:       true,
						hasPrecisionScale: false,
						hasLength:         false,
						databaseType:      "VARCHAR",
						length:            0,
						precision:         0,
						scale:             0,
						scanType:          reflect.TypeOf(""),
					})
				}
			}

			rwsrdr.xmlsx.EndElement = func(xmlsn *xml.XmlSax, space, name string) (done bool) {
				if xmlsn.Level == 2 {
					done = true
					rwsrdr.firstdata = len(rwsrdr.data) > 0
				}
				return
			}

			for {
				if canext, prseerr := rwsrdr.xmlsx.Next(); !canext || prseerr != nil {
					if prseerr != nil {
						err = prseerr
					}
					break
				}
			}
		} else {
			rwsrdr.xmlsx.StartElement = nil

			rwsrdr.xmlsx.ElemData = func(xmlsn *xml.XmlSax, data []byte) {
				if xmlsn.Level == 3 {
					if dtai := rwsrdr.indexOfColumn(xmlsn.LevelNames[xmlsn.Level][1]); dtai > -1 {
						if cl := len(rwsrdr.cls); cl > 0 && len(rwsrdr.data) != cl {
							rwsrdr.data = make([]interface{}, cl)
						}
						rwsrdr.data[dtai] = string(data)
					}
				}
			}

			rwsrdr.xmlsx.EndElement = func(xmlsn *xml.XmlSax, space, name string) (done bool) {
				if xmlsn.Level == 2 {
					done = true
				}
				return
			}

			for !rwsrdr.eof {
				if canext, prseerr := rwsrdr.xmlsx.Next(); !canext || prseerr != nil {
					if prseerr != nil {
						err = prseerr
					}
					break
				}
			}
		}
	}
	return
}

func parseJSONRWS(rwsrdr *RWSReader, rdr io.RuneReader, readcols bool) (err error) {

	if rwsrdr.jsnsx == nil {
		if r, _ := rdr.(io.Reader); r != nil {
			rwsrdr.jsnsx = json.NewJsonSAX(io.Reader(r))
			rwsrdr.jsnsx.Eof = func(jsnsx *json.JsonSax) {
				jsnsx.Close()
				rwsrdr.eof = true
			}
		}
	}
	if rwsrdr.jsnsx != nil {
		if len(rwsrdr.cls) == 0 {
			var data []interface{} = nil
			var cltpe *ColumnType = nil

			rwsrdr.jsnsx.AppendArr = nil

			rwsrdr.jsnsx.SetKeyVal = func(jsnsx *json.JsonSax, k string, val interface{}, vtpe rune) {
				if jsnsx.LevelKeys[2] == "columns" && jsnsx.Level == 3 {
					if cltpe == nil {
						cltpe = &ColumnType{}
					}
					if k == "name" {
						cltpe.name, _ = val.(string)
					} else if k == "title" {
						cltpe.name, _ = val.(string)
					} else if k == "dbtype" {
						cltpe.databaseType, _ = val.(string)
					} else if k == "length" {
						if lngth, _ := val.(float64); lngth > 0 {
							cltpe.length = int64(lngth)
							cltpe.hasLength = true
						}
					} else if k == "precision" {
						if prcsn, _ := val.(float64); prcsn > 0 {
							cltpe.precision = int64(prcsn)
							cltpe.hasPrecisionScale = true
						}
					} else if k == "scale" {
						if scle, _ := val.(float64); scle > 0 {
							cltpe.scale = int64(scle)
							cltpe.hasPrecisionScale = true
						}
					}
				} else if jsnsx.Level == 2 && jsnsx.LevelType[1] == 'A' && jsnsx.LevelType[2] == 'O' {
					if cltpe == nil {
						cltpe = &ColumnType{}
					}
					cltpe.name = k
					data = append(data, val)
					if cltpe != nil {
						rwsrdr.cls = append(rwsrdr.cls, cltpe.name)
						rwsrdr.coltypes = append(rwsrdr.coltypes, cltpe)
						cltpe = nil
					}
				}
			}

			rwsrdr.jsnsx.StartObj = nil

			rwsrdr.jsnsx.EndObj = func(jsnsx *json.JsonSax) (done bool) {
				if jsnsx.Level == 3 && jsnsx.LevelKeys[2] == "columns" {
					if cltpe != nil {
						rwsrdr.cls = append(rwsrdr.cls, cltpe.name)
						rwsrdr.coltypes = append(rwsrdr.coltypes, cltpe)
						cltpe = nil
					}
				} else if jsnsx.Level == 2 && jsnsx.LevelType[1] == 'A' {
					if len(data) > 0 && len(data) == len(rwsrdr.cls) {
						if len(rwsrdr.data) < len(data) {
							rwsrdr.data = data[:]
							rwsrdr.firstdata = true
						}
					}
					done = true
				}
				return
			}

			rwsrdr.jsnsx.StartArr = nil

			rwsrdr.jsnsx.EndArr = func(jsnsx *json.JsonSax) (done bool) {
				if jsnsx.Level == 2 && jsnsx.LevelKeys[jsnsx.Level] == "columns" {

					return true
				} else if jsnsx.Level == 1 && jsnsx.LevelType[jsnsx.Level] == 'A' {
					err = io.EOF
					return true
				}
				return
			}
			for !rwsrdr.eof {
				if canext, prseerr := rwsrdr.jsnsx.Next(); !canext || prseerr != nil {
					if prseerr != nil {
						err = prseerr
					}
					break
				}
			}
		} else {
			if len(rwsrdr.data) == 0 {

				rwsrdr.jsnsx.StartObj = nil
				rwsrdr.jsnsx.SetKeyVal = nil
				rwsrdr.jsnsx.StartArr = nil

				rwsrdr.jsnsx.SetKeyVal = func(jsnsx *json.JsonSax, k string, val interface{}, vtpe rune) {
					if jsnsx.Level == 2 && jsnsx.LevelType[1] == 'A' && jsnsx.LevelType[2] == 'O' {
						if dtai := rwsrdr.indexOfColumn(k); dtai > -1 {
							if cl := len(rwsrdr.cls); cl > 0 && len(rwsrdr.data) != cl {
								rwsrdr.data = make([]interface{}, cl)
							}
							rwsrdr.data[dtai] = val
						}
					}
				}

				rwsrdr.jsnsx.AppendArr = func(jsnsx *json.JsonSax, val interface{}, vtpe rune) {
					if jsnsx.LevelKeys[2] == "data" && jsnsx.Level == 3 {
						if len(rwsrdr.data) < len(rwsrdr.cls) {
							rwsrdr.data = append(rwsrdr.data, val)
						}
					}
				}

				rwsrdr.jsnsx.EndObj = func(jsnsx *json.JsonSax) (done bool) {
					if jsnsx.Level == 2 && jsnsx.LevelType[1] == 'A' {

						done = true
					}
					return
				}

				rwsrdr.jsnsx.EndArr = func(jsnsx *json.JsonSax) (done bool) {
					if jsnsx.LevelKeys[2] == "data" && jsnsx.Level == 3 {

						done = true
					} else if jsnsx.Level == 1 && jsnsx.LevelType[jsnsx.Level] == 'A' {
						err = io.EOF
						return true
					}
					return
				}
				for !rwsrdr.eof {
					if canext, prseerr := rwsrdr.jsnsx.Next(); !canext || prseerr != nil {
						if prseerr != nil {
							err = prseerr
						}
						break
					}
				}
			} else {
				var dtai = 0
				rwsrdr.jsnsx.StartObj = nil
				rwsrdr.jsnsx.SetKeyVal = nil
				rwsrdr.jsnsx.StartArr = nil

				rwsrdr.jsnsx.SetKeyVal = func(jsnsx *json.JsonSax, k string, val interface{}, vtpe rune) {
					if jsnsx.Level == 2 && jsnsx.LevelType[1] == 'A' && jsnsx.LevelType[2] == 'O' {
						if dtai = rwsrdr.indexOfColumn(k); dtai > -1 {
							if cl := len(rwsrdr.cls); cl > 0 && len(rwsrdr.data) != cl {
								rwsrdr.data = make([]interface{}, cl)
							}
							rwsrdr.data[dtai] = val
						}
					}
				}

				rwsrdr.jsnsx.AppendArr = func(jsnsx *json.JsonSax, val interface{}, vtpe rune) {
					if jsnsx.LevelKeys[2] == "data" && jsnsx.Level == 3 {
						if dtai < len(rwsrdr.cls) && len(rwsrdr.data) == len(rwsrdr.cls) {
							rwsrdr.data[dtai] = val
							dtai++
						}
					}
				}

				rwsrdr.jsnsx.EndObj = func(jsnsx *json.JsonSax) (done bool) {
					if jsnsx.Level == 2 && jsnsx.LevelType[1] == 'A' {

						done = true
					}
					return
				}

				rwsrdr.jsnsx.EndArr = func(jsnsx *json.JsonSax) (done bool) {
					if jsnsx.LevelKeys[2] == "data" && jsnsx.Level == 3 {

						done = true
					} else if jsnsx.Level == 1 && jsnsx.LevelType[jsnsx.Level] == 'A' {
						err = io.EOF
						return true
					}
					return
				}
				for !rwsrdr.eof {
					if canext, prseerr := rwsrdr.jsnsx.Next(); !canext || prseerr != nil {
						if prseerr != nil {
							err = prseerr
						}
						break
					}
				}
			}
		}
	}
	return
}

func parseCSVRWS(rwsrdr *RWSReader, headers bool, coldelim []rune, rowdelim []rune, rdr io.RuneReader, readcols bool) (err error) {
	var tmddata []interface{} = nil
	var tmprunedata []rune = nil
	var tmprndtai = 0
	var coldelimi = 0
	var rowdelimi = 0
	var tmpstr = ""

	var txtrne = rune(0)

	var nextcol = true
	var cancol = true

	var flushTmpRuneData = func() {
		if tmprndtai > 0 {
			tmpstr += string(tmprunedata[0:tmprndtai])
			tmprndtai = 0
		}
	}

	var appendTmpRuneData = func(r rune) {
		if nextcol {
			nextcol = false
			if txtrne == rune(0) && (r == '"' || r == '\'') {
				txtrne = r
				cancol = false
				return
			}
		}
		if tmprunedata == nil {
			tmprunedata = make([]rune, 8192)
		}
		tmprunedata[tmprndtai] = r
		tmprndtai++
		if tmprndtai == 8192 {
			flushTmpRuneData()
		}
	}

	var prvr rune = 0
	var prvc rune = 0

	var coldelimfunc = func(r rune) {
		if cancol {
			if coldelimi > 0 && coldelim[coldelimi-1] == prvc && coldelim[coldelimi] != r {
				for _, rn := range rowdelim[0:rowdelimi] {
					appendTmpRuneData(rn)
				}
				coldelimi = 0
			}
			if coldelim[coldelimi] == r {
				coldelimi++
				if coldelimi == len(coldelim) {
					flushTmpRuneData()
					tmpstr = strings.TrimSpace(tmpstr)
					tmddata = append(tmddata, tmpstr+"")
					tmpstr = ""
					prvc = 0
					coldelimi = 0
					nextcol = true
					tmpstr = ""
					txtrne = rune(0)
				} else {
					prvc = r
				}
			} else {
				if coldelimi > 0 {
					for _, rn := range coldelim[0:coldelimi] {
						appendTmpRuneData(rn)
					}
					rowdelimi = 0
				}
				prvc = r
				appendTmpRuneData(r)
			}
		} else {
			if (txtrne == '"' || txtrne == '\'') && txtrne == r {
				if prvc == txtrne {
					prvc = rune(0)
					appendTmpRuneData(r)
				} else {
					cancol = true
				}
			} else {
				prvc = r
				appendTmpRuneData(r)
			}
		}
	}

	var canrow = true

	var wrapupRow = func() {
		flushTmpRuneData()
		tmpstr = strings.TrimSpace(tmpstr)
		if nextcol || tmpstr != "" {
			tmddata = append(tmddata, tmpstr)
			nextcol = false
		}
		if len(tmddata) > 0 {
			if readcols && headers {
				if len(rwsrdr.cls) == 0 && len(tmddata) > 0 {
					rwsrdr.cls = make([]string, len(tmddata))
					rwsrdr.coltypes = make([]*ColumnType, len(tmddata))
					for tmpdtan, tmpdta := range tmddata {
						if tmpstr, _ = tmpdta.(string); tmpstr != "" {
							if tmpstr = strings.TrimSpace(tmpstr); tmpstr != "" {
								rwsrdr.cls[tmpdtan] = tmpstr
							} else {
								rwsrdr.cls[tmpdtan] = fmt.Sprintf("%s%d", "COL", tmpdtan)
							}
						} else {
							rwsrdr.cls[tmpdtan] = fmt.Sprintf("%s%d", "COL", tmpdtan)
						}
						rwsrdr.coltypes[tmpdtan] = &ColumnType{
							name:              rwsrdr.cls[tmpdtan],
							hasNullable:       true,
							hasPrecisionScale: false,
							hasLength:         false,
							databaseType:      "VARCHAR",
							length:            0,
							precision:         0,
							scale:             0,
							scanType:          reflect.TypeOf(""),
						}
					}
				}
			} else {
				if readcols && len(rwsrdr.cls) == 0 && len(tmddata) > 0 {
					rwsrdr.data = nil
					rwsrdr.data = make([]interface{}, len(tmddata))
					copy(rwsrdr.data, tmddata)
					rwsrdr.cls = make([]string, len(tmddata))
					rwsrdr.coltypes = make([]*ColumnType, len(tmddata))
					for tmpdtan := range tmddata {
						rwsrdr.cls[tmpdtan] = fmt.Sprintf("%s%d", "COL", tmpdtan)
						rwsrdr.coltypes[tmpdtan] = &ColumnType{
							name:              rwsrdr.cls[tmpdtan],
							hasNullable:       true,
							hasPrecisionScale: false,
							hasLength:         false,
							databaseType:      "VARCHAR",
							length:            0,
							precision:         0,
							scale:             0,
							scanType:          reflect.TypeOf(""),
						}
					}
				} else if len(rwsrdr.cls) == len(tmddata) {
					rwsrdr.data = nil
					rwsrdr.data = make([]interface{}, len(tmddata))
					copy(rwsrdr.data, tmddata)
				}
			}
		}
		tmddata = nil
		prvr = 0
		prvc = 0
		rowdelimi = 0
		coldelimi = 0
		tmpstr = ""
		txtrne = rune(0)
	}

	var rowdelimfunc = func(r rune) bool {
		if canrow {
			if rowdelimi > 0 && rowdelim[rowdelimi-1] == prvr && rowdelim[rowdelimi] != r {
				for _, rn := range rowdelim[0:rowdelimi] {
					coldelimfunc(rn)
				}
				rowdelimi = 0
			}
			if rowdelim[rowdelimi] == r {
				rowdelimi++
				if rowdelimi == len(rowdelim) {
					wrapupRow()
					return true
				} else {
					prvr = r
				}
			} else {
				if rowdelimi > 0 {
					for _, rn := range rowdelim[0:rowdelimi] {
						coldelimfunc(rn)
					}
					rowdelimi = 0
				}
				prvr = r
				coldelimfunc(r)
			}
		} else {
			coldelimfunc(r)
		}
		return false
	}

	for {
		if r, s, rerr := rdr.ReadRune(); rerr == nil {
			if s > 0 {
				if rowdelimfunc(r) {
					break
				}
			}
		} else {
			if rerr == io.EOF {
				if coldelimi > 0 || nextcol {
					wrapupRow()
				}
			}
			break
		}
	}
	return
}

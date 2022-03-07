package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/evocert/kwe/iorw"
)

type RWSReader struct {
	rdr       io.RuneReader
	lsterr    error
	strmtype  string
	strmstngs map[string]interface{}
	sqlrws    *sql.Rows
	coltypes  []*ColumnType
	cls       []string
	data      []interface{}
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
				if strmtype == "csv" {
					if strings.Contains("row-delim,col-delim,headers", strmk) {
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
		if rdr != nil && ((strmtype == "csv" && len(stngs) > 0) || strmtype == "json") {
			rwsrrdr = &RWSReader{strmstngs: stngs, rdr: rdr, strmtype: strmtype}
		} else {
			err = fmt.Errorf("%s", "Unsupported Data Stream Configutaion")
		}
	} else if sqlrws != nil {
		rwsrrdr = &RWSReader{sqlrws: sqlrws}
	}
	return
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
		} else if rwsrdr.rdr != nil && len(rwsrdr.strmstngs) > 0 && rwsrdr.strmtype != "" {
			if len(rwsrdr.data) > 0 {
				rwsrdr.data = nil
			}
			if len(rwsrdr.cls) == 0 {
				prepRWSColumns(rwsrdr)
			} else {
				rwsrdr.lsterr = populateRWSStreamData(rwsrdr)
			}
			nxt = len(rwsrdr.data) > 0
		}
	}
	return
}

func (rwsrdr *RWSReader) Scan(dest ...interface{}) (err error) {
	if rwsrdr != nil {
		if rwsrdr.sqlrws != nil && len(rwsrdr.strmstngs) == 0 && rwsrdr.rdr == nil {
			err = rwsrdr.sqlrws.Scan(dest...)
		} else if rwsrdr.rdr != nil && len(rwsrdr.strmstngs) > 0 && rwsrdr.strmtype != "" {
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
		} else if rwsrdr.sqlrws == nil && len(rwsrdr.strmstngs) > 0 && rwsrdr.rdr != nil {
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

package database

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/evocert/kwe/iorw/active"
)

//Reader - struct
type Reader struct {
	*Executor
	rws         *sql.Rows
	rownr       int64
	cls         []string
	cltpes      []*ColumnType
	data        []interface{}
	datamap     map[string]interface{}
	dispdata    []interface{}
	dataref     []interface{}
	wg          *sync.WaitGroup
	OnColumns   interface{}
	OnRow       interface{}
	OnValidData interface{}
}

//ColumnTypes return Column types in form of a slice, 'array', of []*ColumnType values
func (rdr *Reader) ColumnTypes() []*ColumnType {
	return rdr.cltpes
}

//Columns return Column names in form of a slice, 'array', of string values
func (rdr *Reader) Columns() []string {
	return rdr.cls
}

//Data return Displayable data in the form of a slice, 'array', of interface{} values
func (rdr *Reader) Data() []interface{} {
	//go func(somethingDone chan bool) {
	//	defer func() {
	//		somethingDone <- true
	//	}()

	for n := range rdr.data {
		rdr.dispdata[n] = castSQLTypeValue(rdr.data[n], rdr.cltpes[n])
	}
	//}(rset.dosomething)
	//<-rset.dosomething

	return rdr.dispdata[:]
}

var emptymap = map[string]interface{}{}

//DataMap return Displayable data in the form of a map[string]interface{} column and values
func (rdr *Reader) DataMap() (datamap map[string]interface{}) {
	if rdr != nil && len(rdr.data) > 0 && len(rdr.cls) == len(rdr.data) {
		displdata := rdr.Data()
		if rdr.datamap == nil {
			rdr.datamap = map[string]interface{}{}
		}
		for cn, c := range rdr.cls {
			rdr.datamap[c] = displdata[cn]
		}
		return rdr.datamap
	}
	return emptymap
}

//Next return true if able to move focus of Reader to the next underlying record
// or false if the end is reached
func (rdr *Reader) Next() (next bool, err error) {
	validdata := true
	if rdr.isRemote() && rdr.jsndcdr != nil {
		for {
			if rdr.tknlvl == 3 {
				if rdr.data == nil {
					rdr.data = make([]interface{}, len(rdr.cls))
					rdr.dispdata = make([]interface{}, len(rdr.cls))
				}
				if dcerr := rdr.jsndcdr.Decode(&rdr.data); dcerr == nil {
					next = true
					if err == nil {
						if validdata = invokeDataValid(rdr.script, rdr.OnValidData, rdr.rownr, rdr); validdata {
							rdr.rownr++
							next = invokeRow(rdr.script, rdr.OnRow, rdr.rownr, rdr)
						}
					}
				} else {
					next = false
					//err = dcerr
				}
				break
			}
			tkn, tknerr := rdr.jsndcdr.Token()
			if tknerr != nil {
				rdr.lasterr = tknerr
				next = false
				break
			} else {
				if dlm, dlmok := tkn.(json.Delim); dlmok {
					if rdr.lastdlm = dlm.String(); rdr.lastdlm == "{" {
						rdr.tknlvl++
					} else if rdr.lastdlm == "}" {
						rdr.tknlvl--
					} else if rdr.lastdlm == "[" {
						rdr.tknlvl++
					} else if rdr.lastdlm == "]" {
						rdr.tknlvl--
					}
				} else {
					if s, sok := tkn.(string); sok {
						if rdr.tknlvl == 1 && s != "" {

						} else {
							if rdr.tknlvl == 2 && s != "" {
								if s == "data" {

								}
							} else {

							}
						}
					}

				}
			}
		}
		if !next {
			rdr.Close()
		}
	} else {
		if next = rdr.rws.Next(); next {
			if rdr.wg == nil {
				rdr.wg = &sync.WaitGroup{}
			}
			rdr.wg.Add(1)
			if rdr.data == nil {
				rdr.data = make([]interface{}, len(rdr.cls))
				rdr.dataref = make([]interface{}, len(rdr.cls))
				rdr.dispdata = make([]interface{}, len(rdr.cls))
			}
			wg := rdr.wg
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				for n := range rdr.data {
					rdr.dataref[n] = &rdr.data[n]
				}
				if scerr := rdr.rws.Scan(rdr.dataref...); scerr != nil {
					rdr.Close()
					err = scerr
					next = false
				}
			}(wg)
			wg.Wait()
			if err == nil {
				if validdata = invokeDataValid(rdr.script, rdr.OnValidData, rdr.rownr, rdr); validdata {
					rdr.rownr++
					next = invokeRow(rdr.script, rdr.OnRow, rdr.rownr, rdr)
				}
			}
			//}(rset.dosomething)
			//<-rset.dosomething
			if !next {
				rdr.Close()
			}
		} else {
			if rseterr := rdr.rws.Err(); rseterr != nil {
				err = rseterr
				//rdr.lasterr=err
				invokeError(rdr.script, err, rdr.OnError)
			}
			rdr.Close()
		}
	}
	return next, err
}

//Close the Reader as well as the underlying Executor related to this Reader
//After this action the Reader is 'empty' or cleaned up in a golang world
func (rdr *Reader) Close() (err error) {
	if rdr.data != nil {
		rdr.data = nil
	}
	if rdr.dataref != nil {
		rdr.dataref = nil
	}
	if rdr.dispdata != nil {
		rdr.dispdata = nil
	}
	if rdr.cltpes != nil {
		rdr.cltpes = nil
		rdr.cls = nil
	}
	if rdr.wg != nil {
		rdr.wg = nil
	}
	if rdr.OnColumns != nil {
		rdr.OnColumns = nil
	}
	if rdr.OnRow != nil {
		rdr.OnRow = nil
	}
	if rdr.rws != nil {
		rdr.rws.Close()
		rdr.rws = nil
	}
	if rdr.Executor != nil {
		if rdr.lasterr != nil {
			err = rdr.lasterr
			rdr.Executor.Close()
		} else {
			err = rdr.Executor.Close()
		}
		rdr.Executor = nil
	}
	return
}

func castSQLTypeValue(valToCast interface{}, colType *ColumnType) (castedVal interface{}) {
	if valToCast != nil {
		if d, dok := valToCast.([]uint8); dok {
			castedVal = string(d)
		} else if sd, dok := valToCast.(string); dok {
			castedVal = sd
		} else if dtime, dok := valToCast.(time.Time); dok {
			castedVal = dtime.Format("2006-01-02T15:04:05")
		} else if djsn, djsnok := valToCast.([]byte); djsnok {
			if dv, dverr := json.Marshal(djsn); dverr == nil {
				castedVal = dv
			} else {
				castedVal = djsn
			}
		} else {
			castedVal = valToCast
		}
	} else {
		castedVal = valToCast
	}
	return castedVal
}

//ColumnType structure defining column definition
type ColumnType struct {
	name string

	hasNullable       bool
	hasLength         bool
	hasPrecisionScale bool

	nullable     bool
	length       int64
	databaseType string
	precision    int64
	scale        int64
	scanType     reflect.Type
}

//Name ColumnType.Name()
func (colType *ColumnType) Name() string {
	return colType.name
}

//Numeric ColumnType is Numeric() bool
func (colType *ColumnType) Numeric() bool {
	if colType.hasPrecisionScale {
		return true
	}
	return strings.Index(colType.databaseType, "CHAR") == -1 && strings.Index(colType.databaseType, "DATE") == -1 && strings.Index(colType.databaseType, "TIME") == -1
}

//HasNullable ColumnType content has NULL able content
func (colType *ColumnType) HasNullable() bool {
	return colType.hasNullable
}

//HasLength ColumnType content has Length definition
func (colType *ColumnType) HasLength() bool {
	return colType.hasLength
}

//HasPrecisionScale ColumnType content has PrecisionScale
func (colType *ColumnType) HasPrecisionScale() bool {
	return colType.hasPrecisionScale
}

//Nullable ColumnType content is Nullable
func (colType *ColumnType) Nullable() bool {
	return colType.nullable
}

//Length ColumnType content lenth must be used in conjunction with HasLength
func (colType *ColumnType) Length() int64 {
	return colType.length
}

//DatabaseType ColumnType underlying db type as defined by driver of Connection
func (colType *ColumnType) DatabaseType() string {
	return colType.databaseType
}

//Precision ColumnType numeric Precision. Used in conjunction with HasPrecisionScale
func (colType *ColumnType) Precision() int64 {
	return colType.precision
}

//Scale ColumnType Scale. Used in conjunction with HasPrecisionScale
func (colType *ColumnType) Scale() int64 {
	return colType.scale
}

//Type ColumnType reflect.Type as specified by golang sql/database
func (colType *ColumnType) Type() reflect.Type {
	return colType.scanType
}

//Field - struct
type Field struct {
	rdr   *Reader
	index int
}

//Value - of Field
func (fld *Field) Value() (val interface{}) {
	return
}

//Name - of Field
func (fld *Field) Name() (nme string) {
	if fld.index >= 0 {
		nme = fld.rdr.cls[fld.index]
	}
	return
}

//Type - ColumnType of Field
func (fld *Field) Type() (tpe *ColumnType) {
	if fld.index >= 0 {
		tpe = fld.rdr.cltpes[fld.index]
	}
	return
}

func newReader(exctr *Executor) (rdr *Reader) {
	rdr = &Reader{Executor: exctr, rownr: 0}
	return
}

//Repeat - repeat last query by repopulating parameters but dont regenerate last statement
func (rdr *Reader) Repeat(a ...interface{}) (err error) {
	rdr.execute()
	err = rdr.lasterr
	return
}

func (rdr *Reader) execute() (err error) {
	if rws, cltpes, cls := rdr.Executor.execute(true); rws != nil {
		if err = rdr.lasterr; err == nil {
			rdr.rownr = 0
			rdr.rws = rws
			if len(cls) > 0 {
				rdr.cls = cls[:]
				rdr.cltpes = cltpes
				invokeSuccess(rdr.script, rdr.OnSuccess, rdr)
				invokeColumns(rdr.script, rdr.OnColumns, rdr)
			}
		} else if err != nil {
			invokeError(rdr.script, err, rdr.OnError)
		}
	} else if rdr.isRemote() && rdr.jsndcdr != nil {
		if err = rdr.lasterr; err == nil {
			if len(cls) > 0 {
				rdr.cls = cls[:]
				rdr.cltpes = cltpes
				invokeSuccess(rdr.script, rdr.OnSuccess, rdr)
				invokeColumns(rdr.script, rdr.OnColumns, rdr)
			}
		}
	}
	return
}

func invokeDataValid(script active.Runtime, ondatavalid interface{}, rownr int64, rdr *Reader) (validdata bool) {
	if ondatavalid != nil {
		if fncdatavalid, fncdatavalidok := ondatavalid.(func(*Reader, int64) bool); fncdatavalidok {
			validdata = fncdatavalid(rdr, rownr)
		} else if script != nil {
			invval := script.InvokeFunction(ondatavalid, rdr, rownr)
			if isvalid, isdoneok := invval.(bool); isdoneok {
				validdata = isvalid
			} else {
				validdata = true
			}
		}
	} else {
		validdata = true
	}
	return
}

func invokeRow(script active.Runtime, onrow interface{}, rownr int64, rdr *Reader) (nextrow bool) {
	if onrow != nil {
		if fncrow, fncrowsok := onrow.(func(*Reader, int64)); fncrowsok {
			fncrow(rdr, rownr)
			nextrow = true
		} else if fncrownext, fncrowsnextok := onrow.(func(*Reader, int64) bool); fncrowsnextok {
			nextrow = !fncrownext(rdr, rownr)
		} else if script != nil {
			invval := script.InvokeFunction(onrow, rdr, rownr)
			if isdone, isdoneok := invval.(bool); isdoneok {
				if isdone {
					nextrow = false
				} else {
					nextrow = true
				}
			} else {
				nextrow = true
			}
		}
	} else {
		nextrow = true
	}
	return
}

func invokeColumns(script active.Runtime, oncolumns interface{}, rdr *Reader) {
	if oncolumns != nil {
		if fnccolumns, fnccolumnssok := oncolumns.(func(*Reader)); fnccolumnssok {
			fnccolumns(rdr)
		} else if script != nil {
			script.InvokeFunction(oncolumns, rdr)
		}
	}
}

func columnTypes(sqlcoltypes []*sql.ColumnType, cls []string) (coltypes []*ColumnType) {
	coltypes = make([]*ColumnType, len(sqlcoltypes))
	for n, ctype := range sqlcoltypes {
		coltype := &ColumnType{}
		coltype.databaseType = ctype.DatabaseTypeName()
		coltype.length, coltype.hasLength = ctype.Length()
		coltype.name = ctype.Name()
		coltype.databaseType = ctype.DatabaseTypeName()
		coltype.nullable, coltype.hasNullable = ctype.Nullable()
		coltype.precision, coltype.scale, coltype.hasPrecisionScale = ctype.DecimalSize()
		coltype.scanType = ctype.ScanType()
		coltypes[n] = coltype
	}
	return coltypes
}

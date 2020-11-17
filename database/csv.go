package database

//CSVReader -
type CSVReader struct {
	rdr        *Reader
	Headers    bool
	ColDelim   string
	RowDelim   string
	IncludeEOF bool
}

//NewCSVReader - over rdr*Reader
func NewCSVReader(rdr *Reader) (csvr *CSVReader) {

	return
}

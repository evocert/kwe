package database

//CSVReader -
type CSVReader struct {
	rdr        *Reader
	err        error
	Headers    bool
	ColDelim   string
	RowDelim   string
	IncludeEOF bool
}

//NewCSVReader - over rdr*Reader
func NewCSVReader(rdr *Reader, err error) (csvr *CSVReader) {

	return
}

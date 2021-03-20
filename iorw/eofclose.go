package iorw

import "io"

type EOFCloseReader struct {
	r  io.Reader
	rc io.Closer
}

func NewEOFCloseReader(r io.Reader) (eofclsr *EOFCloseReader) {
	if r != nil {
		eofclsr = &EOFCloseReader{r: r}
		if rc, rck := r.(io.Closer); rck {
			eofclsr.rc = rc
		}
	}
	return
}

func (eofclsr *EOFCloseReader) Read(p []byte) (n int, err error) {
	if eofclsr == nil {
		err = io.EOF
		return
	}
	if n, err = eofclsr.r.Read(p); err != nil {
		eofclsr.Close()
	}
	return
}

func (eofclsr *EOFCloseReader) Close() (err error) {
	if eofclsr != nil {
		if eofclsr.rc != nil {
			eofclsr.rc.Close()
			eofclsr.rc = nil
		}
		if eofclsr.r != nil {
			eofclsr.r = nil
		}
		eofclsr = nil
	}
	return
}

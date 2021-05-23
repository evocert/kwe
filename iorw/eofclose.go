package iorw

import (
	"bufio"
	"io"
	"strings"
)

type EOFCloseSeekReader struct {
	r    io.Reader
	rc   io.Closer
	rs   io.Seeker
	size int64
	bfr  *bufio.Reader
	//Reader Api
	canclose bool
}

func NewEOFCloseSeekReader(r io.Reader, canclose ...bool) (eofclsr *EOFCloseSeekReader) {
	if r != nil {
		eofclsr = &EOFCloseSeekReader{r: r, size: -1, canclose: len(canclose) == 0 || (len(canclose) > 0 && canclose[0])}
		if rc, rck := r.(io.Closer); rck {
			eofclsr.rc = rc
		}
		if rs, rsk := r.(io.Seeker); rsk {
			if size, skerr := rs.Seek(0, io.SeekEnd); skerr == nil && size > 0 {
				eofclsr.size = size
				rs.Seek(0, io.SeekStart)
			}
			eofclsr.rs = rs
		}
	}
	return
}

func (eofclsr *EOFCloseSeekReader) ReadRune() (r rune, size int, err error) {
	if eofclsr == nil {
		err = io.EOF
	} else {
		if eofclsr.bfr == nil && eofclsr.r != nil {
			eofclsr.bfr = bufio.NewReader(eofclsr)
			r, size, err = eofclsr.bfr.ReadRune()
		} else if eofclsr.bfr != nil {
			r, size, err = eofclsr.bfr.ReadRune()
			if err == io.EOF {
				eofclsr.Close()
			}
		} else {
			r, size, err = 0, 0, io.EOF
		}
	}
	return
}

func (eofclsr *EOFCloseSeekReader) Readln() (s string, err error) {
	rns := make([]rune, 1024)
	rnsi := 0
	for {
		rn, size, rnerr := eofclsr.ReadRune()
		if size > 0 {
			if rn == '\n' {
				if rnsi > 0 {
					s += string(rns[:rnsi])
					rnsi = 0
				}
				break
			}
			rns[rnsi] = rn
			rnsi++
			if rnsi == len(rns) {
				s += string(rns[:rnsi])
				rnsi = 0
			}
		}
		if rnerr != nil {
			err = rnerr
			if rnsi > 0 && (err == nil || err == io.EOF) {
				s += string(rns[:rnsi])
				rnsi = 0
			}
			break
		}
	}
	s = strings.TrimSpace(s)
	return
}

func (eofclsr *EOFCloseSeekReader) Readlines() (lines []string, err error) {
	for {
		ln, lnerr := eofclsr.Readln()
		if lnerr == nil {
			if ln != "" {
				if lines == nil {
					lines = []string{}
				}
				lines = append(lines, ln)
			}
		} else {
			break
		}
	}
	return
}

func (eofclsr *EOFCloseSeekReader) ReadAll() (string, error) {
	return ReaderToString(eofclsr)
}

func (eofclsr *EOFCloseSeekReader) Size() int64 {
	return eofclsr.size
}

func (eofclsr *EOFCloseSeekReader) Seek(offset int64, whence int) (n int64, err error) {
	if eofclsr != nil && eofclsr.r != nil && eofclsr.rs != nil {
		n, err = eofclsr.rs.Seek(offset, whence)
		if eofclsr.bfr != nil {
			eofclsr.bfr.Reset(eofclsr.r)
		}
	} else {
		n = -1
	}
	return
}

func (eofclsr *EOFCloseSeekReader) Read(p []byte) (n int, err error) {
	if eofclsr == nil {
		err = io.EOF
		return
	} else if eofclsr.r != nil {
		if n, err = eofclsr.r.Read(p); err != nil {
			if eofclsr.bfr == nil {
				eofclsr.Close()
			} else {
				eofclsr.disposeReader()
			}
			if n > 0 && err == io.EOF {
				err = nil
			}
		}
	}
	if n == 0 && err == nil {
		err = io.EOF
	}
	return
}

func (eofclsr *EOFCloseSeekReader) disposeReader() (err error) {
	if eofclsr != nil {
		if eofclsr.canclose {
			if eofclsr.rc != nil {
				eofclsr.rc.Close()
				eofclsr.rc = nil
			}
		}
		if eofclsr.rs != nil {
			eofclsr.rs = nil
		}
		if eofclsr.r != nil {
			eofclsr.r = nil
		}
	}
	return
}

func (eofclsr *EOFCloseSeekReader) Close() (err error) {
	if eofclsr != nil {
		eofclsr.disposeReader()
		if eofclsr.bfr != nil {
			eofclsr.bfr = nil
		}
		eofclsr = nil
	}
	return
}

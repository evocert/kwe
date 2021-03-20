package iorw

import "io"

type EOFCloseSeekReader struct {
	r    io.Reader
	rc   io.Closer
	rs   io.Seeker
	size int64
}

func NewEOFCloseSeekReader(r io.Reader) (eofclsr *EOFCloseSeekReader) {
	if r != nil {
		eofclsr = &EOFCloseSeekReader{r: r, size: -1}
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

func (eofclsr *EOFCloseSeekReader) Size() int64 {
	return eofclsr.size
}

func (eofclsr *EOFCloseSeekReader) Seek(offset int64, whence int) (n int64, err error) {
	if eofclsr != nil && eofclsr.rs != nil {
		n, err = eofclsr.rs.Seek(offset, whence)
	} else {
		n = -1
	}
	return
}

func (eofclsr *EOFCloseSeekReader) Read(p []byte) (n int, err error) {
	if eofclsr == nil {
		err = io.EOF
		return
	}
	if n, err = eofclsr.r.Read(p); err != nil {
		eofclsr.Close()
	}
	return
}

func (eofclsr *EOFCloseSeekReader) Close() (err error) {
	if eofclsr != nil {
		if eofclsr.rc != nil {
			eofclsr.rc.Close()
			eofclsr.rc = nil
		}
		if eofclsr.rs != nil {
			eofclsr.rs = nil
		}
		if eofclsr.r != nil {
			eofclsr.r = nil
		}
		eofclsr = nil
	}
	return
}

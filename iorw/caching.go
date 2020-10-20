package iorw

import (
	"bufio"
	"io"
)

//CachingReader - struct
type CachingReader struct {
	refr          io.Reader
	refs          io.Seeker
	refc          io.Closer
	rnr           io.RuneReader
	buffer        *Buffer
	maxbuffersize int64
	choffset      int64
	bufr          *BuffReader
}

//NewCachingReader - instance
func NewCachingReader(r io.Reader, maxbuffersize int64) (chngrdr *CachingReader) {
	chngrdr = &CachingReader{buffer: NewBuffer(), maxbuffersize: maxbuffersize, choffset: -1, refr: r, refc: nil, refs: nil}
	if rs, rsok := r.(io.Seeker); rsok {
		chngrdr.refs = rs
	}
	if rc, rcok := r.(io.Closer); rcok {
		chngrdr.refc = rc
	}
	chngrdr.bufr = chngrdr.buffer.Reader()
	return
}

//Read - refer io.Reader
func (chngrdr *CachingReader) Read(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		if chngrdr != nil {
			if chngrdr.bufr != nil {
				for n < pl {
					if chngrdr.choffset == -1 {
						if sn, serr := chngrdr.Seek(chngrdr.choffset+1, io.SeekStart); serr == nil {
							if sn >= 0 {
								chngrdr.choffset = sn
							} else {
								break
							}
						} else {
							break
						}
					} else {
						if sn, serr := chngrdr.Seek(chngrdr.choffset+1, io.SeekStart); serr == nil {
							if sn >= 0 {
								chngrdr.choffset = sn
							} else {
								break
							}
						} else {
							break
						}
					}
					rn, rnerr := chngrdr.bufr.Read(p[n : n+(pl-n)])
					if rn > 0 {
						var rni = 0
						if cl := (rn - rni); cl >= (pl - n) {
							rni += cl
							chngrdr.choffset += int64(rni)
							n += rni
						} else if cl := (pl - n); cl < (rn - rni) {
							rni += cl
							chngrdr.choffset += int64(rni)
							n += rni
						}
					}
					if rnerr != nil {
						if rnerr != io.EOF {
							break
						}
					}
				}
			}
		}
		if n == 0 && err == nil {
			err = io.EOF
		}
	}
	return
}

//ReadRune - refer io.RuneReader
func (chngrdr *CachingReader) ReadRune() (r rune, size int, err error) {
	if chngrdr != nil {
		if chngrdr.rnr == nil {
			chngrdr.rnr = bufio.NewReader(chngrdr)
		}
		r, size, err = chngrdr.rnr.ReadRune()
	} else {
		err = io.EOF
	}
	return
}

//Close - refer io.Closer
func (chngrdr *CachingReader) Close() (err error) {
	if chngrdr != nil {
		if chngrdr.bufr != nil {
			chngrdr.bufr.Close()
			chngrdr.bufr = nil
		}
		if chngrdr.buffer != nil {
			chngrdr.buffer.Close()
			chngrdr.buffer = nil
		}
		if chngrdr.refc != nil {
			chngrdr.refc.Close()
			chngrdr.refc = nil
		}
		if chngrdr.refs != nil {
			chngrdr.refs = nil
		}
		if chngrdr.refr != nil {
			chngrdr.refr = nil
		}
	}
	return
}

//WriteTo - helper for io.Copy
func (chngrdr *CachingReader) WriteTo(w io.Writer) (n int64, err error) {
	if w != nil && chngrdr != nil {
	}
	return
}

//Seek - refer to io.Seeker
func (chngrdr *CachingReader) Seek(offset int64, whence int) (n int64, err error) {
	if chngrdr != nil {
		if chngrdr.refs != nil {
			if n, err = chngrdr.refs.Seek(offset, whence); err == nil && n >= 0 {
				if chngrdr.choffset != n {
					chngrdr.choffset = n
					chngrdr.buffer.Clear()
				}
			}
		} else {
			n = -1
		}
	} else {
		n = -1
	}
	return
}

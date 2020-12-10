package iorw

import (
	"bufio"
	"io"
	"sync"
)

//Buffer -
type Buffer struct {
	buffer [][]byte
	bytes  []byte
	bytesi int
	lck    *sync.RWMutex
	//bufrs   map[*BuffReader]*BuffReader
	OnClose func(*Buffer)
}

//NewBuffer -
func NewBuffer() (buff *Buffer) {
	buff = &Buffer{lck: &sync.RWMutex{}, buffer: [][]byte{}, bytesi: 0, bytes: make([]byte, 8192) /*bufrs: map[*BuffReader]*BuffReader{}*/}
	return
}

//BuffersLen - return len() of internal byte[][] buffer
func (buff *Buffer) BuffersLen() (s int) {
	return len(buff.buffer)
}

//Print - same as fmt.Print just on buffer
func (buff *Buffer) Print(a ...interface{}) {
	Fprint(buff, a...)
}

//Println - same as fmt.Println just on buffer
func (buff *Buffer) Println(a ...interface{}) {
	Fprintln(buff, a...)
}

//String - return buffer as string value
func (buff *Buffer) String() (s string) {
	s = ""
	if buff != nil {
		if buf := buff.buffer; len(buf) > 0 {
			for _, b := range buf {
				s += string(b)
			}
		}
		if buff.bytesi > 0 {
			s += string(buff.bytes[:buff.bytesi])
		}
	}
	return
}

//Size - total size of Buffer
func (buff *Buffer) Size() (s int64) {
	s = 0
	if len(buff.buffer) > 0 {
		s += (int64(len(buff.buffer)) * int64(len(buff.buffer[0])))
	}
	if buff.bytesi > 0 {
		s += int64(buff.bytesi)
	}
	return s
}

//ReadFrom - fere io.ReaderFrom
func (buff *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	if r != nil {
		var p = make([]byte, 4096)
		for {
			pn, pnerr := r.Read(p)
			if pn > 0 {
				n += int64(pn)
				var pi = 0
				for pi < pn {
					wn, wnerr := buff.Write(p[pi : pi+(pn-pi)])
					if wn > 0 {
						pi += wn
					}
					if wnerr != nil {
						pnerr = wnerr
						break
					}
					if wn == 0 {
						break
					}
				}
			}
			if pnerr != nil {
				err = pnerr
				break
			} else {
				if pn == 0 {
					err = io.EOF
					break
				}
			}
		}
	}
	return
}

//WriteRune - Write singe rune
func (buff *Buffer) WriteRune(r rune) (err error) {
	_, err = buff.Write([]byte(string(r)))
	return
}

//WriteRunes - Write runes
func (buff *Buffer) WriteRunes(p []rune) (err error) {
	if pl := len(p); pl > 0 {
		_, err = buff.Write([]byte(string(p[:pl])))
	}
	return
}

//Write - refer io.Writer
func (buff *Buffer) Write(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		func() {
			buff.lck.Lock()
			defer buff.lck.Unlock()
			for n < pl {
				if tl := (len(buff.bytes) - buff.bytesi); (pl - n) >= tl {
					if cl := copy(buff.bytes[buff.bytesi:buff.bytesi+tl], p[n:n+tl]); cl > 0 {
						n += cl
						buff.bytesi += cl
					}
				} else if tl := (pl - n); tl < (len(buff.bytes) - buff.bytesi) {
					if cl := copy(buff.bytes[buff.bytesi:buff.bytesi+tl], p[n:n+tl]); cl > 0 {
						n += cl
						buff.bytesi += cl
					}
				}
				if buff.bytesi == len(buff.bytes) {
					if buff.buffer == nil {
						buff.buffer = [][]byte{}
					}
					var bts = make([]byte, buff.bytesi)
					copy(bts, buff.bytes[:buff.bytesi])
					buff.buffer = append(buff.buffer, bts)
					buff.bytesi = 0
				}
			}
		}()
	}
	return
}

//Reader -
func (buff *Buffer) Reader() (bufr *BuffReader) {
	bufr = &BuffReader{buffer: buff, roffset: -1}
	//buff.bufrs[bufr] = bufr
	return
}

//Close - refer io.Closer
func (buff *Buffer) Close() (err error) {
	if buff != nil {
		if buff.lck != nil {
			if buff.OnClose != nil {
				buff.OnClose(buff)
				buff.OnClose = nil
			}
			buff.Clear()
			buff.lck = nil
		}
		buff = nil
	}
	return
}

//Clear - Buffer
func (buff *Buffer) Clear() (err error) {
	if buff != nil {
		if buff.lck != nil {
			func() {
				buff.lck.Lock()
				defer buff.lck.Unlock()

				/*if buff.bufrs != nil {
					if len(buff.bufrs) > 0 {
						var bufrs = make([]*BuffReader, len(buff.bufrs))
						var bufrsi = 0
						for bufrsk := range buff.bufrs {
							buff.bufrs[bufrsk] = nil
							bufrs[bufrsi] = bufrsk
							bufrsk.Close()
							bufrsi++
						}
						for _, bufrsk := range bufrs {
							delete(buff.bufrs, bufrsk)
						}
						bufrs = nil
					}
					buff.bufrs = nil
				}*/
				if buff.buffer != nil {
					for len(buff.buffer) > 0 {
						buff.buffer[0] = nil
						buff.buffer = buff.buffer[1:]
					}
					buff.buffer = nil
				}
				if buff.bytesi > 0 {
					buff.bytesi = 0
				}
			}()
		}
	}
	return
}

//BuffReader -
type BuffReader struct {
	buffer   *Buffer
	rnr      *bufio.Reader
	roffset  int64
	rbufferi int
	rbytes   []byte
	rbytesi  int
}

//WriteTo - helper for io.Copy
func (bufr *BuffReader) WriteTo(w io.Writer) (n int64, err error) {
	if w != nil && bufr != nil {
		var r = io.Reader(bufr)
		if bufr.rnr != nil {
			r = bufr.rnr
		}
		var p = make([]byte, 4096)
		for {
			pn, pnerr := r.Read(p)
			if pn > 0 {
				n += int64(pn)
				var pi = 0
				for pi < pn {
					wn, wnerr := w.Write(p[pi : pi+(pn-pi)])
					if wn > 0 {
						pi += wn
					}
					if wnerr != nil {
						pnerr = wnerr
						break
					}
					if wn == 0 {
						break
					}
				}
			}
			if pnerr == nil {
				if pn == 0 {
					pnerr = io.EOF
					break
				}
			}
			if pnerr != nil {
				err = pnerr
				break
			}
		}

	}
	return
}

//Close - refer io.Closer
func (bufr *BuffReader) Close() (err error) {
	if bufr != nil {
		if bufr.buffer != nil {
			/*func() {
				if _, ok := bufr.buffer.bufrs[bufr]; ok {
					bufr.buffer.bufrs[bufr] = nil
					delete(bufr.buffer.bufrs, bufr)
				}
			}()*/
			bufr.buffer = nil
		}
		if bufr.rnr != nil {
			bufr.rnr = nil
		}
		if bufr.rbytes != nil {
			bufr.rbytes = nil
		}
	}
	return
}

//Read - refer io.Reader
func (bufr *BuffReader) Read(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		if bufr != nil {
			for n < pl {
				if len(bufr.rbytes) == 0 || (len(bufr.rbytes) > 0 && len(bufr.rbytes) == bufr.rbytesi) {
					if bufr.roffset == -1 {
						offn, offnerr := bufr.Seek(0, io.SeekStart)
						if offnerr == nil && offn >= 0 {
							bufr.roffset = offn
						} else {
							err = offnerr
							break
						}
					} else {
						if bufr.roffset == bufr.buffer.Size() {
							break
						}
						offn, offnerr := bufr.Seek(bufr.roffset, io.SeekStart)
						if offnerr == nil && offn >= 0 {
							bufr.roffset = offn
						} else {
							err = offnerr
							break
						}
					}
				}
				for (pl > n) && (len(bufr.rbytes) > bufr.rbytesi) {
					if cl := (len(bufr.rbytes) - bufr.rbytesi); (pl - n) >= cl {
						copy(p[n:n+cl], bufr.rbytes[bufr.rbytesi:bufr.rbytesi+cl])
						n += cl
						bufr.roffset += int64(cl)
						bufr.rbytesi += cl
					} else if cl := (pl - n); cl < (len(bufr.rbytes) - bufr.rbytesi) {
						copy(p[n:n+cl], bufr.rbytes[bufr.rbytesi:bufr.rbytesi+cl])
						n += cl
						bufr.roffset += int64(cl)
						bufr.rbytesi += cl
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
func (bufr *BuffReader) ReadRune() (r rune, size int, err error) {
	if bufr != nil {
		if bufr.rnr == nil {
			bufr.rnr = bufio.NewReader(bufr)
		}
		r, size, err = bufr.rnr.ReadRune()
	} else {
		err = io.EOF
	}
	return
}

//Seek - refer to io.Seeker
func (bufr *BuffReader) Seek(offset int64, whence int) (n int64, err error) {
	if bufr != nil && bufr.buffer != nil {
		var adjusted = false
		if bufs := bufr.buffer.Size(); bufs > 0 {
			func() {
				bufr.buffer.lck.RLock()
				defer bufr.buffer.lck.RUnlock()
				var adjustOffsetRead = func() {
					var rnbufi = 0
					var bufl = bufr.buffer.BuffersLen()
					var bflen = 0
					var bufbfs = int64(0)
					if bufl > 0 {
						bflen = len(bufr.buffer.buffer[0])
						bufbfs = (int64(bufl) * int64(bflen))
					}
					if n < (int64(bufl) * int64(bflen)) {
						if n < int64(bflen) {
							rnbufi = 0
						} else {
							rnbufi = int(n / int64(bflen))
						}
						bufr.rbytesi = int(n % int64(bflen))
						bufr.rbytes = bufr.buffer.buffer[rnbufi]
						bufr.rbufferi = rnbufi
					} else if n < bufs {
						bufr.rbufferi = rnbufi
						bufr.rbytesi = int(n % (bufs - bufbfs))
						bufr.rbytes = bufr.buffer.bytes[:bufr.buffer.bytesi]
					}
					bufr.roffset = n
					adjusted = true
				}

				if whence == io.SeekStart {
					if offset >= 0 && offset < bufs {
						n = offset
						adjustOffsetRead()
					}
				} else if whence == io.SeekCurrent {
					if bufr.roffset >= -1 && (bufr.roffset+offset) < bufs {
						if bufr.roffset == -1 {
							n = bufr.roffset + 1 + offset
						} else {
							n = bufr.roffset + offset
						}
						adjustOffsetRead()
					}
				} else if whence == io.SeekEnd {
					if (bufs-offset) >= 0 && (bufs-offset) < bufs {
						n = (bufs - offset)
						adjustOffsetRead()
					}
				}
			}()
		}
		if !adjusted {
			n = -1
		}
	} else {
		n = -1
	}
	return
}

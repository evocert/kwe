package iorw

import (
	"bufio"
	"io"
	"sync"
)

//Buffer -
type Buffer struct {
	buffer  [][]byte
	bytes   []byte
	bytesi  int
	lck     *sync.RWMutex
	bufrs   map[*BuffReader]*BuffReader
	OnClose func(*Buffer)
}

//NewBuffer -
func NewBuffer() (buff *Buffer) {
	buff = &Buffer{lck: &sync.RWMutex{}, buffer: [][]byte{}, bytesi: 0, bytes: make([]byte, 8192), bufrs: map[*BuffReader]*BuffReader{}}
	//runtime.SetFinalizer(buff, bufferFinalize)
	return
}

//BuffersLen - return len() of internal byte[][] buffer
func (buff *Buffer) BuffersLen() (s int) {
	return len(buff.buffer)
}

//Clone - return *Buffer clone
func (buff *Buffer) Clone() (clnbf *Buffer) {
	clnbf = NewBuffer()
	if buff.Size() > 0 {
		if len(buff.buffer) > 0 {
			if clnbf.buffer == nil {
				clnbf.buffer = [][]byte{}
			}
			clnbf.buffer = append(clnbf.buffer, buff.buffer...)
		}
		if buff.bytesi > 0 {
			copy(clnbf.bytes, buff.bytes[:buff.bytesi])
			clnbf.bytesi = buff.bytesi
		}
	}
	return
}

//Print - same as fmt.Print just on buffer
func (buff *Buffer) Print(a ...interface{}) {
	Fprint(buff, a...)
}

//Println - same as fmt.Println just on buffer
func (buff *Buffer) Println(a ...interface{}) {
	Fprintln(buff, a...)
}

//SubString - return buffer as string value based on offset ...int64
func (buff *Buffer) SubString(offset ...int64) (s string) {
	if buff != nil {
		if len(offset) > 0 && len(offset)%2 == 0 {
			if sl := buff.Size(); sl > 0 {
				var bufr *BuffReader = nil
				rns := make([]rune, 1024)
				rnsi := 0
				busy := true
				for len(offset) > 0 && busy {
					if offset[0] <= sl && offset[1] < sl {
						if bufr == nil {
							bufr = buff.Reader()
						}
						bufr.Seek(offset[0], 0)
						for {
							r, rs, rerr := bufr.ReadRune()
							if rs > 0 {
								rns[rnsi] = r
								rnsi++
								if rnsi == len(rns) {
									rnsi = 0
									s += string(rns[:])
								}
							}
							if rerr != nil {
								busy = false
								break
							}
							offset[0]++
							if offset[0] >= offset[1] {
								break
							}
						}
						if busy {
							offset = offset[2:]
						}
					} else {
						break
					}
				}
				if bufr != nil {
					bufr.Close()
					bufr = nil
				}
			}
		}
	}
	return
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
	err = buff.WriteRunes(r)
	return
}

//WriteRunes - Write runes
func (buff *Buffer) WriteRunes(p ...rune) (err error) {
	if pl := len(p); pl > 0 {
		if bs := RunesToUTF8(p[:pl]); len(bs) > 0 {
			//_, err = buff.Write([]byte(string(p[:pl])))
			_, err = buff.Write(bs)
		}
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
	if buff != nil {
		bufr = &BuffReader{buffer: buff, roffset: -1, MaxRead: -1}
		//runtime.SetFinalizer(bufr, buffReaderFinalize)
	}
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
				if buff.bufrs != nil {
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
				}
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
	MaxRead  int64
	roffset  int64
	rbufferi int
	rbytes   []byte
	rbytesi  int
}

//SetMaxRead - set max read implementation for Reader interface compliance
func (bufr *BuffReader) SetMaxRead(maxlen int64) (err error) {
	if bufr != nil {
		if maxlen < 0 {
			maxlen = -1
		}
		bufr.MaxRead = maxlen
	}
	return
}

func (bufr *BuffReader) WriteToFunc(funcw func([]byte) (int, error)) (n int64, err error) {
	if bufr != nil && funcw != nil {
		n, err = WriteToFunc(bufr, funcw)
	}
	return
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
			func() {
				if _, ok := bufr.buffer.bufrs[bufr]; ok {
					bufr.buffer.bufrs[bufr] = nil
					delete(bufr.buffer.bufrs, bufr)
				}
			}()
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

//RuneAt - rune at offset int64
func (bufr *BuffReader) RuneAt(offset int64) (rn rune) {
	rn = -1
	if s := bufr.SubString(offset, offset); s != "" {
		rn = rune(s[0])
	}
	return
}

//LastIndex - Last index of s string - n int64
func (bufr *BuffReader) LastIndex(s string, offset ...int64) int64 {
	if bufr == nil || s == "" {
		return -1
	}
	if len(offset) == 2 {
		return bufr.LastByteIndexWithinOffsets(offset[0], offset[1], []byte(s)...)
	} else if len(offset) == 1 {
		return bufr.LastByteIndexWithinOffsets(-1, offset[0], []byte(s)...)
	}
	return bufr.LastByteIndexWithinOffsets(-1, -1, []byte(s)...)
}

//LastByteIndexWithinOffsets - Last index of bs byte... - n int64 within startoffset and endoffset
func (bufr *BuffReader) LastByteIndexWithinOffsets(startoffset, endoffset int64, bs ...byte) (index int64) {
	index = -1
	if bufr != nil && bufr.buffer != nil && len(bs) > 0 {
		if ls := bufr.buffer.Size(); ls > 0 {
			for i, j := 0, len(bs)-1; i < j; i, j = i+1, j-1 {
				bs[i], bs[j] = bs[j], bs[i]
			}
			prvb := byte(0)
			bsi := 0
			toffset := int64(0)
			if bufr.buffer.bytesi > 0 {
				bti := bufr.buffer.bytesi - 1
				for bti > -1 {
					toffset++
					bt := bufr.buffer.bytes[bti]
					bti--
					if bsi > 0 && bs[bsi-1] == prvb && bs[bsi] != bt {
						bsi = 0
						prvb = byte(0)
					}
					if bs[bsi] == bt {
						bsi++
						if bsi == len(bs) {
							toffset += int64(len(bs))
							index = bufr.buffer.Size() - toffset
							break
						} else {
							prvb = bt
						}
					} else {
						if bsi > 0 {
							bsi = 0
						}
					}
				}
			}
			if index == -1 && len(bufr.buffer.buffer) > 0 {
				bfi := len(bufr.buffer.buffer) - 1
				for bfi > -1 {
					toffset++
					bf := bufr.buffer.buffer[bfi]
					bti := len(bf) - 1
					for bti > -1 {
						bt := bufr.buffer.bytes[bti]
						bti--
						if bsi > 0 && bs[bsi-1] == prvb && bs[bsi] != bt {
							bsi = 0
							prvb = byte(0)
						}
						if bs[bsi] == bt {
							bsi++
							if bsi == len(bs) {
								toffset += int64(len(bs))
								index = bufr.buffer.Size() - toffset
								break
							} else {
								prvb = bt
							}
						} else {
							if bsi > 0 {
								bsi = 0
							}
						}
					}
					if index > -1 {
						break
					}
					bfi--
				}
			}
		}
	}
	return
}

//Index - Index of s string - n int64
func (bufr *BuffReader) Index(s string) int64 {
	if bufr == nil || s == "" {
		return -1
	}
	return bufr.ByteIndex([]byte(s)...)
}

//ByteIndex - Index of bs ...byte - n int64
func (bufr *BuffReader) ByteIndex(bs ...byte) (index int64) {
	index = -1
	if bufr != nil && bufr.buffer != nil && len(bs) > 0 {
		prvb := byte(0)
		bsi := 0
		toffset := int64(-1)
		if len(bufr.buffer.buffer) > 0 {
			for _, bf := range bufr.buffer.buffer {
				for _, bt := range bf {
					toffset++
					if bsi > 0 && bs[bsi-1] == prvb && bs[bsi] != bt {
						bsi = 0
						prvb = byte(0)
					}
					if bs[bsi] == bt {
						bsi++
						if bsi == len(bs) {
							index = toffset - int64(len(bs))
							break
						} else {
							prvb = bt
						}
					} else {
						if bsi > 0 {
							bsi = 0
						}
					}
				}
				if index > -1 {
					break
				}
			}
		}
		if index == -1 && bufr.buffer.bytesi > 0 {
			for _, bt := range bufr.buffer.bytes[:bufr.buffer.bytesi] {
				toffset++
				if bsi > 0 && bs[bsi-1] == prvb && bs[bsi] != bt {
					bsi = 0
					prvb = byte(0)
				}
				if bs[bsi] == bt {
					bsi++
					if bsi == len(bs) {
						index = toffset - int64(len(bs))
						break
					} else {
						prvb = bt
					}
				} else {
					if bsi > 0 {
						bsi = 0
					}
				}
			}
		}
	}
	return
}

//Read - refer io.Reader
func (bufr *BuffReader) Read(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		rl := 0
		if bufr != nil {
			for n < pl && (bufr.MaxRead > 0 || bufr.MaxRead == -1) {

				if len(bufr.rbytes) == 0 || (len(bufr.rbytes) > 0 && len(bufr.rbytes) == bufr.rbytesi) {
					if bufr.roffset == -1 {
						if offn, offnerr := bufr.Seek(0, io.SeekStart); offnerr != nil || offn == -1 {
							err = offnerr
							break
						}
						/*if offnerr == nil && offn >= 0 {
							bufr.roffset = offn
						} else {
							err = offnerr
							break
						}*/
					} else {
						if bufr.roffset == bufr.buffer.Size() {
							break
						}
						if offn, offnerr := bufr.Seek(bufr.roffset, io.SeekStart); offnerr != nil || offn == -1 {
							err = offnerr
							break
						}
						/*if offnerr == nil && offn >= 0 {
							bufr.roffset = offn
						} else {
							err = offnerr
							break
						}*/
					}
				}

				for (bufr.MaxRead > 0 || bufr.MaxRead == -1) && (pl > n) && (len(bufr.rbytes) > bufr.rbytesi) {
					rbtsl := len(bufr.rbytes)
					if bufr.MaxRead > 0 {
						if ln := int64(rbtsl - bufr.rbytesi); ln > bufr.MaxRead {
							rl = int(bufr.MaxRead)
						} else {
							rl = int(ln)
						}
						if (rl + bufr.rbytesi) < rbtsl {
							rbtsl = (rl + bufr.rbytesi)
						}
					}
					if cl := (rbtsl - bufr.rbytesi); (pl - n) >= cl {
						copy(p[n:n+cl], bufr.rbytes[bufr.rbytesi:bufr.rbytesi+cl])
						n += cl
						bufr.roffset += int64(cl)
						bufr.rbytesi += cl
						if bufr.MaxRead > 0 {
							bufr.MaxRead -= int64(cl)
							if bufr.MaxRead < 0 {
								bufr.MaxRead = 0
							}
						}
					} else if cl := (pl - n); cl < (rbtsl - bufr.rbytesi) {
						copy(p[n:n+cl], bufr.rbytes[bufr.rbytesi:bufr.rbytesi+cl])
						n += cl
						bufr.roffset += int64(cl)
						bufr.rbytesi += cl
						if bufr.MaxRead > 0 {
							bufr.MaxRead -= int64(cl)
							if bufr.MaxRead < 0 {
								bufr.MaxRead = 0
							}
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

//SubString - return buffer as string value based on offset ...int64
func (bufr *BuffReader) SubString(offset ...int64) (s string) {
	if bufr != nil && bufr.buffer != nil {
		if len(offset) > 0 && len(offset)%2 == 0 {
			if sl := bufr.buffer.Size(); sl > 0 {
				if offset[0] == -1 {
					offset[0] = 0
				}
				if offset[1] == -1 {
					offset[1] = sl - 1
				}
				rns := make([]rune, 1024)
				rnsi := 0
				busy := true
				for len(offset) > 0 && busy {
					if offset[0] <= sl && offset[1] < sl {
						bufr.Seek(offset[0], 0)
						for {
							r, rs, rerr := bufr.ReadRune()
							if rs > 0 {
								rns[rnsi] = r
								rnsi++
								if rnsi == len(rns) {
									rnsi = 0
									s += string(rns[:])
								}
							}
							if rerr != nil {
								busy = false
								break
							}
							offset[0]++
							if offset[0] >= offset[1] {
								busy = false
								break
							}
						}
						if busy {
							offset = offset[2:]
						}
					} else {
						break
					}
				}
				if rnsi > 0 {
					s += string(rns[:rnsi])
				}
			}
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

func (bufr *BuffReader) Readln() (ln string, err error) {
	ln, err = ReadLine(bufr)
	return
}

func (bufr *BuffReader) Readlines() (lines []string, err error) {
	for {
		ln, lnerr := bufr.Readln()
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

func (bufr *BuffReader) ReadAll() (string, error) {
	return bufr.buffer.String(), nil
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
					if bufl > 0 && n < (int64(bufl)*int64(bflen)) {
						if n < int64(bflen) {
							rnbufi = 0
						} else {
							rnbufi = int(n / int64(bflen))
						}
						bufr.rbytesi = int(n % int64(bflen))
						bufr.rbytes = bufr.buffer.buffer[rnbufi]
						bufr.rbufferi = rnbufi
					} else if n < bufs {
						if bflen > 0 {
							if n < int64(bflen) {
								rnbufi = 0
							}
							if n == (int64(bufl) * int64(bflen)) {
								bufr.rbytesi = 0
							} else {
								bufr.rbytesi = int(n % (bufs - bufbfs))
							}
						} else {
							bufr.rbufferi = rnbufi
							bufr.rbytesi = int(n % (bufs - bufbfs))
						}
						bufr.rbytes = bufr.buffer.bytes[:bufr.buffer.bytesi]
					}
					adjusted = true
					bufr.roffset = n
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
					if (bufs-offset) >= 0 && (bufs-offset) <= bufs {
						if (bufs - offset) < bufs {
							n = (bufs - offset)
						} else {
							n = (bufs - offset)
						}
						adjustOffsetRead()
					}
				}
			}()
		}
		if !adjusted {
			n = -1
		} else {
			if bufr.rnr != nil {
				bufr.rnr.Reset(bufr)
			}
		}
	} else {
		n = -1
	}
	return
}

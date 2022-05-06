package iorw

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

//Printer - interface
type Printer interface {
	Print(a ...interface{}) error
	Println(a ...interface{}) error
	Write(p []byte) (int, error)
}

//Reader - interface
type Reader interface {
	Seek(int64, int) (int64, error)
	SetMaxRead(int64) (err error)
	Read([]byte) (int, error)
	ReadRune() (rune, int, error)
	Readln() (string, error)
	Readlines() ([]string, error)
	ReadAll() (string, error)
}

//PrinterReader - interface
type PrinterReader interface {
	Printer
	Reader
}

//Fprint - refer to fmt.Fprint
func Fprint(w io.Writer, a ...interface{}) (err error) {
	if len(a) > 0 && w != nil {
		for dn := range a {
			if s, sok := a[dn].(string); sok {
				if _, err = w.Write([]byte(s)); err != nil {
					break
				}
			} else if ir, irok := a[dn].(io.Reader); irok {
				if _, err = WriteToFunc(ir, func(b []byte) (int, error) {
					return w.Write(b)
				}); err != nil {
					break
				}
			} else if aa, aaok := a[dn].([]interface{}); aaok {
				if len(aa) > 0 {
					if err = Fprint(w, aa...); err != nil {
						break
					}
				}
			} else {
				if _, err = fmt.Fprint(w, a[dn]); err != nil {
					break
				}
			}
		}
	}
	return
}

func CopyBytes(dest []byte, desti int, src []byte, srci int) (lencopied int, destn int, srcn int) {
	if destl, srcl := len(dest), len(src); (destl > 0 && desti < destl) && (srcl > 0 && srci < srcl) {
		if (srcl - srci) <= (destl - desti) {
			cpyl := copy(dest[desti:desti+(srcl-srci)], src[srci:srci+(srcl-srci)])
			srcn = srci + cpyl
			destn = desti + cpyl
			lencopied = cpyl
		} else if (destl - desti) < (srcl - srci) {
			cpyl := copy(dest[desti:desti+(destl-desti)], src[srci:srci+(destl-desti)])
			srcn = srci + cpyl
			destn = desti + cpyl
			lencopied = cpyl
		}
	}
	return
}

//Fprintln - refer to fmt.Fprintln
func Fprintln(w io.Writer, a ...interface{}) (err error) {
	if len(a) > 0 && w != nil {
		err = Fprint(w, a...)
	}
	if err == nil {
		err = Fprint(w, "\r\n")
	}
	return
}

//ReadLines from r io.Reader as lines []string
func ReadLines(r io.Reader) (lines []string, err error) {
	if r != nil {
		var rnrd io.RuneReader = nil
		if rnr, rnrok := r.(io.RuneReader); rnrok {
			rnrd = rnr
		} else {
			rnrd = bufio.NewReader(r)
		}
		rns := make([]rune, 1024)
		rnsi := 0
		s := ""
		for {
			rn, size, rnerr := rnrd.ReadRune()
			if size > 0 {
				if rn == '\n' {
					if rnsi > 0 {
						s += string(rns[:rnsi])
						rnsi = 0
					}
					if s != "" {
						s = strings.TrimSpace(s)
						if lines == nil {
							lines = []string{}
						}
						lines = append(lines, s)
						s = ""
					}
					continue
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
				if rnsi > 0 {
					s += string(rns[:rnsi])
					rnsi = 0
				}
				if s != "" {
					s = strings.TrimSpace(s)
					if lines == nil {
						lines = []string{}
					}
					lines = append(lines, s)
					s = ""
				}
				if err == io.EOF {
					err = nil
				}
				break
			}
		}
	}
	return
}

//ReadLine from r io.Reader as s string
func ReadLine(r io.Reader) (s string, err error) {
	if r != nil {
		var rnrd io.RuneReader = nil
		if rnr, rnrok := r.(io.RuneReader); rnrok {
			rnrd = rnr
		} else {
			rnrd = bufio.NewReader(r)
		}
		rns := make([]rune, 1024)
		rnsi := 0
		for {
			rn, size, rnerr := rnrd.ReadRune()
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
					if err == io.EOF {
						err = nil
					}
					s += string(rns[:rnsi])
					rnsi = 0
				}
				break
			}
		}
	}
	s = strings.TrimSpace(s)
	return
}

//ReaderToString read reader and return content as string
func ReaderToString(r interface{}) (s string, err error) {
	runes := make([]rune, 1024)
	runesi := 0
	if err = ReadRunesEOFFunc(r, func(rn rune) error {
		runes[runesi] = rn
		runesi++
		if runesi == len(runes) {
			s += string(runes[:runesi])
			runesi = 0
		}
		return nil
	}); err == nil || err == io.EOF {
		if runesi > 0 {
			s += string(runes[:runesi])
			runesi = 0
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

//ReadRunesEOFFunc read runes from r io.Reader and call fncrne func(rune) error
func ReadRunesEOFFunc(r interface{}, fncrne func(rune) error) (err error) {
	if r != nil && fncrne != nil {
		var rnrd io.RuneReader = nil
		if rnr, rnrok := r.(io.RuneReader); rnrok {
			rnrd = rnr
		} else if rdr, rdrok := r.(io.Reader); rdrok {
			rnrd = bufio.NewReader(rdr)
		}
		if rnrd != nil {
			for {
				rn, size, rnerr := rnrd.ReadRune()
				if size > 0 {
					if err = fncrne(rn); err != nil {
						break
					}
				}
				if err == nil && rnerr != nil {
					if rnerr != io.EOF {
						err = rnerr
					}
					break
				}
			}
		}
	}
	return
}

func RunesToUTF8(rs []rune) []byte {
	size := 0
	for rn := range rs {
		size += utf8.RuneLen(rs[rn])
	}
	bs := make([]byte, size)
	count := 0
	for rn := range rs {
		count += utf8.EncodeRune(bs[count:], rs[rn])
	}

	return bs
}

type funcrdrwtr struct {
	funcw func([]byte) (int, error)
	funcr func([]byte) (int, error)
}

func (fncrw *funcrdrwtr) Close() (err error) {
	if fncrw != nil {
		if fncrw.funcr != nil {
			fncrw.funcr = nil
		}
		if fncrw.funcw != nil {
			fncrw.funcw = nil
		}
		fncrw = nil
	}
	return
}

func (fncrw *funcrdrwtr) Write(p []byte) (n int, err error) {
	if fncrw != nil && fncrw.funcw != nil {
		n, err = fncrw.funcw(p)
	}
	return
}

func (fncrw *funcrdrwtr) Read(p []byte) (n int, err error) {
	if fncrw != nil && fncrw.funcr != nil {
		n, err = fncrw.funcr(p)
	}
	return
}

func WriteToFunc(r io.Reader, funcw func([]byte) (int, error)) (n int64, err error) {
	if r != nil && funcw != nil {
		func() {
			n, err = ReadWriteToFunc(funcw, func(b []byte) (int, error) {
				return r.Read(b)
			})
		}()
	}
	return
}

func ReadToFunc(w io.Writer, funcr func([]byte) (int, error)) (n int64, err error) {
	if w != nil && funcr != nil {
		func() {
			n, err = ReadWriteToFunc(func(b []byte) (int, error) {
				return w.Write(b)
			}, funcr)
		}()
	}
	return
}

func ReadHandle(r io.Reader, handle func([]byte), maxrlen int) (n int, err error) {
	if maxrlen < 4096 {
		maxrlen = 4096
	}
	s := make([]byte, maxrlen)
	sn := 0
	si := 0
	sl := len(s)
	serr := error(nil)
	for n < maxrlen && err == nil {
		switch sn, serr = r.Read(s[si : si+(sl-si)]); true {
		case sn < 0:
			err = serr
			break
		case sn == 0: // EOF
			if si > 0 {
				handle(s[:si])
				si = 0
			}
			err = serr
			break
		case sn > 0:
			si += sn
			n += sn
			err = serr
		}
	}
	if si > 0 {
		handle(s[:si])
	}
	if n == 0 && err == nil {
		err = io.EOF
	}
	return
}

func ReadWriteToFunc(funcw func([]byte) (int, error), funcr func([]byte) (int, error)) (n int64, err error) {
	if funcw != nil && funcr != nil {
		fncrw := &funcrdrwtr{funcr: funcr, funcw: funcw}
		func() {
			defer func() {
				if rv := recover(); rv != nil {
					switch x := rv.(type) {
					case string:
						err = errors.New(x)
					case error:
						err = x
					default:
						err = errors.New("unknown panic")
					}
				}
				fncrw.Close()
			}()
			n, err = io.Copy(fncrw, fncrw)
		}()
	}
	return
}

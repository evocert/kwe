package iorw

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

//Printer - interface
type Printer interface {
	Print(a ...interface{})
	Println(a ...interface{})
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
func Fprint(w io.Writer, a ...interface{}) {
	if len(a) > 0 && w != nil {
		for _, d := range a {
			if s, sok := d.(string); sok {
				w.Write([]byte(s))
			} else if ir, irok := d.(io.Reader); irok {
				io.Copy(w, ir)
			} else if aa, aaok := d.([]interface{}); aaok {
				if len(aa) > 0 {
					Fprint(w, aa...)
				}
			} else {
				fmt.Fprint(w, d)
			}
		}
	}
}

//Fprintln - refer to fmt.Fprintln
func Fprintln(w io.Writer, a ...interface{}) {
	if len(a) > 0 && w != nil {
		Fprint(w, a...)
	}
	Fprint(w, "\r\n")
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
	for _, r := range rs {
		size += utf8.RuneLen(r)
	}
	bs := make([]byte, size)
	count := 0
	for _, r := range rs {
		count += utf8.EncodeRune(bs[count:], r)
	}

	return bs
}

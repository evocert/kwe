package iorw

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

//Printer - interface
type Printer interface {
	Print(a ...interface{})
	Println(a ...interface{})
}

//Reader - interface
type Reader interface {
	Read([]byte) (int, error)
	ReadRune() (rune, int, error)
	Readln() (string, error)
	Readlines() ([]string, error)
	ReadAll() (string, error)
}

//Fprint - refer to fmt.Fprint
func Fprint(w io.Writer, a ...interface{}) {
	if len(a) > 0 && w != nil {
		for _, d := range a {
			if s, sok := d.(string); sok {
				w.Write([]byte(s))
			} else if ir, irok := d.(io.Reader); irok {
				io.Copy(w, ir)
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

type WrapReader struct {
	scnr *bufio.Scanner
	rdr  io.Reader
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
				if rnsi > 0 && err == nil {
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
func ReaderToString(r io.Reader) (s string, err error) {
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
				rns[rnsi] = rn
				rnsi++
				if rnsi == len(rns) {
					s += string(rns[:rnsi])
					rnsi = 0
				}
			}
			if rnerr != nil {
				if rnerr != io.EOF {
					err = rnerr
				} else {
					if rnsi > 0 && err == nil {
						s += string(rns[:rnsi])
						rnsi = 0
					}
				}
				break
			}
		}
	}
	return
}

package iorw

import (
	"fmt"
	"io"
)

//Printer - interface
type Printer interface {
	Print(a ...interface{})
	Println(a ...interface{})
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

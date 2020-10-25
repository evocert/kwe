package iorw

import (
	"fmt"
	"io"
)

//Fprint - refer to fmt.Fprint
func Fprint(w io.Writer, a ...interface{}) {
	if len(a) > 0 && w != nil {
		fmt.Fprint(w, a...)
	}
}

//Fprintln - refer to fmt.Fprintln
func Fprintln(w io.Writer, a ...interface{}) {
	if len(a) > 0 && w != nil {
		Fprint(w, a...)
	}
	Fprint(w, "\r\n")
}

package main

import "github.com/evocert/kwe/iorw"

func main() {
	var buf = iorw.NewBuffer()
	if buf.Size() == 0 {
		buf.Close()
	}
}

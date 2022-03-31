package iorw

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type MultiArgsReader struct {
	args  []interface{}
	crntr io.Reader
	rnr   io.RuneReader
}

type multistringreader struct {
	strings []string
}

func NewMultiArgsReader(a ...interface{}) (mltiargsr *MultiArgsReader) {
	mltiargsr = &MultiArgsReader{args: a}
	return
}

func (mltiargsr *MultiArgsReader) nextrdr() (nxtrdr io.Reader) {
	if mltiargsr != nil {
		for nxtrdr == nil && len(mltiargsr.args) > 0 {
			d := mltiargsr.args[0]
			mltiargsr.args = mltiargsr.args[1:]
			if d != nil {
				if s, _ := d.(string); s != "" {
					nxtrdr = strings.NewReader(s)
				} else if rdr, _ := d.(io.Reader); rdr != nil {
					nxtrdr = rdr
				} else {
					nxtrdr = strings.NewReader(fmt.Sprint(d))
				}
			} else {
				continue
			}
		}
	}
	return
}

func (mltiargsr *MultiArgsReader) Read(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		for n < pl && err == nil {
			if mltiargsr != nil {
				if mltiargsr.crntr != nil {
					crntn, cnrterr := mltiargsr.crntr.Read(p[n : n+(pl-n)])
					n += crntn
					if cnrterr != nil {
						if cnrterr == io.EOF {
							mltiargsr.crntr = nil
							if nxtrdr := mltiargsr.nextrdr(); nxtrdr != nil {
								mltiargsr.crntr = nxtrdr
							} else if n == 0 {
								err = cnrterr
							}
							break
						} else {
							mltiargsr.crntr = nil
							err = cnrterr
						}
					}
				} else if mltiargsr.crntr == nil {
					if nxtrdr := mltiargsr.nextrdr(); nxtrdr != nil {
						mltiargsr.crntr = nxtrdr
					} else {
						break
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

func (mltiargsr *MultiArgsReader) ReadRune() (r rune, size int, err error) {
	if mltiargsr != nil {
		if mltiargsr.rnr == nil {
			mltiargsr.rnr = bufio.NewReader(mltiargsr)
		}
		r, size, err = mltiargsr.rnr.ReadRune()
	} else {
		err = io.EOF
	}
	return
}

func (mltiargsr *MultiArgsReader) Close() (err error) {
	if mltiargsr != nil {
		if mltiargsr.crntr != nil {
			mltiargsr.crntr = nil
		}
		if mltiargsr.args != nil {
			if len(mltiargsr.args) > 0 {
				for n, d := range mltiargsr.args {
					mltiargsr.args[n] = nil
					if d != nil {
						d = nil
					}
				}
				mltiargsr.args = nil
			}
		}
		if mltiargsr.rnr != nil {
			mltiargsr.rnr = nil
		}
		mltiargsr = nil
	}
	return
}

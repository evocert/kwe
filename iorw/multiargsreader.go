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
	buf   []byte
	bufi  int
	bufl  int
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

func multiArgsRead(mltiargsr *MultiArgsReader, p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		for n < pl && err == nil {
			if mltiargsr != nil {
				if mltiargsr.crntr != nil {
					crntn, cnrterr := mltiargsr.crntr.Read(p[n : n+(pl-n)])
					n += crntn
					if cnrterr != nil {
						if cnrterr == io.EOF {
							if mltiargsr.crntr = mltiargsr.nextrdr(); mltiargsr.crntr == nil {
								break
							}
						} else {
							mltiargsr.crntr = nil
							err = cnrterr
						}
					}
				} else if mltiargsr.crntr == nil {
					if mltiargsr.crntr = mltiargsr.nextrdr(); mltiargsr.crntr == nil {
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

func (mltiargsr *MultiArgsReader) Read(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		for n < pl && err == nil {
			if mltiargsr != nil {
				if mltiargsr.bufl == 0 || mltiargsr.bufl > 0 && mltiargsr.bufi == mltiargsr.bufl {
					if len(mltiargsr.buf) != 4096 {
						mltiargsr.buf = nil
						mltiargsr.buf = make([]byte, 4096)
					}
					pn, perr := multiArgsRead(mltiargsr, mltiargsr.buf)
					if pn > 0 {
						mltiargsr.buf = mltiargsr.buf[:pn]
						mltiargsr.bufi = 0
						mltiargsr.bufl = pn
					}
					if perr != nil {
						if perr != io.EOF {
							err = perr
							break
						}
					}
					if pn == 0 {
						break
					}
				}
				_, n, mltiargsr.bufi = CopyBytes(p, n, mltiargsr.buf[:mltiargsr.bufl], mltiargsr.bufi)
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
		if mltiargsr.buf != nil {
			mltiargsr.buf = nil
		}
		mltiargsr = nil
	}
	return
}

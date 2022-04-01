package iorw

import "io"

type BulkReader struct {
	r        io.Reader
	rc       io.ReadCloser
	buf      []byte
	bufl     int
	bufi     int
	canclose bool
}

func NewBulkReader(r io.Reader, canclose ...bool) (blkrdr *BulkReader) {
	blkrdr = &BulkReader{r: r, canclose: len(canclose) > 0 && canclose[0]}
	blkrdr.rc, _ = r.(io.ReadCloser)
	return
}

func (blkrdr *BulkReader) Read(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		if blkrdr != nil && blkrdr.r != nil {
			cpyl := 0
			for n < pl {
				if blkrdr.bufl == 0 || (blkrdr.bufl > 0 && blkrdr.bufi == blkrdr.bufl) {
					if len(blkrdr.buf) < 4096 {
						blkrdr.buf = make([]byte, 4096)
					}
					bulki := 0
					bufl, bulkerr := ReadHandle(blkrdr.r, func(b []byte) {
						bulki += copy(blkrdr.buf[bulki:bulki-(4096-bulki)], b)
					}, 4096)
					if bulkerr != nil && bulkerr != io.EOF {
						err = bulkerr
						break
					}
					if blkrdr.bufl = bufl; blkrdr.bufl == 0 {
						break
					}
					blkrdr.bufi = 0
				}
				cpyl, _, blkrdr.bufi = CopyBytes(p, n, blkrdr.buf, blkrdr.bufi)
				n += cpyl
			}
		}
		if n == 0 && err == nil {
			err = io.EOF
		}
	}
	return
}

func (blkrdr *BulkReader) Close() (err error) {
	if blkrdr != nil {
		if blkrdr.r != nil {
			blkrdr.r = nil
		}
		if blkrdr.rc != nil {
			if blkrdr.canclose {
				err = blkrdr.rc.Close()
			}
			blkrdr.rc = nil
		}
		if blkrdr.buf != nil {
			blkrdr.buf = nil
		}
		blkrdr = nil
	}
	return
}

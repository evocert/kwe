package parsing

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/evocert/kwe/iorw"
)

type ParsingReader struct {
	mltiargsrdrdrs []*iorw.MultiArgsReader
	eofrdr         *iorw.EOFCloseSeekReader
	prsng          *Parsing
	started        bool
	tmplts         map[string]*iorw.Buffer
	includes       map[string]*iorw.Buffer
}

func NewParseReader(prsng *Parsing, a ...interface{}) (prsngrdr *ParsingReader) {
	if len(a) > 0 && prsng != nil {
		prsngrdr = &ParsingReader{prsng: prsng, mltiargsrdrdrs: []*iorw.MultiArgsReader{iorw.NewMultiArgsReader(a...)}}
	}
	return
}

var intlbl = [][]rune{[]rune("[@"), []rune("@]")}

func processPreParsingFrase(bufrnwtrr *bufio.Writer, prsng *Parsing, prsngrdr *ParsingReader, frasebuf *iorw.Buffer) {
	if prsng != nil {
		if frasebuf.HasPrefix("include ") {
			if bufr := frasebuf.Reader(); bufr != nil {
				func() {
					defer bufr.Close()
					var ln = len([]byte("include "))
					var frases = strings.TrimSpace(bufr.SubString(int64(ln), frasebuf.Size()))
					if frases != "" {
						if inclbuf, _ := prsngrdr.includes[frases]; inclbuf != nil {
							func() {
								if prsngrdr.mltiargsrdrdrs == nil {
									prsngrdr.mltiargsrdrdrs = []*iorw.MultiArgsReader{}
								}
								prsngrdr.mltiargsrdrdrs = append([]*iorw.MultiArgsReader{iorw.NewMultiArgsReader(inclbuf.Reader())}, prsngrdr.mltiargsrdrdrs...)
								frases = ""
							}()
						} else if prsng.LookupTemplate != nil {
							if r, _ := prsng.LookupTemplate(frases); r != nil {
								func() {
									if prsngrdr.includes == nil {
										prsngrdr.includes = map[string]*iorw.Buffer{}
									}
									prsngrdr.includes[frases] = iorw.NewBuffer()
									prsngrdr.includes[frases].ReadFrom(r)
									r = prsngrdr.includes[frases].Reader()
									if prsngrdr.mltiargsrdrdrs == nil {
										prsngrdr.mltiargsrdrdrs = []*iorw.MultiArgsReader{}
									}
									prsngrdr.mltiargsrdrdrs = append([]*iorw.MultiArgsReader{iorw.NewMultiArgsReader(r)}, prsngrdr.mltiargsrdrdrs...)
									frases = ""
								}()
							}
						} else if prsng.AtvActv != nil {
							if r, _ := prsng.AtvActv.AltLookupTemplate(frases); r != nil {
								func() {
									if prsngrdr.includes == nil {
										prsngrdr.includes = map[string]*iorw.Buffer{}
									}
									prsngrdr.includes[frases] = iorw.NewBuffer()
									prsngrdr.includes[frases].ReadFrom(r)
									r = prsngrdr.includes[frases].Reader()
									if prsngrdr.mltiargsrdrdrs == nil {
										prsngrdr.mltiargsrdrdrs = []*iorw.MultiArgsReader{}
									}
									prsngrdr.mltiargsrdrdrs = append([]*iorw.MultiArgsReader{iorw.NewMultiArgsReader(r)}, prsngrdr.mltiargsrdrdrs...)
									frases = ""
								}()
							}
						}
					}
				}()
			}
		} else if frasebuf.HasPrefix("include-template") {
			func() {
				var tmpltrdr = frasebuf.Reader()
				defer tmpltrdr.Close()
				tmpltrdr.Seek(int64(len("include-template")), io.SeekStart)
				var tmplerr error = nil
				var tmplname = ""
				for tmplerr == nil {
					if tmplname, tmplerr = tmpltrdr.Readln(); tmplname != "" {
						tmplname = strings.TrimSpace(tmplname)

						if tmplname != "" {
							var tmpltbuf = iorw.NewBuffer()
							tmpltbuf.ReadRunesFrom(tmpltrdr)
							if tmpltbuf.Size() > 0 {
								if prsngrdr.tmplts == nil {
									prsngrdr.tmplts = map[string]*iorw.Buffer{}
								}
								if tmplbf, _ := prsngrdr.tmplts[tmplname]; tmplbf != nil {
									tmplbf.Clear()
									tmpltrdr.Close()
									tmpltrdr = tmpltbuf.Reader()
									tmplbf.ReadFrom(tmpltrdr)
									tmpltrdr.Close()
								} else {
									prsngrdr.tmplts[tmplname] = tmpltbuf
								}
							} else {
								tmpltbuf.Close()
							}
						}
						break
					} else {
						break
					}
				}
			}()
		} else if frasebuf.HasPrefix("tmpl:") {
			if bufr := frasebuf.Reader(); bufr != nil {
				func() {
					defer bufr.Close()
					var ln = len([]byte("tmpl:"))
					var frases = strings.TrimSpace(bufr.SubString(int64(ln), frasebuf.Size()))
					if frases != "" {
						if tmplbuf, _ := prsngrdr.tmplts[frases]; tmplbuf != nil {
							func() {
								if prsngrdr.mltiargsrdrdrs == nil {
									prsngrdr.mltiargsrdrdrs = []*iorw.MultiArgsReader{}
								}
								prsngrdr.mltiargsrdrdrs = append([]*iorw.MultiArgsReader{iorw.NewMultiArgsReader(tmplbuf.Reader())}, prsngrdr.mltiargsrdrdrs...)
								frases = ""
							}()
						}
					}
				}()
			}
		}
	}
}

func (prsngrdr *ParsingReader) Start() (rnerdr io.RuneReader) {
	if !prsngrdr.started {
		cntx, cntxcncl := context.WithCancel(context.Background())
		pi, pw := io.Pipe()
		go func() {
			var errrns error = nil
			var frasebuf = iorw.NewBuffer()
			bufrnwtrr := bufio.NewWriter(pw)
			finrunes := make([]rune, 4096)
			finrunesi := 0
			cchdr := rune(0)
			cchdsize := 0
			var prvpr rune = rune(0)
			var intlbli = []int{0, 0}
			var flushfin = func() {
				if finrunesi > 0 {
					bufrnwtrr.WriteString(string(finrunes[:finrunesi]))
					finrunesi = 0
					bufrnwtrr.Flush()
				}
			}
			var fincrn = func(frs ...rune) {
				for _, fr := range frs {
					finrunes[finrunesi] = fr
					finrunesi++
					if finrunesi == len(finrunes) {
						flushfin()
					}
				}
			}
			var intparsechar = func(pr rune) {
				if intlbli[1] == 0 && intlbli[0] < len(intlbl[0]) {
					if intlbli[0] > 0 && intlbl[0][intlbli[0]-1] == prvpr && intlbl[0][intlbli[0]] != pr {
						fincrn(intlbl[0][:intlbli[0]]...)
						intlbli[0] = 0
					}
					if intlbl[0][intlbli[0]] == pr {
						intlbli[0]++
						if intlbli[0] == len(intlbl[0]) {
							prvpr = rune(0)
						} else {
							prvpr = pr
						}
					} else {
						if intlbli[0] > 0 {
							fincrn(intlbl[0][:intlbli[0]]...)
							intlbli[0] = 0
						}
						prvpr = pr
						fincrn(pr)
					}
				} else if intlbli[0] == len(intlbl[0]) && intlbli[1] < len(intlbl[1]) {
					if intlbl[1][intlbli[1]] == pr {
						intlbli[1]++
						if intlbli[1] == len(intlbl[1]) {
							intlbli[0] = 0
							intlbli[1] = 0
							if frasebuf.Size() > 0 {
								flushfin()
								processPreParsingFrase(bufrnwtrr, prsngrdr.prsng, prsngrdr, frasebuf)
								frasebuf.Clear()
							} else {
								fincrn(intlbl[0]...)
								fincrn(intlbl[1]...)
							}
						}
					} else {
						if intlbli[1] > 0 {
							for _, fr := range intlbl[1][:intlbli[1]] {
								frasebuf.WriteRune(fr)
							}
							intlbli[1] = 0
						}
						frasebuf.WriteRune(pr)
					}
				}
			}
			var crntmtmltiargsrdrd = func() (mrgsrdr *iorw.MultiArgsReader) {
				if len(prsngrdr.mltiargsrdrdrs) > 0 {
					mrgsrdr = prsngrdr.mltiargsrdrdrs[0]
				}
				return
			}
			defer func() {
				frasebuf.Close()
				if errrns == nil || errrns == io.EOF {
					bufrnwtrr.Flush()
					pw.Close()
				} else if errrns != nil {
					pw.CloseWithError(errrns)
				}
				bufrnwtrr = nil
				finrunes = nil
				intlbli = nil

				flushfin = nil
				fincrn = nil
				intparsechar = nil
				crntmtmltiargsrdrd = nil
			}()
			cntxcncl()
			for {
				if mltiargsrdrd := crntmtmltiargsrdrd(); mltiargsrdrd == nil {
					break
				} else {
					cchdr, cchdsize, errrns = mltiargsrdrd.ReadRune()
					if cchdsize > 0 && errrns == nil {
						intparsechar(cchdr)
					} else if errrns != nil {
						if cchdsize > 0 {
							intparsechar(cchdr)
						}
						if errrns == io.EOF {
							if len(prsngrdr.mltiargsrdrdrs) > 0 {
								prsngrdr.mltiargsrdrdrs = prsngrdr.mltiargsrdrdrs[1:]
								if len(prsngrdr.mltiargsrdrdrs) > 0 {
									errrns = nil
								} else {
									break
								}
							} else {
								break
							}
						} else {
							break
						}
					}
				}
			}
			flushfin()
		}()
		<-cntx.Done()
		if prsngrdr.eofrdr == nil {
			prsngrdr.eofrdr = iorw.NewEOFCloseSeekReader(pi)
		}
		prsngrdr.started = true
	}
	return prsngrdr.eofrdr
}

func (prsngrdr *ParsingReader) ReadRune() (r rune, size int, err error) {
	if rnrrdr := prsngrdr.Start(); rnrrdr != nil {
		r, size, err = rnrrdr.ReadRune()
	} else {
		err = io.EOF
	}
	return
}

func (prsngrdr *ParsingReader) Close() (err error) {
	if prsngrdr != nil {
		if prsngrdr.eofrdr != nil {
			prsngrdr.eofrdr.Close()
			prsngrdr.eofrdr = nil
		}
		if prsngrdr.mltiargsrdrdrs != nil {
			for _, mltiargsrdr := range prsngrdr.mltiargsrdrdrs {
				mltiargsrdr.Close()
			}
			prsngrdr.mltiargsrdrdrs = nil
		}
		if prsngrdr.prsng != nil {
			prsngrdr.prsng = nil
		}
		if prsngrdr.includes != nil {
			for incld := range prsngrdr.includes {
				prsngrdr.includes[incld].Close()
				prsngrdr.includes[incld] = nil
				delete(prsngrdr.includes, incld)
			}
			prsngrdr.includes = nil
		}
		if prsngrdr.tmplts != nil {
			for tmpl := range prsngrdr.tmplts {
				prsngrdr.tmplts[tmpl].Close()
				prsngrdr.tmplts[tmpl] = nil
				delete(prsngrdr.tmplts, tmpl)
			}
			prsngrdr.tmplts = nil
		}
		prsngrdr = nil
	}
	return
}

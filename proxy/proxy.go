package proxy

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/evocert/kwe/iorw"
	"github.com/evocert/kwe/requesting"
)

func Proxy(rqstpathbase, rqstpath string, rqsttopath string, reqst requesting.RequestAPI, respns requesting.ResponseAPI) (err error) {
	if ishttp, ishttps, isws, iswss := strings.HasPrefix(rqsttopath, "http://"), strings.HasPrefix(rqsttopath, "https://"), strings.HasPrefix(rqsttopath, "ws://"), strings.HasPrefix(rqsttopath, "wss://"); ishttp || ishttps || isws || iswss {
		if ishttp || ishttps {
			var httpurl = rqsttopath + rqstpath[len(rqstpathbase):]
			if qrystring := reqst.QueryString(); qrystring != "" {
				if strings.LastIndex(httpurl, "?") > strings.LastIndex(httpurl, "/") {
					httpurl += "&" + qrystring
				} else {
					httpurl += "?" + qrystring
				}
			}

			if rqst, rqsterr := http.NewRequest(reqst.Method(), httpurl, iorw.NewEOFCloseSeekReader(reqst, false)); rqsterr == nil {
				for _, hdr := range reqst.Headers() {
					hdrv := reqst.Header(hdr)
					rqst.Header.Add(hdr, hdrv)
				}
				func() {
					httpclnt := &http.Client{}
					defer httpclnt.CloseIdleConnections()
					var isgzip = false
					var findBase = true
					var ishtml = false
					var isjs = false
					var iscss = false
					if rspns, rspnserr := httpclnt.Do(rqst); rspnserr == nil {
						respns.SetStatus(rspns.StatusCode)
						for hdr, hdrv := range rspns.Header {
							if hdr == "Link" {
								for hnv := range hdrv {
									if strings.HasPrefix(hdrv[hnv], "</") {
										hdrv[hnv] = "<" + rqstpathbase + hdrv[hnv][2:]
									}
								}
							}
							var tmpv = strings.Join(hdrv, ";")
							if strings.EqualFold(hdr, "Content-Length") {
								continue
							} else if strings.EqualFold(hdr, "Content-Encoding") {
								isgzip = strings.Contains(tmpv, "gzip")
							} else if strings.EqualFold(hdr, "Content-Type") {
								ishtml = strings.Contains(tmpv, "text/html")
								isjs = strings.Contains(tmpv, "application/javascript")
								iscss = strings.Contains(tmpv, "text/css")
								findBase = ishtml || isjs || iscss
							}
							respns.SetHeader(hdr, tmpv)
						}
						if respbdy := rspns.Body; respbdy != nil {
							var r = respbdy

							if findBase {
								ctx, ctxcncl := context.WithCancel(context.Background())
								pi, pw := io.Pipe()
								go func() {
									var perr = error(nil)
									defer func() {
										if (perr != nil && perr == io.EOF) || perr == nil {
											pw.Close()
										} else {
											pw.CloseWithError(perr)
										}
									}()
									var w io.Writer = pw
									if isgzip {
										r, _ = gzip.NewReader(respbdy)
										if !ishtml {
											w = gzip.NewWriter(pw)
										}
									}
									var eofr = iorw.NewEOFCloseSeekReader(r)
									ctxcncl()
									var basebts = []byte(rqsttopath)
									var basebtsi = 0
									var p = make([]byte, 8192)
									var pbts = make([]byte, 1)
									var pn = 0
									var prvb = byte(0)
									var pbuf = iorw.NewBuffer()
									for perr == nil {
										if pn, perr = eofr.Read(p); (perr == nil || perr == io.EOF) && pn > 0 {
											for _, pb := range p[:pn] {
												if basebtsi > 0 && basebts[basebtsi-1] == prvb && basebts[basebtsi] != pb {
													pbuf.Write(basebts[:basebtsi])
													basebtsi = 0
												}
												if basebts[basebtsi] == pb {
													basebtsi++
													if basebtsi == len(basebts) {
														pbuf.Print(rqstpathbase)
														if pbuf.Size() > 0 {
															if isjs {
																var wndwbaseurl = ""
																var rplcewndwbaseurl = ""
																var rcplbasepath = rqstpathbase
																if len(rcplbasepath) > 1 && strings.HasSuffix(rcplbasepath, "/") {
																	rcplbasepath = rcplbasepath[:len(rcplbasepath)-1]
																}

																if wndwbaseurl = `window.BaseURL = ""`; pbuf.Contains(wndwbaseurl) {

																	rplcewndwbaseurl = `window.BaseURL = "` + rcplbasepath + `"`
																} else if wndwbaseurl = `window.BaseURL = ''`; pbuf.Contains(wndwbaseurl) {
																	rplcewndwbaseurl = `window.BaseURL = '` + rcplbasepath + `'`
																} else if wndwbaseurl = `window.BaseURL=""`; pbuf.Contains(wndwbaseurl) {
																	rplcewndwbaseurl = `window.BaseURL="` + rcplbasepath + `"`
																} else if wndwbaseurl = `window.BaseURL=''`; pbuf.Contains(wndwbaseurl) {
																	rplcewndwbaseurl = `window.BaseURL='` + rcplbasepath + `'`
																}
																if rplcewndwbaseurl != "" {
																	iorw.Fprint(w, strings.Replace(pbuf.String(), wndwbaseurl, rplcewndwbaseurl, -1))
																} else {
																	iorw.Fprint(w, iorw.NewEOFCloseSeekReader(pbuf.Reader()))
																}
															} else if iscss {
																var containsRstBasePath = pbuf.Contains(rqstpathbase)
																if pbuf.Contains(`url("/`) {
																	if containsRstBasePath {
																		iorw.Fprint(w, strings.Replace(strings.Replace(pbuf.String(), `url("/`, `url("`+rqstpathbase, -1), rqstpathbase+rqstpathbase[1:], rqstpathbase, -1))
																	} else {
																		iorw.Fprint(w, strings.Replace(pbuf.String(), `url("/`, `url("`+rqstpathbase, -1))
																	}
																} else if pbuf.Contains(`url('/`) {
																	if containsRstBasePath {
																		iorw.Fprint(w, strings.Replace(strings.Replace(pbuf.String(), `url('/`, `url('`+rqstpathbase, -1), rqstpathbase+rqstpathbase[1:], rqstpathbase, -1))
																	} else {
																		iorw.Fprint(w, strings.Replace(pbuf.String(), `url('/`, `url('`+rqstpathbase, -1))
																	}
																} else if pbuf.Contains(`url(/`) {
																	if containsRstBasePath {
																		iorw.Fprint(w, strings.Replace(strings.Replace(pbuf.String(), `url(/`, `url(`+rqstpathbase, -1), rqstpathbase+rqstpathbase[1:], rqstpathbase, -1))
																	} else {
																		iorw.Fprint(w, strings.Replace(pbuf.String(), `url(/`, `url(`+rqstpathbase, -1))
																	}
																} else {
																	iorw.Fprint(w, iorw.NewEOFCloseSeekReader(pbuf.Reader()))
																}
															} else {
																iorw.Fprint(w, iorw.NewEOFCloseSeekReader(pbuf.Reader()))
															}
															pbuf.Clear()
														}
														basebtsi = 0
														prvb = 0
													} else {
														prvb = pb
													}
												} else {
													if basebtsi > 0 {
														pbuf.Write(basebts[:basebtsi])
														basebtsi = 0
													}
													prvb = pb
													pbts[0] = pb
													pbuf.Write(pbts)
												}
											}
										} else if pn == 0 {
											perr = io.EOF
										}
									}
									if pbuf.Size() > 0 {
										if isjs {
											var wndwbaseurl = ""
											var rplcewndwbaseurl = ""
											var rcplbasepath = rqstpathbase
											if len(rcplbasepath) > 1 && strings.HasSuffix(rcplbasepath, "/") {
												rcplbasepath = rcplbasepath[:len(rcplbasepath)-1]
											}

											if wndwbaseurl = `window.BaseURL = ""`; pbuf.Contains(wndwbaseurl) {

												rplcewndwbaseurl = `window.BaseURL = "` + rcplbasepath + `"`
											} else if wndwbaseurl = `window.BaseURL = ''`; pbuf.Contains(wndwbaseurl) {
												rplcewndwbaseurl = `window.BaseURL = '` + rcplbasepath + `'`
											} else if wndwbaseurl = `window.BaseURL=""`; pbuf.Contains(wndwbaseurl) {
												rplcewndwbaseurl = `window.BaseURL="` + rcplbasepath + `"`
											} else if wndwbaseurl = `window.BaseURL=''`; pbuf.Contains(wndwbaseurl) {
												rplcewndwbaseurl = `window.BaseURL='` + rcplbasepath + `'`
											}
											if rplcewndwbaseurl != "" {
												iorw.Fprint(w, strings.Replace(pbuf.String(), wndwbaseurl, rplcewndwbaseurl, -1))
											} else {
												iorw.Fprint(w, iorw.NewEOFCloseSeekReader(pbuf.Reader()))
											}
										} else if iscss {
											var containsRstBasePath = pbuf.Contains(rqstpathbase)
											if pbuf.Contains(`url("/`) {
												if containsRstBasePath {
													iorw.Fprint(w, strings.Replace(strings.Replace(pbuf.String(), `url("/`, `url("`+rqstpathbase, -1), rqstpathbase+rqstpathbase[1:], rqstpathbase, -1))
												} else {
													iorw.Fprint(w, strings.Replace(pbuf.String(), `url("/`, `url("`+rqstpathbase, -1))
												}
											} else if pbuf.Contains(`url('/`) {
												if containsRstBasePath {
													iorw.Fprint(w, strings.Replace(strings.Replace(pbuf.String(), `url('/`, `url('`+rqstpathbase, -1), rqstpathbase+rqstpathbase[1:], rqstpathbase, -1))
												} else {
													iorw.Fprint(w, strings.Replace(pbuf.String(), `url('/`, `url('`+rqstpathbase, -1))
												}
											} else if pbuf.Contains(`url(/`) {
												if containsRstBasePath {
													iorw.Fprint(w, strings.Replace(strings.Replace(pbuf.String(), `url(/`, `url(`+rqstpathbase, -1), rqstpathbase+rqstpathbase[1:], rqstpathbase, -1))
												} else {
													iorw.Fprint(w, strings.Replace(pbuf.String(), `url(/`, `url(`+rqstpathbase, -1))
												}
											} else {
												iorw.Fprint(w, iorw.NewEOFCloseSeekReader(pbuf.Reader()))
											}
										} else {
											iorw.Fprint(w, iorw.NewEOFCloseSeekReader(pbuf.Reader()))
										}

										pbuf.Clear()
									}
									if gzpw, _ := w.(*gzip.Writer); gzpw != nil {
										gzpw.Flush()
									}
								}()
								<-ctx.Done()
								r = pi
							}
							if ishtml {
								pihtml, pwhtml := io.Pipe()
								ctxhtml, ctxcnclhtml := context.WithCancel(context.Background())
								go func() {
									var perr = error(nil)
									defer func() {
										if (perr != nil && perr == io.EOF) || perr == nil {
											pwhtml.Close()
										} else {
											pwhtml.CloseWithError(perr)
										}
									}()
									var w io.Writer = pwhtml

									if isgzip {
										w = gzip.NewWriter(pwhtml)
									}
									var eofr = iorw.NewEOFCloseSeekReader(r)
									ctxcnclhtml()
									var headbts = []byte("<head>")
									var headbtsl = len(headbts)
									var headendbts = []byte("</head>")
									var headendbtsl = len(headendbts)
									var headbtsi = []int{0, 0}
									var p = make([]byte, 8192)
									var pbts = make([]byte, 1)
									var pn = 0
									var prvb = byte(0)
									var pbuf = iorw.NewBuffer()
									var pbufheader = iorw.NewBuffer()
									for perr == nil {
										if pn, perr = eofr.Read(p); (perr == nil || perr == io.EOF) && pn > 0 {
											for _, pb := range p[:pn] {
												if headbtsi[1] == 0 && headbtsi[0] < headbtsl {
													if headbtsi[0] > 0 && headbts[headbtsi[0]-1] == prvb && headbts[headbtsi[0]] != pb {
														pbuf.Write(headbts[:headbtsi[0]])
														headbtsi[0] = 0
													}
													if headbts[headbtsi[0]] == pb {
														headbtsi[0]++
														if headbtsi[0] == headbtsl {
															pbuf.Write(headbts)
															if pbuf.Size() > 0 {
																iorw.Fprint(w, iorw.NewEOFCloseSeekReader(pbuf.Reader()))
																pbuf.Clear()
															}
															prvb = 0
														} else {
															prvb = pb
														}
													} else {
														if headbtsi[0] > 0 {
															pbuf.Write(headbts[:headbtsi[0]])
															headbtsi[0] = 0
														}
														prvb = pb
														pbts[0] = pb
														pbuf.Write(pbts)
													}
												} else if headbtsi[0] == headbtsl && headbtsi[1] < headendbtsl {
													if headendbts[headbtsi[1]] == pb {
														headbtsi[1]++
														if headbtsi[1] == headendbtsl {
															if pbufheader.Size() > 0 {
																func() {
																	var containsRstBasePath = pbufheader.Contains(rqstpathbase)
																	if pbufheader.Contains(`"/`) {
																		if containsRstBasePath {
																			pbuf.Print(strings.Replace(strings.Replace(pbufheader.String(), `"/`, `"`+rqstpathbase, -1), rqstpathbase+rqstpathbase[1:], rqstpathbase, -1))
																		} else {
																			pbuf.Print(strings.Replace(pbufheader.String(), `"/`, `"`+rqstpathbase, -1))
																		}
																	} else if pbufheader.Contains(`'/`) {
																		if containsRstBasePath {
																			pbuf.Print(strings.Replace(strings.Replace(pbufheader.String(), `'/`, `'`+rqstpathbase, -1), rqstpathbase+rqstpathbase[1:], rqstpathbase, -1))
																		} else {
																			pbuf.Print(strings.Replace(pbufheader.String(), `'/`, `'`+rqstpathbase, -1))
																		}
																	} else {
																		pbuhhdrrdr := pbufheader.Reader()
																		defer pbuhhdrrdr.Close()
																		pbuf.Print(pbuhhdrrdr)
																	}
																}()
															}
															pbuf.Write(headendbts)
															headbtsi[0] = 0
															headbtsi[1] = 0
															prvb = 0
														}
													} else {
														if headbtsi[1] > 0 {
															pbufheader.Write(headendbts[:headbtsi[1]])
															headbtsi[1] = 0
														}
														pbts[0] = pb
														pbufheader.Write(pbts)
													}
												}
											}
										} else if pn == 0 {
											perr = io.EOF
										}
									}
									if pbuf.Size() > 0 {
										iorw.Fprint(w, iorw.NewEOFCloseSeekReader(pbuf.Reader()))
										pbuf.Clear()
									}
									if gzpw, _ := w.(*gzip.Writer); gzpw != nil {
										gzpw.Flush()
									}
								}()
								<-ctxhtml.Done()
								r = pihtml
							}
							respns.Print(r)
						}
					} else {
						err = rspnserr
					}
				}()
			}
		}
	}
	return
}

package requesting

import "io"

type Requester struct {
	rqst  RequestAPI
	rspns ResponseAPI
}

func (rqstr *Requester) Request() (rqst RequestAPI) {
	if rqstr != nil {
		rqst = rqstr.rqst
	}
	return
}

func (rqstr *Requester) Response() (rspns ResponseAPI) {
	if rqstr != nil {
		rspns = rqstr.rspns
	}
	return
}

func (rqstr *Requester) IsValid() (valid bool, err error) {
	if rqstr != nil {
		if rqstr.rqst != nil {
			if valid, err = rqstr.rqst.IsValid(); valid && err == nil {
				if rqstr.rspns != nil {
					valid, err = rqstr.rspns.IsValid()
				}
			}
		} else if rqstr.rspns != nil {
			if valid, err = rqstr.rspns.IsValid(); valid && err == nil {
				if rqstr.rqst != nil {
					valid, err = rqstr.rqst.IsValid()
				}
			}
		} else {
			valid, err = true, nil
		}
	}
	return
}

func (rqstr *Requester) Close() (err error) {
	if rqstr != nil {
		if rqstr.rqst != nil {
			rqstr.rqst.Close()
			rqstr.rqst = nil
		}
		if rqstr.rspns != nil {
			rqstr.rspns.Close()
			rqstr.rspns = nil
		}
		rqstr = nil
	}
	return
}

type RequestInvokerFunc func(rdr io.Reader, a ...interface{}) RequestAPI
type ResponseInvokerFunc func(w io.Writer, a ...interface{}) ResponseAPI

func NewRequester(a ...interface{}) (rqstor *Requester) {
	var rspns ResponseAPI = nil
	var rqst RequestAPI = nil
	var rqstargs []interface{} = nil
	var rspnsargs []interface{} = nil
	var rqstinvoker RequestInvokerFunc
	var rspnsinvoker ResponseInvokerFunc
	var wtr io.Writer = nil
	var rdr io.Reader = nil
	if len(a) > 0 {
		for _, d := range a {
			if rqstd, _ := d.(RequestAPI); rqstd != nil {
				if rqst == nil {
					rqst = rqstd
				}
			} else if rdrd, _ := d.(io.Reader); rdrd != nil {
				if rdr == nil {
					rdr = rdrd
				}
			} else if rspnsd, _ := d.(ResponseAPI); rspnsd != nil {
				if rspns == nil {
					rspns = rspnsd
				}
			} else if wtrd, _ := d.(io.Writer); wtrd != nil {
				if wtr == nil {
					wtr = wtrd
				}
			} else if rqstinvokerd, _ := d.(func(rdr io.Reader, a ...interface{}) RequestAPI); rqstinvokerd != nil {
				if rqstinvoker == nil {
					rqstinvoker = rqstinvokerd
				}
			} else if rspnsinvokerd, _ := d.(func(w io.Writer, a ...interface{}) ResponseAPI); rspnsinvokerd != nil {
				if rspnsinvoker == nil {
					rspnsinvoker = rspnsinvokerd
				}
			} else if rqstinvokerd, _ := d.(RequestInvokerFunc); rqstinvokerd != nil {
				if rqstinvoker == nil {
					rqstinvoker = rqstinvokerd
				}
			} else if rspnsinvokerd, _ := d.(ResponseInvokerFunc); rspnsinvokerd != nil {
				if rspnsinvoker == nil {
					rspnsinvoker = rspnsinvokerd
				}
			} else if argsd, _ := d.([]interface{}); len(argsd) > 0 {
				if len(rqstargs) == 0 {
					rqstargs = argsd
				} else if len(rspnsargs) == 0 {
					rspnsargs = argsd
				}
			}
		}
	}
	if rspns == nil && rspnsinvoker == nil && len(rspnsargs) > 0 {
		rspnsinvoker = DefaultResponseInvoker
	}

	if rqst == nil && rqstinvoker == nil && len(rqstargs) > 0 {
		rqstinvoker = DefaultRequestInvoker
	}

	if rqst == nil && rqstinvoker != nil {
		rqst = rqstinvoker(rdr, rqstargs...)
	}

	if rspns == nil && rspnsinvoker != nil {
		if rqst != nil {
			rspnsargs = append([]interface{}{rqst}, rspnsargs...)
		}
		rspns = rspnsinvoker(wtr, rspnsargs...)
	}

	rqstor = &Requester{rspns: rspns, rqst: rqst}
	return
}

var DefaultRequestInvoker RequestInvokerFunc = nil
var DefaultResponseInvoker ResponseInvokerFunc = nil

type RequestorHandler interface {
	ServeREQUEST(RequesterAPI)
}

type RequestorHandlerFunc func(RequesterAPI)

func (rqstrhdnlrfnc RequestorHandlerFunc) ServeREQUEST(rqstr RequesterAPI) {
	defer rqstr.Close()
	rqstrhdnlrfnc(rqstr)
}

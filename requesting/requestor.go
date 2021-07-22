package requesting

import (
	"github.com/evocert/kwe/parameters"
)

type RequestorHandler interface {
	Serve(string, RequestAPI, ResponseAPI) error
	ServeRequest(...interface{}) error
}

type RequestorHandlerFunc func(string, RequestAPI, ResponseAPI) error

func (rqstrhndlfunc RequestorHandlerFunc) Serve(path string, rqst RequestAPI, rspns ResponseAPI) (err error) {
	err = rqstrhndlfunc(path, rqst, rspns)
	return
}

type RequestorAPI interface {
	Request() RequestAPI
	Response() ResponseAPI
	IsValid() (bool, error)
	Close() error
	LoadParameters(*parameters.Parameters)
}

type Requestor struct {
	rqst  RequestAPI
	rspns ResponseAPI
}

type RequestInvokerFunc func(path string, r interface{}) RequestAPI
type ResponseInvokerFunc func(w interface{}, a ...RequestAPI) ResponseAPI

func NewRequestor(a ...interface{}) (rqstor *Requestor) {
	var rspns ResponseAPI = nil
	var rqst RequestAPI = nil
	var path string = ""
	var rqstarg interface{} = nil
	var rspnsarg interface{} = nil
	var rqstinvoker RequestInvokerFunc
	var rspnsinvoker ResponseInvokerFunc
	if len(a) > 0 {
		for _, d := range a {
			if sd, _ := d.(string); sd != "" {
				if path == "" {
					path = sd
				}
			} else if rspnsd, _ := d.(ResponseAPI); rspnsd != nil {
				if rspns == nil {
					rspns = rspnsd
				}
			} else if rqstd, _ := d.(RequestAPI); rqstd != nil {
				if rqst == nil {
					rqst = rqstd
				}
			} else if rqstinvokerd, _ := d.(func(path string, r interface{}) RequestAPI); rqstinvokerd != nil {
				if rqstinvoker == nil {
					rqstinvoker = rqstinvokerd
				}
			} else if rspnsinvokerd, _ := d.(func(w interface{}, a ...RequestAPI) ResponseAPI); rspnsinvokerd != nil {
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
			} else if rqstarg == nil {
				rqstarg = d
			} else if rspnsarg == nil {
				rspnsarg = d
			}
		}
	}
	if rspns == nil && rspnsinvoker == nil && rspnsarg != nil {
		rspnsinvoker = DefaultResponseInvoker
	}

	if rqst == nil && rqstinvoker == nil && (rqstarg != nil || path != "") {
		rqstinvoker = DefaultRequestInvoker
	}

	if rqst == nil && rqstinvoker != nil {
		rqst = rqstinvoker(path, rqstarg)
	}

	if rspns == nil && rspnsinvoker != nil {
		rspns = rspnsinvoker(rspnsarg, rqst)
	}

	if rqst == nil && rspns != nil {
		rqst = rspns.Request()
	}
	rqstor = &Requestor{rspns: rspns, rqst: rqst}
	return
}

func (rqstor *Requestor) LoadParameters(prms *parameters.Parameters) {
	if rqstor != nil && rqstor.rqst != nil {
		rqstor.rqst.LoadParameters(prms)
	}
}

func (rqstor *Requestor) Request() (rwrqst RequestAPI) {
	if rqstor != nil {
		rwrqst = rqstor.rqst
	}
	return
}

func (rqstor *Requestor) IsValid() (valid bool, err error) {
	if rqstor != nil {
		if rqstor.rqst != nil {
			if valid, err = rqstor.rqst.IsValid(); valid && err == nil {
				if rqstor.rspns != nil {
					valid, err = rqstor.rspns.IsValid()
				}
			}
		} else if rqstor.rspns != nil {
			if valid, err = rqstor.rspns.IsValid(); valid && err == nil {
				if rqstor.rqst != nil {
					valid, err = rqstor.rqst.IsValid()
				}
			}
		}
	} else {
		valid = true
	}
	return
}

func (rqstor *Requestor) Response() (rwrspns ResponseAPI) {
	if rqstor != nil {
		rwrspns = rqstor.rspns
	}
	return
}

func (rqstor *Requestor) Close() (err error) {
	if rqstor != nil {
		if rqstor.rqst != nil {
			rqstor.rqst.Close()
			rqstor.rqst = nil
		}
		if rqstor.rspns != nil {
			rqstor.rspns.Close()
			rqstor.rspns = nil
		}
		rqstor = nil
	}
	return
}

var DefaultRequestInvoker RequestInvokerFunc = nil
var DefaultResponseInvoker ResponseInvokerFunc = nil

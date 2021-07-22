package requesting

import (
	"github.com/evocert/kwe/parameters"
)

type RequestorHandler interface {
	Serve(string, RequestAPI, ResponseAPI) error
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

func NewRequestor(rwrqst RequestAPI, rwrspns ResponseAPI) (rqstor *Requestor) {
	rqstor = &Requestor{rspns: rwrspns, rqst: rwrqst}
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

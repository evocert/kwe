package fasthttp

import (
	"context"

	"github.com/evocert/kwe/chnls"
	"github.com/evocert/kwe/parameters"
	"github.com/evocert/kwe/requesting"
	"github.com/valyala/fasthttp"
)

type FasttHttpRequestor struct {
	ctxvalid     context.Context
	fstctx       *fasthttp.RequestCtx
	fsthttprqst  *FastHttpRequest
	fsthttpsrpns *FastHttpResponse
}

func iniFastHttpRequestor(fstctx *fasthttp.RequestCtx) (fsthttprqstor *FasttHttpRequestor) {
	fsthttprqstor = &FasttHttpRequestor{fstctx: fstctx}
	fsthttprqstor.fsthttprqst = iniFastHttpRequest(fsthttprqstor)
	fsthttprqstor.fsthttpsrpns = iniFastHttpResponse(fsthttprqstor)
	return
}
func NewFastHttpRequestor(fstctx *fasthttp.RequestCtx) (rqstor requesting.RequestorAPI) {
	rqstor = iniFastHttpRequestor(fstctx)
	return
}

func (fsthttprqstor *FasttHttpRequestor) LoadParameters(prms *parameters.Parameters) {
	if fsthttprqstor != nil && fsthttprqstor.fsthttprqst != nil {
		fsthttprqstor.fsthttprqst.LoadParameters(prms)
	}
}

func (fsthttprqstor *FasttHttpRequestor) Request() (rqst requesting.RequestAPI) {
	if fsthttprqstor != nil {
		rqst = fsthttprqstor.fsthttprqst
	}
	return
}

func (fsthttprqstor *FasttHttpRequestor) Response() (rqst requesting.ResponseAPI) {
	if fsthttprqstor != nil {
		rqst = fsthttprqstor.fsthttpsrpns
	}
	return
}

func (fsthttprqstor *FasttHttpRequestor) IsValid() (valid bool, err error) {
	if fsthttprqstor != nil && fsthttprqstor.ctxvalid != nil {
		select {
		case <-fsthttprqstor.ctxvalid.Done():
			valid, err = false, fsthttprqstor.ctxvalid.Err()
		default:
			valid = true
		}
	} else {
		valid = true
	}
	return
}

func (fsthttprqstor *FasttHttpRequestor) Close() (err error) {
	if fsthttprqstor != nil {
		if fsthttprqstor.fstctx != nil {
			fsthttprqstor.fstctx = nil
		}
		if fsthttprqstor.fsthttprqst != nil {
			fsthttprqstor.fsthttprqst.Close()
			fsthttprqstor.fsthttprqst = nil
		}
		if fsthttprqstor.fsthttpsrpns != nil {
			fsthttprqstor.fsthttpsrpns.Close()
			fsthttprqstor.fsthttpsrpns = nil
		}
	}
	return
}

type FastHttpRequestorHandler func(ctx *fasthttp.RequestCtx)

var DefaultFastHttpRequestHandler func(ctx *fasthttp.RequestCtx)

func init() {
	DefaultFastHttpRequestHandler = FastHttpRequestorHandler(func(ctx *fasthttp.RequestCtx) {
		chnls.GLOBALCHNL().ServeRequest(iniFastHttpRequestor(ctx))
	})
}

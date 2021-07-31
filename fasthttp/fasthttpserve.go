package fasthttp

import (
	"log"

	"github.com/valyala/fasthttp"
)

func ListenAndServe(listenAddr string) {
	if err := fasthttp.ListenAndServe(listenAddr, DefaultFastHttpRequestHandler); err != nil {
		log.Fatalf("error in ListenAndServe: %s", err)
	}
}

package api

import (
	"github.com/evocert/kwe/channeling/channelingapi"
)

var FAFExecute func(ssn channelingapi.SessionAPI, a ...interface{}) (err error) = nil

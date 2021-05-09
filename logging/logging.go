package logging

import (
	"github.com/rs/zerolog"

	_ "github.com/rs/zerolog/log"
)

type Log struct {
}

func (lg *Log) Write(p []byte) (n int, err error) {

	return
}

func (lg *Log) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {

	return
}

func NewLogger(lg *Log) (lggr interface{}) {
	lggr = zerolog.New(lg)
	return
}

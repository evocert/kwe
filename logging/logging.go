package logging

import "log"

//"github.com/rs/zerolog"

//_ "github.com/rs/zerolog/log"

type LogWriter interface {
	Write([]byte) (int, error)
	Print(a ...interface{})
	Println(a ...interface{})
	LOG(prefx string, postfx string, a ...interface{}) LogWriter
	INFO(a ...interface{}) LogWriter
	WARN(a ...interface{}) LogWriter
	ERR(err error, a ...interface{}) LogWriter
	DEBUG(a ...interface{}) LogWriter
	TRACE(a ...interface{}) LogWriter
}

type logWriter struct {
	lggr    *log.Logger
	PREFIX  string
	POSTFIX string
}

func (lg *logWriter) Write(p []byte) (n int, err error) {

	return
}

func (lg *logWriter) Print(a ...interface{}) {

}

func (lg *logWriter) Println(a ...interface{}) {

}

func (lg *logWriter) LOG(prefx string, postfx string, a ...interface{}) LogWriter {

	return lg
}

func (lg *logWriter) INFO(a ...interface{}) LogWriter {

	return lg
}

func (lg *logWriter) WARN(a ...interface{}) LogWriter {

	return lg
}
func (lg *logWriter) ERR(err error, a ...interface{}) LogWriter {

	return lg
}

func (lg *logWriter) DEBUG(a ...interface{}) LogWriter {

	return lg
}

func (lg *logWriter) TRACE(a ...interface{}) LogWriter {

	return lg
}

func NewLogger(lgwriter ...LogWriter) (lggr interface{}) {
	if len(lgwriter) == 1 && lgwriter[0] != nil {
		lggr = lgwriter[0]
	}
	return
}

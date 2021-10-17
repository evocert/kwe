package requesting

import "github.com/evocert/kwe/parameters"

type RequesterAPI interface {
	Request() RequestAPI
	Response() ResponseAPI
	IsValid() (bool, error)
	Close() error
}

type RequestAPI interface {
	Parameters() parameters.ParametersAPI
	Proto() string
	Method() string
	Path() string
	RangeType() string
	RangeOffset() int64
	Headers() []string
	Header(string) string
	RemoteAddr() string
	LocalAddr() string
	Read([]byte) (int, error)
	ReadRune() (rune, int, error)
	Readln() (string, error)
	ReadLines() ([]string, error)
	ReadAll() (string, error)
	IsValid() (bool, error)
	Close() error
	Response() ResponseAPI
}

type ResponseAPI interface {
	IsValid() (bool, error)
	Headers() []string
	Header(string) string
	SetHeader(string, string)
	//SetContentType(string)
	//ContentType() string
	SetStatus(int)
	Print(...interface{})
	Println(...interface{})
	Write([]byte) (int, error)
	Flush()
	Close() error
}

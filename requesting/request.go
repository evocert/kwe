package requesting

import "github.com/evocert/kwe/parameters"

type RequestAPI interface {
	IsValid() (bool, error)
	Proto() string
	Method() string
	Path() string
	LoadParameters(*parameters.Parameters)
	Headers() []string
	Header(string) string
	RemoteAddr() string
	LocalAddr() string
	StartedReading() error
	Readln() (string, error)
	Readlines() ([]string, error)
	ReadAll() (string, error)
	Read([]byte) (int, error)
	SetMaxRead(int64) error
	Seek(int64, int) (int64, error)
	ReadRune() (rune, int, error)
	Close() error
}

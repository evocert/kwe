package requesting

type ResponseAPI interface {
	Request() RequestAPI
	IsValid() (bool, error)
	Headers() []string
	Header(string) string
	SetHeader(string, string)
	SetContentType(string)
	ContentType() string
	SetStatus(int)
	StartedWriting(...bool) error
	Print(...interface{})
	Println(...interface{})
	Write([]byte) (int, error)
	Flush()
	Close() error
}

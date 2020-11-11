package active

//Runtime - interface
type Runtime interface {
	InvokeFunction(interface{}, ...interface{})
}

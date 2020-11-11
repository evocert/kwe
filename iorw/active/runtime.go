package active

type Runtime interface {
	InvokeFunction(interface{}, ...interface{})
}

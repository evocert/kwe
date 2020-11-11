package active

//Runetime - interface
type Runtime interface {
	InvokeFunction(interface{}, ...interface{})
}

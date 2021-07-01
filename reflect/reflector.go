package reflect

type RefletorAPI interface {
}

type Reflector struct {
	ownerref interface{}
	owner    interface{}
	mthds    map[string]*method
	flds     map[string]*field
}

type method struct {
	rflctr *Reflector
}

type field struct {
	rfltr *Reflector
}

func NewReflector(owner interface{}) (rflctr *Reflector) {
	if owner != nil {
		rflctr.owner = owner
		rflctr.ownerref = &owner
	}
	return
}

func call(rflctr *Reflector, callname string, args ...interface{}) (rval interface{}, err error) {

	return
}

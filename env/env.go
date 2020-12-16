package env

var wrapupcalls []func()

//ShutdownEnvironment - cleanup (shutdown) environment
func ShutdownEnvironment() {
	if len(wrapupcalls) > 0 {
		for _, wrpupcall := range wrapupcalls {
			wrpupcall()
		}
	}
}

//WrapupCall - set WrapupCall
func WrapupCall(wrpupcall ...func()) {
	if len(wrpupcall) > 0 {
		if len(wrapupcalls) == 0 {
			wrapupcalls = []func(){}
		}
		wrapupcalls = append(wrapupcalls, wrpupcall...)
	}
}

func init() {
	//network.RegisterShutdownEnv(ShutdownEnvironment)
}

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

type EnvAPI interface {
	Set(name string, val interface{})
	Get(name string) (val interface{})
	Keys() (keys []string)
}

type envrmnt struct {
	envsettions map[string]interface{}
}

func (env *envrmnt) Set(name string, val interface{}) {
	if env != nil {
		if env.envsettions[name] != nil {
			env.envsettions[name] = nil
		}
		env.envsettions[name] = val
	}
}

func (env *envrmnt) Get(name string) (val interface{}) {
	if env != nil {
		val = env.envsettions[name]
	}
	return
}

func (env *envrmnt) Keys() (keys []string) {
	if env != nil {
		if len(env.envsettions) > 0 {
			keys = make([]string, len(env.envsettions))
			keysi := 0
			for key := range env.envsettions {
				keys[keysi] = key
				keysi++
			}
		}
	}
	return
}

var env *envrmnt = nil

func Env() *envrmnt {
	return env
}

func init() {
	env = &envrmnt{envsettions: map[string]interface{}{}}
}

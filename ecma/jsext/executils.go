package jsext
import(
	"github.com/dop251/goja"
)
func Register_jsext_executils(vm*goja.Runtime){
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json",true))
	type Version struct{
		Major int `json:"major"`
		Minor int`json:"minor"`
		Bump int`json:"bump"`
	}
	//todo: namespace everything kwe.fsutils.etcetcetc
	//first test for kwe then do set fsutils on kwe
	vm.Set("executils",struct{
		Version Version`json:"version"`
		About func()(string)`json:"about"`
	}{
		Version:Version{
			Major:0,
			Minor:0,
			Bump:0,
		},
		About:func()(string){
			return "executils contains various utility functions for process execution"
		},
	})
}

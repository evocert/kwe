package jsext
import(
	"github.com/dop251/goja"
	"os/exec"
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
		//Exec func(string,...string)(string)`json:"exec"`
		Exec func(...string)(string)`json:"exec"`
	}{
		Version:Version{
			Major:0,
			Minor:0,
			Bump:0,
		},
		About:func()(string){
			return "executils contains various utility functions for process execution"
		},
		Exec:func(args...string)(string){
			cmdpath:=args[0]
			cmdargs:=args[1:]
			cmd:=exec.Command(cmdpath,cmdargs...)
			out,err:=cmd.Output()
			if err!=nil{
				panic(vm.ToValue("Execution error"))
			}
			return string(out)
		},
	})
}

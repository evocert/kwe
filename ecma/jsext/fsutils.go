package jsext
import(
	"github.com/dop251/goja"
	"io/ioutil"
	"io"
	"os"
)
func Register_jsext_fsutils(vm*goja.Runtime){
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json",true))
	type Version struct{
		Major int `json:"major"`
		Minor int`json:"minor"`
		Bump int`json:"bump"`
	}
	//todo: namespace everything kwe.fsutils.etcetcetc
	//first test for kwe then do set fsutils on kwe
	vm.Set("fsutils",struct{
		Version Version`json:"version"`
		File2String func(string)(string)`json:"file2string"`
		String2File func(string,string)`json:"string2file"`
	}{
		Version:Version{
			Major:0,
			Minor:0,
			Bump:2,
		},
		File2String:func(path string)(string){
			content,err:=ioutil.ReadFile(path)
			if err!=nil{
				panic(vm.ToValue("Failed to open file"))
			}
			text:=string(content)
			return text
		},
		String2File:func(path string,contents string){
				file, err := os.Create(path)
				if err != nil {
					panic(vm.ToValue("Failed to create file"))
				}
				defer file.Close()
				_, err = io.WriteString(file,contents)
				if err != nil {
					panic(vm.ToValue("Failed to write to file"))
				}
				if file.Sync()!=nil{
					panic(vm.ToValue("Failed to sync file"))
				}
		},
	})
}

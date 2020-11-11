package jsext
import(
	"github.com/dop251/goja"
)
func Register(vm*goja.Runtime){
	Register_jsext_fsutils(vm)
}

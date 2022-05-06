package mlgoja

import (
	"github.com/evocert/kwe/iorw/active"
	_ "github.com/evocert/kwe/ml/tenserflow"
	"github.com/evocert/kwe/resources"
)

func init() {
	active.LoadGlobalModule("tf.js", resources.GLOBALRSNG().FS().CAT("/tflw/js/tf.min.js"))
}

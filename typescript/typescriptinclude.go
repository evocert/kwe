package typescript

import "github.com/evocert/kwe/resources"

func init() {
	gblrs := resources.GLOBALRSNG()
	gblrs.FS().MKDIR("/typescript", "")
	gblrs.FS().SET("/typescript/typescript.js", TypescriptJS())
}

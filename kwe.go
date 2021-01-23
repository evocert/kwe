package main

import (
	"os"

	_ "github.com/evocert/kwe/database/mysql"
	_ "github.com/evocert/kwe/database/ora"
	_ "github.com/evocert/kwe/database/postgres"
	_ "github.com/evocert/kwe/database/sqlserver"
	"github.com/evocert/kwe/service"
)

func main() {
	/*if cmd, cmderr := osprc.NewCommand("cmd"); cmderr == nil {
		for {
			if rs, _ := cmd.Readln(); rs != "" {
				fmt.Println(rs)
			} else {
				break
			}
		}
		cmd.Println("go get -u -v github.com/dop251/goja")
		for {
			if rs, _ := cmd.Readln(); rs != "" {
				fmt.Println(rs)
			} else {
				break
			}
		}
		fmt.Println(cmd.Dir())
	}*/
	service.RunService(os.Args...)
}

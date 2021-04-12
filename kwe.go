package main

import (
	"os"

	//_ "github.com/evocert/kwe/database/db2"

	_ "github.com/evocert/kwe/database/mysql"
	_ "github.com/evocert/kwe/database/ora"
	_ "github.com/evocert/kwe/database/postgres"
	_ "github.com/evocert/kwe/database/sqlserver"
	"github.com/evocert/kwe/service"
)

func main() {
	/*if cmd, cmderr := osprc.NewCommand("cmd", "-"); cmderr == nil {
		for {
			if rs, _ := cmd.Readln(); rs != "" {
				fmt.Println(rs)
			} else {
				break
			}
		}
		cmd.Println("whoami")
		for {
			if rs, _ := cmd.ReadAll(); rs != "" {
				fmt.Println(rs)
			} else {
				break
			}
		}
		cmd.Println("date")
		for {
			if rs, _ := cmd.ReadAll(); rs != "" {
				fmt.Println(rs)
			} else {
				break
			}
		}
		cmd.Println("time")
		for {
			if rs, _ := cmd.ReadAll(); rs != "" {
				fmt.Println(rs)
			} else {
				break
			}
		}
		fmt.Println(cmd.Dir())
		cmd.Close()
	}*/
	/*go func() {
		time.Sleep(time.Second * 5)
		cnt := float64(0)

		strtm := time.Now()
		tlrqst := float64(1000)

		for cnt < tlrqst {
			chnls.GLOBALCHNL().DefaultServeRW(nil, "/mem/test.hl", nil)
			cnt++
		}
		s := time.Now().Sub(strtm).Seconds()
		fmt.Println(s)
		fmt.Println()
		fmt.Println((tlrqst / s))
		fmt.Println()
	}()*/
	service.RunService(os.Args...)

}

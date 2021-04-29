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

	/*mphndl := caching.NewMapHandler()

	mphndl.Put("mp1", map[string]interface{}{"ka": 7879879, "kb": "hjgghhjg"})
	if vfound := mphndl.Find("mp1"); vfound != nil {
		fmt.Println(vfound)
	}

	fmt.Println(mphndl)

	//strt := time.Now()

	for i := 0; i < 30000; i++ {
		k := fmt.Sprintf("k%d", i)
		mphndl.Put(k, 80*(i+1))
	}
	//fmt.Println(time.Now().Sub(strt).Milliseconds(), ":ADD VALUES")
	for i := 0; i < 30000; i++ {
		k := fmt.Sprintf("k%d", i)
		mphndl.Put(k, 60*(i+1))
	}

	for i := 0; i < 30000; i++ {
		k := fmt.Sprintf("k%d", i)
		k1 := fmt.Sprintf("s%d", (i + 1))
		mphndl.ReplaceKey(k, k1)
	}

	//fmt.Println(time.Now().Sub(strt).Milliseconds(), ":REPACE KEYS")

	for _, k := range mphndl.Keys() {
		if k != "" {
			continue
		}
	}

	for _, k := range mphndl.Values() {
		if k != "" {
			continue
		}
	}

	//fmt.Println(time.Now().Sub(strt).Milliseconds(), ":READ Values")

	mphndl.Clear()

	for i := 0; i < 30000; i++ {
		k := fmt.Sprintf("s%d", (i + 1))
		//lst.Add(nil, nil, k)
		mphndl.Remove(k)
		//iorw.Fprintln(os.Stdout, k, ":", mphndl.Find(k))
	}

	//fmt.Println(time.Now().Sub(strt).Milliseconds(), ":REMOVED ALL")
	*/
	service.RunService(os.Args...)
}

package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/evocert/kwe/caching"
)

func main() {
	var mp caching.MapAPI = caching.NewMapHandler(caching.NewMap())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		i := 0
		for {
			go func(key string) {
				d := mp.Find(key)
				v, _ := d.(int)
				v++
				if v > 100 {
					v = 1
				}
				mp.Put(key, v)
			}(fmt.Sprintf("TEST%d", i))
			i++
			if i >= 10 {
				i = 0
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	go func() {
		i := 0
		for {
			go func(key string) {
				mp.Remove(key)
			}(fmt.Sprintf("TEST%d", i))
			i++
			if i >= 10 {
				i = 0
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	go func() {
		i := 0
		for {
			go func() {
				fmt.Println(mp.String())
			}()
			i++
			if i >= 10 {
				i = 0
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	go func() {
		i := 0
		for {
			mp.Clear()
			time.Sleep(2 * time.Second)
			i++
			if i >= 10 {
				mp.Close()
				i = 0
			}
		}
	}()

	go func() {
		time.Sleep(10 * time.Minute)
		wg.Done()
	}()
	wg.Wait()
}

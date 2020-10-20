package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func test() {
	resp, err := http.Get("https://api.github.com/users/tensorflow")
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	fmt.Print(string(body))
}
func main() {
	fmt.Print("#Hello#")
	test()
}

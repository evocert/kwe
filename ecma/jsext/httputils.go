package jsext

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func Register_jsext_httputils(lclobjmp map[string]interface{}) {
	if lclobjmp == nil {
		return
	}
	//vm.SetFieldNameMapper(goja.TagFieldNameMapper("json",true))
	type Version struct {
		Major int `json:"major"`
		Minor int `json:"minor"`
		Bump  int `json:"bump"`
	}
	//todo: namespace everything kwe.fsutils.etcetcetc
	//first test for kwe then do set fsutils on kwe
	lclobjmp["httputils"] = struct {
		Version Version                             `json:"version"`
		Get     func(string) string                 `json:"get"`
		Post    func(string, string, string) string `json:"post"`
	}{
		Version: Version{
			Major: 0,
			Minor: 0,
			Bump:  2,
		},
		Get: func(path string) string {
			resp, err := http.Get(path)
			if err != nil {
				panic("Failed to open url")
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic("Failed to read from url")
			}
			return string(body)
		},
		Post: func(path string, contenttype string, body string) string {

			resp, err := http.Post(path, contenttype, bytes.NewBufferString(body))
			if err != nil {
				panic("Failed to open url")
			}
			defer resp.Body.Close()
			resbody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic("Failed to read from url")
			}
			return string(resbody)

		},
	}
}

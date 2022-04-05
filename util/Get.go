package util

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func splitVmess(s string) []string {
	var sep = map[rune]bool{
		' ':  true,
		'\n': true,
		',':  true,
		';':  true,
		'\t': true,
		'\f': true,
		'\v': true,
		'\r': true,
	}
	return strings.FieldsFunc(s, func(r rune) bool {
		return sep[r]
	})
}

// GetVmessList 从某个URL中读取Vmess列表
func GetVmessList(url string) ([]string, error) {
	resp, err := http.Get(url) //请求base64Vmess
	if err != nil {
		return []string{}, err
	}
	base64Vmess, err := ioutil.ReadAll(resp.Body) //读取base64Vmess
	if err != nil {
		return nil, err
	}
	strVmess, err := Base64VmessListDecode(string(base64Vmess)) //解码base64Vmess为strVmess
	if err != nil {
		return nil, err
	}
	return splitVmess(strVmess), nil //分割strVmess
}

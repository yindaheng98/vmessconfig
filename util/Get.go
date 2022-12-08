package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os/exec"
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

// defalutGetVmessList 默认的GetVmessList
func defalutGetVmessList(url string) ([]string, error) {
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

// GetVmessList 从某个URL中读取Vmess列表，可自定义
var GetVmessList = defalutGetVmessList

// WgetGetVmessList 通过wget读取Vmess列表的GetVmessList
// http.Get 没法直接读取 octet-stream 遂出此下策
func WgetGetVmessList(url string) ([]string, error) {
	cmd := exec.Command("wget", "-O", "-", url)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return []string{}, err
	}
	strVmess, err := Base64VmessListDecode(stdout.String()) //解码base64Vmess为strVmess
	if err != nil {
		return nil, err
	}
	return splitVmess(strVmess), nil //分割strVmess
}

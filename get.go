package vmessconfig

import (
	"bytes"
	"github.com/yindaheng98/vmessconfig/util"
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

// DefalutGetVmessList 默认的GetVmessList
func DefalutGetVmessList(url string) ([]string, error) {
	resp, err := http.Get(url) //请求base64Vmess
	if err != nil {
		return []string{}, err
	}
	base64Vmess, err := ioutil.ReadAll(resp.Body) //读取base64Vmess
	if err != nil {
		return nil, err
	}
	strVmess, err := util.Base64VmessListDecode(string(base64Vmess)) //解码base64Vmess为strVmess
	if err != nil {
		return nil, err
	}
	return splitVmess(strVmess), nil //分割strVmess
}

// getVmessList 从某个URL中读取Vmess列表，可自定义
var getVmessList = DefalutGetVmessList

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
	strVmess, err := util.Base64VmessListDecode(stdout.String()) //解码base64Vmess为strVmess
	if err != nil {
		return nil, err
	}
	return splitVmess(strVmess), nil //分割strVmess
}

func CustomizeGetVmessList(f func(url string) ([]string, error)) {
	getVmessList = f
}

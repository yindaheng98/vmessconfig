package util

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetVmessList(t *testing.T) {
	vl, err := GetVmessList("https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range vl {
		fmt.Println(v)
	}
	vcl := VmessListParse(vl, false, false)
	for vm, vc := range vcl {
		fmt.Println(vm)
		j, err := json.Marshal(vc)
		if err != nil {
			t.Error(err)
			continue
		}
		fmt.Println(string(j))
	}
}

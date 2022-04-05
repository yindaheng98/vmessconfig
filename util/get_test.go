package util

import (
	"fmt"
	"testing"
)

func TestGetVmessList(t *testing.T) {
	vl, err := GetVmessList("https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex")
	if err == nil {
		for _, v := range vl {
			fmt.Println(v)
		}
	}
}

package vmessconfig

import (
	"context"
	"fmt"
	"github.com/yindaheng98/vmessconfig/util"
	"testing"
)

func TestVmessPingAll(t *testing.T) {
	vmesslist, err := getVmessList("https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex")
	if err != nil {
		t.Error(err)
	}
	vmessstats := VmessPingAll(vmesslist, DefaultPingConfig(), context.Background())
	vmesss := util.VmessSort(vmessstats)
	for _, vmessstr := range vmesss {
		fmt.Println(vmessstr)
		fmt.Printf("%+v\n", vmessstats[vmessstr])
	}
}

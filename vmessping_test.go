package vmessconfig

import (
	"fmt"
	"github.com/yindaheng98/vmessconfig/util"
	"os"
	"testing"
)

func TestVmessPingAll(t *testing.T) {
	vmesslist, err := util.GetVmessList("https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex")
	if err != nil {
		t.Error(err)
	}
	vmessstats := VmessPingAll(vmesslist, DefaultPingConfig(), make(chan os.Signal))
	vmesss := util.VmessSort(vmessstats)
	for _, vmessstr := range vmesss {
		fmt.Println(vmessstr)
		fmt.Printf("%+v\n", vmessstats[vmessstr])
	}
}

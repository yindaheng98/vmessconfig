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
	vmessstats := VmessPingAll(vmesslist, &PingConfig{
		Dest:  "http://www.google.com/gen_204",
		Count: 4, Timeoutsec: 8, Inteval: 1, Quit: 0,
		ShowNode: true, Verbose: false, UseMux: false, AllowInsecure: false,
		Threads: 16,
	}, make(chan os.Signal))
	vmesss := util.VmessSort(vmessstats)
	for _, vmessstr := range vmesss {
		fmt.Println(vmessstr)
		fmt.Printf("%+v\n", vmessstats[vmessstr])
	}
}

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
		dest:  "http://www.google.com/gen_204",
		count: 8, timeoutsec: 4, inteval: 1, quit: 0,
		showNode: true, verbose: false, useMux: false, allowInsecure: false,
	}, 16, make(chan os.Signal))
	for vmessstr, vmessstat := range vmessstats {
		fmt.Println(vmessstr)
		fmt.Printf("%+v\n", vmessstat)
	}
}

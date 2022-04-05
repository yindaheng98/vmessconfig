package vmessconfig

import (
	"encoding/json"
	"fmt"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
	"os"
	"testing"
)

func TestVmessConfigBalancer(t *testing.T) {
	template := &conf.Config{}
	err := json.Unmarshal([]byte(defaultTemplate), template)
	if err != nil {
		t.Error(err)
		return
	}
	bconf := DefaultBalancerConfig()
	vconf, err := VmessConfigBalancer("https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex", template, bconf, make(<-chan os.Signal))
	if err != nil {
		t.Error(err)
		return
	}
	j, err := json.Marshal(vconf)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(j))
}

package vmessconfig

import (
	"encoding/json"
	"fmt"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
	"os"
	"testing"
)

func printVmessConfig(t *testing.T, conf *conf.Config) {
	j, err := json.MarshalIndent(conf, "", " ")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(j))

}

func TestVmessConfigBalancer(t *testing.T) {
	template := &conf.Config{}
	err := json.Unmarshal([]byte(DefaultBalancerTemplate), template)
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
	printVmessConfig(t, vconf)
}

func TestVmessConfigSingleNode(t *testing.T) {
	template := &conf.Config{}
	err := json.Unmarshal([]byte(DefaultSingleNodeTemplate), template)
	if err != nil {
		t.Error(err)
		return
	}
	bconf := DefaultSingleNodeConfig()
	vconf, err := VmessConfigSingleNode("https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex", template, bconf, make(<-chan os.Signal))
	if err != nil {
		t.Error(err)
		return
	}
	printVmessConfig(t, vconf)
}

func TestVmessConfig(t *testing.T) {
	template := &conf.Config{}
	err := json.Unmarshal([]byte(DefaultBalancerTemplate), template)
	if err != nil {
		t.Error(err)
		return
	}
	bconf := DefaultBalancerConfig()
	bconf.PingConfig.Count = 1
	vconf, err := VmessConfig("https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex", template, bconf, make(<-chan os.Signal))
	if err != nil {
		t.Error(err)
		return
	}
	printVmessConfig(t, vconf)

	template = &conf.Config{}
	err = json.Unmarshal([]byte(DefaultSingleNodeTemplate), template)
	if err != nil {
		t.Error(err)
		return
	}
	sconf := DefaultSingleNodeConfig()
	sconf.PingConfig.Count = 1
	vconf, err = VmessConfig("https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex", template, sconf, make(<-chan os.Signal))
	if err != nil {
		t.Error(err)
		return
	}
	printVmessConfig(t, vconf)
}

package vmessconfig

import (
	"context"
	"fmt"
	"testing"
)

func TestVmessConfigBalancer(t *testing.T) {
	template := []byte(DefaultBalancerTemplate)
	bconf := DefaultBalancerConfig()
	vconf, err := VmessConfigBalancer([]string{
		"https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex",
		"https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex",
	}, template, bconf, context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(vconf))
}

func TestVmessWgetConfigBalancer(t *testing.T) {
	template := []byte(DefaultBalancerTemplate)
	bconf := DefaultBalancerConfig()
	CustomizeGetVmessList(WgetGetVmessList)
	vconf, err := VmessConfigBalancer([]string{
		"https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex",
		"https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex",
	}, template, bconf, context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(vconf))
}

func TestVmessConfigSingleNode(t *testing.T) {
	template := []byte(DefaultSingleNodeTemplate)
	bconf := DefaultSingleNodeConfig()
	vconf, err := VmessConfigSingleNode([]string{
		"https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex",
		"https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex",
	}, template, bconf, context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(vconf))
}

func TestVmessConfig(t *testing.T) {
	template := []byte(DefaultBalancerTemplate)
	bconf := DefaultBalancerConfig()
	bconf.PingConfig.Count = 1
	vconf, err := VmessConfig([]string{
		"https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex",
		"https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex",
	}, template, bconf, context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(vconf))

	template = []byte(DefaultSingleNodeTemplate)
	sconf := DefaultSingleNodeConfig()
	sconf.PingConfig.Count = 1
	vconf, err = VmessConfig([]string{
		"https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex",
		"https://get.cloudv2.net/osubscribe.php?sid=128958&token=MDByRw64Cnex",
	}, template, sconf, context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(vconf))
}

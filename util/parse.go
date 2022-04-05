package util

import (
	"fmt"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
	"github.com/v2fly/vmessping/miniv2ray"
	"github.com/v2fly/vmessping/vmess"
)

func VmessParse(vms string, useMux, allowInsecure bool) (*conf.OutboundDetourConfig, error) {
	vml, err := vmess.ParseVmess(vms)
	if err != nil {
		return nil, err
	}
	outbound, err := miniv2ray.Vmess2OutboundDetour(vml, useMux, allowInsecure, &conf.OutboundDetourConfig{})
	if err != nil {
		return nil, err
	}
	return outbound, nil
}

func VmessListParse(vmesslist []string, useMux, allowInsecure bool) []*conf.OutboundDetourConfig {
	var outbounds []*conf.OutboundDetourConfig
	for _, vms := range vmesslist {
		outbound, err := VmessParse(vms, useMux, allowInsecure)
		if err != nil {
			fmt.Printf("Cannot parse :%s\n%+v\n", vms, err)
			continue
		}
		outbounds = append(outbounds, outbound)
	}
	return outbounds
}

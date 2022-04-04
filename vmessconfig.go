package vmessconfig

import (
	core "github.com/v2fly/v2ray-core/v4"
	"github.com/v2fly/vmessping"
	"github.com/v2fly/vmessping/miniv2ray"
	"github.com/v2fly/vmessping/vmess"
	"github.com/yindaheng98/vmessconfig/util"
)

type Outbound struct {
	config *core.OutboundHandlerConfig
	stats  *vmessping.PingStat
}

func VmessOutboundConfig(url string, template *core.OutboundHandlerConfig, useMux, allowInsecure bool) ([]*Outbound, error) {
	vl, err := util.GetVmessList(url)
	if err != nil {
		return nil, err
	}
	outbounds := make([]*Outbound, len(vl))
	for i := 0; i < len(vl); i++ {
		vmess, err := vmess.ParseVmess(vl[i])
		if err != nil {
			return nil, err
		}
		pingstats, err := vmessping.Ping(vl[i])
		if err != nil {
			return nil, err
		}
		o, err := miniv2ray.Vmess2Outbound(vmess, useMux, allowInsecure)
		if err != nil {
			return nil, err
		}
		config, err := util.ConfigMerge(o, template)
		if err != nil {
			return nil, err
		}
		outbounds[i] = &Outbound{config, pingstats}
	}
	return outbounds, nil
}

package vmessconfig

import (
	"github.com/remeh/sizedwaitgroup"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
	"github.com/v2fly/vmessping"
	"os"
)

type PingConfig struct {
	dest                                     string
	count, timeoutsec, inteval, quit         uint
	showNode, verbose, useMux, allowInsecure bool
}

func VmessPingOne(vmessstr string, pingconfig *PingConfig, stopCh <-chan os.Signal) (*vmessping.PingStat, error) {
	pingstat, err := vmessping.Ping(
		vmessstr,
		pingconfig.count,
		pingconfig.dest,
		pingconfig.timeoutsec,
		pingconfig.inteval,
		pingconfig.quit,
		stopCh,
		pingconfig.showNode,
		pingconfig.verbose,
		pingconfig.useMux,
		pingconfig.allowInsecure,
	)
	if err != nil {
		return nil, err
	}
	return pingstat, nil
}

type Outbound struct {
	config *conf.OutboundDetourConfig
	stats  *vmessping.PingStat
}

func VmessPingAll(vmesslist []string, pingconfig *PingConfig, threads int, stopCh <-chan os.Signal) map[string]*vmessping.PingStat {
	vmessstats := make(map[string]*vmessping.PingStat)
	wg := sizedwaitgroup.New(threads)
	for _, vmessstr := range vmesslist {
		wg.Add()
		go func(vmessstr string) {
			pingstat, err := VmessPingOne(vmessstr, pingconfig, stopCh)
			if err == nil {
				vmessstats[vmessstr] = pingstat
			}
			wg.Done()
		}(vmessstr)
	}
	return vmessstats
}

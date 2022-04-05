package vmessconfig

import (
	"fmt"
	"github.com/remeh/sizedwaitgroup"
	"github.com/v2fly/vmessping"
	"os"
)

type PingConfig struct {
	Dest          string `json:"destination"`
	Count         uint   `json:"count"`
	Timeoutsec    uint   `json:"timeout"`
	Inteval       uint   `json:"interval"`
	Quit          uint   `json:"quit"`
	ShowNode      bool   `json:"showNode"`
	Verbose       bool   `json:"verbose"`
	UseMux        bool   `json:"useMux"`
	AllowInsecure bool   `json:"allowInsecure"`
	Threads       int    `json:"threads"`
}

func VmessPingOne(vmessstr string, pingconfig *PingConfig, stopCh <-chan os.Signal) (*vmessping.PingStat, error) {
	pingstat, err := vmessping.Ping(
		vmessstr,
		pingconfig.Count,
		pingconfig.Dest,
		pingconfig.Timeoutsec,
		pingconfig.Inteval,
		pingconfig.Quit,
		stopCh,
		pingconfig.ShowNode,
		pingconfig.Verbose,
		pingconfig.UseMux,
		pingconfig.AllowInsecure,
	)
	if err != nil {
		return nil, err
	}
	return pingstat, nil
}

func VmessPingAll(vmesslist []string, pingconfig *PingConfig, stopCh <-chan os.Signal) map[string]*vmessping.PingStat {
	vmessstats := make(map[string]*vmessping.PingStat)
	wg := sizedwaitgroup.New(pingconfig.Threads)
	for _, vmessstr := range vmesslist {
		wg.Add()
		go func(vmessstr string) {
			pingstat, err := VmessPingOne(vmessstr, pingconfig, stopCh)
			if err != nil {
				fmt.Printf("Cannot stat :%s\n%+v\n", vmessstr, err)
			} else {
				vmessstats[vmessstr] = pingstat
			}
			wg.Done()
		}(vmessstr)
	}
	wg.Wait()
	return vmessstats
}

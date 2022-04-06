package vmessconfig

import (
	"context"
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
	Threads       uint   `json:"threads"`
}

func DefaultPingConfig() *PingConfig {
	return &PingConfig{
		Dest:  "http://www.google.com/gen_204",
		Count: 4, Timeoutsec: 8, Inteval: 1, Quit: 0,
		ShowNode: true, Verbose: false, UseMux: false, AllowInsecure: false,
		Threads: 16,
	}
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

func VmessPingOneWithContext(vmessstr string, pingconfig *PingConfig, ctx context.Context) (*vmessping.PingStat, error) {
	statCh := make(chan *vmessping.PingStat)
	erroCh := make(chan error)
	stopCh := make(chan os.Signal)
	go func(statCh chan<- *vmessping.PingStat, erroCh chan<- error, stopCh <-chan os.Signal) {
		vmessstat, err := VmessPingOne(vmessstr, pingconfig, stopCh)
		if err != nil {
			erroCh <- err
		} else {
			statCh <- vmessstat
		}
	}(statCh, erroCh, stopCh)
	select {
	case <-ctx.Done():
		close(stopCh)
		return nil, fmt.Errorf("Context exited :%s\n", vmessstr)
	case vmessstat := <-statCh:
		return vmessstat, nil
	case err := <-erroCh:
		return nil, err
	}
}

func VmessPingAll(vmesslist []string, pingconfig *PingConfig, ctx context.Context) map[string]*vmessping.PingStat {
	vmessstats := make(map[string]*vmessping.PingStat)
	wg := sizedwaitgroup.New(int(pingconfig.Threads))
	for _, vmessstr := range vmesslist {
		err := wg.AddWithContext(ctx)
		if err != nil {
			return vmessstats
		}
		go func(vmessstr string) {
			pingstat, err := VmessPingOneWithContext(vmessstr, pingconfig, ctx)
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

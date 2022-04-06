package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
	"github.com/yindaheng98/vmessconfig"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

var (
	verbose       = flag.Bool("v", false, "verbose (debug log)")
	showNode      = flag.Bool("n", false, "show node location/outbound ip")
	useMux        = flag.Bool("m", false, "use mux outbound")
	allowInsecure = flag.Bool("allow-insecure", false, "allow insecure TLS connections")
	desturl       = flag.String("dest", "http://www.google.com/gen_204", "the test destination url, need 204 for success return")
	count         = flag.Uint("c", 4, "Count. Stop after sending COUNT requests")
	timeout       = flag.Uint("o", 10, "timeout seconds for each request")
	inteval       = flag.Uint("i", 1, "inteval seconds between pings")
	quit          = flag.Uint("q", 0, "fast quit on error counts")

	threads                 = flag.Uint("k", 8, "Threads. At most THREADS pinging at the same time")
	outboundInsertBeforeTag = flag.String("outbound-insert-before", "vmessconfig-outbound-insert", "OutboundInsertBeforeTag. insert outbound before the exists outbound whose tag is this")
	tagFormat               = flag.String("tag-format", "vmessconfig-autogenerated-%d", "TagFormat. format of the auto-generated outbounds' tag")
	balancerInsertToTag     = flag.String("balancer-insert-into", "vmessconfig-autogenerated-balancer", "BalancerInsertToTag. insert the balancer.selector into the balancer whose tag is this")
	maxSelect               = flag.Uint("max-select", 8, "MaxSelect. how many outbounds do you want to put into")

	tempath = flag.String("t", "", "Path to your template")
	writeto = flag.String("w", "-", "Path to write the auto-generated v2ray config json")
)

func exit(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Println(os.Args[0], "single|balancer http|https:// ...")
	flag.Usage()
	os.Exit(1)
}

func main() {
	flag.Parse()

	var mode string
	var urls []string
	if flag.NArg() < 2 {
		exit(nil)
	}
	mode = flag.Args()[0]
	urls = flag.Args()[1:]

	template := &conf.Config{}
	if *tempath != "" {
		data, err := ioutil.ReadFile(*tempath)
		if err != nil {
			exit(err)
		}
		err = json.Unmarshal(data, template)
		if err != nil {
			exit(err)
		}
	} else {
		if mode == "single" {
			err := json.Unmarshal([]byte(vmessconfig.DefaultSingleNodeTemplate), template)
			if err != nil {
				exit(err)
			}
		} else if mode == "balancer" {
			err := json.Unmarshal([]byte(vmessconfig.DefaultBalancerTemplate), template)
			if err != nil {
				exit(err)
			}
		} else {
			exit(nil)
		}
	}

	pingconfig := &vmessconfig.PingConfig{
		Dest:          *desturl,
		Count:         *count,
		Timeoutsec:    *timeout,
		Inteval:       *inteval,
		Quit:          *quit,
		ShowNode:      *showNode,
		Verbose:       *verbose,
		UseMux:        *useMux,
		AllowInsecure: *allowInsecure,
		Threads:       *threads,
	}

	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)

	var config interface{}
	if mode == "single" {
		config = &vmessconfig.SingleNodeConfig{
			BaseConfig: vmessconfig.BaseConfig{
				PingConfig: pingconfig,
			},
			OutboundInsertBeforeTag: *outboundInsertBeforeTag,
		}
	} else if mode == "balancer" {
		config = &vmessconfig.BalancerConfig{
			BaseConfig: vmessconfig.BaseConfig{
				PingConfig: pingconfig,
			},
			OutboundInsertBeforeTag: *outboundInsertBeforeTag,
			TagFormat:               *tagFormat,
			BalancerInsertToTag:     *balancerInsertToTag,
			MaxSelect:               *maxSelect,
		}
	} else {
		exit(nil)
	}
	result, err := vmessconfig.VmessConfig(urls, template, config, ctx)
	if err != nil {
		exit(err)
	}
	j, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		exit(err)
	}
	if *writeto != "" && *writeto != "-" {
		err := ioutil.WriteFile(*writeto, j, 0777)
		if err != nil {
			exit(err)
		}
	} else {
		fmt.Println(string(j))
	}
}
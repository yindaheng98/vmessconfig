package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
	"github.com/yindaheng98/vmessconfig"
	"github.com/yindaheng98/vmessconfig/cmd/args"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

func exit(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%s [balancer|single] -urls https://... -urls https://...\n", os.Args[0])
	args.BalancerFlagSet.Usage()
	args.SingleFlagSet.Usage()
	os.Exit(1)
}

func main() {
	err := args.GenerateCmdArgs()
	if err != nil {
		exit(err)
	}
	config, err := args.ParseCmdArgs(os.Args[1:])
	if err != nil {
		exit(err)
	}

	template := &conf.Config{}
	if config.TemplateConfig.From == "" {
		err = json.Unmarshal([]byte(config.TemplateConfig.DefaultTemplate()), template)
		if err != nil {
			exit(err)
		}
	} else {
		data, err := ioutil.ReadFile(config.TemplateConfig.From)
		if err != nil {
			exit(err)
		}
		err = json.Unmarshal(data, template)
		if err != nil {
			exit(err)
		}
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)

	v2config, err := vmessconfig.VmessConfig(config.Urls, template, config.Config, ctx)

	j, err := json.MarshalIndent(v2config, "", " ")
	if err != nil {
		exit(err)
	}
	if config.TemplateConfig.To != "" && config.TemplateConfig.To != "-" {
		err := ioutil.WriteFile(config.TemplateConfig.To, j, 0777)
		if err != nil {
			exit(err)
		}
	} else {
		fmt.Println(string(j))
	}
}

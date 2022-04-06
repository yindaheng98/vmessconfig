package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/octago/sflags/gen/gflag"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
	"github.com/yindaheng98/vmessconfig"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

type TemplateConfig struct {
	From string `desc:"Where the template file is"`
	To   string `desc:"Where the v2ray json config file should write to"`
}

type BalancerCmdConfig struct {
	*vmessconfig.BalancerConfig
	TemplateConfig *TemplateConfig
	Urls           []string `desc:"List of your subscription urls"`
}

type SingleCmdConfig struct {
	*vmessconfig.SingleNodeConfig
	TemplateConfig *TemplateConfig
	Urls           []string `desc:"List of your subscription urls"`
}

var (
	templateConfig = &TemplateConfig{
		From: "",
		To:   "",
	}

	balancerFlagSet   = flag.NewFlagSet("balancer", flag.ExitOnError)
	balancerCmdConfig = &BalancerCmdConfig{
		BalancerConfig: vmessconfig.DefaultBalancerConfig(),
		TemplateConfig: templateConfig,
		Urls:           []string{},
	}

	singleFlagSet   = flag.NewFlagSet("single", flag.ExitOnError)
	singleCmdConfig = &SingleCmdConfig{
		SingleNodeConfig: vmessconfig.DefaultSingleNodeConfig(),
		TemplateConfig:   templateConfig,
		Urls:             []string{},
	}
)

func exit(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%s [balancer|single] -urls https://... -urls https://...\n", os.Args[0])
	balancerFlagSet.Usage()
	singleFlagSet.Usage()
	os.Exit(1)
}

func main() {
	var err error

	err = gflag.ParseTo(balancerCmdConfig, balancerFlagSet)
	if err != nil {
		exit(err)
	}
	err = gflag.ParseTo(singleCmdConfig, singleFlagSet)
	if err != nil {
		exit(err)
	}

	if len(os.Args) < 2 {
		exit(nil)
	}

	template := &conf.Config{}
	switch os.Args[1] {
	case "balancer":
		err := balancerFlagSet.Parse(os.Args[2:])
		if err != nil {
			exit(err)
		}
		err = json.Unmarshal([]byte(vmessconfig.DefaultBalancerTemplate), template)
		if err != nil {
			exit(err)
		}

	case "single":
		err := singleFlagSet.Parse(os.Args[2:])
		if err != nil {
			exit(err)
		}
		err = json.Unmarshal([]byte(vmessconfig.DefaultSingleNodeTemplate), template)
		if err != nil {
			exit(err)
		}
	default:
		exit(nil)
	}

	if templateConfig.From != "" {
		data, err := ioutil.ReadFile(templateConfig.From)
		if err != nil {
			exit(err)
		}
		err = json.Unmarshal(data, template)
		if err != nil {
			exit(err)
		}
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)

	var v2config *conf.Config
	switch os.Args[1] {
	case "balancer":
		v2config, err = vmessconfig.VmessConfig(balancerCmdConfig.Urls, template, balancerCmdConfig.BalancerConfig, ctx)
		if err != nil {
			exit(err)
		}
	case "single":
		v2config, err = vmessconfig.VmessConfig(singleCmdConfig.Urls, template, singleCmdConfig.SingleNodeConfig, ctx)
		if err != nil {
			exit(err)
		}
	default:
		exit(nil)
	}

	j, err := json.MarshalIndent(v2config, "", " ")
	if err != nil {
		exit(err)
	}
	if templateConfig.To != "" && templateConfig.To != "-" {
		err := ioutil.WriteFile(templateConfig.To, j, 0777)
		if err != nil {
			exit(err)
		}
	} else {
		fmt.Println(string(j))
	}
}

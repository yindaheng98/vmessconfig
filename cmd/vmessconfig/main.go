package main

import (
	"context"
	"fmt"
	"github.com/yindaheng98/vmessconfig"
	"github.com/yindaheng98/vmessconfig/cmd/args"
	"os"
	"os/signal"
	"syscall"
)

func exit(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%s [balancer|single] -urls https://... -urls https://...\n", os.Args[0])
	args.PrintUsage()
	os.Exit(1)
}

func main() {
	config := args.NewCmdConfig()
	err := config.GenerateCmdArgs()
	if err != nil {
		exit(err)
	}
	err = config.ParseCmdArgs(os.Args[1:])
	if err != nil {
		exit(err)
	}

	template, err := config.TemplateConfig.Template()
	if err != nil {
		exit(err)
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)

	v2config, err := vmessconfig.VmessConfig(config.Urls, template, config.Config, ctx)
	if err != nil {
		exit(err)
	}

	err = config.TemplateConfig.Write(v2config)
	if err != nil {
		exit(err)
	}
}

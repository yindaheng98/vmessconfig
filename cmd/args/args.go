package args

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/octago/sflags/gen/gflag"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
	"github.com/yindaheng98/vmessconfig"
	"io/ioutil"
)

type TemplateConfig struct {
	From        string `desc:"Where the template file is"`
	To          string `desc:"Where the v2ray json config file should write to"`
	defaultTemp string
}

func (c TemplateConfig) DefaultTemplate() string {
	return c.defaultTemp
}

func (c TemplateConfig) Template() (*conf.Config, error) {
	template := &conf.Config{}
	if c.From == "" {
		err := json.Unmarshal([]byte(c.DefaultTemplate()), template)
		if err != nil {
			return nil, err
		}
	} else {
		data, err := ioutil.ReadFile(c.From)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, template)
		if err != nil {
			return nil, err
		}
	}
	return template, nil
}

type CmdConfig struct {
	vmessconfig.Config `flag:"-"`
	TemplateConfig     *TemplateConfig
	Urls               []string `desc:"List of your subscription urls"`
}

var (
	BalancerFlagSet   = flag.NewFlagSet("balancer", flag.ExitOnError)
	balancerCmdConfig = &CmdConfig{
		Config: vmessconfig.DefaultBalancerConfig(),
		TemplateConfig: &TemplateConfig{
			From:        "",
			To:          "",
			defaultTemp: vmessconfig.DefaultBalancerTemplate,
		},
		Urls: []string{},
	}

	SingleFlagSet   = flag.NewFlagSet("single", flag.ExitOnError)
	singleCmdConfig = &CmdConfig{
		Config: vmessconfig.DefaultSingleNodeConfig(),
		TemplateConfig: &TemplateConfig{
			From:        "",
			To:          "",
			defaultTemp: vmessconfig.DefaultSingleNodeTemplate,
		},
		Urls: []string{},
	}
)

func GenerateCmdArgs() error {
	errb := genBalancerCmdArgs()
	errs := genSingleCmdArgs()
	if errb != nil && errs != nil {
		return nil
	} else if errb != nil {
		return errb
	} else if errs != nil {
		return errs
	}
	return nil
}

func genBalancerCmdArgs() error {
	err := gflag.ParseTo(balancerCmdConfig.Config, BalancerFlagSet)
	if err != nil {
		return err
	}
	err = gflag.ParseTo(balancerCmdConfig, BalancerFlagSet)
	if err != nil {
		return err
	}
	return nil
}

func genSingleCmdArgs() error {
	err := gflag.ParseTo(singleCmdConfig.Config, SingleFlagSet)
	if err != nil {
		return err
	}
	err = gflag.ParseTo(singleCmdConfig, SingleFlagSet)
	if err != nil {
		return err
	}
	return nil
}

func ParseCmdArgs(args []string) (*CmdConfig, error) {
	if args[0] == "balancer" {
		err := BalancerFlagSet.Parse(args[1:])
		if err != nil {
			return nil, err
		}
		return balancerCmdConfig, nil
	} else if args[0] == "single" {
		err := SingleFlagSet.Parse(args[1:])
		if err != nil {
			return nil, err
		}
		return singleCmdConfig, nil
	} else {
		return nil, fmt.Errorf("INVALID ARGS: %+v", args)
	}
}

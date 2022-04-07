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

func (c TemplateConfig) Write(v2config *conf.Config) error {
	j, err := json.MarshalIndent(v2config, "", " ")
	if err != nil {
		return err
	}
	if c.To != "" && c.To != "-" {
		err := ioutil.WriteFile(c.To, j, 0777)
		if err != nil {
			return err
		}
	} else {
		fmt.Println(string(j))
	}
	return nil
}

type CmdConfig struct {
	vmessconfig.Config `flag:"-"`
	TemplateConfig     *TemplateConfig
	Urls               []string `desc:"List of your subscription urls"`
	balancerConfig     *vmessconfig.BalancerConfig
	singleNodeConfig   *vmessconfig.SingleNodeConfig
}

func NewCmdConfig() *CmdConfig {
	return &CmdConfig{
		balancerConfig:   vmessconfig.DefaultBalancerConfig(),
		singleNodeConfig: vmessconfig.DefaultSingleNodeConfig(),
	}
}

var (
	BalancerFlagSet = flag.NewFlagSet("balancer", flag.ExitOnError)
	SingleFlagSet   = flag.NewFlagSet("single", flag.ExitOnError)
)

func (config *CmdConfig) GenerateCmdArgs() error {
	errb := gflag.ParseTo(config, BalancerFlagSet)
	if errb != nil {
		return errb
	}
	errb = gflag.ParseTo(config.balancerConfig, BalancerFlagSet)
	if errb != nil {
		return errb
	}

	errc := gflag.ParseTo(config, SingleFlagSet)
	if errc != nil {
		return errc
	}
	errc = gflag.ParseTo(config.singleNodeConfig, SingleFlagSet)
	if errc != nil {
		return errc
	}
	return nil
}

func AddCmdArgs(cfg interface{}) error {
	errb := gflag.ParseTo(cfg, BalancerFlagSet)
	if errb != nil {
		return errb
	}
	errc := gflag.ParseTo(cfg, SingleFlagSet)
	if errc != nil {
		return errc
	}
	return nil
}

func PrintUsage() {
	BalancerFlagSet.Usage()
	SingleFlagSet.Usage()
}

func (config *CmdConfig) ParseCmdArgs(args []string) error {
	if args[0] == "balancer" {
		err := BalancerFlagSet.Parse(args[1:])
		if err != nil {
			return err
		}
		config.Config = config.balancerConfig
		config.TemplateConfig.defaultTemp = vmessconfig.DefaultBalancerTemplate
		return nil
	} else if args[0] == "single" {
		err := SingleFlagSet.Parse(args[1:])
		if err != nil {
			return err
		}
		config.Config = config.singleNodeConfig
		config.TemplateConfig.defaultTemp = vmessconfig.DefaultSingleNodeTemplate
		return nil
	} else {
		return fmt.Errorf("INVALID ARGS: %+v", args)
	}
}

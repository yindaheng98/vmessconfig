package util

import (
	"fmt"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
)

func insertOutboundConfig(outboundDetourConfigs []conf.OutboundDetourConfig, template *conf.Config, insertBeforeTag string) *conf.Config {
	for i, outbound := range template.OutboundConfigs {
		if outbound.Tag == insertBeforeTag {
			template.OutboundConfigs = append(append(template.OutboundConfigs[0:i], outboundDetourConfigs...), template.OutboundConfigs[i:]...)
			return template
		}
	}
	template.OutboundConfigs = append(template.OutboundConfigs, outboundDetourConfigs...)
	return template
}

func insertBalancerTags(tags []string, template *conf.Config, insertToTag string) *conf.Config {
	for _, balancer := range template.RouterConfig.Balancers {
		if balancer.Tag == insertToTag {
			balancer.Selectors = append(balancer.Selectors, tags...)
		}
	}
	return template
}

// VmessBalancerConfigMerge 将一系列OutboundDetourConfig写入负载均衡配置的模板
// tagFormat: outbound.Tag的格式
// outboundInsertBeforeTag: 在模板outbound列表的何处插入outbounds配置，找不到位置就插在最后
// balancerInsertToTag: 在模板中的哪个balancer中插入outbounds tag列表，找不到位置就不插
func VmessBalancerConfigMerge(
	outboundDetourConfigs []*conf.OutboundDetourConfig, template *conf.Config,
	tagFormat, outboundInsertBeforeTag, balancerInsertToTag string,
) *conf.Config {
	tags := make([]string, len(outboundDetourConfigs))
	outboundConfigs := make([]conf.OutboundDetourConfig, len(outboundDetourConfigs))
	for i, outbound := range outboundDetourConfigs {
		tag := fmt.Sprintf(tagFormat, i)
		outbound.Tag = tag
		tags[i] = tag
		outboundConfigs[i] = *outbound
	}
	template = insertOutboundConfig(outboundConfigs, template, outboundInsertBeforeTag)
	template = insertBalancerTags(tags, template, balancerInsertToTag)
	return template
}

func VmessSingleNodeConfigMerge(outboundDetourConfig *conf.OutboundDetourConfig, template *conf.Config, outboundInsertBeforeTag string) *conf.Config {
	return insertOutboundConfig([]conf.OutboundDetourConfig{*outboundDetourConfig}, template, outboundInsertBeforeTag)
}

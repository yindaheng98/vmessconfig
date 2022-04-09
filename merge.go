package vmessconfig

import (
	"encoding/json"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
)

type V2Config json.RawMessage
type OutboundDetourConfig json.RawMessage

func insertOutboundConfig(outboundDetourConfigs []json.RawMessage, templateRaw V2Config, insertBeforeTag string) V2Config {
	template := make(map[string]json.RawMessage)  // 模板是JSONObject
	err := json.Unmarshal(templateRaw, &template) // 解码模板
	if err != nil {
		return templateRaw
	}
	outboundsRaw, ok := template["outbounds"] // 从模板中取出outbounds
	var outbounds []json.RawMessage           // outbounds是JSONArray
	if !ok {                                  // 模板中没有outbounds就直接赋值
		outbounds = outboundDetourConfigs
	} else { // 模板中有outbounds就先插入再赋值
		err = json.Unmarshal(outboundsRaw, &outbounds) // 解码outbounds
		if err != nil {
			return templateRaw
		}
		for i, outboundRaw := range outbounds { // 遍历outbounds里的所有outbound
			outbound := make(map[string]json.RawMessage)  // outbound是JSONObject
			err := json.Unmarshal(outboundRaw, &outbound) // 解码outbound
			if err != nil {
				continue
			}
			tagRaw, ok := outbound["tag"] // 从outbound里取出tag
			if !ok {
				continue
			}
			var tag string                     // tag是string
			err = json.Unmarshal(tagRaw, &tag) // 解码tag
			if err != nil {
				continue
			}
			if tag == insertBeforeTag { // 如果正是要找的tag
				out := make([]json.RawMessage, 0) // 就执行插入操作
				out = append(out, outbounds[:i]...)
				out = append(out, outboundDetourConfigs...)
				out = append(out, outbounds[i:]...)
				outbounds = out // 插入完就赋值给outbounds
				break           // 插入完就退出
			}
		}
	}
	outboundsRaw, err = json.MarshalIndent(outbounds, "", " ") // 把outbounds编码回去
	if err != nil {
		return templateRaw
	}
	template["outbounds"] = outboundsRaw                         // 然后赋值给outbounds
	templateRawNew, err := json.MarshalIndent(template, "", " ") // 然后把template编码回去
	if err != nil {
		return templateRaw
	}
	return templateRawNew // 然后返回
}

func insertBalancerTags(tags []string, templateRaw V2Config, insertToTag string) V2Config {
	template := make(map[string]json.RawMessage)  // 模板是JSONObject
	err := json.Unmarshal(templateRaw, &template) // 解码模板
	if err != nil {
		return templateRaw
	}
	routingRaw, ok := template["routing"] // 从模板中取出routing
	if !ok {
		return templateRaw
	}
	var routing map[string]json.RawMessage     // routing是JSONObject
	err = json.Unmarshal(routingRaw, &routing) // 解码routing
	if err != nil {
		return templateRaw
	}
	balancersRaw, ok := routing["balancers"]       // 从模板routing取出balancers
	var balancers []json.RawMessage                // balancers是JSONArray
	err = json.Unmarshal(balancersRaw, &balancers) // 解码balancers
	if err != nil {
		return templateRaw
	}
	for i, balancerRaw := range balancers {
		balancer := make(map[string]json.RawMessage)  // balancer是JSONObject
		err := json.Unmarshal(balancerRaw, &balancer) // 解码balancer
		if err != nil {
			continue
		}
		tagRaw, ok := balancer["tag"] // 从balancer里取出tag
		if !ok {
			continue
		}
		var tag string                     // tag是string
		err = json.Unmarshal(tagRaw, &tag) // 解码tag
		if err != nil {
			continue
		}
		if tag == insertToTag {
			var selector []string                             // selector是string数组
			if selectorRaw, ok := balancer["selector"]; !ok { //如果没有selector就直接赋值
				selector = tags
			} else {
				err = json.Unmarshal(selectorRaw, &selector) // 解码selector
				if err != nil {
					continue
				}
				selector = append(selector, tags...) //然后拼接
			}
			selectorRaw, err := json.MarshalIndent(selector, "", " ") //把selector编码回去
			if err != nil {
				continue
			}
			balancer["selector"] = selectorRaw
			balancerRaw, err = json.MarshalIndent(balancer, "", " ") //把balancer编码回去
			if err != nil {
				continue
			}
			balancers[i] = balancerRaw //赋值给balancers中的原位置
		}
	}
	balancersRaw, err = json.MarshalIndent(balancers, "", " ") //把balancers编码回去
	if err != nil {
		return templateRaw
	}
	routing["balancers"] = balancersRaw
	routingRaw, err = json.MarshalIndent(routing, "", " ") //把routing编码回去
	if err != nil {
		return templateRaw
	}
	template["routing"] = routingRaw
	templateRawNew, err := json.MarshalIndent(template, "", " ") //把template编码回去
	if err != nil {
		return templateRaw
	}
	return templateRawNew
}

// VmessBalancerConfigMerge 将一系列OutboundDetourConfig写入负载均衡配置的模板
// tagFormat: outbound.Tag的格式
// outboundInsertBeforeTag: 在模板outbound列表的何处插入outbounds配置，找不到位置就插在最后
// balancerInsertToTag: 在模板中的哪个balancer中插入outbounds tag列表，找不到位置就不插
func VmessBalancerConfigMerge(
	outboundDetourConfigs []*conf.OutboundDetourConfig, template V2Config,
	outboundInsertBeforeTag, balancerInsertToTag string,
) V2Config {
	tagSet := make(map[string]bool)
	var outboundConfigs []json.RawMessage
	for _, outbound := range outboundDetourConfigs {
		tagSet[outbound.Tag] = true
		outboundRaw, err := json.MarshalIndent(*outbound, "", " ")
		if err != nil {
			continue
		}
		outboundConfigs = append(outboundConfigs, outboundRaw)
	}
	tags := make([]string, len(tagSet))
	i := 0
	for tag := range tagSet {
		tags[i] = tag
		i++
	}
	template = insertOutboundConfig(outboundConfigs, template, outboundInsertBeforeTag)
	template = insertBalancerTags(tags, template, balancerInsertToTag)
	return template
}

func VmessSingleNodeConfigMerge(outboundDetourConfig *conf.OutboundDetourConfig, template V2Config, outboundInsertBeforeTag string) V2Config {
	outboundRaw, err := json.MarshalIndent(*outboundDetourConfig, "", " ")
	if err != nil {
		return template
	}
	return insertOutboundConfig([]json.RawMessage{outboundRaw}, template, outboundInsertBeforeTag)
}

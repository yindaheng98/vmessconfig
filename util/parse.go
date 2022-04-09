package util

import (
	"fmt"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
	"github.com/v2fly/vmessping/miniv2ray"
	"github.com/v2fly/vmessping/vmess"
	"reflect"
)

func getFieldValue(obj *vmess.VmessLink, fieldName string) string {
	v := reflect.ValueOf(*obj).FieldByName(fieldName)
	if v.Type().Name() == "string" {
		return v.String()
	}
	fmt.Println("Supported fields: Add, Aid, Host, ID, Net, Path, Ps, TLS, Type, number.")
	return "[unknown]"
}

func VmessParse(vms string, tagFormat, fieldName string, useMux, allowInsecure bool) (*conf.OutboundDetourConfig, error) {
	vml, err := vmess.ParseVmess(vms)
	if err != nil {
		return nil, err
	}
	outbound, err := miniv2ray.Vmess2OutboundDetour(vml, useMux, allowInsecure, &conf.OutboundDetourConfig{})
	if err != nil {
		return nil, err
	}
	if tagFormat == "fixed" {
		outbound.Tag = fieldName
	} else if fieldName == "number" {
		outbound.Tag = tagFormat
	} else {
		outbound.Tag = fmt.Sprintf(tagFormat, getFieldValue(vml, fieldName))
	}
	return outbound, nil
}

func VmessListParse(vmesslist []string, tagFormat, fieldName string, useMux, allowInsecure bool) []*conf.OutboundDetourConfig {
	var outbounds []*conf.OutboundDetourConfig
	for i, vms := range vmesslist {
		tf := tagFormat
		if fieldName == "number" {
			tf = fmt.Sprintf(tagFormat, fmt.Sprintf("%d", i))
		}
		outbound, err := VmessParse(vms, tf, fieldName, useMux, allowInsecure)
		if err != nil {
			fmt.Printf("Cannot parse :%s\n%+v\n", vms, err)
			continue
		}
		outbounds = append(outbounds, outbound)
	}
	return outbounds
}

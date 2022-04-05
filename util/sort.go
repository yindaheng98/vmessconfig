package util

import (
	"github.com/v2fly/vmessping"
	"sort"
)

// VmessSort 排除连不上的节点并且给剩下的节点按通信质量排序
func VmessSort(vmessstats map[string]*vmessping.PingStat) []string {
	var vmesss []string
	for vmess, vmessstat := range vmessstats {
		if vmessstat.ErrCounter >= vmessstat.ReqCounter {
			continue
		}
		vmesss = append(vmesss, vmess)
	}
	sort.SliceStable(vmesss, func(i, j int) bool {
		if vmessstats[vmesss[i]].ErrCounter == vmessstats[vmesss[j]].ErrCounter {
			return vmessstats[vmesss[i]].AvgMs < vmessstats[vmesss[j]].AvgMs
		}
		return vmessstats[vmesss[i]].ErrCounter < vmessstats[vmesss[j]].ErrCounter
	})
	return vmesss
}

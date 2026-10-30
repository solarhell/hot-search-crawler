package source

import (
	"fmt"
	"strings"
)

// All 返回内置的全部数据源，顺序决定输出文件里的榜单顺序。
func All() []Source {
	return []Source{
		NewWeibo(),
		NewDouyin(),
		NewXiaohongshu(),
	}
}

// Parse 按逗号分隔的平台标识选择数据源。空字符串表示全部启用。
// 重复的平台只保留一次，未知标识返回错误。
func Parse(ids string) ([]Source, error) {
	if strings.TrimSpace(ids) == "" {
		return All(), nil
	}

	available := make(map[string]Source, len(All()))
	for _, src := range All() {
		available[src.ID()] = src
	}

	seen := make(map[string]bool)
	var sources []Source
	for _, raw := range strings.Split(ids, ",") {
		id := strings.TrimSpace(raw)
		if id == "" || seen[id] {
			continue
		}
		src, ok := available[id]
		if !ok {
			return nil, fmt.Errorf("未知平台 %q，可选值: %s", id, strings.Join(IDs(), ", "))
		}
		seen[id] = true
		sources = append(sources, src)
	}
	if len(sources) == 0 {
		return nil, fmt.Errorf("未选择任何平台")
	}
	return sources, nil
}

// IDs 返回内置平台标识，按展示顺序排列。
func IDs() []string {
	all := All()
	ids := make([]string, 0, len(all))
	for _, src := range all {
		ids = append(ids, src.ID())
	}
	return ids
}

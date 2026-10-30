package source

import (
	"strconv"
	"strings"
)

// ParseHotText 把平台热度文本换算成整数。
//
// 小红书返回 "947.5w"、"886w"，部分平台会返回 "1.2亿" 或带千分位的数字。
// 无法识别的文本返回 0，调用方仍可通过 HotText 保留原始值。
func ParseHotText(text string) int64 {
	s := strings.TrimSpace(text)
	if s == "" {
		return 0
	}
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, " ", "")

	multiplier := float64(1)
	switch {
	case strings.HasSuffix(s, "亿"):
		multiplier = 1e8
		s = strings.TrimSuffix(s, "亿")
	case strings.HasSuffix(s, "万"):
		multiplier = 1e4
		s = strings.TrimSuffix(s, "万")
	case strings.HasSuffix(s, "w"), strings.HasSuffix(s, "W"):
		multiplier = 1e4
		s = s[:len(s)-1]
	case strings.HasSuffix(s, "k"), strings.HasSuffix(s, "K"):
		multiplier = 1e3
		s = s[:len(s)-1]
	}

	s = strings.TrimSuffix(s, "+")
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(value * multiplier)
}

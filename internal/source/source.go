// Package source 定义热搜数据源接口，并实现微博、抖音、小红书三个平台。
package source

import (
	"context"
	"errors"
	"strconv"
	"time"
)

// Item 是统一后的热搜条目。不同平台能提供的字段不完全相同，
// 缺失的字段留空，平台特有信息放进 Extra。
type Item struct {
	// Rank 是榜单名次，由抓取层在过滤广告后重新编号，从 1 开始。
	Rank int `json:"rank"`
	// Title 是热搜原文。
	Title string `json:"title"`
	// URL 是可直接打开的搜索或话题链接。
	URL string `json:"url,omitempty"`
	// Hot 是归一化后的热度数值。小红书等平台的 "947.5w" 会换算成整数。
	Hot int64 `json:"hot,omitempty"`
	// HotText 是平台原始热度文本，例如 "947.5w"。
	HotText string `json:"hot_text,omitempty"`
	// Tag 是平台标注，例如 "新"、"热"、"爆"。
	Tag string `json:"tag,omitempty"`
	// Translated 是可选的中文译文，国内平台通常为空。
	Translated string `json:"translated,omitempty"`
	// Extra 保存平台特有字段，例如抖音的视频数、小红书的榜位变化。
	Extra map[string]any `json:"extra,omitempty"`
}

// Source 是一个热搜数据源。Fetch 返回按榜单顺序排列的条目，
// limit 是期望的最大条数；实现应跳过广告和空标题后再计数。
type Source interface {
	ID() string
	Name() string
	Translatable() bool
	Fetch(ctx context.Context, limit int) ([]Item, error)
}

// ErrThrottled 表示请求被平台限流。可以重试，但应立即退避而不是连续重试。
var ErrThrottled = errors.New("被平台限流")

// ThrottledError 携带限流平台和平台建议的等待时长。
type ThrottledError struct {
	Platform   string
	Code       int
	Message    string
	RetryAfter time.Duration
}

func (e *ThrottledError) Error() string {
	return e.Platform + "热搜被限流: code=" + strconv.Itoa(e.Code) + ", msg=" + e.Message
}

func (e *ThrottledError) Unwrap() error { return ErrThrottled }

// IsThrottled 判断错误是否属于限流。
func IsThrottled(err error) bool { return errors.Is(err, ErrThrottled) }

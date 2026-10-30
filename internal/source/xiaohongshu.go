package source

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	xiaohongshuEndpoint = "https://edith.xiaohongshu.com/api/sns/v1/search/hot_list"

	// xiaohongshuThrottledCode 是小红书“访问频繁，请稍后再试”的返回码。
	xiaohongshuThrottledCode = 300013

	// xiaohongshuThrottleRetry 被限流后的建议等待时长。
	xiaohongshuThrottleRetry = 10 * time.Minute
)

// xiaohongshuHeaders 是一份抓包得到的 iOS 客户端请求指纹。
//
// 小红书没有公开热榜接口，这里复用一个固定的设备指纹。同一指纹被间歇性
// 限流是常态而不是故障，代码按限流处理。指纹可能随客户端升级失效，届时
// 需要重新抓包更新，或通过 --sources 关闭该平台。
var xiaohongshuHeaders = map[string]string{
	"User-Agent":       "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.7(0x18000733) NetType/WIFI Language/zh_CN",
	"Referer":          "https://app.xhs.cn/",
	"xy-direction":     "22",
	"xy-platform-info": "platform=iOS&version=8.7&build=8070515&deviceId=C323D3A5-6A27-4CE6-AA0E-51C9D4C26A24&bundle=com.xingin.discover",
	"xy-common-params": "app_id=ECFAAF02&build=8070515&channel=AppStore&deviceId=C323D3A5-6A27-4CE6-AA0E-51C9D4C26A24&device_fingerprint=20230920120211bd7b71a80778509cf4211099ea911000010d2f20f6050264&device_fingerprint1=20230920120211bd7b71a80778509cf4211099ea911000010d2f20f6050264&device_model=phone&fid=1695182528-0-0-63b29d709954a1bb8c8733eb2fb58f29&gid=7dc4f3d168c355f1a886c54a898c6ef21fe7b9a847359afc77fc24ad&identifier_flag=0&lang=zh-Hans&launch_id=716882697&platform=iOS&project_id=ECFAAF&sid=session.1695189743787849952190&t=1695190591&teenager=0&tz=Asia/Shanghai&uis=light&version=8.7",
	"shield":           "XYAAAAAQAAAAEAAABTAAAAUzUWEe4xG1IYD9/c+qCLOlKGmTtFa+lG434Oe+FTRagxxoaz6rUWSZ3+juJYz8RZqct+oNMyZQxLEBaBEL+H3i0RhOBVGrauzVSARchIWFYwbwkV",
}

type xiaohongshuSource struct {
	endpoint string
}

// NewXiaohongshu 创建小红书数据源。
func NewXiaohongshu() Source {
	return &xiaohongshuSource{endpoint: xiaohongshuEndpoint}
}

func (s *xiaohongshuSource) ID() string         { return "xiaohongshu" }
func (s *xiaohongshuSource) Name() string       { return "小红书" }
func (s *xiaohongshuSource) Translatable() bool { return false }

func (s *xiaohongshuSource) Fetch(ctx context.Context, limit int) ([]Item, error) {
	body, err := getBody(ctx, s.endpoint, requestOptions{
		impersonate: false,
		headers:     xiaohongshuHeaders,
	})
	if err != nil {
		return nil, fmt.Errorf("请求小红书热搜失败: %w", err)
	}
	return parseXiaohongshu(body, limit)
}

type xiaohongshuResponse struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Msg     string `json:"msg"`
	Data    struct {
		Items []xiaohongshuItem `json:"items"`
	} `json:"data"`
}

type xiaohongshuItem struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Score      string `json:"score"`
	WordType   string `json:"word_type"`
	Type       string `json:"type"`
	RankChange int    `json:"rank_change"`
}

func (r *xiaohongshuResponse) errorText() string {
	if r.Msg != "" {
		return r.Msg
	}
	return r.Message
}

func parseXiaohongshu(body []byte, limit int) ([]Item, error) {
	var resp xiaohongshuResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析小红书热搜响应失败: %w", err)
	}

	if !resp.Success {
		message := resp.errorText()
		if resp.Code == xiaohongshuThrottledCode || strings.Contains(message, "访问频繁") {
			return nil, &ThrottledError{
				Platform:   "小红书",
				Code:       resp.Code,
				Message:    message,
				RetryAfter: xiaohongshuThrottleRetry,
			}
		}
		return nil, fmt.Errorf("小红书热搜返回异常: code=%d, msg=%s", resp.Code, message)
	}

	items := make([]Item, 0, limit)
	for _, raw := range resp.Data.Items {
		if len(items) >= limit {
			break
		}
		if raw.Title == "" {
			continue
		}

		tag := raw.WordType
		if tag == "无" {
			tag = ""
		}

		extra := make(map[string]any)
		if raw.ID != "" {
			extra["id"] = raw.ID
		}
		if raw.RankChange != 0 {
			extra["rank_change"] = raw.RankChange
		}

		items = append(items, Item{
			Rank:    len(items) + 1,
			Title:   raw.Title,
			URL:     "https://www.xiaohongshu.com/search_result?keyword=" + url.QueryEscape(raw.Title),
			Hot:     ParseHotText(raw.Score),
			HotText: raw.Score,
			Tag:     tag,
			Extra:   extra,
		})
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("小红书热搜返回数据为空")
	}
	return items, nil
}

package source

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

const douyinEndpoint = "https://aweme.snssdk.com/aweme/v1/hot/search/list/"

// douyinSource 抓取抖音热榜。
//
// word_list 是完整热榜，trending_list 是“实时上升热点”，两者语义不同，
// 这里只取 word_list，避免把两类数据混进同一个榜单。
type douyinSource struct {
	endpoint string
}

// NewDouyin 创建抖音数据源。
func NewDouyin() Source {
	return &douyinSource{endpoint: douyinEndpoint}
}

func (s *douyinSource) ID() string         { return "douyin" }
func (s *douyinSource) Name() string       { return "抖音" }
func (s *douyinSource) Translatable() bool { return false }

func (s *douyinSource) Fetch(ctx context.Context, limit int) ([]Item, error) {
	body, err := getBody(ctx, s.endpoint, requestOptions{
		impersonate: true,
		headers: map[string]string{
			"Referer": "https://www.douyin.com/",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("请求抖音热榜失败: %w", err)
	}
	return parseDouyin(body, limit)
}

type douyinResponse struct {
	StatusCode int `json:"status_code"`
	Data       struct {
		WordList []douyinItem `json:"word_list"`
	} `json:"data"`
}

type douyinItem struct {
	Word              string `json:"word"`
	HotValue          int64  `json:"hot_value"`
	Position          int    `json:"position"`
	Label             int    `json:"label"`
	EventTime         int64  `json:"event_time"`
	VideoCount        int    `json:"video_count"`
	DiscussVideoCount int    `json:"discuss_video_count"`
	SentenceID        string `json:"sentence_id"`
}

func parseDouyin(body []byte, limit int) ([]Item, error) {
	var resp douyinResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析抖音热榜响应失败: %w", err)
	}
	if resp.StatusCode != 0 {
		return nil, fmt.Errorf("抖音热榜返回异常状态: status_code=%d", resp.StatusCode)
	}

	items := make([]Item, 0, limit)
	for _, raw := range resp.Data.WordList {
		if len(items) >= limit {
			break
		}
		if raw.Word == "" {
			continue
		}

		extra := make(map[string]any)
		if raw.Label != 0 {
			extra["label"] = raw.Label
		}
		if raw.VideoCount > 0 {
			extra["video_count"] = raw.VideoCount
		}
		if raw.DiscussVideoCount > 0 {
			extra["discuss_video_count"] = raw.DiscussVideoCount
		}
		if raw.SentenceID != "" {
			extra["sentence_id"] = raw.SentenceID
		}
		if raw.EventTime > 0 {
			extra["event_time"] = raw.EventTime
		}

		items = append(items, Item{
			Rank:  len(items) + 1,
			Title: raw.Word,
			URL:   "https://www.douyin.com/search/" + url.PathEscape(raw.Word),
			Hot:   raw.HotValue,
			Extra: extra,
		})
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("抖音热榜返回数据为空")
	}
	return items, nil
}

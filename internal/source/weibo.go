package source

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

const weiboEndpoint = "https://weibo.com/ajax/side/hotSearch"

// weiboSource 抓取微博热搜榜。
//
// 接口返回的 realtime 列表里混有广告（is_ad=1，且 realpos 为空），
// 必须过滤掉，否则广告会占据榜单名次。
type weiboSource struct {
	endpoint string
}

// NewWeibo 创建微博数据源。
func NewWeibo() Source {
	return &weiboSource{endpoint: weiboEndpoint}
}

func (s *weiboSource) ID() string         { return "weibo" }
func (s *weiboSource) Name() string       { return "微博" }
func (s *weiboSource) Translatable() bool { return false }

func (s *weiboSource) Fetch(ctx context.Context, limit int) ([]Item, error) {
	body, err := getBody(ctx, s.endpoint, requestOptions{
		impersonate: true,
		headers: map[string]string{
			"Referer": "https://weibo.com/",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("请求微博热搜失败: %w", err)
	}
	return parseWeibo(body, limit)
}

type weiboResponse struct {
	OK   int `json:"ok"`
	Data struct {
		Realtime []weiboItem `json:"realtime"`
	} `json:"data"`
}

type weiboItem struct {
	Word       string `json:"word"`
	WordScheme string `json:"word_scheme"`
	Num        int64  `json:"num"`
	LabelName  string `json:"label_name"`
	FlagDesc   string `json:"flag_desc"`
	IsAd       int    `json:"is_ad"`
}

func parseWeibo(body []byte, limit int) ([]Item, error) {
	var resp weiboResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析微博热搜响应失败: %w", err)
	}
	if resp.OK != 1 {
		return nil, fmt.Errorf("微博热搜返回异常状态: ok=%d", resp.OK)
	}

	items := make([]Item, 0, limit)
	for _, raw := range resp.Data.Realtime {
		if len(items) >= limit {
			break
		}
		if raw.IsAd == 1 || raw.Word == "" {
			continue
		}

		scheme := raw.WordScheme
		if scheme == "" {
			scheme = raw.Word
		}

		tag := raw.LabelName
		if tag == "" {
			tag = raw.FlagDesc
		}

		items = append(items, Item{
			Rank:  len(items) + 1,
			Title: raw.Word,
			URL:   "https://s.weibo.com/weibo?q=" + url.QueryEscape(scheme),
			Hot:   raw.Num,
			Tag:   tag,
		})
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("微博热搜返回数据为空")
	}
	return items, nil
}

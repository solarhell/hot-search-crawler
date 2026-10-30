package source

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatalf("读取测试数据 %s 失败: %v", name, err)
	}
	return data
}

func TestParseHotText(t *testing.T) {
	cases := []struct {
		input string
		want  int64
	}{
		{"947.5w", 9_475_000},
		{"886w", 8_860_000},
		{"934.8W", 9_348_000},
		{"1.2亿", 120_000_000},
		{"3.5万", 35_000},
		{"12,345", 12_345},
		{"2.1k", 2_100},
		{"123", 123},
		{"", 0},
		{"暂无", 0},
	}
	for _, tc := range cases {
		if got := ParseHotText(tc.input); got != tc.want {
			t.Errorf("ParseHotText(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

func TestParseWeibo(t *testing.T) {
	items, err := parseWeibo(readFixture(t, "weibo.json"), 10)
	if err != nil {
		t.Fatalf("parseWeibo: %v", err)
	}

	// 广告（is_ad=1）必须被过滤。
	if len(items) != 2 {
		t.Fatalf("items = %d, want 2", len(items))
	}
	if items[0].Rank != 1 || items[0].Title != "TES战胜BLG" || items[0].Hot != 2066277 {
		t.Errorf("items[0] = %+v", items[0])
	}
	if items[0].Tag != "新" {
		t.Errorf("items[0].Tag = %q, want 新", items[0].Tag)
	}
	if !strings.Contains(items[0].URL, "weibo.com") || !strings.Contains(items[0].URL, "TES") {
		t.Errorf("items[0].URL = %q", items[0].URL)
	}
	if items[1].Rank != 2 {
		t.Errorf("items[1].Rank = %d, want 2", items[1].Rank)
	}
}

func TestParseWeiboLimitAndStatus(t *testing.T) {
	items, err := parseWeibo(readFixture(t, "weibo.json"), 1)
	if err != nil {
		t.Fatalf("parseWeibo: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("items = %d, want 1", len(items))
	}

	if _, err := parseWeibo([]byte(`{"ok":0}`), 10); err == nil {
		t.Fatal("ok=0 应返回错误")
	}
}

func TestParseDouyin(t *testing.T) {
	items, err := parseDouyin(readFixture(t, "douyin.json"), 10)
	if err != nil {
		t.Fatalf("parseDouyin: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("items = %d, want 2", len(items))
	}
	if items[0].Rank != 1 || items[0].Hot != 12159378 {
		t.Errorf("items[0] = %+v", items[0])
	}
	if items[0].Extra["label"] != 3 {
		t.Errorf("items[0].Extra[label] = %v, want 3", items[0].Extra["label"])
	}
	if items[0].Extra["video_count"] != 7 {
		t.Errorf("items[0].Extra[video_count] = %v, want 7", items[0].Extra["video_count"])
	}
	if items[1].Rank != 2 || items[1].Title != "中国U17女足首进世界杯8强" {
		t.Errorf("items[1] = %+v", items[1])
	}
}

func TestParseXiaohongshu(t *testing.T) {
	items, err := parseXiaohongshu(readFixture(t, "xiaohongshu.json"), 10)
	if err != nil {
		t.Fatalf("parseXiaohongshu: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("items = %d, want 2", len(items))
	}
	if items[0].Rank != 1 || items[0].Hot != 9_475_000 || items[0].HotText != "947.5w" {
		t.Errorf("items[0] = %+v", items[0])
	}
	if items[0].Tag != "热" {
		t.Errorf("items[0].Tag = %q, want 热", items[0].Tag)
	}
	if items[0].Extra["id"] != "dora_2178732" {
		t.Errorf("items[0].Extra[id] = %v", items[0].Extra["id"])
	}
	if items[1].Tag != "" {
		t.Errorf("items[1].Tag = %q, want empty", items[1].Tag)
	}
	if items[1].Hot != 9_348_000 {
		t.Errorf("items[1].Hot = %d", items[1].Hot)
	}
}

func TestParseXiaohongshuThrottled(t *testing.T) {
	_, err := parseXiaohongshu(readFixture(t, "xiaohongshu_throttled.json"), 10)
	if err == nil {
		t.Fatal("限流响应应返回错误")
	}
	if !errors.Is(err, ErrThrottled) {
		t.Fatalf("err = %v, want ErrThrottled", err)
	}

	var throttled *ThrottledError
	if !errors.As(err, &throttled) {
		t.Fatalf("err = %T, want *ThrottledError", err)
	}
	if throttled.RetryAfter != 10*time.Minute {
		t.Errorf("RetryAfter = %v, want 10m", throttled.RetryAfter)
	}
}

func TestParseRegistry(t *testing.T) {
	all, err := Parse("")
	if err != nil {
		t.Fatalf("Parse(empty): %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("all = %d, want 3", len(all))
	}

	selected, err := Parse("weibo, xiaohongshu,weibo")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(selected) != 2 || selected[0].ID() != "weibo" || selected[1].ID() != "xiaohongshu" {
		t.Fatalf("selected = %+v", selected)
	}

	if _, err := Parse("weibo,unknown"); err == nil {
		t.Fatal("未知平台应返回错误")
	}
}

package crawl

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/solarhell/hot-search-crawler/internal/source"
)

type fakeSource struct {
	id           string
	name         string
	translatable bool
	items        []source.Item
	err          error
	failures     int
	calls        int
}

func (f *fakeSource) ID() string         { return f.id }
func (f *fakeSource) Name() string       { return f.name }
func (f *fakeSource) Translatable() bool { return f.translatable }
func (f *fakeSource) Fetch(context.Context, int) ([]source.Item, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	if f.calls <= f.failures {
		return nil, errors.New("temporary failure")
	}
	return f.items, nil
}

type fakeTranslator struct {
	results []string
	err     error
	calls   int
}

func (f *fakeTranslator) Translate(context.Context, []string) ([]string, error) {
	f.calls++
	return f.results, f.err
}

func quietLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func testOptions(now time.Time) Options {
	return Options{
		Top:      10,
		Timeout:  time.Second,
		Location: time.FixedZone("Asia/Shanghai", 8*3600),
		Logger:   quietLogger(),
		Now:      func() time.Time { return now },
	}
}

func TestRunKeepsSourceOrderAndRanks(t *testing.T) {
	now := time.Date(2026, 10, 30, 9, 5, 0, 0, time.UTC)
	sources := []source.Source{
		&fakeSource{id: "weibo", name: "微博", items: []source.Item{
			{Title: "a", Hot: 100},
			{Title: "b", Hot: 90},
		}},
		&fakeSource{id: "douyin", name: "抖音", items: []source.Item{
			{Title: "c", Hot: 80},
		}},
	}

	snapshot := Run(context.Background(), sources, testOptions(now))

	if len(snapshot.Boards) != 2 {
		t.Fatalf("boards = %d, want 2", len(snapshot.Boards))
	}
	if snapshot.Boards[0].Platform != "weibo" || snapshot.Boards[1].Platform != "douyin" {
		t.Fatalf("board order = %s, %s", snapshot.Boards[0].Platform, snapshot.Boards[1].Platform)
	}
	if !snapshot.Boards[0].OK || len(snapshot.Boards[0].Items) != 2 {
		t.Fatalf("weibo board = %+v", snapshot.Boards[0])
	}
	if snapshot.Boards[0].Items[1].Rank != 2 {
		t.Errorf("rank = %d, want 2", snapshot.Boards[0].Items[1].Rank)
	}
	if got := snapshot.GeneratedAt.In(time.FixedZone("Asia/Shanghai", 8*3600)).Hour(); got != 17 {
		t.Errorf("generated hour = %d, want 17 (UTC+8)", got)
	}
}

func TestRunTopLimit(t *testing.T) {
	now := time.Date(2026, 10, 30, 9, 5, 0, 0, time.UTC)
	src := &fakeSource{id: "weibo", name: "微博", items: []source.Item{
		{Title: "a"}, {Title: "b"}, {Title: "c"},
	}}
	opts := testOptions(now)
	opts.Top = 2

	snapshot := Run(context.Background(), []source.Source{src}, opts)
	if len(snapshot.Boards[0].Items) != 2 {
		t.Fatalf("items = %d, want 2", len(snapshot.Boards[0].Items))
	}
}

func TestRunPartialFailure(t *testing.T) {
	now := time.Date(2026, 10, 30, 9, 5, 0, 0, time.UTC)
	sources := []source.Source{
		&fakeSource{id: "weibo", name: "微博", items: []source.Item{{Title: "a"}}},
		&fakeSource{id: "douyin", name: "抖音", err: errors.New("boom")},
	}

	snapshot := Run(context.Background(), sources, testOptions(now))

	if !snapshot.Boards[0].OK {
		t.Error("weibo 应成功")
	}
	if snapshot.Boards[1].OK || snapshot.Boards[1].Error != "boom" {
		t.Errorf("douyin board = %+v", snapshot.Boards[1])
	}
	if snapshot.Boards[1].Items == nil {
		t.Error("失败榜单的 Items 应为空数组而不是 nil")
	}
}

func TestRunEmptyItemsIsFailure(t *testing.T) {
	now := time.Date(2026, 10, 30, 9, 5, 0, 0, time.UTC)
	src := &fakeSource{id: "weibo", name: "微博", items: nil}

	snapshot := Run(context.Background(), []source.Source{src}, testOptions(now))
	if snapshot.Boards[0].OK {
		t.Fatal("空结果应视为失败")
	}
}

func TestRunRetriesTemporaryFailure(t *testing.T) {
	oldBase := retryBase
	retryBase = time.Millisecond
	defer func() { retryBase = oldBase }()

	now := time.Date(2026, 10, 30, 9, 5, 0, 0, time.UTC)
	src := &fakeSource{
		id:       "weibo",
		name:     "微博",
		items:    []source.Item{{Title: "a"}},
		failures: 1,
	}
	opts := testOptions(now)
	opts.Retries = 1

	snapshot := Run(context.Background(), []source.Source{src}, opts)
	if !snapshot.Boards[0].OK {
		t.Fatalf("重试后应成功: %+v", snapshot.Boards[0])
	}
	if src.calls != 2 {
		t.Errorf("calls = %d, want 2", src.calls)
	}
}

func TestRunDoesNotRetryThrottle(t *testing.T) {
	now := time.Date(2026, 10, 30, 9, 5, 0, 0, time.UTC)
	src := &fakeSource{
		id:   "xiaohongshu",
		name: "小红书",
		err: &source.ThrottledError{
			Platform:   "小红书",
			Code:       300013,
			Message:    "访问频繁",
			RetryAfter: 10 * time.Minute,
		},
	}
	opts := testOptions(now)
	opts.Retries = 3

	snapshot := Run(context.Background(), []source.Source{src}, opts)
	if snapshot.Boards[0].OK {
		t.Fatal("限流应失败")
	}
	if src.calls != 1 {
		t.Errorf("calls = %d, want 1（限流不应重试）", src.calls)
	}
}

func TestRunTranslatesTranslatableSource(t *testing.T) {
	now := time.Date(2026, 10, 30, 9, 5, 0, 0, time.UTC)
	src := &fakeSource{
		id:           "twitter",
		name:         "Twitter",
		translatable: true,
		items:        []source.Item{{Title: "hello"}},
	}
	tr := &fakeTranslator{results: []string{"你好"}}
	opts := testOptions(now)
	opts.Translator = tr

	snapshot := Run(context.Background(), []source.Source{src}, opts)
	if tr.calls != 1 {
		t.Fatalf("translator calls = %d, want 1", tr.calls)
	}
	if got := snapshot.Boards[0].Items[0].Translated; got != "你好" {
		t.Errorf("translated = %q, want 你好", got)
	}
}

func TestRunKeepsOriginalWhenTranslationFails(t *testing.T) {
	now := time.Date(2026, 10, 30, 9, 5, 0, 0, time.UTC)
	src := &fakeSource{
		id:           "twitter",
		name:         "Twitter",
		translatable: true,
		items:        []source.Item{{Title: "hello"}},
	}
	tr := &fakeTranslator{err: errors.New("model down")}
	opts := testOptions(now)
	opts.Translator = tr

	snapshot := Run(context.Background(), []source.Source{src}, opts)
	if !snapshot.Boards[0].OK {
		t.Fatal("翻译失败不应导致抓取失败")
	}
	if snapshot.Boards[0].Items[0].Translated != "" {
		t.Error("翻译失败时应保留空译文")
	}
}

func TestRetryWait(t *testing.T) {
	oldBase := retryBase
	retryBase = 2 * time.Second
	defer func() { retryBase = oldBase }()

	cases := map[int]time.Duration{
		0: 2 * time.Second,
		1: 4 * time.Second,
		2: 8 * time.Second,
		3: 10 * time.Second,
		9: 10 * time.Second,
	}
	for attempt, want := range cases {
		if got := retryWait(attempt); got != want {
			t.Errorf("retryWait(%d) = %v, want %v", attempt, got, want)
		}
	}
}

package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/solarhell/hot-search-crawler/internal/crawl"
	"github.com/solarhell/hot-search-crawler/internal/source"
)

var shanghai = time.FixedZone("Asia/Shanghai", 8*3600)

func snapshotAt(t *testing.T, store *Store, hour int, boards ...crawl.Board) crawl.Snapshot {
	t.Helper()
	now := time.Date(2026, 10, 30, hour, 5, 0, 0, shanghai)
	return crawl.Snapshot{
		SchemaVersion: crawl.SchemaVersion,
		GeneratedAt:   now,
		Timezone:      shanghai.String(),
		Boards:        boards,
	}
}

func board(platform, name string, items ...source.Item) crawl.Board {
	return crawl.Board{
		Platform: platform,
		Name:     name,
		OK:       true,
		Items:    items,
	}
}

func newTestStore(t *testing.T) *Store {
	t.Helper()
	root := t.TempDir()
	return &Store{
		DataDir:    filepath.Join(root, "data"),
		ArchiveDir: filepath.Join(root, "archives"),
		Location:   shanghai,
		Now:        func() time.Time { return time.Date(2026, 10, 30, 12, 0, 0, 0, shanghai) },
	}
}

func TestWriteSnapshotAndLatest(t *testing.T) {
	store := newTestStore(t)
	snap := snapshotAt(t, store, 9, board("weibo", "微博", source.Item{Title: "a", Rank: 1}))

	path, err := store.WriteSnapshot(snap)
	if err != nil {
		t.Fatalf("WriteSnapshot: %v", err)
	}
	want := filepath.Join(store.DataDir, "2026-10-30", "09.json")
	if path != want {
		t.Errorf("path = %q, want %q", path, want)
	}
	if _, err := os.Stat(want); err != nil {
		t.Errorf("stat: %v", err)
	}

	latest, err := store.WriteLatest(snap)
	if err != nil {
		t.Fatalf("WriteLatest: %v", err)
	}
	if latest != filepath.Join(store.DataDir, "latest.json") {
		t.Errorf("latest = %q", latest)
	}
}

func TestBuildSummaryMergesHours(t *testing.T) {
	store := newTestStore(t)

	morning := snapshotAt(t, store, 9,
		board("weibo", "微博",
			source.Item{Title: "A", Rank: 1, Hot: 100},
			source.Item{Title: "B", Rank: 2, Hot: 80},
		),
		board("douyin", "抖音",
			source.Item{Title: "D", Rank: 1, Hot: 50},
		),
	)
	noon := snapshotAt(t, store, 12,
		board("weibo", "微博",
			source.Item{Title: "A", Rank: 2, Hot: 120},
			source.Item{Title: "C", Rank: 1, Hot: 90},
		),
	)

	for _, snap := range []crawl.Snapshot{morning, noon} {
		if _, err := store.WriteSnapshot(snap); err != nil {
			t.Fatalf("WriteSnapshot: %v", err)
		}
	}

	summary, err := store.BuildSummary("2026-10-30")
	if err != nil {
		t.Fatalf("BuildSummary: %v", err)
	}

	if len(summary.Boards) != 2 {
		t.Fatalf("boards = %d, want 2", len(summary.Boards))
	}
	weibo := summary.Boards[0]
	if weibo.Platform != "weibo" {
		t.Fatalf("board[0] = %+v", weibo)
	}
	if len(weibo.Items) != 3 {
		t.Fatalf("weibo items = %d, want 3", len(weibo.Items))
	}

	// A 出现两次，取最高热度、最好名次。
	a := weibo.Items[0]
	if a.Title != "A" || a.Hot != 120 || a.BestRank != 1 || a.Appearances != 2 {
		t.Errorf("A = %+v", a)
	}
	// 其余按热度排序。
	if weibo.Items[1].Title != "C" || weibo.Items[2].Title != "B" {
		t.Errorf("order = %s, %s", weibo.Items[1].Title, weibo.Items[2].Title)
	}
	if summary.UpdatedAt.Location() != shanghai {
		t.Errorf("UpdatedAt location = %v", summary.UpdatedAt.Location())
	}
}

func TestBuildSummaryMissingDate(t *testing.T) {
	store := newTestStore(t)
	summary, err := store.BuildSummary("2026-10-31")
	if err != nil {
		t.Fatalf("BuildSummary: %v", err)
	}
	if len(summary.Boards) != 0 {
		t.Errorf("boards = %d, want 0", len(summary.Boards))
	}
}

func TestBuildSummarySkipsFailedBoards(t *testing.T) {
	store := newTestStore(t)
	snap := snapshotAt(t, store, 9,
		crawl.Board{Platform: "weibo", Name: "微博", OK: false, Error: "boom", Items: []source.Item{}},
		board("douyin", "抖音", source.Item{Title: "D", Rank: 1, Hot: 10}),
	)
	if _, err := store.WriteSnapshot(snap); err != nil {
		t.Fatalf("WriteSnapshot: %v", err)
	}

	summary, err := store.BuildSummary("2026-10-30")
	if err != nil {
		t.Fatalf("BuildSummary: %v", err)
	}
	if len(summary.Boards) != 1 || summary.Boards[0].Platform != "douyin" {
		t.Fatalf("boards = %+v", summary.Boards)
	}
}

func TestWriteSummaryAndArchive(t *testing.T) {
	store := newTestStore(t)
	summary := Summary{
		SchemaVersion: SchemaVersion,
		Date:          "2026-10-30",
		Timezone:      shanghai.String(),
		UpdatedAt:     time.Date(2026, 10, 30, 12, 0, 0, 0, shanghai),
		Boards: []SummaryBoard{{
			Platform: "weibo",
			Name:     "微博",
			Items: []SummaryItem{{
				Title:       "A|B",
				URL:         "https://example.com/search?q=a%20b",
				Hot:         12_345,
				BestRank:    1,
				Appearances: 2,
			}},
		}},
	}

	path, err := store.WriteSummary(summary)
	if err != nil {
		t.Fatalf("WriteSummary: %v", err)
	}
	if path != filepath.Join(store.DataDir, "2026-10-30", "summary.json") {
		t.Errorf("path = %q", path)
	}

	archivePath, err := store.WriteArchive(summary)
	if err != nil {
		t.Fatalf("WriteArchive: %v", err)
	}
	data, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	text := string(data)
	if !strings.Contains(text, "# 2026-10-30 热搜归档") {
		t.Errorf("archive = %s", text)
	}
	if !strings.Contains(text, `A\|B`) {
		t.Errorf("标题中的竖线未转义: %s", text)
	}
	if !strings.Contains(text, "1.2万") {
		t.Errorf("热度格式 = %s", text)
	}
}

func TestUpdateReadmeReplacesSection(t *testing.T) {
	store := newTestStore(t)
	path := filepath.Join(t.TempDir(), "README.md")
	initial := "# Demo\n\nintro\n\n" + beginMarker + "\nold\n" + endMarker + "\n\nfooter\n"
	if err := os.WriteFile(path, []byte(initial), 0o644); err != nil {
		t.Fatal(err)
	}

	summary := Summary{
		Date:      "2026-10-30",
		Timezone:  shanghai.String(),
		UpdatedAt: time.Date(2026, 10, 30, 12, 0, 0, 0, shanghai),
		Boards: []SummaryBoard{{
			Platform: "weibo",
			Name:     "微博",
			Items: []SummaryItem{{
				Title: "A",
				URL:   "https://example.com/a",
				Hot:   100_000,
			}},
		}},
	}

	if err := store.UpdateReadme(path, summary); err != nil {
		t.Fatalf("UpdateReadme: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)

	if !strings.Contains(text, "# Demo") || !strings.Contains(text, "footer") {
		t.Error("标记区间外的内容被破坏")
	}
	if strings.Contains(text, "\nold\n") {
		t.Error("旧热榜内容未被替换")
	}
	if !strings.Contains(text, "### 微博") || !strings.Contains(text, "10万") {
		t.Errorf("新内容缺失: %s", text)
	}
	if strings.Count(text, beginMarker) != 1 || strings.Count(text, endMarker) != 1 {
		t.Error("标记数量异常")
	}
}

func TestUpdateReadmeAppendsWhenMissing(t *testing.T) {
	store := newTestStore(t)
	path := filepath.Join(t.TempDir(), "README.md")
	if err := os.WriteFile(path, []byte("# Demo"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := store.UpdateReadme(path, Summary{Date: "2026-10-30", Timezone: shanghai.String()}); err != nil {
		t.Fatalf("UpdateReadme: %v", err)
	}
	data, _ := os.ReadFile(path)
	text := string(data)
	if !strings.Contains(text, beginMarker) || !strings.Contains(text, endMarker) {
		t.Errorf("未追加标记: %s", text)
	}
}

func TestFormatHot(t *testing.T) {
	cases := map[int64]string{
		0:           "0",
		999:         "999",
		10_000:      "1万",
		12_345:      "1.2万",
		9_475_000:   "947.5万",
		120_000_000: "1.2亿",
	}
	for value, want := range cases {
		if got := formatHot(value); got != want {
			t.Errorf("formatHot(%d) = %q, want %q", value, got, want)
		}
	}
}

func TestIsHourFile(t *testing.T) {
	valid := []string{"00.json", "09.json", "23.json"}
	for _, name := range valid {
		if !isHourFile(name) {
			t.Errorf("%q 应为小时文件", name)
		}
	}
	invalid := []string{"summary.json", "9.json", "24.txt", "latest.json", "a9.json"}
	for _, name := range invalid {
		if isHourFile(name) {
			t.Errorf("%q 不应为小时文件", name)
		}
	}
}

func TestSnapshotJSONRoundTrip(t *testing.T) {
	store := newTestStore(t)
	snap := snapshotAt(t, store, 9, board("weibo", "微博", source.Item{Title: "A", Rank: 1}))
	path, err := store.WriteSnapshot(snap)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var decoded crawl.Snapshot
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("JSON 无法解析: %v", err)
	}
	if decoded.Boards[0].Items[0].Title != "A" {
		t.Errorf("decoded = %+v", decoded)
	}
}

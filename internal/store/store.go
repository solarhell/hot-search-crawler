// Package store 负责把抓取快照落盘，并生成日报和 README 热榜。
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/solarhell/hot-search-crawler/internal/crawl"
)

// SchemaVersion 是 summary.json 的结构版本。
const SchemaVersion = 1

// Summary 是按天聚合后的热搜数据。
type Summary struct {
	SchemaVersion int            `json:"schema_version"`
	Date          string         `json:"date"`
	Timezone      string         `json:"timezone"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Boards        []SummaryBoard `json:"boards"`
}

// SummaryBoard 是单个平台当天的聚合榜单。
type SummaryBoard struct {
	Platform string        `json:"platform"`
	Name     string        `json:"name"`
	Items    []SummaryItem `json:"items"`
}

// SummaryItem 是当天出现过的热搜词及其在榜统计。
type SummaryItem struct {
	Title       string    `json:"title"`
	URL         string    `json:"url,omitempty"`
	Hot         int64     `json:"hot,omitempty"`
	HotText     string    `json:"hot_text,omitempty"`
	Tag         string    `json:"tag,omitempty"`
	Translated  string    `json:"translated,omitempty"`
	BestRank    int       `json:"best_rank"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
	Appearances int       `json:"appearances"`
}

// Store 管理 data/ 与 archives/ 目录。
type Store struct {
	DataDir    string
	ArchiveDir string
	Location   *time.Location
	// Now 用于测试注入，为空时使用 time.Now。
	Now func() time.Time
}

// New 创建 Store。
func New(dataDir, archiveDir string, location *time.Location) *Store {
	if location == nil {
		location = time.Local
	}
	return &Store{
		DataDir:    dataDir,
		ArchiveDir: archiveDir,
		Location:   location,
		Now:        time.Now,
	}
}

// Date 返回快照在目标时区下的日期。
func (s *Store) Date(snapshot crawl.Snapshot) string {
	return snapshot.GeneratedAt.In(s.Location).Format(time.DateOnly)
}

// WriteSnapshot 写入 data/<date>/<HH>.json，返回文件路径。
func (s *Store) WriteSnapshot(snapshot crawl.Snapshot) (string, error) {
	local := snapshot.GeneratedAt.In(s.Location)
	path := filepath.Join(s.DataDir, local.Format(time.DateOnly), local.Format("15")+".json")
	if err := writeJSON(path, snapshot); err != nil {
		return "", err
	}
	return path, nil
}

// WriteLatest 写入 data/latest.json，供静态页面直接消费。
func (s *Store) WriteLatest(snapshot crawl.Snapshot) (string, error) {
	path := filepath.Join(s.DataDir, "latest.json")
	if err := writeJSON(path, snapshot); err != nil {
		return "", err
	}
	return path, nil
}

// WriteSummary 写入 data/<date>/summary.json。
func (s *Store) WriteSummary(summary Summary) (string, error) {
	path := filepath.Join(s.DataDir, summary.Date, "summary.json")
	if err := writeJSON(path, summary); err != nil {
		return "", err
	}
	return path, nil
}

// BuildSummary 读取某一天的所有小时快照并重新聚合。
//
// 每次从原始小时文件重建而不是增量合并，保证重复运行、补抓或删除
// 单个快照后，summary 仍然是当天数据的准确投影。
func (s *Store) BuildSummary(date string) (Summary, error) {
	summary := Summary{
		SchemaVersion: SchemaVersion,
		Date:          date,
		Timezone:      s.Location.String(),
		UpdatedAt:     s.Now().In(s.Location),
		Boards:        []SummaryBoard{},
	}

	dir := filepath.Join(s.DataDir, date)
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return summary, nil
	}
	if err != nil {
		return Summary{}, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if isHourFile(entry.Name()) {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	type accumulator struct {
		platform string
		name     string
		items    map[string]*SummaryItem
	}
	boards := make(map[string]*accumulator)
	var boardOrder []*accumulator

	for _, name := range files {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return Summary{}, err
		}

		var snapshot crawl.Snapshot
		if err := json.Unmarshal(data, &snapshot); err != nil {
			return Summary{}, fmt.Errorf("解析快照 %s/%s 失败: %w", date, name, err)
		}

		seenAt := snapshot.GeneratedAt.In(s.Location)
		for _, board := range snapshot.Boards {
			if !board.OK {
				continue
			}
			acc, ok := boards[board.Platform]
			if !ok {
				acc = &accumulator{
					platform: board.Platform,
					name:     board.Name,
					items:    make(map[string]*SummaryItem),
				}
				boards[board.Platform] = acc
				boardOrder = append(boardOrder, acc)
			}

			seen := make(map[string]bool)
			for _, item := range board.Items {
				if item.Title == "" || seen[item.Title] {
					continue
				}
				seen[item.Title] = true

				existing, ok := acc.items[item.Title]
				if !ok {
					acc.items[item.Title] = &SummaryItem{
						Title:       item.Title,
						URL:         item.URL,
						Hot:         item.Hot,
						HotText:     item.HotText,
						Tag:         item.Tag,
						Translated:  item.Translated,
						BestRank:    item.Rank,
						FirstSeen:   seenAt,
						LastSeen:    seenAt,
						Appearances: 1,
					}
					continue
				}

				if item.Hot > existing.Hot {
					existing.Hot = item.Hot
				}
				if item.HotText != "" {
					existing.HotText = item.HotText
				}
				if existing.URL == "" {
					existing.URL = item.URL
				}
				if existing.Tag == "" {
					existing.Tag = item.Tag
				}
				if item.Translated != "" {
					existing.Translated = item.Translated
				}
				if item.Rank > 0 && (existing.BestRank == 0 || item.Rank < existing.BestRank) {
					existing.BestRank = item.Rank
				}
				if seenAt.Before(existing.FirstSeen) {
					existing.FirstSeen = seenAt
				}
				if seenAt.After(existing.LastSeen) {
					existing.LastSeen = seenAt
				}
				existing.Appearances++
			}
		}
	}

	for _, acc := range boardOrder {
		items := make([]SummaryItem, 0, len(acc.items))
		for _, item := range acc.items {
			items = append(items, *item)
		}
		sortItems(items)
		summary.Boards = append(summary.Boards, SummaryBoard{
			Platform: acc.platform,
			Name:     acc.name,
			Items:    items,
		})
	}

	return summary, nil
}

// WriteArchive 写入 archives/<date>.md。
func (s *Store) WriteArchive(summary Summary) (string, error) {
	path := filepath.Join(s.ArchiveDir, summary.Date+".md")
	if err := writeFileAtomic(path, []byte(renderArchive(summary)), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// HasData 判断聚合结果里是否有任何榜单条目。
func HasData(summary Summary) bool {
	for _, board := range summary.Boards {
		if len(board.Items) > 0 {
			return true
		}
	}
	return false
}

func sortItems(items []SummaryItem) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Hot != items[j].Hot {
			return items[i].Hot > items[j].Hot
		}
		if items[i].BestRank != items[j].BestRank {
			return items[i].BestRank < items[j].BestRank
		}
		return items[i].Title < items[j].Title
	})
}

func isHourFile(name string) bool {
	if len(name) != len("15.json") || name[2] != '.' || name[3:] != "json" {
		return false
	}
	return name[0] >= '0' && name[0] <= '9' && name[1] >= '0' && name[1] <= '9'
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return writeFileAtomic(path, data, 0o644)
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Chmod(perm); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

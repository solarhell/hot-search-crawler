package store

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	beginMarker = "<!-- BEGIN HOTSEARCH -->"
	endMarker   = "<!-- END HOTSEARCH -->"

	readmeTopN  = 10
	archiveTopN = 100
)

// UpdateReadme 把今日热榜写进 README 的标记区间。
// 标记不存在时追加到文件末尾，方便使用者把代码复制到已有 README。
func (s *Store) UpdateReadme(path string, summary Summary) error {
	content, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	text := string(content)
	section := renderReadmeSection(summary)

	begin := strings.Index(text, beginMarker)
	end := strings.Index(text, endMarker)
	if begin >= 0 && end > begin {
		text = text[:begin] + beginMarker + "\n" + section + "\n" + endMarker + text[end+len(endMarker):]
	} else {
		if text != "" && !strings.HasSuffix(text, "\n") {
			text += "\n"
		}
		text += "\n" + beginMarker + "\n" + section + "\n" + endMarker + "\n"
	}

	return writeFileAtomic(path, []byte(text), 0o644)
}

func renderReadmeSection(summary Summary) string {
	var b strings.Builder

	fmt.Fprintf(&b, "> 数据更新于 **%s**（%s），展示今日榜单前 %d。\n\n",
		summary.UpdatedAt.Format("2006-01-02 15:04"), summary.Timezone, readmeTopN)

	for _, board := range summary.Boards {
		fmt.Fprintf(&b, "### %s\n\n", board.Name)
		if len(board.Items) == 0 {
			b.WriteString("暂无数据。\n\n")
			continue
		}

		n := min(len(board.Items), readmeTopN)
		for i := 0; i < n; i++ {
			item := board.Items[i]
			title := escapeMarkdown(item.Title)
			if item.URL != "" {
				title = "[" + title + "](" + escapeURL(item.URL) + ")"
			}

			line := fmt.Sprintf("%d. %s", i+1, title)
			if hot := hotLabel(item); hot != "" {
				line += " — " + hot
			}
			if item.Tag != "" {
				line += " · " + escapeMarkdown(item.Tag)
			}
			b.WriteString(line + "\n")
		}
		b.WriteString("\n")
	}

	fmt.Fprintf(&b, "[完整快照](data/latest.json) · [今日归档](archives/%s.md)\n", summary.Date)
	return strings.TrimRight(b.String(), "\n")
}

func renderArchive(summary Summary) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# %s 热搜归档\n\n", summary.Date)
	fmt.Fprintf(&b, "> 时区：%s · 更新：%s\n\n", summary.Timezone,
		summary.UpdatedAt.Format("2006-01-02 15:04"))

	for _, board := range summary.Boards {
		fmt.Fprintf(&b, "## %s\n\n", board.Name)
		if len(board.Items) == 0 {
			b.WriteString("暂无数据。\n\n")
			continue
		}

		b.WriteString("| # | 热搜 | 热度 | 最高排名 | 在榜次数 |\n")
		b.WriteString("|---:|---|---:|---:|---:|\n")

		n := min(len(board.Items), archiveTopN)
		for i := 0; i < n; i++ {
			item := board.Items[i]
			title := escapeTable(item.Title)
			if item.URL != "" {
				title = "[" + title + "](" + escapeURL(item.URL) + ")"
			}

			bestRank := "-"
			if item.BestRank > 0 {
				bestRank = strconv.Itoa(item.BestRank)
			}
			fmt.Fprintf(&b, "| %d | %s | %s | %s | %d |\n",
				i+1, title, hotLabel(item), bestRank, item.Appearances)
		}

		if len(board.Items) > n {
			fmt.Fprintf(&b, "\n> 仅展示热度前 %d 条，完整数据见 `data/%s/summary.json`。\n",
				archiveTopN, summary.Date)
		}
		b.WriteString("\n")
	}

	return strings.TrimRight(b.String(), "\n") + "\n"
}

func hotLabel(item SummaryItem) string {
	if item.Hot > 0 {
		return formatHot(item.Hot)
	}
	return item.HotText
}

func formatHot(value int64) string {
	switch {
	case value >= 100_000_000:
		return trimTrailingZero(strconv.FormatFloat(float64(value)/1e8, 'f', 1, 64)) + "亿"
	case value >= 10_000:
		return trimTrailingZero(strconv.FormatFloat(float64(value)/1e4, 'f', 1, 64)) + "万"
	default:
		return strconv.FormatInt(value, 10)
	}
}

func trimTrailingZero(s string) string {
	return strings.TrimSuffix(s, ".0")
}

var markdownEscaper = strings.NewReplacer(
	`\`, `\\`,
	`[`, `\[`,
	`]`, `\]`,
	`*`, `\*`,
	`_`, `\_`,
	"`", "\\`",
	`<`, `\<`,
	`>`, `\>`,
	"\n", " ",
)

func escapeMarkdown(s string) string {
	return markdownEscaper.Replace(s)
}

func escapeTable(s string) string {
	return strings.ReplaceAll(escapeMarkdown(s), "|", `\|`)
}

func escapeURL(s string) string {
	replacer := strings.NewReplacer(
		" ", "%20",
		"(", "%28",
		")", "%29",
	)
	return replacer.Replace(s)
}

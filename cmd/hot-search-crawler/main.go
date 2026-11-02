// Command hot-search-crawler 抓取微博、抖音、小红书热搜并输出 JSON。
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/solarhell/hot-search-crawler/internal/crawl"
	"github.com/solarhell/hot-search-crawler/internal/source"
	"github.com/solarhell/hot-search-crawler/internal/store"
	"github.com/solarhell/hot-search-crawler/internal/translate"
)

// version 可在构建时通过 -ldflags "-X main.version=v0.1.0" 覆盖。
var version = "dev"

const (
	defaultSources  = "weibo,douyin,xiaohongshu"
	defaultTop      = 20
	defaultTimeout  = 25 * time.Second
	defaultRetries  = 2
	defaultInterval = 30 * time.Minute
)

// 文件名和日报统一使用北京时间，避免服务器时区不同导致目录漂移。
var location = time.FixedZone("Asia/Shanghai", 8*3600)

type config struct {
	sources      string
	top          int
	dataDir      string
	archiveDir   string
	readme       string
	updateReadme bool
	timeout      time.Duration
	retries      int
	translate    bool
	dryRun       bool
	debug        bool
	live         bool
	interval     time.Duration
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	switch args[0] {
	case "crawl":
		return runCrawl(args[1:], stdout, stderr, false)
	case "daemon":
		return runCrawl(args[1:], stdout, stderr, true)
	case "sources":
		return runSources(stdout)
	case "version", "--version", "-v":
		_, _ = fmt.Fprintln(stdout, version)
		return 0
	case "help", "-h", "--help":
		printUsage(stdout)
		return 0
	default:
		_, _ = fmt.Fprintf(stderr, "未知命令 %q\n\n", args[0])
		printUsage(stderr)
		return 2
	}
}

func runCrawl(args []string, stdout, stderr io.Writer, daemon bool) int {
	cfg := &config{}
	fs := flag.NewFlagSet("crawl", flag.ContinueOnError)
	fs.SetOutput(stderr)
	bindCommonFlags(fs, cfg)
	if daemon {
		fs.DurationVar(&cfg.interval, "interval", defaultInterval, "常驻模式下的抓取间隔")
	}
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() > 0 {
		_, _ = fmt.Fprintf(stderr, "多余的参数: %v\n", fs.Args())
		return 2
	}

	logger := newLogger(stderr, cfg.debug)

	sources, err := source.Parse(cfg.sources)
	if err != nil {
		logger.Error("平台配置错误", "error", err)
		return 2
	}

	translator, err := buildTranslator(cfg, sources, logger)
	if err != nil {
		logger.Error("翻译配置错误", "error", err)
		return 2
	}

	crawlOpts := crawl.Options{
		Top:        cfg.top,
		Timeout:    cfg.timeout,
		Retries:    cfg.retries,
		Translator: translator,
		Location:   location,
		Logger:     logger,
	}
	st := store.New(cfg.dataDir, cfg.archiveDir, location)

	if cfg.dryRun {
		snapshot := crawl.Run(context.Background(), sources, crawlOpts)
		data, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			logger.Error("序列化快照失败", "error", err)
			return 1
		}
		_, _ = fmt.Fprintln(stdout, string(data))
		if okCount(snapshot) == 0 {
			return 1
		}
		return 0
	}

	if daemon {
		return runDaemon(sources, crawlOpts, cfg, st, logger)
	}
	return runOnce(context.Background(), sources, crawlOpts, cfg, st, logger, cfg.live)
}

func runOnce(
	ctx context.Context,
	sources []source.Source,
	crawlOpts crawl.Options,
	cfg *config,
	st *store.Store,
	logger *slog.Logger,
	live bool,
) int {
	snapshot := crawl.Run(ctx, sources, crawlOpts)
	okBoards := okCount(snapshot)
	if okBoards == 0 {
		logger.Error("所有平台抓取失败", "boards", len(snapshot.Boards))
	} else {
		logger.Info("抓取完成", "ok_boards", okBoards, "total_boards", len(snapshot.Boards))
	}

	// 实时模式只写 latest.json，不生成小时快照、日报和 README，
	// 避免高频运行时往 git 历史里灌入大量重复数据。
	if live {
		if okBoards == 0 {
			return 1
		}
		latestPath, err := st.WriteLatest(snapshot)
		if err != nil {
			logger.Error("写入实时快照失败", "error", err)
			return 1
		}
		logger.Info("实时快照已写入", "path", latestPath)
		return 0
	}

	snapshotPath, err := st.WriteSnapshot(snapshot)
	if err != nil {
		logger.Error("写入小时快照失败", "error", err)
		return 1
	}
	latestPath, err := st.WriteLatest(snapshot)
	if err != nil {
		logger.Error("写入 latest.json 失败", "error", err)
		return 1
	}

	summary, err := st.BuildSummary(st.Date(snapshot))
	if err != nil {
		logger.Error("聚合日报失败", "error", err)
		return 1
	}
	summaryPath, err := st.WriteSummary(summary)
	if err != nil {
		logger.Error("写入日报失败", "error", err)
		return 1
	}
	logger.Info("数据已写入", "snapshot", snapshotPath, "latest", latestPath, "summary", summaryPath)

	if store.HasData(summary) {
		archivePath, err := st.WriteArchive(summary)
		if err != nil {
			logger.Error("写入归档 Markdown 失败", "error", err)
			return 1
		}
		logger.Info("归档已写入", "path", archivePath)

		if cfg.updateReadme {
			if err := st.UpdateReadme(cfg.readme, summary); err != nil {
				logger.Error("更新 README 失败", "error", err)
				return 1
			}
			logger.Info("README 已更新", "path", cfg.readme)
		}
	}

	if okBoards == 0 {
		return 1
	}
	return 0
}

func runDaemon(
	sources []source.Source,
	crawlOpts crawl.Options,
	cfg *config,
	st *store.Store,
	logger *slog.Logger,
) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	interval := cfg.interval
	if interval <= 0 {
		interval = defaultInterval
	}

	failures := 0
	for {
		if ctx.Err() != nil {
			return 0
		}

		code := runOnce(ctx, sources, crawlOpts, cfg, st, logger, cfg.live)

		var wait time.Duration
		if code == 0 {
			failures = 0
			wait = interval
		} else {
			failures++
			wait = daemonBackoff(failures)
		}

		logger.Info("等待下一轮", "in", wait, "consecutive_failures", failures)
		select {
		case <-ctx.Done():
			logger.Info("收到退出信号")
			return 0
		case <-time.After(wait):
		}
	}
}

func daemonBackoff(failures int) time.Duration {
	wait := time.Minute
	for i := 1; i < failures && wait < 15*time.Minute; i++ {
		wait *= 2
	}
	if wait > 15*time.Minute {
		wait = 15 * time.Minute
	}
	return wait
}

func bindCommonFlags(fs *flag.FlagSet, cfg *config) {
	fs.StringVar(&cfg.sources, "sources", defaultSources, "启用的平台，逗号分隔: weibo,douyin,xiaohongshu")
	fs.IntVar(&cfg.top, "top", defaultTop, "每个平台保留的条数")
	fs.StringVar(&cfg.dataDir, "data", "data", "快照输出目录")
	fs.StringVar(&cfg.archiveDir, "archives", "archives", "日报输出目录")
	fs.StringVar(&cfg.readme, "readme", "README.md", "README 路径")
	fs.BoolVar(&cfg.updateReadme, "update-readme", false, "把今日热榜写进 README")
	fs.DurationVar(&cfg.timeout, "timeout", defaultTimeout, "单次请求超时")
	fs.IntVar(&cfg.retries, "retries", defaultRetries, "失败重试次数")
	fs.BoolVar(&cfg.translate, "translate", false, "启用翻译（读取 TRANSLATE_* 环境变量）")
	fs.BoolVar(&cfg.dryRun, "dry-run", false, "只把快照 JSON 输出到标准输出，不写文件")
	fs.BoolVar(&cfg.live, "live", false, "只写 latest.json，不生成小时快照/日报/README（用于高频发布）")
	fs.BoolVar(&cfg.debug, "debug", false, "输出调试日志")
}

func buildTranslator(cfg *config, sources []source.Source, logger *slog.Logger) (translate.Translator, error) {
	if !cfg.translate {
		return nil, nil
	}

	needed := false
	for _, src := range sources {
		if src.Translatable() {
			needed = true
			break
		}
	}
	if !needed {
		logger.Info("当前启用的平台都是中文平台，翻译不会执行")
		return nil, nil
	}

	baseURL := strings.TrimSpace(os.Getenv("TRANSLATE_BASE_URL"))
	apiKey := strings.TrimSpace(os.Getenv("TRANSLATE_API_KEY"))
	model := strings.TrimSpace(os.Getenv("TRANSLATE_MODEL"))
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if apiKey == "" || model == "" {
		return nil, errors.New("启用 --translate 需要设置 TRANSLATE_API_KEY 和 TRANSLATE_MODEL")
	}

	logger.Info("翻译已启用", "model", model, "base_url", baseURL)
	return translate.NewOpenAI(translate.OpenAIConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   model,
	}), nil
}

func runSources(stdout io.Writer) int {
	_, _ = fmt.Fprintln(stdout, "ID\t名称\t需要翻译")
	for _, src := range source.All() {
		_, _ = fmt.Fprintf(stdout, "%s\t%s\t%t\n", src.ID(), src.Name(), src.Translatable())
	}
	return 0
}

func okCount(snapshot crawl.Snapshot) int {
	count := 0
	for _, board := range snapshot.Boards {
		if board.OK {
			count++
		}
	}
	return count
}

func newLogger(w io.Writer, debug bool) *slog.Logger {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: level}))
}

func printUsage(w io.Writer) {
	_, _ = fmt.Fprint(w, `hot-search-crawler — 微博 / 抖音 / 小红书热搜抓取

用法:
  hot-search-crawler crawl [flags]     抓取一次
  hot-search-crawler daemon [flags]    常驻定时抓取
  hot-search-crawler sources           列出内置平台
  hot-search-crawler version           显示版本

示例:
  hot-search-crawler crawl --update-readme
  hot-search-crawler crawl --sources weibo,douyin --top 10 --dry-run
  hot-search-crawler crawl --live --data dist
  hot-search-crawler daemon --interval 30m

运行 "hot-search-crawler crawl -h" 查看全部参数。
代理通过 HTTP_PROXY / HTTPS_PROXY 环境变量设置。
`)
}

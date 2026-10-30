// Package crawl 并发调度多个数据源，产出一次完整快照。
package crawl

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/solarhell/hot-search-crawler/internal/source"
	"github.com/solarhell/hot-search-crawler/internal/translate"
)

// SchemaVersion 是输出 JSON 的结构版本。
const SchemaVersion = 1

// Board 是单个平台在一次抓取中的结果。抓取失败时 OK 为 false，
// Error 保留原因，Items 为空数组，便于下游统一处理。
type Board struct {
	Platform  string        `json:"platform"`
	Name      string        `json:"name"`
	OK        bool          `json:"ok"`
	FetchedAt time.Time     `json:"fetched_at"`
	ElapsedMS int64         `json:"elapsed_ms"`
	Error     string        `json:"error,omitempty"`
	Items     []source.Item `json:"items"`
}

// Snapshot 是一次抓取的全部平台结果。
type Snapshot struct {
	SchemaVersion int       `json:"schema_version"`
	GeneratedAt   time.Time `json:"generated_at"`
	Timezone      string    `json:"timezone"`
	Boards        []Board   `json:"boards"`
}

// Options 控制抓取行为。
type Options struct {
	// Top 是每个平台保留的最大条目数，默认 20。
	Top int
	// Timeout 是单次请求的超时，默认 25 秒。
	Timeout time.Duration
	// Retries 是单次抓取失败后的重试次数，默认 0。
	Retries int
	// Translator 为空表示不做翻译。
	Translator translate.Translator
	// Location 用于 GeneratedAt 和文件名，为空时使用本地时区。
	Location *time.Location
	// Logger 为空时使用 slog.Default()。
	Logger *slog.Logger
	// Now 用于测试注入，为空时使用 time.Now。
	Now func() time.Time
}

func (o *Options) applyDefaults() {
	if o.Top <= 0 {
		o.Top = 20
	}
	if o.Timeout <= 0 {
		o.Timeout = 25 * time.Second
	}
	if o.Retries < 0 {
		o.Retries = 0
	}
	if o.Location == nil {
		o.Location = time.Local
	}
	if o.Logger == nil {
		o.Logger = slog.Default()
	}
	if o.Now == nil {
		o.Now = time.Now
	}
}

// Run 并发抓取所有数据源，返回按数据源顺序排列的快照。
// 单个平台失败不会中断其他平台。
func Run(ctx context.Context, sources []source.Source, opts Options) Snapshot {
	opts.applyDefaults()

	snapshot := Snapshot{
		SchemaVersion: SchemaVersion,
		GeneratedAt:   opts.Now().In(opts.Location),
		Timezone:      opts.Location.String(),
		Boards:        make([]Board, len(sources)),
	}

	var wg sync.WaitGroup
	for i, src := range sources {
		wg.Add(1)
		go func(i int, src source.Source) {
			defer wg.Done()
			snapshot.Boards[i] = fetchOne(ctx, src, opts)
		}(i, src)
	}
	wg.Wait()

	return snapshot
}

func fetchOne(ctx context.Context, src source.Source, opts Options) Board {
	start := opts.Now()
	board := Board{
		Platform: src.ID(),
		Name:     src.Name(),
		Items:    []source.Item{},
	}

	logger := opts.Logger.With("platform", src.ID())
	logger.Debug("开始抓取")

	items, err := fetchWithRetry(ctx, src, opts, logger)
	board.FetchedAt = opts.Now().In(opts.Location)
	board.ElapsedMS = opts.Now().Sub(start).Milliseconds()

	if err != nil {
		board.Error = err.Error()
		if source.IsThrottled(err) {
			logger.Warn("抓取被限流", "error", err)
		} else {
			logger.Error("抓取失败", "error", err, "elapsed_ms", board.ElapsedMS)
		}
		return board
	}

	if len(items) == 0 {
		board.Error = "返回数据为空"
		logger.Warn("抓取结果为空")
		return board
	}

	if opts.Translator != nil && src.Translatable() {
		items = translateItems(ctx, opts, logger, items)
	}

	if len(items) > opts.Top {
		items = items[:opts.Top]
	}
	for i := range items {
		items[i].Rank = i + 1
	}

	board.OK = true
	board.Items = items
	logger.Info("抓取成功", "count", len(items), "elapsed_ms", board.ElapsedMS)
	return board
}

func fetchWithRetry(
	ctx context.Context,
	src source.Source,
	opts Options,
	logger *slog.Logger,
) ([]source.Item, error) {
	attempts := opts.Retries + 1
	var lastErr error

	for attempt := 0; attempt < attempts; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, opts.Timeout)
		items, err := src.Fetch(attemptCtx, opts.Top)
		cancel()

		if err == nil {
			return items, nil
		}
		lastErr = err

		// 限流时连续重试只会让平台把门关得更久，直接交给调用方。
		if source.IsThrottled(err) {
			return nil, err
		}
		if attempt == attempts-1 {
			break
		}

		wait := retryWait(attempt)
		logger.Warn("抓取失败，准备重试",
			"attempt", attempt+1,
			"retry_in", wait,
			"error", err,
		)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(wait):
		}
	}

	return nil, lastErr
}

// retryBase 是指数退避的起点，测试中会调小。
var retryBase = 2 * time.Second

func retryWait(attempt int) time.Duration {
	wait := retryBase
	for i := 0; i < attempt && wait < 10*time.Second; i++ {
		wait *= 2
	}
	if wait > 10*time.Second {
		wait = 10 * time.Second
	}
	return wait
}

func translateItems(
	ctx context.Context,
	opts Options,
	logger *slog.Logger,
	items []source.Item,
) []source.Item {
	titles := make([]string, len(items))
	for i := range items {
		titles[i] = items[i].Title
	}

	translateCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	translated, err := opts.Translator.Translate(translateCtx, titles)
	if err != nil {
		// 翻译是尽力而为，失败时保留原文。
		logger.Warn("翻译失败，保留原文", "error", err)
		return items
	}
	for i := range items {
		if i < len(translated) {
			items[i].Translated = translated[i]
		}
	}
	return items
}

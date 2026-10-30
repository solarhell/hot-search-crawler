package source

import (
	"context"
	"fmt"
	"time"

	"github.com/enetx/g"
	"github.com/enetx/surf"
)

// defaultTimeout 单次 HTTP 请求的超时。三个平台的接口都在国内，
// 20 秒足够覆盖 TLS 握手和偶发抖动。
const defaultTimeout = 20 * time.Second

type requestOptions struct {
	// impersonate 使用 Chrome 的 TLS/HTTP2 指纹。微博和抖音会对
	// 明显的 Go 客户端特征返回空数据或风控页。
	impersonate bool
	headers     map[string]string
}

// httpStatusError 表示服务端返回了非 2xx 状态码。
type httpStatusError struct {
	StatusCode int
	Body       string
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Body)
}

// getBody 发起 GET 请求并返回响应体。统一处理请求错误、读取错误和状态码。
func getBody(ctx context.Context, url string, opts requestOptions) ([]byte, error) {
	builder := surf.NewClient().Builder().Timeout(defaultTimeout)
	if opts.impersonate {
		builder = builder.Impersonate().Chrome()
	}

	client := builder.Build().Unwrap()
	defer func() { _ = client.Close() }()

	req := client.Get(g.String(url)).WithContext(ctx)
	for key, value := range opts.headers {
		req = req.AddHeaders(key, value)
	}

	resp := req.Do()
	if resp.IsErr() {
		return nil, fmt.Errorf("请求失败: %w", resp.Err())
	}

	ok := resp.Ok()
	body := ok.Body.String()
	if body.IsErr() {
		return nil, fmt.Errorf("读取响应体失败: %w", body.Err())
	}

	data := []byte(body.Ok())
	if ok.StatusCode < 200 || ok.StatusCode >= 300 {
		return nil, &httpStatusError{
			StatusCode: int(ok.StatusCode),
			Body:       truncate(string(data), 512),
		}
	}
	return data, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

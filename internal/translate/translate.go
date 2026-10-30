// Package translate 提供可选的热搜标题翻译能力。
//
// 当前内置的三个国内平台不需要翻译；保留这个接口，是为了让新增海外
// 数据源时可以直接复用，而不必把某个 LLM 厂商的 SDK 绑进爬虫。
package translate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Translator 把一个标题列表翻译成中文，返回结果与输入按索引对齐。
type Translator interface {
	Translate(ctx context.Context, texts []string) ([]string, error)
}

const systemPrompt = `你是热搜词翻译助手。用户会给出一个待翻译的热搜词列表（JSON 数组），请把每个热搜词翻译成简体中文，并以 JSON 对象返回。
对象只包含原文到译文的映射，不要输出解释或 Markdown 代码块。`

// OpenAIConfig 配置一个 OpenAI 兼容的 Chat Completions 接口。
type OpenAIConfig struct {
	// BaseURL 例如 https://api.openai.com/v1，也接受完整的 /chat/completions 地址。
	BaseURL string
	APIKey  string
	Model   string
	// HTTPClient 为空时使用 30 秒超时的默认客户端。
	HTTPClient *http.Client
}

// OpenAI 通过 OpenAI 兼容接口批量翻译。
type OpenAI struct {
	endpoint string
	apiKey   string
	model    string
	client   *http.Client
}

// NewOpenAI 创建翻译器。
func NewOpenAI(cfg OpenAIConfig) *OpenAI {
	endpoint := strings.TrimRight(cfg.BaseURL, "/")
	if !strings.HasSuffix(endpoint, "/chat/completions") {
		endpoint += "/chat/completions"
	}
	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	return &OpenAI{
		endpoint: endpoint,
		apiKey:   cfg.APIKey,
		model:    cfg.Model,
		client:   client,
	}
}

type chatRequest struct {
	Model          string        `json:"model"`
	Messages       []chatMessage `json:"messages"`
	Temperature    float64       `json:"temperature"`
	ResponseFormat struct {
		Type string `json:"type"`
	} `json:"response_format"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Translate 调用模型并解析 JSON 映射。模型漏掉的标题对应位置返回空串。
func (t *OpenAI) Translate(ctx context.Context, texts []string) ([]string, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("待翻译标题为空")
	}

	input, err := json.Marshal(texts)
	if err != nil {
		return nil, fmt.Errorf("序列化待翻译标题失败: %w", err)
	}

	reqBody := chatRequest{
		Model:       t.model,
		Temperature: 0,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: string(input)},
		},
	}
	reqBody.ResponseFormat.Type = "json_object"

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化翻译请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.endpoint, strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("创建翻译请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.apiKey)

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求翻译接口失败: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("读取翻译响应失败: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("翻译接口返回 HTTP %d: %s", resp.StatusCode, truncate(string(body), 256))
	}

	var parsed chatResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("解析翻译响应失败: %w", err)
	}
	if parsed.Error != nil {
		return nil, fmt.Errorf("翻译接口返回错误: %s", parsed.Error.Message)
	}
	if len(parsed.Choices) == 0 {
		return nil, fmt.Errorf("翻译接口返回空结果")
	}

	content := strings.TrimSpace(parsed.Choices[0].Message.Content)
	content = stripCodeFence(content)

	var mapping map[string]string
	if err := json.Unmarshal([]byte(content), &mapping); err != nil {
		return nil, fmt.Errorf("解析翻译结果失败: %w", err)
	}

	results := make([]string, len(texts))
	for i, text := range texts {
		results[i] = mapping[text]
	}
	return results, nil
}

func stripCodeFence(s string) string {
	if !strings.HasPrefix(s, "```") {
		return s
	}
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

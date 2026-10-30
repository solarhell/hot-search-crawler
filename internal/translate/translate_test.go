package translate

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenAITranslate(t *testing.T) {
	var gotAuth string
	var gotBody chatRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("解析请求失败: %v", err)
		}
		content := `{"hello":"你好","world":"世界"}`
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{"content": content}},
			},
		})
	}))
	defer server.Close()

	translator := NewOpenAI(OpenAIConfig{
		BaseURL: server.URL + "/v1",
		APIKey:  "test-key",
		Model:   "test-model",
	})

	results, err := translator.Translate(context.Background(), []string{"hello", "world", "missing"})
	if err != nil {
		t.Fatalf("Translate: %v", err)
	}

	if gotAuth != "Bearer test-key" {
		t.Errorf("Authorization = %q", gotAuth)
	}
	if gotBody.Model != "test-model" {
		t.Errorf("model = %q", gotBody.Model)
	}
	if gotBody.ResponseFormat.Type != "json_object" {
		t.Errorf("response_format = %q", gotBody.ResponseFormat.Type)
	}
	if len(results) != 3 || results[0] != "你好" || results[1] != "世界" || results[2] != "" {
		t.Errorf("results = %v", results)
	}
}

func TestOpenAITranslateStripsCodeFence(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content := "```json\n{\"hello\":\"你好\"}\n```"
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{"content": content}},
			},
		})
	}))
	defer server.Close()

	translator := NewOpenAI(OpenAIConfig{BaseURL: server.URL, APIKey: "k", Model: "m"})
	results, err := translator.Translate(context.Background(), []string{"hello"})
	if err != nil {
		t.Fatalf("Translate: %v", err)
	}
	if results[0] != "你好" {
		t.Errorf("results[0] = %q", results[0])
	}
}

func TestOpenAITranslateErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"bad key"}}`))
	}))
	defer server.Close()

	translator := NewOpenAI(OpenAIConfig{BaseURL: server.URL, APIKey: "k", Model: "m"})
	_, err := translator.Translate(context.Background(), []string{"hello"})
	if err == nil {
		t.Fatal("HTTP 401 应返回错误")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("error = %v", err)
	}
}

func TestOpenAITranslateEmpty(t *testing.T) {
	translator := NewOpenAI(OpenAIConfig{})
	if _, err := translator.Translate(context.Background(), nil); err == nil {
		t.Fatal("空输入应返回错误")
	}
}

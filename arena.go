package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Arena runs the model comparison.
type Arena struct {
	models []ModelConfig
	client *http.Client
}

func NewArena(models []ModelConfig) *Arena {
	return &Arena{
		models: models,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

// Run sends the prompt to all models concurrently.
func (a *Arena) Run(prompt string, stream bool) []ArenaResult {
	var wg sync.WaitGroup
	results := make([]ArenaResult, len(a.models))
	msg := ChatMessage{Role: "user", Content: prompt}

	for i, m := range a.models {
		wg.Add(1)
		go func(idx int, cfg ModelConfig) {
			defer wg.Done()
			results[idx] = a.callModel(cfg, msg, stream)
		}(i, m)
	}
	wg.Wait()
	return results
}

func (a *Arena) callModel(cfg ModelConfig, msg ChatMessage, stream bool) ArenaResult {
	start := time.Now()
	res := ArenaResult{Model: cfg.Name, Alias: cfg.Alias, Provider: cfg.Provider}

	apiKey := expandEnv(cfg.APIKey)
	if apiKey == "" || strings.HasPrefix(apiKey, "${") {
		res.Error = fmt.Sprintf("API key not set for %s", cfg.Name)
		return res
	}

	body, _ := json.Marshal(ChatRequest{
		Model:    cfg.Name,
		Messages: []ChatMessage{msg},
		Stream:   stream,
	})

	url := strings.TrimRight(cfg.BaseURL, "/") + "/chat/completions"
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	httpResp, err := a.client.Do(req)
	if err != nil {
		res.Error = fmt.Sprintf("request failed: %v", err)
		return res
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		respBody, _ := io.ReadAll(httpResp.Body)
		res.Error = fmt.Sprintf("HTTP %d: %s", httpResp.StatusCode, strings.TrimSpace(string(respBody)))
		return res
	}

	if stream {
		res.Content = a.readStream(httpResp.Body)
	} else {
		var chatResp ChatResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&chatResp); err != nil {
			res.Error = fmt.Sprintf("decode: %v", err)
			return res
		}
		if len(chatResp.Choices) > 0 {
			res.Content = chatResp.Choices[0].Message.Content
		}
		if chatResp.Usage != nil {
			res.PromptTokens = chatResp.Usage.PromptTokens
			res.CompletionTokens = chatResp.Usage.CompletionTokens
			res.TotalTokens = chatResp.Usage.TotalTokens
		}
	}

	res.LatencyMs = time.Since(start).Milliseconds()
	if res.TotalTokens == 0 {
		res.TotalTokens = res.PromptTokens + res.CompletionTokens
	}
	res.Success = true
	res.CalculateCost()
	return res
}

func (a *Arena) readStream(body io.ReadCloser) string {
	defer body.Close()
	var content strings.Builder
	scanner := bufio.NewScanner(body)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		json.Unmarshal([]byte(data), &chunk)
		for _, c := range chunk.Choices {
			content.WriteString(c.Delta.Content)
		}
	}
	return content.String()
}

func expandEnv(s string) string {
	return os.Expand(s, func(name string) string {
		if v := os.Getenv(name); v != "" {
			return v
		}
		return os.Getenv(name)
	})
}

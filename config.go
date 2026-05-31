package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ArenaConfig is the top-level YAML config structure.
type ArenaConfig struct {
	Models []ModelConfig `yaml:"models"`
}

// DefaultModels returns the built-in model configurations using environment variables.
func DefaultModels() []ModelConfig {
	return []ModelConfig{
		{
			Name: "deepseek-chat", Alias: "ds", Provider: "deepseek",
			BaseURL: "https://api.deepseek.com",
			APIKey:  "${DEEPSEEK_API_KEY}",
		},
		{
			Name: "qwen-plus", Alias: "qwen", Provider: "alibaba",
			BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
			APIKey:  "${QWEN_API_KEY}",
		},
		{
			Name: "glm-4-flash", Alias: "glm", Provider: "zhipu",
			BaseURL: "https://open.bigmodel.cn/api/paas/v4",
			APIKey:  "${ZHIPU_API_KEY}",
		},
		{
			Name: "ernie-speed", Alias: "baidu", Provider: "baidu",
			BaseURL: "https://qianfan.baidubce.com/v2",
			APIKey:  "${BAIDU_API_KEY}",
		},
		{
			Name: "doubao-pro-32k", Alias: "volc", Provider: "volcengine",
			BaseURL: "https://ark.cn-beijing.volces.com/api/v3",
			APIKey:  "${VOLC_API_KEY}",
		},
	}
}

// LoadConfig loads model configuration from a YAML file.
func LoadConfig(path string) ([]ModelConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	// Expand environment variables in the raw YAML
	expanded := os.Expand(string(data), func(name string) string {
		if v := os.Getenv(name); v != "" {
			return v
		}
		return "${" + name + "}"
	})

	var cfg ArenaConfig
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if len(cfg.Models) == 0 {
		return nil, fmt.Errorf("no models defined in config")
	}

	return cfg.Models, nil
}

// splitAndTrim splits a comma-separated string and trims whitespace.
func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

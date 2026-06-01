package main

import (
	"os"
	"testing"
)

func TestDefaultModels(t *testing.T) {
	models := DefaultModels()
	if len(models) == 0 {
		t.Fatal("DefaultModels() returned empty")
	}
	// Should include all major providers
	providers := make(map[string]bool)
	for _, m := range models {
		providers[m.Provider] = true
	}
	for _, p := range []string{"deepseek", "alibaba", "zhipu", "baidu", "volcengine"} {
		if !providers[p] {
			t.Errorf("missing provider: %s", p)
		}
	}
}

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"a,b,c", 3},
		{"a, b, c", 3},
		{"single", 1},
		{"", 0},
		{",,", 0},
	}
	for _, tt := range tests {
		got := splitAndTrim(tt.input)
		if len(got) != tt.want {
			t.Errorf("splitAndTrim(%q) = %d items, want %d", tt.input, len(got), tt.want)
		}
	}
}

func TestExpandEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "hello")
	defer os.Unsetenv("TEST_VAR")

	got := expandEnv("${TEST_VAR}")
	if got != "hello" {
		t.Errorf("expandEnv = %q, want %q", got, "hello")
	}

	// Unset variable returns empty string
	got2 := expandEnv("${NONEXISTENT}")
	if got2 != "" {
		t.Errorf("expandEnv(nonexistent) = %q, want empty string", got2)
	}
}

func TestMin(t *testing.T) {
	if min(5, 3) != 3 {
		t.Errorf("min(5,3) = %d, want 3", min(5, 3))
	}
	if min(3, 5) != 3 {
		t.Errorf("min(3,5) = %d, want 3", min(3, 5))
	}
	if min(-1, 1) != -1 {
		t.Errorf("min(-1,1) = %d, want -1", min(-1, 1))
	}
}

func TestGetPricing(t *testing.T) {
	// Known model
	p := GetPricing("deepseek-chat")
	if p == nil {
		t.Fatal("GetPricing(deepseek-chat) returned nil")
	}
	if p.InputPrice <= 0 {
		t.Errorf("deepseek-chat input price should be > 0, got %f", p.InputPrice)
	}

	// Unknown model
	p2 := GetPricing("nonexistent-model")
	if p2 != nil {
		t.Errorf("GetPricing(nonexistent) = %v, want nil", p2)
	}
}

func TestCalculateCost(t *testing.T) {
	r := &ArenaResult{
		Model:            "deepseek-chat",
		PromptTokens:     1000,
		CompletionTokens: 500,
		Success:          true,
	}
	r.CalculateCost()
	if r.TotalCost <= 0 {
		t.Errorf("CalculateCost should produce positive cost, got %f", r.TotalCost)
	}
	// deepseek-chat: input 0.001/1K, output 0.002/1K
	// cost = (1000 * 0.001 / 1000) + (500 * 0.002 / 1000) = 0.001 + 0.001 = 0.002
	expected := 0.002
	if r.TotalCost != expected {
		t.Errorf("CalculateCost = %f, want %f", r.TotalCost, expected)
	}
	if r.InputCost != 0.001 {
		t.Errorf("InputCost = %f, want 0.001", r.InputCost)
	}
	if r.OutputCost != 0.001 {
		t.Errorf("OutputCost = %f, want 0.001", r.OutputCost)
	}
}

func TestCalculateCostUnknownModel(t *testing.T) {
	r := &ArenaResult{
		Model:            "unknown-model",
		PromptTokens:     1000,
		CompletionTokens: 500,
		Success:          true,
	}
	r.CalculateCost()
	if r.TotalCost != 0 {
		t.Errorf("CalculateCost(unknown) = %f, want 0", r.TotalCost)
	}
}

func TestNewArena(t *testing.T) {
	models := []ModelConfig{
		{Name: "test-model", Provider: "test", BaseURL: "https://test.com", APIKey: "${TEST_KEY}"},
	}
	a := NewArena(models)
	if a == nil {
		t.Fatal("NewArena returned nil")
	}
	if len(a.models) != 1 {
		t.Errorf("NewArena has %d models, want 1", len(a.models))
	}
}

func TestModelConfigFiltering(t *testing.T) {
	// Simulate the filtering logic from main.go
	allModels := []ModelConfig{
		{Name: "deepseek-chat", Alias: "ds"},
		{Name: "qwen-plus", Alias: "qwen"},
		{Name: "glm-4-flash", Alias: "glm"},
	}

	wanted := map[string]bool{"ds": true, "qwen": true}
	var filtered []ModelConfig
	for _, m := range allModels {
		if wanted[m.Name] || wanted[m.Alias] {
			filtered = append(filtered, m)
		}
	}
	if len(filtered) != 2 {
		t.Errorf("filtered = %d models, want 2", len(filtered))
	}
}

func TestColorOutput(t *testing.T) {
	// Test colored output
	useColor = true
	g := green("test")
	if g == "" {
		t.Error("green() returned empty")
	}

	// Test no-color mode
	useColor = false
	g2 := green("test")
	if g2 != "test" {
		t.Errorf("green(no-color) = %q, want %q", g2, "test")
	}
	useColor = true
}

package main

// Pricing 模型定价（¥/1K tokens）
// 数据来源：各厂商官网公开定价（2026年5月更新）
type Pricing struct {
	InputPrice  float64 `json:"input_price"`  // ¥ per 1K input tokens
	OutputPrice float64 `json:"output_price"` // ¥ per 1K output tokens
}

// pricingTable 国产主流模型定价表
var pricingTable = map[string]Pricing{
	// DeepSeek
	"deepseek-chat":            {InputPrice: 0.001, OutputPrice: 0.002},
	"deepseek-reasoner":        {InputPrice: 0.004, OutputPrice: 0.016},
	"deepseek-v4-pro-260425":   {InputPrice: 0.004, OutputPrice: 0.016},
	"deepseek-v4-flash-260425": {InputPrice: 0.001, OutputPrice: 0.002},

	// 通义千问
	"qwen-turbo-latest":    {InputPrice: 0.0008, OutputPrice: 0.002},
	"qwen-plus-latest":     {InputPrice: 0.0008, OutputPrice: 0.002},
	"qwen-max-latest":      {InputPrice: 0.002, OutputPrice: 0.006},
	"qwen2.5-72b-instruct": {InputPrice: 0.004, OutputPrice: 0.012},

	// 智谱 GLM
	"glm-4-plus":  {InputPrice: 0.005, OutputPrice: 0.005},
	"glm-4-air":   {InputPrice: 0.001, OutputPrice: 0.001},
	"glm-4-flash": {InputPrice: 0.0001, OutputPrice: 0.0001},

	// 百度 ERNIE
	"ernie-4.0-8k-latest": {InputPrice: 0.012, OutputPrice: 0.012},
	"ernie-3.5-8k-latest": {InputPrice: 0.0008, OutputPrice: 0.002},

	// 火山引擎 Doubao
	"doubao-pro-32k":                 {InputPrice: 0.0008, OutputPrice: 0.002},
	"doubao-seed-2-0-pro-260215":     {InputPrice: 0.005, OutputPrice: 0.02},
	"doubao-seed-1-6-flash-250615":   {InputPrice: 0.0005, OutputPrice: 0.001},
	"doubao-seed-1-6-250615":         {InputPrice: 0.0008, OutputPrice: 0.002},
	"doubao-seed-2-0-mini-260215":    {InputPrice: 0.001, OutputPrice: 0.004},
	"doubao-seed-2-0-lite-260215":    {InputPrice: 0.0003, OutputPrice: 0.0006},
	"doubao-1-5-thinking-pro-250415": {InputPrice: 0.004, OutputPrice: 0.016},
	"doubao-seed-1-8-251228":         {InputPrice: 0.0008, OutputPrice: 0.002},

	// 百度 ERNIE (alias)
	"ernie-speed":      {InputPrice: 0.0008, OutputPrice: 0.002},
	"ernie-speed-128k": {InputPrice: 0.0015, OutputPrice: 0.003},
	"ernie-lite-8k":    {InputPrice: 0.0, OutputPrice: 0.0}, // 免费
}

// GetPricing 获取模型定价，找不到返回 nil
func GetPricing(modelName string) *Pricing {
	if p, ok := pricingTable[modelName]; ok {
		return &p
	}
	return nil
}

// CalculateCost 根据 token 用量计算成本（元）
func (r *ArenaResult) CalculateCost() {
	if p := GetPricing(r.Model); p != nil {
		inputCost := p.InputPrice * float64(r.PromptTokens) / 1000
		outputCost := p.OutputPrice * float64(r.CompletionTokens) / 1000
		r.InputCost = inputCost
		r.OutputCost = outputCost
		r.TotalCost = inputCost + outputCost
	}
}

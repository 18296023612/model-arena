package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	prompt := runCmd.String("prompt", "", "Prompt to send to all models")
	models := runCmd.String("models", "", "Comma-separated model list (e.g. deepseek-chat,qwen-plus)")
	config := runCmd.String("config", "", "Config file path (YAML)")
	stream := runCmd.Bool("stream", false, "Stream output in real-time")
	jsonOut := runCmd.Bool("json", false, "JSON output")

	flag.Usage = func() {
		fmt.Println(bold("Model Arena v" + version))
		fmt.Println()
		fmt.Println(dim("多模型实时对比工具 — 同一个 prompt，看谁回答好、谁快、谁便宜"))
		fmt.Println()
		fmt.Println(bold("Usage:"))
		fmt.Println("  model-arena run --prompt <text> [options]")
		fmt.Println("  model-arena config                    Show default config template")
		fmt.Println("  model-arena version                   Show version")
		fmt.Println()
		fmt.Println(bold("Examples:"))
		fmt.Println(`  model-arena run --prompt "你好"`)
		fmt.Println(`  model-arena run --prompt "写首诗" --models deepseek-chat,qwen-plus --stream`)
		fmt.Println(`  model-arena run --prompt "Hi" --config arena.yaml --json`)
		fmt.Println()
		fmt.Println(bold("Environment:"))
		fmt.Println("  DEEPSEEK_API_KEY    API key for DeepSeek")
		fmt.Println("  QWEN_API_KEY        API key for Alibaba Qwen")
		fmt.Println("  ZHIPU_API_KEY       API key for Zhipu GLM")
		fmt.Println("  BAIDU_API_KEY       API key for Baidu ERNIE")
		fmt.Println("  VOLC_API_KEY        API key for Volcengine Doubao")
		fmt.Println()
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "run":
		runCmd.Parse(os.Args[2:])
		cmdRun(*prompt, *models, *config, *stream, *jsonOut)
	case "config":
		printConfigTemplate()
	case "version", "--version", "-V":
		fmt.Printf("model-arena v%s\n", version)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		flag.Usage()
		os.Exit(1)
	}
}

func cmdRun(prompt, models, config string, stream, jsonOut bool) {
	if prompt == "" {
		fmt.Fprintln(os.Stderr, red("Error: --prompt is required"))
		os.Exit(1)
	}

	// Load model configuration
	var modelCfgs []ModelConfig
	if config != "" {
		var err error
		modelCfgs, err = LoadConfig(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, red("Error loading config: %v\n"), err)
			os.Exit(1)
		}
	} else {
		modelCfgs = DefaultModels()
	}

	// Filter by --models if specified
	if models != "" {
		wanted := splitAndTrim(models)
		filtered := make([]ModelConfig, 0)
		for _, m := range modelCfgs {
			for _, w := range wanted {
				if m.Name == w || m.Alias == w {
					filtered = append(filtered, m)
					break
				}
			}
		}
		modelCfgs = filtered
	}

	if len(modelCfgs) == 0 {
		fmt.Fprintln(os.Stderr, red("Error: no models configured"))
		fmt.Fprintln(os.Stderr, dim("  Set environment variables or use --config"))
		os.Exit(1)
	}

	// Run the arena
	arena := NewArena(modelCfgs)
	results := arena.Run(prompt, stream)

	if jsonOut {
		printJSON(results)
	} else {
		printResults(results)
	}
}

func printConfigTemplate() {
	fmt.Println(bold("# Model Arena Configuration"))
	fmt.Println("# Save as arena.yaml and use: model-arena run --prompt \"hi\" --config arena.yaml")
	fmt.Println()
	fmt.Println(`models:
  - name: deepseek-chat
    alias: ds
    provider: deepseek
    base_url: https://api.deepseek.com
    api_key: ${DEEPSEEK_API_KEY}

  - name: qwen-plus
    alias: qwen
    provider: alibaba
    base_url: https://dashscope.aliyuncs.com/compatible-mode/v1
    api_key: ${QWEN_API_KEY}

  - name: glm-4-flash
    alias: glm
    provider: zhipu
    base_url: https://open.bigmodel.cn/api/paas/v4
    api_key: ${ZHIPU_API_KEY}

  - name: ernie-speed
    alias: baidu
    provider: baidu
    base_url: https://qianfan.baidubce.com/v2
    api_key: ${BAIDU_API_KEY}

  - name: doubao-pro-32k
    alias: volc
    provider: volcengine
    base_url: https://ark.cn-beijing.volces.com/api/v3
    api_key: ${VOLC_API_KEY}
`)
}

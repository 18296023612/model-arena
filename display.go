package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ── ANSI color helpers ──

var useColor = true

func init() {
	if os.Getenv("NO_COLOR") != "" {
		useColor = false
	}
}

func color(s, code string) string {
	if !useColor {
		return s
	}
	return "\033[" + code + "m" + s + "\033[0m"
}

func green(s string) string  { return color(s, "32") }
func red(s string) string    { return color(s, "31") }
func yellow(s string) string { return color(s, "33") }
func blue(s string) string   { return color(s, "34") }
func bold(s string) string   { return color(s, "1") }
func dim(s string) string    { return color(s, "2") }

// ── Output ──

func printResults(results []ArenaResult) {
	fmt.Println()
	fmt.Println(bold("╔══════════════════════════════════════════════════════════════╗"))
	fmt.Println(bold("║             🤖  Model Arena — 模型对比报告               ║"))
	fmt.Println(bold("╚══════════════════════════════════════════════════════════════╝"))
	fmt.Println()

	// Summary table
	printSummaryTable(results)
	fmt.Println()

	// Detailed outputs
	for i, r := range results {
		if i > 0 {
			fmt.Println(dim(strings.Repeat("─", 72)))
		}
		printModelOutput(r)
		fmt.Println()
	}
}

func printSummaryTable(results []ArenaResult) {
	// Header
	fmt.Printf("  %-18s %-8s %-10s %-7s %-9s %s\n",
		bold("Model"), bold("Status"), bold("Latency"), bold("Tokens"), bold("Cost(¥)"), bold("Provider"))
	fmt.Printf("  %s\n", dim(strings.Repeat("─", 80)))

	// Rows
	for _, r := range results {
		status := green("✅")
		latency := fmt.Sprintf("%dms", r.LatencyMs)
		tokens := fmt.Sprintf("%d", r.TotalTokens)
		cost := dim("-")
		name := r.Model

		if !r.Success {
			status = red("❌")
			latency = red(r.Error[:min(len(r.Error), 18)])
			tokens = "-"
		} else if r.TotalCost > 0 {
			costFmt := fmt.Sprintf("%.4f", r.TotalCost)
			if r.TotalCost < 0.001 {
				cost = green(costFmt)
			} else if r.TotalCost < 0.01 {
				cost = yellow(costFmt)
			} else {
				cost = red(costFmt)
			}
		}

		if r.LatencyMs > 5000 && r.Success {
			latency = yellow(latency)
		}

		alias := r.Alias
		if alias != "" {
			name = alias
		}

		fmt.Printf("  %-18s %-8s %-10s %-7s %-9s %s\n",
			name, status, latency, tokens, cost, dim(r.Provider))
	}
}

func printModelOutput(r ArenaResult) {
	name := r.Model
	if r.Alias != "" {
		name = fmt.Sprintf("%s (%s)", r.Alias, r.Model)
	}

	if !r.Success {
		fmt.Printf("  %s %s\n", red("✖"), bold(name))
		fmt.Printf("  %s\n", red("  Error: "+r.Error))
		return
	}

	fmt.Printf("  %s %s  %s  %s\n",
		green("✓"), bold(name),
		dim(fmt.Sprintf("%dms", r.LatencyMs)),
		dim(fmt.Sprintf("%dtokens", r.TotalTokens)))

	// Truncate long output for display
	content := r.Content
	if len(content) > 300 {
		content = content[:300] + dim("...")
	}
	fmt.Println()
	for _, line := range strings.Split(content, "\n") {
		fmt.Printf("    %s\n", line)
	}
	fmt.Println()
}

func printJSON(results []ArenaResult) {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

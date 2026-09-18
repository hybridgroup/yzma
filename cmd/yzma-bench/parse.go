package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// Units that the benchmarks report with b.ReportMetric.
const (
	unitTokensPerSecond = "tokens/s"
	unitTTFT            = "ttft_ms"
	unitTotal           = "total_ms"
	unitPromptTokens    = "prompt_tokens"
)

// benchmarkResult is what the output of one go test run gives.
type benchmarkResult struct {
	arch            string
	cpu             string
	llamaCPP        string
	runs            []float64
	tokensPerSecond float64
	ttftMs          float64
	totalMs         float64
	promptTokens    float64
}

// parseBenchmark reads the output of go test -bench. It takes the median of
// each metric on its own. Only tokens/s is needed. The comparison suite adds
// the time to the first token and the time of a whole request.
func parseBenchmark(raw string) (benchmarkResult, error) {
	var result benchmarkResult
	passed := false
	var ttft, total, prompt []float64

	for line := range strings.SplitSeq(raw, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "goarch:"):
			result.arch = strings.TrimSpace(strings.TrimPrefix(line, "goarch:"))
		case strings.HasPrefix(line, "cpu:"):
			result.cpu = strings.TrimSpace(strings.TrimPrefix(line, "cpu:"))
		// The WebAssembly benchmarks put the tag of their build in the output,
		// because it does not come from the library directory of the machine.
		case strings.HasPrefix(line, "llama.cpp:"):
			result.llamaCPP = strings.TrimSpace(strings.TrimPrefix(line, "llama.cpp:"))
		case strings.HasPrefix(line, "PASS"):
			passed = true
		case strings.HasPrefix(line, "Benchmark"):
			if value, ok := metricValue(line, unitTokensPerSecond); ok {
				result.runs = append(result.runs, value)
			}
			if value, ok := metricValue(line, unitTTFT); ok {
				ttft = append(ttft, value)
			}
			if value, ok := metricValue(line, unitTotal); ok {
				total = append(total, value)
			}
			if value, ok := metricValue(line, unitPromptTokens); ok {
				prompt = append(prompt, value)
			}
		}
	}

	if !passed {
		return result, fmt.Errorf("the output has no PASS line, the run failed")
	}
	if len(result.runs) == 0 {
		return result, fmt.Errorf("the output has no tokens/s value")
	}

	result.tokensPerSecond = median(result.runs)
	if len(ttft) > 0 {
		result.ttftMs = median(ttft)
	}
	if len(total) > 0 {
		result.totalMs = median(total)
	}
	if len(prompt) > 0 {
		result.promptTokens = median(prompt)
	}

	return result, nil
}

// metricValue takes the value before the given unit of a benchmark line.
func metricValue(line, unit string) (float64, bool) {
	fields := strings.Fields(line)
	for i, f := range fields {
		if f != unit || i == 0 {
			continue
		}

		value, err := strconv.ParseFloat(fields[i-1], 64)
		if err != nil {
			return 0, false
		}

		return value, true
	}

	return 0, false
}

func median(values []float64) float64 {
	sorted := slices.Clone(values)
	slices.Sort(sorted)

	middle := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return sorted[middle]
	}

	return (sorted[middle-1] + sorted[middle]) / 2
}

package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// benchmarkResult is what the output of one go test run gives.
type benchmarkResult struct {
	arch            string
	cpu             string
	runs            []float64
	tokensPerSecond float64
}

// parseBenchmark reads the output of go test -bench. It takes the median of the
// tokens/s values, which the benchmarks report with b.ReportMetric.
func parseBenchmark(raw string) (benchmarkResult, error) {
	var result benchmarkResult
	passed := false

	for line := range strings.SplitSeq(raw, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "goarch:"):
			result.arch = strings.TrimSpace(strings.TrimPrefix(line, "goarch:"))
		case strings.HasPrefix(line, "cpu:"):
			result.cpu = strings.TrimSpace(strings.TrimPrefix(line, "cpu:"))
		case strings.HasPrefix(line, "PASS"):
			passed = true
		case strings.HasPrefix(line, "Benchmark"):
			value, ok := tokensPerSecond(line)
			if ok {
				result.runs = append(result.runs, value)
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

	return result, nil
}

// tokensPerSecond takes the value before the tokens/s unit of a benchmark line.
func tokensPerSecond(line string) (float64, bool) {
	fields := strings.Fields(line)
	for i, f := range fields {
		if f != "tokens/s" || i == 0 {
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

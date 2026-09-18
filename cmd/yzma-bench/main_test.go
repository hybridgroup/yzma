package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readTestdata(t *testing.T, name string) string {
	t.Helper()

	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}

	return string(b)
}

func TestParseBenchmark(t *testing.T) {
	result, err := parseBenchmark(readTestdata(t, "text-cuda.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if result.arch != "amd64" {
		t.Errorf("arch = %q, want amd64", result.arch)
	}
	if result.cpu != "13th Gen Intel(R) Core(TM) i9-13900HX" {
		t.Errorf("cpu = %q", result.cpu)
	}
	if len(result.runs) != 5 {
		t.Errorf("runs = %d, want 5", len(result.runs))
	}
	if result.tokensPerSecond != 842.5 {
		t.Errorf("tokens/s = %v, want 842.5", result.tokensPerSecond)
	}
}

func TestParseBenchmarkFailedRun(t *testing.T) {
	if _, err := parseBenchmark(readTestdata(t, "failed.txt")); err == nil {
		t.Fatal("a run without PASS must give an error")
	}
}

func TestMedian(t *testing.T) {
	if got := median([]float64{3, 1, 2}); got != 2 {
		t.Errorf("median = %v, want 2", got)
	}
	if got := median([]float64{4, 1, 2, 3}); got != 2.5 {
		t.Errorf("median = %v, want 2.5", got)
	}
}

// newDocument copies the base file, because put and buildTables change it.
func newDocument(t *testing.T) *document {
	t.Helper()

	path := filepath.Join(t.TempDir(), "linux.md")
	if err := os.WriteFile(path, []byte(readTestdata(t, "base.md")), 0o644); err != nil {
		t.Fatal(err)
	}

	doc, err := loadDocument(path)
	if err != nil {
		t.Fatal(err)
	}

	return doc
}

func put(t *testing.T, doc *document, m meta) {
	t.Helper()

	if err := doc.put(m, renderSection(m, "", "", "output")); err != nil {
		t.Fatal(err)
	}
	if err := doc.buildTables(); err != nil {
		t.Fatal(err)
	}
}

func cudaMeta() meta {
	return meta{
		Suite: "text", Backend: "cuda", Arch: "amd64", Machine: "rtx-4070",
		Label: "NVIDIA GeForce RTX 4070", CPU: "Intel i9-13900HX",
		TokensPerSecond: 842.5, LlamaCPP: "b10964", Yzma: "1.27.0", Date: "2026-09-16",
	}
}

func TestPutAddsSectionAndTableRow(t *testing.T) {
	doc := newDocument(t)
	put(t, doc, cudaMeta())

	out := doc.String()
	if !strings.Contains(out, markerStart+"text/cuda/amd64/rtx-4070"+markerClose) {
		t.Error("the section is not there")
	}
	if !strings.Contains(out, "| CUDA | amd64 | NVIDIA GeForce RTX 4070 | - | 842.5 | b10964 | 2026-09-16 |") {
		t.Errorf("the table row is not there:\n%s", out)
	}
	if strings.Contains(out, "| CUDA | amd64 | NVIDIA GeForce RTX 4070 | - | 842.5 | b10964 | 2026-09-16 |\n| CUDA") {
		t.Error("the row is there two times")
	}
}

func TestPutReplacesTheSameKey(t *testing.T) {
	doc := newDocument(t)
	put(t, doc, cudaMeta())

	newer := cudaMeta()
	newer.TokensPerSecond = 900.1
	newer.LlamaCPP = "b11000"
	newer.Date = "2026-10-01"
	put(t, doc, newer)

	out := doc.String()
	if n := strings.Count(out, markerStart+"text/cuda/amd64/rtx-4070"+markerClose); n != 1 {
		t.Errorf("the file has %d sections of the same key, want 1", n)
	}
	if strings.Contains(out, "842.5") {
		t.Error("the old numbers are still there")
	}
	if !strings.Contains(out, "| CUDA | amd64 | NVIDIA GeForce RTX 4070 | - | 900.1 | b11000 | 2026-10-01 |") {
		t.Errorf("the table did not follow the section:\n%s", out)
	}
}

func TestPutKeepsTheOrderOfTheBackends(t *testing.T) {
	doc := newDocument(t)
	put(t, doc, cudaMeta())

	cpu := cudaMeta()
	cpu.Backend = "cpu"
	cpu.Label = "Intel i9-13900HX"
	cpu.TokensPerSecond = 270.5
	put(t, doc, cpu)

	out := doc.String()
	if strings.Index(out, "text/cpu/amd64") > strings.Index(out, "text/cuda/amd64") {
		t.Errorf("cpu must come before cuda:\n%s", out)
	}

	rows := tableRowsOf(out)
	if len(rows) != 2 || !strings.HasPrefix(rows[0], "| CPU |") {
		t.Errorf("the rows are in the wrong order: %v", rows)
	}
}

func TestPutKeepsTheSuitesApart(t *testing.T) {
	doc := newDocument(t)
	put(t, doc, cudaMeta())

	mm := cudaMeta()
	mm.Suite = "multimodal"
	mm.TokensPerSecond = 500
	put(t, doc, mm)

	out := doc.String()
	text := out[strings.Index(out, "## Text"):strings.Index(out, "## Multimodal")]
	if strings.Contains(text, "multimodal/cuda") {
		t.Errorf("the multimodal section is in the text part:\n%s", out)
	}
	if !strings.Contains(out[strings.Index(out, "## Multimodal"):], "multimodal/cuda/amd64/rtx-4070") {
		t.Errorf("the multimodal section is not in its part:\n%s", out)
	}
}

func TestBuildTablesIsStable(t *testing.T) {
	doc := newDocument(t)
	put(t, doc, cudaMeta())

	before := doc.String()
	if err := doc.buildTables(); err != nil {
		t.Fatal(err)
	}
	if doc.String() != before {
		t.Error("buildTables changed a file that is already correct")
	}
}

func TestRemoveDeletesTheSectionAndTheRow(t *testing.T) {
	doc := newDocument(t)
	put(t, doc, cudaMeta())

	cpu := cudaMeta()
	cpu.Backend = "cpu"
	cpu.Label = "Intel i9-13900HX"
	cpu.TokensPerSecond = 270.5
	put(t, doc, cpu)

	if !doc.remove("text/cuda/amd64/rtx-4070") {
		t.Fatal("remove did not find the section")
	}
	if err := doc.buildTables(); err != nil {
		t.Fatal(err)
	}

	out := doc.String()
	if strings.Contains(out, "text/cuda/amd64/rtx-4070") {
		t.Errorf("the section is still there:\n%s", out)
	}
	if strings.Contains(out, "842.5") {
		t.Errorf("the row is still there:\n%s", out)
	}
	if rows := tableRowsOf(out); len(rows) != 1 || !strings.HasPrefix(rows[0], "| CPU |") {
		t.Errorf("the other row went away too: %v", rows)
	}
	if doc.remove("text/cuda/amd64/rtx-4070") {
		t.Error("remove found a section that is not there")
	}
}

func TestParseBenchmarkTakesTheLlamaCppTag(t *testing.T) {
	result, err := parseBenchmark("goarch: wasm\nllama.cpp: b10964\nBenchmarkInference-1\t1\t1 ns/op\t13.8 tokens/s\nPASS\n")
	if err != nil {
		t.Fatal(err)
	}
	if result.llamaCPP != "b10964" {
		t.Errorf("llama.cpp = %q, want b10964", result.llamaCPP)
	}
}

func TestMetaNeedsTheLlamaCppTag(t *testing.T) {
	m := cudaMeta()
	m.LlamaCPP = ""
	if err := m.validate(); err == nil {
		t.Error("a section without the tag of the build must give an error")
	}
}

func TestMetaNeedsEveryPartOfTheKey(t *testing.T) {
	m := cudaMeta()
	m.Machine = ""
	if err := m.validate(); err == nil {
		t.Error("a section without a machine must give an error")
	}

	m = cudaMeta()
	m.Machine = "two words"
	if err := m.validate(); err == nil {
		t.Error("a machine with a space must give an error")
	}
}

func tableRowsOf(out string) []string {
	var rows []string
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "| ") && !strings.HasPrefix(line, "| Backend") &&
			!strings.HasPrefix(line, "| Engine") && !strings.HasPrefix(line, "| ---") {
			rows = append(rows, line)
		}
	}

	return rows
}

func compareMeta() meta {
	return meta{
		Suite: suiteCompareMultimodal, Backend: "yzma", Arch: "amd64", Machine: "rtx-4070",
		Model: "qwen3-vl-4b", Label: "NVIDIA GeForce RTX 4070", CPU: "Intel i9-13900HX",
		TokensPerSecond: 212.4, TTFTMs: 31.2, TotalMs: 144.9, PromptTokens: 706,
		LlamaCPP: "b10964", Yzma: "1.27.0", Date: "2026-09-18",
	}
}

func TestParseBenchmarkTakesEveryMetric(t *testing.T) {
	result, err := parseBenchmark(readTestdata(t, "compare-ollama.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if result.tokensPerSecond != 58.3 {
		t.Errorf("tokens/s = %v, want 58.3", result.tokensPerSecond)
	}
	if result.ttftMs != 93.6 {
		t.Errorf("ttft_ms = %v, want 93.6", result.ttftMs)
	}
	if result.totalMs != 411.6 {
		t.Errorf("total_ms = %v, want 411.6", result.totalMs)
	}
	// The count of the prompt tokens says if the engines do the same work.
	if result.promptTokens != 215 {
		t.Errorf("prompt_tokens = %v, want 215", result.promptTokens)
	}
}

// The suites that came before report tokens/s only, thus the new metrics stay
// empty and the parser must not fail.
func TestParseBenchmarkWithoutTheNewMetrics(t *testing.T) {
	result, err := parseBenchmark(readTestdata(t, "text-cuda.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if result.ttftMs != 0 || result.totalMs != 0 {
		t.Errorf("ttft_ms = %v, total_ms = %v, want 0 and 0", result.ttftMs, result.totalMs)
	}
}

func TestCompareSuiteWritesItsOwnTable(t *testing.T) {
	doc := newDocument(t)
	put(t, doc, compareMeta())

	out := doc.String()
	if !strings.Contains(out, markerStart+"compare-multimodal/yzma/amd64/rtx-4070/qwen3-vl-4b"+markerClose) {
		t.Errorf("the section is not there:\n%s", out)
	}

	want := "| yzma, in process | amd64 | NVIDIA GeForce RTX 4070 | qwen3-vl-4b | 706 | 212.4 | 31.2 | 144.9 | b10964 | 2026-09-18 |"
	if !strings.Contains(out, want) {
		t.Errorf("the table row is not there:\n%s", out)
	}
	if !strings.Contains(out, "| Engine | Arch | Machine | Model | Prompt tokens |") {
		t.Errorf("the comparison table has the wrong header:\n%s", out)
	}
}

func TestCompareSuitePutsTheEnginesOfOneModelTogether(t *testing.T) {
	doc := newDocument(t)

	dmr := compareMeta()
	dmr.Backend = "dmr"
	dmr.EngineVersion = "v0.1.44"
	dmr.LlamaCPP = ""
	dmr.TokensPerSecond = 51.0
	put(t, doc, dmr)

	put(t, doc, compareMeta())

	other := compareMeta()
	other.Model = "gemma4-e4b"
	other.TokensPerSecond = 180.2
	put(t, doc, other)

	rows := tableRowsOf(doc.String())
	if len(rows) != 3 {
		t.Fatalf("the table has %d rows, want 3: %v", len(rows), rows)
	}
	if !strings.Contains(rows[0], "gemma4-e4b") {
		t.Errorf("the models are in the wrong order: %v", rows)
	}
	// yzma comes before dmr inside one model.
	if !strings.Contains(rows[1], "yzma") || !strings.Contains(rows[2], "Docker Model Runner") {
		t.Errorf("the engines are in the wrong order: %v", rows)
	}
}

// A model server has no llama.cpp tag of ours, thus the release of the engine
// takes that place.
func TestCompareSuiteTakesTheReleaseOfTheEngine(t *testing.T) {
	m := compareMeta()
	m.Backend = "ollama"
	m.LlamaCPP = ""
	if err := m.validate(); err == nil {
		t.Error("a comparison section without a version must give an error")
	}

	m.EngineVersion = "0.17.2"
	if err := m.validate(); err != nil {
		t.Errorf("the release of the engine must be enough: %v", err)
	}
	if !strings.Contains(tableRow(m), "0.17.2") {
		t.Errorf("the table does not show the release: %s", tableRow(m))
	}
}

func TestCompareSuiteNeedsAModel(t *testing.T) {
	m := compareMeta()
	m.Model = ""
	if err := m.validate(); err == nil {
		t.Error("a comparison section without a model must give an error")
	}
}

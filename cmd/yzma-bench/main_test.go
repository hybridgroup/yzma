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
		if strings.HasPrefix(line, "| ") && !strings.HasPrefix(line, "| Backend") && !strings.HasPrefix(line, "| ---") {
			rows = append(rows, line)
		}
	}

	return rows
}

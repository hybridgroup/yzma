package main

import (
	"fmt"
	"strings"
)

// backendNames gives the name to show for a backend.
var backendNames = map[string]string{
	"cpu":         "CPU",
	"cpu-threads": "CPU, more threads",
	"cuda":        "CUDA",
	"metal":       "Metal",
	"rocm":        "ROCm",
	"vulkan":      "Vulkan",
	"webgpu":      "WebGPU",
}

func backendName(backend string) string {
	if name, ok := backendNames[backend]; ok {
		return name
	}

	return backend
}

// buildTables makes every table block of the file again from the sections.
func (d *document) buildTables() error {
	for pass := 0; ; pass++ {
		block, ok := d.nextTable()
		if !ok {
			return nil
		}
		if pass > len(d.lines) {
			return fmt.Errorf("%s: the tables do not settle, check the markers", d.path)
		}

		d.replace(block.first+1, block.last, d.tableRows(block.suite))
	}
}

type tableBlock struct {
	suite string
	first int
	last  int
	rows  []string
}

// nextTable gives the first table block whose rows are not correct yet.
func (d *document) nextTable() (tableBlock, bool) {
	var blocks []tableBlock
	var current *tableBlock

	for i, line := range d.lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, markerTableEnd):
			if current != nil && markerValue(trimmed, markerTableEnd) == current.suite {
				current.last = i
				blocks = append(blocks, *current)
				current = nil
			}
		case strings.HasPrefix(trimmed, markerTable):
			current = &tableBlock{suite: markerValue(trimmed, markerTable), first: i}
		case current != nil:
			current.rows = append(current.rows, line)
		}
	}

	for _, block := range blocks {
		if !d.tableIsCurrent(block) {
			return block, true
		}
	}

	return tableBlock{}, false
}

func (d *document) tableIsCurrent(block tableBlock) bool {
	want := d.tableRows(block.suite)

	return strings.Join(block.rows, "\n") == strings.Join(want, "\n")
}

func (d *document) tableRows(suite string) []string {
	rows := []string{
		"| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |",
		"| --- | --- | --- | --- | --- | --- | --- |",
	}
	for _, s := range d.sections() {
		if s.meta.Suite != suite {
			continue
		}
		rows = append(rows, tableRow(s.meta))
	}

	return rows
}

func tableRow(m meta) string {
	label := m.Label
	if label == "" {
		label = m.Machine
	}

	device := m.Device
	if device == "" {
		device = "-"
	}

	return fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |",
		backendName(m.Backend), m.Arch, label, device,
		formatRate(m.TokensPerSecond), orUnknown(m.LlamaCPP), orUnknown(m.Date))
}

func formatRate(value float64) string {
	if value == 0 {
		return "unknown"
	}

	return fmt.Sprintf("%.1f", value)
}

func orUnknown(value string) string {
	if value == "" {
		return "unknown"
	}

	return value
}

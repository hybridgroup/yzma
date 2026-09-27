package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	markerStart    = "<!-- yzma:bench start "
	markerEnd      = "<!-- yzma:bench end "
	markerMeta     = "<!-- yzma:bench meta "
	markerTable    = "<!-- yzma:bench table "
	markerTableEnd = "<!-- yzma:bench table end "
	markerClose    = " -->"
)

// The suites that measure yzma against a model server. Their backend field
// holds the engine, and they report several metrics.
const (
	suiteCompareText       = "compare-text"
	suiteCompareMultimodal = "compare-multimodal"
	suiteCompareEmbeddings = "compare-embeddings"
)

// isCompare reports whether a suite compares engines with each other.
func isCompare(suite string) bool {
	return suite == suiteCompareText || suite == suiteCompareMultimodal ||
		suite == suiteCompareEmbeddings
}

// meta is the machine readable part of a section. The tables come from it.
type meta struct {
	Suite           string  `json:"suite"`
	Backend         string  `json:"backend"`
	Arch            string  `json:"arch"`
	Machine         string  `json:"machine"`
	Device          string  `json:"device,omitempty"`
	Model           string  `json:"model,omitempty"`
	Label           string  `json:"label,omitempty"`
	CPU             string  `json:"cpu,omitempty"`
	TokensPerSecond float64 `json:"tokens_per_second"`
	TTFTMs          float64 `json:"ttft_ms,omitempty"`
	TotalMs         float64 `json:"total_ms,omitempty"`
	PromptTokens    float64 `json:"prompt_tokens,omitempty"`
	LlamaCPP        string  `json:"llamacpp,omitempty"`
	EngineVersion   string  `json:"engine_version,omitempty"`
	Yzma            string  `json:"yzma,omitempty"`
	Date            string  `json:"date,omitempty"`
}

// key names the section. It includes the device, because one machine can have
// several devices on the same backend, and the model, because the comparison
// suite runs several models on one machine.
func (m meta) key() string {
	parts := []string{m.Suite, m.Backend, m.Arch, m.Machine}
	if m.Device != "" {
		parts = append(parts, strings.ToLower(m.Device))
	}
	if m.Model != "" {
		parts = append(parts, strings.ToLower(m.Model))
	}

	return strings.Join(parts, "/")
}

// version returns the build or release that produced the result.
func (m meta) version() string {
	if isCompare(m.Suite) && m.EngineVersion != "" {
		return m.EngineVersion
	}

	return m.LlamaCPP
}

func (m meta) validate() error {
	for name, value := range map[string]string{
		"suite": m.Suite, "backend": m.Backend, "arch": m.Arch, "machine": m.Machine,
	} {
		if value == "" {
			return fmt.Errorf("the section needs a %s", name)
		}
		if strings.ContainsAny(value, "/ ") {
			return fmt.Errorf("the %s %q must have no space and no slash", name, value)
		}
	}
	if strings.ContainsAny(m.Device, "/ ") {
		return fmt.Errorf("the device %q must have no space and no slash", m.Device)
	}
	if strings.ContainsAny(m.Model, "/ ") {
		return fmt.Errorf("the model %q must have no space and no slash", m.Model)
	}
	// Each result must record which build produced it, or the table can't
	// compare it with the others. A model server has no llama.cpp tag of ours,
	// so the comparison suite uses the engine release instead.
	if isCompare(m.Suite) {
		if m.Model == "" {
			return fmt.Errorf("the comparison suite needs a model, pass --model")
		}
		if m.version() == "" {
			return fmt.Errorf("the section needs the engine release, pass --engine-version")
		}

		return nil
	}
	if m.LlamaCPP == "" {
		return fmt.Errorf("the section needs the llama.cpp build tag, pass --llamacpp")
	}

	return nil
}

// backendOrder is the order of the backends in a file. A backend not listed
// here goes last, in alphabetical order.
var backendOrder = []string{
	"cpu", "cpu-threads", "metal", "cuda", "rocm", "vulkan", "webgpu",
	"yzma", "ollama", "dmr",
}

func (m meta) rank() (int, string, string) {
	i := slices.Index(backendOrder, m.Backend)
	if i < 0 {
		i = len(backendOrder)
	}

	// The comparison suite groups the engines for one model together, because
	// the table compares them with each other.
	if isCompare(m.Suite) {
		return i, m.Model + "/" + m.Machine, m.Arch
	}

	return i, m.Backend + "/" + m.Arch, m.Machine + "/" + m.Device + "/" + m.Model
}

func (m meta) less(other meta) bool {
	ai, ab, am := m.rank()
	bi, bb, bm := other.rank()

	// The comparison suite sorts by model first, so each model gets one group
	// of rows with the engines lined up inside it.
	if isCompare(m.Suite) && isCompare(other.Suite) {
		switch {
		case ab != bb:
			return ab < bb
		case ai != bi:
			return ai < bi
		default:
			return am < bm
		}
	}

	switch {
	case ai != bi:
		return ai < bi
	case ab != bb:
		return ab < bb
	default:
		return am < bm
	}
}

type section struct {
	meta  meta
	key   string
	first int // index of the start marker
	last  int // index of the end marker
}

type document struct {
	path  string
	lines []string
}

func loadDocument(path string) (*document, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &document{path: path, lines: strings.Split(string(b), "\n")}, nil
}

func (d *document) String() string {
	return strings.Join(d.lines, "\n")
}

func (d *document) save() error {
	return os.WriteFile(d.path, []byte(d.String()), 0o644)
}

// sections returns every marked section in file order.
func (d *document) sections() []section {
	var found []section
	var current *section

	for i, line := range d.lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, markerStart):
			key := markerValue(trimmed, markerStart)
			current = &section{key: key, first: i, last: -1}
		case current != nil && strings.HasPrefix(trimmed, markerMeta):
			raw := markerValue(trimmed, markerMeta)
			_ = json.Unmarshal([]byte(raw), &current.meta)
		case current != nil && strings.HasPrefix(trimmed, markerEnd):
			current.last = i
			found = append(found, *current)
			current = nil
		}
	}

	return found
}

func markerValue(line, prefix string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, prefix), markerClose))
}

// put replaces the section with the same key, or adds it in the right place
// within its suite.
func (d *document) put(m meta, body string) error {
	lines := strings.Split(body, "\n")

	for _, s := range d.sections() {
		if s.key == m.key() {
			d.replace(s.first, s.last+1, lines)
			return nil
		}
	}

	at, err := d.insertPoint(m)
	if err != nil {
		return err
	}

	d.replace(at, at, append(lines, ""))

	return nil
}

// remove deletes the section for a key and reports whether it found one.
func (d *document) remove(key string) bool {
	for _, s := range d.sections() {
		if s.key != key {
			continue
		}

		last := s.last + 1
		// A section has one empty line after it, which is removed with it.
		if last < len(d.lines) && strings.TrimSpace(d.lines[last]) == "" {
			last++
		}
		d.replace(s.first, last, nil)

		return true
	}

	return false
}

// insertPoint returns the line where a new section for the suite goes.
func (d *document) insertPoint(m meta) (int, error) {
	last := -1
	for _, s := range d.sections() {
		if s.meta.Suite != m.Suite {
			continue
		}
		if m.less(s.meta) {
			return s.first, nil
		}
		last = s.last
	}

	if last >= 0 {
		return min(last+2, len(d.lines)), nil
	}

	// The suite has no section yet, so the new one goes after its table.
	for i, line := range d.lines {
		if strings.TrimSpace(line) == markerTableEnd+m.Suite+markerClose {
			return min(i+2, len(d.lines)), nil
		}
	}

	return 0, fmt.Errorf("%s has no table block for the suite %q, add %s%s%s to it",
		d.path, m.Suite, markerTable, m.Suite, markerClose)
}

func (d *document) replace(from, to int, lines []string) {
	updated := slices.Clone(d.lines[:from])
	updated = append(updated, lines...)
	updated = append(updated, d.lines[to:]...)
	d.lines = updated
}

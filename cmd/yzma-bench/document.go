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

// meta is the machine readable part of a section. The tables come from it.
type meta struct {
	Suite           string  `json:"suite"`
	Backend         string  `json:"backend"`
	Arch            string  `json:"arch"`
	Machine         string  `json:"machine"`
	Device          string  `json:"device,omitempty"`
	Label           string  `json:"label,omitempty"`
	CPU             string  `json:"cpu,omitempty"`
	TokensPerSecond float64 `json:"tokens_per_second"`
	LlamaCPP        string  `json:"llamacpp,omitempty"`
	Yzma            string  `json:"yzma,omitempty"`
	Date            string  `json:"date,omitempty"`
}

// key names the section. The device is part of it, because one machine can
// have more than one device of the same backend.
func (m meta) key() string {
	parts := []string{m.Suite, m.Backend, m.Arch, m.Machine}
	if m.Device != "" {
		parts = append(parts, strings.ToLower(m.Device))
	}

	return strings.Join(parts, "/")
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
	// Each result must say which build made it, or the table cannot compare it
	// with the others.
	if m.LlamaCPP == "" {
		return fmt.Errorf("the section needs the tag of the llama.cpp build, give --llamacpp")
	}

	return nil
}

// backendOrder is the order of the backends in a file. A backend that is not
// here goes last, in alphabetical order.
var backendOrder = []string{"cpu", "cpu-threads", "metal", "cuda", "rocm", "vulkan", "webgpu"}

func (m meta) rank() (int, string, string) {
	i := slices.Index(backendOrder, m.Backend)
	if i < 0 {
		i = len(backendOrder)
	}

	return i, m.Backend + "/" + m.Arch, m.Machine + "/" + m.Device
}

func (m meta) less(other meta) bool {
	ai, ab, am := m.rank()
	bi, bb, bm := other.rank()
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

// sections gives every marked section of the file, in the order of the file.
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

// put replaces the section that has the same key, or adds it in the correct
// place of its suite.
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

// remove deletes the section of a key and says if it found one.
func (d *document) remove(key string) bool {
	for _, s := range d.sections() {
		if s.key != key {
			continue
		}

		last := s.last + 1
		// A section has one empty line after it, which goes away with it.
		if last < len(d.lines) && strings.TrimSpace(d.lines[last]) == "" {
			last++
		}
		d.replace(s.first, last, nil)

		return true
	}

	return false
}

// insertPoint gives the line where a new section of the suite goes.
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

	// The suite has no section yet, thus the new one goes after its table.
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

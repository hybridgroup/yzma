package decide

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

// toJSON writes a JSON value as the tojson filter of llama.cpp does, with
// ", " and ": " separators and floats as C++ streams print them.
func toJSON(data []byte, sortKeys bool) (string, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	v, err := readJSON(dec)
	if err != nil {
		return "", err
	}
	if _, err := dec.Token(); err != io.EOF {
		return "", errors.New("decide: trailing data after JSON value")
	}

	var b strings.Builder
	writeJSON(&b, v, sortKeys)
	return b.String(), nil
}

// stateText passes a string through and writes any other value with [toJSON].
func stateText(state any, sortKeys bool) (string, error) {
	if s, ok := state.(string); ok {
		return s, nil
	}

	data, ok := state.(json.RawMessage)
	if !ok {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(state); err != nil {
			return "", fmt.Errorf("state: %w", err)
		}
		data = buf.Bytes()
	}
	return toJSON(data, sortKeys)
}

func writeJSON(b *strings.Builder, n *jsonNode, sortKeys bool) {
	switch n.kind {
	case '[':
		b.WriteByte('[')
		for i, v := range n.vals {
			if i > 0 {
				b.WriteString(", ")
			}
			writeJSON(b, v, sortKeys)
		}
		b.WriteByte(']')
	case '{':
		order := make([]int, len(n.keys))
		for i := range order {
			order[i] = i
		}
		if sortKeys {
			slices.SortStableFunc(order, func(x, y int) int { return strings.Compare(n.keys[x], n.keys[y]) })
		}
		b.WriteByte('{')
		for k, i := range order {
			if k > 0 {
				b.WriteString(", ")
			}
			writeJSONString(b, n.keys[i])
			b.WriteString(": ")
			writeJSON(b, n.vals[i], sortKeys)
		}
		b.WriteByte('}')
	default:
		var s string
		if len(n.raw) > 0 && n.raw[0] == '"' && json.Unmarshal(n.raw, &s) == nil {
			writeJSONString(b, s)
			return
		}
		b.WriteString(jsonNumber(string(n.raw)))
	}
}

// jsonNumber keeps an integer, true, false and null as written and prints a
// float with 6 significant digits.
func jsonNumber(s string) string {
	if !strings.ContainsAny(s, ".eE") || s == "true" || s == "false" {
		return s
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	return strconv.FormatFloat(f, 'g', 6, 64)
}

func writeJSONString(b *strings.Builder, s string) {
	b.WriteByte('"')
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if c < 0x20 {
				fmt.Fprintf(b, `\u%04x`, c)
			} else {
				b.WriteByte(c)
			}
		}
	}
	b.WriteByte('"')
}

package decide

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
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
	data, err := stateJSON(state)
	if err != nil {
		return "", err
	}
	return toJSON(data, sortKeys)
}

// stateJSON returns a state that is not a string as JSON.
func stateJSON(state any) ([]byte, error) {
	if data, ok := state.(json.RawMessage); ok {
		return data, nil
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(state); err != nil {
		return nil, fmt.Errorf("state: %w", err)
	}
	return buf.Bytes(), nil
}

// kevText passes a string through and writes any other value as text, as the
// Kev training code does, with the keys of an object kept as labels.
func kevText(state any) (string, error) {
	if s, ok := state.(string); ok {
		return s, nil
	}
	data, err := stateJSON(state)
	if err != nil {
		return "", err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	n, err := readJSON(dec)
	if err != nil {
		return "", fmt.Errorf("state: %w", err)
	}
	return kevRender(n, 0), nil
}

func kevRender(n *jsonNode, indent int) string {
	pad := strings.Repeat("  ", indent)
	var b strings.Builder
	switch n.kind {
	case '[':
		for i, v := range n.vals {
			if i > 0 {
				b.WriteByte('\n')
			}
			b.WriteString(pad + "- " + strings.TrimLeft(kevRender(v, indent+1), " \t\n\r"))
		}
	case '{':
		for i, k := range n.keys {
			if i > 0 {
				b.WriteByte('\n')
			}
			if n.vals[i].kind != 0 {
				b.WriteString(pad + k + ":\n" + kevRender(n.vals[i], indent+1))
			} else {
				b.WriteString(pad + k + ": " + kevRender(n.vals[i], 0))
			}
		}
	default:
		var s string
		switch raw := string(n.raw); {
		case raw == "null":
		case raw == "true":
			b.WriteString("True")
		case raw == "false":
			b.WriteString("False")
		case strings.HasPrefix(raw, `"`) && json.Unmarshal(n.raw, &s) == nil:
			b.WriteString(s)
		default:
			b.WriteString(dumpNumber(raw))
		}
	}
	return b.String()
}

// dumpNumber writes a number as nlohmann::json dump does. A float always has
// a dot or an exponent, and an exponent is used outside 1e-4 to 1e15.
func dumpNumber(s string) string {
	if !strings.ContainsAny(s, ".eE") {
		return s
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	sign := ""
	if math.Signbit(f) {
		sign, f = "-", -f
	}
	if f == 0 {
		return sign + "0.0"
	}

	e := strconv.FormatFloat(f, 'e', -1, 64)
	mant, exp, _ := strings.Cut(e, "e")
	digits := strings.Replace(mant, ".", "", 1)
	x, _ := strconv.Atoi(exp)
	k, p := len(digits), x+1
	switch {
	case k <= p && p <= 15:
		return sign + digits + strings.Repeat("0", p-k) + ".0"
	case 0 < p && p <= 15:
		return sign + digits[:p] + "." + digits[p:]
	case -4 < p && p <= 0:
		return sign + "0." + strings.Repeat("0", -p) + digits
	}
	m := digits[:1]
	if k > 1 {
		m += "." + digits[1:]
	}
	es := "+"
	if x < 0 {
		es, x = "-", -x
	}
	return fmt.Sprintf("%s%se%s%02d", sign, m, es, x)
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

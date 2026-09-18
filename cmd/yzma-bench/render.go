package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// renderSection makes the markdown of one result, between the markers that let
// an update find it again.
func renderSection(m meta, notes, device, output string) string {
	encoded, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s%s%s\n", markerStart, m.key(), markerClose)
	heading := fmt.Sprintf("### %s, %s, %s", backendName(m.Backend), m.Arch, m.Label)
	if m.Device != "" {
		heading += ", " + m.Device
	}
	if m.Model != "" {
		heading += ", " + m.Model
	}
	fmt.Fprintf(&b, "%s\n", heading)
	fmt.Fprintf(&b, "%s%s%s\n\n", markerMeta, encoded, markerClose)

	if m.CPU != "" {
		fmt.Fprintf(&b, "%s. %s tokens a second.", m.CPU, formatRate(m.TokensPerSecond))
		if m.TTFTMs != 0 {
			fmt.Fprintf(&b, " %s ms to the first token.", formatRate(m.TTFTMs))
		}
		if m.TotalMs != 0 {
			fmt.Fprintf(&b, " %s ms for a request.", formatRate(m.TotalMs))
		}
		if m.PromptTokens != 0 {
			fmt.Fprintf(&b, " %s prompt tokens.", formatCount(m.PromptTokens))
		}
		b.WriteString("\n\n")
	}
	if notes != "" {
		fmt.Fprintf(&b, "%s\n\n", notes)
	}

	if device != "" {
		b.WriteString("<details><summary>The device</summary>\n\n")
		fmt.Fprintf(&b, "```\n%s\n```\n\n", strings.TrimRight(device, "\n"))
		b.WriteString("</details>\n\n")
	}

	b.WriteString("<details><summary>The output of go test</summary>\n\n")
	fmt.Fprintf(&b, "```\n%s\n```\n\n", strings.TrimRight(output, "\n"))
	b.WriteString("</details>\n")

	fmt.Fprintf(&b, "%s%s%s", markerEnd, m.key(), markerClose)

	return b.String()
}

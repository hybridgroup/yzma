package decide

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

const (
	configFormat   = "macjev-readout-v1"
	configTemplate = "macjev-render-v1"
	configReadout  = "verdict"
)

// SlotToken is a token the readout depends on.
type SlotToken struct {
	Text string `json:"text"`
	ID   int32  `json:"id"`
}

// Group is a fitted calibration temperature.
type Group struct {
	T float64 `json:"T"`
}

// Config is the readout_config.json that ships with a Jev-style model.
type Config struct {
	Format   string `json:"format"`
	Template string `json:"template"`
	Readout  string `json:"readout"`

	SlotTokens struct {
		Yes         SlotToken `json:"yes"`
		No          SlotToken `json:"no"`
		VerdictSlot SlotToken `json:"verdict_slot"`
	} `json:"slot_tokens"`

	// Families maps a category pattern such as "general_*" to a calibration family.
	Families map[string]string `json:"families"`

	Budgets struct {
		MaxLen  int `json:"max_len"`
		HeadMax int `json:"head_max"`
	} `json:"budgets"`

	Temperatures struct {
		Global float64          `json:"global"`
		Clamp  []float64        `json:"clamp"`
		Groups map[string]Group `json:"groups"`
	} `json:"temperatures"`

	prefixes []familyPrefix
}

type familyPrefix struct {
	prefix, family string
}

// LoadConfig reads and checks a readout_config.json file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseConfig(data)
}

// ParseConfig parses and checks the contents of a readout_config.json file.
func ParseConfig(data []byte) (*Config, error) {
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("readout config: %w", err)
	}

	if c.Format != configFormat || c.Template != configTemplate {
		return nil, fmt.Errorf("readout config is %q / %q, want %q / %q", c.Format, c.Template, configFormat, configTemplate)
	}
	if c.Readout != configReadout {
		return nil, fmt.Errorf("readout config uses readout %q, only %q is supported", c.Readout, configReadout)
	}
	if c.SlotTokens.Yes.Text == "" || c.SlotTokens.No.Text == "" || c.SlotTokens.VerdictSlot.Text == "" {
		return nil, errors.New("readout config is missing slot tokens")
	}
	if c.Budgets.MaxLen <= 0 || c.Budgets.HeadMax <= 0 || c.Budgets.HeadMax > c.Budgets.MaxLen {
		return nil, fmt.Errorf("readout config has bad budgets max_len %d head_max %d", c.Budgets.MaxLen, c.Budgets.HeadMax)
	}
	if c.Temperatures.Global <= 0 {
		return nil, errors.New("readout config has no global temperature")
	}
	if len(c.Temperatures.Clamp) != 2 || c.Temperatures.Clamp[0] <= 0 || c.Temperatures.Clamp[0] > c.Temperatures.Clamp[1] {
		return nil, fmt.Errorf("readout config has bad temperature clamp %v", c.Temperatures.Clamp)
	}

	for pattern, fam := range c.Families {
		if p, ok := strings.CutSuffix(pattern, "*"); ok && p != "" {
			c.prefixes = append(c.prefixes, familyPrefix{p, fam})
		}
	}
	// Longest prefix first, so a more specific pattern wins.
	sort.Slice(c.prefixes, func(i, j int) bool {
		if len(c.prefixes[i].prefix) != len(c.prefixes[j].prefix) {
			return len(c.prefixes[i].prefix) > len(c.prefixes[j].prefix)
		}
		return c.prefixes[i].prefix < c.prefixes[j].prefix
	})

	return &c, nil
}

// Family returns the calibration family of a category, or "other".
func (c *Config) Family(category string) string {
	for _, p := range c.prefixes {
		if strings.HasPrefix(category, p.prefix) {
			return p.family
		}
	}
	return "other"
}

// Temperature returns the calibration temperature for a question.
// An empty category uses the global temperature.
func (c *Config) Temperature(category string, t Type, nOptions int) float64 {
	temp := c.Temperatures.Global
	if category != "" {
		key := c.Family(category) + "|" + string(t) + "|" + optionBucket(nOptions)
		if g, ok := c.Temperatures.Groups[key]; ok {
			temp = g.T
		}
	}
	return math.Min(c.Temperatures.Clamp[1], math.Max(c.Temperatures.Clamp[0], temp))
}

func optionBucket(k int) string {
	switch {
	case k <= 2:
		return "2"
	case k <= 5:
		return "3-5"
	case k <= 10:
		return "6-10"
	case k <= 20:
		return "11-20"
	}
	return "21+"
}

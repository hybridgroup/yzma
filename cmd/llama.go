package cmd

import (
	"fmt"

	"github.com/hybridgroup/yzma/pkg/download"
	"github.com/urfave/cli/v2"
)

var LlamaCmd = &cli.Command{
	Name:  "llama",
	Usage: "Show most recent llama.cpp version",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "short",
			Aliases: []string{"s"},
			Usage:   "suppress additional output",
			Value:   false,
		},
	},
	Action: func(c *cli.Context) error {
		return showLlama(c)
	},
}

func showLlama(c *cli.Context) error {
	version, err := download.LlamaLatestVersion()
	if err != nil {
		return fmt.Errorf("could not obtain latest version: %w", err)
	}

	if c.Bool("short") {
		fmt.Println(version)
		return nil
	}

	fmt.Printf("llama.cpp most recent version is %s\n", version)
	return nil
}

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "yzma",
		Usage:    "YZMA command line tool",
		Commands: buildCommands(),
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func buildCommands() []*cli.Command {
	return []*cli.Command{
		installCmd,
		systemCmd,
		llamaCmd,
		modelCmd,
		versionCmd,
		infoCmd,
	}
}

var versionCmd = &cli.Command{
	Name:  "version",
	Usage: "Show yzma version",
	Action: func(c *cli.Context) error {
		return runShowVersion(c)
	},
}

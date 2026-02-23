package main

import (
	"fmt"
	"os"

	"github.com/hybridgroup/yzma/cmd"
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
		cmd.InstallCmd,
		cmd.SystemCmd,
		cmd.LlamaCmd,
		cmd.ModelCmd,
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

func runShowVersion(c *cli.Context) error {
	return showYzmaVersion()
}

func showYzmaVersion() error {
	fmt.Printf("yzma version %s\n", Version())
	return nil
}

func runShowInfo(c *cli.Context) error {
	cmd.ShowInfo(c)
	return showYzmaVersion()
}

var infoCmd = &cli.Command{
	Name:  "info",
	Usage: "Show yzma version",
	Action: func(c *cli.Context) error {
		return runShowInfo(c)
	},
}

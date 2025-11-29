package main

import (
	"fmt"

	"github.com/hybridgroup/yzma"
	"github.com/urfave/cli/v2"
)

func runShowVersion(c *cli.Context) error {
	return showYzmaVersion()
}

func showYzmaVersion() error {
	fmt.Printf("yzma version %s\n", yzma.Version())
	return nil
}

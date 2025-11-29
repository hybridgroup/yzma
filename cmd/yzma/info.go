package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var infoCmd = &cli.Command{
	Name:  "info",
	Usage: "Show yzma version",
	Action: func(c *cli.Context) error {
		return showInfo(c)
	},
}

func showInfo(c *cli.Context) error {
	fmt.Println(logo)
	fmt.Println("Go local inference based on llama.cpp")
	showYzmaVersion()

	return nil
}

const logo = `
 __ __  _____  ___ ___   ____ 
|  |  ||     ||   |   | /    |
|  |  ||__/  || _   _ ||  o  |
|  ~  ||   __||  \_/  ||     |
|___, ||  /  ||   |   ||  _  |
|     ||     ||   |   ||  |  |
|____/ |_____||___|___||__|__|
`

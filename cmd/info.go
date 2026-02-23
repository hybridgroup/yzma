package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func ShowInfo(c *cli.Context) error {
	fmt.Println(logo)
	fmt.Println()
	fmt.Println("Local inference in Go using llama.cpp with hardware acceleration")

	return nil
}

const logo = `
 __ __  _____  ___ ___   ____ 
|  |  ||     ||   |   | /    |
|  |  ||__/  || _   _ ||  o  |
|  ~  ||   __||  \_/  ||     |
|___, ||  /  ||   |   ||  _  |
|     ||     ||   |   ||  |  |
|____/ |_____||___|___||__|__|`

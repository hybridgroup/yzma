package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hybridgroup/yzma/pkg/download"
	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/urfave/cli/v2"
)

var ModelCmd = &cli.Command{
	Name:  "model",
	Usage: "Manage models",
	Subcommands: []*cli.Command{
		modelInfoCmd,
		modelGetCmd,
	},
}

var modelInfoCmd = &cli.Command{
	Name:  "info",
	Usage: "Show information about a model",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "model",
			Aliases:  []string{"m"},
			Usage:    "path to the model file",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "lib",
			Aliases: []string{"l"},
			Usage:   "path to llama.cpp compiled library files",
			EnvVars: []string{"YZMA_LIB"},
		},
	},
	Action: func(c *cli.Context) error {
		return runModelInfo(c)
	},
}

func runModelInfo(c *cli.Context) error {
	return showModelInfo(c)
}

func showModelInfo(c *cli.Context) error {
	llama.Load(c.String("lib"))
	llama.LogSet(llama.LogSilent())

	llama.Init()

	model, err := llama.ModelLoadFromFile(c.String("model"), llama.ModelDefaultParams())
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to load model from file %s: %v\n", c.String("model"), err)
		os.Exit(1)
	}

	defer llama.ModelFree(model)

	desc := llama.ModelDesc(model)
	fmt.Printf("Model Description: %s\n", desc)

	size := llama.ModelSize(model)
	fmt.Printf("Model Size: %d tensors\n", size)

	encoder := llama.ModelHasEncoder(model)
	fmt.Printf("Model Has Encoder: %v\n", encoder)

	decoder := llama.ModelHasDecoder(model)
	fmt.Printf("Model Has Decoder: %v\n", decoder)

	recurrent := llama.ModelIsRecurrent(model)
	fmt.Printf("Model Is Recurrent: %v\n", recurrent)

	hybrid := llama.ModelIsHybrid(model)
	fmt.Printf("Model Is Hybrid: %v\n", hybrid)

	count := llama.ModelMetaCount(model)
	fmt.Printf("Model Metadata (%d entries):\n", count)
	for i := int32(0); i < count; i++ {
		key, ok := llama.ModelMetaKeyByIndex(model, i)
		if !ok {
			fmt.Printf("Error getting key for index %d\n", i)
			continue
		}
		value, ok := llama.ModelMetaValStrByIndex(model, i)
		if !ok {
			fmt.Printf("Error getting value for index %d\n", i)
			continue
		}
		fmt.Printf("  %s: %s\n", key, value)
	}

	return nil
}

var modelGetCmd = &cli.Command{
	Name:  "get",
	Usage: "Download a model from a URL",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "url",
			Aliases:  []string{"u"},
			Usage:    "URL of the model to download",
			Required: true,
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       "Path to save the downloaded model",
			Value:       download.DefaultModelsDir(),
			DefaultText: "~/models",
		},
		&cli.BoolFlag{
			Name:    "yes",
			Aliases: []string{"y"},
			Usage:   "Automatically answer yes to prompts",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:  "show-progress",
			Usage: "Show download progress",
			Value: true,
		},
	},
	Action: func(c *cli.Context) error {
		return runModelDownload(c)
	},
}

func runModelDownload(c *cli.Context) error {
	url := c.String("url")
	output := c.String("output")
	autoYes := c.Bool("yes")
	showProgress := c.Bool("show-progress")

	// Check if the output directory exists
	if _, err := os.Stat(output); os.IsNotExist(err) {
		if !autoYes {
			fmt.Printf("Directory %s does not exist.\n", output)
			fmt.Print("Would you like to create it? [y/N]: ")

			var response string
			fmt.Scanln(&response)

			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Download cancelled.")
				return nil
			}
		}

		if err := os.MkdirAll(output, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error creating output directory: %v\n", err)
			return err
		}
		fmt.Printf("Created directory %s\n", output)
	}

	fmt.Printf("Downloading model from %s to %s...\n", url, output)

	if !showProgress {
		download.ProgressTracker = nil
	}

	if err := download.GetModel(url, output); err != nil {
		fmt.Fprintf(os.Stderr, "error downloading model: %v\n", err)
		return err
	}

	fmt.Println("Download completed successfully.")
	return nil
}

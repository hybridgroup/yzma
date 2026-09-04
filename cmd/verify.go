package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/hybridgroup/yzma/pkg/download"
	"github.com/urfave/cli/v2"
)

var VerifyCmd = &cli.Command{
	Name:  "verify",
	Usage: "Check the installed llama.cpp libraries against their published digests",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "lib",
			Aliases: []string{"l"},
			Usage:   "path to llama.cpp compiled library files",
			EnvVars: []string{"YZMA_LIB"},
		},
		&cli.StringFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "the llama.cpp version that must be installed, optionally as VERSION@sha256:DIGEST to pin the digests (leave empty to take the installed one)",
			Value:   "",
		},
		&cli.BoolFlag{
			Name:  "strict",
			Usage: "also fail when a file in the directory is not part of the install",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "json",
			Usage: "write the report as JSON",
			Value: false,
		},
	},
	Action: func(c *cli.Context) error {
		return runVerify(c)
	},
}

func runVerify(c *cli.Context) error {
	libPath := c.String("lib")
	if libPath == "" {
		return fmt.Errorf("missing lib flag or YZMA_LIB env var")
	}

	report, err := download.VerifyInstall(context.Background(), libPath, c.String("version"))
	switch {
	case errors.Is(err, download.ErrNoInstallRecord):
		return fmt.Errorf("%w. Install with a yzma release that writes one, then verify", err)
	case errors.Is(err, download.ErrNoFileDigests):
		return fmt.Errorf("%w. The publisher gives no file digests for this release", err)
	case err != nil:
		return err
	}

	if c.Bool("json") {
		body, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(body))
	} else {
		showVerifyReport(report)
	}

	if !report.OK() || (c.Bool("strict") && report.Unexpected > 0) {
		return cli.Exit("", 1)
	}

	return nil
}

// showVerifyReport writes the report for a person to read. Only the files that need
// attention are named, because an install holds many files.
func showVerifyReport(report *download.VerifyReport) {
	out := os.Stdout

	fmt.Fprintf(out, "llama.cpp %s in %s\n", report.Tag, report.LibPath)

	for _, file := range report.Files {
		if file.State == download.FileVerified {
			continue
		}
		fmt.Fprintf(out, "  %-10s %s\n", file.State, file.Name)
	}

	fmt.Fprintf(out, "%d verified, %d changed, %d missing, %d not part of this install\n",
		report.Verified, report.Changed, report.Missing, report.Unexpected)

	if report.OK() {
		fmt.Fprintln(out, "ok.")
	}
}

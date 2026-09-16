// yzma-bench puts the result of a benchmark run into the markdown file of the
// platform. See benchmarks/README.md.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "update":
		err = runUpdate(os.Args[2:])
	case "check":
		err = runCheck(os.Args[2:])
	case "-h", "--help", "help":
		usage()
		return
	default:
		usage()
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "yzma-bench:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `usage:
  yzma-bench update --file FILE --suite SUITE --backend BACKEND --machine NAME [flags]
  yzma-bench check FILE...

update puts one result in the file of the platform and makes the tables again.
check verifies that the tables agree with the sections.
`)
}

func runUpdate(args []string) error {
	fs := flag.NewFlagSet("update", flag.ExitOnError)
	file := fs.String("file", "", "markdown file of the platform")
	suite := fs.String("suite", "", "suite of the benchmark (text, multimodal)")
	backend := fs.String("backend", "", "backend (cpu, cuda, rocm, vulkan, metal, webgpu)")
	arch := fs.String("arch", "", "architecture, empty to take goarch from the output")
	machine := fs.String("machine", "", "short name of the machine, used in the key")
	device := fs.String("device", "", "device of the run, as go test takes it (CUDA0, VULKAN1)")
	label := fs.String("label", "", "name of the machine to show, empty to take the short name")
	llamacpp := fs.String("llamacpp", "", "tag of the llama.cpp build")
	yzma := fs.String("yzma", "", "version of yzma")
	date := fs.String("date", "", "date of the run, today when the flag is not there")
	output := fs.String("output", "-", "file with the output of go test, - for stdin")
	deviceInfo := fs.String("device-info", "", "file with the information of the device")
	notes := fs.String("notes", "", "one line of text to put above the output")
	dryRun := fs.Bool("dry-run", false, "print the result and change no file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	switch {
	case *file == "":
		return fmt.Errorf("--file is needed")
	case *suite == "":
		return fmt.Errorf("--suite is needed")
	case *backend == "":
		return fmt.Errorf("--backend is needed")
	case *machine == "":
		return fmt.Errorf("--machine is needed")
	}

	// A date that the flag does not give is today. The old results that came
	// from BENCHMARKS.md give "", which the table shows as unknown.
	if !isSet(fs, "date") {
		*date = time.Now().Format("2006-01-02")
	}

	raw, err := readInput(*output)
	if err != nil {
		return err
	}

	result, err := parseBenchmark(raw)
	if err != nil {
		return err
	}

	deviceText := ""
	if *deviceInfo != "" {
		b, err := os.ReadFile(*deviceInfo)
		if err != nil {
			return err
		}
		deviceText = string(b)
	}

	m := meta{
		Suite:           *suite,
		Backend:         *backend,
		Arch:            firstOf(*arch, result.arch),
		Machine:         *machine,
		Device:          *device,
		Label:           firstOf(*label, *machine),
		CPU:             result.cpu,
		TokensPerSecond: result.tokensPerSecond,
		LlamaCPP:        *llamacpp,
		Yzma:            *yzma,
		Date:            *date,
	}
	if err := m.validate(); err != nil {
		return err
	}

	doc, err := loadDocument(*file)
	if err != nil {
		return err
	}

	body := renderSection(m, *notes, deviceText, raw)
	if err := doc.put(m, body); err != nil {
		return err
	}
	if err := doc.buildTables(); err != nil {
		return err
	}

	if *dryRun {
		fmt.Print(doc.String())
		return nil
	}

	return doc.save()
}

func runCheck(args []string) error {
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() == 0 {
		return fmt.Errorf("give one file or more")
	}

	bad := 0
	for _, name := range fs.Args() {
		doc, err := loadDocument(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
			bad++
			continue
		}

		before := doc.String()
		if err := doc.buildTables(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
			bad++
			continue
		}
		if doc.String() != before {
			fmt.Fprintf(os.Stderr, "%s: the tables do not agree with the sections\n", name)
			bad++
			continue
		}
		fmt.Printf("%s: ok, %d sections\n", name, len(doc.sections()))
	}

	if bad > 0 {
		return fmt.Errorf("%d file(s) need an update, run yzma-bench update or make the tables again", bad)
	}

	return nil
}

func isSet(fs *flag.FlagSet, name string) bool {
	found := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})

	return found
}

func readInput(name string) (string, error) {
	if name == "-" {
		b, err := io.ReadAll(os.Stdin)
		return string(b), err
	}

	b, err := os.ReadFile(name)
	return string(b), err
}

func firstOf(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}

	return ""
}

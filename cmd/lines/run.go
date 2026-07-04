package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/moderrek/lines/pkg/lines"
)

func run(stdout, stderr io.Writer, args []string) error {
	opts, fs, err := parseFlags(stderr, args)
	if err != nil {
		return err
	}

	isTerminal := isatty.IsTerminal(os.Stdout.Fd())
	useColor := (isTerminal || opts.color) && !opts.noColor
	color.NoColor = !useColor

	if opts.version {
		fmt.Printf("%s version %s created by %s\n", PROGRAM_NAME, VERSION, AUTHOR)
		return nil
	}

	if opts.help {
		fmt.Fprintf(stderr, "Usage: %s [options]\n", PROGRAM_NAME)
		fs.PrintDefaults()
		return nil
	}

	startTime := time.Now()
	if isTerminal && !opts.json {
		fmt.Fprintf(stderr, "Analyzing ...\n")
	}

	config := lines.Config{
		IncludeHidden: opts.hidden,
		NumWorkers:    int(opts.jobs),
	}
	counter := lines.NewCounter(config)
	result, err := counter.Run(opts.dir)
	if err != nil {
		return err
	}

	if opts.json {
		return printJSONOutput(stdout, result)
	}

	printHumanOutput(stdout, result, opts)

	if isTerminal {
		color.New(color.FgGreen).Fprintf(stderr, "\nTime taken: %v to analyze files\n", time.Since(startTime))
	}

	return nil
}

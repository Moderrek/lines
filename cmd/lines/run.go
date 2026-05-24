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
		fmt.Fprintln(stdout, "Lines version 1.2.0 created by @Moderrek")
		return nil
	}

	if opts.help {
		fmt.Fprintln(stderr, "Usage: lines [options]")
		fs.PrintDefaults()
		return nil
	}

	startTime := time.Now()
	if isTerminal && !opts.json {
		fmt.Fprintf(stderr, "Analyzing.. %s\n\n", opts.dir)
	}

	config := lines.Config{
		IncludeHidden: opts.hidden,
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

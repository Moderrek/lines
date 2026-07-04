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

	isStdoutTerminal := isatty.IsTerminal(os.Stdout.Fd())
	isStderrTerminal := isatty.IsTerminal(os.Stderr.Fd())

	useColor := (isStdoutTerminal || opts.color) && !opts.noColor
	color.NoColor = !useColor

	if opts.version {
		fmt.Printf("%s version %s created by %s\n", PROGRAM_NAME, VERSION, AUTHOR)
		return nil
	}

	if opts.help {
		fmt.Fprintf(stdout, "Usage: %s [options]\n", PROGRAM_NAME)
		fs.SetOutput(stdout)
		fs.PrintDefaults()
		return nil
	}

	config := lines.Config{
		IncludeHidden: opts.hidden,
		NumWorkers:    int(opts.jobs),
	}
	counter := lines.NewCounter(config)

	stopProgress := make(chan struct{})
	defer close(stopProgress)

	startTime := time.Now()

	showProgress := isStderrTerminal && !opts.json
	if showProgress {
		go startProgressReporter(stderr, stopProgress, startTime, counter)
	}

	result, err := counter.Run(opts.dir)
	if err != nil {
		return err
	}

	if showProgress {
		stopProgress <- struct{}{}
		reportProgress(stderr, startTime, counter)
		fmt.Fprintf(stderr, "\n")
	}

	if opts.json {
		return printJSONOutput(stdout, result)
	}

	printHumanOutput(stdout, result, opts)

	return nil
}

func startProgressReporter(w io.Writer, stop chan struct{}, startTime time.Time, counter *lines.Counter) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			reportProgress(w, startTime, counter)
		}
	}
}

func reportProgress(w io.Writer, startTime time.Time, counter *lines.Counter) {
	processed := counter.FilesProcessed.Load()
	found := counter.FilesFound.Load()

	elapsed := time.Since(startTime)

	inQueue := found - processed
	if inQueue < 0 {
		inQueue = 0
	}

	fmt.Fprintf(w, "\r\033[KProcessed: %d | In Queue: %d | Elapsed: %v", processed, inQueue, elapsed)
}

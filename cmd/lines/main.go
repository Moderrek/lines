package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/moderrek/lines/pkg/lines"
)

var options struct {
	dir     string
	version bool
	help    bool
	hidden  bool
	top     uint
	noColor bool
}

func main() {
	flag.StringVar(&options.dir, "dir", ".", "The directory to analyze")
	flag.BoolVar(&options.version, "version", false, "Print the version and exit")
	flag.BoolVar(&options.help, "help", false, "Print the help message and exit")
	flag.BoolVar(&options.hidden, "hidden", false, "Allows to analize hidden files")
	flag.UintVar(&options.top, "top", 0, "Print the top N extensions")
	flag.BoolVar(&options.noColor, "no-color", false, "Disable color output")
	flag.Parse()

	if options.noColor {
		color.NoColor = true
	}

	if options.version {
		color.Green("Lines version 1.1.0 created by @Moderrek")
		return
	}

	if options.help {
		color.Yellow("Usage: lines [options]")
		flag.PrintDefaults()
		return
	}

	startTime := time.Now()
	fmt.Printf("Analyzing.. %s\n\n", options.dir)

	config := lines.Config{
		IncludeHidden: options.hidden,
	}
	counter := lines.NewCounter(config)
	result, err := counter.Run(options.dir)
	if err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}

	lineMap := result.LinesByExtension
	sortedKeys := make([]string, 0, len(lineMap))
	for key := range lineMap {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		return lineMap[sortedKeys[i]] > lineMap[sortedKeys[j]]
	})

	printCounter := uint(0)
	for _, key := range sortedKeys {
		if options.top > 0 && printCounter >= options.top {
			break
		}
		linesCount := lineMap[key]
		if linesCount == 0 {
			continue
		}
		printCounter++
		color.New(color.Bold).Printf("%d. %s | Lines of code: %d\n", printCounter, key, linesCount)
	}

	color.Green("\nTime taken: %v to analyze files\n", time.Since(startTime))
}

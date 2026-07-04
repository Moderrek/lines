package main

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/fatih/color"
	"github.com/moderrek/lines/pkg/lines"
)

func printJSONOutput(w io.Writer, result *lines.Result) error {
	jsonOutput, err := json.MarshalIndent(result.LinesByExtension, "", "  ")
	if err != nil {
		return fmt.Errorf("error generating JSON: %w", err)
	}
	fmt.Fprintln(w, string(jsonOutput))
	return nil
}

func printHumanOutput(w io.Writer, result *lines.Result, opts *cliOptions) {
	lineMap := result.LinesByExtension

	sortedKeys := make([]string, 0, len(lineMap))
	for key := range lineMap {
		sortedKeys = append(sortedKeys, key)
	}

	sort.Slice(sortedKeys, func(i, j int) bool {
		return lineMap[sortedKeys[i]] > lineMap[sortedKeys[j]]
	})

	extColor := color.New(color.Bold)
	linesColor := color.New(color.FgWhite)

	for i, key := range sortedKeys {
		if opts.top > 0 && uint(i) >= opts.top {
			break
		}

		linesCount := lineMap[key]
		if linesCount == 0 {
			continue
		}

		extColor.Fprintf(w, "%s", key)
		fmt.Fprint(w, "\t")
		linesColor.Fprintf(w, "%d\n", linesCount)
	}
}

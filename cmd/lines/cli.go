package main

import (
	"flag"
	"io"
)

type cliOptions struct {
	dir     string
	version bool
	help    bool
	hidden  bool
	top     uint
	noColor bool
	color   bool
	json    bool
}

func parseFlags(stderr io.Writer, args []string) (*cliOptions, *flag.FlagSet, error) {
	opts := &cliOptions{}
	fs := flag.NewFlagSet("lines", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.StringVar(&opts.dir, "dir", ".", "The directory to analyze")
	fs.BoolVar(&opts.version, "version", false, "Print the version and exit")
	fs.BoolVar(&opts.help, "help", false, "Print the help message and exit")
	fs.BoolVar(&opts.hidden, "hidden", false, "Allows to analize hidden files")
	fs.UintVar(&opts.top, "top", 0, "Print the top N extensions")
	fs.BoolVar(&opts.noColor, "no-color", false, "Disable color output")
	fs.BoolVar(&opts.color, "color", false, "Force color output (e.g. when piping)")
	fs.BoolVar(&opts.json, "json", false, "Output results in JSON format")

	err := fs.Parse(args[1:])
	if err != nil {
		return nil, nil, err
	}

	return opts, fs, nil
}

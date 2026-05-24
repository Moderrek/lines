package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(os.Stdout, os.Stderr, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

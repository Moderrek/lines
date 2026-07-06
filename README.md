

<div align="center">



   # ⚡ Blazingly FAST Line Counter
   
   ![GitHub License](https://img.shields.io/github/license/Moderrek/lines)
   [![Go](https://github.com/Moderrek/lines/actions/workflows/go.yml/badge.svg)](https://github.com/Moderrek/lines/actions/workflows/go.yml)
   ![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/Moderrek/lines/total)


</div>

A concurrent non-blank line counter for source code directories, written in Go.
It recursively walks a directory, concurrently analyzes files, and reports the number of non-blank lines of code, grouped by file extension.
The tool is designed for performance, utilizing goroutines to process files in parallel.

## Installation

To install the `lines` command-line tool, ensure you have [Go](https://go.dev/doc/install) installed and configured, then run:

```shell
go install github.com/moderrek/lines/cmd/lines@latest
```

This will download the source, compile it, and place the `lines` binary in your Go bin directory (`$GOPATH/bin` or `$HOME/go/in`).

## Usage

The `lines` command accepts the following flags:

```
Usage: lines [options]

Options:
  -color
        Force color output (e.g. when piping)
  -no-color
        Disable color output
  -help
        Print the help message
  -hidden
        Allows to analyze hidden files
  -jobs uint
        Specifies the number of jobs
  -json
        Output results in JSON format
  -top uint
        Print the top N extensions
  -verbose
        Verbose output
  -version
        Print the version
```

### Example

To analyze the directory `~/projects/my-app` and display the top 5 extensions:

```shell
lines --top 5 ~/projects/my-app
```

To get the output in JSON format, which can be piped to other tools like `jq`:

```shell
lines --json ~/projects/my-app
```

Example output (`--json`):
```json
{
    ".css": 1122,
    ".go": 15230,
    ".html": 4357,
    ".js": 8828,
    ".mod": 4980
}
```

## Library Usage

The core counting logic is available as a library.
It can be imported into other Go projects.

```go
import "github.com/moderrek/lines/pkg/lines"
```

### Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/moderrek/lines/pkg/lines"
)

func main() {
	// Configure the counter with default settings.
	config := lines.Config{
		IncludeHidden: false,
	}
	counter := lines.NewCounter(config)

	// Run the analysis on the current directory.
	result, err := counter.Run(".")
	if err != nil {
		log.Fatalf("Analysis failed: %v", err)
	}

	// Print results.
	for ext, count := range result.LinesByExtension {
		fmt.Printf("Extension: %s, Lines: %d\n", ext, count)
	}
}
```

### Custom Configuration

You can customize which directories and file extensions to ignore:

```go
config := lines.Config{
	IncludeHidden: false,
	IgnoredDirs: map[string]struct{}{
		"node_modules": {},
		".git":         {},
	},
	IgnoredExtensions: map[string]struct{}{
		".exe": {},
		".env": {},
	},
}
counter := lines.NewCounter(config)
result, err := counter.Run("./src")
```

If `IgnoredDirs` or `IgnoredExtensions` are not provided, the library uses sensible defaults.

## Building from Source

1. Clone the repository:
   ```shell
   git clone https://github.com/Moderrek/lines.git
   ```
2. Navigate to the project directory:
   ```shell
   cd lines
   ```
3. Build the binary:
   ```shell
   go build ./cmd/lines
   ```
   This will create a `lines` executable in the current directory.

## License

This project is licensed under the MIT License.
See the [LICENSE](./LICENSE) file for details.



<div align="center">



   # ⚡ Blazingly FAST Line Counter
   
   ![GitHub License](https://img.shields.io/github/license/Moderrek/lines)
   [![Go](https://github.com/Moderrek/lines/actions/workflows/go.yml/badge.svg)](https://github.com/Moderrek/lines/actions/workflows/go.yml)
   ![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/Moderrek/lines/total)


</div>

A concurrent, non-blank line counter for source code directories, written in GO.

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
  -dir string
        The directory to analyze (default ".")
  -hidden
        Include hidden files and directories in the analysis
  -top uint
        Show only the top N extensions by line count
  -no-color
        Disable colorized output
  -version
        Print version information and exit
  -help
        Show this help message and exit
```

### Example

To analyze the directory `~/projects/my-app` and dispaly the top 5 extensions:

```shell
lines --dir ~/projects/my-app --top 5
```

Example output:
```
Analyzing.. ~/projects/my-app

1. .java | Lines of code: 24304
2. .json | Lines of code: 8828
3. .yaml | Lines of code 4980
4. .tsx  | Lines of code: 4290
5. .yml  | Lines of code: 1122
```

## Library Usage

The core couting logic is available as a library.
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
   // Configure the counter.
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
See the [LICENSE](LICENSE) file for details.

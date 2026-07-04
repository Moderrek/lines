package lines

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

// countNonBlankLines reads a file and counts non-blank, non-comment lines.
// Lines starting with '//', '#' or '--' are treated as comments and skipped.
// path specifies the file path to count non-blank lines.
// bufferInitialSize specifies the initial scanner buffer size.
// bufferMaxSize specifies the maximum scanner buffer size.
func countNonBlankLines(path string, bufferInitialSize, bufferMaxSize int) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warn: failed to close file %q: %v\n", path, err)
		}
	}()

	scanner := bufio.NewScanner(file)
	buffer := make([]byte, 0, bufferInitialSize)
	scanner.Buffer(buffer, bufferMaxSize)

	commentDoubleSlash := []byte("//")
	commentHash := []byte("#")
	commentDoubleDash := []byte("--")

	lineCount := 0

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())

		// Skip empty lines.
		if len(line) == 0 {
			continue
		}

		// Skip comment lines: //, #, or --.
		if bytes.HasPrefix(line, commentDoubleSlash) || bytes.HasPrefix(line, commentHash) || bytes.HasPrefix(line, commentDoubleDash) {
			continue
		}

		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

package lines

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

// analyzeFile reads a file and counts non-blank, non-comment lines.
// Lines starting with '//', '#' or '--' are treated as comments and skipped.
// path specifies the file path to count non-blank lines.
// bufferInitialSize specifies the initial scanner buffer size.
// bufferMaxSize specifies the maximum scanner buffer size.
func analyzeFile(path string, readerInitialBufferSize int) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warn: failed to close file %q: %v\n", path, err)
		}
	}()

	var r *bufio.Reader
	if readerInitialBufferSize > 0 {
		r = bufio.NewReaderSize(file, readerInitialBufferSize)
	} else {
		r = bufio.NewReader(file)
	}

	commentDoubleSlash := []byte("//")
	commentHash := []byte("#")
	commentDoubleDash := []byte("--")

	lineCount := 0
	isInsideLongLine := false

	for {
		line, isPrefix, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}

		if isInsideLongLine {
			if !isPrefix {
				isInsideLongLine = false
			}
			continue
		}

		if isPrefix {
			isInsideLongLine = true
		}

		cleaned := bytes.TrimSpace(line)

		// Skip empty lines.
		if len(cleaned) == 0 {
			continue
		}

		// Skip comment lines: //, #, or --.
		if bytes.HasPrefix(cleaned, commentDoubleSlash) || bytes.HasPrefix(cleaned, commentHash) || bytes.HasPrefix(cleaned, commentDoubleDash) {
			continue
		}

		lineCount++
	}

	return lineCount, nil
}

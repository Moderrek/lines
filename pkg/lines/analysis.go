package lines

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// analyzeFile reads a file and counts non-blank, non-comment lines.
// Lines starting with '//', '#' or '--' are treated as comments and skipped.
// path specifies the file path to count non-blank lines.
// bufferInitialSize specifies the initial scanner buffer size.
// bufferMaxSize specifies the maximum scanner buffer size.
func analyzeFile(path string, readerInitialBufferSize int) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("failed to analyze %q: %v", filepath.ToSlash(path), err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warn: failed to close file %q: %v\n", filepath.ToSlash(path), err)
		}
	}()

	reader := newReader(file, readerInitialBufferSize)
	lineCount := 0
	isInsideLongLine := false

	for {
		line, err := readLine(reader, path, &isInsideLongLine)
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}

		if len(line) == 0 {
			continue
		}

		if isCommentLine(line) {
			continue
		}

		lineCount++
	}

	return lineCount, nil
}

func readLine(reader *bufio.Reader, path string, isInsideLongLine *bool) ([]byte, error) {
	line, isPrefix, err := reader.ReadLine()
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("failed to analyze %q: %v", filepath.ToSlash(path), err)
	}

	if *isInsideLongLine {
		if !isPrefix {
			*isInsideLongLine = false
		}
		return nil, nil
	}

	if isPrefix {
		*isInsideLongLine = true
	}

	trimmedLine := bytes.TrimSpace(line)
	return trimmedLine, nil
}

func isCommentLine(line []byte) bool {
	doubleSlash := []byte("//")
	hash := []byte("#")
	doubleDash := []byte("--")

	return bytes.HasPrefix(line, doubleSlash) || bytes.HasPrefix(line, hash) || bytes.HasPrefix(line, doubleDash)
}

func newReader(file *os.File, initialBufferSize int) *bufio.Reader {
	if initialBufferSize > 0 {
		return bufio.NewReaderSize(file, initialBufferSize)
	}
	return bufio.NewReader(file)
}

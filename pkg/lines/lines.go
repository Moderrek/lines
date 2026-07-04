package lines

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Result represents the results of the line counting process.
type Result struct {
	// LinesByExtension maps file extensions to their total line counts.
	LinesByExtension map[string]int
}

// Run analyzes the given directory and returns the results.
// It recursively walks the directory tree using goroutines for performance.
func (c *Counter) Run(dir string) (*Result, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %q does not exist", dir)
	}

	c.workers.Add(1)
	go c.walkDir(dir)
	c.workers.Wait()

	return &Result{
		LinesByExtension: c.lines.Items(),
	}, nil
}

// countLinesInFile counts non-blank lines in a file and updates results.
func (c *Counter) countLinesInFile(path string) {
	ext := strings.ToLower(filepath.Ext(path))
	c.workers.Add(1)

	go func() {
		defer c.workers.Done()

		lineCount, err := countNonBlankLines(path, c.config.BufferInitialSize, c.config.BufferMaxSize)
		if err != nil {
			if _, writeErr := fmt.Fprintf(os.Stderr, "warn: failed to count lines in %q: %v\n", path, err); writeErr != nil {
				return
			}
			return
		}

		if lineCount <= 0 {
			return
		}

		// Updates result.
		c.lines.Upsert(ext, lineCount, func(exists bool, valueInMap int, newValue int) int {
			if exists {
				return valueInMap + newValue
			}
			return newValue
		})
	}()
}

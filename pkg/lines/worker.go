package lines

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (c *Counter) worker() {
	defer c.workers.Done()

	for path := range c.filesToAnalyze {
		lineCount, err := countNonBlankLines(path, c.config.BufferInitialSize, c.config.BufferMaxSize)
		c.FilesProcessed.Add(1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: failed to count lines in %q: %v\n", path, err)
			continue
		}

		if lineCount > 0 {
			ext := strings.ToLower(filepath.Ext(path))
			c.addLineCount(ext, lineCount)
		}
	}
}

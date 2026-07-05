package lines

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (c *Counter) analyzeFilesWorker() {
	defer c.workers.Done()

	for path := range c.filesToAnalyze {
		lineCount, err := analyzeFile(path, c.Config.ReaderInitialBufferSize)
		c.FilesProcessed.Add(1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: failed to count lines in %q: %v\n", filepath.ToSlash(path), err)
			continue
		}

		if lineCount > 0 {
			ext := strings.ToLower(filepath.Ext(path))
			c.addLineCount(ext, lineCount)
			c.logVerbosef("file %q had %d lines", filepath.ToSlash(path), lineCount)

		}
	}
}

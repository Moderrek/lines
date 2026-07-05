package lines

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// walkDir recursively walks the directory tree and enqueues files for the analyzeFilesWorker pool to analyze.
func (c *Counter) walkDir(dir string) error {
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: failed to access path %q: %v\n", filepath.ToSlash(path), err)
			return nil
		}

		if d.IsDir() {
			if path == dir {
				return nil
			}
			name := d.Name()
			if !c.Config.IncludeHidden && strings.HasPrefix(name, ".") {
				c.logVerbosef("skipping hidden directory: %q", filepath.ToSlash(path))
				return filepath.SkipDir
			}
			if c.isIgnoredDir(name) {
				c.logVerbosef("skipping ignored directory: %q", filepath.ToSlash(path))
				return filepath.SkipDir
			}
			return nil
		}

		if d.Type().IsRegular() && c.shouldAnalyzeFile(path, d.Name()) {
			c.logVerbosef("found file to analyze: %q", filepath.ToSlash(path))
			c.FilesFound.Add(1)
			c.filesToAnalyze <- path
		} else {
			c.logVerbosef("skipping file: %q", filepath.ToSlash(path))
		}

		return nil
	})
	return err
}

// shouldAnalyzeFile determines if a file should be analyzed.
// Returns false if the file is hidden (when IncludeHidden is false),
// has no extension, or has an ignored extension.
func (c *Counter) shouldAnalyzeFile(path, filename string) bool {
	if !c.Config.IncludeHidden && strings.HasPrefix(filename, ".") {
		return false
	}

	ext := filepath.Ext(path)
	if len(ext) == 0 {
		// probably its binary file
		return false
	}

	return !c.isIgnoredExtension(ext)
}

// isIgnoredDir checks if a directory should be ignored.
func (c *Counter) isIgnoredDir(dirname string) bool {
	_, ok := c.Config.IgnoredDirs[dirname]
	return ok
}

// isIgnoredExtension checks if a file extension should be ignored for line counting.
// The comparison is case-insensitive. Extension should begin with dot.
func (c *Counter) isIgnoredExtension(ext string) bool {
	_, ok := c.Config.IgnoredExtensions[strings.ToLower(ext)]
	return ok
}

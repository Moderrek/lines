package lines

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// walkDir recursively walks the directory tree and enqueues files for the worker pool to analyze.
func (c *Counter) walkDir(dir string) {
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			if path == dir {
				return nil
			}
			name := d.Name()
			if !c.config.IncludeHidden && strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			if c.isIgnoredDir(name) {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Type().IsRegular() && c.needToAnalyze(path, d.Name()) {
			c.FilesFound.Add(1)
			c.filesToAnalyze <- path
		}

		return nil
	})
}

// needToAnalyze determines if a file should be analyzed.
// Returns false if the file is hidden (when IncludeHidden is false),
// has no extension, or has an ignored extension.
func (c *Counter) needToAnalyze(path, filename string) bool {
	if !c.config.IncludeHidden && filename[0] == '.' {
		return false
	}

	extension := filepath.Ext(path)
	if len(extension) == 0 {
		return false
	}

	return !c.isIgnoredExtension(extension)
}

// isIgnoredDir checks if a directory should be ignored.
func (c *Counter) isIgnoredDir(dirname string) bool {
	for _, ignored := range c.config.IgnoredDirs {
		if dirname == ignored {
			return true
		}
	}
	return false
}

// isIgnoredExtension checks if a file extension should be ignored for line counting.
// The comparison is case-insensitive.
func (c *Counter) isIgnoredExtension(ext string) bool {
	_, ok := c.config.IgnoredExtensions[strings.ToLower(ext)]
	return ok
}

package lines

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// walkDir recursively walks the directory tree and counts lines in files.
// It spawns goroutines for each subdirectory to achieve parallel processing.
func (c *Counter) walkDir(dir string) {
	defer c.workers.Done()

	filepath.WalkDir(dir, c.walkFn(dir))
}

func (c *Counter) walkFn(root string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return c.handleWalkError(path, err)
		}

		if d.IsDir() {
			return c.handleDirectory(root, path)
		}

		return c.handleFile(path, d)
	}
}

func (c *Counter) handleWalkError(path string, err error) error {
	fmt.Fprintf(os.Stderr, "error: cannot read %q: %v\n", path, err)
	return err
}

func (c *Counter) handleDirectory(root string, path string) error {
	if path == root {
		return nil
	}

	name := filepath.Base(path)

	if !c.config.IncludeHidden && strings.HasPrefix(name, ".") {
		return filepath.SkipDir
	}

	if c.isIgnoredDir(name) {
		return filepath.SkipDir
	}

	c.workers.Add(1)
	go c.walkDir(path)

	return filepath.SkipDir
}

func (c *Counter) handleFile(path string, d fs.DirEntry) error {
	if !d.Type().IsRegular() {
		return nil
	}

	if c.needToAnalyze(path) {
		c.countLinesInFile(path)
	}

	return nil
}

// needToAnalyze determines if a file should be analyzed.
// Returns false if the file is hidden (when IncludeHidden is false),
// has no extension, or has an ignored extension.
func (c *Counter) needToAnalyze(path string) bool {
	if !c.config.IncludeHidden && filepath.Base(path)[0] == '.' {
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

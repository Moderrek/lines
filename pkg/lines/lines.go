package lines

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	cmap "github.com/orcaman/concurrent-map/v2"
)

// Config holds settings for the line counting process.
type Config struct {
	// IncludeHidden analyzes hidden files and directories (starting with '.').
	IncludeHidden bool
	// IgnoredDirs are directories to skip during analysis. Defaults to ["node_modules", "vendor", ".git", "target"].
	IgnoredDirs []string
	// IgnoredExtensions are file extensions to skip. Defaults to common binary and media formats.
	IgnoredExtensions map[string]struct{}
	// BufferInitialSize is the initial buffer size for the scanner. Defaults to 64KB.
	BufferInitialSize int
	// BufferMaxSize is the maximum buffer size for the scanner. Defaults to 1MB.
	BufferMaxSize int
}

// Represents the results of the line counting process.
type Result struct {
	// LinesByExtension maps file extensions to their total line counts.
	LinesByExtension map[string]int
}

// Counter analyzes directories and counts non-blank lines of code.
// NOTE: Counter is safe for concurrent use and uses goroutines internally.
type Counter struct {
	config  Config
	lines   cmap.ConcurrentMap[string, int]
	workers sync.WaitGroup
}

// NewCounter creates a new Counter with the given configuration.
// If IgnoredDirs or IgnoredExtensions are empty, sensible defaults are used.
func NewCounter(config Config) *Counter {
	// Use sensible defaults if lists are empty.
	if len(config.IgnoredDirs) == 0 {
		config.IgnoredDirs = DefaultIgnoredDirs()
	}
	if len(config.IgnoredExtensions) == 0 {
		config.IgnoredExtensions = DefaultIgnoredExtensions()
	}
	if config.BufferInitialSize == 0 {
		config.BufferInitialSize = 64 * 1024
	}
	if config.BufferMaxSize == 0 {
		config.BufferMaxSize = 1024 * 1024
	}

	return &Counter{
		config: config,
		lines:  cmap.New[int](),
	}
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

// Run analyzes the given directory and returns the results.
// It recursively walks the directory tree using goroutines for performance.
func (c *Counter) Run(dir string) (*Result, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory '%s' does not exist", dir)
	}

	c.workers.Add(1)
	go c.walkDir(dir)
	c.workers.Wait()

	result := &Result{
		LinesByExtension: c.lines.Items(),
	}
	return result, nil
}

// walkDir recursively walks the directory tree and counts lines in files.
// It spawns goroutines for each subdirectory to achieve parallel processing.
func (c *Counter) walkDir(dir string) {
	defer c.workers.Done()

	visit := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			// NOTE: Log access errors but continue with other directories.
			fmt.Fprintf(os.Stderr, "ERROR: cannot access path %q: %v\n", path, err)
			return err
		}
		if f.IsDir() && path != dir {
			dirname := filepath.Base(path)
			if !c.config.IncludeHidden && dirname[0] == '.' {
				return filepath.SkipDir
			}
			if c.isIgnoredDir(dirname) {
				return filepath.SkipDir
			}
			c.workers.Add(1)
			go c.walkDir(path)
			return filepath.SkipDir
		}
		if f.Mode().IsRegular() {
			if c.needToAnalyze(path) {
				c.fastLineCounter(path)
			}
		}
		return nil
	}
	filepath.Walk(dir, visit)
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
	if c.isIgnoredExtension(extension) {
		return false
	}
	return true
}

// fastLineCounter counts non-blank lines in a file and updates results.
// TODO: Consider caching results for frequently accessed files.
func (c *Counter) fastLineCounter(path string) {
	extension := strings.ToLower(filepath.Ext(path))
	c.workers.Add(1)
	go func() {
		defer c.workers.Done()
		countedLines, err := countNonBlankLines(path, c.config.BufferInitialSize, c.config.BufferMaxSize)
		if err != nil {
			// NOTE: Silently skip files with read/encoding issues.
			fmt.Fprintf(os.Stderr, "WARNING: failed to count lines in %q: %v\n", path, err)
			return
		}
		if countedLines > 0 {
			c.lines.Upsert(extension, countedLines, func(exists bool, valueInMap int, newValue int) int {
				if exists {
					return valueInMap + newValue
				}
				return newValue
			})
		}
	}()
}

// countNonBlankLines reads a file and counts non-blank, non-comment lines.
// Lines starting with '//' or '#' are treated as comments and skipped.
// bufferInitialSize specifies the initial scanner buffer size.
// bufferMaxSize specifies the maximum scanner buffer size.
func countNonBlankLines(path string, bufferInitialSize, bufferMaxSize int) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buffer := make([]byte, 0, bufferInitialSize)
	scanner.Buffer(buffer, bufferMaxSize)

	lineCounter := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines.
		if line == "" {
			continue
		}
		// Skip comment lines: //, #, or --.
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "--") {
			continue
		}
		lineCounter++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCounter, nil
}

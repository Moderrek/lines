package lines

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

// Counter analyzes directories and counts non-blank lines of code.
type Counter struct {
	config Config

	linesLock sync.Mutex
	lines     map[string]int

	workers        sync.WaitGroup
	filesToAnalyze chan string

	FilesFound     atomic.Int64
	FilesProcessed atomic.Int64
}

// NewCounter creates a new Counter with the given configuration.
// If IgnoredDirs or IgnoredExtensions are empty, sensible defaults are used.
func NewCounter(cfg Config) *Counter {
	if len(cfg.IgnoredDirs) == 0 {
		cfg.IgnoredDirs = DefaultIgnoredDirs()
	}
	if len(cfg.IgnoredExtensions) == 0 {
		cfg.IgnoredExtensions = DefaultIgnoredExtensions()
	}
	if cfg.BufferInitialSize == 0 {
		cfg.BufferInitialSize = 64 * 1024
	}
	if cfg.BufferMaxSize == 0 {
		cfg.BufferMaxSize = 1024 * 1024
	}
	if cfg.NumWorkers <= 0 {
		cfg.NumWorkers = runtime.NumCPU() * 2
	}

	return &Counter{
		config: cfg,
		lines:  make(map[string]int),
	}
}

// Run analyzes the given directory and returns the results.
// It recursively walks the directory tree using goroutines for performance.
func (c *Counter) Run(dir string) (*Result, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %q does not exist", dir)
	}

	numWorkers := c.config.NumWorkers
	c.filesToAnalyze = make(chan string, numWorkers*4)

	for i := 0; i < numWorkers; i++ {
		c.workers.Add(1)
		go c.worker()
	}

	c.walkDir(dir)
	close(c.filesToAnalyze)
	c.workers.Wait()

	return &Result{
		LinesByExtension: c.lines,
	}, nil
}

func (c *Counter) addLineCount(ext string, count int) {
	c.linesLock.Lock()
	c.lines[ext] += count
	c.linesLock.Unlock()
}

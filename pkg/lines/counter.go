package lines

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
)

// Counter analyzes directories and counts non-blank lines of code.
type Counter struct {
	Config Config

	linesLock sync.Mutex
	lines     map[string]int

	workers        sync.WaitGroup
	filesToAnalyze chan string

	FilesFound     atomic.Int64
	FilesProcessed atomic.Int64
}

// NewCounter creates a new Counter with the given configuration.
// If IgnoredDirs or IgnoredExtensions are empty, sensible defaults are used.
func NewCounter(config Config) *Counter {
	return &Counter{
		Config:         configWithDefaults(config),
		linesLock:      sync.Mutex{},
		lines:          make(map[string]int),
		workers:        sync.WaitGroup{},
		filesToAnalyze: nil,
		FilesFound:     atomic.Int64{},
		FilesProcessed: atomic.Int64{},
	}
}

func configWithDefaults(config Config) Config {
	if len(config.IgnoredDirs) == 0 {
		config.IgnoredDirs = DefaultIgnoredDirs()
	}
	if len(config.IgnoredExtensions) == 0 {
		config.IgnoredExtensions = DefaultIgnoredExtensions()
	}
	if config.ReaderInitialBufferSize == 0 {
		config.ReaderInitialBufferSize = 64 * 1024
	}
	if config.NumWorkers <= 0 {
		config.NumWorkers = runtime.NumCPU() * 2
	}
	return config
}

// Run analyzes the given directory and returns the results.
// It recursively walks the directory tree using goroutines for performance.
func (c *Counter) Run(dir string) (*Result, error) {
	c.reset()

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %q does not exist", dir)
	}

	numWorkers := c.Config.NumWorkers
	maximumWaitingWork := numWorkers * 4
	c.filesToAnalyze = make(chan string, maximumWaitingWork)

	c.logVerbosef("creating %d workers with work queue of capacity %d", numWorkers, maximumWaitingWork)
	for i := 0; i < numWorkers; i++ {
		c.workers.Add(1)
		go c.analyzeFilesWorker()
	}

	c.logVerbosef("starting analysis at: %q", filepath.ToSlash(dir))
	err := c.walkDir(dir)
	if err != nil {
		return nil, err
	}
	close(c.filesToAnalyze)
	c.workers.Wait()

	// Makes copy of result.
	c.logVerbosef("coping result")
	linesByExt := make(map[string]int)
	for ext, count := range c.lines {
		linesByExt[ext] = count
	}

	return &Result{
		LinesByExtension: linesByExt,
	}, nil
}

func (c *Counter) addLineCount(ext string, count int) {
	c.linesLock.Lock()
	c.lines[ext] += count
	c.linesLock.Unlock()
}

func (c *Counter) reset() {
	c.workers = sync.WaitGroup{}
	c.FilesFound.Store(0)
	c.FilesProcessed.Store(0)

	c.linesLock.Lock()
	c.lines = make(map[string]int)
	c.linesLock.Unlock()
}

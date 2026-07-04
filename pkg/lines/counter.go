package lines

import (
	"sync"

	cmap "github.com/orcaman/concurrent-map/v2"
)

// TODO: migrate from cmap to stdlib map
// TODO: create workers pool

// Counter analyzes directories and counts non-blank lines of code.
type Counter struct {
	config  Config
	lines   cmap.ConcurrentMap[string, int]
	workers sync.WaitGroup
}

// NewCounter creates a new Counter with the given configuration.
// If IgnoredDirs or IgnoredExtensions are empty, sensible defaults are used.
func NewCounter(config Config) *Counter {
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

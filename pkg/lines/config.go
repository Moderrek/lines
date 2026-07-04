package lines

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
	// NumWorkers is the number of workers to use for file analysis.
	NumWorkers int
}

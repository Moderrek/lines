package lines

// Config holds settings for the line counting process.
type Config struct {
	// IncludeHidden analyzes hidden files and directories (starting with '.').
	IncludeHidden bool
	// IgnoredDirs are directories to skip during analysis. Defaults to ["node_modules", "vendor", ".git", "target"].
	IgnoredDirs map[string]struct{}
	// IgnoredExtensions are file extensions to skip. Defaults to common binary and media formats.
	IgnoredExtensions map[string]struct{}
	// ReaderInitialBufferSize is the initial buffer size for the scanner.
	ReaderInitialBufferSize int
	// NumWorkers is the number of workers to use for file analysis.
	NumWorkers int
	// Verbose enables logging while file analysis.
	Verbose bool
}

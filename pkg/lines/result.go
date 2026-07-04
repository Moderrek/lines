package lines

// Result represents the results of the line counting process.
type Result struct {
	// LinesByExtension maps file extensions to their total line counts.
	LinesByExtension map[string]int
}

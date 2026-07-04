package lines

// DefaultIgnoredDirs returns default directories to ignore.
func DefaultIgnoredDirs() []string {
	return []string{
		"node_modules", "vendor", ".git", "target",
	}
}

// DefaultIgnoredExtensions returns default file extensions to ignore.
func DefaultIgnoredExtensions() []string {
	return []string{
		".exe", ".dll", ".so", ".dylib",
		".zip", ".tar", ".gz", ".bz2", ".xz",
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg", ".ico",
		".mp3", ".wav", ".flac", ".ogg", ".aac",
		".mp4", ".mkv", ".avi", ".mov", ".wmv",
		".pdf", ".doc", ".docx", ".xls", ".xlsx",
		".icns", ".ttf", ".otf", ".woff", ".woff2",
		".eot", ".svgz", ".uasset", ".plist",
		".url", ".pbxproj", ".sln",
		".vcxproj", ".csproj", ".vcproj", ".tlog",
		".tmp", ".filters", ".idb", ".lock", ".rc",
		".sqlite", ".gdb", ".node", ".rmeta",
		".rlib", ".mcmeta", ".iml", ".map", ".natvis",
		".d", ".dat_old", ".storyboard", ".ilk", ".ppt",
		".pptx", ".odt", ".ods", ".odp", ".odg", ".mca",
		".psd", ".bin", ".jar", ".pdb", ".dox", ".db",
		".schem", ".lnk", ".mod", ".lib", ".o", ".obj",
		".a", ".class", ".pyc", ".pyo", ".whl", ".log",
		".in", ".dat", ".TAG", ".repositories", ".MF",
	}
}

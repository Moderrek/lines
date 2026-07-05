package lines

import "strings"

// DefaultIgnoredDirs returns slice of default directories to ignore.
func DefaultIgnoredDirs() map[string]struct{} {
	return map[string]struct{}{
		"node_modules": {},
		"vendor":       {},
		".git":         {},
		"target":       {},
	}
}

// DefaultIgnoredExtensions returns set of default file extensions to ignore.
func DefaultIgnoredExtensions() map[string]struct{} {
	return makeExtensionSet(
		"exe", "dll", "so", "dylib", "msi", "mui", "mun",
		"zip", "tar", "gz", "bz2", "xz",
		"jpg", "jpeg",
		"png", "dng", "heic",
		"gif", "bmp", "webp", "svg", "ico",
		"mp3", "wav", "flac", "ogg", "aac",
		"mp4", "mkv", "avi", "mov", "wmv",
		"pdf", "doc", "docx", "xls", "xlsx",
		"icns", "ttf", "otf", "woff", "woff2",
		"eot", "svgz", "uasset", "plist",
		"url", "pbxproj", "sln",
		"vcxproj", "csproj", "vcproj", "tlog",
		"tmp", "filters", "idb", "lock", "rc",
		"sqlite", "gdb", "node", "rmeta",
		"rlib", "mcmeta", "iml", "map", "natvis",
		"d", "dat_old", "storyboard", "ilk", "ppt",
		"pptx", "odt", "ods", "odp", "odg", "mca",
		"psd", "bin", "jar", "pdb", "dox", "db",
		"schem", "lnk", "mod", "lib", "o", "obj",
		"a", "class", "pyc", "pyo", "whl", "log",
		"in", "dat", "TAG", "repositories", "MF",
	)
}

// makeExtensionSet makes set of file extensions for faster search.
// The extensions are stored in lowercase and prefixed with a dot.
// NOTE: Extensions are prefixed with dot to avoid stripping out the dot for every file.
func makeExtensionSet(items ...string) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set["."+strings.ToLower(item)] = struct{}{}
	}
	return set
}

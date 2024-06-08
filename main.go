package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

var options struct {
	dir     string
	version bool
	help    bool
	hidden  bool
	top     uint
}

var workers sync.WaitGroup
var lines = cmap.New[int]()

var notAllowedExtensions = map[string]bool{
	".exe": true, ".dll": true, ".so": true, ".dylib": true,
	".zip": true, ".tar": true, ".gz": true, ".bz2": true, ".xz": true,
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true, ".webp": true, ".svg": true, ".ico": true,
	".mp3": true, ".wav": true, ".flac": true, ".ogg": true, ".aac": true,
	".mp4": true, ".mkv": true, ".avi": true, ".mov": true, ".wmv": true,
	".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
	".icns": true, ".ttf": true, ".otf": true, ".woff": true, ".woff2": true,
	".eot": true, ".svgz": true, ".uasset": true, ".plist": true,
	".url": true, ".pbxproj": true, ".sln": true,
	".vcxproj": true, ".csproj": true, ".vcproj": true, ".tlog": true,
	".tmp": true, ".filters": true, ".idb": true, ".lock": true, ".rc": true,
	".sqlite": true, ".gdb": true, ".node": true, ".rmeta": true,
	".rlib": true, ".mcmeta": true, ".iml": true, ".map": true, ".natvis": true,
	".d": true, ".dat_old": true, ".storyboard": true, ".ilk": true, ".ppt": true,
	".pptx": true, ".odt": true, ".ods": true, ".odp": true, ".odg": true, ".mca": true,
	".psd": true, ".bin": true, ".jar": true, ".pdb": true, ".dox": true, ".db": true,
	".schem": true, ".lnk": true, ".mod": true, ".lib": true, ".o": true, ".obj": true,
	".a": true, ".class": true, ".pyc": true, ".pyo": true, ".whl": true, ".log": true,
	".in": true, "idb": true, ".dat": true, ".TAG": true, ".repositories": true, ".MF": true,
}

var notAllowedDirs = map[string]bool{
	"node_modules": true, "vendor": true, ".git": true, "target": true,
}

func countNonBlankLines(path string) int {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	buffer := make([]byte, 0, 64*1024)
	scanner.Buffer(buffer, 1024*1024)

	line_counter := 0

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		// Skip empty lines
		if len(strings.TrimSpace(line)) > 0 {
			line_counter++
		}
		// Skip comments
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "*") || strings.HasSuffix(trimmed, "*/") || strings.HasPrefix(trimmed, "#") {
			// Skip comments
			_ = trimmed
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to count lines %s: %s\n", path, err)
	}

	return line_counter
}

func fastLineCounter(path string) {
	extension := filepath.Ext(path)
	workers.Add(1)
	go func() {
		defer workers.Done()
		countedLines := countNonBlankLines(path)
		if val, ok := lines.Get(extension); ok {
			lines.Set(extension, val+countedLines)
		} else {
			lines.Set(extension, countedLines)
		}
	}()
}

func needToAnalyze(path string) bool {
	// Skip hidden files
	if !options.hidden && filepath.Base(path)[0] == '.' {
		return false
	}
	extension := filepath.Ext(path)
	// Skip files without extension
	if len(extension) == 0 {
		return false
	}
	// Skip not allowed extensions
	if _, ok := notAllowedExtensions[extension]; ok {
		return false
	}
	return true
}

func walkDir(dir string) {
	defer workers.Done()

	visit := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && path != dir {
			dirname := filepath.Base(path)
			// Skip hidden directory
			if !options.hidden && dirname[0] == '.' {
				return filepath.SkipDir
			}
			// Skip node_modules, vendor, .git, target directories
			if _, ok := notAllowedDirs[dirname]; ok {
				return filepath.SkipDir
			}
			// Walk the directory
			workers.Add(1)
			go walkDir(path)

			return filepath.SkipDir
		}
		if f.Mode().IsRegular() {
			if needToAnalyze(path) {
				fastLineCounter(path)
			}
		}
		return nil
	}
	filepath.Walk(dir, visit)
}

func main() {
	// Parse the command line flags
	flag.StringVar(&options.dir, "dir", ".", "The directory to analyze")
	flag.BoolVar(&options.version, "version", false, "Print the version and exit")
	flag.BoolVar(&options.help, "help", false, "Print the help message and exit")
	flag.BoolVar(&options.hidden, "hidden", false, "Allows to analize hidden files")
	flag.UintVar(&options.top, "top", 0, "Print the top N files")

	flag.Parse()

	// If the version flag is set, print the version and exit
	if options.version {
		fmt.Println("Lines installed version 1.0.0")
		return
	}

	// If the help flag is set, print the help message and exit
	if options.help {
		fmt.Println("Usage: lines [options]")
		flag.PrintDefaults()
		return
	}

	if _, err := os.Stat(options.dir); os.IsNotExist(err) {
		fmt.Printf("Directory %s does not exist\n", options.dir)
		return
	}

	// Get current time. Will be used to calculate the time taken to analyze files
	startTime := time.Now()

	fmt.Printf("Analyzing.. %s\n\n", options.dir)

	// Start the analysis
	workers.Add(1)
	walkDir(options.dir)
	workers.Wait()

	// Convert lines to a regular map
	lineMap := make(map[string]int)
	for _, key := range lines.Keys() {
		lineMap[key], _ = lines.Get(key)
	}

	// Sort the map by value
	sortedKeys := make([]string, 0, len(lineMap))
	for key := range lineMap {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		return lineMap[sortedKeys[i]] > lineMap[sortedKeys[j]]
	})

	// Print the top N extensions
	counter := uint(0)
	for _, key := range sortedKeys {
		if options.top > 0 {
			if counter >= options.top {
				break
			}
			counter++
		}
		var lines int = lineMap[key]
		if lines == 0 {
			continue
		}
		fmt.Printf("%s | Lines of code: %d\n", key, lines)
	}
	fmt.Printf("\nTime taken: %v to analyze files\n", time.Since(startTime))
}

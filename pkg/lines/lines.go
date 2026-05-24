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

type Config struct {
	IncludeHidden bool
}

type Result struct {
	LinesByExtension map[string]int
}

type Counter struct {
	config               Config
	lines                cmap.ConcurrentMap[string, int]
	workers              sync.WaitGroup
	notAllowedDirs       map[string]bool
	notAllowedExtensions map[string]bool
}

func NewCounter(config Config) *Counter {
	return &Counter{
		config: config,
		lines:  cmap.New[int](),
		notAllowedDirs: map[string]bool{
			"node_modules": true, "vendor": true, ".git": true, "target": true,
		},
		notAllowedExtensions: map[string]bool{
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
		},
	}
}

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

func (c *Counter) walkDir(dir string) {
	defer c.workers.Done()

	visit := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: cannot access path %q: %v\n", path, err)
			return err
		}
		if f.IsDir() && path != dir {
			dirname := filepath.Base(path)
			if !c.config.IncludeHidden && dirname[0] == '.' {
				return filepath.SkipDir
			}
			if _, ok := c.notAllowedDirs[dirname]; ok {
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

func (c *Counter) needToAnalyze(path string) bool {
	if !c.config.IncludeHidden && filepath.Base(path)[0] == '.' {
		return false
	}
	extension := filepath.Ext(path)
	if len(extension) == 0 {
		return false
	}
	if _, ok := c.notAllowedExtensions[extension]; ok {
		return false
	}
	return true
}

func (c *Counter) fastLineCounter(path string) {
	extension := filepath.Ext(path)
	c.workers.Add(1)
	go func() {
		defer c.workers.Done()
		countedLines, err := countNonBlankLines(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: cannot count lines in file %q: %v\n", path, err)
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

func countNonBlankLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buffer := make([]byte, 0, 64*1024)
	scanner.Buffer(buffer, 1024*1024)

	lineCounter := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}
		// Skip comment lines (common in many programming languages)
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

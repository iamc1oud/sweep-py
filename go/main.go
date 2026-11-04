package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

var rules = map[string][]string{
	"package.json":    {"node_modules", ".next", ".nuxt", ".turbo", ".parcel-cache"},
	"pyproject.toml":  {".venv", "venv", "__pycache__", ".pytest_cache", ".mypy_cache", ".ruff_cache", "*.egg-info"},
	"requirements.txt": {".venv", "venv", "__pycache__", ".pytest_cache", ".mypy_cache", ".ruff_cache"},
	"setup.py":        {".venv", "venv", "__pycache__", "*.egg-info", "build", "dist"},
	"Cargo.toml":      {"target"},
	"pom.xml":         {"target"},
	"build.gradle":    {"build", ".gradle"},
	"composer.json":   {"vendor"},
	"*.csproj":        {"bin", "obj"},
	"pubspec.yaml":    {".dart_tool", "build"},
	"Package.swift":   {".build"},
	"Gemfile":         {".bundle"},
}

var skip = map[string]bool{".git": true}

type junkFolder struct {
	path string
	size int64
}

// globMatch checks if name matches any pattern in patterns
func globMatch(name string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}
	}
	return false
}

// findJunk yields junk directories under root. Matched dirs are not descended into.
func findJunk(root string) ([]junkFolder, error) {
	var junk []junkFolder
	findJunkRecursive(root, &junk)
	return junk, nil
}

func findJunkRecursive(path string, junk *[]junkFolder) {
	// List files in current directory to check for marker files
	entries, err := os.ReadDir(path)
	if err != nil {
		return // skip on permission errors
	}

	var files []string
	var dirs []string
	for _, entry := range entries {
		name := entry.Name()
		if skip[name] {
			continue
		}
		if entry.IsDir() {
			dirs = append(dirs, name)
		} else {
			files = append(files, name)
		}
	}

	// Find all applicable junk patterns based on marker files
	patterns := make(map[string]bool)
	for markerPattern, markerPatterns := range rules {
		// Check if any file in this directory matches the marker pattern
		for _, file := range files {
			if matched, _ := filepath.Match(markerPattern, file); matched {
				for _, p := range markerPatterns {
					patterns[p] = true
				}
				break
			}
		}
	}

	// Check if any directories match junk patterns
	// We need to track which ones to skip
	toSkip := make(map[string]bool)
	for _, dir := range dirs {
		dirPath := filepath.Join(path, dir)

		// Check if it's a symlink
		fi, err := os.Lstat(dirPath)
		if err != nil || fi.Mode()&os.ModeSymlink != 0 {
			continue
		}

		// Check if directory name matches any junk pattern
		for pattern := range patterns {
			if matched, _ := filepath.Match(pattern, dir); matched {
				size, _ := dirSize(dirPath)
				*junk = append(*junk, junkFolder{dirPath, size})
				toSkip[dir] = true
				break
			}
		}
	}

	// Recursively process non-junk directories
	for _, dir := range dirs {
		if !toSkip[dir] {
			dirPath := filepath.Join(path, dir)
			findJunkRecursive(dirPath, junk)
		}
	}
}

// dirSize returns the total size of all files in a directory (non-recursive for symlinks)
func dirSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() && info.Mode()&os.ModeSymlink == 0 {
			size += info.Size()
		}

		return nil
	})

	return size, err
}

func main() {
	// Manual argument parsing to handle args in any order (like Python argparse)
	root := "."
	var index *int
	deleteFlag := false

	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--delete" {
			deleteFlag = true
		} else if arg[0] != '-' {
			// Positional argument
			if idx, err := strconv.Atoi(arg); err == nil {
				// It's a digit - could be index in current dir
				if index == nil && root == "." {
					// First positional arg that's a digit -> treat as index
					index = &idx
				} else if root == "." {
					// Second numeric arg with root still default -> still index
					index = &idx
				}
			} else {
				// It's not a digit - it's the root path
				if root == "." {
					root = arg
				}
			}
		}
	}

	// If root is a digit string, treat it as index in current dir
	if rootIdx, err := strconv.Atoi(root); err == nil && index == nil {
		index = &rootIdx
		root = "."
	}

	// Find all junk folders
	junk, err := findJunk(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Sort by path for consistent output
	sort.Slice(junk, func(i, j int) bool {
		return junk[i].path < junk[j].path
	})

	var total int64
	for _, j := range junk {
		total += j.size
	}

	// Delete specific index or all
	if index != nil {
		if *index >= 0 && *index < len(junk) {
			j := junk[*index]
			if deleteFlag {
				err := os.RemoveAll(j.path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error deleting %s: %v\n", j.path, err)
					os.Exit(1)
				}
				fmt.Printf("deleted %.1f MB  %s\n", float64(j.size)/1e6, j.path)
			} else {
				fmt.Printf("would free %.1f MB  %s\n", float64(j.size)/1e6, j.path)
			}
		} else {
			if len(junk) == 0 {
				fmt.Println("No junk folders found")
			} else {
				fmt.Printf("Index %d out of range (0-%d)\n", *index, len(junk)-1)
			}
		}
	} else {
		// List all with indices
		for idx, j := range junk {
			fmt.Printf("[%d]  %10.1f MB  %s\n", idx, float64(j.size)/1e6, j.path)
			if deleteFlag {
				err := os.RemoveAll(j.path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error deleting %s: %v\n", j.path, err)
					continue
				}
			}
		}
		action := "would free"
		if deleteFlag {
			action = "deleted"
		}
		fmt.Printf("%s %.1f MB\n", action, float64(total)/1e6)
	}
}

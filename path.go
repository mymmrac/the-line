package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func parsePatterns(patterns []string) (paths []string, err error) {
	var matched []string
	for _, pattern := range patterns {
		matched, err = filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("glob: %w", err)
		}

		paths = append(paths, matched...)
	}

	return paths, nil
}

func parsePaths(paths []string, recursive bool) (files []string, err error) {
	files, dirs, err := splitFilesAndDirs(paths)
	if err != nil {
		return nil, fmt.Errorf("split files & dirs: %w", err)
	}

	if recursive && len(dirs) > 0 {
		for _, dir := range dirs {
			err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return fmt.Errorf("file walk: %w", err)
				}

				if !info.IsDir() {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("walk: %w", err)
			}
		}
	}

	files, err = absPaths(files)
	if err != nil {
		return nil, fmt.Errorf("abs paths: %w", err)
	}

	return files, nil
}

func splitFilesAndDirs(paths []string) (files, dirs []string, err error) {
	var fileInfo os.FileInfo
	for _, path := range paths {
		fileInfo, err = os.Stat(path)
		if err != nil {
			return nil, nil, fmt.Errorf("file info: %w", err)
		}

		if fileInfo.IsDir() {
			dirs = append(dirs, path)
		} else {
			files = append(files, path)
		}
	}

	return files, dirs, nil
}

func absPaths(paths []string) ([]string, error) {
	var err error
	absP := make([]string, len(paths))
	for i, p := range paths {
		absP[i], err = filepath.Abs(p)
		if err != nil {
			return nil, fmt.Errorf("path %q: %w", p, err)
		}
	}

	return absP, nil
}

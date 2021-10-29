package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type filesData struct {
	filePaths    []string
	usedFiles    int
	skippedFiles int
}

func processPatterns(patterns []string, isRecursive, isDotFiles bool) (*filesData, error) {
	paths, err := pathsFromPatterns(patterns)
	if err != nil {
		return nil, fmt.Errorf("paths: %w", err)
	}

	filePaths, err := filesFromPaths(paths, isRecursive, isDotFiles)
	if err != nil {
		return nil, fmt.Errorf("file info: %w", err)
	}

	textFilePaths, err := textFilesFromFiles(filePaths)
	if err != nil {
		return nil, fmt.Errorf("filter text files: %w", err)
	}

	return &filesData{
		filePaths:    textFilePaths,
		usedFiles:    len(textFilePaths),
		skippedFiles: len(filePaths) - len(textFilePaths),
	}, nil
}

func pathsFromPatterns(patterns []string) (paths []string, err error) {
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

func filesFromPaths(paths []string, recursive, dotFiles bool) (files []string, err error) {
	files, dirs, err := filesAndDirsFromPaths(paths)
	if err != nil {
		return nil, fmt.Errorf("split files & dirs: %w", err)
	}

	if recursive && len(dirs) > 0 {
		for _, dir := range dirs {
			err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
				if !dotFiles && isDotPath(path) {
					return nil
				}

				if err != nil {
					return fmt.Errorf("file walk: {%s} %w", dirs, err)
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

	files, err = absPathsFromPaths(files)
	if err != nil {
		return nil, fmt.Errorf("abs paths: %w", err)
	}

	if !dotFiles {
		files = excludeDotPaths(files)
	}

	return files, nil
}

func filesAndDirsFromPaths(paths []string) (files, dirs []string, err error) {
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

func absPathsFromPaths(paths []string) (absPaths []string, err error) {
	absPaths = make([]string, len(paths))

	for i, p := range paths {
		absPaths[i], err = filepath.Abs(p)
		if err != nil {
			return nil, fmt.Errorf("path %q: %w", p, err)
		}
	}

	return absPaths, nil
}

const dotPath = "/."

func isDotPath(path string) bool {
	return strings.Contains(path, dotPath)
}

func excludeDotPaths(paths []string) []string {
	var excluded []string

	for _, path := range paths {
		if !isDotPath(path) {
			excluded = append(excluded, path)
		}
	}

	return excluded
}

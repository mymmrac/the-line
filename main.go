package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/godoc/util"
)

func main() {
	recursive := flag.Bool("r", false, "Recursively search in dirs matched by pattern")
	flag.Parse()

	if len(flag.Args()) < 1 {
		exitWithPrint("Please enter pattern(s) as arguments")
	}
	patterns := flag.Args()

	var paths []string
	var err error
	for _, pattern := range patterns {
		onePatternPaths, err := filepath.Glob(pattern)
		if err != nil {
			exitWithPrint("Pattern:", err)
		}
		paths = append(paths, onePatternPaths...)
	}

	filePaths, err := parsePaths(paths, *recursive)
	if err != nil {
		exitWithPrint("Parse file info:", err)
	}

	textFilesPaths, err := filterTextFiles(filePaths)
	if err != nil {
		exitWithPrint("Failed to filter text files:", err)
	}

	fmt.Println("File count:", len(textFilesPaths))
	fmt.Println("Skipped files:", len(filePaths)-len(textFilesPaths))

	for _, path := range textFilesPaths {
		lc, err := parseFile(path)
		if err != nil {
			exitWithPrint("Reading file:", err)
		}

		fmt.Println(lc.Any, lc.Blank)
	}
}

type lineCount struct {
	Any   int
	Blank int
}

func parseFile(path string) (*lineCount, error) {
	//nolint:gosec
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	lc := &lineCount{}
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()

		lc.Any++
		if strings.TrimSpace(line) == "" {
			lc.Blank++
		}
	}
	if err = sc.Err(); err != nil {
		return nil, err
	}

	return lc, nil
}

func filterTextFiles(files []string) (textFiles []string, err error) {
	for _, file := range files {
		isText, err := isTextFile(file)
		if err != nil {
			return nil, err
		}

		if isText {
			textFiles = append(textFiles, file)
		}
	}
	return textFiles, nil
}

func isTextFile(filename string) (bool, error) {
	//nolint:gosec
	f, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			exitWithPrint("Fail to close file:", closeErr)
		}
	}()

	var buf [1024]byte
	n, err := f.Read(buf[0:])
	if err != nil {
		return false, err
	}

	return util.IsText(buf[0:n]), nil
}

func parsePaths(paths []string, recursive bool) (files []string, err error) {
	files, dirs, err := splitFilesAndDirs(paths)
	if err != nil {
		return nil, err
	}

	if recursive && len(dirs) > 0 {
		for _, dir := range dirs {
			err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return files, nil
}

func splitFilesAndDirs(paths []string) (files, dirs []string, err error) {
	for _, path := range paths {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return nil, nil, err
		}

		if fileInfo.IsDir() {
			dirs = append(dirs, path)
		} else {
			files = append(files, path)
		}
	}

	return files, dirs, nil
}

func exitWithPrint(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}

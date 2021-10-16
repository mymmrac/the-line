package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	recursive := flag.Bool("r", false, "Recursively search in dirs matched by pattern")
	flag.Parse()

	if len(flag.Args()) < 1 {
		exitWithPrint("Please enter pattern(s) as arguments")
	}
	patterns := flag.Args()

	paths, err := parsePatterns(patterns)
	if err != nil {
		exitWithPrint("Pattern(s):", err)
	}

	filePaths, err := parsePaths(paths, *recursive)
	if err != nil {
		exitWithPrint("Parse file info:", err)
	}

	textFilePaths, err := filterTextFiles(filePaths)
	if err != nil {
		exitWithPrint("Failed to filter text files:", err)
	}

	usedFiles := len(textFilePaths)
	skippedFiles := len(filePaths) - usedFiles

	fmt.Println("Used files:", usedFiles)
	fmt.Println("Skipped files:", skippedFiles)
	fmt.Println()

	counters, err := countLines(textFilePaths)
	if err != nil {
		exitWithPrint("Counting lines:", err)
	}

	for _, counter := range counters {
		fmt.Println(counter.filename)
		fmt.Println(counter.matched)
		fmt.Println()
	}

	// TODO: Add sum of all
}

func exitWithPrint(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}

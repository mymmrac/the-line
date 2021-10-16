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

func exitWithPrint(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}

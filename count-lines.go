package main

import (
	"bufio"
	"fmt"
	"os"
)

func countLinesInFiles(files []string, profs profiles) (lineCounters, error) {
	lcs := make(lineCounters, len(files))
	for i, path := range files {
		lc, err := countLineInFile(path, profs)
		if err != nil {
			return nil, fmt.Errorf("parsing file: %w", err)
		}

		lcs[i] = *lc
	}

	return lcs, nil
}

func countLineInFile(path string, profs profiles) (*lineCounter, error) {
	//nolint:gosec
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	lc := newLineCounter(path)

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()

		lc.countLine(line, profs)
	}

	if err = sc.Err(); err != nil {
		return nil, fmt.Errorf("scan file: %w", err)
	}

	return lc, nil
}
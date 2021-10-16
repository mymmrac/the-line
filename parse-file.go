package main

import (
	"bufio"
	"fmt"
	"os"
)

func parseFile(path string, profs profiles) (*lineCounter, error) {
	//nolint:gosec
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	lc := newLineCounter(path, profs)

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()

		lc.count(line)
	}
	if err = sc.Err(); err != nil {
		return nil, fmt.Errorf("scan file: %w", err)
	}

	return lc, nil
}

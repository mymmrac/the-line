package main

import "fmt"

type lineCounter struct {
	filename string
	rules    []*rule
	matched  map[string]int
}

type lineCounters []lineCounter

func newLineCounter(filename string, rules []*rule) *lineCounter {
	return &lineCounter{
		filename: filename,
		rules:    rules,
		matched:  make(map[string]int),
	}
}

func (l *lineCounter) count(line string) {
	for _, r := range l.rules {
		if r.checkPath(l.filename) && r.checkLine(line) {
			l.matched[r.name]++
		}
	}
}

func countLines(files []string) (lineCounters, error) {
	lcs := make(lineCounters, len(files))
	for i, path := range files {
		lc, err := parseFile(path)
		if err != nil {
			return nil, fmt.Errorf("parsing file: %w", err)
		}

		lcs[i] = *lc
	}

	return lcs, nil
}

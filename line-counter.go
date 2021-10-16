package main

import "fmt"

type lineCounter struct {
	filename string
	profiles profiles

	matched map[string]map[string]int
}

type lineCounters []lineCounter

func newLineCounter(filename string, profiles profiles) *lineCounter {
	return &lineCounter{
		filename: filename,
		profiles: profiles,

		matched: make(map[string]map[string]int),
	}
}

func (l *lineCounter) count(line string) {
	for profName, p := range l.profiles {
		if p.checkPath(l.filename) {
			for ruleName, r := range p.Rules {
				if r.checkPath(l.filename) && r.checkLine(line) {
					if _, ok := l.matched[profName]; !ok {
						l.matched[profName] = make(map[string]int)
					}

					l.matched[profName][ruleName]++
				}
			}
		}
	}
}

func countLines(files []string, profs profiles) (lineCounters, error) {
	lcs := make(lineCounters, len(files))
	for i, path := range files {
		lc, err := parseFile(path, profs)
		if err != nil {
			return nil, fmt.Errorf("parsing file: %w", err)
		}

		lcs[i] = *lc
	}

	return lcs, nil
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type anyFilter struct{}

func (f anyFilter) filter(_ string) bool {
	return true
}

type blankFilter struct{}

func (f blankFilter) filter(line string) bool {
	return line == ""
}

type trimSpaceModifier struct{}

func (t trimSpaceModifier) modify(line string) string {
	return strings.TrimSpace(line)
}

func parseFile(path string) (*lineCounter, error) {
	//nolint:gosec
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	anyRule, _ := newRule("any", ".+", nil, []filterer{anyFilter{}})
	blankRule, _ := newRule("blank", ".+", nil, []filterer{blankFilter{}})
	blankTrimmedRule, _ := newRule("blank trimmed", ".+",
		[]modifier{trimSpaceModifier{}}, []filterer{blankFilter{}})

	lc := newLineCounter(path, []*rule{anyRule, blankRule, blankTrimmedRule})

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

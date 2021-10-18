package main

type countByRule map[string]int

type countByProfile map[string]countByRule

type lineCounter struct {
	filename string
	count    countByProfile
}

func newLineCounter(filename string) *lineCounter {
	return &lineCounter{
		filename: filename,
		count:    make(countByProfile),
	}
}

func (l *lineCounter) countLine(line string, profiles profiles) {
	for profName, p := range profiles {
		if p.checkPath(l.filename) {
			if _, ok := l.count[profName]; !ok {
				l.count[profName] = make(countByRule)
			}

			for ruleName, r := range p.Rules {
				if _, ok := l.count[profName][ruleName]; !ok {
					l.count[profName][ruleName] = 0
				}

				if r.checkPath(l.filename) && r.checkLine(line) {
					l.count[profName][ruleName]++
				}
			}
		}
	}
}

type lineCounters []lineCounter

func (l lineCounters) totalCount() countByProfile {
	count := make(countByProfile, len(l))

	for _, lc := range l {
		for profileName, profileCount := range lc.count {
			if _, ok := count[profileName]; !ok {
				count[profileName] = make(countByRule, len(profileCount))
			}

			for ruleName, ruleCount := range profileCount {
				count[profileName][ruleName] += ruleCount
			}
		}
	}

	return count
}

package main

type lineCounter struct {
	filename string
	count    map[string]map[string]int // TODO: Crate type alias
}

func newLineCounter(filename string) *lineCounter {
	return &lineCounter{
		filename: filename,
		count:    make(map[string]map[string]int),
	}
}

func (l *lineCounter) countLine(line string, profiles profiles) {
	for profName, p := range profiles {
		if p.checkPath(l.filename) {
			if _, ok := l.count[profName]; !ok {
				l.count[profName] = make(map[string]int)
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

func (l lineCounters) totalCount() map[string]map[string]int {
	count := make(map[string]map[string]int, len(l))

	for _, lc := range l {
		for profileName, profileCount := range lc.count {
			if _, ok := count[profileName]; !ok {
				count[profileName] = make(map[string]int, len(profileCount))
			}

			for ruleName, ruleCount := range profileCount {
				count[profileName][ruleName] += ruleCount
			}
		}
	}

	return count
}

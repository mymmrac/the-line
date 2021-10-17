package main

type lineCounter struct {
	filename string
	profiles profiles

	matched map[string]map[string]int
}

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
			if _, ok := l.matched[profName]; !ok {
				l.matched[profName] = make(map[string]int)
			}
			for ruleName, r := range p.Rules {
				if _, ok := l.matched[profName][ruleName]; !ok {
					l.matched[profName][ruleName] = 0
				}

				if r.checkPath(l.filename) && r.checkLine(line) {
					l.matched[profName][ruleName]++
				}
			}
		}
	}
}

type lineCounters []lineCounter

func (l lineCounters) sum() map[string]map[string]int {
	matched := make(map[string]map[string]int)
	for _, lc := range l {
		for profileName, p := range lc.matched {
			if _, ok := matched[profileName]; !ok {
				matched[profileName] = make(map[string]int)
			}

			for ruleName, count := range p {
				matched[profileName][ruleName] += count
			}
		}
	}

	return matched
}

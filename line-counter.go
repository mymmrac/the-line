package main

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

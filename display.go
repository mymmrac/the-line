package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func displayCounts(usedFiles, skippedFiles int, counters lineCounters) {
	fmt.Println("Used files:", usedFiles)
	fmt.Println("Skipped files:", skippedFiles)
	fmt.Println()

	for _, counter := range counters {
		fmt.Printf(" ==== %s ==== \n", counter.filename)
		displayCount(counter.count)
		fmt.Println()
	}

	fmt.Println(" ==== TOTAL ==== ")
	displayCount(counters.totalCount())
}

type ruleSorter struct {
	rulesNames []string
	counts     []string
}

func (r ruleSorter) Len() int {
	return len(r.rulesNames)
}

func (r ruleSorter) Less(i, j int) bool {
	return r.rulesNames[i] < r.rulesNames[j]
}

func (r *ruleSorter) Swap(i, j int) {
	r.rulesNames[i], r.rulesNames[j] = r.rulesNames[j], r.rulesNames[i]
	r.counts[i], r.counts[j] = r.counts[j], r.counts[i]
}

type namedProfile struct {
	name      string
	ruleMatch map[string]int
}

const minRuleNameLength = 6

func displayCount(matched map[string]map[string]int) {
	i := 0
	np := make([]namedProfile, len(matched))
	for profName, ruleMatch := range matched {
		np[i] = namedProfile{
			name:      profName,
			ruleMatch: ruleMatch,
		}
		i++
	}

	sort.Slice(np, func(i, j int) bool {
		return np[i].name < np[j].name
	})

	for j := range np {
		profName, ruleMatch := np[j].name, np[j].ruleMatch
		k := 0
		rulesList := make([]string, len(ruleMatch))
		countList := make([]string, len(ruleMatch))
		for ruleName, count := range ruleMatch {
			rulesList[k] = ruleName
			countList[k] = strconv.Itoa(count)
			k++
		}

		rs := ruleSorter{
			rulesNames: rulesList,
			counts:     countList,
		}

		sort.Sort(&rs)

		profRow := profName + ": "
		rulesRow := strings.Repeat(" ", len(profRow))

		for l, rn := range rs.rulesNames {
			r := rn + "  "
			if len(r) < minRuleNameLength {
				r += strings.Repeat(" ", minRuleNameLength-len(r))
			}

			rulesRow += r
			profRow += rs.counts[l] + strings.Repeat(" ", len(r)-len(rs.counts[l]))
		}

		fmt.Println(rulesRow)
		fmt.Println(profRow)
	}
}

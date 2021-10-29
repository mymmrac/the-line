package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	bracesColor  = lipgloss.Color("#185ADB")
	successColor = lipgloss.Color("#6ECB63")
	errorColor   = lipgloss.Color("#FF5C58")
	numbersColor = lipgloss.Color("#3EDBF0")
)

var (
	checkStr = lipgloss.NewStyle().SetString("✓").Foreground(successColor).String()
	crossStr = lipgloss.NewStyle().SetString("✗").Foreground(errorColor).String()

	bracesStyle   = lipgloss.NewStyle().Foreground(bracesColor)
	leftBraceStr  = bracesStyle.SetString("❮").String()
	rightBraceStr = bracesStyle.SetString("❯").String()

	numbersStyle = lipgloss.NewStyle().Foreground(numbersColor)

	titleStyle = lipgloss.NewStyle().Underline(true)
)

func displayCounts(usedFiles, skippedFiles int, counters lineCounters, verbose bool) string {
	res := &bytes.Buffer{}

	fmt.Fprintln(res, "Files:")
	fmt.Fprintln(res, checkStr+" Used   ", numbersStyle.SetString(strconv.Itoa(usedFiles)))
	fmt.Fprintln(res, crossStr+" Skipped", numbersStyle.SetString(strconv.Itoa(skippedFiles)))
	fmt.Fprintln(res)

	for _, counter := range counters {
		if counter.tooLongLine {
			res.WriteString(fmt.Sprintf("%s: File %q has too long line (64K line limit)\n\n", lipgloss.NewStyle().Foreground(errorColor).SetString("WARN"), counter.filename))
		}
		if verbose {
			displayBlockCount(counter.filename, counter.count, res)
			fmt.Fprintln(res)
		}
	}

	displayBlockCount("TOTAL", counters.totalCount(), res)

	return res.String()
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func displayBlockCount(title string, count countByProfile, res *bytes.Buffer) {
	maxProfLen := -1
	for profName := range count {
		maxProfLen = max(maxProfLen, len(profName))
	}

	resTotal := &bytes.Buffer{}
	l := displayCount(count, maxProfLen, resTotal)

	title = fmt.Sprintf(" %s %s %s ", leftBraceStr, titleStyle.SetString(title), rightBraceStr)
	if l > lipgloss.Width(title) {
		title = lipgloss.NewStyle().Width(l).Align(lipgloss.Center).SetString(title).String()
	}

	fmt.Fprintln(res, title)
	res.Write(resTotal.Bytes())
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

func displayCount(matched countByProfile, maxProfLen int, res *bytes.Buffer) int {
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

	maxLen := -1

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

		profRow := profName + " "
		if len(profRow) < maxProfLen {
			profRow += strings.Repeat(" ", maxProfLen-len(profRow)+1)
		}
		rulesRow := strings.Repeat(" ", len(profRow))

		for l, rn := range rs.rulesNames {
			r := rn + "  "
			if len(r) < minRuleNameLength {
				r += strings.Repeat(" ", minRuleNameLength-len(r))
			}

			rulesRow += r
			profRow += numbersStyle.SetString(rs.counts[l]).String() + strings.Repeat(" ", len(r)-len(rs.counts[l]))
		}

		fmt.Fprintln(res, rulesRow)
		fmt.Fprintln(res, profRow)

		maxLen = max(maxLen, len(rulesRow))
	}

	return maxLen
}

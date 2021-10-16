package main

import (
	"regexp"
	"strings"
)

type filterer interface {
	filter(line string) bool
}

type anyFilter struct{}

func (f anyFilter) filter(_ string) bool {
	return true
}

type blankFilter struct{}

func (f blankFilter) filter(line string) bool {
	return line == ""
}

type matchFilter struct {
	Line string
}

func (f matchFilter) filter(line string) bool {
	return line == f.Line
}

type containsFilter struct {
	Line string
}

func (f containsFilter) filter(line string) bool {
	return strings.Contains(line, f.Line)
}

type lengthFilter struct {
	Length int
}

func (f lengthFilter) filter(line string) bool {
	return len(line) == f.Length
}

type regexpFilter struct {
	Reg *regexp.Regexp
}

func (f regexpFilter) filter(line string) bool {
	return f.Reg.MatchString(line)
}

type multiLineFilter struct {
	startFilter filterer
	endFilter   filterer
	matching    bool
}

func (f *multiLineFilter) filter(line string) bool {
	if f.startFilter.filter(line) {
		f.matching = true
		return true
	}

	if f.endFilter.filter(line) {
		f.matching = false
		return true
	}

	return f.matching
}

type unionFilter struct {
	filterA filterer
	filterB filterer
}

func (f unionFilter) filter(line string) bool {
	return f.filterA.filter(line) || f.filterB.filter(line)
}

type intersectionFilter struct {
	filterA filterer
	filterB filterer
}

func (f intersectionFilter) filter(line string) bool {
	return f.filterA.filter(line) && f.filterB.filter(line)
}

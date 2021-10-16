package main

import "strings"

type lineCounter struct {
	Any   int
	Blank int
}

func (l lineCounter) count(line string) {
	l.Any++

	if strings.TrimSpace(line) == "" {
		l.Blank++
	}
}

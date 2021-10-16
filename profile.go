package main

import (
	"fmt"
	"regexp"
)

type profile struct {
	PathFormat *regexp.Regexp `json:"path_format"`
	Rules      rules          `json:"rules"`
}

func (p *profile) checkPath(filename string) bool {
	return p.PathFormat.MatchString(filename)
}

type profiles map[string]profile

func (p profiles) filter(profileNames []string) (profiles, error) {
	filtered := make(profiles)
	for _, profName := range profileNames {
		prof, ok := p[profName]
		if !ok {
			return nil, fmt.Errorf("profile not found: %q", profName)
		}

		filtered[profName] = prof
	}

	return filtered, nil
}

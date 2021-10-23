package main

import (
	"fmt"
	"regexp"
)

type profile struct {
	PathFormat *regexp.Regexp `yaml:"path-format"`
	Rules      rules          `yaml:"rules"`
}

func (p *profile) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type uProfile struct {
		PathFormat string `yaml:"path-format"`
		Rules      rules  `yaml:"rules"`
	}

	u := uProfile{}
	err := unmarshal(&u)
	if err != nil {
		return err
	}

	p.Rules = u.Rules
	p.PathFormat, err = regexp.Compile(u.PathFormat)

	return err
}

func (p *profile) checkPath(filename string) bool {
	return p.PathFormat.MatchString(filename)
}

type profiles map[string]profile

func filterProfiles(p profiles, profileNames []string) (profiles, error) {
	filtered := make(profiles, len(profileNames))

	for _, profName := range profileNames {
		prof, ok := p[profName]
		if !ok {
			return nil, fmt.Errorf("profile not found: %q", profName)
		}

		filtered[profName] = prof
	}

	return filtered, nil
}

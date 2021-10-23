package main

import (
	_ "embed"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

//go:embed default.yaml
var defaultConfigBytes []byte

type config struct {
	Profiles profiles `yaml:"profiles"`
}

func embeddedConfig() (*config, error) {
	var conf config
	err := yaml.Unmarshal(defaultConfigBytes, &conf)
	if err != nil {
		return nil, fmt.Errorf("decoding: %w", err)
	}

	return &conf, nil
}

func userConfig(configFilename string) (*config, error) {
	//nolint:gosec
	configFile, err := os.Open(configFilename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	var conf config
	err = yaml.NewDecoder(configFile).Decode(&conf)
	if err != nil {
		return nil, fmt.Errorf("decoding: %w", err)
	}

	return &conf, nil
}

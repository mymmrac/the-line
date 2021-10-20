package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type config struct {
	Profiles profiles `yaml:"profiles"`
}

const defaultProfile = "default"

const separator = ","

//go:embed default.yaml
var defaultConfigBytes []byte

func main() {
	// Reading embedded config
	conf, err := embeddedConfig()
	if err != nil {
		exitWithPrint("Embedded config:", err)
	}

	// Reading user input
	args, err := userInput()
	if err != nil {
		exitWithPrint("User input:", err)
	}

	// Reading user config
	if args.configFilename != "" {
		conf, err = userConfig(args.configFilename)
		if err != nil {
			exitWithPrint("User config:", err)
		}
	}

	// Processing configs
	// TODO: Refactor
	profs, err := conf.Profiles.filter(args.profileNames)
	if err != nil {
		exitWithPrint("Getting profiles:", err)
	}

	// Processing files
	data, err := processPatterns(args.patterns, args.isRecursive)
	if err != nil {
		exitWithPrint("Processing files:", err)
	}

	// Counting lines
	counters, err := countLinesInFiles(data.filePaths, profs)
	if err != nil {
		exitWithPrint("Counting lines:", err)
	}

	// Displaying output
	displayCounts(data.usedFiles, data.skippedFiles, counters)
}

func embeddedConfig() (*config, error) {
	var conf config
	err := yaml.Unmarshal(defaultConfigBytes, &conf)
	if err != nil {
		return nil, fmt.Errorf("decoding: %w", err)
	}

	return &conf, nil
}

type userArgs struct {
	isRecursive    bool
	profileNames   []string
	configFilename string
	patterns       []string
}

func userInput() (userArgs, error) {
	// TODO: Add arguments validation
	// TODO: Flag for include hidden files

	recursiveFlag := flag.Bool("r", false, "Recursively search in dirs matched by pattern")
	profNamesFlag := flag.String("p", "", "Profiles to use")
	configFileFlag := flag.String("c", "", "User defined config file")

	flag.Parse()

	args := userArgs{
		isRecursive:    *recursiveFlag,
		configFilename: *configFileFlag,
	}

	args.profileNames = []string{
		defaultProfile,
	}

	if len(*profNamesFlag) > 0 {
		args.profileNames = strings.Split(*profNamesFlag, separator)
	}

	if len(flag.Args()) < 1 {
		return userArgs{}, fmt.Errorf("no patterns given")
	}
	args.patterns = flag.Args()

	return args, nil
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

type filesData struct {
	filePaths    []string
	usedFiles    int
	skippedFiles int
}

func processPatterns(patterns []string, isRecursive bool) (filesData, error) {
	paths, err := pathsFromPatterns(patterns)
	if err != nil {
		return filesData{}, fmt.Errorf("paths: %w", err)
	}

	filePaths, err := filesFromPaths(paths, isRecursive)
	if err != nil {
		return filesData{}, fmt.Errorf("file info: %w", err)
	}

	textFilePaths, err := textFilesFromFiles(filePaths)
	if err != nil {
		return filesData{}, fmt.Errorf("filter text files: %w", err)
	}

	return filesData{
		filePaths:    textFilePaths,
		usedFiles:    len(textFilePaths),
		skippedFiles: len(filePaths) - len(textFilePaths),
	}, nil
}

func exitWithPrint(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}

package main

import (
	"flag"
	"fmt"
	"strings"
)

const defaultProfile = "default"

const separator = ","

type userArgs struct {
	isRecursive    bool
	isDotFiles     bool
	profileNames   []string
	configFilename string
	patterns       []string
}

func userInput() (*userArgs, error) {
	// TODO: Add arguments validation

	recursiveFlag := flag.Bool("r", false, "Recursively search in dirs matched by pattern")
	dotFilesFlag := flag.Bool("d", false, "Include dot files/folders")
	profNamesFlag := flag.String("p", "", "Profiles to use")
	configFileFlag := flag.String("c", "", "User defined config file")

	flag.Parse()

	args := &userArgs{
		isRecursive:    *recursiveFlag,
		isDotFiles:     *dotFilesFlag,
		configFilename: *configFileFlag,
	}

	args.profileNames = []string{
		defaultProfile,
	}

	if len(*profNamesFlag) > 0 {
		args.profileNames = strings.Split(*profNamesFlag, separator)
	}

	if len(flag.Args()) < 1 {
		return nil, fmt.Errorf("no patterns given")
	}
	args.patterns = flag.Args()

	return args, nil
}

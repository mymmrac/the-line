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
	conf := &config{}
	err := yaml.Unmarshal(defaultConfigBytes, conf)
	if err != nil {
		exitWithPrint("Embedded config:", err)
	}

	// Reading arguments
	// TODO: Add arguments validation

	recursiveFlag := flag.Bool("r", false, "Recursively search in dirs matched by pattern")
	profNamesFlag := flag.String("p", "", "Profiles to use")
	configFileFlag := flag.String("c", "", "User defined config file")

	flag.Parse()

	profileNames := []string{
		defaultProfile,
	}

	if len(*profNamesFlag) > 0 {
		profileNames = strings.Split(*profNamesFlag, separator)
	}

	if len(flag.Args()) < 1 {
		exitWithPrint("Please enter pattern(s) as arguments")
	}
	patterns := flag.Args()

	// Reading user config
	if *configFileFlag != "" {
		var configFile *os.File
		configFile, err = os.Open(*configFileFlag)
		if err != nil {
			exitWithPrint("Reading user config:", err)
		}

		err = yaml.NewDecoder(configFile).Decode(conf)
		if err != nil {
			exitWithPrint("User config:", err)
		}
	}

	// Processing configs
	profs, err := conf.Profiles.filter(profileNames)
	if err != nil {
		exitWithPrint("Getting profiles:", err)
	}

	// Processing files
	paths, err := pathsFromPatterns(patterns)
	if err != nil {
		exitWithPrint("Pattern(s):", err)
	}

	filePaths, err := filesFromPaths(paths, *recursiveFlag)
	if err != nil {
		exitWithPrint("Parse file info:", err)
	}

	textFilePaths, err := textFilesFromFiles(filePaths)
	if err != nil {
		exitWithPrint("Failed to filter text files:", err)
	}

	usedFiles := len(textFilePaths)
	skippedFiles := len(filePaths) - usedFiles

	counters, err := countLinesInFiles(textFilePaths, profs)
	if err != nil {
		exitWithPrint("Counting lines:", err)
	}

	// Displaying output
	displayCounts(usedFiles, skippedFiles, counters)
}

func exitWithPrint(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}

/*

Desired output:

<filename #>
	- 		 	<rule #>	<rule #> ...
<profile #> 	  123		  123
	- 			<rule #>	<rule #> 	<rule #> ...
<profile #> 	  123		  123		  123
	- 			<rule #>	<rule #> ...
<profile #> 	  123		  123

<filename #>
	- 		 	<rule #>	<rule #> ...
<profile #> 	  123		  123
	- 		 	<rule #> ...
<profile #> 	  123

========================================================

Total
	- 		 	<rule #>	<rule #> ...
<profile #> 	  123		  123
	- 			<rule #>	<rule #> 	<rule #> ...
<profile #> 	  123		  123		  123
	- 			<rule #>	<rule #> ...
<profile #> 	  123		  123
	- 		 	<rule #> ...
<profile #> 	  123

...

*/

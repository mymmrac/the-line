package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type config struct {
	Profiles profiles `yaml:"profiles"`
}

const defaultProfile = "Default"

const separator = ","

func main() {
	// Reading configs
	configFile, _ := os.Open("default.yaml")

	conf := &config{}
	err := yaml.NewDecoder(configFile).Decode(conf)
	if err != nil {
		exitWithPrint("Config:", err)
	}

	// Reading arguments
	recursive := flag.Bool("r", false, "Recursively search in dirs matched by pattern")
	profNameFlag := flag.String("p", "", "Profiles to use")

	flag.Parse()

	profileNames := []string{
		defaultProfile,
	}

	if len(*profNameFlag) > 0 {
		profileNames = strings.Split(*profNameFlag, separator)
	}

	profs, err := conf.Profiles.filter(profileNames)
	if err != nil {
		exitWithPrint("Getting profiles:", err)
	}

	if len(flag.Args()) < 1 {
		exitWithPrint("Please enter pattern(s) as arguments")
	}
	patterns := flag.Args()

	// Processing files
	paths, err := parsePatterns(patterns)
	if err != nil {
		exitWithPrint("Pattern(s):", err)
	}

	filePaths, err := parsePaths(paths, *recursive)
	if err != nil {
		exitWithPrint("Parse file info:", err)
	}

	textFilePaths, err := filterTextFiles(filePaths)
	if err != nil {
		exitWithPrint("Failed to filter text files:", err)
	}

	usedFiles := len(textFilePaths)
	skippedFiles := len(filePaths) - usedFiles

	counters, err := countLines(textFilePaths, profs)
	if err != nil {
		exitWithPrint("Counting lines:", err)
	}

	// Displaying output
	displayCounts(usedFiles, skippedFiles, counters)
}

func displayCounts(usedFiles, skippedFiles int, counters lineCounters) {
	fmt.Println("Used files:", usedFiles)
	fmt.Println("Skipped files:", skippedFiles)
	fmt.Println()

	for _, counter := range counters {
		fmt.Println(counter.filename)
		for profName, ruleMatch := range counter.matched {
			rulesRow := "          "
			profRow := profName + ": "
			for ruleName, count := range ruleMatch {
				rulesRow += ruleName + "   "
				profRow += fmt.Sprintf(" %d ", count)
			}

			fmt.Println(rulesRow)
			fmt.Println(profRow)
		}
		fmt.Println()
	}

	// TODO: Add sum of all
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

package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const defaultProfile = "default"

type userArgs struct {
	recursive      bool
	dotFiles       bool
	verbose        bool
	configFilename string
	profileNames   []string
	patterns       []string
}

func userInput() (*userArgs, bool, error) {
	// TODO: Add arguments validation
	var (
		isExecuting = false

		cmdArgs []string

		recursiveFlag  bool
		dotFilesFlag   bool
		verboseFlag    bool
		configFileFlag string
		profNames      []string
	)

	rootCmd := &cobra.Command{
		Use:   "{patterns}...",
		Short: "Count lines in matched files",
		Long: `Counts lines in files matched by glob patterns.
Displays count information based in profiles & rules specified.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			isExecuting = true
			cmdArgs = args
		},
		Version: "v1.0.2",
	}

	rootCmd.SetVersionTemplate(`{{printf "The-Line. Version %s" .Version}}
`)

	rootCmd.Flags().BoolVarP(&recursiveFlag, "recursive", "r", false, "Recursively search in dirs matched by pattern")
	rootCmd.Flags().BoolVarP(&dotFilesFlag, "dot-files", "d", false, "Include dot files/folders")
	rootCmd.Flags().StringVarP(&configFileFlag, "config", "c", "", "User defined config file")
	rootCmd.Flags().StringSliceVarP(&profNames, "profiles", "p", []string{defaultProfile}, "Profiles to use")
	rootCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Verbose output")

	err := rootCmd.Execute()
	if err != nil {
		return nil, false, fmt.Errorf("cmd: %w", err)
	}

	args := &userArgs{
		recursive:      recursiveFlag,
		dotFiles:       dotFilesFlag,
		verbose:        verboseFlag,
		configFilename: configFileFlag,
		profileNames:   profNames,
		patterns:       cmdArgs,
	}

	return args, isExecuting, nil
}

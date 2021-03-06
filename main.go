package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type model struct {
	err        error
	countReady bool
	data       *filesData
	counters   lineCounters
	verbose    bool
}

type parameters struct {
	profs profiles
	args  *userArgs
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		p := parameters{}

		// Reading embedded config
		conf, err := embeddedConfig()
		if err != nil {
			return fmt.Errorf("embedded config: %w", err)
		}

		// Reading user input
		isExecuting := false
		p.args, isExecuting, err = userInput()
		if !isExecuting {
			return tea.Quit()
		}
		if err != nil {
			return fmt.Errorf("user input: %w", err)
		}

		// Reading user config
		if p.args.configFilename != "" {
			conf, err = userConfig(p.args.configFilename)
			if err != nil {
				return fmt.Errorf("user config: %w", err)
			}
		}

		//log.Println(conf)

		// Filter profiles
		p.profs, err = filterProfiles(conf.Profiles, p.args.profileNames)
		if err != nil {
			return fmt.Errorf("getting profiles: %w", err)
		}

		return p
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//log.Println(reflect.TypeOf(msg), msg)

	switch msg := msg.(type) {
	case error:
		m.err = msg
		return m, tea.Quit
	case parameters:
		var err error
		// Processing files
		m.data, err = processPatterns(msg.args.patterns, msg.args.recursive, msg.args.dotFiles)
		if err != nil {
			m.err = fmt.Errorf("processing files: %w", err)
			return m, tea.Quit
		}

		// Counting lines
		m.counters, err = countLinesInFiles(m.data.filePaths, msg.profs)
		if err != nil {
			m.err = fmt.Errorf("counting lines: %w", err)
			return m, tea.Quit
		}

		m.verbose = msg.args.verbose
		m.countReady = true
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Something went wrong: %s\n", m.err)
	}

	if m.countReady {
		// Displaying output
		return displayCounts(m.data.usedFiles, m.data.skippedFiles, m.counters, m.verbose)
	}

	return ""
}

func main() {
	// Logging setup
	//f, err := tea.LogToFile("debug.log", "debug")
	//if err != nil {
	//	fmt.Println("Fatal:", err)
	//	os.Exit(1)
	//}
	//_ = f.Truncate(0)

	// Program
	if err := tea.NewProgram(&model{}).Start(); err != nil {
		fmt.Printf("Opss, something went really wrong: %v", err)
		os.Exit(1)
	}

	// Closing log file
	//_ = f.Close()
}

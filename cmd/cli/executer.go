package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/theckman/yacspin"
)

type Executer struct {
	parser   Parser
	lockfile Lockfile
	spinner  *yacspin.Spinner
}

var spinnerCfg = yacspin.Config{
	Frequency:         100 * time.Millisecond,
	Colors:            []string{"fgYellow"},
	CharSet:           yacspin.CharSets[11],
	Suffix:            " ",
	SuffixAutoColon:   true,
	Message:           "Running commands",
	StopCharacter:     "✓",
	StopColors:        []string{"fgGreen"},
	StopMessage:       "Done",
	StopFailCharacter: "✗",
	StopFailColors:    []string{"fgRed"},
	StopFailMessage:   "Failed",
}

// Executer constructor.
func NewExecuter(p Parser, l Lockfile) Executer {
	spinner, _ := yacspin.New(spinnerCfg)

	return Executer{
		parser:   p,
		lockfile: l,
		spinner:  spinner,
	}
}

// Executes all command strings under given taskName.
// Each call happens in its own go routine.
func (e *Executer) Execute(taskName string, initialRun bool) {
	if _, ok := e.parser.Tasks[taskName]; !ok {
		fmt.Printf("command '%s' not found\n", taskName)
		os.Exit(1)
	}

	task := e.parser.Tasks[taskName]

	if initialRun {
		e.spinner.Start()
	}

	if !initialRun || e.shouldDispatch(task) {
		e.dispatchCommands(task, initialRun)
	} else {
		e.spinner.StopMessage("Nothing to run")
	}

	if initialRun {
		e.spinner.Stop()
	}
}

// Checks whether files have changed since the last run.
// Also updates the lockfile if files did get modified.
// If no "files" key is present in the task, simply returns true.
func (e *Executer) shouldDispatch(task Task) bool {
	if len(task.Files) == 0 {
		return true
	}

	dispatchCh := make(chan bool)
	go e.shouldDispatchRoutine(task, dispatchCh)
	dispatch := <-dispatchCh

	if dispatch {
		e.lockfile.UpdateTimestampsForFiles(task.Files)
	}

	return dispatch
}

// Go Routine function that determines whether the stored
// mtime is greater  than mtime if the file at this moment.
func (e *Executer) shouldDispatchRoutine(task Task, ch chan bool) {
	lockedModTimes := e.lockfile.GetCurrentProject()

	for _, f := range task.Files {
		fo, err := os.Stat(f)
		if err != nil {
			log.Fatalln(err)
		}

		modTimeNow := fo.ModTime().Unix()
		if lockedModTimes[f] < modTimeNow {
			ch <- true
			break
		}
	}

	ch <- false
}

// Dispatches the individual commands of the current task,
// including and events that need to be run.
func (e *Executer) dispatchCommands(task Task, initialRun bool) {
	outputs := make(chan string)

	if initialRun {
		for _, beforeEachCmd := range e.parser.Global.Events.BeforeEachTask {
			e.runSysOrRecurse(beforeEachCmd, &outputs)
		}
	}

	for _, mainCmd := range task.Run {
		if initialRun {
			for _, beforeEachCmd := range e.parser.Global.Events.BeforeEachRun {
				e.runSysOrRecurse(beforeEachCmd, &outputs)
			}
		}

		e.runSysOrRecurse(mainCmd, &outputs)

		if initialRun {
			for _, afterEachCmd := range e.parser.Global.Events.AfterEachRun {
				e.runSysOrRecurse(afterEachCmd, &outputs)
			}
		}
	}

	for _, afterEachCmd := range e.parser.Global.Events.AfterEachTask {
		e.runSysOrRecurse(afterEachCmd, &outputs)
	}
}

// Determine what to execute: system command or another declared task in goke.yml.
func (e *Executer) runSysOrRecurse(cmd string, ch *chan string) {
	if _, ok := e.parser.Tasks[cmd]; ok {
		e.Execute(cmd, false)
	} else {
		go e.runSysCommand(cmd, *ch)
		fmt.Print(<-*ch)
	}
}

// Executes the given string in the underlying OS.
func (e *Executer) runSysCommand(c string, outChan chan string) {
	splitCmd, err := e.parseCommandLine(c)

	if err != nil {
		log.Fatalln(err)
	}

	e.spinner.Message(fmt.Sprintf("Running: %s", c))
	out, err := exec.Command(splitCmd[0], splitCmd[1:]...).Output()

	if err != nil {
		log.Fatalln(err)
	}

	outChan <- "\n" + string(out) + "\n"
}

// Parses the command string into an array of [command, args, args]...
func (e *Executer) parseCommandLine(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true

	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, fmt.Errorf("unclosed quote in command line: %s", command)
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}

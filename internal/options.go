package internal

import (
	"github.com/docopt/docopt-go"
)

const CURRENT_VERSION = "0.2.0"

const usage = `Goke

Usage:
  goke <task> [-w|--watch] [-c|--no-cache] [-f|--force] [-q|--quiet] [-a|--args=<a>...]
  goke -i | --init
  goke -h | --help
  goke -v | --version

Options:
  -h --help      Show this screen
  -v --version   Show version
  -i --init      Creates a goke.yaml file in the current directory
  -w --watch     Run task in watch mode
  -c --no-cache  Clears the program's cache
  -f --force     Runs the task even if files have not been changed
  -a --args=<a>  The arguments and options to pass to the underlying commands
  -q --quiet     Suppresses all output from tasks`

type Options struct {
	TaskName string   `docopt:"<task>"`
	Watch    bool     `docopt:"-w,--watch"`
	NoCache  bool     `docopt:"-c,--no-cache"`
	Force    bool     `docopt:"-f,--force"`
	Quiet    bool     `docopt:"-q,--quiet"`
	Args     []string `docopt:"-a,--args"`
	Init     bool     `docopt:"-i,--init"`
}

func NewCliOptions() Options {
	var opts Options

	parsedDoc, _ := docopt.ParseArgs(usage, nil, CURRENT_VERSION)
	parsedDoc.Bind(&opts)

	return opts
}

func (opts *Options) InitHandler() error {
	if !opts.Init {
		return nil
	}

	err := CreateGokeConfig()
	if err != nil && !opts.Quiet {
		return err
	}

	return nil
}

package internal

import (
	"log"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dugajean/goke/internal/cli"
	"gopkg.in/yaml.v3"
)

type Parseable interface {
	Bootstrap()
	GetGlobal() *Global
	GetTask(string) (Task, bool)
	GetFilePaths() []string
	parseTasks() error
	parseGlobal() error
	expandFilePaths(string) ([]string, error)
	getTempFileName() string
	shouldClearCache(string) bool
}

type Task struct {
	Name  string
	Files []string          `yaml:"files,omitempty"`
	Run   []string          `yaml:"run"`
	Env   map[string]string `yaml:"env,omitempty"`
}

type Global struct {
	Shared struct {
		Env    map[string]string `yaml:"environment,omitempty"`
		Events struct {
			BeforeEachRun  []string `yaml:"before_each_run,omitempty"`
			AfterEachRun   []string `yaml:"after_each_run,omitempty"`
			BeforeEachTask []string `yaml:"before_each_task,omitempty"`
			AfterEachTask  []string `yaml:"after_each_task,omitempty"`
		} `yaml:"events,omitempty"`
	} `yaml:"global,omitempty"`
}

type parser struct {
	Tasks     taskList
	FilePaths []string
	config    string
	options   Options
	fs        FileSystem
	Global
}

type taskList map[string]Task

var osCommandRegexp = regexp.MustCompile(`\$\((.+)\)`)
var parserString string

// NewParser creates a parser instance which can be either a blank one,
// or one provided  from the cache, which gets deserialized.
func NewParser(cfg string, opts *Options, fs FileSystem) Parseable {
	p := parser{}
	p.fs = fs
	p.config = cfg
	p.options = *opts

	tempFile := path.Join(p.fs.TempDir(), p.getTempFileName())

	if p.shouldClearCache(tempFile) {
		_ = p.fs.Remove(tempFile)
	}

	if !p.fs.FileExists(tempFile) {
		return &p
	}

	pBytes, err := p.fs.ReadFile(tempFile)
	if err != nil && !opts.Quiet {
		log.Fatal(err)
	}

	pStr := string(pBytes)
	parserString = pStr

	p = GOBDeserialize(pStr, &p)
	return &p
}

// Bootstrap does the parsing process or skip if cached.
func (p *parser) Bootstrap() {
	// Nothing too bootstrap if cached.
	if parserString != "" {
		return
	}

	err := p.parseGlobal()
	if err != nil && !p.options.Quiet {
		log.Fatal(err)
	}

	err = p.parseTasks()
	if err != nil && !p.options.Quiet {
		log.Fatal(err)
	}

	pStr := GOBSerialize(p)
	err = p.fs.WriteFile(path.Join(p.fs.TempDir(), p.getTempFileName()), []byte(pStr), 0644)

	if err != nil && !p.options.Quiet {
		log.Fatal(err)
	}
}

func (p *parser) GetGlobal() *Global {
	return &p.Global
}

func (p *parser) GetTask(taskName string) (Task, bool) {
	task, ok := p.Tasks[taskName]
	return task, ok
}

func (p *parser) GetFilePaths() []string {
	return p.FilePaths
}

// Parses the individual user defined tasks in the YAML config,
// and processes the dynamic parts of both "run" and "files" sections.
func (p *parser) parseTasks() error {
	var tasks taskList

	if err := yaml.Unmarshal([]byte(p.config), &tasks); err != nil {
		return err
	}

	allFilesPaths := []string{}

	for k, c := range tasks {
		filePaths := []string{}
		for i := range c.Files {
			cli.ReplaceEnvironmentVariables(osCommandRegexp, &tasks[k].Files[i])
			expanded, err := p.expandFilePaths(tasks[k].Files[i])

			if err != nil {
				return err
			}

			filePaths = append(filePaths, expanded...)
			allFilesPaths = append(allFilesPaths, expanded...)
		}

		c.Files = filePaths
		tasks[k] = c

		for i, r := range c.Run {
			tasks[k].Run[i] = strings.Replace(r, "{FILES}", strings.Join(c.Files, " "), -1)
			cli.ReplaceEnvironmentVariables(osCommandRegexp, &tasks[k].Run[i])
		}

		if len(c.Env) != 0 {
			vars, err := cli.SetEnvVariables(c.Env)
			if err != nil {
				return err
			}
			c.Env = vars
		}
		c.Name = k
		tasks[k] = c
	}

	p.FilePaths = allFilesPaths
	p.Tasks = tasks

	return nil
}

// Parses the "global" key in the yaml config and adds it to the parser.
// Also sets all variables under global.environment as OS environment variables.
func (p *parser) parseGlobal() error {
	var g Global

	if err := yaml.Unmarshal([]byte(p.config), &g); err != nil {
		return err
	}

	vars, err := cli.SetEnvVariables(g.Shared.Env)
	if err != nil {
		return nil
	}

	g.Shared.Env = vars
	p.Global = g

	return nil
}

// Expand the path glob and returns all paths in an array
func (p *parser) expandFilePaths(file string) ([]string, error) {
	filePaths := []string{}

	if strings.Contains(file, "*") {
		files, err := p.fs.Glob(file)
		if err != nil {
			return nil, err
		}

		if len(files) > 0 {
			filePaths = append(filePaths, files...)
		}
	} else if p.fs.FileExists(file) {
		filePaths = append(filePaths, file)
	}

	return filePaths, nil
}

// Retrieves the temp file name
func (p *parser) getTempFileName() string {
	cwd, _ := p.fs.Getwd()
	return "goke-" + strings.Replace(cwd, string(filepath.Separator), "-", -1)
}

// Determines whether the parser cache should be cleaned or not
func (p *parser) shouldClearCache(tempFile string) bool {
	tempFileExists := p.fs.FileExists(tempFile)
	mustCleanCache := false

	if !p.options.NoCache && tempFileExists {
		tempStat, _ := p.fs.Stat(tempFile)
		tempModTime := tempStat.ModTime().Unix()

		configStat, _ := p.fs.Stat(CurrentConfigFile())
		configModTime := configStat.ModTime().Unix()

		mustCleanCache = tempModTime < configModTime
	}

	if p.options.NoCache && tempFileExists {
		mustCleanCache = true
	}

	return mustCleanCache
}

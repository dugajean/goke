package internal

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var yamlConfig = ReadYamlConfig()

type Task struct {
	Name  string
	Files []string `yaml:"files,omitempty"`
	Run   []string `yaml:"run"`
}

type Global struct {
	Shared struct {
		Environment map[string]string `yaml:"environment,omitempty"`
		Events      struct {
			BeforeEachRun  []string `yaml:"before_each_run,omitempty"`
			AfterEachRun   []string `yaml:"after_each_run,omitempty"`
			BeforeEachTask []string `yaml:"before_each_task,omitempty"`
			AfterEachTask  []string `yaml:"after_each_task,omitempty"`
		} `yaml:"events,omitempty"`
	} `yaml:"global,omitempty"`
}

type taskList map[string]Task

const osCommandRegexp = `\$\(([\w\d]+)\)`

var parserString string

type Parser struct {
	Tasks     taskList
	FilePaths []string
	Global
}

func NewParser(clearCache bool) Parser {
	tempFile := path.Join(os.TempDir(), getTempFileName())

	if clearCache && FileExists(tempFile) {
		os.Remove(tempFile)
	}

	if !FileExists(tempFile) {
		return Parser{}
	}

	pBytes, err := os.ReadFile(tempFile)
	if err != nil {
		log.Fatal(err)
	}

	pStr := string(pBytes)
	parserString = pStr

	return GOBDeserialize(pStr, &Parser{})
}

func (p *Parser) Bootstrap() {
	if parserString != "" {
		return
	}

	p.parseGlobal()
	p.parseTasks()
	pStr := GOBSerialize(*p)

	err := os.WriteFile(path.Join(os.TempDir(), getTempFileName()), []byte(pStr), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Parses the individual user defined tasks in the YAML config,
// and processes the dynamic parts of both "run" and "files" sections.
func (p *Parser) parseTasks() {
	var tasks taskList

	if err := yaml.Unmarshal([]byte(yamlConfig), &tasks); err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(osCommandRegexp)
	allFilesPaths := []string{}

	for k, c := range tasks {
		filePaths := []string{}
		for i := range c.Files {
			p.replaceEnvironmentVariables(re, &tasks[k].Files[i])
			expanded := p.expandFilePaths(tasks[k].Files[i])
			filePaths = append(filePaths, expanded...)
			allFilesPaths = append(allFilesPaths, expanded...)
		}

		c.Files = filePaths
		tasks[k] = c

		for i, r := range c.Run {
			tasks[k].Run[i] = strings.Replace(r, "$(files)", strings.Join(c.Files, " "), -1)
			p.replaceEnvironmentVariables(re, &tasks[k].Run[i])
		}

		c.Name = k
		tasks[k] = c
	}

	p.FilePaths = allFilesPaths
	p.Tasks = tasks
}

// Parses the "global" key in the yaml config and adds it to the parser.
// Also sets all variables under global.environment as OS environment variables.
func (p *Parser) parseGlobal() {
	var g Global
	if err := yaml.Unmarshal([]byte(yamlConfig), &g); err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(osCommandRegexp)
	for k, v := range g.Shared.Environment {
		_, cmd := p.parseSystemCmd(re, v)

		if cmd == "" {
			continue
		}

		out, err := exec.Command(cmd).Output()

		if err != nil {
			continue
		}

		os.Setenv(k, string(out))
	}

	p.Global = g
}

// Parses the interpolated system commands, ie. "Hello $(echo 'World')" and returns it.
// Returns the command wrapper in $() and without the wrapper.
func (p *Parser) parseSystemCmd(re *regexp.Regexp, str string) (string, string) {
	match := re.FindAllStringSubmatch(str, -1)

	if len(match) > 0 && len(match[0]) > 0 {
		return match[0][0], match[0][1]
	}

	return "", ""
}

// Replace the placeholders with actual environment variable values in string pointer.
// Given that a string pointer must be provided, the replacement happens in place.
func (p *Parser) replaceEnvironmentVariables(re *regexp.Regexp, str *string) {
	resolved := *str
	raw, env := p.parseSystemCmd(re, resolved)

	if raw != "" && env != "" {
		*str = strings.Replace(resolved, raw, os.Getenv(env), -1)
	}
}

func (p *Parser) expandFilePaths(file string) []string {
	filePaths := []string{}

	if strings.Contains(file, "*") {
		files, err := filepath.Glob(file)
		if err != nil {
			log.Fatal(err)
		}

		if len(files) > 0 {
			filePaths = append(filePaths, files...)
		}
	} else if FileExists(file) {
		filePaths = append(filePaths, file)
	}

	return filePaths
}

func getTempFileName() string {
	cwd, _ := os.Getwd()
	return "goke-" + strings.Replace(cwd, string(filepath.Separator), "-", -1)
}

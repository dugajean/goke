package internal

import (
	"testing"

	"github.com/dugajean/goke/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var yamlConfigStub = `
global:
  environment:
    FOO: "foo"
    BAR: "$(echo 'bar')"
    BAZ: "baz"

events:
  before_each_run:
    - "echo 'before each 1'"
    - "echo 'before each 2'"
  after_each_run:
    - "echo 'after each 1'"
    - "greet-lisha"
  before_each_task:
    - "echo 'before task'"
  after_each_task:
    - "echo 'after task'"

greet-lisha:
  run:
    - "echo 'Hello Lisha!'"

greet-loki:
  run:
    - 'echo "Hello Boki"'

greet-cats:
  files: [cmd/cli/*]
  run:
    - 'echo "Hello Frey"'
    - 'echo "Hello Sunny"'
    - "greet-loki"`

var options = Options{
	ClearCache: true,
}

func TestNewParser(t *testing.T) {
	mockStdlib := tests.GetStdlibMock(t).(StdlibWrapper)

	parser := NewParser(yamlConfigStub, &options, mockStdlib)
	assert.NotNil(t, parser)
}

func TestTaskParsing(t *testing.T) {
	mockStdlib := tests.GetStdlibMock(t).(StdlibWrapper)

	parser := NewParser(yamlConfigStub, &options, mockStdlib)
	assert.NotNil(t, parser.Tasks)

	parser.parseTasks()
	assert.NotNil(t, parser.Tasks)

	assert.NotNil(t, parser.Tasks["greet-loki"])
	assert.NotNil(t, parser.Tasks["greet-cats"])
	assert.NotNil(t, parser.Tasks["greet-lisha"])
}

func TestGlobalsParsing(t *testing.T) {
	mockStdlib := tests.NewStdlibWrapper(t)
	mockStdlib.On("TempDir").Return("path/to/temp")
	mockStdlib.On("Getwd").Return("path/to/cwd", nil)
	mockStdlib.On("FileExists", mock.Anything).Return(false)
	mockStdlib.On("Setenv", mock.Anything, mock.Anything).Return(nil)

	parser := NewParser(yamlConfigStub, &options, mockStdlib)

	parser.parseGlobal()
	assert.Equal(t, "foo", parser.Global.Shared.Environment["FOO"])
	assert.Equal(t, "$(echo 'bar')", parser.Global.Shared.Environment["BAR"])
	assert.Equal(t, "baz", parser.Global.Shared.Environment["BAZ"])
	mockStdlib.AssertNumberOfCalls(t, "Setenv", 3)
}

func TestTaskFilesExpansion(t *testing.T) {
	mockStdlib := tests.GetStdlibMock(t).(StdlibWrapper)

	parser := NewParser(yamlConfigStub, &options, mockStdlib)
	parser.parseTasks()

	greetCatsTask := parser.Tasks["greet-cats"]
	assert.NotNil(t, greetCatsTask)
}

func TestParserWithoutCache(t *testing.T) {
	mockStdlib := tests.NewStdlibWrapper(t)
	mockStdlib.On("TempDir").Return("path/to/temp")
	mockStdlib.On("Getwd").Return("path/to/cwd", nil)
	mockStdlib.On("FileExists", mock.Anything).Return(false)

	parser := NewParser(yamlConfigStub, &options, mockStdlib)
	parser.parseTasks()

	greetCatsTask := parser.Tasks["greet-cats"]
	assert.NotNil(t, greetCatsTask)
	mockStdlib.AssertExpectations(t)
	mockStdlib.AssertNumberOfCalls(t, "TempDir", 1)
	mockStdlib.AssertNumberOfCalls(t, "Getwd", 1)
	mockStdlib.AssertNumberOfCalls(t, "FileExists", 2)
}

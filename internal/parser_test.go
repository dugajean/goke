package internal

import (
	"os"
	"strings"
	"testing"

	"github.com/dugajean/goke/internal/tests"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var clearCacheOpts = Options{
	ClearCache: true,
}

var baseOptions = Options{}

func mockCacheDoesNotExist(t *testing.T) *tests.FileSystem {
	fsMock := tests.NewFileSystem(t)
	fsMock.On("TempDir").Return("path/to/temp")
	fsMock.On("Getwd").Return("path/to/cwd", nil)
	fsMock.On("FileExists", mock.Anything).Return(false).Twice()

	return fsMock
}

func mockCacheDoesNotExistOnce(t *testing.T) *tests.FileSystem {
	fsMock := tests.NewFileSystem(t)
	fsMock.On("TempDir").Return("path/to/temp")
	fsMock.On("Getwd").Return("path/to/cwd", nil)
	fsMock.On("FileExists", mock.Anything).Return(false).Once()
	fsMock.On("FileExists", mock.Anything).Return(true).Once()

	return fsMock
}

func mockCacheExists(t *testing.T) *tests.FileSystem {
	fsMock := tests.NewFileSystem(t)
	fsMock.On("TempDir").Return("path/to/temp")
	fsMock.On("Getwd").Return("path/to/cwd", nil)
	fsMock.On("FileExists", mock.Anything).Return(true).Twice()

	return fsMock
}

func TestNewParserWithoutCache(t *testing.T) {
	fsMock := mockCacheDoesNotExist(t)
	parser := NewParser(tests.YamlConfigStub, &clearCacheOpts, fsMock)
	require.NotNil(t, parser)
}

func TestNewParserWithCache(t *testing.T) {
	fsMock := mockCacheDoesNotExistOnce(t)
	fsMock.On("ReadFile", mock.Anything).Return([]byte(tests.ReadFileBase64), nil)

	parser := NewParser(tests.YamlConfigStub, &clearCacheOpts, fsMock)
	require.NotNil(t, parser)
}

func TestNewParserWithCacheAndWithoutClearCacheFlag(t *testing.T) {
	fsMock := mockCacheExists(t)
	fsMock.On("Stat", mock.Anything).Return(tests.MemFileInfo{}, nil).Twice()
	fsMock.On("ReadFile", mock.Anything).Return([]byte(tests.ReadFileBase64), nil).Once()

	parser := NewParser(tests.YamlConfigStub, &baseOptions, fsMock)
	require.NotNil(t, parser)
}

func TestNewParserWithShouldClearCacheTrue(t *testing.T) {
	fsMock := tests.NewFileSystem(t)
	fsMock.On("TempDir").Return("path/to/temp")
	fsMock.On("Getwd").Return("path/to/cwd", nil)
	fsMock.On("FileExists", mock.Anything).Return(true).Once()
	fsMock.On("FileExists", mock.Anything).Return(false).Once()
	fsMock.On("Remove", mock.Anything).Return(nil)

	parser := NewParser(tests.YamlConfigStub, &clearCacheOpts, fsMock)
	require.NotNil(t, parser)
}

func TestTaskParsing(t *testing.T) {
	fsMock := mockCacheDoesNotExist(t)
	fsMock.On("Glob", mock.Anything).Return([]string{"foo", "bar"}, nil).Once()
	parser := NewParser(tests.YamlConfigStub, &clearCacheOpts, fsMock)

	parser.parseTasks()

	greetLoki, _ := parser.GetTask("greet-loki")
	greetCats, _ := parser.GetTask("greet-cats")
	greetLisha, _ := parser.GetTask("greet-lisha")
	require.NotNil(t, greetLoki)
	require.NotNil(t, greetCats)
	require.NotNil(t, greetLisha)
}

func TestGlobalsParsing(t *testing.T) {
	fsMock := mockCacheDoesNotExist(t)
	parser := NewParser(tests.YamlConfigStub, &clearCacheOpts, fsMock)

	parser.parseGlobal()

	require.Equal(t, "foo", os.Getenv("FOO"))
	require.True(t, strings.Contains(os.Getenv("BAR"), "bar"))
	require.Equal(t, "baz", os.Getenv("BAZ"))

	global := parser.GetGlobal()
	require.Equal(t, "foo", global.Shared.Environment["FOO"])
	require.True(t, strings.Contains(global.Shared.Environment["BAR"], "bar"))
	require.Equal(t, "baz", global.Shared.Environment["BAZ"])
}

func TestTaskGlobFilesExpansion(t *testing.T) {
	fsMock := mockCacheDoesNotExist(t)
	fsMock.On("Glob", mock.Anything).Return(tests.ExpectedGlob, nil)
	parser := NewParser(tests.YamlConfigStub, &clearCacheOpts, fsMock)

	parser.parseTasks()
	greetCatsTask, _ := parser.GetTask("greet-cats")

	require.Equal(t, tests.ExpectedGlob, greetCatsTask.Files)
}

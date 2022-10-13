package internal

import (
	"testing"

	"github.com/dugajean/goke/internal/tests"
	"github.com/stretchr/testify/mock"
)

func getDependencies(t *testing.T) (*Parser, *Lockfile, *tests.Process) {
	fsMock := mockCacheDoesNotExist(t)
	fsMock.On("FileExists", mock.Anything).Return(false)
	fsMock.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	fsMock.On("Stat", mock.Anything).Return(tests.MemFileInfo{}, nil)
	fsMock.On("ReadFile", mock.Anything).Return([]byte(dotGokeFile), nil)
	fsMock.On("Glob", mock.Anything).Return(expectedGlob, nil)

	process := tests.NewProcess(t)

	parser := NewParser(yamlConfigStub, &clearCacheOpts, fsMock)
	lockfile := NewLockfile(files, &clearCacheOpts, fsMock)

	parser.Bootstrap()
	lockfile.Bootstrap()

	return &parser, &lockfile, process
}

func TestExecutingCommandOnce(t *testing.T) {
	parser, lockfile, process := getDependencies(t)

	process.On("Execute", mock.Anything, mock.AnythingOfType("string")).Return([]byte("foo"), nil)
	process.On("Fprint", mock.Anything, mock.AnythingOfType("string")).Return(10, nil)

	executor := NewExecutor(parser, lockfile, &clearCacheOpts, process)
	executor.Start("greet-loki")

	process.AssertNumberOfCalls(t, "Execute", 1)
	process.AssertNumberOfCalls(t, "Fprint", 1)

	process.AssertExpectations(t)
}

package internal

import (
	"testing"

	"github.com/dugajean/goke/internal/tests"
	"github.com/stretchr/testify/mock"
)

func getDependencies(t *testing.T) (*Parser, *Lockfile) {
	fsMock := mockCacheDoesNotExist(t)
	fsMock.On("FileExists", mock.Anything).Return(false)
	fsMock.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	fsMock.On("Stat", mock.Anything).Return(tests.MemFileInfo{}, nil)
	fsMock.On("ReadFile", mock.Anything).Return([]byte(dotGokeFile), nil)
	fsMock.On("Glob", mock.Anything).Return(expectedGlob, nil)

	parser := NewParser(yamlConfigStub, &clearCacheOpts, fsMock)
	lockfile := NewLockfile(files, &clearCacheOpts, fsMock)

	parser.Bootstrap()
	lockfile.Bootstrap()

	return &parser, &lockfile
}

func TestExecutingCommandOnce(t *testing.T) {
	parser, lockfile := getDependencies(t)
	executor := NewExecutor(parser, lockfile, &clearCacheOpts)
	executor.Start("greet-loki")
}

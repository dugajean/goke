package internal

import (
	"context"
	"testing"
	"time"

	"github.com/dugajean/goke/internal/tests"
	"github.com/stretchr/testify/mock"
)

func getDependencies(t *testing.T, opts *Options) (*Parseable, *Lockfile, *tests.Process, FileSystem) {
	fsMock := mockCacheDoesNotExist(t)
	fsMock.On("FileExists", mock.Anything).Return(false)
	fsMock.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	fsMock.On("Stat", mock.Anything).Return(tests.MemFileInfo{}, nil)
	fsMock.On("ReadFile", mock.Anything).Return([]byte(dotGokeFile), nil)
	fsMock.On("Glob", mock.Anything).Return(tests.ExpectedGlob, nil)

	process := tests.NewProcess(t)

	parser := NewParser(tests.YamlConfigStub, opts, fsMock)
	lockfile := NewLockfile(files, opts, fsMock)

	parser.Bootstrap()
	lockfile.Bootstrap()

	return &parser, &lockfile, process, fsMock
}

func TestStartNonWatch(t *testing.T) {
	parser, lockfile, process, fsMock := getDependencies(t, &clearCacheOpts)

	process.On("Execute", mock.Anything, mock.AnythingOfType("string")).Return([]byte("foo"), nil)
	process.On("Fprint", mock.Anything, mock.AnythingOfType("string")).Return(10, nil)

	ctx := context.Background()
	executor := NewExecutor(parser, lockfile, &clearCacheOpts, process, fsMock, &ctx)
	executor.Start("greet-loki")

	process.AssertNumberOfCalls(t, "Execute", 1)
	process.AssertNumberOfCalls(t, "Fprint", 1)

	process.AssertExpectations(t)
}

func TestStartWatchWithNoFiles(t *testing.T) {
	watchOpts := Options{
		Watch:   true,
		NoCache: true,
	}

	parser, lockfile, process, fsMock := getDependencies(t, &watchOpts)
	process.On("Exit", mock.Anything).Return()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	executor := NewExecutor(parser, lockfile, &watchOpts, process, fsMock, &ctx)
	executor.Start("greet-loki")
	cancel()

	process.AssertNotCalled(t, "Execute")
	process.AssertNotCalled(t, "Fprint")
	process.AssertNumberOfCalls(t, "Exit", 1)
}

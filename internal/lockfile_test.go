package internal

import (
	"testing"

	"github.com/dugajean/goke/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var files = []string{"foo", "bar", "baz"}

func TestNewLockfile(t *testing.T) {
	stdlibMock := tests.NewStdlibWrapper(t)
	lockfile := NewLockfile(files, stdlibMock)

	assert.NotNil(t, lockfile)
	assert.Equal(t, files, lockfile.files)
}

func TestGenerateLockfileWithTrue(t *testing.T) {
	stdlibMock := tests.NewStdlibWrapper(t)
	stdlibMock.On("Stat", mock.Anything).Return(MemFileInfo{}, nil)
	stdlibMock.On("Getwd").Return("path/to/wd", nil)
	stdlibMock.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	lockfile := NewLockfile(files, stdlibMock)
	err := lockfile.generateLockfile(true)

	assert.Nil(t, err)
	stdlibMock.AssertExpectations(t)
	stdlibMock.AssertNumberOfCalls(t, "Stat", len(files))
	stdlibMock.AssertNumberOfCalls(t, "WriteFile", 1)
	stdlibMock.AssertNumberOfCalls(t, "Getwd", 1)
}

func TestGenerateLockfileWithFalse(t *testing.T) {
	stdlibMock := tests.NewStdlibWrapper(t)
	stdlibMock.On("Stat", mock.Anything).Return(MemFileInfo{}, nil)
	stdlibMock.On("Getwd").Return("path/to/wd", nil)
	stdlibMock.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	lockfile := NewLockfile(files, stdlibMock)
	err := lockfile.generateLockfile(true)

	assert.Nil(t, err)
	stdlibMock.AssertExpectations(t)
}

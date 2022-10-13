package internal

import (
	"fmt"
	"io"
	"os/exec"
)

type Process interface {
	Execute(name string, args ...string) ([]byte, error)
	Fprint(w io.Writer, a ...any) (n int, err error)
}

type ShellProcess struct{}

func (sp *ShellProcess) Execute(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

func (sp *ShellProcess) Fprint(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprint(w, a...)
}

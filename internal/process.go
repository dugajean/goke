package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Process interface {
	Execute(name string, args ...string) ([]byte, error)
	Fprint(w io.Writer, a ...any) (n int, err error)
	Exit(code int)
}

type ShellProcess struct{}

func (sp *ShellProcess) Execute(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

func (sp *ShellProcess) Fprint(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprint(w, a...)
}

func (sp *ShellProcess) Exit(code int) {
	os.Exit(code)
}

package internal

import "os/exec"

type Process interface {
	Execute(name string, args ...[]string) ([]byte, error)
}

type ShellProcess struct{}

func (sp *ShellProcess) Execute(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

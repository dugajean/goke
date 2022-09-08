package internal

import (
	"fmt"
	"os"
)

type Ref[T comparable] struct {
	Value T
	Error error
}

func (e *Ref[T]) Equal(value T) bool {
	return e.Value == value
}

func ReadYamlConfig() string {
	content, err := os.ReadFile("goke.yml")

	if err != nil {
		fmt.Println("no presence of goke sighted")
		os.Exit(1)
	}

	return string(content)
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

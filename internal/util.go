package internal

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
)

func GokeFiles() []string {
	return []string{"goke.yml", "goke.yaml"}
}

func CurrentConfigFile() string {
	for _, f := range GokeFiles() {
		if FileExists(f) {
			return f
		}
	}

	return ""
}

func ReadYamlConfig() (string, error) {
	for _, f := range GokeFiles() {
		content, err := os.ReadFile(f)

		if err == nil && len(content) > 0 {
			return string(content), nil
		}
	}

	return "", errors.New("no presence of goke.yml sighted")
}

func CreateGokeConfig() error {
	const sampleConfig = `global:
environment:
  MY_BINARY: "my_binary"

build: 
  files: [cmd/cli/*.go, internal/*]
  run:
    - "go build -o ./build/${MY_BINARY} ./cmd/cli"
`

	for _, f := range GokeFiles() {
		if FileExists(f) {
			return fmt.Errorf("%s already present in this directory", f)
		}
	}

	return os.WriteFile("goke.yml", []byte(sampleConfig), 0644)
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Serialize a struct
func GOBSerialize[T any](structInstance T) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(structInstance)

	if err != nil {
		log.Fatal("failed gob encode", err)
	}

	return base64.StdEncoding.EncodeToString(b.Bytes())
}

// Deserialize a struct
func GOBDeserialize[T any](structStr string, structShell *T) T {
	by, err := base64.StdEncoding.DecodeString(structStr)

	if err != nil {
		log.Fatal("failed base64 decode", err)
	}

	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(structShell)

	if err != nil {
		log.Fatal("failed gob decode", err)
	}

	return *structShell
}

func PermutateArgs(args []string) int {
	args = args[1:]
	optind := 0

	for i := range args {
		if args[i][0] == '-' {
			tmp := args[i]
			args[i] = args[optind]
			args[optind] = tmp
			optind++
		}
	}

	return optind + 1
}

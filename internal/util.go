package internal

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io/fs"
	"log"
	"os"
	"testing"
	"time"

	"github.com/dugajean/goke/internal/mocks"
	"github.com/stretchr/testify/mock"
)

type Options struct {
	ClearCache bool
	Watch      bool
	Force      bool
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

func GetMock(t *testing.T) StdlibWrapper {
	mockStdlib := mocks.NewStdlibWrapper(t)

	mockStdlib.On("TempDir").Return("path/to/temp")
	mockStdlib.On("Getwd").Return("path/to/cwd", nil)
	mockStdlib.On("FileExists", mock.Anything).Return(true)
	mockStdlib.On("Remove", mock.Anything).Return(nil)
	mockStdlib.On("ReadFile", mock.Anything).Return([]byte("Sf+BAwEBBlBhcnNlcgH/ggABBAEFVGFza3MB/4gAAQlGaWxlUGF0aHMB/4YAAQpZQU1MQ29uZmlnAQwAAQZHbG9iYWwB/4oAAAAZ/4cEAQEIdGFza0xpc3QB/4gAAQwB/4QAACn/gwMBAv+EAAEDAQROYW1lAQwAAQVGaWxlcwH/hgABA1J1bgH/hgAAABb/hQIBAQhbXXN0cmluZwH/hgABDAAAIP+JAwEBBkdsb2JhbAH/igABAQEGU2hhcmVkAf+MAAAA/gGY/4sDAQH+AWtzdHJ1Y3QgeyBFbnZpcm9ubWVudCBtYXBbc3RyaW5nXXN0cmluZyAieWFtbDpcImVudmlyb25tZW50LG9taXRlbXB0eVwiIjsgRXZlbnRzIHN0cnVjdCB7IEJlZm9yZUVhY2hSdW4gW11zdHJpbmcgInlhbWw6XCJiZWZvcmVfZWFjaF9ydW4sb21pdGVtcHR5XCIiOyBBZnRlckVhY2hSdW4gW11zdHJpbmcgInlhbWw6XCJhZnRlcl9lYWNoX3J1bixvbWl0ZW1wdHlcIiI7IEJlZm9yZUVhY2hUYXNrIFtdc3RyaW5nICJ5YW1sOlwiYmVmb3JlX2VhY2hfdGFzayxvbWl0ZW1wdHlcIiI7IEFmdGVyRWFjaFRhc2sgW11zdHJpbmcgInlhbWw6XCJhZnRlcl9lYWNoX3Rhc2ssb21pdGVtcHR5XCIiIH0gInlhbWw6XCJldmVudHMsb21pdGVtcHR5XCIiIH0B/4wAAQIBC0Vudmlyb25tZW50Af+OAAEGRXZlbnRzAf+QAAAAIf+NBAEBEW1hcFtzdHJpbmddc3RyaW5nAf+OAAEMAQwAAP4BWP+PAwEB//1zdHJ1Y3QgeyBCZWZvcmVFYWNoUnVuIFtdc3RyaW5nICJ5YW1sOlwiYmVmb3JlX2VhY2hfcnVuLG9taXRlbXB0eVwiIjsgQWZ0ZXJFYWNoUnVuIFtdc3RyaW5nICJ5YW1sOlwiYWZ0ZXJfZWFjaF9ydW4sb21pdGVtcHR5XCIiOyBCZWZvcmVFYWNoVGFzayBbXXN0cmluZyAieWFtbDpcImJlZm9yZV9lYWNoX3Rhc2ssb21pdGVtcHR5XCIiOyBBZnRlckVhY2hUYXNrIFtdc3RyaW5nICJ5YW1sOlwiYWZ0ZXJfZWFjaF90YXNrLG9taXRlbXB0eVwiIiB9Af+QAAEEAQ1CZWZvcmVFYWNoUnVuAf+GAAEMQWZ0ZXJFYWNoUnVuAf+GAAEOQmVmb3JlRWFjaFRhc2sB/4YAAQ1BZnRlckVhY2hUYXNrAf+GAAAA/gMv/4IBBQZnbG9iYWwBBmdsb2JhbAAGZXZlbnRzAQZldmVudHMAC2dyZWV0LWxpc2hhAQtncmVldC1saXNoYQIBE2VjaG8gJ0hlbGxvIExpc2hhIScACmdyZWV0LWxva2kBCmdyZWV0LWxva2kCARFlY2hvICJIZWxsbyBCb2tpIgAKZ3JlZXQtY2F0cwEKZ3JlZXQtY2F0cwEBD2NtZC9jbGkvbWFpbi5nbwEDEWVjaG8gIkhlbGxvIEZyZXkiEmVjaG8gIkhlbGxvIFN1bm55IgpncmVldC1sb2tpAAEBD2NtZC9jbGkvbWFpbi5nbwH+AhIKZ2xvYmFsOgogIGVudmlyb25tZW50OgogICAgRk9POiAiZm9vIgogICAgQkFSOiAiJChlY2hvICdiYXInKSIKICAgIEJBWjogImJheiIKCmV2ZW50czoKICBiZWZvcmVfZWFjaF9ydW46CiAgICAtICJlY2hvICdiZWZvcmUgZWFjaCAxJyIKICAgIC0gImVjaG8gJ2JlZm9yZSBlYWNoIDInIgogIGFmdGVyX2VhY2hfcnVuOgogICAgLSAiZWNobyAnYWZ0ZXIgZWFjaCAxJyIKICAgIC0gImdyZWV0LWxpc2hhIgogIGJlZm9yZV9lYWNoX3Rhc2s6CiAgICAtICJlY2hvICdiZWZvcmUgdGFzayciCiAgYWZ0ZXJfZWFjaF90YXNrOgogICAgLSAiZWNobyAnYWZ0ZXIgdGFzayciCgpncmVldC1saXNoYToKICBydW46CiAgICAtICJlY2hvICdIZWxsbyBMaXNoYSEnIgoKZ3JlZXQtbG9raToKICBydW46CiAgICAtICdlY2hvICJIZWxsbyBCb2tpIicKCmdyZWV0LWNhdHM6CiAgZmlsZXM6IFtjbWQvY2xpLypdCiAgcnVuOgogICAgLSAnZWNobyAiSGVsbG8gRnJleSInCiAgICAtICdlY2hvICJIZWxsbyBTdW5ueSInCiAgICAtICJncmVldC1sb2tpIgEBAQMDQkFaA2JhegNGT08DZm9vA0JBUg0kKGVjaG8gJ2JhcicpAQAAAAA="), nil)

	return mockStdlib
}

// Parses the command string into an array of [command, args, args]...
func ParseCommandLine(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true

	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, fmt.Errorf("unclosed quote in command: %s", command)
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}

type MemFileInfo struct{}

func (fi MemFileInfo) Name() string {
	return "foo"
}

func (fi MemFileInfo) Size() int64 {
	return 10000
}

func (fi MemFileInfo) Mode() fs.FileMode {
	return 0644
}

func (fi MemFileInfo) ModTime() time.Time {
	return time.Date(2022, time.December, 24, 1, 1, 1, 1, time.UTC)
}

func (fi MemFileInfo) IsDir() bool {
	return false
}

func (fi MemFileInfo) Sys() any {
	return nil
}

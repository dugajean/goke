package internal

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

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

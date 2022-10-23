package cli

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSystemCmd(t *testing.T) {

	cmds := []string{
		"$(echo 'Hello Thor')",
		"hello world",
	}

	want := [][]string{
		{"$(echo 'Hello Thor')", "echo 'Hello Thor'"},
		{"", ""},
	}

	for i, cmd := range cmds {
		got0, got1 := parseSystemCmd(osCommandRegexp, cmd)
		assert.Equal(t, want[i][0], got0, "expected "+want[i][0]+", got ", got0)
		assert.Equal(t, want[i][1], got1, "expected "+want[i][1]+", got ", got1)
	}

}

func TestSetEnvVariables(t *testing.T) {

	values := map[string]string{
		"THOR":     "Lord of thunder",
		"THOR_CMD": "$(echo 'Hello Thor')",
	}

	want := map[string]string{
		"THOR":     "Lord of thunder",
		"THOR_CMD": "Hello Thor",
	}

	got, _ := SetEnvVariables(values)
	assert.Equal(t, want["THOR"], os.Getenv("THOR"))
	assert.Equal(t, want["THOR_CMD"], os.Getenv("THOR_CMD"))

	for k := range got {
		assert.Equal(t, want[k], got[k])
	}
}

func TestParseCommandLine(t *testing.T) {
	t.Skip()
}

func TestReplaceEnvironmentVariables(t *testing.T) {
	values := map[string]string{
		"THOR": "Lord of thunder",
		"LOKI": "Lord of deception",
	}

	for k, v := range values {
		t.Setenv(k, v)
	}

	str := "I am ${THOR}"
	want := "I am Lord of thunder"

	ReplaceEnvironmentVariables(osEnvRegexp, &str)

	assert.Equal(t, want, str, "wrong env value is injected")
}

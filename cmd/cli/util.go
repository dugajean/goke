package main

import (
	"fmt"
	"os"

	app "github.com/dugajean/goke/internal"
)

// handleGlobalOptions executes the global options logic.
func handleGlobalOptions(opts *app.Options, p *app.Parseable) {
	if code, err := opts.InitHandler(); code != -1 {
		if err != nil {
			fmt.Println(err)
		}

		os.Exit(code)
	}

	if code := opts.TasksHandler(p); code != -1 {
		os.Exit(code)
	}
}

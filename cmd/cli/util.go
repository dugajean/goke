package main

import (
	"fmt"
	"os"

	app "github.com/dugajean/goke/internal"
)

// handleGlobalOptions executes the global options logic.
func handleGlobalOptions(opts *app.Options, p *app.Parseable) {
	for _, fn := range opts.Handlers(p) {
		if code, err := fn(p); code != -1 {
			if err != nil {
				fmt.Println(err)
			}

			os.Exit(code)
		}
	}
}

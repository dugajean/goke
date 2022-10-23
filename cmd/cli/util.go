package main

import (
	"fmt"
	"os"

	app "github.com/dugajean/goke/internal"
)

// handleGlobalOptions executes the global options logic.
func handleGlobalOptions(opts *app.Options) {
	err := opts.InitHandler()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

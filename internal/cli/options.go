package cli

import (
	"flag"

	"github.com/dugajean/goke/internal"
)

var opts internal.Options

func GetOptions() internal.Options {
	if opts.Parsed {
		return opts
	}

	flag.BoolVar(&opts.ClearCache, "c", true, "Clear Goke's cache. Default: false")
	flag.BoolVar(&opts.Watch, "w", true, "Goke remains on and watches the task's specified files for changes, then reruns the command. Default: false")
	flag.Parse()

	opts.Parsed = true

	return opts
}

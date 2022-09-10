package cli

import (
	"flag"

	"github.com/dugajean/goke/internal"
)

func GetOptions() internal.Options {
	var opts internal.Options

	flag.BoolVar(&opts.ClearCache, "no-cache", false, "Clear Goke's cache. Default: false")
	flag.BoolVar(&opts.Watch, "watch", false, "Goke remains on and watches the task's specified files for changes, then reruns the command. Default: false")
	flag.Parse()

	return opts
}

package internal

import (
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"strings"
)

const GITHUB_TAGS_ENDPOINT = "https://api.github.com/repos/dugajean/goke/git/refs/tags"

type Options struct {
	ClearCache bool
	Watch      bool
	Force      bool
	Init       bool
	Quiet      bool
	Version    bool
}

func GetOptions() Options {
	var opts Options

	flag.BoolVar(&opts.ClearCache, "no-cache", false, "Clear Goke's cache. Default: false")
	flag.BoolVar(&opts.Watch, "watch", false, "Goke remains on and watches the task's specified files for changes, then reruns the command. Default: false")
	flag.BoolVar(&opts.Force, "force", false, "Executes the task regardless whether the files have changed or not. Default: false")
	flag.BoolVar(&opts.Init, "init", false, "Initializes a goke.yml file in the current directory")
	flag.BoolVar(&opts.Quiet, "quiet", false, "Disables all output to the console. Default: false")
	flag.BoolVar(&opts.Version, "version", false, "Prints the current Goke version")
	flag.Parse()

	return opts
}

func (opts *Options) InitHandler() error {
	if !opts.Init {
		return nil
	}

	err := CreateGokeConfig()
	if err != nil && !opts.Quiet {
		return err
	}

	return nil
}

func (opts *Options) VersionHandler() (string, error) {
	if !opts.Version {
		return "", nil
	}

	res, err := http.Get(GITHUB_TAGS_ENDPOINT)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var refs []struct {
		Ref string `json:"ref,omitempty"`
	}

	err = json.NewDecoder(res.Body).Decode(&refs)
	if err != nil {
		return "", err
	}

	if len(refs) == 0 {
		return "", errors.New("could not retrieve version")
	}

	lastRef := refs[len(refs)-1].Ref
	lastRefSplit := strings.Split(lastRef, "/")
	lastRef = lastRefSplit[len(lastRefSplit)-1]

	return lastRef, nil
}

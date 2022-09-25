package internal

import (
	"encoding/json"
	"errors"
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

// Package github is a wrapper for some of the GitHub API.
// See https://docs.github.com/en/rest?apiVersion=2022-11-28
package github

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/version-release/avr/pkg/lib"
	"io"
	"net/http"
	"os"
)

const (
	EnvActions         = "GITHUB_ACTIONS"
	EnvActor           = "GITHUB_ACTOR"
	EnvApiUrl          = "GITHUB_API_URL"
	EnvRepository      = "GITHUB_REPOSITORY"
	EnvRepositoryOwner = "GITHUB_REPOSITORY_OWNER"
	EnvToken           = "GH_TOKEN"
	Name               = "github"
)

var (
	envs         map[string]string
	PublicApiUrl = "https://api.github.com"
	PublicServer = "github.com"
)

func init() {
	processEnvironment()
}

// processEnvironment Load the environment variables and panic when any are not set.
func processEnvironment() {
	var e error
	vars := []string{}

	if os.Getenv(EnvActions) == "true" {
		// were not in the CircleCI environment, so exit.
		vars = append(vars,
			EnvActor,
			EnvRepository,
			EnvRepositoryOwner,
		)
	}

	envs, e = lib.GetRequiredEnvVars(append(vars, EnvApiUrl, EnvToken))

	if e != nil {
		panic(e.Error())
	}
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func KnownSshKeys(apiUrl, host string, client *http.Client) (string, error) {
	res, e1 := client.Get(apiUrl)
	if e1 != nil {
		return "", fmt.Errorf(stderr.GetMeta, e1.Error())
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf(stderr.GetMeta)
	}

	content, e2 := io.ReadAll(res.Body)
	if e2 != nil {
		return "", fmt.Errorf(stderr.GetMetaBody, e2.Error())
	}

	ghMeta := &Meta{}
	if e := json.Unmarshal(content, ghMeta); e != nil {
		return "", fmt.Errorf(stderr.MergeWaitTimeout, e.Error())
	}

	known := ""
	for _, key := range ghMeta.SshKeys {
		known += fmt.Sprintf("%v %v\n", host, key)
	}

	return known, nil
}

type Meta struct {
	SshKeys []string `json:"ssh_keys"`
}

// Package github is a wrapper for some of the GitHub API.
// See https://docs.github.com/en/rest?apiVersion=2022-11-28
package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	EnvMergeType = "PARAM_MERGE_TYPE"
	EnvToken     = "GITHUB_TOKEN"
	EnvServer    = "PARAM_GH_SERVER"
	Name         = "github"
	Server       = "github.com"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	envToken = ""
)

func init() {
	// set the token from the environment.
	val, k := os.LookupEnv(EnvToken)
	if k {
		envToken = val
	}
}

func NewRequest(uri, method, token string) (*http.Request, error) {
	req, err1 := http.NewRequest(method, uri, nil)
	if err1 != nil {
		return nil, fmt.Errorf(stderr.CouldNotPrepareRequest, err1.Error())
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	return req, nil
}

func parseRepositoryUri(uri string) (string, string, string) {
	//https://github.com/kohirens/version-release
	//git@github.com:kohirens/version-release.git
	re := regexp.MustCompile(`^(https://|git@)([^/:]+)[/:]([^/]+)/(.+)`)
	m := re.FindAllStringSubmatch(uri, -1)

	if m != nil {
		return m[0][2], m[0][3], strings.Replace(m[0][4], ".git", "", 1)
	}

	return "", "", ""
}

func KnownSshKeys(client *http.Client) (string, error) {
	res, e1 := client.Get("https://api.github.com/meta")
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
		known += fmt.Sprintf("github.com %v\n", key)
	}

	return known, nil
}

type Meta struct {
	SshKeys []string `json:"ssh_keys"`
}

// Package github is a small subset wrapper for the GitHub API.
// See https://docs.github.com/en/rest?apiVersion=2022-11-28
package github

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
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
	//https://github.com/kohirens/version-release-orb
	//git@github.com:kohirens/version-release-orb.git
	re := regexp.MustCompile(`^(https://|git@)([^/:]+)[/:]([^/]+)/(.+)`)
	m := re.FindAllStringSubmatch(uri, -1)

	if m != nil {
		return m[0][2], m[0][3], strings.Replace(m[0][4], ".git", "", 1)
	}

	return "", "", ""
}

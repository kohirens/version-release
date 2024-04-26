package gittoolbelt

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/log"
	"strings"
)

const (
	cmdGitToolBelt = "git-tool-belt"
)

type SemverInfo struct {
	CommitHash        string `json:"CommitHash"`
	CurrentVersion    string `json:"currentVersion"`
	NextVersion       string `json:"nextVersion"`
	NextVersionReason string `json:"nextVersionReason"`
}

// IsTaggable Looks at the current directory as a repo and branch, then checks
// return a boolean to indicate if any of the commits since the last release
// can be tagged (based on the conventional commit message).
func IsTaggable(wd string) bool {
	so, se, _, co := cli.RunCommand(
		wd,
		cmdGitToolBelt,
		[]string{"taggable"},
	)

	log.Infof(stdout.Co, co)

	if se != nil {
		log.Logf(stdout.CommitNotTagged, so, se.Error())
		return false
	}

	return strings.Trim(string(so), "\n") == "true"
}

// Semver Determine semantic (current and possibly the next) version info using
// git-tool-belt.
func Semver(wd string) (*SemverInfo, error) {
	so, se, _, _ := cli.RunCommand(
		wd,
		cmdGitToolBelt,
		[]string{"semver"},
	)
	if se != nil {
		return nil, fmt.Errorf(stderr.CouldNotDetermineSemver, wd, se.Error())
	}

	si := &SemverInfo{}
	if e := json.Unmarshal(so, si); e != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonDecode, e.Error())
	}

	return si, nil
}

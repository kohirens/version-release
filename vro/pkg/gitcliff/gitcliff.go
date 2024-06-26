package gitcliff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release/vro/internal/util"
	"os"
	"regexp"
	"strings"
)

const (
	CliffConfigName = "cliff.toml"
	Cmd             = "git-cliff"
)

// BuildChangelog Update the change log file with any unreleased changes.
//
//	semVer allows you to set the tag manually
func BuildChangelog(wd, chgLogFile, semVer string) ([]string, error) {
	configFile := wd + "/" + CliffConfigName
	files := []string{chgLogFile}

	if !fsio.Exist(configFile) { // make a config when none present.
		if e := os.WriteFile(configFile, []byte(cliffConfig), 0776); e != nil {
			return nil, e
		}
		files = append(files, CliffConfigName)
	}

	// build new
	args := []string{"--bump", "--unreleased", "--output", chgLogFile}

	// prepend changes.
	if fsio.Exist(wd + stdlib.PS + chgLogFile) {
		args = []string{"--bump", "--unreleased", "--prepend", chgLogFile}
	}
	if semVer != "" {
		args = append(args, "--tag", semVer)
	}

	_, se, _, cs := cli.RunCommand(
		wd,
		Cmd,
		args,
	)

	log.Infof(stdout.Cs, cs)

	if se != nil && strings.Contains(se.Error(), "WARN") {
		return nil, fmt.Errorf(stderr.UpdateChgLog, se.Error())
	}

	return files, nil
}

// UnreleasedMessage Get unreleased commits changes without header and footer.
func UnreleasedMessage(wd string) ([]byte, error) {
	args := []string{"--bump", "--strip", "all", "--unreleased"}
	so, se, _, cs := cli.RunCommand(
		wd,
		Cmd,
		args,
	)

	log.Infof(stdout.Cs, cs)

	if se != nil && strings.Contains(se.Error(), "WARN") {
		return so, fmt.Errorf(stderr.UnreleasedMsg, se.Error())
	}

	return so, nil
}

// Bump Returns the next semantic version if there are unreleased changes.
//
//	This will return an empty string if there are no released changes, however
//	that does NOT mean the changelog is up-to-date.
func Bump(wd string) string {
	so, se, _, cs := cli.RunCommand(
		wd,
		Cmd,
		[]string{"--bumped-version"},
	)

	log.Infof(stdout.Wd, wd)
	log.Infof(stdout.Cs, cs)

	if se != nil {
		log.Errf(stderr.BumpedVersion, se.Error())
		return ""
	}

	if bytes.Contains(so, []byte("There is nothing to bump")) { // every thing is up-to-date
		return ""
	}

	re := regexp.MustCompile(util.CheckSemVer)

	var found []byte
	out := bytes.Split(bytes.Trim(so, "\n"), []byte("\n"))

	if len(out) > 1 { // this means no previous releases found and there are unreleased changes to bump
		found = re.Find(out[1])
	} else { // this means a previous release exist and there are unreleased changes to bump
		found = re.Find(out[0])
	}

	return string(found)
}

// NextVersion Calculated semantic version.
func NextVersion(wd string) (string, error) {
	var version string

	u, e1 := UnreleasedChanges(wd)
	if e1 != nil {
		return "", e1
	}

	if len(u) > 0 {
		version = u[0].Version
	}

	if version == "" {
		return "", fmt.Errorf(stderr.NoVersionTag)
	}

	return version, nil
}

// HasUnreleasedChanges Indicate there are changes in the current branch that
// needed to be added to the changelog and tagged.
//
//	This makes use of the --context flag to return any unreleased commit.
func HasUnreleasedChanges(wd string) (bool, error) {
	so, se, _, cs := cli.RunCommand(
		wd,
		Cmd,
		[]string{"--unreleased", "--context"},
	)

	log.Infof(stdout.Cs, cs)

	if se != nil {
		return false, fmt.Errorf(stderr.UpdateChgLog, se.Error())
	}

	if len(so) < 1 || bytes.Equal(bytes.Trim(so, "\n"), []byte("[]")) {
		return false, nil
	}

	var u []Unreleased

	if e := json.Unmarshal(so, &u); e != nil {
		return false, fmt.Errorf(stderr.CannotDecodeJson, e.Error())
	}

	return len(u[0].Commits) > 0, nil
}

// UnreleasedChanges Indicate there are changes in the current branch that
// needed to be added to the changelog and tagged.
//
//	This makes use of the --context flag to return any unreleased commit.
func UnreleasedChanges(wd string) ([]Unreleased, error) {
	so, se, _, cs := cli.RunCommand(
		wd,
		Cmd,
		[]string{"--unreleased", "--context", "--bump"},
	)

	log.Infof(stdout.Cs, cs)

	if se != nil {
		return nil, fmt.Errorf(stderr.UpdateChgLog, se.Error())
	}

	if len(so) < 1 {
		return nil, nil
	}

	// remove extra INFO|WARN git-cliff output
	soClean := so
	idx := bytes.IndexByte(so, '[')
	if idx == -1 {
		return nil, nil
	}

	soClean = so[idx:]
	var u []Unreleased

	if e := json.Unmarshal(soClean, &u); e != nil {
		return nil, fmt.Errorf(stderr.CannotDecodeJson, e.Error())
	}

	return u, nil
}

type Unreleased struct {
	Version string `json:"version"`
	Commits []struct {
		Id                  string        `json:"id"`
		Message             string        `json:"message"`
		Body                interface{}   `json:"body"`
		Footers             []interface{} `json:"footers"`
		Group               string        `json:"group"`
		BreakingDescription interface{}   `json:"breaking_description"`
		Breaking            bool          `json:"breaking"`
		Scope               interface{}   `json:"scope"`
		Links               []interface{} `json:"links"`
		Author              struct {
			Name      string `json:"name"`
			Email     string `json:"email"`
			Timestamp int    `json:"timestamp"`
		} `json:"author"`
		Committer struct {
			Name      string `json:"name"`
			Email     string `json:"email"`
			Timestamp int    `json:"timestamp"`
		} `json:"committer"`
		Conventional bool `json:"conventional"`
		MergeCommit  bool `json:"merge_commit"`
		Github       struct {
			Username    interface{}   `json:"username"`
			PrTitle     interface{}   `json:"pr_title"`
			PrNumber    interface{}   `json:"pr_number"`
			PrLabels    []interface{} `json:"pr_labels"`
			IsFirstTime bool          `json:"is_first_time"`
		} `json:"github"`
	} `json:"commits"`
	CommitId  interface{} `json:"commit_id"`
	Timestamp int         `json:"timestamp"`
	Previous  interface{} `json:"previous"`
	Github    struct {
		Contributors []interface{} `json:"contributors"`
	} `json:"github"`
}

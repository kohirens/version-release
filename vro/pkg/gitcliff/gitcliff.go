package gitcliff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/log"
	"regexp"
	"strings"
)

const (
	Cmd = "git-cliff"
)

// BuildChangelog Runs git-cliff to update the change log file.
func BuildChangelog(wd, chgLogFile string) error {
	// Note footer links may not be generated with these methods.
	// build new: git-cliff --output CHANGELOG.md
	args := []string{"--bump", "--output", chgLogFile}
	if stdlib.PathExist(wd + stdlib.PS + chgLogFile) {
		// update existing: git-cliff --unreleased --bump --prepend CHANGELOG.md
		args = []string{"--bump", "--unreleased", "--prepend", chgLogFile}
	}
	_, se, _, cs := cli.RunCommand(
		wd,
		Cmd,
		args,
	)

	log.Infof(stdout.Cs, cs)

	if se != nil && strings.Contains(se.Error(), "WARN") {
		return fmt.Errorf(stderr.UpdateChgLog, se.Error())
	}

	return nil
}

// Bump Returns the next semantic version
func Bump(wd string) string {
	so, se, _, cs := cli.RunCommand(
		wd,
		Cmd,
		[]string{"--bumped-version"},
	)

	log.Infof(stdout.Cs, cs)

	if se != nil {
		log.Errf(stderr.UpdateChgLog, se.Error())
		return ""
	}

	re := regexp.MustCompile(`^v?\d\.\d\.\d`)

	var found []byte
	out := bytes.Split(bytes.Trim(so, "\n"), []byte("\n"))

	if len(out) > 1 {
		found = re.Find(out[1])
	} else {
		found = re.Find(out[0])
	}

	return string(found)
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

	var u []unreleased

	if e := json.Unmarshal(so, &u); e != nil {
		return false, fmt.Errorf(stderr.CannotDecodeJson, e.Error())
	}

	return len(u) > 0, nil
}

type unreleased struct {
	Version interface{} `json:"version"`
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

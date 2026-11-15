package gitcliff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kohirens/stdlib"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/stdlib/str"
	"github.com/kohirens/version-release/avr/pkg/lib"
)

const (
	CliffConfigName      = "cliff.toml"
	unreleasedConfigName = "unreleased-cliff.toml"
	Cmd                  = "git-cliff"
)

var log = logger.Standard{}

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

// BuildChangelog Update the change log file with any unreleased changes.
// This writes the changelog file, overwriting/updating an existing file, or making a new one.
//
//	semVer while optional will allow you to set the tag manually
func BuildChangelog(wd, chgLogFile, semVer string) error {
	// build new
	args := []string{"--bump", "--unreleased", "--output", chgLogFile}

	// prepend changes when a changelog already exist.
	if fsio.Exist(wd + stdlib.PS + chgLogFile) {
		log.Dbugf(stdout.PrependToChangelog)
		args = []string{"--bump", "--unreleased", "--prepend", chgLogFile}
	}
	if semVer != "" {
		log.Dbugf(stdout.AddTagManually, semVer)
		args = append(args, "--tag", semVer)
	}

	_, se, _, cs := cli.RunCommand(
		wd,
		Cmd,
		args,
	)

	log.Infof(stdout.Cs, cs)

	if se != nil {
		return fmt.Errorf(stderr.UpdateChgLog, se.Error())
	}

	return nil
}

// BuildConfig Build a default config file when none exist.
func BuildConfig(configFile string) (bool, error) {
	if fsio.Exist(configFile) { // a config already exist.
		return false, nil
	}

	if e := os.WriteFile(configFile, []byte(cliffConfig), 0664); e != nil {
		// Attempted but failed to write the config file.
		return false, e
	}

	// A config file was written successfully.
	return true, nil
}

// Bump Returns the next semantic version if there are unreleased changes.
//
//	This will return an empty string if there are no released changes, however
//	that does NOT mean the changelog is up-to-date.
func Bump(wd string, enableTagVPrefix bool) string {
	so, se, _, cs := cli.RunCommand(
		wd,
		Cmd,
		[]string{"--bumped-version"},
	)

	log.Infof(stdout.Wd, Cmd, wd)
	log.Infof(stdout.Cs, cs)

	if se != nil {
		log.Errf(stderr.BumpedVersion, se.Error())
		return ""
	}

	if bytes.Contains(so, []byte(stdout.NoCommitsToBump)) { // every thing is up-to-date
		return ""
	}

	log.Dbugf("so: %v", so)

	re := regexp.MustCompile(lib.CheckSemVer)

	log.Dbugf(stdout.Fmt, cs, so)

	soClean := bytes.Trim(scrubOutput(so), "\n")
	var found []byte
	out := bytes.Split(soClean, []byte("\n"))

	if len(out) > 1 { // this means no previous releases found and there are unreleased changes to bump
		found = re.Find(out[1])
	} else { // this means a previous release exist and there are unreleased changes to bump
		found = re.Find(out[0])
	}

	sv := string(found)

	if enableTagVPrefix && !strings.HasPrefix(sv, "v") {
		return "v" + string(found)
	}

	return sv
}

func NextVersion(wd, nextVer string, enableTagVPrefix bool) string {
	// check if a version has been provided as input.
	if nextVer == "" {
		nextVer = Bump(wd, enableTagVPrefix)
	}

	log.Infof(stdout.NextSemVer, nextVer)

	return nextVer
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

	log.Dbugf(stdout.Fmt, cs, so)

	// remove extra INFO|WARN git-cliff output
	soClean := scrubOutput(so)
	var u []Unreleased

	if e := json.Unmarshal(soClean, &u); e != nil {
		return nil, fmt.Errorf(stderr.CannotDecodeJson, e.Error())
	}

	return u, nil
}

// UnreleasedChangesHash Calculate a SHA256 hash for the changes to be added to the changelog.
func UnreleasedChangesHash(wd string) ([]byte, []byte, error) {
	changes, e2 := unreleasedChangesCommits(wd)
	if e2 != nil {
		return nil, nil, e2
	}

	// use a special git-cliff config that does not contain dates. This should allow for a hash that only changes when
	// a user makes changes to it. Not just because it picked the current date.
	hash, e3 := str.Sha256(string(changes))
	if e3 != nil {
		return nil, nil, fmt.Errorf(stderr.CannotCalcHash, e3.Error())
	}

	return changes, hash, nil
}

// unreleasedChangesCommits Get unreleased commits changes without header and footer.
func unreleasedChangesCommits(wd string) ([]byte, error) {
	_, e1 := BuildConfig(wd + "/" + unreleasedConfigName)
	if e1 != nil {
		return nil, e1
	}

	// use a config generated from the built-in and remove the header and footer
	args := []string{"--bump", "--strip", "all", "--unreleased", "--config", unreleasedConfigName}
	so, se, _, cs := cli.RunCommand(
		wd,
		Cmd,
		args,
	)

	log.Infof(stdout.Cs, cs)

	soClean := scrubOutput(so)

	if se != nil {
		return soClean, fmt.Errorf(stderr.UnreleasedMsg, se.Error())
	}

	return soClean, nil
}

// scrubOutput Git-Cliff unfortunately adds diagnostic messages like warnings to
// the standard output along with program data, making it hard to parse. This
// trims that output.
func scrubOutput(so []byte) []byte {
	if len(so) < 1 {
		return []byte{}
	}

	lines := bytes.Split(so, []byte("\n"))
	re := regexp.MustCompile(`(INFO|ERROR|WARN)\s*.*$`)
	start := 0

	for i, line := range lines {
		log.Dbugf(stdout.Scrubbing, fmt.Sprintf("line %v", i), line)
		if !re.Match(line) {
			break
		}

		// Skip over GitCliff log output.
		start = i + 1
	}

	if start > 0 {
		trimmed := bytes.Join(lines[start:], []byte("\n"))
		log.Dbugf(stdout.Trimmed, trimmed)
		return trimmed
	}

	return so
}

package gitcliff

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/log"
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

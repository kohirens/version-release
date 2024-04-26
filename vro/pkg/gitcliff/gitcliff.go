package gitcliff

import (
	"fmt"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/log"
)

const (
	Cmd = "git-cliff"
)

// RebuildChangelog Runs git-cliff to update the change log file.
func RebuildChangelog(wd, chgLogFile string) error {
	_, se, _, cs := cli.RunCommand(
		wd,
		Cmd,
		[]string{"--output", chgLogFile, "--bump"},
	)

	log.Infof(stdout.Cs, cs)

	if se != nil {
		return fmt.Errorf(stderr.UpdateChgLog, se.Error())
	}

	return nil
}

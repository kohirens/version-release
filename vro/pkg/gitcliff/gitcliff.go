package gitcliff

import (
	"fmt"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release-orb/vro/pkg/gittoolbelt"
)

const (
	Cmd = "git-cliff"
)

// RebuildChangelog Runs git-cliff to update the change log file.
func RebuildChangelog(wd, chgLogFile string, si *gittoolbelt.SemverInfo) error {
	_, se, _, cs := cli.RunCommand(
		wd,
		Cmd,
		[]string{"--output", chgLogFile, "--tag", si.NextVersion},
	)

	log.Infof(stdout.Cs, cs)

	if se != nil {
		return fmt.Errorf(stderr.CouldNotUpdateChgLog, se.Error())
	}

	return nil
}

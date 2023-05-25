package gitchglog

import (
	"fmt"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release-orb/vro/pkg/gittoolbelt"
)

// RebuildChangelog Runs git-chglog to update the change log file.
func RebuildChangelog(wd, chgLogFile string, si *gittoolbelt.SemverInfo) error {
	_, se, _, cs := cli.RunCommand(
		wd,
		"git-chglog",
		[]string{"--output", chgLogFile, "--next-tag", si.NextVersion},
	)

	log.Infof(stdout.Cs, cs)

	if se != nil {
		return fmt.Errorf(stderr.CouldNotUpdateChgLog, se.Error())
	}

	return nil
}

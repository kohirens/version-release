package github

import (
	"fmt"
	"github.com/kohirens/stdlib/log"
	"os"
)

const (
	gaOutputFile = "GITHUB_OUTPUT"
)

// AddOutputVar Add a variable to the GitHUB Actions output file.
//
//	We set an output variable that will cause the next workflow to execute if its condition is met.
//	The workflow has three possible values: none|publish-changelog|publish-release-tag.
//	These are the same values used to trigger the equivalent CircleCI workflows.
//	For GitHub Actions we merely need to set an output variable to continue onto the next workflow.
func AddOutputVar(name, nameWorkflow string) error {
	outFile, ok := os.LookupEnv(gaOutputFile)
	if !ok {
		return fmt.Errorf("env var %s not set", gaOutputFile)
	}

	fileHandle, e2 := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if e2 != nil {
		return fmt.Errorf(stderr.OpenFile, outFile, e2)
	}

	_, e3 := fileHandle.Write([]byte(name + "=" + nameWorkflow + "\n"))

	return e3
}

func DumpOutput() error {
	outFile, ok := os.LookupEnv(gaOutputFile)
	if !ok {
		return fmt.Errorf("env var %s not set", gaOutputFile)
	}

	contents, e1 := os.ReadFile(outFile)
	if e1 != nil {
		return fmt.Errorf(stderr.ReadFile, outFile, e1)
	}

	log.Dbugf(stdout.GitHubOutputsFile, outFile, string(contents))

	return nil
}

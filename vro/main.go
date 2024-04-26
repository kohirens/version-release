//go:generate git-tool-belt semver -save info.go -format go -packageName main -varName af

package main

import (
	"flag"
	"fmt"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release-orb/vro/pkg/circleci"
	"github.com/kohirens/version-release-orb/vro/pkg/github"
	"github.com/kohirens/version-release-orb/vro/pkg/gittoolbelt"
	"net/http"
	"os"
	"time"
)

const (
	maxRequestTimeout         = 5 * time.Second
	publishChgLogWorkflow     = "publish-changelog"
	publishReleaseTagWorkflow = "publish-release-tag"
	workflowSelector          = "workflow-selector"
)

type appFlags struct {
	help           bool
	version        bool
	CommitHash     string
	CurrentVersion string
}

// envVars Values pulled from their environment variables equivalent. See GetRequiredEnvVars
type envVars map[string]string

var af = &appFlags{}

func init() {
	flag.BoolVar(&af.help, "help", false, um["help"])
	flag.BoolVar(&af.version, "version", false, um["version"])
}

func main() {
	var mainErr error

	defer func() {
		if mainErr != nil {
			log.Errf("%v", mainErr.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}()

	flag.Parse()

	if af.help {
		flag.Usage()
		return
	}

	if af.version {
		log.Logf(stdout.CurrentVersion, af.CurrentVersion, af.CommitHash)
		return
	}

	ca := flag.Args()

	if len(ca) < 2 {
		fmt.Println("nothing to do, bye!")
		os.Exit(0)
		return
	}

	client := &http.Client{
		Timeout: maxRequestTimeout,
	}

	switch ca[0] {
	case publishReleaseTagWorkflow:
		log.Logf(stdout.StartWorkflow, publishChgLogWorkflow)
		// Grab all the environment variables and alert if any are not set.
		eVars, err1 := getRequiredEnvVars([]string{
			"CIRCLE_TOKEN",
			"GH_TOKEN",
			"CIRCLE_REPOSITORY_URL",
			"PARAM_GH_SERVER",
		})
		if err1 != nil {
			mainErr = err1
			return
		}

		if len(ca) < 3 {
			log.Logf(stderr.PublishReleaseTagArgs)
			os.Exit(1)
			return
		}

		branch := ca[2]
		wd := ca[3]

		gh := github.NewClient(eVars["CIRCLE_REPOSITORY_URL"], eVars["GH_TOKEN"], eVars["PARAM_GH_SERVER"], client)
		wf := NewWorkflow(eVars["CIRCLE_TOKEN"], gh)

		mainErr = wf.PublishReleaseTag(branch, wd)

	case publishChgLogWorkflow:
		log.Logf(stdout.StartWorkflow, publishChgLogWorkflow)

		// Grab all the environment variables and alert if any are not set.
		eVars, err1 := getRequiredEnvVars([]string{
			"CIRCLE_REPOSITORY_URL",
			"CIRCLE_TOKEN",
			"CIRCLE_USERNAME",
			"GH_TOKEN",
			"PARAM_GH_SERVER",
			"PARAM_GIT_CHGLOG_CONFIG_FILE",
			"PARAM_MERGE_TYPE",
		})
		if err1 != nil {
			mainErr = err1
			return
		}

		if len(ca) < 3 {
			log.Logf(stderr.PublishChangelogArgs)
			os.Exit(1)
			return
		}

		chgLogFile := ca[1]
		branch := ca[2]
		wd := ca[3]

		gh := github.NewClient(eVars["CIRCLE_REPOSITORY_URL"], eVars["GH_TOKEN"], eVars["PARAM_GH_SERVER"], client)
		gh.MergeMethod = eVars["PARAM_MERGE_TYPE"]
		gh.Username = eVars["CIRCLE_USERNAME"]
		gh.ChglogConfigFile = eVars["PARAM_GIT_CHGLOG_CONFIG_FILE"]
		wf := NewWorkflow(eVars["CIRCLE_TOKEN"], gh)

		mainErr = wf.PublishChangelog(wd, chgLogFile, branch)

	case workflowSelector:
		log.Logf(stdout.StartWorkflow, workflowSelector)

		// Step 1: Grab all the environment variables and alert if any are not
		// set. See https://circleci.com/docs/variables/#built-in-environment-variables
		eVars, err1 := getRequiredEnvVars([]string{
			"CIRCLE_TOKEN",
			"CIRCLE_PROJECT_REPONAME",
			"CIRCLE_PROJECT_USERNAME",
			"PARAM_CIRCLECI_API_HOST",
			"PARAM_CIRCLECI_APP_HOST",
			"PARAM_VCS_TYPE",
		})
		if err1 != nil {
			mainErr = err1
			return
		}

		if len(ca) < 5 {
			mainErr = fmt.Errorf(stderr.FiveArgsRequired)
			return
		}

		chgLogFile := ca[1]
		branch := ca[2]
		wd := ca[3]
		commit := ca[4]

		// Step 2: When the changelog has updates, then trigger the changelog
		// workflow and return.
		thereAreChanges, err5 := IsChangelogUpToDate(wd, chgLogFile)
		if err5 != nil {
			mainErr = err5
			return
		}

		if thereAreChanges {
			// Step 2.a: Set pipeline parameters and trigger the workflow: publish changelog.
			pp, e1 := circleci.GetPipelineParameters(branch, publishChgLogWorkflow)
			if e1 != nil {
				mainErr = e1
				return
			}

			mainErr = circleci.TriggerPipeline(pp, client, eVars)

			return
		}

		// Step 4: No changelog updates, then verify the commit is not tagged.
		if IsCommitTagged(wd, commit) {
			log.Logf(stderr.CommitAlreadyTagged, commit)
			return
		}

		// Step 5: Verify that the range of commits contain a message to
		// indicate they should be tagged.
		if !gittoolbelt.IsTaggable(wd) {
			log.Logf(stdout.NoCommitsToTag)
			return
		}

		// Step 6: Build pipeline parameters to trigger the tag-and-release
		// workflow
		pp, errY1 := circleci.GetPipelineParameters(branch, publishReleaseTagWorkflow)
		if errY1 != nil {
			mainErr = errY1
			return
		}

		// Step 7: Trigger the workflow
		mainErr = circleci.TriggerPipeline(pp, client, eVars)
	}
}

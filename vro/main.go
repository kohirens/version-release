//go:generate git-tool-belt semver -save info.go -format go -packageName main -varName af

package main

import (
	"flag"
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release/vro/internal/util"
	"github.com/kohirens/version-release/vro/pkg/circleci"
	"github.com/kohirens/version-release/vro/pkg/git"
	"github.com/kohirens/version-release/vro/pkg/gitcliff"
	"github.com/kohirens/version-release/vro/pkg/github"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	autoReleaseHeader         = "auto: Release %v"
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
	SemVer         string
	TagAndRelease  struct {
		Flags  flag.FlagSet
		Help   bool
		SemVer string
	}
}

// envVars Values pulled from their environment variables equivalent. See GetRequiredEnvVars
type envVars map[string]string

var af = &appFlags{}

func init() {
	flag.BoolVar(&af.help, "help", false, um["help"])
	flag.BoolVar(&af.version, "version", false, um["version"])
	flag.StringVar(&af.SemVer, "semver", "", um["semver"])

	af.TagAndRelease.Flags = flag.FlagSet{}
	af.TagAndRelease.Flags.BoolVar(
		&af.TagAndRelease.Help,
		"help",
		false,
		um["tag_and_release_help"],
	)
	af.TagAndRelease.Flags.StringVar(
		&af.TagAndRelease.SemVer,
		"semver",
		"",
		um["tag_and_release_semver"],
	)
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

	vl, vlFound := os.LookupEnv("VERBOSITY_LEVEL")
	if vlFound {
		num, _ := strconv.ParseInt(vl, 10, 64)
		log.VerbosityLevel = int(num)
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

	semVer := ""
	if af.SemVer != "" {
		if regexp.MustCompile(util.CheckSemVer).MatchString(af.SemVer) {
			semVer = af.SemVer
		} else {
			mainErr = fmt.Errorf(stderr.InvalidSemVer, af.SemVer)
			return
		}
	}
	log.Infof("semantic version set to %v", semVer)

	switch ca[0] {
	case publishReleaseTagWorkflow:
		if e := af.TagAndRelease.Flags.Parse(ca[1:]); e != nil {
			af.TagAndRelease.Flags.Usage()
			mainErr = e
			return
		}

		if af.TagAndRelease.Help {
			af.TagAndRelease.Flags.PrintDefaults()
			return
		}

		log.Logf(stdout.StartWorkflow, publishReleaseTagWorkflow)

		// Grab all the environment variables and alert if any are not set.
		eVars, e1 := getRequiredEnvVars([]string{
			"CIRCLE_TOKEN",
			"GH_TOKEN",
			"CIRCLE_REPOSITORY_URL",
			"PARAM_GH_SERVER",
		})
		if e1 != nil {
			mainErr = e1
			return
		}

		subCa := af.TagAndRelease.Flags.Args()

		if len(subCa) < 2 {
			mainErr = fmt.Errorf(stderr.PublishReleaseTagArgs)
			return
		}

		branch := subCa[0]
		wd := subCa[1]

		log.Infof("branch %v", branch)
		log.Infof("working directory %v", wd)

		if af.TagAndRelease.SemVer != "" {
			if regexp.MustCompile(util.CheckSemVer).MatchString(af.TagAndRelease.SemVer) {
				semVer = af.TagAndRelease.SemVer
			} else {
				mainErr = fmt.Errorf(stderr.InvalidSemVer, af.TagAndRelease.SemVer)
				return
			}
		}

		gh := github.NewClient(eVars["CIRCLE_REPOSITORY_URL"], eVars["GH_TOKEN"], eVars["PARAM_GH_SERVER"], client)

		wf := NewWorkflow(eVars["CIRCLE_TOKEN"], gh)

		nextVer, e2 := nextVersion(semVer, wd)
		if e2 != nil {
			mainErr = e2
			return
		}

		mainErr = wf.PublishReleaseTag(branch, nextVer)

	case publishChgLogWorkflow:
		log.Logf(stdout.StartWorkflow, publishChgLogWorkflow)

		// Grab all the environment variables and alert if any are not set.
		eVars, err1 := getRequiredEnvVars([]string{
			"CIRCLE_REPOSITORY_URL",
			"CIRCLE_TOKEN",
			"CIRCLE_USERNAME",
			"GH_TOKEN",
			"PARAM_GH_SERVER",
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
		wf := NewWorkflow(eVars["CIRCLE_TOKEN"], gh)

		mainErr = wf.PublishChangelog(wd, chgLogFile, branch, semVer)

	case workflowSelector:
		log.Logf(stdout.StartWorkflow, workflowSelector)

		// Step 1: Grab all the environment variables and alert if any are not
		// set. See https://circleci.com/docs/variables/#built-in-environment-variables
		eVars, e1 := getRequiredEnvVars([]string{
			"CIRCLE_TOKEN",
			"CIRCLE_PROJECT_REPONAME",
			"CIRCLE_PROJECT_USERNAME",
			"PARAM_CIRCLECI_API_HOST",
			"PARAM_CIRCLECI_APP_HOST",
			"PARAM_VCS_TYPE",
		})
		if e1 != nil {
			mainErr = e1
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

		if !git.IsCommit(wd, commit) {
			mainErr = fmt.Errorf(stderr.InvalidCommit, commit)
			return
		}
		hasSemverTag := git.HasSemverTag(wd, commit)

		// NOTE: nextVersion is equivalent to this check, so does it make sense to run this as it seems to be no benefit.
		if hasSemverTag { // Do nothing when the commit is semver tagged.
			log.Logf(stderr.CommitAlreadyTagged, commit)
		}

		// only consider tagging if HEAD has no tag and the commit message
		// contains the expected auto-release header.
		if !hasSemverTag {
			nextVer, e2 := nextVersion(semVer, wd)
			if e2 != nil {
				mainErr = e2
				return
			}

			// Skip commits that are NOT released and also should NOT be tagged.
			if !strings.Contains(git.Log(wd, commit), fmt.Sprintf(autoReleaseHeader, nextVer)) {
				return
			}

			// Build pipeline parameters to trigger the tag-and-release workflow.
			pp, e3 := circleci.GetPipelineParameters(branch, publishReleaseTagWorkflow)
			if e3 != nil {
				mainErr = e3
				return
			}

			log.Logf(stdout.TriggerWorkflow, publishReleaseTagWorkflow)

			//  Trigger the workflow
			mainErr = circleci.TriggerPipeline(pp, client, eVars)
			return
		}

		hasUnreleasedChanges, e4 := gitcliff.UnreleasedChanges(wd)
		if e4 != nil {
			mainErr = e4
			return
		}

		if len(hasUnreleasedChanges) > 0 {
			// Scan the changelog to verify it does not already contain the unreleased entries.
			// git-cliff just blindly prepends commit to the CHANGELOG without verify they were already added. So we want to prevent duplicate entries.
			if fsio.Exist(chgLogFile) {
				// Exit when the change is update-to-date or an error occurred
				if containUnreleased, e5 := changelogContains(&hasUnreleasedChanges[0], wd, chgLogFile); containUnreleased || e5 != nil {
					mainErr = e5
					return
				}
			}

			// Build pipeline parameters for the publish-changelog workflow.
			pp, e6 := circleci.GetPipelineParameters(branch, publishChgLogWorkflow)
			if e6 != nil {
				mainErr = e6
				return
			}

			log.Logf(stdout.TriggerWorkflow, publishChgLogWorkflow)

			// Trigger the publish-changelog workflow.
			mainErr = circleci.TriggerPipeline(pp, client, eVars)

			return
		}
	}
}

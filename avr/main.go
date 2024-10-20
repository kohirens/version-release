//go:generate git-tool-belt semver -save info.go -format go -packageName main -varName clo

package main

import (
	"flag"
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release/avr/pkg/circleci"
	"github.com/kohirens/version-release/avr/pkg/git"
	"github.com/kohirens/version-release/avr/pkg/github"
	"github.com/kohirens/version-release/avr/pkg/lib"
	"github.com/kohirens/version-release/vro/pkg/gitcliff"
	"net/http"
	"os"
	"path/filepath"
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

var clo = &commandLineOptions{}

func init() {
	defineOptions(clo)
	definePublishChangelogOptions(clo)
	defineTagAndReleaseOptions(clo)
	defineWorkflowSelectorOptions(clo)
	defineKnownSshKeysOptions(clo)
}

func main() {
	var mainErr error

	// Run this when this function returns.
	defer func() {
		if mainErr != nil {
			log.Errf("%v", mainErr.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}()

	flag.Parse()

	// Display the help message of this executable and exit 0.
	if clo.help {
		flag.Usage()
		return
	}

	// Print the current version of this executable and exit 0.
	if clo.version {
		log.Logf(stdout.CurrentVersion, clo.CurrentVersion, clo.CommitHash)
		return
	}

	// Set the logging level via an environment variable.
	vl, vlFound := os.LookupEnv("VERBOSITY_LEVEL")
	if vlFound {
		num, _ := strconv.ParseInt(vl, 10, 64)
		log.VerbosityLevel = int(num)
	}

	// Get all command line arguments.
	cla := flag.Args()

	// Exit 0 if you do not pass any command line arguments.
	if len(cla) < 1 {
		fmt.Println(stdout.Nothing)
		os.Exit(0)
		return
	}

	semVer := ""
	if clo.SemVer != "" {
		if regexp.MustCompile(CheckSemVer).MatchString(clo.SemVer) {
			semVer = clo.SemVer
		} else {
			mainErr = fmt.Errorf(stderr.InvalidSemVer, clo.SemVer)
			return
		}
		log.Infof(stdout.SemVer, semVer)
	}

	log.Dbugf(stdout.Wd, clo.WorkDir)
	// clean up the working directory
	workDir, e1 := filepath.Abs(clo.WorkDir)
	if e1 != nil {
		mainErr = fmt.Errorf(stderr.WorkDir, e1.Error())
		return
	}
	log.Dbugf(stdout.Wd, workDir)

	// An HTTP client is also required for everything below.
	client := &http.Client{
		Timeout: maxRequestTimeout,
	}

	if e := git.ConfigGlobal(workDir, "safe.directory", workDir); e != nil {
		mainErr = e
		return
	}

	switch cla[0] {
	case publishReleaseTagWorkflow:
		if e := clo.TagAndRelease.Flags.Parse(cla[1:]); e != nil {
			clo.TagAndRelease.Flags.Usage()
			mainErr = e
			return
		}

		// Print subcommand help and exit 0.
		if clo.TagAndRelease.Help {
			clo.TagAndRelease.Flags.PrintDefaults()
			return
		}

		// Let the user know we are starting down the tag and release subcommand.
		log.Logf(stdout.StartWorkflow, publishReleaseTagWorkflow)

		branch := clo.Branch

		log.Infof(stdout.Branch, branch)

		var gh *github.Client

		switch clo.CiCd {
		case circleci.Name:
			eVars, e2 := lib.GetRequiredEnvVars([]string{
				github.EnvToken,
				circleci.EnvProjectRepoName,
				circleci.EnvProjectUsername,
				github.EnvApiUrl,
			})
			if e2 != nil {
				mainErr = e2
				return
			}
			gh = github.NewClient(
				eVars[circleci.EnvProjectUsername]+"/"+eVars[circleci.EnvProjectRepoName],
				eVars[github.EnvToken],
				eVars[github.EnvApiUrl],
				client,
			)
		case github.Name:
			eVars, e2 := lib.GetRequiredEnvVars([]string{
				github.EnvApiUrl,
				github.EnvRepository,
				github.EnvRepositoryOwner,
				github.EnvToken,
			})
			if e2 != nil {
				mainErr = e2
				return
			}
			gh = github.NewClient(
				eVars[github.EnvRepository],
				eVars[github.EnvToken],
				eVars[github.EnvApiUrl],
				client,
			)
		}

		wf := NewWorkflow(gh)

		nextVer, e3 := nextVersion(semVer, workDir)
		if e3 != nil {
			mainErr = e3
			return
		}

		// Publish a new tag on GitHub.
		ghRelease, e4 := wf.GitHubClient.TagAndRelease(branch, nextVer)
		if e4 != nil {
			mainErr = e4
			return
		}

		log.Logf(stdout.ReleaseTag, ghRelease.Name)
		mainErr = wf.PublishReleaseTag(branch, nextVer)

	case publishChgLogWorkflow:
		//# publish-changelog <path-to-changelog> [<merge-type>]
		//# publish-changelog CHANGELOG.md rebase
		if e := clo.PublishChangelog.Flags.Parse(cla[1:]); e != nil {
			clo.PublishChangelog.Flags.Usage()
			mainErr = e
			return
		}

		// Print subcommand help and exit 0.
		if clo.PublishChangelog.Help {
			clo.PublishChangelog.Flags.PrintDefaults()
			return
		}

		log.Logf(stdout.StartWorkflow, publishChgLogWorkflow)

		github.PublicServer = clo.PublishChangelog.GitHubServer

		// Get command line arguments.
		subCla := clo.PublishChangelog.Flags.Args()
		log.Dbugf("subCla: %v", subCla)
		if len(subCla) < 1 {
			mainErr = fmt.Errorf(stderr.PublishChangelogArgs)
			return
		}

		changelog := subCla[0]

		mergeType, e3 := validateMergeType(clo.PublishChangelog.MergeType)
		if e3 != nil {
			mainErr = e3
			return
		}

		nextVer, e2 := nextVersion(semVer, workDir)
		if e2 != nil {
			mainErr = e2
			return
		}

		gitHubToken := lib.GetEnv(github.EnvToken)
		gitHubApiUrl := lib.GetEnv(github.EnvApiUrl)

		var gh *github.Client
		var repo string

		switch clo.CiCd {
		case circleci.Name:
			log.Logf(stdout.CciChangelog, clo.CiCd)

			eVars, err1 := lib.GetRequiredEnvVars([]string{
				circleci.EnvProjectRepoName,
				circleci.EnvToken,
				circleci.EnvProjectUsername,
			})
			if err1 != nil {
				mainErr = err1
				return
			}

			repo = eVars[circleci.EnvProjectUsername] + "/" + eVars[circleci.EnvProjectRepoName]
		case github.Name:
			log.Logf(stdout.GaChangelog)

			eVars, err1 := lib.GetRequiredEnvVars([]string{
				github.EnvActor,
				github.EnvRepository,
				github.EnvRepositoryOwner,
			})
			if err1 != nil {
				mainErr = err1
				return
			}
			repo = eVars[github.EnvRepository]
		}

		gh = github.NewClient(repo, gitHubToken, gitHubApiUrl, client)
		gh.MergeMethod = mergeType
		wf := NewWorkflow(gh)
		mainErr = wf.PublishChangelog(workDir, changelog, clo.Branch, nextVer)

	case workflowSelector:
		// For GitHib Actions default to none
		if clo.CiCd == github.Name {
			if e := github.AddOutputVar("workflow", "none"); e != nil {
				mainErr = e
				return
			}
		}

		if e := clo.WorkflowSelector.Flags.Parse(cla[1:]); e != nil {
			clo.WorkflowSelector.Flags.Usage()
			mainErr = e
			return
		}

		// Print subcommand help and exit 0.
		if clo.WorkflowSelector.Help {
			clo.WorkflowSelector.Flags.PrintDefaults()
			return
		}

		log.Logf(stdout.StartWorkflow, workflowSelector)

		// Get command line arguments.
		subCla := clo.WorkflowSelector.Flags.Args()
		log.Dbugf("subCla: %v", subCla)
		if len(subCla) < 2 {
			mainErr = fmt.Errorf(stderr.WorkflowSelectorInput)
			return
		}

		changelog := subCla[0]
		commit := subCla[1]

		// Validate the commit is in the history of this repository. For times when you rebase and run a workflow where the commit was removed.
		if !git.IsCommit(workDir, commit) {
			mainErr = fmt.Errorf(stderr.InvalidCommit, commit)
			return
		}

		hasSemverTag := git.HasSemverTag(workDir, commit)

		// Log that the commit already has a tag.
		if hasSemverTag {
			log.Logf(stderr.CommitAlreadyTagged, commit)
		}

		// Only consider tagging if HEAD has no tag and the commit message
		// contains the expected auto-release header.
		if !hasSemverTag {
			nextVer, e2 := nextVersion(semVer, workDir)
			if e2 != nil { // No version to tag, then check for changelog updates.
				log.Logf(e2.Error())
				goto changLog
			}

			l := git.Log(workDir, commit)
			log.Dbugf(stdout.DbgCommitLog, l)

			// Skip tagging when the latest commit is NOT a changelog update.
			if !strings.Contains(l, fmt.Sprintf(autoReleaseHeader, nextVer)) {
				goto changLog
			}

			// Start the "publish a release" workflow based on the platform.
			switch clo.CiCd {
			case circleci.Name:
				cci := circleci.NewClient("gh", client)
				mainErr = cci.RunWorkflow(clo.Branch, publishReleaseTagWorkflow)
			case github.Name:
				log.Logf("trigger a tag_and_release workflow on GitHub Actions")
				// For GitHub Actions we merely need to set an output variable to continue onto the next workflow.
				mainErr = github.AddOutputVar("workflow", publishReleaseTagWorkflow)
			default:
				mainErr = fmt.Errorf("invlaid platform %v", clo.CiCd)
			}
			return
		}

	changLog:
		hasUnreleasedChanges, e4 := gitcliff.UnreleasedChanges(workDir)
		if e4 != nil {
			mainErr = e4
			return
		}

		if len(hasUnreleasedChanges) > 0 {
			// Scan the changelog to verify it does not already contain the unreleased entries.
			// git-cliff just blindly prepends commit to the CHANGELOG without verify they were already added. So we want to prevent duplicate entries.
			if fsio.Exist(changelog) {
				// Exit when the changelog is update-to-date or an error occurred
				// TODO: Fix the name of this function it is not clear what it does.
				isUpToDate, e5 := changelogContains(&hasUnreleasedChanges[0], workDir, changelog)
				if e5 != nil {
					mainErr = e5
					return
				}

				if isUpToDate {
					log.Logf(stdout.ChgLogUpToDate)
					return
				}
			}

			switch clo.CiCd {
			case circleci.Name:
				cci := circleci.NewClient("gh", client)
				// Trigger the publish-changelog workflow.
				mainErr = cci.RunWorkflow(clo.Branch, publishChgLogWorkflow)
			case github.Name:
				log.Dbugf("triggered a publish-changelog workflow on GitHub Actions")
				// For GitHub Actions we merely need to set an output variable to continue onto the next workflow.
				mainErr = github.AddOutputVar("workflow", publishChgLogWorkflow)
			}

			return
		}
	case "known-sshkeys":
		subCla := clo.KnownSshKeys.Flags.Args()
		log.Dbugf("subCla: %v", subCla)
		if len(subCla) < 2 {
			mainErr = fmt.Errorf(stderr.KnownSshKeys)
			return
		}
		known, e2 := github.KnownSshKeys(subCla[0], subCla[1], client)
		if e2 != nil {
			mainErr = e2
			return
		}

		fmt.Print(known)
	}
}

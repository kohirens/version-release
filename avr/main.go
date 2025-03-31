//go:generate git-tool-belt semver -save info.go -format go -packageName main -varName clo

package main

import (
	"flag"
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release/avr/pkg/circleci"
	"github.com/kohirens/version-release/avr/pkg/git"
	"github.com/kohirens/version-release/avr/pkg/gitcliff"
	"github.com/kohirens/version-release/avr/pkg/github"
	"github.com/kohirens/version-release/avr/pkg/lib"
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
	autoReleaseFooter         = "changes sha256 %v"
	maxRequestTimeout         = 10 * time.Second
	publishChgLogWorkflow     = "publish-changelog"
	publishReleaseTagWorkflow = "publish-release-tag"
	workflowSelector          = "workflow-selector"
)

var (
	clo  = &commandLineOptions{}
	logs = log.StdLogger{}
)

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
			logs.Errf("%v", mainErr.Error())
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
		logs.Logf(stdout.CurrentVersion, clo.CurrentVersion, clo.CommitHash)
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
		if regexp.MustCompile(lib.CheckSemVer).MatchString(clo.SemVer) {
			semVer = clo.SemVer
		} else {
			mainErr = fmt.Errorf(stderr.InvalidSemVer, clo.SemVer)
			return
		}
		logs.Infof(stdout.SemVer, semVer)
	}

	// clean up the working directory
	workDir, e1 := filepath.Abs(clo.WorkDir)
	if e1 != nil {
		mainErr = fmt.Errorf(stderr.WorkDir, e1.Error())
		return
	}

	logs.Dbugf(stdout.WorkDir, workDir)

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

		// Let the user know we are starting down the tag and release path.
		logs.Logf(stdout.StartWorkflow, publishReleaseTagWorkflow)

		branch := clo.Branch

		logs.Infof(stdout.Branch, branch)

		gh, e2 := newGitHubClient(client)
		if e2 != nil {
			mainErr = e2
			return
		}

		// required tag for the release
		nextVer := gitcliff.NextVersion(workDir, semVer)
		if nextVer == "" { // this is definitely an error at this point.
			mainErr = fmt.Errorf(stderr.NoSemverTag)
			return
		}

		// Publish a new tag on GitHub.
		ghRelease, e4 := gh.TagAndRelease(branch, nextVer)
		if e4 != nil {
			mainErr = e4
			return
		}

		logs.Logf(stdout.ReleaseTag, ghRelease.Name)

	case publishChgLogWorkflow:
		// publish-changelog <path-to-changelog> [<merge-type>]
		// publish-changelog CHANGELOG.md rebase
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

		logs.Logf(stdout.StartWorkflow, publishChgLogWorkflow)

		github.PublicServer = clo.PublishChangelog.GitHubServer

		// Get command line arguments.
		subCla := clo.PublishChangelog.Flags.Args()
		logs.Dbugf(stdout.SubCla, subCla)
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

		// required semver tag for the changelog header
		nextVer := gitcliff.NextVersion(workDir, semVer)
		if nextVer == "" {
			mainErr = fmt.Errorf(stderr.NoSemverTag)
			return
		}

		gitHubToken := lib.GetEnv(github.EnvToken)
		logs.Dbugf("token: ***%v", gitHubToken[len(gitHubToken)-5:])
		gitHubApiUrl := lib.GetEnv(github.EnvApiUrl)

		var gh *github.Client
		var repo string

		switch clo.CiCd {
		case circleci.Name:
			logs.Logf(stdout.CciChangelog, clo.CiCd)

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
			logs.Logf(stdout.GaChangelog)

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
		mainErr = PublishChangelog(workDir, changelog, clo.Branch, nextVer, gh)

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

		logs.Logf(stdout.StartWorkflow, workflowSelector)

		// Get command line arguments.
		subCla := clo.WorkflowSelector.Flags.Args()
		logs.Dbugf(stdout.SubCla, subCla)
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
			logs.Logf(stderr.CommitAlreadyTagged, commit)
		}

		// Only consider tagging if:
		// 1. HEAD has no tag
		// 2. the commit message contains the expected auto-release header.
		// 3. There are conventional commits to tag.
		if !hasSemverTag {
			nextVer := gitcliff.NextVersion(workDir, semVer)
			if nextVer == "" { // No version to tag, then check for changelog updates.
				logs.Logf(stdout.NoChangesToTag)
				goto changLog
			}

			l := git.Log(workDir, commit)
			logs.Dbugf(stdout.DbgCommitLog, l)

			// Skip tagging when the latest commit is NOT a changelog update.
			if !strings.Contains(l, fmt.Sprintf(autoReleaseHeader, nextVer)) {
				goto changLog
			}

			// Start a "publish a release" workflow based on the platform.
			switch clo.CiCd {
			case circleci.Name:
				cci := circleci.NewClient("gh", client)
				mainErr = cci.RunWorkflow(clo.Branch, publishReleaseTagWorkflow)
			case github.Name:
				logs.Logf("trigger a tag_and_release workflow on GitHub Actions")
				// For GitHub Actions we merely need to set an output variable to continue onto the next workflow.
				mainErr = github.AddOutputVar("workflow", publishReleaseTagWorkflow)
				mainErr = github.AddOutputVar("next_semver", nextVer)
				mainErr = github.AddOutputVar("changelog_hash", "")
				_ = github.DumpOutput()
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

		// verify there are no unreleased changes.
		if len(hasUnreleasedChanges) == 0 {
			logs.Logf(stdout.NoWorkflowSelected)
			return
		}

		// Scan the changelog to verify it does not already contain the
		// unreleased entries.
		// git-cliff blindly prepends commits to the CHANGELOG without
		// verifying if they were already added. So we want to prevent duplicate
		// entries.
		var changeLogHash64 string

		if fsio.Exist(changelog) {
			// Exit when the changelog is update-to-date or an error occurred
			isUpToDate, hash64, e5 := isChangelogCurrent(workDir, changelog)
			if e5 != nil {
				mainErr = e5
				return
			}
			changeLogHash64 = hash64

			if isUpToDate {
				// Verify the tag has been published.
				gh, e6 := newGitHubClient(client)
				if e6 != nil {
					mainErr = e6
					return
				}

				nextVer := gitcliff.NextVersion(workDir, semVer)
				if nextVer == "" { // No version to tag, then check for changelog updates.
					logs.Logf(stdout.NoChangesToTag)
					return
				}

				rr, _ := gh.ReleaseByTag(nextVer)
				if rr == nil {
					logs.Logf(stdout.ChgLogUpToDate)
					return
				}
			}
		}

		switch clo.CiCd {
		case circleci.Name:
			cci := circleci.NewClient("gh", client)
			// Trigger the publish-changelog workflow.
			mainErr = cci.RunWorkflow(clo.Branch, publishChgLogWorkflow)
		case github.Name:
			logs.Dbugf(stdout.GhPublishChgLog)
			// For GitHub Actions we merely need to set an output variable
			// to continue onto the next workflow.
			mainErr = github.AddOutputVar("workflow", publishChgLogWorkflow)
			mainErr = github.AddOutputVar("changelog_hash", changeLogHash64)
			mainErr = github.AddOutputVar("next_semver", "")
			_ = github.DumpOutput()
		}

	case "known-sshkeys":
		subCla := clo.KnownSshKeys.Flags.Args()
		logs.Dbugf(stdout.SubCla, subCla)
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

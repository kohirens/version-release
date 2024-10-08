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

type commandLineOptions struct {
	help           bool
	version        bool
	Branch         string
	CiCd           string
	CommitHash     string
	CurrentVersion string
	GitHubServer   string
	GitHubToken    string
	Project        string
	SemVer         string
	Username       string
	WorkDir        string
	TagAndRelease  struct {
		Flags flag.FlagSet
		Help  bool
	}
	WorkflowSelector struct {
		Flags         flag.FlagSet
		ChangeLogFile string
		Help          bool
	}
}

// envVars Values pulled from their environment variables equivalent. See GetRequiredEnvVars
type envVars map[string]string

var clo = &commandLineOptions{}

func init() {
	flag.BoolVar(&clo.help, "help", false, um["help"])
	flag.BoolVar(&clo.version, "version", false, um["version"])
	flag.StringVar(&clo.Branch, "branch", "main", um["branch"])
	flag.StringVar(&clo.CiCd, "cicd", circleci.Name, um["cicd"])
	flag.StringVar(&clo.GitHubServer, "gh-server", github.Server, um["gh_server"])
	flag.StringVar(&clo.GitHubToken, "gh-token", "", um["gh_token"])
	flag.StringVar(&clo.Project, "project", "", um["project"])
	flag.StringVar(&clo.SemVer, "semver", "", um["semver"])
	flag.StringVar(&clo.Username, "username", "", um["username"])
	flag.StringVar(&clo.WorkDir, "wd", ".", um["wd"])

	clo.TagAndRelease.Flags = flag.FlagSet{}
	clo.TagAndRelease.Flags.BoolVar(
		&clo.TagAndRelease.Help,
		"help",
		false,
		um["tag_and_release_help"],
	)

	clo.WorkflowSelector.Flags = flag.FlagSet{}
	clo.WorkflowSelector.Flags.BoolVar(
		&clo.WorkflowSelector.Help,
		"help",
		false,
		um["workflow_selector_help"],
	)
	clo.WorkflowSelector.Flags.StringVar(
		&clo.WorkflowSelector.ChangeLogFile,
		"changelog",
		"CHANGELOG.md",
		um["workflow_selector_help"],
	)
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

	// Validate GitHub Server was set.
	if clo.GitHubServer == "" {
		mainErr = fmt.Errorf(stderr.GitHubServer)
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
	clo.WorkDir = workDir
	log.Dbugf(stdout.Wd, workDir)

	// An HTTP client is also required for everything below.
	client := &http.Client{
		Timeout: maxRequestTimeout,
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

		// Grab required environment variables, if any are not set then exit 1.
		eVars, e1 := getRequiredEnvVars([]string{
			github.EnvToken,
			circleci.EnvRepoUrl,
			github.EnvServer,
		})
		if e1 != nil {
			mainErr = e1
			return
		}

		subCa := clo.TagAndRelease.Flags.Args()

		if len(subCa) < 2 {
			mainErr = fmt.Errorf(stderr.PublishReleaseTagArgs)
			return
		}

		branch := clo.Branch
		wd := clo.WorkDir

		log.Infof(stdout.Branch, branch)
		log.Infof(stdout.Wd, wd)

		gh := github.NewClient(eVars[circleci.EnvRepoUrl], clo.GitHubToken, eVars[github.EnvServer], client)

		wf := NewWorkflow(eVars[circleci.EnvToken], gh)

		nextVer, e2 := nextVersion(semVer, wd)
		if e2 != nil {
			mainErr = e2
			return
		}

		// Publish a new tag on GitHub.
		ghRelease, e2 := wf.GitHubClient.TagAndRelease(branch, semVer)
		if e2 != nil {
			mainErr = e2
			return
		}

		log.Logf(stdout.ReleaseTag, ghRelease.Name)
		mainErr = wf.PublishReleaseTag(branch, nextVer)

	case publishChgLogWorkflow:
		log.Logf(stdout.StartWorkflow, publishChgLogWorkflow)

		// Grab all the environment variables and alert if any are not set.
		eVars, err1 := getRequiredEnvVars([]string{
			circleci.EnvRepoUrl,
			circleci.EnvToken,
			circleci.EnvUsername,
			github.EnvToken,
			github.EnvServer,
			github.EnvMergeType,
		})
		if err1 != nil {
			mainErr = err1
			return
		}

		if len(cla) < 1 {
			log.Logf(stderr.PublishChangelogArgs)
			os.Exit(1)
			return
		}

		branch := clo.Branch
		wd := clo.WorkDir
		chgLogFile := cla[0]

		gh := github.NewClient(eVars[circleci.EnvRepoUrl], clo.GitHubToken, eVars[github.EnvServer], client)
		gh.MergeMethod = eVars[github.EnvMergeType]
		gh.Username = eVars[circleci.EnvUsername]
		wf := NewWorkflow(eVars[circleci.EnvToken], gh)

		nextVer, e2 := nextVersion(semVer, wd)
		if e2 != nil {
			mainErr = e2
			return
		}

		mainErr = wf.PublishChangelog(wd, chgLogFile, branch, nextVer)

	case workflowSelector:
		log.Logf(stdout.StartWorkflow, workflowSelector)

		// For GitHib Actions default to none
		if clo.CiCd == github.Name {
			if e := os.Setenv("workflow", "none"); e != nil {
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

		// Get command line arguments.
		subCla := clo.WorkflowSelector.Flags.Args()
		log.Dbugf("subCla: %v", subCla)
		if len(subCla) < 1 {
			mainErr = fmt.Errorf(stderr.WorkflowSelectorInput)
			return
		}

		commit := subCla[0]

		// Get command line options.
		wd := clo.WorkDir
		chgLogFile := clo.WorkflowSelector.ChangeLogFile

		// Validate the commit is in the history of this repository. For times when you rebase and run a workflow where the commit was removed.
		if !git.IsCommit(wd, commit) {
			mainErr = fmt.Errorf(stderr.InvalidCommit, commit)
			return
		}

		hasSemverTag := git.HasSemverTag(wd, commit)

		// Log that the commit already has a tag.
		if hasSemverTag {
			log.Logf(stderr.CommitAlreadyTagged, commit)
		}

		// Only consider tagging if HEAD has no tag and the commit message
		// contains the expected auto-release header.
		if !hasSemverTag {
			nextVer, e2 := nextVersion(semVer, wd)
			if e2 != nil { // No version to tag, then check for changelog updates.
				log.Logf(e2.Error())
				goto changLog
			}

			l := git.Log(wd, commit)
			log.Dbugf(stdout.DbgCommitLog, l)

			// Skip commits that are NOT released and also should NOT be tagged.
			if !strings.Contains(l, fmt.Sprintf(autoReleaseHeader, nextVer)) {
				goto changLog
			}

			// Start the "publish a release" workflow based on the platform.
			switch clo.CiCd {
			case circleci.Name:
				cci := circleci.NewClient("gh", clo.Project, clo.Username, client)
				mainErr = cci.RunWorkflow(clo.Branch, publishReleaseTagWorkflow)
			case github.Name:
				log.Logf("trigger a tag_and_release workflow on GitHub Actions")
				// For GitHub Actions we merely need to set an output variable to continue onto the next workflow.
				// We set an output variable that will cause the next workflow to execute if its condition is met.
				// The workflow has three possible values: none|publish-changelog|publish-release-tag.
				// These are the same values used to trigger the equivalent CircleCI workflows.
				mainErr = os.Setenv("workflow", publishReleaseTagWorkflow)
			default:
				mainErr = fmt.Errorf("invlaid platform %v", clo.CiCd)
			}
			return
		}

	changLog:
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

			switch clo.CiCd {
			case circleci.Name:
				cci := circleci.NewClient("gh", clo.Project, clo.Username, client)
				// Trigger the publish-changelog workflow.
				mainErr = cci.RunWorkflow(clo.Branch, publishChgLogWorkflow)
			case github.Name:
				log.Logf("trigger a publish-changelog workflow on GitHub Actions")
				// For GitHub Actions we merely need to set an output variable to continue onto the next workflow.
				mainErr = os.Setenv("workflow", publishChgLogWorkflow)
			}

			return
		}
	case "known-sshkeys":
		known, e1 := github.KnownSshKeys(client)
		if e1 != nil {
			mainErr = e1
			return
		}

		fmt.Print(known)
	}
}

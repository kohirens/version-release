// main Define all the command line interface options and functionality here.
package main

import (
	"flag"
	"github.com/kohirens/version-release/avr/pkg/circleci"
	"github.com/kohirens/version-release/avr/pkg/github"
)

type commandLineOptions struct {
	help             bool
	version          bool
	Branch           string
	CiCd             string
	CommitHash       string
	CurrentVersion   string
	EnableTagVPrefix bool
	GitHubAPIURL     string
	GitHubServer     string
	SemVer           string
	WorkDir          string
	PublishChangelog struct {
		Flags     flag.FlagSet
		Help      bool
		MergeType string
	}
	TagAndRelease struct {
		Flags flag.FlagSet
		Help  bool
	}
	WorkflowSelector struct {
		Flags flag.FlagSet
		Help  bool
	}
	KnownSshKeys struct {
		Flags flag.FlagSet
		Help  bool
	}
}

func defineOptions(options *commandLineOptions) {
	flag.BoolVar(&options.help, "help", false, um["help"])
	flag.BoolVar(&options.version, "version", false, um["version"])
	flag.StringVar(&options.Branch, "branch", "main", um["branch"])
	flag.BoolVar(&options.EnableTagVPrefix, "enable-tag-v-prefix", false, um["enable-tag-v-prefix"])
	flag.StringVar(&options.CiCd, "cicd", circleci.Name, um["cicd"])
	flag.StringVar(&options.SemVer, "semver", "", um["semver"])
	flag.StringVar(&options.WorkDir, "wd", ".", um["wd"])
	flag.StringVar(&options.GitHubServer, "github-server", github.Server, um["gh_server"])
	flag.StringVar(&options.GitHubAPIURL, "github-api-url", github.APIURL, um["gh_api_url"])
}

func definePublishChangelogOptions(options *commandLineOptions) {
	options.PublishChangelog.Flags = flag.FlagSet{}
	options.PublishChangelog.Flags.BoolVar(
		&options.PublishChangelog.Help,
		"help",
		false,
		um["help"],
	)
	options.PublishChangelog.Flags.StringVar(
		&options.PublishChangelog.MergeType,
		"merge-type",
		"rebase",
		um["publish_changelog_merge_type"],
	)
}

func defineTagAndReleaseOptions(options *commandLineOptions) {
	options.TagAndRelease.Flags = flag.FlagSet{}
	options.TagAndRelease.Flags.BoolVar(
		&options.TagAndRelease.Help,
		"help",
		false,
		um["help"],
	)
}

func defineWorkflowSelectorOptions(options *commandLineOptions) {
	options.WorkflowSelector.Flags = flag.FlagSet{}
	options.WorkflowSelector.Flags.BoolVar(
		&options.WorkflowSelector.Help,
		"help",
		false,
		um["help"],
	)
}

func defineKnownSshKeysOptions(options *commandLineOptions) {
	options.KnownSshKeys.Flags = flag.FlagSet{}
	options.KnownSshKeys.Flags.BoolVar(
		&options.KnownSshKeys.Help,
		"help",
		false,
		um["help"],
	)
}

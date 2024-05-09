package main

var stderr = struct {
	ChgLogAutoUpdate      string
	ChangelogChanges      string
	CommitAlreadyTagged   string
	CouldNotGetVersion    string
	CouldNotUpdateChglog  string
	FiveArgsRequired      string
	GitStatus             string
	InvalidSemVer         string
	MissingEnv            string
	NoSemverInfo          string
	NoChangelogChanges    string
	NothingToTag          string
	NoVersionTag          string
	OpenFile              string
	PublishChangelogArgs  string
	PublishReleaseTagArgs string
	SemverWithChgLogCheck string
}{
	ChgLogAutoUpdate:      "abort, the last commit contains an auto update to the CHANGELOG.md",
	ChangelogChanges:      "the changelog has changes to be committed",
	CommitAlreadyTagged:   "commit %v is already tagged, so there is nothing to do",
	CouldNotGetVersion:    "could not get version info; %s",
	CouldNotUpdateChglog:  "could not update changelog: %s; %v",
	FiveArgsRequired:      "5 arguments are required to run this command, see -help",
	GitStatus:             "git status failed: %s; %s",
	InvalidSemVer:         "invalid semantic version given %v",
	MissingEnv:            "%s environment variable is not set",
	NoChangelogChanges:    "the changelog has no changes to be committed",
	NoSemverInfo:          "could not get version info; %s",
	NothingToTag:          "next semantic version is empty, so we cannot release a tag",
	NoVersionTag:          "could not find semantic version tag",
	OpenFile:              "could not open file %v: %v",
	PublishChangelogArgs:  "3 arguments are required to run this command, see -help",
	PublishReleaseTagArgs: "3 arguments are required to run this command, see -help",
	SemverWithChgLogCheck: "while checking for changes in the %s; could not get semver info: %s",
}

var stdout = struct {
	Cs             string
	CurrentVersion string
	CurrentVer     string
	FoundChgInFile string
	GitStatus      string
	Match          string
	NoCommitsToTag string
	NextVersion    string
	NoChanges      string
	NoTags         string
	ReleaseTag     string
	StartWorkflow  string
	Wd             string
}{
	Cs:             "cs: %s",
	CurrentVersion: "%v, %v",
	CurrentVer:     "current version %v",
	FoundChgInFile: "running %s produced changes in the %s",
	GitStatus:      "Git status: %v",
	Match:          "entry for %v in the changelog was found, we assume this means the changelog is up-to-date",
	NextVersion:    "next version %v",
	NoChanges:      "no changes to release",
	NoCommitsToTag: "no commits to tag",
	ReleaseTag:     "releasing %v",
	StartWorkflow:  "starting %v workflow",
	Wd:             "wd: %s",
}

var um = map[string]string{
	"help":                        "display this help",
	"version":                     "display version information",
	"tag_and_release_semver_help": "publish-release-tag -help displays this help",
	"tag_and_release_semver":      "provide a semantic version to tag and release",
}

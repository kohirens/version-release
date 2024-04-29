package main

var stderr = struct {
	ChgLogAutoUpdate      string
	ChangelogChanges      string
	CommitAlreadyTagged   string
	CouldNotGetVersion    string
	CouldNotUpdateChglog  string
	FiveArgsRequired      string
	GitStatus             string
	MissingEnv            string
	NoSemverInfo          string
	NoChangelogChanges    string
	PublishChangelogArgs  string
	PublishReleaseTagArgs string
	SemverWithChgLogCheck string
}{
	ChgLogAutoUpdate:      "the last commit contains an auto update to the CHANGELOG.md",
	ChangelogChanges:      "the changelog has changes to be committed",
	CommitAlreadyTagged:   "commit %v is already tagged, so there is nothing to do",
	CouldNotGetVersion:    "could not get version info; %s",
	CouldNotUpdateChglog:  "could not update changelog: %s; %v",
	FiveArgsRequired:      "5 arguments are required to run this command, see -help",
	GitStatus:             "git status failed: %s; %s",
	MissingEnv:            "%s environment variable is not set",
	NoChangelogChanges:    "the changelog has no changes to be committed",
	NoSemverInfo:          "could not get version info; %s",
	PublishChangelogArgs:  "3 arguments are required to run this command, see -help",
	PublishReleaseTagArgs: "3 arguments are required to run this command, see -help",
	SemverWithChgLogCheck: "while checking for changes in the %s; could not get semver info: %s",
}

var stdout = struct {
	Cs             string
	CurrentVersion string
	FoundChgInFile string
	GitStatus      string
	NoCommitsToTag string
	NoTags         string
	ReleaseTag     string
	StartWorkflow  string
	Wd             string
}{
	Cs:             "cs: %s",
	CurrentVersion: "%v, %v",
	FoundChgInFile: "running %s produced changes in the %s",
	GitStatus:      "Git status: %v",
	NoCommitsToTag: "no commits to tag",
	ReleaseTag:     "releasing %v",
	StartWorkflow:  "starting %v workflow",
	Wd:             "wd: %s",
}

var um = map[string]string{
	"help":    "display this help",
	"version": "display version information",
}

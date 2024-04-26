package main

var stderr = struct {
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
	AlreadyReleased string
	Cs              string
	CurrentVersion  string
	FoundChgInFile  string
	GitStatus       string
	NoCommitsToTag  string
	ReleaseTag      string
	StartWorkflow   string
	Wd              string
}{
	AlreadyReleased: "commit %v already has a release tag %s",
	Cs:              "cs: %s",
	CurrentVersion:  "%v, %v",
	FoundChgInFile:  "running %s produced changes in the %s",
	GitStatus:       "Git status: %v",
	NoCommitsToTag:  "no commits to tag",
	ReleaseTag:      "released tag %q",
	StartWorkflow:   "starting %v workflow",
	Wd:              "wd: %s",
}

var um = map[string]string{
	"help":    "display this help",
	"version": "display version information",
}

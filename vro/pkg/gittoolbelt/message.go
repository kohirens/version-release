package gittoolbelt

var stderr = struct {
	CouldNotAddConfig       string
	CouldNotDetermineSemver string
	CouldNotJsonDecode      string
	CouldNotUpdateChgLog    string
}{
	CouldNotAddConfig:       "could not add a git-chglog config: %s; %v\n",
	CouldNotDetermineSemver: "could not determine semantic version for %s: %v\n",
	CouldNotJsonDecode:      "could not decode JSON: %s\n",
	CouldNotUpdateChgLog:    "cannot update changelog: %v\n",
}

var stdout = struct {
	CommitNotTagged string
	Co              string
	Cs              string
}{
	CommitNotTagged: "no tag for commit found: %s; %v\n",
	Co:              "co: %s\n",
	Cs:              "cs: %s\n",
}

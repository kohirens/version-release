package gitcliff

var stderr = struct {
	BumpedVersion,
	CannotCalcHash,
	CannotDecodeJson,
	NoVersionTag,
	UnreleasedMsg,
	UpdateChgLog string
}{
	BumpedVersion:    "cannot get bumped version: %v",
	CannotCalcHash:   "could not calculate a hash for the changes to be added to the %v",
	CannotDecodeJson: "could not decode JSON %v: %v",
	NoVersionTag:     "could not find semantic version tag",
	UnreleasedMsg:    "could not build unreleased message: %v",
	UpdateChgLog:     "could not upgrade changelog: %v",
}

var stdout = struct {
	AddTagManually,
	Cs,
	DiagnosticsFound,
	Fmt,
	NextSemVer,
	NoCommitsToBump,
	NoMoreDiagnostics,
	NothingToTag,
	PrependToChangelog,
	Trimmed,
	TrimmedStart,
	Wd string
}{
	AddTagManually:     "adding manual semver tag %v",
	Cs:                 "exec command string: %s",
	DiagnosticsFound:   "found diagnostic message: %v",
	Fmt:                "%v stdout: \n%v",
	NextSemVer:         "next semantic version: %v",
	NoCommitsToBump:    "There is nothing to bump",
	NoMoreDiagnostics:  "no more matches found",
	NothingToTag:       "next semantic version is empty, so we cannot release a tag",
	PrependToChangelog: "will prepend unreleased to the changelog.",
	Trimmed:            "trimmed:\n%v",
	TrimmedStart:       "stdout len: %v, start at index: %v",
	Wd:                 "%v working directory: %s",
}

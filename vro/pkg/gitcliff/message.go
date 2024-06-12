package gitcliff

var stderr = struct {
	BumpedVersion    string
	CannotDecodeJson string
	NoVersionTag     string
	UnreleasedMsg    string
	UpdateChgLog     string
}{
	BumpedVersion:    "cannot get bumped version: %v",
	CannotDecodeJson: "could not decode JSON %v: %v",
	NoVersionTag:     "could not find semantic version tag",
	UnreleasedMsg:    "could not build unreleased message: %v",
	UpdateChgLog:     "could not upgrade changelog: %v",
}

var stdout = struct {
	Cs string
	Wd string
}{
	Cs: "cs: %s",
	Wd: "wd: %s",
}

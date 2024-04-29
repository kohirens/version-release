package gitcliff

var stderr = struct {
	CannotDecodeJson string
	UpdateChgLog     string
}{
	CannotDecodeJson: "could not decode JSON %v: %v",
	UpdateChgLog:     "cannot update changelog: %v",
}

var stdout = struct {
	Cs string
	Wd string
}{
	Cs: "cs: %s",
	Wd: "wd: %s",
}

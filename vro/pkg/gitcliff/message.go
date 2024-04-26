package gitcliff

var stderr = struct {
	UpdateChgLog string
}{
	UpdateChgLog: "cannot update changelog: %v\n",
}

var stdout = struct {
	Cs string
	Wd string
}{
	Cs: "cs: %s\n",
	Wd: "wd: %s",
}

package gitchglog

var stderr = struct {
	CouldNotUpdateChgLog string
}{
	CouldNotUpdateChgLog: "cannot update changelog: %v\n",
}

var stdout = struct {
	Cs string
}{
	Cs: "cs: %s\n",
}

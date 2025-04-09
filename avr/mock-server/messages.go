package main

var stderr = struct {
	CannotLoad404Page,
	CouldNotDecodeJson,
	CouldNotReadBody,
	FileNotFound,
	InvalidLoginState,
	MockExist,
	NoMatch string
}{
	CannotLoad404Page:  "cannot load 404 page %q",
	CouldNotDecodeJson: "could not decode JSON: %s",
	CouldNotReadBody:   "could not read body: %s",
	FileNotFound:       "%q not found",
	MockExist:          "mock %v exists %v",
	NoMatch:            "re %V did not match %v",
}

var stdout = struct {
	LoadingFile string
}{
	LoadingFile: "loading %s file %s",
}

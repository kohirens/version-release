package main

var Stderr = struct {
	CannotLoad404Page,
	CouldNotDecodeJson,
	CouldNotReadBody,
	FileNotFound,
	InvalidLoginState,
	NoMatch string
}{
	CannotLoad404Page:  "cannot load 404 page %q",
	CouldNotDecodeJson: "could not decode JSON: %s",
	CouldNotReadBody:   "could not read body: %s",
	FileNotFound:       "%q not found",
	NoMatch:            "re %V did not match %v",
}

var Stdout = struct {
	LoadingFile string
}{
	LoadingFile: "loading %s file %s",
}

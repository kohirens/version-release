package main

var Stderr = struct {
	CannotLoad404Page  string
	CouldNotDecodeJson string
	CouldNotReadBody   string
	FileNotFound       string
	InvalidLoginState  string
}{
	CannotLoad404Page:  "cannot load 404 page %q",
	CouldNotDecodeJson: "could not decode JSON: %s",
	CouldNotReadBody:   "could not read body: %s",
	FileNotFound:       "%q not found",
}

var Stdout = struct {
	LoadingFile string
}{
	LoadingFile: "loading %s file %s",
}

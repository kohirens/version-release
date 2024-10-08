package lib

var stderr = struct {
	MissingEnv string
}{
	MissingEnv: "%s environment variable is not set",
}

var stdout = struct {
	RepoDetails string
}{
	RepoDetails: "found owner %v and repo name %v from %v",
}

package main

import (
	"bytes"
	"fmt"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/git"
	"github.com/kohirens/stdlib/log"
	"io"
	"net/http"
	"os"
)

const gitRoot = "cache/git"

// gitHttpBackend Takes a request and proxies it through git http-backend, then
// writes the output as a http.Response. This is in order to mock responses for
// git CLI when it sends request to remote origins like GitHub.com.
// This may seem sketchy, but this is a legit process that works well and uses
// http-backend as it was intended, but for testing purposes.
//
//	For more details about git http-backend, see https://git-scm.com/docs/git-http-backend
func gitHttpBackend(w http.ResponseWriter, r *http.Request) error {
	service := r.URL.Query().Get("service")
	ct := ""

	if service != "" {
		ct = fmt.Sprintf("application/x-%v-advertisement", service)
	}

	switch r.Header.Get("Content-Type") {
	case "application/x-git-receive-pack-request":
		ct = "application/x-git-receive-pack-result"
		service = "git-receive-pack"
	}

	log.Dbugf("service = %v", service)
	log.Dbugf("response Content-Type = %v", ct)

	bodyBytes, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		return fmt.Errorf("could not read body: %v\n", err1.Error())
	}

	// configure git http-backend request
	// See https://git-scm.com/docs/git-http-backend
	cgiVars := map[string]string{ // see https://git-scm.com/docs/git-http-backend#_environment
		"GIT_HTTP_EXPORT_ALL": "1",
		"GIT_PROTOCOL":        "2",
		"GIT_PROJECT_ROOT":    gitRoot, // This and PATH_INFO will be concatenated to form the repository path.
		"PATH_INFO":           r.URL.Path,
		"REMOTE_USER":         "git",
		"REMOTE_ADDR":         "github.com",
		"CONTENT_TYPE":        r.Header.Get("Content-Type"),
		"CONTENT_LENGTH":      r.Header.Get("Content-Length"),
		"QUERY_STRING":        r.URL.RawQuery,
		"REQUEST_METHOD":      r.Method,
	}

	log.Dbugf("cgiVars: %v", cgiVars)
	// we configure the http-backend request with environment variables that the CGI will read
	setEnv(cgiVars)
	// clear out any http-backend request environment variables.
	defer unsetEnv(cgiVars)

	so, se, _, _ := cli.RunCommandWithInput(
		".",
		"git",
		[]string{"http-backend"},
		bodyBytes,
	)

	if se != nil {
		return fmt.Errorf("problem with %v: %v\n", service, se.Error())
	}

	// Debug: print the response we got back from git http-backend
	log.Dbugf("so:\n%v\n", bytes.NewBuffer(so).String())

	// write the http-backend body to the response:
	bb := bytes.SplitAfter(so, []byte{13, 10, 13, 10})

	// just cheat for now and rebuild the headers; leaving out the expires.
	w.Header().Set("Content-Type", ct)
	w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
	w.Header().Set("Pragma", "no-cache")

	if _, e := w.Write(bb[1]); e != nil {
		return fmt.Errorf("%s\n", e.Error())
	}

	return nil
}

// unbundleRepo Clone a git bundle.
func unbundleRepo(bundle, owner string) {
	cloneToDir := cacheDir + "/git/" + owner
	if fsio.Exist(cloneToDir) {
		if e := os.RemoveAll(cloneToDir); e != nil {
			log.Logf("could not remove %s: %s", cloneToDir, e.Error())
		}
	}
	git.CloneFromBundle(bundle, cloneToDir, "testdata", "/")
}

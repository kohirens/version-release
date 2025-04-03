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
	"os/exec"
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
	default:

	}

	log.Dbugf("service = %v", service)
	log.Dbugf("response Content-Type = %v", ct)

	// configure git http-backend request
	// See https://git-scm.com/docs/git-http-backend
	httpBackendInput := map[string]string{ // see https://git-scm.com/docs/git-http-backend#_environment
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

	log.Dbugf("httpBackendInput: %v", httpBackendInput)
	// we configure the http-backend request with environment variables that the CGI will read
	//setEnv(httpBackendInput)
	//// clear out any http-backend request environment variables.
	//defer unsetEnv(httpBackendInput)

	//bodyBytes, e1 := io.ReadAll(r.Body)
	//if e1 != nil {
	//	return fmt.Errorf("could not read body: %v\n", e1.Error())
	//}

	//so, se, _, _ := cli.RunCommandWithInput(
	//	".",
	//	"git",
	//	[]string{"http-backend"},
	//	bodyBytes,
	//)
	//so, se, _, _ := cli.RunCommandWithInputAndEnv(
	//	".",
	//	"git",
	//	[]string{"http-backend"},
	//	bodyBytes,
	//	httpBackendInput,
	//)

	outputBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	e2 := RunPipedCommandWithInputAndEnv(
		".",
		"git",
		[]string{"http-backend"},
		httpBackendInput,
		r.Body,
		outputBuf,
		errBuf,
		os.Stdout,
	)

	if e2 != nil {
		return fmt.Errorf("problem proxying request through git http-backend: %v", e2.Error())
	}

	if errBuf != nil {
		return fmt.Errorf("problem with %v: %v\n", service, errBuf.String())
	}

	so := outputBuf.Bytes()
	// Debug: print the response we got back from git http-backend
	log.Dbugf("so:\n%v\n", string(so))

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

// RunPipedCommandWithInputAndEnv run an external program in a sub process, with
// input and setting environment variables in the sub process. It
// will pass in the os.Environ(), overwriting key=value pairs from env map,
// comparison for the key (variable name) is case-sensitive.
func RunPipedCommandWithInputAndEnv(
	wd,
	program string,
	args []string,
	env map[string]string,
	input io.Reader,
	output,
	errput,
	logger io.Writer,
) error {
	cmd := exec.Command(program, args...)
	cmd.Dir = wd
	ce := os.Environ()

	// overwrite or set environment variables
	if env != nil {
		ce = cli.AmendStringAry(ce, env)
	}

	cmd.Env = ce

	cmdIn, e1 := cmd.StdinPipe()
	if e1 != nil {
		return fmt.Errorf("could not get pipe to stdin: %v", e1.Error())
	}

	cmdOut, e2 := cmd.StdoutPipe()
	if e2 != nil {
		return e2
	}

	cmdErr, e3 := cmd.StderrPipe()
	if e3 != nil {
		return e3
	}

	if e := cmd.Start(); e != nil {
		return fmt.Errorf("could not run the command: %v", e.Error())
	}

	go func() {
		defer func() {
			_ = cmdIn.Close()
		}()

		// stream the input
		_, e := io.Copy(cmdIn, input)
		if e != nil {
			_, _ = fmt.Fprintf(logger, "problem piping input: %v", e.Error())
		}
	}()

	go func() {
		defer func() {
			_ = cmdOut.Close()
		}()

		_, e := io.Copy(output, cmdOut)
		if e != nil {
			_, _ = fmt.Fprintf(logger, "problem piping output: %v", e.Error())
		}
	}()

	go func() {
		defer func() {
			_ = cmdErr.Close()
		}()

		_, e := io.Copy(errput, cmdErr)
		if e != nil {
			_, _ = fmt.Fprintf(logger, "problem piping error: %v", e.Error())
		}
	}()

	return cmd.Wait()
}

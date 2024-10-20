package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/git"
	"github.com/kohirens/stdlib/log"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var (
	serverDir = Abs("mock-server")
)

// Abs Return the absolute path if it exists or the directory entered.
func Abs(loc string) string {
	full, e1 := filepath.Abs(loc)

	if e1 != nil {
		return loc
	}

	return full
}

func main() {
	handlers := &Handlers{
		"/": LoadMock,
	}

	// Register HTTP request handlers
	for endpoint, handler := range *handlers {
		http.HandleFunc(endpoint, handler)
	}

	// run the web server
	mainErr := http.ListenAndServeTLS(
		":443",
		"mock-server/ssl/certs/ca-cert-mock-server-CA.pem",
		"mock-server/ssl/private/mock-server-server.key",
		nil,
	)

	if mainErr != nil {
		log.Fatf("%v", mainErr.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

type Handlers map[string]func(w http.ResponseWriter, r *http.Request)

type Parameters struct {
	TriggeredFlow string `json:"triggered_flow"`
}

type PipelineParams struct {
	Branch     string     `json:"branch"`
	Parameters Parameters `json:"parameters"`
}

type tmplVars struct {
	Data map[string]string
}

type postData map[string]string

func LoadMock(w http.ResponseWriter, r *http.Request) {
	vars := &tmplVars{
		Data: map[string]string{
			"Path":  r.URL.Path,
			"Query": r.URL.RawQuery,
		},
	}
	mock := "does-not-exist.json"

	logToFile(serverDir+"/request-access.log", r.URL.String()+" : "+r.URL.RawQuery)

	switch r.URL.Path {
	case "/":
		mock = "health.json"

	case "/api/v2/project/gh/kohirens/version-release/pipeline":
		bodyBytes, err1 := io.ReadAll(r.Body)
		if err1 != nil {
			vars.Data["Error1"] = err1.Error()
			log.Errf(Stderr.CouldNotReadBody, err1.Error())
		}

		vars.Data["Post"] = string(bodyBytes)
		pp := &PipelineParams{}
		if e := json.Unmarshal(bodyBytes, pp); e != nil {
			vars.Data["Error2"] = e.Error()
			log.Errf(Stderr.CouldNotDecodeJson, e.Error())
		}
		mock = fmt.Sprintf("%s.json", pp.Parameters.TriggeredFlow)
		vars.Data["Mock"] = mock
		w.WriteHeader(201)
	case "/repos/kohirens/version-release/releases":
		b, _ := io.ReadAll(r.Body)
		var data postData

		_ = json.Unmarshal(b, &data)
		log.Logf("post data = %v", data)
		w.WriteHeader(201)
		vars.Data["TagNameDate"] = time.Now().Format("2006-01-02")

		mock = fmt.Sprintf("%v-releases.json", data["tag_name"])
		vars.Data["Mock"] = mock
	case "/kohirens/repo-01/info/refs":
		// Setup fixture
		fixedRepo := "/tmp/repo-01"
		if fsio.Exist(fixedRepo) {
			if e := os.RemoveAll(fixedRepo); e != nil {
				log.Logf("could not remove %s: %s", fixedRepo, e.Error())
			}
		}
		git.CloneFromBundle("repo-01", "/tmp", "testdata", "/")

		service, _, cgiVars, _ := gitHttpBackendProxy("request-repo-01/git-receive-pack", w, r)

		so, se, _, _ := cli.RunCommand(".", "git", []string{"http-backend"})
		if se != nil {
			log.Logf("problem with %s: %s", service, se.Error())
		}

		bb := bytes.SplitAfter(so, []byte{13, 10, 13, 10})
		if len(bb) > 1 {
			saveFile(serverDir+"/tmp/response-repo-01-info-refs.txt", bb[1])
		}

		if _, e := w.Write(bb[1]); e != nil {
			log.Logf("%s", e.Error())
		}

		unsetEnv(cgiVars)

		return
	case "/kohirens/repo-01/git-receive-pack":
		service, body, cgiVars, _ := gitHttpBackendProxy("request-repo-01/git-receive-pack", w, r)

		so, se, _, _ := cli.RunCommandWithInput(
			".",
			"git",
			[]string{"http-backend", "/tmp/repo-01"},
			body,
		)
		if se != nil {
			log.Logf("problem with %s: %s\n", service, se.Error())
		}

		log.Logf("so:\n%v\n", bytes.NewBuffer(so).String())

		bb := bytes.SplitAfter(so, []byte{13, 10, 13, 10})
		if len(bb) > 1 {
			saveFile(serverDir+"/tmp/response-repo-01-git-receive-pack.txt", bb[1])
		}

		if _, e := w.Write(bb[1]); e != nil {
			log.Logf("%s\n", e.Error())
		}

		unsetEnv(cgiVars)

		return
	case "/repos/kohirens/version-release/pulls":
		mock = "make-pr.json"
		vars.Data["Mock"] = mock
		w.WriteHeader(http.StatusCreated)
	case "/repos/kohirens/version-release/pulls/1/merge":
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusNoContent)
			return
		case "PUT":
			mock = "merge-pr.json"
			w.WriteHeader(http.StatusOK)
		}
		vars.Data["Mock"] = mock
	case "/repos/kohirens/repo-01/pulls":
		mock = "make-pr.json"
		vars.Data["Mock"] = mock
		w.WriteHeader(http.StatusCreated)
	case "/repos/kohirens/repo-01/pulls/1/merge":
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusNoContent)
			return
		case "PUT":
			mock = "merge-pr.json"
			w.WriteHeader(http.StatusOK)
		}
		vars.Data["Mock"] = mock
	}

	log.Infof(Stdout.LoadingFile, "JSON", mock)

	tFile := serverDir + "/responses/" + mock
	if e := LoadTemplate(tFile, w, vars); e != nil {
		log.Errf("%v", e.Error())
		logToFile(serverDir+"/request-404.log", r.URL.String())
		Load404Page(serverDir+"/responses/not-found.json", w, vars)
	}
}

func LoadTemplate(tFile string, w io.Writer, vars *tmplVars) error {
	if !fsio.Exist(tFile) {
		return fmt.Errorf(Stderr.FileNotFound, tFile)
	}

	t, err1 := template.ParseFiles(tFile)
	if err1 != nil {
		return err1
	}

	return t.Execute(w, vars)
}

func Load404Page(tFile string, w http.ResponseWriter, vars *tmplVars) {
	w.WriteHeader(404)
	err1 := LoadTemplate(tFile, w, vars)
	if err1 != nil {
		panic(fmt.Errorf(Stderr.CannotLoad404Page))
	}
}

func logToFile(filename string, line string) {
	f, e1 := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0744)
	if e1 != nil {
		log.Logf("%s\n", e1.Error())
	}

	l := fmt.Sprintf("%s\n", line)
	_, e2 := f.Write([]byte(l))
	if e2 != nil {
		log.Logf("%s\n", e2.Error())
	}
}

func gitHttpBackendProxy(prefix string, w http.ResponseWriter, r *http.Request) (string, []byte, map[string]string, string) {
	name := fmt.Sprintf("%s/tmp/%s", serverDir, strings.Replace(prefix, "/", "-", -1))
	saveHeadersToFile(name+"-headers.txt", r.Header)

	log.Logf("request method = %s\n", r.Method)
	log.Logf("request query = %s\n", r.URL.RawQuery)
	log.Logf("request Content-Type = %s\n", r.Header.Get("Content-Type"))

	service := r.URL.Query().Get("service")
	ct := ""

	if service != "" {
		ct = fmt.Sprintf("application/x-%s-advertisement", service)
	}

	switch r.Header.Get("Content-Type") {
	case "application/x-git-receive-pack-request":
		ct = "application/x-git-receive-pack-result"
		service = "git-receive-pack"
	}

	log.Logf("service = %s\n", service)
	log.Logf("response Content-Type = %s\n", ct)

	bodyBits, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		log.Logf("%s\n", err1.Error())
	}

	saveFile(fmt.Sprintf("%s-%s-body.txt", name, service), bodyBits)

	cgiVars := map[string]string{ // see https://git-scm.com/docs/git-http-backend/en#_environment
		"GIT_HTTP_EXPORT_ALL": "1",
		"GIT_PROTOCOL":        "2",
		"GIT_PROJECT_ROOT":    "/tmp",
		"PATH_INFO":           strings.Replace(r.URL.Path, "/kohirens", "", 1),
		"REMOTE_USER":         "git",
		"REMOTE_ADDR":         "github.com",
		"CONTENT_TYPE":        r.Header.Get("Content-Type"),
		"CONTENT_LENGTH":      r.Header.Get("Content-Length"),
		"QUERY_STRING":        r.URL.RawQuery,
		"REQUEST_METHOD":      r.Method,
	}

	w.Header().Set("Content-Type", ct)
	w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
	w.Header().Set("Pragma", "no-cache")

	setEnv(cgiVars)

	return service, bodyBits, cgiVars, ct
}

// saveFile Save data to a file overwriting if it exists.
func saveFile(filename string, b []byte) {
	f, e1 := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0744)
	if e1 != nil {
		log.Logf("%s\n", e1.Error())
	}

	_, e3 := f.Write(b)
	if e3 != nil {
		log.Logf("%s\n", e3.Error())
	}
}

func saveHeadersToFile(filename string, header http.Header) {
	f, e1 := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0744)
	if e1 != nil {
		log.Logf("%s\n", e1.Error())
	}

	h := fmt.Sprintf("--- %s\n", time.Now().Format("2006-01-02"))
	_, e2 := f.Write([]byte(h))
	if e2 != nil {
		log.Logf("%s\n", e2.Error())
	}

	for k, v := range header {
		l := fmt.Sprintf("%s: %s\n", k, v)
		_, e3 := f.Write([]byte(l))
		if e3 != nil {
			log.Logf("%s\n", e3.Error())
		}
	}
}

func setEnv(vars map[string]string) {
	for k, v := range vars {
		if e := os.Setenv(k, v); e != nil {
			log.Logf("could not set environment variable %q: %s\n", k, e.Error())
		}
	}
}

func unsetEnv(vars map[string]string) {
	for k, _ := range vars {
		if e := os.Unsetenv(k); e != nil {
			log.Logf("could not set environment variable %q: %s\n", k, e.Error())
		}
	}
}

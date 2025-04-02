package main

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/git"
	"github.com/kohirens/stdlib/log"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	serverDir = abs("mock-server")
	cacheDir  = abs("cache")
)

func main() {
	// Set the logging level via an environment variable.
	vl, vlFound := os.LookupEnv("VERBOSITY_LEVEL")
	if vlFound {
		num, _ := strconv.ParseInt(vl, 10, 64)
		log.VerbosityLevel = int(num)
	}

	// Register HTTP request handlers
	handler := http.NewServeMux()
	handler.HandleFunc("/", LoadMock)
	handler.HandleFunc("/{owner}/{repo}/info/refs", gitHttpBackendProxy)
	handler.HandleFunc("/{owner}/{repo}/git-receive-pack", gitHttpBackendProxy)

	// run the web server
	mainErr := http.ListenAndServeTLS(
		":443",
		"mock-server/ssl/certs/ca-cert-mock-server-CA.pem",
		"mock-server/ssl/private/mock-server-server.key",
		handler,
	)

	if mainErr != nil {
		log.Fatf("%v", mainErr.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func LoadMock(w http.ResponseWriter, r *http.Request) {
	vars := &tmplVars{
		Data: map[string]string{
			"Path":  r.URL.Path,
			"Query": r.URL.RawQuery,
		},
	}
	mock := "does-not-exist.json"
	log.Infof("r.URL.Path = %v", r.URL.Path)
	// used this to capture the request for developing mocks for new request that have not been handled.
	captureRequestInfo(r.URL.Path, r)

	logToFile(cacheDir+"/request/access.log", r.URL.String()+" : "+r.URL.RawQuery)

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
		fixedRepo := gitRoot + "/kohirens/repo-01"
		if fsio.Exist(fixedRepo) {
			if e := os.RemoveAll(fixedRepo); e != nil {
				log.Logf("could not remove %s: %s", fixedRepo, e.Error())
			}
		}
		git.CloneFromBundle("repo-01", gitRoot, "testdata", "/")

		if e := gitHttpBackend(w, r); e != nil {
			log.Logf("%s", e.Error())
		}

		return
	case "/kohirens/repo-01/git-receive-pack":
		if e := gitHttpBackend(w, r); e != nil {
			log.Logf("%s", e.Error())
		}

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

		if e := getBody(r, "prm"); e != nil {
			panic(e.Error())
			return
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
	if e := loadTemplate(tFile, w, vars); e != nil {
		log.Errf(e.Error())
		logToFile(cacheDir+"/request/404.log", r.URL.String())
		load404Page(cacheDir+"/responses/not-found.json", w, vars)
	}
}

// gitHttpBackendProxy the http.HandleFunc that proxies request through git http-backend.
// For more details about git http-backend, see https://git-scm.com/docs/git-http-backend
func gitHttpBackendProxy(w http.ResponseWriter, r *http.Request) {
	log.Infof("\ngit http-backend proxy begin")
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")

	// used this to capture the request for developing mocks for new request that have not been handled.
	captureRequestInfo(r.URL.Path, r)

	log.Infof("r.URL.Path = %v", r.URL.Path)

	unbundleRepo(repo, owner)

	if e := gitHttpBackend(w, r); e != nil {
		log.Logf(e.Error())
	}

	log.Infof("git http-backend proxy end\n")
}

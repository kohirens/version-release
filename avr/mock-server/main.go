package main

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/log"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	serverDir  = abs("mock-server")
	cacheDir   = abs("cache")
	requestDir = cacheDir + "/request"
)

func main() {
	// Set the logging level via an environment variable.
	vl, vlFound := os.LookupEnv("VERBOSITY_LEVEL")
	if vlFound {
		num, _ := strconv.ParseInt(vl, 10, 64)
		log.VerbosityLevel = int(num)
	}

	// Register HTTP request handlers
	// Also read up at https://go.dev/blog/routing-enhancements
	handler := http.NewServeMux()
	handler.HandleFunc("/", notFound)
	handler.HandleFunc("api.circleci.com/", mockCircleCi)
	handler.HandleFunc("circleci.com/", mockCircleCi)
	handler.HandleFunc("/health", healthCheck)
	handler.HandleFunc("/{owner}/{repo}/info/refs", gitHttpBackendProxy)
	handler.HandleFunc("/{owner}/{repo}/git-receive-pack", gitHttpBackendProxy)
	handler.HandleFunc("/{owner}/{repo}/HEAD", gitHttpBackendProxy)
	handler.HandleFunc("api.github.com/repos/{owner}/{repo}/", mockGitHubApi)

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

func mockCircleCi(w http.ResponseWriter, r *http.Request) {
	vars := &tmplVars{
		Data: map[string]string{
			"Path":  r.URL.Path,
			"Query": r.URL.RawQuery,
		},
	}
	mock := ""

	log.Infof("mockCircleCi controller")
	log.Infof("r.URL.Path = %v", r.URL.Path)

	logToFile(cacheDir+"/request/access.log", r.URL.String()+" : "+r.URL.RawQuery)

	switch r.URL.Path {
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
	}

	log.Infof(Stdout.LoadingFile, "JSON", mock)

	if mock == "" {
		load404Page(serverDir+"/responses/not-found.json", w, vars)
		return
	}

	if e := loadTemplate(serverDir+"/responses/"+mock, w, vars); e != nil {
		log.Errf(e.Error())
		logToFile(cacheDir+"/request/404.log", r.URL.String())
	}
}

// gitHttpBackendProxy the http.HandleFunc that proxies request through git http-backend.
// For more details about git http-backend, see https://git-scm.com/docs/git-http-backend
func gitHttpBackendProxy(w http.ResponseWriter, r *http.Request) {
	log.Infof("\ngit http-backend proxy begin")
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")

	log.Infof("r.URL.Path = %v", r.URL.Path)

	unbundleRepo(repo, owner)

	if e := gitHttpBackend(w, r); e != nil {
		log.Logf(e.Error())
	}

	log.Infof("git http-backend proxy end\n")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	vars := &tmplVars{
		Data: map[string]string{
			"Path":  r.URL.Path,
			"Query": r.URL.RawQuery,
		},
	}

	if e := loadTemplate(serverDir+"/responses/health.json", w, vars); e != nil {
		log.Errf(e.Error())
		logToFile(cacheDir+"/request/404.log", r.URL.String())
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	vars := &tmplVars{
		Data: map[string]string{
			"Path":  r.URL.Path,
			"Query": r.URL.RawQuery,
		},
	}

	log.Dbugf("notFound controller")
	log.Infof("r.URL.Path = %v", r.URL.Path)

	// used this to capture the request for developing mocks for new request that have not been handled.
	captureRequestInfo(r.URL.Path, r)

	load404Page(serverDir+"/responses/not-found.json", w, vars)
}

func mockGitHubApi(w http.ResponseWriter, r *http.Request) {
	vars := &tmplVars{
		Data: map[string]string{
			"Path":  r.URL.Path,
			"Query": r.URL.RawQuery,
		},
	}
	mock := ""

	log.Infof("\nmockGitHubApi begin")
	log.Infof("r.URL.Path = %v", r.URL.Path)

	logToFile(cacheDir+"/request/access.log", r.URL.String()+" : "+r.URL.RawQuery)

	owner := r.PathValue("owner")
	repo := r.PathValue("repo")
	p := strings.Replace(r.URL.Path, "/repos/"+owner+"/"+repo, "", 1)

	log.Dbugf("p = %v", p)

	switch p {
	case "/releases":
		b, _ := io.ReadAll(r.Body)
		var data postData

		_ = json.Unmarshal(b, &data)
		log.Logf("post data = %v", data)
		w.WriteHeader(201)
		vars.Data["TagNameDate"] = time.Now().Format("2006-01-02")

		mock = fmt.Sprintf("%v-releases.json", data["tag_name"])
		vars.Data["Mock"] = mock

	case "/pulls":
		mock = "make-pr.json"
		vars.Data["Mock"] = mock
		w.WriteHeader(http.StatusCreated)
	case "/pulls/1/merge":
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusNoContent)
			return
		case "PUT":
			mock = "merge-pr.json"
			w.WriteHeader(http.StatusOK)
		}
		vars.Data["Mock"] = mock
		if strings.Contains(p, "version-release") {
			if e := getBody(r, "prm"); e != nil {
				panic(e.Error())
				return
			}
		}
	}

	log.Infof(Stdout.LoadingFile, "JSON", mock)

	if mock == "" {
		load404Page(serverDir+"/responses/not-found.json", w, vars)
		return
	}

	if e := loadTemplate(serverDir+"/responses/"+mock, w, vars); e != nil {
		log.Errf(e.Error())
		logToFile(cacheDir+"/request/404.log", r.URL.String())
	}
}

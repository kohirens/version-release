package main

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/version-release/avr/pkg/github"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	wd         = rootDir()
	serverDir  = "mock-server"
	cacheDir   = "cache"
	requestDir = cacheDir + "/request"
	log        = logger.Standard{}
	mockDir    = wd + "mock-server/responses/"
)

func main() {
	// Set the logging level via an environment variable.
	vl, vlFound := os.LookupEnv("VERBOSITY_LEVEL")
	if vlFound {
		num, _ := strconv.ParseInt(vl, 10, 64)
		logger.VerbosityLevel = int(num)
	}

	log.Dbugf("working dir: %s", wd)

	// Register HTTP request handlers
	// Also read up at https://go.dev/blog/routing-enhancements
	handler := http.NewServeMux()
	handler.HandleFunc("/", notFound)
	handler.HandleFunc("api.circleci.com/", mockCircleCi)
	handler.HandleFunc("circleci.com/", mockCircleCi)
	handler.HandleFunc("/health", healthCheck)
	handler.HandleFunc("github.com/{owner}/{repo}/", gitHttpBackendProxy)
	handler.HandleFunc("api.github.com/repos/{owner}/{repo}/{remain...}", mockGitHubApi)

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
			log.Errf(stderr.CouldNotReadBody, err1.Error())
		}

		vars.Data["Post"] = string(bodyBytes)
		pp := &PipelineParams{}
		if e := json.Unmarshal(bodyBytes, pp); e != nil {
			vars.Data["Error2"] = e.Error()
			log.Errf(stderr.CouldNotDecodeJson, e.Error())
		}
		mock = fmt.Sprintf("%s.json", pp.Parameters.TriggeredFlow)
		vars.Data["Mock"] = mock
		w.WriteHeader(201)
	}

	log.Infof(stdout.LoadingFile, "JSON", mock)

	if mock == "" {
		notFound(w, r)
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
	logToFile(cacheDir+"/request/access.log", r.URL.String()+" : "+r.URL.RawQuery)

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

	w.WriteHeader(404)

	e1 := loadTemplate(serverDir+"/responses/not-found.json", w, vars)
	if e1 != nil {
		panic(fmt.Errorf(stderr.CannotLoad404Page))
	}
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

	repo := r.PathValue("repo")
	remain := r.PathValue("remain")

	log.Dbugf("remain = %v", remain)

	switch remain {
	case "releases":
		b, _ := io.ReadAll(r.Body)
		var data postData

		_ = json.Unmarshal(b, &data)
		log.Logf("post data = %v", data)
		w.WriteHeader(201)
		vars.Data["TagNameDate"] = time.Now().Format("2006-01-02")

		mock = fmt.Sprintf("%v-releases.json", data["tag_name"])
		vars.Data["Mock"] = mock

	case "pulls":
		mock = "make-pr.json"
		vars.Data["Mock"] = mock
		w.WriteHeader(http.StatusCreated)
	case "pulls/1/merge":
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusNoContent)
			return
		case "PUT":
			mock = "merge-pr.json"
			w.WriteHeader(http.StatusOK)
		}
		vars.Data["Mock"] = mock
		if strings.Contains(remain, "version-release") {
			if e := getBody(r, "prm"); e != nil {
				panic(e.Error())
				return
			}
		}
	default:
		if repo == "readme-only-no-tags" && remain == "git/trees" {
			log.Dbugf("speicial case, we are checking for additional files.")
			// TODO: Make sure the new tree contains the additional files
			b, _ := io.ReadAll(r.Body)
			var data github.RequestTree

			_ = json.Unmarshal(b, &data)
			log.Dbugf("body posted: %v", string(b))
			for _, tree := range data.Tree {
				log.Dbugf("tree paths posted: %v", tree.Path)
				if tree.Path == "generated-file.txt" {
					// The addition file was found in the list of files to add to the commit.
					log.Dbugf("additonal file %v found", "generated-file.txt")
					goto go2go
				}
			}
			return
		}
	go2go:

		e := getResponseMock(repo, remain, w)
		if e == nil {
			return
		}
		log.Errf(e.Error())
	}

	log.Infof(stdout.LoadingFile, "JSON", mock)

	if mock == "" {
		notFound(w, r)
		return
	}

	if e := loadTemplate(serverDir+"/responses/"+mock, w, vars); e != nil {
		log.Errf(e.Error())
		logToFile(cacheDir+"/request/404.log", r.URL.String())
	}
}

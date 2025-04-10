package main

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

// Abs Return the absolute path if it exists or the directory entered.
func rootDir() string {
	w, e1 := os.Getwd()
	if e1 != nil {
		return ""
	}

	ps := string(os.PathSeparator)
	sep := ps + "avr"
	idx := strings.Index(w, sep) + 4
	_wd := w[:idx]

	return _wd + "/"
}

// captureRequestInfo For debugging purposes.
func captureRequestInfo(prefix string, r *http.Request) {
	if prefix == "" || prefix == "/" {
		return
	}

	if prefix[:1] == "/" {
		prefix = prefix[1:]
	}

	// Build a filename based on the url path.
	name := requestDir + "/" + strings.TrimRight(strings.Replace(prefix, "/", "-", -1), "-") + "-http.txt"
	log.Dbugf(stdout.CaptureRequest, name)

	headers := getHeadersAsString(r.Header)

	bodyBits, e1 := io.ReadAll(r.Body)
	if e1 != nil {
		log.Errf(stderr.CouldNotReadBody, e1.Error())
		return
	}

	// build an HTTP request string
	req := fmt.Sprintf("%v %v %v\r\nHost: %v\r\n%v\r\n%s\r\n", r.Method, r.URL.String(), r.Proto, r.Host, headers, bodyBits)

	saveFile(name, []byte(req))
}

func getBody(r *http.Request, prefix string) error {
	name := fmt.Sprintf("%v/tmp/%v", serverDir, prefix)

	bodyBits, e1 := io.ReadAll(r.Body)
	if e1 != nil {
		return fmt.Errorf(stderr.CouldNotReadBody, r.URL, e1.Error())
	}

	fname := fmt.Sprintf("%v-body.txt", name)
	saveFile(fname, bodyBits)

	return nil
}

// getResponseMock Parse a URL pack into a file path, returning the path if it exists.
func getResponseMock(repo, p string, w http.ResponseWriter) error {
	// Convert the path to a filename.
	filename := mockDir + repo + "/" + strings.Replace(p, "/", "-", -1) + ".json"
	exists := fsio.Exist(filename)

	log.Dbugf(stderr.MockExist, filename, exists)

	if !exists {

		return fmt.Errorf(stderr.FindMock, filename)
	}

	content, e1 := os.ReadFile(filename)
	if e1 != nil {
		return fmt.Errorf(stderr.MockLoad, filename, e1.Error())
	}

	log.Dbugf(stdout.MockLoad, filename)

	obj := struct {
		Status string `json:"status,omitempty"`
	}{}

	if e := json.Unmarshal(content, &obj); e != nil {
		return e
	}

	if w != nil && obj.Status != "" {
		log.Dbugf(stdout.MockStatus, obj.Status)
		sc, e := strconv.ParseInt(obj.Status, 10, 0)
		if e != nil {
			return fmt.Errorf(stderr.MockStatus, e.Error())
		}
		w.WriteHeader(int(sc))
	}

	if w != nil {
		_, e2 := w.Write(content)
		if e2 != nil {
			return fmt.Errorf(stderr.MockWrite, e2.Error())
		}
	}

	return nil
}

func loadTemplate(tFile string, w io.Writer, vars *tmplVars) error {
	if !fsio.Exist(tFile) {
		log.Dbugf(stderr.TemplateFind, tFile)
		return fmt.Errorf(stderr.FileNotFound, tFile)
	}

	t, e1 := template.ParseFiles(tFile)
	if e1 != nil {
		log.Dbugf(stderr.TemplateLoad, tFile)
		return e1
	}

	log.Dbugf(stdout.TemplateLoad, tFile)
	return t.Execute(w, vars)
}

func logToFile(filename string, line string) {
	dirname := filepath.Dir(filename)
	if os.MkdirAll(dirname, 0755) != nil {
		log.Errf(stderr.MakeDir, dirname)
	}

	f, e1 := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0774)
	if e1 != nil {
		log.Logf(e1.Error())
	}

	l := fmt.Sprintf("%s\n", line)
	_, e2 := f.Write([]byte(l))
	if e2 != nil {
		log.Logf("%s\n", e2.Error())
	}
}

// saveFile Save data to a file overwriting if it exists.
func saveFile(filename string, b []byte) {
	f, e1 := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0744)
	if e1 != nil {
		log.Errf(stderr.FileOpen, e1.Error())
		return
	}

	_, e2 := f.Write(b)
	if e2 != nil {
		log.Errf(stderr.FileWrite, e2.Error())
	}
}

// getHeadersAsString Format headers as a string as they would appear in an HTTP request
func getHeadersAsString(header http.Header) string {
	h := ""
	for k, v := range header {
		for _, vv := range v {
			h += fmt.Sprintf("%v: %v\r\n", k, vv)
		}
	}

	return h
}

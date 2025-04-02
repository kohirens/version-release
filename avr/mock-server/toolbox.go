package main

import (
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/log"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Abs Return the absolute path if it exists or the directory entered.
func abs(loc string) string {
	full, e1 := filepath.Abs(loc)

	if e1 != nil {
		return loc
	}

	return full
}

// captureRequestInfo For debugging purposes.
func captureRequestInfo(prefix string, r *http.Request) {
	if prefix == "" || prefix == "/" {
		return
	}

	if prefix[:1] == "/" {
		prefix = prefix[1:]
	}

	//build a filename based on the url path
	name := cacheDir + "/" + strings.Replace(prefix, "/", "-", -1)
	log.Dbugf("capture request info into %v", name)

	headers := getHeadersAsString(r.Header)

	bodyBits, e1 := io.ReadAll(r.Body)
	if e1 != nil {
		log.Errf("could not read request body: %v", e1.Error())
		return
	}

	// build an HTTP request string
	req := fmt.Sprintf("%v %v\r\n%v\r\n%s\r\n", r.Method, r.URL.String(), headers, bodyBits)

	saveFile(name, []byte(req))
}

func getBody(r *http.Request, prefix string) error {
	name := fmt.Sprintf("%v/tmp/%v", serverDir, prefix)

	bodyBits, e1 := io.ReadAll(r.Body)
	if e1 != nil {
		return fmt.Errorf("cannot read body of %v: %v\n", r.URL, e1.Error())
	}

	fname := fmt.Sprintf("%v-body.txt", name)
	saveFile(fname, bodyBits)

	return nil
}

func loadTemplate(tFile string, w io.Writer, vars *tmplVars) error {
	if !fsio.Exist(tFile) {
		return fmt.Errorf(Stderr.FileNotFound, tFile)
	}

	t, e1 := template.ParseFiles(tFile)
	if e1 != nil {
		return e1
	}

	return t.Execute(w, vars)
}

func load404Page(tFile string, w http.ResponseWriter, vars *tmplVars) {
	w.WriteHeader(404)
	e1 := loadTemplate(tFile, w, vars)
	if e1 != nil {
		panic(fmt.Errorf(Stderr.CannotLoad404Page))
	}
}

func logToFile(filename string, line string) {
	dirname := filepath.Dir(filename)
	if os.MkdirAll(dirname, 0755) != nil {
		log.Errf("could not create dir: %v", dirname)
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
		log.Errf("could not open file %v", e1.Error())
		return
	}

	_, e2 := f.Write(b)
	if e2 != nil {
		log.Errf("could not write file %v", e2.Error())
	}
}

func getHeadersAsString(header http.Header) string {
	h := ""
	for k, v := range header {
		h += fmt.Sprintf("%v: %v\r\n", k, v)
	}

	return h
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

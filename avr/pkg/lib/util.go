package lib

import (
	"fmt"
	"github.com/kohirens/stdlib/logger"
	"os"
	"regexp"
	"strings"
)

const CheckSemVer = `v?\d+\.\d+\.\d+(-.+)?`

var log = logger.Standard{}

// To get environment variables that are already available see
// https://circleci.com/docs/variables/#built-in-environment-variables
// and also https://circleci.com/docs/env-vars/#order-of-precedence
func GetRequiredEnvVars(eVarNames []string) (map[string]string, error) {
	errStr := ""
	eVars := map[string]string{}

	for _, name := range eVarNames {
		val, ok := os.LookupEnv(name)
		if !ok { // collect all the errors, then report
			errStr = errStr + fmt.Sprintf(stderr.MissingEnv, name) + "\n"
			continue
		}

		eVars[name] = val
	}

	if len(errStr) > 0 {
		return nil, fmt.Errorf(errStr)
	}

	return eVars, nil
}

func GetEnv(name string) string {
	value, exists := os.LookupEnv(name)
	if !exists {
		panic(fmt.Sprintf("Environment variable %s not found", name))
	}
	return value
}

func ParseRepositoryUri(uri string) (string, string, string) {
	//https://github.com/kohirens/version-release
	//git@github.com:kohirens/version-release.git
	re := regexp.MustCompile(`^(https://|git@)([^/:]+)[/:]([^/]+)/(.+)`)
	m := re.FindAllStringSubmatch(uri, -1)

	if m != nil {
		log.Dbugf(stdout.RepoDetails, m[0][2], m[0][3], uri)

		return m[0][2], m[0][3], strings.Replace(m[0][4], ".git", "", 1)
	}

	return "", "", ""
}

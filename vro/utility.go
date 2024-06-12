package main

import (
	"fmt"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release/vro/pkg/gitcliff"
	"os"
)

// To get environment variables that are already available see
// https://circleci.com/docs/variables/#built-in-environment-variables
// and also https://circleci.com/docs/env-vars/#order-of-precedence
func getRequiredEnvVars(eVarNames []string) (envVars, error) {
	errStr := ""
	eVars := envVars{}

	for _, name := range eVarNames {
		val, ok := os.LookupEnv(name)
		if !ok { // collect all the errors, then report
			errStr = errStr + fmt.Sprintf(stderr.MissingEnv, name)
			continue
		}

		eVars[name] = val
	}

	if len(errStr) > 0 {
		return nil, fmt.Errorf(errStr)
	}

	return eVars, nil
}

func nextVersion(semVer string, wd string) (string, error) {
	// check if a version has been provided as input.
	nextVer := semVer
	if nextVer == "" {
		nextVer = gitcliff.Bump(wd)
	}

	log.Infof("semVer = %v", nextVer)

	if nextVer == "" {
		return "", fmt.Errorf(stderr.NothingToTag)
	}

	return nextVer, nil
}

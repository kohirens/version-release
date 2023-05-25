package main

import (
	"fmt"
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

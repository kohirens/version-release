package lib

import (
	"fmt"
	"os"
)

const CheckSemVer = `v?\d+\.\d+\.\d+(-.+)?`

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

func GetEnv(m *map[string]string, name string) {
	val, k := os.LookupEnv(name)
	if k {
		(*m)[name] = val
	}
}

// GetVal Get a value from the environment or use the default. Panic if both are empty.
func GetVal(name, def string) string {
	val, k := os.LookupEnv(name)
	if k {
		return val
	}

	if def != "" {
		return def
	}

	panic(fmt.Sprintf("could not retrieve value for %v", name))
}

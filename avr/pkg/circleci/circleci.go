// Package circleci is a small subset wrapper for the CircleCI API.
// See https://circleci.com/docs/api/v2/index.html
package circleci

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/version-release/avr/pkg/lib"
	"os"
	"time"
)

type ParameterMap struct {
	TriggeredFlow string `json:"triggered_flow"`
}

type PipelineParams struct {
	Branch     string        `json:"branch"`
	Parameters *ParameterMap `json:"parameters"`
}

type TriggerPipelineResponse struct {
	Id        string    `json:"id"`
	State     string    `json:"state"`
	Number    int       `json:"number"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

const (
	EnvApiHost         = "PARAM_CIRCLECI_API_HOST"
	EnvAppHost         = "PARAM_CIRCLECI_APP_HOST"
	EnvCircleCI        = "CIRCLECI"
	EnvProjectRepoName = "CIRCLE_PROJECT_REPONAME"
	EnvProjectUsername = "CIRCLE_PROJECT_USERNAME"
	EnvRepoUrl         = "CIRCLE_REPOSITORY_URL"
	EnvToken           = "CIRCLE_TOKEN"
	EnvUsername        = "CIRCLE_USERNAME"
	Name               = "circleci"
)

var (
	envs = map[string]string{}
	log  = logger.Standard{}
)

func init() {
	processEnvironment()
}

// processEnvironment Load the environment variables and panic when any are not set.
func processEnvironment() {
	var e error
	if os.Getenv(EnvCircleCI) != "true" {
		// were not in the CircleCI environment, so exit.
		return
	}

	envs, e = lib.GetRequiredEnvVars([]string{
		EnvProjectRepoName,
		EnvProjectUsername,
		EnvToken,
	})

	if e != nil {
		panic(e.Error())
	}
}

func PipelineParameters(branch, flow string) ([]byte, error) {
	params := &PipelineParams{
		Branch:     branch,
		Parameters: &ParameterMap{TriggeredFlow: flow},
	}

	b, err1 := json.Marshal(params)
	if err1 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonEncode, params, err1.Error())
	}

	log.Infof(stdout.PipelineParams, b)

	return b, nil
}

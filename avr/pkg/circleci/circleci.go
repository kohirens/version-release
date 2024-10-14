// Package circleci is a small subset wrapper for the CircleCI API.
// See https://circleci.com/docs/api/v2/index.html
package circleci

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release/avr/pkg/github"
	"github.com/kohirens/version-release/avr/pkg/lib"
	"time"
)

// TODO: remove this reference to the github package.
type GithubClient interface {
	TagAndRelease(revision, tag string) (*github.ReleasesResponse, error)
	PublishChangelog(wd, branch, header, msg string, files []string) error
}

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
	EnvProjectReponame = "CIRCLE_PROJECT_REPONAME"
	EnvProjectUsername = "CIRCLE_PROJECT_USERNAME"
	EnvRepoUrl         = "CIRCLE_REPOSITORY_URL"
	EnvToken           = "CIRCLE_TOKEN"
	EnvUsername        = "CIRCLE_USERNAME"
	EnvVcsType         = "PARAM_VCS_TYPE"
	Name               = "circleci"
)

var (
	env map[string]string
)

func init() {
	processEnvironment()
}

// processEnvironment Load variable from the environment.
func processEnvironment() {
	lib.GetEnv(&env, EnvApiHost)
	lib.GetEnv(&env, EnvAppHost)
	lib.GetEnv(&env, EnvToken)
	lib.GetEnv(&env, EnvProjectUsername)
	lib.GetEnv(&env, EnvProjectReponame)
	lib.GetEnv(&env, EnvUsername)
	lib.GetEnv(&env, EnvVcsType)
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

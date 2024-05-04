// Package circleci is a small subset wrapper for the CircleCI API.
// See https://circleci.com/docs/api/v2/index.html
package circleci

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release/vro/pkg/github"
	"io"
	"net/http"
	"time"
)

type GithubClient interface {
	TagAndRelease(revision, tag string) (*github.ReleasesResponse, error)
	PublishChangelog(wd, branch, chaneLogFile string) error
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

func GetPipelineParameters(branch, flow string) ([]byte, error) {
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

// TriggerPipeline Trigger a pipeline based on parameters set. See docs
// https://circleci.com/docs/api/v2/index.html#tag/Pipeline; also see an
// example:
// https://circleci.com/docs/triggers-overview/#run-a-pipeline-using-the-api
// Note: keep in mind that you have to use a personal API token; project tokens
// are currently not supported on CircleCI API (v2) at this time.
// see: https://circleci.com/docs/api/v2/#operation/triggerPipeline
func TriggerPipeline(pp []byte, client *http.Client, eVars map[string]string) error {
	VcsType := eVars["PARAM_VCS_TYPE"]
	CircleProjectUsername := eVars["CIRCLE_PROJECT_USERNAME"]
	CircleProjectReponame := eVars["CIRCLE_PROJECT_REPONAME"]

	projectUrl := fmt.Sprintf(
		"%s/api/v2/project/%s/%s/%s/pipeline",
		eVars["PARAM_CIRCLECI_API_HOST"],
		VcsType,
		CircleProjectUsername,
		CircleProjectReponame,
	)

	log.Logf(stdout.CircleProjectUrl, projectUrl)

	req, err1 := http.NewRequest("POST", projectUrl, bytes.NewReader(pp))
	if err1 != nil {
		return fmt.Errorf(stderr.CouldNotPostRequest, err1.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Circle-Token", eVars["CIRCLE_TOKEN"])
	req.Header.Set(
		"Authorization",
		"Basic "+base64.StdEncoding.EncodeToString([]byte(eVars["CIRCLE_TOKEN"]+":")))

	res, err2 := client.Do(req)
	if err2 != nil {
		return fmt.Errorf(stderr.Request, err2.Error())
	}

	if res.StatusCode != 201 {
		return fmt.Errorf(stderr.ResponseCode, res.StatusCode, projectUrl)
	}

	b, err3 := io.ReadAll(res.Body)
	if err3 != nil {
		return fmt.Errorf(stderr.CouldNotReadResponse, err3.Error())
	}

	log.Infof(stdout.ResponseBody, b)

	tpr := &TriggerPipelineResponse{}
	if e := json.Unmarshal(b, tpr); e != nil {
		return fmt.Errorf(stderr.CouldNotJsonDecode, e.Error())
	}

	log.Logf(stdout.TriggerPipeline,
		eVars["PARAM_CIRCLECI_APP_HOST"],
		VcsType,
		CircleProjectUsername,
		CircleProjectReponame,
		tpr.Number,
	)

	return nil
}

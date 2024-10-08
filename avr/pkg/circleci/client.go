package circleci

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/env"
	"github.com/kohirens/stdlib/log"
	"io"
	"net/http"
)

const (
	pipelineFmt = "%s/api/v2/project/%s/%s/%s/pipeline"
)

var ( // Set the default here and update them in init with environment variables.
	// ApiHost CircleCI API hostname, defaults to the public CircleCI host.
	ApiHost = "https://api.circleci.com"
	// AppHost CircleCI API hostname, defaults to the public CircleCI host.
	AppHost = "https://app.circleci.com"
)

type Client struct {
	ApiHost         string
	AppHost         string
	Http            *http.Client
	Project         string
	Token           string
	ProjectUsername string
	VcsType         string
}

func NewClient(vcs string, client *http.Client) *Client {
	repo := envs[EnvProjectRepoName]
	owner := envs[EnvProjectUsername]
	token := envs[EnvToken]

	return &Client{
		ApiHost:         env.Get(EnvApiHost, ApiHost),
		AppHost:         env.Get(EnvAppHost, AppHost),
		Http:            client,
		Project:         repo,
		Token:           token,
		ProjectUsername: owner,
		VcsType:         vcs,
	}
}

// RunWorkflow Trigger a pipeline based on parameters set. See docs
// https://circleci.com/docs/api/v2/index.html#tag/Pipeline; also see an
// example:
// https://circleci.com/docs/triggers-overview/#run-a-pipeline-using-the-api
// Note: keep in mind that you have to use a personal API token; project tokens
// are currently not supported on CircleCI API (v2) at this time.
// see: https://circleci.com/docs/api/v2/#operation/triggerPipeline
func (c *Client) RunWorkflow(branch, nameWorkflow string) error {
	// Build pipeline parameters to trigger the tag-and-release workflow.
	pp, e1 := PipelineParameters(branch, nameWorkflow)
	if e1 != nil {
		return e1
	}

	log.Logf(stdout.TriggerWorkflow, nameWorkflow)

	projectUrl := fmt.Sprintf(
		pipelineFmt,
		c.ApiHost,
		c.VcsType,
		c.ProjectUsername,
		c.Project,
	)

	log.Logf(stdout.CircleProjectUrl, projectUrl)

	req, err1 := http.NewRequest("POST", projectUrl, bytes.NewReader(pp))
	if err1 != nil {
		return fmt.Errorf(stderr.CouldNotPostRequest, err1.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Circle-Token", c.Token)
	req.Header.Set(
		"Authorization",
		"Basic "+base64.StdEncoding.EncodeToString([]byte(c.Token+":")))

	res, err2 := c.Http.Do(req)
	if err2 != nil {
		return fmt.Errorf(stderr.Request, err2.Error())
	}

	if res.StatusCode != 201 {
		return fmt.Errorf(stderr.ResponseCode, res.StatusCode, projectUrl, res.Status)
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
		c.AppHost,
		c.VcsType,
		c.ProjectUsername,
		c.Project,
		tpr.Number,
	)

	return nil
}

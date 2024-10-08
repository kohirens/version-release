package circleci

var stderr = struct {
	CouldNotJsonEncode,
	CouldNotJsonDecode,
	CouldNotPostRequest,
	CouldNotReadResponse,
	Request,
	ResponseCode,
	TokenNotSet string
}{
	CouldNotJsonEncode:   "could not encode %t to JSON: %v",
	CouldNotJsonDecode:   "could not decode JSON: %s",
	CouldNotPostRequest:  "could not POST request: %v",
	CouldNotReadResponse: "could not read response body: %v",
	Request:              "could not make request: %s",
	ResponseCode:         "got a %d response from %v: %v",
	TokenNotSet:          "CircleCI Token not set",
}

var stdout = struct {
	CircleProjectUrl,
	PipelineParams,
	ResponseBody,
	TriggerPipeline,
	TriggerWorkflow string
}{
	CircleProjectUrl: "circleci project URL: %s",
	PipelineParams:   "pipeline parameters are %s",
	ResponseBody:     "response body:\n%s",
	TriggerPipeline:  "triggered pipeline %s/jobs/%s/%s/%s/%d",
	TriggerWorkflow:  "trigger workflow %v",
}

package circleci

var stderr = struct {
	CouldNotJsonEncode   string
	CouldNotJsonDecode   string
	CouldNotPostRequest  string
	CouldNotReadResponse string
	Request              string
	ResponseCode         string
}{
	CouldNotJsonEncode:   "could not encode %t to JSON: %v",
	CouldNotJsonDecode:   "could not decode JSON: %s",
	CouldNotPostRequest:  "could not POST request: %v",
	CouldNotReadResponse: "could not read response body: %v",
	Request:              "could not make request: %s",
	ResponseCode:         "got a %d response from %v: %v",
}

var stdout = struct {
	CircleProjectUrl string
	PipelineParams   string
	ResponseBody     string
	TriggerPipeline  string
}{
	CircleProjectUrl: "circleci project URL: %s",
	PipelineParams:   "pipeline parameters are %s",
	ResponseBody:     "response body:\n%s",
	TriggerPipeline:  "triggered pipeline %s/jobs/%s/%s/%s/%d",
}

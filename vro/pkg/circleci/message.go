package circleci

var stderr = struct {
	CouldNotJsonEncode   string
	CouldNotJsonDecode   string
	CouldNotPostRequest  string
	CouldNotReadResponse string
	Request              string
	ResponseCode         string
}{
	CouldNotJsonEncode:   "could not encode %t to JSON: %v\n",
	CouldNotJsonDecode:   "could not decode JSON: %s\n",
	CouldNotPostRequest:  "could not POST request: %v\n",
	CouldNotReadResponse: "could not read response body: %v\n",
	Request:              "could not make request: %s\n",
	ResponseCode:         "got a %d response from %v: %v",
}

var stdout = struct {
	CircleProjectUrl string
	PipelineParams   string
	ResponseBody     string
	TriggerPipeline  string
}{
	CircleProjectUrl: "circleci project URL: %s\n",
	PipelineParams:   "pipeline parameters are %s\n",
	ResponseBody:     "response body:\n%s\n",
	TriggerPipeline:  "triggered pipeline %s/jobs/%s/%s/%s/%d\n",
}

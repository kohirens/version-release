package main

type Parameters struct {
	TriggeredFlow string `json:"triggered_flow"`
}

type PipelineParams struct {
	Branch     string     `json:"branch"`
	Parameters Parameters `json:"parameters"`
}

type tmplVars struct {
	Data map[string]string
}

type postData map[string]string

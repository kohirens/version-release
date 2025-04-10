package main

var stderr = struct {
	CannotLoad404Page,
	CouldNotDecodeJson,
	CouldNotReadBody,
	FileNotFound,
	FindMock,
	FileOpen,
	FileWrite,
	InvalidLoginState,
	MakeDir,
	MockExist,
	MockLoad,
	MockStatus,
	MockWrite,
	NoMatch,
	TemplateFind,
	TemplateLoad string
}{
	CannotLoad404Page:  "cannot load 404 page %q",
	CouldNotDecodeJson: "could not decode JSON: %s",
	CouldNotReadBody:   "could not read body: %s",
	FileNotFound:       "%q not found",
	FindMock:           "could not find mock response: %v",
	FileOpen:           "could not open file %v",
	FileWrite:          "could not write file %v",
	MockExist:          "mock %v exists %v",
	MakeDir:            "could not make dir: %v",
	MockLoad:           "could not load mock %v: %v",
	MockStatus:         "failed to convert status string to int: %v",
	MockWrite:          "could not write mock response: %v",
	NoMatch:            "re %V did not match %v",
	TemplateFind:       "could not find template: %v",
	TemplateLoad:       "could not load template: %v",
}

var stdout = struct {
	CaptureRequest,
	LoadingFile,
	MockLoad,
	MockStatus,
	TemplateLoad string
}{
	CaptureRequest: "capture request info into %v",
	LoadingFile:    "loading %s file %s",
	MockLoad:       "loaded mock response: %v",
	MockStatus:     "loaded mock response Status: %v",
	TemplateLoad:   "loaded template: %v",
}

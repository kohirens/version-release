package main

import (
	"github.com/kohirens/stdlib/git"
	"github.com/kohirens/version-release/vro/pkg/github"
	"os"
	"testing"
)

type ghMock struct {
}

func (gh *ghMock) TagAndRelease(branch, tag string) (*github.ReleasesResponse, error) {
	return &github.ReleasesResponse{
		Name: "0.1.0 - 2023-06-04",
	}, nil
}

func (gh *ghMock) PublishChangelog(wd, branch, chaneLogFile, msg string) error {
	return nil
}

func TestWorkflow_PublishReleaseTag(t *testing.T) {
	mc := &ghMock{}
	repo := git.CloneFromBundle("repo-06", "tmp", "testdata", string(os.PathSeparator))
	wf := NewWorkflow("citoken", mc)

	err1 := wf.PublishReleaseTag("main", repo, "")
	if err1 == nil {
		t.Errorf("PublishReleaseTag() expected error, got = %v", err1)
	}
}

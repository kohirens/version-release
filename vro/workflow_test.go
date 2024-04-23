package main

import (
	"github.com/kohirens/stdlib/test"
	"github.com/kohirens/version-release-orb/vro/pkg/github"
	"github.com/kohirens/version-release-orb/vro/pkg/gittoolbelt"
	"os"
	"testing"
)

type ghMock struct {
}

func (gh *ghMock) TagAndRelease(branch string, info *gittoolbelt.SemverInfo) (*github.ReleasesResponse, error) {
	return &github.ReleasesResponse{
		Name: "0.1.0 - 2023-06-04",
	}, nil
}

func (gh *ghMock) PublishChangelog(wd, branch, chaneLogFile string) error {
	return nil
}

func TestWorkflow_PublishReleaseTag(t *testing.T) {
	mc := &ghMock{}
	repo := test.SetupARepository("repo-06", "tmp", "testdata", string(os.PathSeparator))
	wf := NewWorkflow("citoken", mc)

	err1 := wf.PublishReleaseTag("main", repo)
	if err1 != nil {
		t.Errorf("PublishReleaseTag() error = %v", err1)
	}
}

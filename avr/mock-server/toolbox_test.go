package main

import "testing"

func Test_getResponseMock(t *testing.T) {
	cases := []struct {
		name    string
		repo    string
		p       string
		wantErr bool
	}{
		{
			"url1",
			"github-repo-commit-message",
			"git/ref/heads/auto-update-changelog",
			false,
		},
		{
			"url2",
			"non-existent",
			"does-not-exist",
			true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getResponseMock(c.repo, c.p, nil)

			if (got != nil) != c.wantErr {
				t.Errorf("loadMockResponse() got error %v, want %v", got, c.wantErr)
				return
			}
		})
	}
}

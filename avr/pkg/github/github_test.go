package github

import "testing"

func Test_parseRepositoryUri(t *testing.T) {
	tests := []struct {
		name  string
		uri   string
		want  string
		want1 string
		want2 string
	}{
		{"https", "https://github.com/kohirens/version-release", "github.com", "kohirens", "version-release"},
		{"git", "git@github.com:kohirens/version-release.git", "github.com", "kohirens", "version-release"},
		{"git", "git@github.com:kohirens/git-tool-belt.git", "github.com", "kohirens", "git-tool-belt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := parseRepositoryUri(tt.uri)
			if got != tt.want {
				t.Errorf("parseRepositoryUri() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseRepositoryUri() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("parseRepositoryUri() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

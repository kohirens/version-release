package lib

import (
	"os"
	"testing"
)

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
			got, got1, got2 := ParseRepositoryUri(tt.uri)
			if got != tt.want {
				t.Errorf("ParseRepositoryUri() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseRepositoryUri() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("ParseRepositoryUri() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_getRequiredEnvVars(t *testing.T) {
	_ = os.Setenv("TEST_VAR1", "123")

	tests := []struct {
		name      string
		eVarNames []string
		want      string
		wantErr   bool
	}{
		{"hasEnvVar", []string{"TEST_VAR1"}, "123", false},
		{"doesNotHaveEnvVar", []string{"TEST_VAR2"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRequiredEnvVars(tt.eVarNames)

			if tt.wantErr != (err != nil) {
				t.Errorf("GetRequiredEnvVars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {

				for _, k := range tt.eVarNames {
					v, ok := got[k]
					if !ok {
						t.Errorf("key %s not in map", k)
					}
					if v != tt.want {
						t.Errorf("got %s, want %s", v, k)
					}
				}
			}
		})
	}
}

// This file container mock test that needs to run under the mock-server
// virtualization container.
package main

import (
	"bytes"
	"fmt"
	git2 "github.com/kohirens/stdlib/git"
	"github.com/kohirens/version-release/vro/pkg/circleci"
	"github.com/kohirens/version-release/vro/pkg/git"
	"github.com/kohirens/version-release/vro/pkg/github"
	"os"
	"strings"
	"testing"

	help "github.com/kohirens/stdlib/test"
)

const (
	dirMode = 0777
	ps      = string(os.PathSeparator)
	// subCmdFlags space separated list of command line flags.
	subEnvVarName = "RECURSIVE_TEST_FLAGS"
)

var (
	fixtureDir = "testdata"
	tmpDir     = help.AbsPath("tmp")
)

func TestMain(m *testing.M) {
	// Only runs when this environment variable is set.
	if _, ok := os.LookupEnv(subEnvVarName); ok {
		runAppMain()
	}

	help.ResetDir(tmpDir, dirMode)

	// Run all tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestCallingMain(tester *testing.T) {
	var tests = []struct {
		name     string
		wantCode int
		args     []string
	}{
		{"versionFlag", 0, []string{"-version"}},
		{"helpFlag", 0, []string{"-help"}},
	}

	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {
			cmd := help.GetTestBinCmd(subEnvVarName, test.args)

			_, _ = help.VerboseSubCmdOut(cmd.CombinedOutput())

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if got != test.wantCode {
				t.Errorf("got %q, want %q", got, test.wantCode)
			}
		})
	}
}

// Will trigger the workflow to publish the changelog.
func TestWorkflowSelector_PublishChangelog(t *testing.T) {
	wantCode := 0

	repo := git2.CloneFromBundle("repo-01", tmpDir, fixtureDir, ps)
	// This git commit has changes where the change log needs updating and there is something to tag
	fixedArgs := []string{"workflow-selector", "CHANGELOG.md", "main", repo, "e9daa3786634ae4ae1346c65dd46247c08eb8416"}

	cmd := help.GetTestBinCmd(subEnvVarName, fixedArgs)

	_, _ = help.VerboseSubCmdOut(cmd.CombinedOutput())

	// get exit code.
	got := cmd.ProcessState.ExitCode()

	if got != wantCode {
		t.Errorf("got %d, want %d", got, wantCode)
	}
}

// Will trigger the workflow to publish a release tag
func TestWorkflowSelector_TagRelease(t *testing.T) {
	wantCode := 0

	repo := git2.CloneFromBundle("repo-02", tmpDir, fixtureDir, ps)
	// This git commit has no changelog updates but there is a commit to tag
	fixedArgs := []string{"workflow-selector", "CHANGELOG.md", "main", repo, "9a9d4945706632b75f3ed6f1df93d8a166472455"}

	cmd := help.GetTestBinCmd(subEnvVarName, fixedArgs)

	so, _ := help.VerboseSubCmdOut(cmd.CombinedOutput())

	// get exit code.
	got := cmd.ProcessState.ExitCode()

	if got != wantCode {
		t.Errorf("got %d, want %d", got, wantCode)
	}

	expectSo := []byte(fmt.Sprintf(stdout.TriggerWorkflow, publishReleaseTagWorkflow))
	if !bytes.Contains(so, expectSo) {
		t.Errorf("stdout did not contain expected %q", expectSo)
	}
}

func TestTriggeredPublishReleaseTagWorkflow(t *testing.T) {
	_ = os.Setenv(circleci.EnvRepoUrl, "git@github.com:kohirens/version-release.git")
	defer os.Unsetenv(circleci.EnvRepoUrl)

	wantCode := 1

	repo := git2.CloneFromBundle("repo-03", tmpDir, fixtureDir, ps)
	// This git commit has no changelog updates but there is a commit to tag
	fixedArgs := []string{"publish-release-tag", "main", repo}

	cmd := help.GetTestBinCmd(subEnvVarName, fixedArgs)

	_, _ = help.VerboseSubCmdOut(cmd.CombinedOutput())

	// get exit code.
	got := cmd.ProcessState.ExitCode()

	if got != wantCode {
		t.Errorf("got %d, want %d", got, wantCode)
	}
}

// Should pick the workflow to publish a release tag
func TestTriggeredPublishChangelogWorkflow(t *testing.T) {
	wantCode := 0

	repo := git2.CloneFromBundle("repo-01", tmpDir, fixtureDir, ps)
	oldUrl, _ := git.RemoteGetUrl(repo, "origin")
	_ = git.RemoteSetUrl(repo, "origin", "https://github.com/kohirens/repo-01", oldUrl)

	// This git commit has changes where the change log needs updating
	fixedArgs := []string{"publish-changelog", "CHANGELOG.md", "main", repo}
	_ = os.Setenv(circleci.EnvRepoUrl, "git@github.com:kohirens/version-release.git")
	_ = os.Setenv(github.EnvMergeType, "rebase")
	_ = os.Setenv(circleci.EnvUsername, "test")

	// run the test
	cmd := help.GetTestBinCmd(subEnvVarName, fixedArgs)

	_, _ = help.VerboseSubCmdOut(cmd.CombinedOutput())

	// get exit code.
	got := cmd.ProcessState.ExitCode()

	// assert
	if got != wantCode {
		t.Errorf("got %d, want %d", got, wantCode)
	}
}

// Should not trigger a pipeline.
// No change log changes to commit or commits to tag.
func TestTriggerPipeline_NoChangelogOrTaggableChanges(t *testing.T) {
	wantCode := 0

	repo := git2.CloneFromBundle("repo-03", tmpDir, fixtureDir, ps)

	fixedArgs := []string{"workflow-selector", "CHANGELOG.md", "main", repo, "ec1556c852c18794f749c5fd67e461e9a142cd03"}

	// run the test
	cmd := help.GetTestBinCmd(subEnvVarName, fixedArgs)

	so, _ := help.VerboseSubCmdOut(cmd.CombinedOutput())

	// get exit code.
	got := cmd.ProcessState.ExitCode()

	// assert
	if got != wantCode {
		t.Errorf("got %d, want %d", got, wantCode)
	}

	want2 := "commit ec1556c852c18794f749c5fd67e461e9a142cd03 is already tagged"
	if !strings.Contains(string(so), want2) {
		t.Errorf("did not find expected message %q in output", want2)
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
			got, err := getRequiredEnvVars(tt.eVarNames)

			if tt.wantErr != (err != nil) {
				t.Errorf("getRequiredEnvVars() error = %v, wantErr %v", err, tt.wantErr)
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

func TestPublishReleaseTagWorkflows(t *testing.T) {
	_ = os.Setenv(circleci.EnvRepoUrl, "git@github.com:kohirens/version-release.git")
	defer os.Unsetenv(circleci.EnvRepoUrl)

	fixedArgs := []string{"publish-release-tag", "main"}
	tests := []struct {
		name     string
		bundle   string
		args     []string
		wantCode int
		contains string
	}{
		{"no-commits-to-tag", "repo-03", fixedArgs, 1, "version is empty"},
		{"no-commits-to-tag-2", "repo-05", fixedArgs, 1, "version is empty"},
		{"first-release", "repo-07", fixedArgs, 0, "releasing 0.1.0"},
		{"has-commits-to-tag", "repo-08", fixedArgs, 0, "releasing 0.1.2"},
		// TODO: Remove `publish-release-tag", "-semver` in the next major release.
		{"specify-a-release", "repo-08", []string{"publish-release-tag", "-semver", "1.0.0", "main"}, 0, "releasing 1.0.0"},
		{"semver-2.0.0", "repo-08", []string{"-semver", "2.0.0", "publish-release-tag", "main"}, 0, "releasing 2.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := git2.CloneFromBundle(tt.bundle, tmpDir, fixtureDir, ps)

			tt.args = append(tt.args, repo)

			cmd := help.GetTestBinCmd(subEnvVarName, tt.args)

			out, _ := help.VerboseSubCmdOut(cmd.CombinedOutput())

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if got != tt.wantCode {
				t.Errorf("PublishReleaseTagWorkflow got %d, want %d", got, tt.wantCode)
			}

			if !bytes.Contains(out, []byte(tt.contains)) {
				t.Errorf("PublishReleaseTagWorkflow got %d, want %d", got, tt.wantCode)
			}
		})
	}
}

// Used for running the application's main function from test in a sub process.
func runAppMain() {
	args := strings.Split(os.Getenv(subEnvVarName), " ")
	os.Args = append([]string{os.Args[0]}, args...)

	// Cannot use testing.Verbose() here since flag.Parse() has not been called.
	// Debug stmt, uncomment when needed.
	//fmt.Printf("\nsub os.Args = %v\n", os.Args)

	main()
}

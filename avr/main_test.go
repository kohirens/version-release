// This file container mock test that needs to run under the mock-server
// virtualization container.
package main

import (
	"bytes"
	"fmt"
	git2 "github.com/kohirens/stdlib/git"
	"github.com/kohirens/version-release/avr/pkg/circleci"
	"github.com/kohirens/version-release/avr/pkg/git"
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
		help.RunMain(subEnvVarName, main)
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
func TestWorkflowSelector_PublishChangelog(runner *testing.T) {
	repo := git2.CloneFromBundle("test-workflow-selector-publish-chglog", tmpDir, fixtureDir, ps)
	expectSo := []byte(fmt.Sprintf("trigger workflow %v", publishChgLogWorkflow))
	// This git commit has changes where the change log needs updating and there is something to tag
	cases := []struct {
		name     string
		args     []string
		wantCode int
		want     []byte
	}{
		{"circleci", []string{"-wd", repo, "workflow-selector", "CHANGELOG.md", "e9daa3786634ae4ae1346c65dd46247c08eb8416"}, 0, expectSo},
		{"github", []string{"-wd", repo, "workflow-selector", "CHANGELOG.md", "e9daa3786634ae4ae1346c65dd46247c08eb8416"}, 0, expectSo},
	}

	for _, c := range cases {
		runner.Run(c.name, func(t *testing.T) {
			t.Setenv(circleci.EnvRepoUrl, fmt.Sprintf("git@github.com:kohirens/test-workflow-selector-publish-chglog.git"))
			cmd := help.GetTestBinCmd(subEnvVarName, c.args)

			so, _ := help.VerboseSubCmdOut(cmd.CombinedOutput())

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if got != c.wantCode {
				t.Errorf("got %d, want %d", got, c.wantCode)
			}

			if !bytes.Contains(so, c.want) {
				t.Errorf("stdout did not contain expected %q", c.want)
			}
		})
	}
}

// Will trigger the workflow to publish a release tag
func TestWorkflowSelector_TagRelease(runner *testing.T) {
	repo := git2.CloneFromBundle("repo-02", tmpDir, fixtureDir, ps)
	argsFixture := []string{"-wd", repo, "workflow-selector", "CHANGELOG.md", "9a9d4945706632b75f3ed6f1df93d8a166472455"}
	expectSo := []byte(fmt.Sprintf("trigger workflow %v", publishReleaseTagWorkflow))

	cases := []struct {
		name     string
		args     []string
		wantCode int
		want     []byte
	}{
		{"circleci", argsFixture, 0, expectSo},
		{"github", argsFixture, 0, expectSo},
	}

	for _, c := range cases {
		runner.Run(c.name, func(t *testing.T) {
			// This git commit has no changelog updates but there is a commit to tag

			cmd := help.GetTestBinCmd(subEnvVarName, c.args)

			so, _ := help.VerboseSubCmdOut(cmd.CombinedOutput())

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if got != c.wantCode {
				t.Errorf("got %d, want %d", got, c.wantCode)
			}

			if !bytes.Contains(so, c.want) {
				t.Errorf("stdout did not contain expected %q", c.want)
			}
		})
	}
}

func TestTriggeredPublishReleaseTagWorkflow(t *testing.T) {
	t.Setenv(circleci.EnvRepoUrl, "git@github.com:kohirens/repo-03.git")

	wantCode := 1

	repo := git2.CloneFromBundle("repo-03", tmpDir, fixtureDir, ps)
	// This git commit has no changelog updates but there is a commit to tag
	fixedArgs := []string{"-wd", repo, "publish-release-tag"}

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
	help.Chdir(t, repo)
	oldUrl, _ := git.RemoteGetUrl(repo, "origin")
	_ = git.RemoteSetUrl(repo, "origin", "https://github.com/kohirens/repo-01", oldUrl)

	// This git commit has changes where the change log needs updating
	fixedArgs := []string{"-wd", repo, "publish-changelog", "CHANGELOG.md"}
	t.Setenv(circleci.EnvRepoUrl, "git@github.com:kohirens/repo-01.git")
	t.Setenv(circleci.EnvProjectRepoName, "repo-01")
	t.Setenv(circleci.EnvProjectUsername, "kohirens")

	// run the test
	cmd := help.GetTestBinCmd(subEnvVarName, fixedArgs)

	_, _ = help.VerboseSubCmdOut(cmd.CombinedOutput())

	// get exit code.
	got := cmd.ProcessState.ExitCode()

	// assert
	if got != wantCode {
		t.Errorf("got %d, want %d", got, wantCode)
		return
	}
}

// Should not trigger a pipeline.
// No change log changes to commit or commits to tag.
func TestTriggerPipeline_NoChangelogOrTaggableChanges(t *testing.T) {
	wantCode := 0

	repo := git2.CloneFromBundle("repo-03", tmpDir, fixtureDir, ps)

	fixedArgs := []string{"-wd", repo, "workflow-selector", "CHANGELOG.md", "ec1556c852c18794f749c5fd67e461e9a142cd03"}

	// run the test
	cmd := help.GetTestBinCmd(subEnvVarName, fixedArgs)

	so, _ := help.VerboseSubCmdOut(cmd.CombinedOutput())

	// get exit code.
	got := cmd.ProcessState.ExitCode()

	// assert
	if got != wantCode {
		t.Errorf("got %d, want %d", got, wantCode)
		return
	}

	want2 := "commit ec1556c852c18794f749c5fd67e461e9a142cd03 is already tagged"
	if !strings.Contains(string(so), want2) {
		t.Errorf("did not find expected message %q in output", want2)
		return
	}
}

func TestPublishReleaseTagWorkflows(t *testing.T) {
	t.Setenv(circleci.EnvRepoUrl, "git@github.com:kohirens/version-release.git")

	fixedArgs := []string{"-wd", "", "publish-release-tag"}
	ghFixedArgs := []string{"-wd", "", "-cicd", "github", "publish-release-tag"}
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
		{"specify-a-release", "repo-08", []string{"-wd", "", "-semver", "1.0.0", "publish-release-tag"}, 0, "releasing 1.0.0"},
		{"semver-2.0.0", "repo-08", []string{"-wd", "", "-semver", "2.0.0", "publish-release-tag"}, 0, "releasing 2.0.0"},
		{"gh-no-commits-to-tag", "repo-03", ghFixedArgs, 1, "version is empty"},
		{"gh-no-commits-to-tag-2", "repo-05", ghFixedArgs, 1, "version is empty"},
		{"gh-first-release", "repo-07", ghFixedArgs, 0, "releasing 0.1.0"},
		{"gh-has-commits-to-tag", "repo-08", ghFixedArgs, 0, "releasing 0.1.2"},
		{"gh-specify-a-release", "repo-08", []string{"-wd", "", "-semver", "1.0.0", "publish-release-tag"}, 0, "releasing 1.0.0"},
		{"gh-semver-2.0.0", "repo-08", []string{"-wd", "", "-semver", "2.0.0", "publish-release-tag"}, 0, "releasing 2.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := git2.CloneFromBundle(tt.bundle, tmpDir, fixtureDir, ps)

			tt.args[1] = repo

			cmd := help.GetTestBinCmd(subEnvVarName, tt.args)

			out, _ := help.VerboseSubCmdOut(cmd.CombinedOutput())

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if got != tt.wantCode {
				t.Errorf("PublishReleaseTagWorkflow got %d, want %d", got, tt.wantCode)
				return
			}

			if !bytes.Contains(out, []byte(tt.contains)) {
				t.Errorf("PublishReleaseTagWorkflow got %v, want %v", string(out), tt.contains)
				return
			}
		})
	}
}

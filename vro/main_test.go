// This file container mock test that needs to run under the mock-server
// virtualization container.
package main

import (
	"github.com/kohirens/version-release-orb/vro/pkg/git"
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

// Will pick the workflow to publish the changelog.
func TestWorkflowSelector_PublishChangelog(t *testing.T) {
	wantCode := 0

	repo := help.SetupARepository("repo-01", tmpDir, fixtureDir, ps)
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

// Should pick the workflow to publish a release tag
func TestWorkflowSelector_TagRelease(t *testing.T) {
	wantCode := 0

	repo := help.SetupARepository("repo-02", tmpDir, fixtureDir, ps)
	// This git commit has no changelog updates but there is a commit to tag
	fixedArgs := []string{"workflow-selector", "CHANGELOG.md", "main", repo, "0541de58335c312459f0783aab09e0797fc824e5"}

	cmd := help.GetTestBinCmd(subEnvVarName, fixedArgs)

	_, _ = help.VerboseSubCmdOut(cmd.CombinedOutput())

	// get exit code.
	got := cmd.ProcessState.ExitCode()

	if got != wantCode {
		t.Errorf("got %d, want %d", got, wantCode)
	}
}

func TestTriggeredPublishReleaseTagWorkflow(t *testing.T) {
	wantCode := 0

	repo := help.SetupARepository("repo-03", tmpDir, fixtureDir, ps)
	// This git commit has no changelog updates but there is a commit to tag
	fixedArgs := []string{"publish-release-tag", "CHANGELOG.md", "main", repo}
	_ = os.Setenv("CIRCLE_REPOSITORY_URL", "git@github.com:kohirens/version-release-orb.git")

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

	repo := help.SetupARepository("repo-01", tmpDir, fixtureDir, ps)
	oldUrl, _ := git.RemoteGetUrl(repo, "origin")
	_ = git.RemoteSetUrl(repo, "origin", "https://github.com/kohirens/repo-01", oldUrl)

	// This git commit has changes where the change log needs updating
	fixedArgs := []string{"publish-changelog", "CHANGELOG.md", "main", repo}
	_ = os.Setenv("CIRCLE_REPOSITORY_URL", "git@github.com:kohirens/version-release-orb.git")
	_ = os.Setenv("PARAM_MERGE_TYPE", "rebase")
	_ = os.Setenv("CIRCLE_USERNAME", "test")

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

	repo := help.SetupARepository("repo-03", tmpDir, fixtureDir, ps)

	fixedArgs := []string{"workflow-selector", "CHANGELOG.md", "main", repo, "285543c691b57f334644d1c29f9288b52645cd08"}

	// run the test
	cmd := help.GetTestBinCmd(subEnvVarName, fixedArgs)

	so, _ := help.VerboseSubCmdOut(cmd.CombinedOutput())

	// get exit code.
	got := cmd.ProcessState.ExitCode()

	// assert
	if got != wantCode {
		t.Errorf("got %d, want %d", got, wantCode)
	}

	want2 := "no commits to tag"
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

// Used for running the application's main function from test in a sub process.
func runAppMain() {
	args := strings.Split(os.Getenv(subEnvVarName), " ")
	os.Args = append([]string{os.Args[0]}, args...)

	// Cannot use testing.Verbose() here since flag.Parse() has not been called.
	// Debug stmt, uncomment when needed.
	//fmt.Printf("\nsub os.Args = %v\n", os.Args)

	main()
}

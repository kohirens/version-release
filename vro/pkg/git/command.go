package git

import (
	"bytes"
	"fmt"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release/vro/internal/util"
	"regexp"
	"strings"
)

const (
	cmdGit = "git"
)

// CheckoutBranch check out a branch, making it if it does not exist.
// git checkout -b <branch_name>
func CheckoutBranch(wd, branch string) error {
	status, se, _, _ := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"checkout", "-b", branch},
	)

	if se != nil {
		return fmt.Errorf(stderr.CouldNotCheckoutBranch, status, se.Error())
	}

	log.Logf(stdout.Status, status)

	return nil
}

// Commit Commits any current staged changes.
// git commit -m "${mergeBranchCommitMsg}" -m "automated update of ${PARAM_CHANGELOG_FILE}"
func Commit(wd string, messages ...string) error {
	ms := []string{}
	for _, m := range messages {
		msg := fmt.Sprintf("-m %q", m)
		ms = append(ms, msg)
	}

	status, se, _, _ := cli.RunCommand(
		wd,
		cmdGit,
		append([]string{"commit"}, ms...),
	)

	if se != nil {
		return fmt.Errorf(stderr.CouldNotCommit, status, se.Error())
	}

	log.Logf(stdout.Status, status)

	return nil
}

// ConfigGlobal Set or return a config global value.
// git config --global <key> <value>
func ConfigGlobal(wd, key, val string) error {
	so, se, _, _ := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"config", "--global", key, val},
	)

	if se != nil {
		return fmt.Errorf(stderr.CouldNotSetGlobalConfig, so, se.Error())
	}

	log.Logf(stdout.SetGitGlobalConfig, key)

	return nil
}

// DoesBranchExistRemotely return the result of:
// git ls-remote --heads <repository_url> <branch_name>
func DoesBranchExistRemotely(wd, uri, branch string) bool {
	status, se, _, _ := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"ls-remote", "--heads", uri, branch},
	)
	if se != nil {
		log.Logf(stderr.CouldNotGitListRemote, status, se.Error())
		return false
	}

	//3f6010eae0788e74f2ae724a7f449fbdcc8d78bf        refs/heads/main
	log.Logf(stdout.FoundRemoteBranch, status)

	return len(bytes.Trim(status, "\n")) > 0
}

func GetCurrentTag(wd string) string {
	so, se, _, _ := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"tag", "-l", "--sort", "-creatordate"},
	)

	if se != nil {
		return ""
	}

	if len(so) < 1 {
		return ""
	}

	tags := bytes.Split(so, []byte("\n"))
	if len(tags) < 0 {
		return ""
	}
	return string(tags[0])
}

// HasSemverTag Indicates when a commit is tagged.
//
//	Uses git describe to finds the most recent tag that is reachable from a
//	commit. Only shows annotated tags.
func HasSemverTag(wd, commit string) bool {
	so, se, _, _ := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"describe", "--contains", commit},
	)

	if se != nil {
		if bytes.Contains(so, []byte("cannot describe")) { // handle expected error
			log.Logf(stdout.NoTags, commit)
		} else {
			log.Errf(stderr.GitDescribeContains, so, se.Error())
		}
		return false
	}

	log.Logf(stdout.TagsInfo, so)

	re := regexp.MustCompile(util.CheckSemVer)

	return re.Match(so)
}

// IsCommit Verify a hash is a commit.
func IsCommit(wd, commit string) bool {
	so, se, _, _ := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"cat-file", "-t", commit},
	)

	if se != nil {
		log.Errf(stderr.CatFile, commit, se.Error())
		return false
	}

	log.Infof(stdout.CatFile, so)

	return bytes.Contains(so, []byte("commit"))
}

// LastLog Return the last commit log.
func LastLog(wd string) string {
	so, se, _, _ := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"log", "-1"},
	)

	if se != nil {
		log.Errf(stderr.LastLog, se.Error())
		return ""
	}

	return string(so)
}

// Log Return a commit log.
func Log(wd, refId string) string {
	so, se, _, _ := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"log", refId, "-1"},
	)

	if se != nil {
		log.Errf(stderr.CommitLog, se.Error())
		return ""
	}

	return string(so)
}

// Push Pushes changes to an origin.
// git push origin <branch_name>
func Push(wd, origin, branch string) error {
	status, se, _, _ := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"push", origin, branch},
	)
	if se != nil {
		return fmt.Errorf(stderr.CouldNotPush, status, se.Error())
	}

	fmt.Println(status)

	return nil
}

// RemoteOriginAdd runs `git remote add <origin> <uri>`
func RemoteOriginAdd(wd, name, uri string) error {
	so, se, _, _ := cli.RunCommand(wd, cmdGit, []string{"remote", "add", name, uri})
	if se != nil {
		return fmt.Errorf(stderr.CouldNotAddOrigin, name, so, se.Error())
	}

	return nil
}

// RemoteOriginRemove Will run the command: `git remote remove <origin>`
func RemoteOriginRemove(wd, name string) error {
	so, se, _, _ := cli.RunCommand(wd, cmdGit, []string{"remote", "remove", name})
	if se != nil {
		return fmt.Errorf(stderr.CouldNotRemoveOrigin, name, so, se.Error())
	}

	return nil
}

// RemoteGetUrl runs `git remote get-url --push origin <name>`
func RemoteGetUrl(wd, name string, flags ...string) (string, error) {
	// reorganize the command line arguments by placing the flags before any arguments.
	args := append(flags, []string{"remote", "get-url", name}...)

	so, se, _, _ := cli.RunCommand(wd, cmdGit, args)
	if se != nil {
		return "", fmt.Errorf(stderr.CouldNotGetRemoteUrl, so, se.Error())
	}

	return strings.Trim(string(so), "\n"), nil
}

// RemoteSetUrl runs `git remote set-url [--push] <name> <newurl> [oldUrl]`
func RemoteSetUrl(wd, name, newUrl, oldUrl string, flags ...string) error {
	// reorganize the command line arguments by placing the flags before any arguments.
	args := append(flags, []string{"remote", "set-url", name, newUrl, oldUrl}...)

	so, se, _, _ := cli.RunCommand(wd, cmdGit, args)

	if se != nil {
		return fmt.Errorf(stderr.CouldNotSetRemoteUrl, so, se.Error())
	}

	return nil
}

// StageFiles git add CHANGELOG.md "${gitChgLogCon}gDir}"
func StageFiles(wd string, files ...string) error {
	status, se, _, _ := cli.RunCommand(wd, cmdGit, append([]string{"add"}, files...))
	if se != nil {
		return fmt.Errorf("%s, %v", status, se.Error())
	}

	log.Logf(stdout.StagedFiles, strings.Join(files, ","))

	return nil
}

// PrintStatus Outputs the git status of the specified directory.
func PrintStatus(wd string) {
	status, e1 := StatusWithOptions(wd, []string{})

	if e1 == nil {
		fmt.Println(string(status))
	}
}

// StatusWithOptions Print git status, pass in some optional flags if needed.
func StatusWithOptions(wd string, options []string) ([]byte, error) {
	opts := append([]string{"status"}, options...)

	status, se, _, _ := cli.RunCommand(
		wd,
		cmdGit,
		opts,
	)

	if se != nil {
		return nil, fmt.Errorf(stderr.CouldNotDisplayGitStatus, status, se.Error())
	}

	return status, nil
}

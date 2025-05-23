package git

import (
	"bytes"
	"fmt"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/version-release/avr/pkg/lib"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	cmdGit = "git"
)

var log = logger.Standard{}

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

	log.Logf(stdout.Status, string(status))

	return nil
}

// CheckoutFileFrom Check out a file at a specific reference.
func CheckoutFileFrom(wd, treeIsh, filename string) error {
	if !IsRef(wd, treeIsh) {
		return fmt.Errorf(stderr.NotACommit, treeIsh)
	}

	so, se, ec, _ := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"checkout", "-f", treeIsh, "--", filename},
	)

	if se != nil {
		return fmt.Errorf(stderr.CouldNotCheckoutFile, filename, se.Error())
	}

	if ec != 0 {
		return fmt.Errorf(stderr.CouldNotCheckoutFile, filename, ec)
	}

	log.Logf(stdout.Status, string(so))

	return nil
}

// Commit Commits any current staged changes.
// git commit -m "${mergeBranchCommitMsg}" -m "automated update of ${PARAM_CHANGELOG_FILE}"
func Commit(wd string, messages string) error {
	log.Dbugf("msg in length = %v", len(messages))

	// commit message via stdin
	sop, sep, sip, e1 := RunCommandInteractive(
		wd,
		cmdGit,
		[]string{"commit", "--file", "-"},
	)
	if e1 != nil {
		return fmt.Errorf(stderr.CouldNotCommit, e1.Error())
	}

	if _, e := io.WriteString(sip, messages); e != nil {
		return fmt.Errorf(stderr.WriteCommit, e.Error())
	}

	if e := sip.Close(); e != nil {
		return fmt.Errorf(stderr.WriteCommit, e.Error())
	}

	se, e3 := io.ReadAll(sep)
	if e3 != nil {
		return fmt.Errorf(stderr.CouldNotCommit, e3.Error())
	}
	if len(se) > 0 {
		return fmt.Errorf(stderr.CouldNotCommit, string(se))
	}

	status, _ := io.ReadAll(sop)
	if len(status) > 0 {
		log.Logf(stdout.Status, string(status))
	}

	return nil
}

// Config Set or return a config global value.
// git config <key> <value>
func Config(wd, key, val string) error {
	so, se, _, _ := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"config", key, val},
	)

	if se != nil {
		return fmt.Errorf(stderr.CouldNotSetGlobalConfig, so, se.Error())
	}

	log.Logf(stdout.SetGitGlobalConfig, key)

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
	status, se, _, cs := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"ls-remote", "--heads", uri, branch},
	)

	log.Dbugf(stdout.Cs, cs)

	if se != nil {
		//log.Errf(stderr.GitListRemote, exitCode, se.Error())
		return false
	}

	log.Dbugf(stdout.FoundRemoteBranch, status)

	return len(bytes.Trim(status, "\n")) > 0
}

func GetCurrentTag(wd string) string {
	so, se, _, _ := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"tag", "--list", "--sort", "-creatordate"},
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

// https://git-scm.com/docs/git-status
type File struct {
	Stage,
	Submodule,
	FileModeHead,
	FileModeIndex,
	FileModeWorktree,
	ObjectNameHead,
	ObjectNameIndex,
	Path string
}

type StatusPorcelainFiles struct {
	Files   []*File
	Deleted []string
}

func (s *StatusPorcelainFiles) Add(xy, sub, mH, mI, mW, hH, hI, path string) {
	s.Files = append(s.Files, &File{
		xy,
		sub,
		mH,
		mI,
		mW,
		hH,
		hI,
		path,
	})
}

// Status Git status, using the porcelain option, of file changes.
// For details see https://git-scm.com/docs/git-status
func Status(wd string) *StatusPorcelainFiles {
	so, se, _, _ := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"status", "--porcelain=2"},
	)

	if se != nil {
		return nil
	}

	if len(so) < 1 {
		return nil
	}

	files := &StatusPorcelainFiles{}

	for _, status := range bytes.Split(so, []byte("\n")) {
		parts := strings.Split(string(status), " ")
		if len(parts) < 9 || parts[0] == "u" {
			continue
		}
		switch parts[0] {
		case "2": // renamed
			// NOTE: We are throwing away the rename or copy score parts[8].
			// parts[9] is of the form <path><sep><origPath>, a tab (ASCII 0x09)
			// byte separates them since `-z` option is NOT used.
			files.Add(parts[1], parts[2], parts[3], parts[4], parts[5], parts[6], parts[7], parts[9])
		default:
			files.Add(parts[1], parts[2], parts[3], parts[4], parts[5], parts[6], parts[7], parts[8])
		}
	}

	return files
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

	re := regexp.MustCompile(lib.CheckSemVer)

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

// IsRef Verify a hash is a commit.
func IsRef(wd, commit string) bool {
	so, se, ec, _ := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"cat-file", "-t", commit},
	)

	if ec != 0 {
		log.Errf(stderr.CatFileExit, commit, ec)
		return false
	}

	if se != nil {
		log.Errf(stderr.CatFile, commit, se.Error())
		return false
	}

	log.Infof(stdout.CatFile, so)

	return true
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
	status, se, _, cmdStr := cli.RunCommand(
		wd,
		cmdGit,
		[]string{"push", origin, branch},
	)

	log.Dbugf(stdout.PushCmd, cmdStr)

	if se != nil {
		return fmt.Errorf(stderr.Push, origin, branch, string(status), se.Error())
	}

	log.Logf(stdout.Push, string(status))

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

// RunCommandInteractive run an external program via CLI interactively, in a
// sub process; passing environment variables along.
//
//	It will pass in the os.Environ(), overwriting key=value pairs from env map,
//	comparison for the key (variable name) is case-sensitive.
func RunCommandInteractive(
	wd,
	program string,
	args []string,
) (io.ReadCloser, io.ReadCloser, io.WriteCloser, error) {
	cmd := exec.Command(program, args...)
	cmd.Dir = wd
	ce := os.Environ()

	// overwrite or set environment variables

	cmd.Env = ce

	cmdIn, err1 := cmd.StdinPipe()
	if err1 != nil {
		return nil, nil, nil, err1
	}

	cmdOut, err2 := cmd.StdoutPipe()
	if err2 != nil {
		return nil, nil, nil, err2
	}

	cmdErr, err3 := cmd.StderrPipe()
	if err3 != nil {
		return nil, nil, nil, err3
	}

	if e := cmd.Start(); e != nil {
		return nil, nil, nil, e
	}

	return cmdOut, cmdErr, cmdIn, nil
}

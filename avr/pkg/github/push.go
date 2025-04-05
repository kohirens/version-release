package github

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Push Construct a new branch with changes and push to GitHub using the REST
// API. Parameters branch and parent must be fully qualified reference (ie:
// refs/heads/master). If they don't start with 'refs' and have at least two
// slashes, it will be rejected.
func Push(name, parent, commitMessage string, files []string, gh *Client) error {
	// Support passing in just the branch name (ie: main),
	parentBranch := parent
	if !strings.Contains(parent, "heads/") {
		parentBranch = "heads/" + parent
	}

	// Get the sha for the base reference.
	parentRef, e1 := GetReference(parentBranch, gh)
	if e1 != nil {
		return e1
	}

	// Get the sha for the base reference.
	parentTree, e1 := GetTree(parentBranch, gh)
	if e1 != nil {
		return e1
	}

	// Make a new Git tree, based off of a parent tree, with new changes.
	// We found out that we had to use a tree instead of just a reference.
	// Otherwise, the new branch would not contain any of bases history.
	newTree, e2 := NewTree(parentTree.Sha, files, gh)
	if e2 != nil {
		return e2
	}

	// Pass in commiter info of the integrators choosing.
	commiterGpg, committer, e3 := GetCommiter()
	if e3 != nil {
		return e3
	}

	// Make a new commit that points to the new tree with changes.
	commit, e4 := NewCommit(nil, committer, commitMessage, commiterGpg, newTree.Sha, []string{parentRef.Object.Sha}, gh)
	if e4 != nil {
		return e4
	}

	Log.Dbugf("new commit sha: %v", commit.Sha)
	Log.Dbugf("new commit tree: %v", commit.Tree)
	Log.Dbugf("new commit comitter: %v", commit.Committer)

	// Make the new branch based on the main truck.
	// The name of the fully qualified reference (ie: refs/heads/my-feature).
	// If it doesn't start with 'refs' and have at least two slashes, it will
	// be rejected.
	newRef, e5 := NewReference(fullRefPrefix+name, commit.Sha, gh)
	if e5 != nil {
		return e5
	}

	Log.Dbugf("newRef.Ref: %v", newRef.Ref)
	Log.Dbugf("newRef.Sha: %v", newRef.Object.Sha)

	return nil
}

func GetCommiter() (string, *UserInfo, error) {
	commiterName, ok1 := os.LookupEnv("PARAM_COMMITTER_NAME")
	if !ok1 {
		return "", nil, fmt.Errorf(stderr.EnvNotSet, "PARAM_COMMITTER_NAME")
	}
	commiterEmail, ok2 := os.LookupEnv("PARAM_COMMITTER_EMAIL")
	if !ok2 {
		return "", nil, fmt.Errorf(stderr.EnvNotSet, "PARAM_COMMITTER_EMAIL")
	}
	commiterGpg, _ := os.LookupEnv("PARAM_COMMITTER_GPG_KEY")

	var commiter *UserInfo
	if commiterGpg != "" && commiterEmail != "" {
		commiter = &UserInfo{
			Name:  commiterName,
			Email: commiterEmail,
			Date:  time.Now().UTC(), //YYYY-MM-DDTHH:MM:SSZ
		}
	}

	return commiterGpg, commiter, nil
}

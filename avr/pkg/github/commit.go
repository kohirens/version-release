package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type GitTree struct {
	Sha       string `json:"sha"`
	Url       string `json:"url"`
	Tree      []Tree `json:"tree"`
	Truncated bool   `json:"truncated"`
}

type Tree struct {
	Path    string `json:"path,omitempty"`
	Mode    string `json:"mode,omitempty"`
	Type    string `json:"type,omitempty"`
	Size    int    `json:"size,omitempty"`
	Sha     string `json:"sha,omitempty"`
	Content string `json:"content,omitempty"`
	Url     string `json:"url,omitempty"`
}

type RequestTree struct {
	BaseTree string  `json:"base_tree"`
	Tree     []*Tree `json:"tree"`
}

type RequestCommit struct {
	Author    *UserInfo `json:"author,omitempty"`
	Committer *UserInfo `json:"committer,omitempty"`
	Message   string    `json:"message"`
	Signature string    `json:"signature,omitempty"`
	Tree      string    `json:"tree"`
	Parents   []string  `json:"parents"`
}

type ResponseCommit struct {
	Sha          string        `json:"sha"`
	NodeId       string        `json:"node_id"`
	Url          string        `json:"url"`
	HtmlUrl      string        `json:"html_url"`
	Author       *UserInfo     `json:"author"`
	Committer    *UserInfo     `json:"committer"`
	Tree         *Tree         `json:"tree"`
	Message      string        `json:"message"`
	Parents      []*Tree       `json:"parents"`
	Verification *Verification `json:"verification"`
}

// GetTree Get a tree using the GitHub REST API
// `get /repos/{owner}/{repo}/git/trees/{tree_sha}`
// Returns a single tree using the SHA1 value or ref name for that tree.
// For details see
// https://docs.github.com/en/rest/git/trees?apiVersion=2022-11-28#get-a-tree
func GetTree(ref string, gh *Client) (*GitTree, error) {
	uri := fmt.Sprintf(epTree+"/"+ref, gh.Host, gh.Repository)

	Log.Dbugf("get tree uri: %v", uri)
	Log.Dbugf("ref: %v", ref)

	res, e1 := gh.Send(uri, "GET", nil)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.CouldNotMakeRequest, e1.Error())
	}

	resBody, e2 := io.ReadAll(res.Body)
	if e2 != nil {
		return nil, fmt.Errorf(stderr.CouldNotReadResponse, e2.Error())
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf(stderr.GetTree, res.Status, string(resBody))
	}

	gt := &GitTree{}
	if e3 := json.Unmarshal(resBody, gt); e3 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonDecode, e3.Error())
	}

	return gt, nil
}

// NewTree Make a new git tree using the GitHub REST API
// `post /repos/{owner}/{repo}/git/trees` endpoint.
// The parentSha is a SHA1 of an existing Git tree object which will be used as
// the base for the new tree. If provided, a new Git tree object will be made
// from entries in the Git tree object pointed to by parentSha and entries
// defined in the tree parameter. Entries defined in the tree parameter will
// overwrite items from base_tree with the same path.
// For details see
// https://docs.github.com/en/rest/git/trees?apiVersion=2022-11-28#create-a-tree
func NewTree(parentSha string, files []string, gh *Client) (*GitTree, error) {
	uri := fmt.Sprintf(epTree, gh.Host, gh.Repository)

	Log.Dbugf("new tree uri: %v", uri)
	Log.Dbugf("new tree parent sha: %v", parentSha)
	Log.Dbugf("files with changes: %v", files)

	body := &RequestTree{
		BaseTree: parentSha,
	}

	// Add contents of the commit to the request body.
	for _, filePath := range files {
		content, e1 := os.ReadFile(filePath)
		if e1 != nil {
			return nil, fmt.Errorf(e1.Error())
		}

		fi, e2 := os.Stat(filePath)
		if e2 != nil {
			return nil, fmt.Errorf(e2.Error())
		}

		// determine the file mode
		fileMode := "100644"
		if fi.IsDir() {
			fileMode = "160000"
		}

		body.Tree = append(body.Tree, &Tree{
			Path:    filePath,
			Mode:    fileMode,
			Type:    "blob",
			Content: string(content),
		})
	}

	reqBody, e1 := json.Marshal(body)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonEncode, body, e1.Error())
	}

	res, e2 := gh.Send(uri, "POST", bytes.NewReader(reqBody))
	if e2 != nil {
		return nil, fmt.Errorf(stderr.CouldNotMakeRequest, e2.Error())
	}

	resBody, e3 := io.ReadAll(res.Body)
	if e3 != nil {
		return nil, fmt.Errorf(stderr.CouldNotReadResponse, e3.Error())
	}

	if res.StatusCode != 201 {
		return nil, fmt.Errorf(stderr.MakeTree, res.Status, string(resBody))
	}

	gt := &GitTree{}
	if e4 := json.Unmarshal(resBody, gt); e4 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonDecode, e4.Error())
	}

	return gt, nil
}

// NewCommit Make a new commit using the GitHub REST API
// `post /repos/{owner}/{repo}/git/commits` endpoint.
// For details see
// https://docs.github.com/en/rest/git/commits?apiVersion=2022-11-28#create-a-commit
//
//	author - By default, the author will be the authenticated user and the
//	current date.
//	committer - By default, committer will use the information set in author.
//	treeSha - The SHA of the tree object this commit points to.
//	parentShas - The full SHAs of the commits that were the parents of this
//	commit. If omitted or empty, the commit will be written as a root commit.
//	For a single parent, an array of one SHA should be provided; for a merge
//	commit, an array of more than one should be provided.
func NewCommit(author, committer *UserInfo, message, gpg, treeSha string, parentShas []string, gh *Client) (*ResponseCommit, error) {
	uri := fmt.Sprintf(epCommit, gh.Host, gh.Repository)

	Log.Dbugf("new commit uri: %v", uri)
	Log.Dbugf("treeSha: %v", treeSha)
	Log.Dbugf("parentSha: %v", parentShas)
	Log.Dbugf("author: %v", author)
	Log.Dbugf("committer: %v", committer)

	body := &RequestCommit{
		Committer: committer,
		Message:   message,
		Tree:      treeSha,
		Signature: gpg,
		Author:    author,
		Parents:   parentShas,
	}

	reqBody, e1 := json.Marshal(body)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonEncode, body, e1.Error())
	}

	res, e2 := gh.Send(uri, "POST", bytes.NewReader(reqBody))
	if e2 != nil {
		return nil, fmt.Errorf(stderr.CouldNotMakeRequest, e2.Error())
	}

	resBody, e3 := io.ReadAll(res.Body)
	if e3 != nil {
		return nil, fmt.Errorf(stderr.CouldNotReadResponse, e3.Error())
	}

	if res.StatusCode != 201 {
		return nil, fmt.Errorf(stderr.NewCommit, res.Status, string(resBody))
	}

	rc := &ResponseCommit{}
	if e4 := json.Unmarshal(resBody, rc); e4 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonDecode, e4.Error())
	}

	return rc, nil
}

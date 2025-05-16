package github

import (
	"fmt"
	jwt "github.com/kohirens/json-web-token"
	"os"
	"time"
)

type TokenPermissions struct {
	Checks       string `json:"checks"`
	Contents     string `json:"contents"`
	Metadata     string `json:"metadata"`
	PullRequests string `json:"pull_requests"`
	Statuses     string `json:"statuses"`
}

type InstallationToken struct {
	Token               string           `json:"token"`
	ExpiresAt           time.Time        `json:"expires_at"`
	Permissions         TokenPermissions `json:"permissions"`
	RepositorySelection string           `json:"repository_selection"`
}

// BuildJWT
// https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/generating-a-json-web-token-jwt-for-a-github-app
// Also see https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/managing-private-keys-for-github-apps#generating-private-keys
func BuildJWT() (string, error) {
	clientId, ok := os.LookupEnv("APP_CLIENT_ID")
	if !ok {
		return "", fmt.Errorf("environment variables APP_CLIENT_ID not set")
	}

	// TODO: Pull this from a vault secrets manager.
	privateKeyPem, ok := os.LookupEnv("PRIVATE_KEY_PEM")
	if !ok {
		return "", fmt.Errorf("environment variables APP_CLIENT_ID not set")
	}

	return jwt.GitHub(clientId, privateKeyPem)
}

//// Installation Get the app installation ID
//func Installation(gh Client) (string, error) {
//	uri := fmt.Sprintf(epAppToken, gh.Host, installationId)
//
//	res, e1 := gh.Send(uri, "POST", nil)
//	if e1 != nil {
//		return "", fmt.Errorf(stderr.InstallationToken, e1.Error())
//	}
//
//	resBody, e2 := io.ReadAll(res.Body)
//	if e2 != nil {
//		return "", fmt.Errorf(stderr.CouldNotReadResponse, e2.Error())
//	}
//
//	token := &InstallationToken{}
//	if e := json.Unmarshal(resBody, token); e != nil {
//		return "", fmt.Errorf(stderr.InstallationToken, e)
//	}
//
//	return installationId, nil
//}
//
//// Token An installation token to interact with the repo.
//func Token(installationId string, gh *Client) (string, error) {
//	uri := fmt.Sprintf(epAppToken, gh.Host, installationId)
//
//	res, e1 := gh.Send(uri, "POST", nil)
//	if e1 != nil {
//		return "", fmt.Errorf(stderr.InstallationToken, e1.Error())
//	}
//
//	resBody, e2 := io.ReadAll(res.Body)
//	if e2 != nil {
//		return "", fmt.Errorf(stderr.CouldNotReadResponse, e2.Error())
//	}
//
//	token := &InstallationToken{}
//	if e := json.Unmarshal(resBody, token); e != nil {
//		return "", fmt.Errorf(stderr.InstallationToken, e)
//	}
//
//	return token.Token, nil
//}

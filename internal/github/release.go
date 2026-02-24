package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ReleaseClient creates GitHub releases via the REST API.
type ReleaseClient struct {
	Token  string
	Repo   string // "owner/repo" from GITHUB_REPOSITORY
	APIURL string // from GITHUB_API_URL, defaults to "https://api.github.com"
}

// NewReleaseClient creates a client from environment variables.
func NewReleaseClient(token string) *ReleaseClient {
	repo := os.Getenv("GITHUB_REPOSITORY")
	apiURL := os.Getenv("GITHUB_API_URL")
	if apiURL == "" {
		apiURL = "https://api.github.com"
	}
	return &ReleaseClient{
		Token:  token,
		Repo:   repo,
		APIURL: apiURL,
	}
}

type createReleaseRequest struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Body       string `json:"body"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
}

// CreateRelease creates a GitHub release for the given tag.
func (c *ReleaseClient) CreateRelease(tag, name, body string, draft, prerelease bool) error {
	url := fmt.Sprintf("%s/repos/%s/releases", c.APIURL, c.Repo)

	payload := createReleaseRequest{
		TagName:    tag,
		Name:       name,
		Body:       body,
		Draft:      draft,
		Prerelease: prerelease,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal release request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("create release: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create release failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

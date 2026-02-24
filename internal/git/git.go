package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/netwarlan/action-semantic-versioning/internal/semver"
)

const commitDelimiter = "---SEMVER-COMMIT-END---"

// RawCommit holds a commit hash and full message.
type RawCommit struct {
	Hash    string
	Message string
}

// Client wraps git operations.
type Client struct {
	WorkDir string
}

// IsShallowRepository checks if the current repo is a shallow clone.
func (c *Client) IsShallowRepository() (bool, error) {
	out, err := c.run("rev-parse", "--is-shallow-repository")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) == "true", nil
}

// FindLatestSemverTag finds the highest semver tag with the given prefix.
func (c *Client) FindLatestSemverTag(prefix string) (string, error) {
	out, err := c.run("tag", "--list", prefix+"*", "--sort=-version:refname")
	if err != nil {
		return "", err
	}

	var latest *semver.Version
	var latestTag string

	for _, line := range strings.Split(out, "\n") {
		tag := strings.TrimSpace(line)
		if tag == "" {
			continue
		}
		v, err := semver.Parse(tag)
		if err != nil {
			continue // skip non-semver tags
		}
		if latest == nil || v.Compare(*latest) > 0 {
			latest = &v
			latestTag = tag
		}
	}

	return latestTag, nil
}

// ListCommitsSince lists all commits since the given tag (or all commits if tag is empty).
func (c *Client) ListCommitsSince(tag string) ([]RawCommit, error) {
	format := fmt.Sprintf("%%H%%n%%B%%n%s", commitDelimiter)

	var args []string
	if tag == "" {
		args = []string{"log", "--format=" + format, "HEAD"}
	} else {
		args = []string{"log", "--format=" + format, tag + "..HEAD"}
	}

	out, err := c.run(args...)
	if err != nil {
		return nil, err
	}

	return parseCommits(out), nil
}

// CreateTag creates a lightweight tag.
func (c *Client) CreateTag(tag string) error {
	_, err := c.run("tag", tag)
	return err
}

// PushTag pushes a tag to the remote.
func (c *Client) PushTag(tag string) error {
	_, err := c.run("push", "origin", tag)
	return err
}

func (c *Client) run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if c.WorkDir != "" {
		cmd.Dir = c.WorkDir
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, string(out))
	}
	return string(out), nil
}

func parseCommits(output string) []RawCommit {
	blocks := strings.Split(output, commitDelimiter)
	var commits []RawCommit

	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		// First line is the hash, rest is the message.
		i := strings.IndexByte(block, '\n')
		if i < 0 {
			// Hash only, no message body.
			commits = append(commits, RawCommit{Hash: block})
			continue
		}

		hash := strings.TrimSpace(block[:i])
		msg := strings.TrimSpace(block[i+1:])
		if hash != "" {
			commits = append(commits, RawCommit{Hash: hash, Message: msg})
		}
	}

	return commits
}

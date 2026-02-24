package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
	}

	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("setup %v: %v\n%s", args, err, out)
		}
	}

	return dir
}

func makeCommit(t *testing.T, dir, message string) {
	t.Helper()
	// Create or modify a file to have something to commit.
	f := filepath.Join(dir, "file.txt")
	data, _ := os.ReadFile(f)
	data = append(data, []byte(message+"\n")...)
	if err := os.WriteFile(f, data, 0644); err != nil {
		t.Fatal(err)
	}

	for _, args := range [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", message},
	} {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("commit %v: %v\n%s", args, err, out)
		}
	}
}

func createTag(t *testing.T, dir, tag string) {
	t.Helper()
	cmd := exec.Command("git", "tag", tag)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("tag %s: %v\n%s", tag, err, out)
	}
}

func TestFindLatestSemverTagEmpty(t *testing.T) {
	dir := setupTestRepo(t)
	makeCommit(t, dir, "initial")

	c := &Client{WorkDir: dir}
	tag, err := c.FindLatestSemverTag("v")
	if err != nil {
		t.Fatal(err)
	}
	if tag != "" {
		t.Errorf("expected empty tag, got %q", tag)
	}
}

func TestFindLatestSemverTag(t *testing.T) {
	dir := setupTestRepo(t)
	makeCommit(t, dir, "first")
	createTag(t, dir, "v1.0.0")
	makeCommit(t, dir, "second")
	createTag(t, dir, "v1.1.0")
	makeCommit(t, dir, "third")
	createTag(t, dir, "v2.0.0")

	c := &Client{WorkDir: dir}
	tag, err := c.FindLatestSemverTag("v")
	if err != nil {
		t.Fatal(err)
	}
	if tag != "v2.0.0" {
		t.Errorf("expected v2.0.0, got %q", tag)
	}
}

func TestFindLatestSemverTagSkipsNonSemver(t *testing.T) {
	dir := setupTestRepo(t)
	makeCommit(t, dir, "first")
	createTag(t, dir, "v1.0.0")
	makeCommit(t, dir, "second")
	createTag(t, dir, "release-1")

	c := &Client{WorkDir: dir}
	tag, err := c.FindLatestSemverTag("v")
	if err != nil {
		t.Fatal(err)
	}
	if tag != "v1.0.0" {
		t.Errorf("expected v1.0.0, got %q", tag)
	}
}

func TestListCommitsSinceTag(t *testing.T) {
	dir := setupTestRepo(t)
	makeCommit(t, dir, "feat: first feature")
	createTag(t, dir, "v1.0.0")
	makeCommit(t, dir, "fix: a bug")
	makeCommit(t, dir, "feat: second feature")

	c := &Client{WorkDir: dir}
	commits, err := c.ListCommitsSince("v1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(commits))
	}
	// Most recent first.
	if commits[0].Message != "feat: second feature" {
		t.Errorf("commit[0].Message = %q", commits[0].Message)
	}
	if commits[1].Message != "fix: a bug" {
		t.Errorf("commit[1].Message = %q", commits[1].Message)
	}
}

func TestListCommitsSinceEmpty(t *testing.T) {
	dir := setupTestRepo(t)
	makeCommit(t, dir, "feat: initial")
	makeCommit(t, dir, "fix: something")

	c := &Client{WorkDir: dir}
	commits, err := c.ListCommitsSince("")
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(commits))
	}
}

func TestCreateTag(t *testing.T) {
	dir := setupTestRepo(t)
	makeCommit(t, dir, "initial")

	c := &Client{WorkDir: dir}
	if err := c.CreateTag("v1.0.0"); err != nil {
		t.Fatal(err)
	}

	// Verify the tag exists.
	cmd := exec.Command("git", "tag", "--list", "v1.0.0")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(out)); got != "v1.0.0" {
		t.Errorf("tag list = %q, want v1.0.0", got)
	}
}

func TestMultilineCommitMessage(t *testing.T) {
	dir := setupTestRepo(t)

	msg := "feat: add user system\n\nThis adds the complete user management.\n\nBREAKING CHANGE: removed legacy auth"
	makeCommit(t, dir, msg)

	c := &Client{WorkDir: dir}
	commits, err := c.ListCommitsSince("")
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}
	if !strings.Contains(commits[0].Message, "BREAKING CHANGE") {
		t.Errorf("expected message to contain BREAKING CHANGE footer, got: %q", commits[0].Message)
	}
}

func TestIsShallowRepository(t *testing.T) {
	dir := setupTestRepo(t)
	makeCommit(t, dir, "initial")

	c := &Client{WorkDir: dir}
	shallow, err := c.IsShallowRepository()
	if err != nil {
		t.Fatal(err)
	}
	if shallow {
		t.Error("expected non-shallow repo")
	}
}

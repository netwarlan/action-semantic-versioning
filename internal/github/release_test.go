package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateReleaseSuccess(t *testing.T) {
	var received createReleaseRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization = %q", got)
		}
		if got := r.URL.Path; got != "/repos/owner/repo/releases" {
			t.Errorf("path = %q", got)
		}

		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte(`{"id": 1}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := &ReleaseClient{
		Token:  "test-token",
		Repo:   "owner/repo",
		APIURL: server.URL,
	}

	err := client.CreateRelease("v1.2.3", "v1.2.3", "changelog", true, false)
	if err != nil {
		t.Fatal(err)
	}

	if received.TagName != "v1.2.3" {
		t.Errorf("TagName = %q", received.TagName)
	}
	if received.Name != "v1.2.3" {
		t.Errorf("Name = %q", received.Name)
	}
	if received.Body != "changelog" {
		t.Errorf("Body = %q", received.Body)
	}
	if !received.Draft {
		t.Error("Draft should be true")
	}
	if received.Prerelease {
		t.Error("Prerelease should be false")
	}
}

func TestCreateReleaseAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if _, err := w.Write([]byte(`{"message": "Validation Failed"}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := &ReleaseClient{
		Token:  "test-token",
		Repo:   "owner/repo",
		APIURL: server.URL,
	}

	err := client.CreateRelease("v1.0.0", "v1.0.0", "", false, false)
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); got == "" {
		t.Error("expected non-empty error message")
	}
}

func TestCreateReleasePrerelease(t *testing.T) {
	var received createReleaseRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte(`{"id": 1}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := &ReleaseClient{
		Token:  "test-token",
		Repo:   "owner/repo",
		APIURL: server.URL,
	}

	err := client.CreateRelease("v1.0.0-beta", "v1.0.0-beta", "", false, true)
	if err != nil {
		t.Fatal(err)
	}
	if !received.Prerelease {
		t.Error("Prerelease should be true")
	}
}

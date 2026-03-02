package ghfd

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-github/v71/github"
)

type mockGitHubClient struct {
	// users maps login -> full profile, used by Get.
	users     map[string]*github.User
	getCounts map[string]int
	followers []*github.User
	following []*github.User
	err       error
}

func (m *mockGitHubClient) Get(ctx context.Context, login string) (*github.User, *github.Response, error) {
	m.getCounts[login]++
	if m.err != nil {
		return nil, nil, m.err
	}
	if u, ok := m.users[login]; ok {
		return u, nil, nil
	}
	return &github.User{Login: github.Ptr(login)}, nil, nil
}

func (m *mockGitHubClient) ListFollowers(ctx context.Context, user string, opts *github.ListOptions) ([]*github.User, *github.Response, error) {
	return m.followers, nil, m.err
}

func (m *mockGitHubClient) ListFollowing(ctx context.Context, user string, opts *github.ListOptions) ([]*github.User, *github.Response, error) {
	return m.following, nil, m.err
}

func newMockGitHubClient(followers []*github.User, following []*github.User, err error) *mockGitHubClient {
	return &mockGitHubClient{
		// Full profiles returned by Get, keyed by login.
		users: map[string]*github.User{
			"me":         {Login: github.Ptr("me"), Name: github.Ptr("My Name")},
			"follower1":  {Login: github.Ptr("follower1"), Name: github.Ptr("Follower One")},
			"both":       {Login: github.Ptr("both"), Name: github.Ptr("")},
			"following1": {Login: github.Ptr("following1"), Name: github.Ptr("Following One")},
		},
		getCounts: make(map[string]int),
		followers: followers,
		following: following,
		err:       err,
	}
}

func TestDiffFollowers(t *testing.T) {
	followers := map[string]string{"user1": "Name 1", "user2": ""}
	following := map[string]string{"user2": "", "user3": "Name 3"}

	notFollowedBack, notFollowingBack := DiffFollowers(followers, following)

	if len(notFollowedBack) != 1 || notFollowedBack["user1"] != "Name 1" {
		t.Errorf("expected notFollowedBack to contain user1, got %v", notFollowedBack)
	}

	if len(notFollowingBack) != 1 || notFollowingBack["user3"] != "Name 3" {
		t.Errorf("expected notFollowingBack to contain user3, got %v", notFollowingBack)
	}
}

func TestRunWithClient(t *testing.T) {
	mock := newMockGitHubClient(
		[]*github.User{
			{Login: github.Ptr("follower1")},
			{Login: github.Ptr("both")},
		},
		[]*github.User{
			{Login: github.Ptr("following1")},
			{Login: github.Ptr("both")},
		},
		nil,
	)

	var out strings.Builder
	err := runWithClient(context.Background(), mock, "me", false, false, &out)
	if err != nil {
		t.Fatalf("runWithClient failed: %v", err)
	}

	output := out.String()
	expectedSubstrings := []string{
		"User: My Name (me)",
		"Followers not followed back:",
		" - follower1",
		"Following not following back:",
		" - following1",
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(output, s) {
			t.Errorf("expected output to contain %q, but it didn't", s)
		}
	}

	if strings.Contains(output, " - both") {
		t.Errorf("expected output not to contain 'both', but it did")
	}
}

func TestRunWithClient_WithList(t *testing.T) {
	mock := newMockGitHubClient(
		[]*github.User{
			{Login: github.Ptr("follower1")},
			{Login: github.Ptr("both")},
		},
		[]*github.User{
			{Login: github.Ptr("following1")},
			{Login: github.Ptr("both")},
		},
		nil,
	)

	var out strings.Builder
	err := runWithClient(context.Background(), mock, "me", true, false, &out)
	if err != nil {
		t.Fatalf("runWithClient failed: %v", err)
	}

	output := out.String()
	expectedSubstrings := []string{
		"Followers:",
		" - follower1",
		"Following:",
		" - following1",
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(output, s) {
			t.Errorf("expected output to contain %q, but it didn't", s)
		}
	}

	if !strings.Contains(output, " - both") {
		t.Errorf("expected output to contain 'both', but it didn't")
	}
}

func TestRunWithClient_WithName(t *testing.T) {
	mock := newMockGitHubClient(
		[]*github.User{
			{Login: github.Ptr("follower1")},
			{Login: github.Ptr("both")},
		},
		[]*github.User{
			{Login: github.Ptr("following1")},
			{Login: github.Ptr("both")},
		},
		nil,
	)

	var out strings.Builder
	err := runWithClient(context.Background(), mock, "me", true, true, &out)
	if err != nil {
		t.Fatalf("runWithClient failed: %v", err)
	}

	output := out.String()
	expectedSubstrings := []string{
		"User: My Name (me)",
		"Followers not followed back:",
		" - Follower One (follower1)",
		"Following not following back:",
		" - Following One (following1)",
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(output, s) {
			t.Errorf("expected output to contain %q, but it didn't", s)
		}
	}

	if !strings.Contains(output, " - both") {
		t.Errorf("expected output to contain 'both', but it didn't")
	}
}

func TestRunWithClient_Deduplicates(t *testing.T) {
	mock := newMockGitHubClient(
		[]*github.User{
			{Login: github.Ptr("follower1")},
			{Login: github.Ptr("both")},
		},
		[]*github.User{
			{Login: github.Ptr("following1")},
			{Login: github.Ptr("both")},
		},
		nil,
	)

	var out strings.Builder
	err := runWithClient(context.Background(), mock, "me", false, true, &out)
	if err != nil {
		t.Fatalf("runWithClient failed: %v", err)
	}

	// "both" appears in both lists; must be resolved exactly once.
	if mock.getCounts["both"] != 1 {
		t.Errorf("expected Get(\"both\") called 1 time, got %d", mock.getCounts["both"])
	}
	// "follower1" is in followers list.
	if mock.getCounts["follower1"] != 1 {
		t.Errorf("expected Get(\"follower1\") called 1 time, got %d", mock.getCounts["follower1"])
	}
	// "following1" is in following list.
	if mock.getCounts["following1"] != 1 {
		t.Errorf("expected Get(\"following1\") called 1 time, got %d", mock.getCounts["following1"])
	}
	// "me" is the target user lookup.
	if mock.getCounts["me"] != 1 {
		t.Errorf("expected Get(\"me\") called 1 time, got %d", mock.getCounts["me"])
	}
}

func TestRunWithClient_None(t *testing.T) {
	mock := newMockGitHubClient(
		[]*github.User{},
		[]*github.User{},
		nil,
	)

	var out strings.Builder
	err := runWithClient(context.Background(), mock, "me", false, false, &out)
	if err != nil {
		t.Fatalf("runWithClient failed: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "- None") {
		t.Errorf("expected output to contain '- None', but it didn't")
	}
}

func TestRunWithClient_Error(t *testing.T) {
	mock := newMockGitHubClient(
		[]*github.User{},
		[]*github.User{},
		fmt.Errorf("api error"),
	)

	var out strings.Builder
	err := runWithClient(context.Background(), mock, "me", false, false, &out)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to get user") {
		t.Errorf("expected error message to contain 'failed to get user', got %v", err)
	}
}

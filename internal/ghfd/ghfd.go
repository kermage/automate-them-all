package ghfd

import (
	"context"
	"fmt"
	"io"

	"github.com/google/go-github/v71/github"
)

type GitHubClient interface {
	Get(ctx context.Context, user string) (*github.User, *github.Response, error)
	ListFollowers(ctx context.Context, user string, opts *github.ListOptions) ([]*github.User, *github.Response, error)
	ListFollowing(ctx context.Context, user string, opts *github.ListOptions) ([]*github.User, *github.Response, error)
}

type ghClientWrapper struct {
	client *github.Client
}

func (w *ghClientWrapper) Get(ctx context.Context, user string) (*github.User, *github.Response, error) {
	return w.client.Users.Get(ctx, user)
}

func (w *ghClientWrapper) ListFollowers(ctx context.Context, user string, opts *github.ListOptions) ([]*github.User, *github.Response, error) {
	return w.client.Users.ListFollowers(ctx, user, opts)
}

func (w *ghClientWrapper) ListFollowing(ctx context.Context, user string, opts *github.ListOptions) ([]*github.User, *github.Response, error) {
	return w.client.Users.ListFollowing(ctx, user, opts)
}

func Run(ctx context.Context, token string, username string, list bool, resolveName bool, out io.Writer) error {
	client := github.NewClient(nil).WithAuthToken(token)
	gh := &ghClientWrapper{client: client}
	return runWithClient(ctx, gh, username, list, resolveName, out)
}

func runWithClient(ctx context.Context, gh GitHubClient, username string, list bool, resolveName bool, out io.Writer) error {
	user, _, err := gh.Get(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	actualUsername := user.GetLogin()
	name := user.GetName()
	if name == "" {
		fmt.Fprintf(out, "User: %s\n", actualUsername)
	} else {
		fmt.Fprintf(out, "User: %s (%s)\n", name, actualUsername)
	}

	followersList, _, err := gh.ListFollowers(ctx, actualUsername, &github.ListOptions{PerPage: 100})
	if err != nil {
		return fmt.Errorf("failed to list followers: %w", err)
	}

	followingList, _, err := gh.ListFollowing(ctx, actualUsername, &github.ListOptions{PerPage: 100})
	if err != nil {
		return fmt.Errorf("failed to list following: %w", err)
	}

	users := prepareList(followersList, followingList)

	if resolveName {
		users, err = resolveNames(ctx, gh, users)
		if err != nil {
			return fmt.Errorf("failed to resolve profiles: %w", err)
		}
	}

	followers := userMap(followersList, users)
	following := userMap(followingList, users)

	if list {
		fmt.Fprintln(out, "\nFollowers:")
		printMap(out, followers)

		fmt.Fprintln(out, "\nFollowing:")
		printMap(out, following)
	}

	notFollowedBack, notFollowingBack := DiffFollowers(followers, following)

	fmt.Fprintln(out, "\nFollowers not followed back:")
	printMap(out, notFollowedBack)

	fmt.Fprintln(out, "\nFollowing not following back:")
	printMap(out, notFollowingBack)

	return nil
}

func DiffFollowers(followers, following map[string]string) (notFollowedBack, notFollowingBack map[string]string) {
	notFollowedBack = make(map[string]string)
	for login, name := range followers {
		if _, ok := following[login]; !ok {
			notFollowedBack[login] = name
		}
	}

	notFollowingBack = make(map[string]string)
	for login, name := range following {
		if _, ok := followers[login]; !ok {
			notFollowingBack[login] = name
		}
	}
	return
}

func prepareList(users ...[]*github.User) map[string]string {
	unique := make(map[string]string)
	for _, list := range users {
		for _, user := range list {
			unique[user.GetLogin()] = ""
		}
	}
	return unique
}

// Fetch full profiles for the union of all provided users.
// A per-user Get failure is tolerated: the user is recorded with an empty name.
func resolveNames(ctx context.Context, gh GitHubClient, users map[string]string) (map[string]string, error) {
	m := make(map[string]string, len(users))
	for login := range users {
		profile, _, err := gh.Get(ctx, login)
		if err != nil {
			m[login] = ""
			continue
		}
		m[login] = profile.GetName()
	}
	return m, nil
}

// Build a login→name map for a single user list using a pre-resolved names cache.
func userMap(users []*github.User, names map[string]string) map[string]string {
	m := make(map[string]string, len(users))
	for _, u := range users {
		login := u.GetLogin()
		m[login] = names[login]
	}
	return m
}

func printMap(out io.Writer, m map[string]string) {
	if len(m) == 0 {
		fmt.Fprintln(out, " - None")
		return
	}
	for login, name := range m {
		if name == "" {
			fmt.Fprintf(out, " - %s\n", login)
		} else {
			fmt.Fprintf(out, " - %s (%s)\n", name, login)
		}
	}
}

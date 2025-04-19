package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v71/github"
)

type GitHub struct {
	*github.Client
}

var ctx = context.Background()
var users = make(map[string]string, 0)

func main() {
	args := os.Args[1:]

	if len(args) == 0 || len(args) > 2 {
		fmt.Println("Usage: gh-user <token> [username]")

		if len(args) > 2 {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

	token := args[0]
	username := ""

	if len(args) == 2 {
		username = args[1]
	}

	gh := &GitHub{github.NewClient(nil).WithAuthToken(token)}
	user, res, err := gh.Users.Get(ctx, username)

	maybeFailed(res, err)
	printUser("User: ", user.GetLogin(), gh.getName(username))

	followers := make(map[string]string, 0)
	following := make(map[string]string, 0)

	fmt.Println("Followers:")
	for _, user := range gh.getFollowers(username, 100) {
		username := user.GetLogin()
		displayName := gh.getName(username)
		followers[username] = displayName

		printUser(" - ", username, displayName)
	}

	fmt.Println("Following:")
	for _, user := range gh.getFollowing(username, 100) {
		username := user.GetLogin()
		displayName := gh.getName(username)
		following[username] = displayName

		printUser(" - ", username, displayName)
	}

	fmt.Println("Followers not followed back:")
	for follower, name := range followers {
		if _, ok := following[follower]; !ok {
			printUser(" - ", follower, name)
		}
	}

	fmt.Println("Following not following back:")
	for following, name := range following {
		if _, ok := followers[following]; !ok {
			printUser(" - ", following, name)
		}
	}
}

func maybeFailed(response *github.Response, err error) {
	if err != nil {
		fmt.Println(response.Body)
		os.Exit(1)
	}
}

func printUser(prefix string, username string, displayName string) {
	if displayName == "" {
		fmt.Printf("%s%s\n", prefix, username)
	} else {
		fmt.Printf("%s%s (%s)\n", prefix, displayName, username)
	}
}

func (gh *GitHub) getFollowers(user string, number int) []*github.User {
	followers, response, err := gh.Users.ListFollowers(ctx, user, &github.ListOptions{PerPage: number})

	maybeFailed(response, err)

	return followers
}

func (gh *GitHub) getFollowing(user string, number int) []*github.User {
	following, response, err := gh.Users.ListFollowing(ctx, user, &github.ListOptions{PerPage: number})

	maybeFailed(response, err)

	return following
}

func (gh *GitHub) getName(username string) string {
	if _, ok := users[username]; ok {
		return users[username]
	}

	user, response, err := gh.Users.Get(ctx, username)

	maybeFailed(response, err)

	name := user.GetName()
	users[username] = name

	return name
}

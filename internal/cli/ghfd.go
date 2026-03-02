package cli

import (
	"context"
	"gatas/internal/ghfd"
	"os"

	"github.com/spf13/cobra"
)

var (
	username    string
	list        bool
	resolveName bool
)

var GhfdCmd = &cobra.Command{
	Use:   "ghfd <token>",
	Short: "GitHub followers/following diff tool",
	Long:  `ghfd compares your GitHub followers and following lists and identifies those you don't follow back and those who don't follow you back.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := args[0]
		return ghfd.Run(context.Background(), token, username, list, resolveName, os.Stdout)
	},
}

func SetupGhfdCmd() {
	GhfdCmd.Flags().StringVarP(&username, "username", "u", "", "GitHub username (defaults to authenticated user)")
	GhfdCmd.Flags().BoolVarP(&list, "list", "l", false, "Print the full lists of followers and following")
	GhfdCmd.Flags().BoolVarP(&resolveName, "name", "n", false, "Resolve display names via extra API calls (slower)")
}

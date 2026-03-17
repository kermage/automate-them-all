package cli

import (
	"fmt"
	"gatas/internal/cloak"
	"os"

	"github.com/spf13/cobra"
)

var (
	undo bool
)

var CloakCmd = &cobra.Command{
	Use:   "cloak [target...]",
	Short: "Git local exclusion from tracking",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		return cloak.Run(path, args, !undo, os.Stdout)
	},
}

func SetupCloakCmd() {
	CloakCmd.Flags().BoolVarP(&undo, "undo", "u", false, "Remove from local exclusion")
}

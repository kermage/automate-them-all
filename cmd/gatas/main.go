package main

import (
	"gatas/internal/cli"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gatas",
	Short: "A consolidation of utility scripts.",
}

func init() {
	cli.SetupGhfdCmd()
	cli.SetupCloakCmd()
	cli.SetupInjectCmd()
	rootCmd.AddCommand(cli.GhfdCmd)
	rootCmd.AddCommand(cli.SrtfCmd)
	rootCmd.AddCommand(cli.CloakCmd)
	rootCmd.AddCommand(cli.InjectCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

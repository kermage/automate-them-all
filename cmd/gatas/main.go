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
	rootCmd.AddCommand(cli.GhfdCmd)
	rootCmd.AddCommand(cli.SrtfCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

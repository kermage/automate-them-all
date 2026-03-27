package cli

import (
	"gatas/internal/inject"
	"os"

	"github.com/spf13/cobra"
)

var (
	source string
	target string
)

var InjectCmd = &cobra.Command{
	Use:   "inject [json]",
	Short: "Inject content from a file into JSON",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		json := args[0]
		return inject.Run(json, source, target, os.Stdout)
	},
}

func SetupInjectCmd() {
	InjectCmd.Flags().StringVarP(&source, "source", "s", "", "The file path to the content file")
	InjectCmd.MarkFlagRequired("source")
	InjectCmd.Flags().StringVarP(&target, "target", "t", "", "The dot-syntax path used within the JSON")
}

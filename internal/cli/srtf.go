package cli

import (
	"gatas/internal/srtf"
	"os"

	"github.com/spf13/cobra"
)

var SrtfCmd = &cobra.Command{
	Use:   "srtf [path] [filename]",
	Short: "SRT subtitle file finder and copier",
	Long:  `srtf recursively walks the directory tree to find matching SRT files and copies each match one level up, renamed to the parent directory name.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return srtf.Run(args[0], args[1], os.Stdout)
	},
}

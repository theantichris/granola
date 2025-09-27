package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var timeout time.Duration

// exportCmd is the command to export Granola notes.
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Granola notes to Markdown.",
	Long:  "Export Granola notes to Markdown.",
	RunE:  runExport,
	Args:  cobra.ArbitraryArgs,
}

// init initializes the export command.
func init() {
	RootCmd.AddCommand(exportCmd)

	exportCmd.Flags().DurationVar(&timeout, "timout", 2*time.Minute, "HTTP timeout for API requests")
}

func runExport(cmd *cobra.Command, args []string) error {
	// Check for supabase config value
	// Read supabase file
	// Get token
	// Print token to stdout
	return nil
}

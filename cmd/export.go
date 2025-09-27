package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/theantichris/granola/internal/api"
)

var (
	ErrSupabaseEmpty  = errors.New("supabase cannot be empty")
	ErrSupabaseRead   = errors.New("failed to read supabase.json")
	ErrDocumentExport = errors.New("failed to export documents")
)

// exportCmd is the command to export Granola notes.
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Granola notes to Markdown.",
	Long:  "Export Granola notes to Markdown.",
	RunE:  runExport,
}

// init initializes the export command.
func init() {
	RootCmd.AddCommand(exportCmd)

	var timeout time.Duration

	exportCmd.Flags().DurationVar(&timeout, "timeout", 2*time.Minute, "HTTP timeout for API requests")

	_ = viper.BindPFlag("timeout", exportCmd.Flags().Lookup("timeout"))
}

func runExport(cmd *cobra.Command, args []string) error {
	supabaseFile := viper.GetString("supabase")

	// Check file
	if strings.TrimSpace(supabaseFile) == "" {
		return fmt.Errorf("%w: set the path to supabase.json via flag, config file, or env variable", ErrSupabaseEmpty)
	}

	// Get supabase.json file contents
	supabaseContent, err := getSupabaseContent(supabaseFile)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDocumentExport, err)
	}

	// TODO: Add URL to config.
	timeout := viper.GetDuration("timeout")
	httpClient := http.Client{Timeout: timeout}
	documents, err := api.GetDocuments("https://api.granola.ai/v2/get-documents", supabaseContent, &httpClient)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDocumentExport, err)
	}

	// Print documents to stdout
	fmt.Printf("%v", documents)

	return nil
}

// getSupabaseContent reads the supabase.json file and returns the content as []byte.
func getSupabaseContent(filename string) ([]byte, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrSupabaseRead, err)
	}

	return content, nil
}

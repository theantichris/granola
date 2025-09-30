package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/theantichris/granola/internal/api"
	"github.com/theantichris/granola/internal/writer"
)

var appFS = afero.NewOsFs()

var (
	ErrExportCmdInit  = errors.New("failed to initialize the export command")
	ErrSupabaseEmpty  = errors.New("supabase cannot be empty")
	ErrDocumentExport = errors.New("failed to export documents")
)

// NewExportCmd creates a new ExportCommand and binds its flags.
func NewExportCmd(logger *log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export Granola notes to Markdown.",
		Long:  "Export Granola notes to Markdown files in the specified output directory.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlag("timeout", cmd.Flags().Lookup("timeout")); err != nil {
				return fmt.Errorf("%w: %s", ErrExportCmdInit, err)
			}
			if err := viper.BindPFlag("output", cmd.Flags().Lookup("output")); err != nil {
				return fmt.Errorf("%w: %s", ErrExportCmdInit, err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return exportNotes(logger)
		},
	}

	var timeout time.Duration
	cmd.Flags().DurationVar(&timeout, "timeout", 2*time.Minute, "HTTP timeout for API requests, default 2 minutes")

	var output string
	cmd.Flags().StringVar(&output, "output", "./notes", "Output directory for exported Markdown files")

	return cmd
}

// exportNotes loads the contents of supabase.json and uses it to call and retrieve
// the documents from the Granola API, then writes them to Markdown files.
func exportNotes(logger *log.Logger) error {
	filename := viper.GetString("supabase")

	if strings.TrimSpace(filename) == "" {
		return fmt.Errorf("%w: set the path to supabase.json via flag, config file, or env variable", ErrSupabaseEmpty)
	}

	logger.Info("Reading supabase configuration", "file", filename)
	supabaseContent, err := afero.ReadFile(appFS, filename)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDocumentExport, err)
	}

	// TODO: Add URL to config.
	timeout := viper.GetDuration("timeout")
	logger.Info("Fetching documents from Granola API", "timeout", timeout)
	httpClient := http.Client{Timeout: timeout}
	documents, err := api.GetDocuments("https://api.granola.ai/v2/get-documents", supabaseContent, &httpClient)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDocumentExport, err)
	}

	logger.Info("Retrieved documents", "count", len(documents))

	outputDir := viper.GetString("output")
	logger.Info("Writing documents to Markdown files", "output", outputDir)

	if err := writer.Write(documents, outputDir, appFS); err != nil {
		return fmt.Errorf("%w: %s", ErrDocumentExport, err)
	}

	logger.Info("Export completed successfully", "files", len(documents))

	return nil
}

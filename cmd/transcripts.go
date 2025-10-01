package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/theantichris/granola/internal/cache"
	"github.com/theantichris/granola/internal/transcript"
)

var (
	ErrTranscriptCmdInit  = errors.New("failed to initialize the transcripts command")
	ErrTranscriptExport   = errors.New("failed to export transcripts")
	ErrCacheFileNotFound  = errors.New("cache file not found")
	invalidCharsRegex     = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
)

// NewTranscriptsCmd creates a new transcripts command and binds its flags.
func NewTranscriptsCmd(logger *log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transcripts",
		Short: "Export Granola transcripts to text files.",
		Long:  "Export raw Granola transcripts with timestamps to plain text files in the specified output directory.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlag("transcript-output", cmd.Flags().Lookup("output")); err != nil {
				return fmt.Errorf("%w: %s", ErrTranscriptCmdInit, err)
			}
			if err := viper.BindPFlag("cache-file", cmd.Flags().Lookup("cache")); err != nil {
				return fmt.Errorf("%w: %s", ErrTranscriptCmdInit, err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return writeTranscripts(logger)
		},
	}

	var output string
	cmd.Flags().StringVar(&output, "output", "./transcripts", "Output directory for exported transcript files")

	var cacheFile string
	defaultCachePath := cache.GetDefaultCachePath()
	cmd.Flags().StringVar(&cacheFile, "cache", defaultCachePath, "Path to Granola cache file")

	return cmd
}

// writeTranscripts reads the local cache file and exports raw transcripts with timestamps.
func writeTranscripts(logger *log.Logger) error {
	cacheFile := viper.GetString("cache-file")

	fmt.Println("Reading Granola cache file...")
	logger.Info("Reading Granola cache file", "file", cacheFile)
	cacheData, err := cache.ReadCache(cacheFile)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrTranscriptExport, err)
	}

	logger.Info("Loaded cache data", "documents", len(cacheData.Documents), "transcripts", len(cacheData.Transcripts))

	outputDir := viper.GetString("transcript-output")
	fmt.Printf("Exporting %d transcripts to %s...\n", len(cacheData.Transcripts), outputDir)
	logger.Info("Writing transcripts to files", "output", outputDir)

	// Create output directory
	if err := appFS.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("%w: failed to create output directory: %s", ErrTranscriptExport, err)
	}

	usedFilenames := make(map[string]bool)
	count := 0

	for docID, segments := range cacheData.Transcripts {
		// Skip if no segments
		if len(segments) == 0 {
			continue
		}

		// Get document info
		doc, ok := cacheData.Documents[docID]
		if !ok {
			// Use docID if document not found
			doc = cache.Document{
				ID:    docID,
				Title: docID,
			}
		}

		// Generate filename
		filename := doc.Title
		if filename == "" {
			filename = doc.ID
		}
		filename = sanitizeFilename(filename)
		filename = makeUnique(filename, usedFilenames)
		usedFilenames[filename] = true

		filePath := filepath.Join(outputDir, filename+".txt")

		// Check if file needs updating
		if !shouldUpdateFile(doc, filePath) {
			continue
		}

		// Format transcript
		content := transcript.FormatTranscript(doc, segments)
		if content == "" {
			continue
		}

		// Write file
		if err := afero.WriteFile(appFS, filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("%w: failed to write file %s: %s", ErrTranscriptExport, filePath, err)
		}

		count++
	}

	fmt.Println("✓ Export completed successfully")
	logger.Info("Export completed successfully", "files", count)

	return nil
}

// sanitizeFilename removes invalid characters and limits length.
func sanitizeFilename(name string) string {
	name = invalidCharsRegex.ReplaceAllString(name, "_")
	if len(name) > 100 {
		name = name[:100]
	}
	return name
}

// makeUnique ensures filename is unique by appending a number if needed.
func makeUnique(filename string, used map[string]bool) string {
	if !used[filename] {
		return filename
	}

	counter := 2
	for {
		uniqueName := fmt.Sprintf("%s_%d", filename, counter)
		if !used[uniqueName] {
			return uniqueName
		}
		counter++
	}
}

// shouldUpdateFile checks if the file needs to be updated based on timestamps.
func shouldUpdateFile(doc cache.Document, filePath string) bool {
	fileInfo, err := appFS.Stat(filePath)
	if err != nil {
		// File doesn't exist or other error, write it
		return true
	}

	// Parse document's updated_at timestamp
	docUpdated, err := time.Parse(time.RFC3339, doc.UpdatedAt)
	if err != nil {
		// Can't parse timestamp, write the file to be safe
		return true
	}

	// If document is newer than file, update it
	return docUpdated.After(fileInfo.ModTime())
}

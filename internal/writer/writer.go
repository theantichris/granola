// Package writer provides functionality for writing Granola documents to Markdown files.
package writer

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/afero"
	"github.com/theantichris/granola/internal/api"
	"github.com/theantichris/granola/internal/converter"
)

var invalidFileChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)

// Write writes documents to Markdown files in the specified output directory.
// It only writes files if they don't exist or if the document's updated_at timestamp
// is newer than the existing file's modification time.
func Write(docs []api.Document, outputDir string, fs afero.Fs) error {
	if err := fs.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	usedFilenames := make(map[string]int)

	for _, doc := range docs {
		filename := sanitizeFilename(doc.Title, doc.ID)
		filename = makeUnique(filename, usedFilenames)
		usedFilenames[filename]++

		filePath := filepath.Join(outputDir, filename+".md")

		// Check if file exists and compare timestamps
		shouldWrite, err := shouldUpdateFile(fs, filePath, doc.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to check file status for %s: %w", filePath, err)
		}

		if !shouldWrite {
			continue
		}

		markdown, err := converter.ToMarkdown(doc)
		if err != nil {
			return fmt.Errorf("failed to convert document %s: %w", doc.ID, err)
		}

		if err := afero.WriteFile(fs, filePath, []byte(markdown), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}

// shouldUpdateFile checks if a file should be written based on whether it exists
// and if the document's updated_at timestamp is newer than the file's modification time.
func shouldUpdateFile(fs afero.Fs, filePath string, updatedAt string) (bool, error) {
	// Check if file exists
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		return false, err
	}

	// If file doesn't exist, we should write it
	if !exists {
		return true, nil
	}

	// Parse the document's updated_at timestamp
	docUpdatedAt, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		// If we can't parse the timestamp, write the file to be safe
		return true, nil
	}

	// Get the file's modification time
	fileInfo, err := fs.Stat(filePath)
	if err != nil {
		return false, err
	}

	// Write the file if the document is newer than the existing file
	return docUpdatedAt.After(fileInfo.ModTime()), nil
}

// sanitizeFilename removes invalid characters from a filename and falls back to ID if empty.
func sanitizeFilename(title, id string) string {
	// Use title if available, otherwise use ID
	name := strings.TrimSpace(title)
	if name == "" {
		name = id
	}

	// Replace invalid characters with underscores
	name = invalidFileChars.ReplaceAllString(name, "_")

	// Replace multiple consecutive underscores with a single one
	name = regexp.MustCompile(`_+`).ReplaceAllString(name, "_")

	// Trim underscores from start and end
	name = strings.Trim(name, "_")

	// Ensure we have something
	if name == "" {
		name = "untitled"
	}

	// Limit length to 100 characters for filesystem compatibility
	if len(name) > 100 {
		name = name[:100]
	}

	return name
}

// makeUnique appends a number to a filename if it already exists.
func makeUnique(filename string, used map[string]int) string {
	if count, exists := used[filename]; exists {
		return fmt.Sprintf("%s_%d", filename, count+1)
	}
	return filename
}
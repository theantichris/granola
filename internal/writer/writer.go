package writer

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"
	"github.com/theantichris/granola/internal/api"
	"github.com/theantichris/granola/internal/converter"
)

var invalidFileChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)

// Write writes documents to Markdown files in the specified output directory.
func Write(docs []api.Document, outputDir string, fs afero.Fs) error {
	if err := fs.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	usedFilenames := make(map[string]int)

	for _, doc := range docs {
		markdown, err := converter.ToMarkdown(doc)
		if err != nil {
			return fmt.Errorf("failed to convert document %s: %w", doc.ID, err)
		}

		filename := sanitizeFilename(doc.Title, doc.ID)
		filename = makeUnique(filename, usedFilenames)
		usedFilenames[filename]++

		filepath := filepath.Join(outputDir, filename+".md")
		if err := afero.WriteFile(fs, filepath, []byte(markdown), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filepath, err)
		}
	}

	return nil
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
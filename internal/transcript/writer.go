package transcript

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/afero"
	"github.com/theantichris/granola/internal/api"
	"github.com/theantichris/granola/internal/prosemirror"
)

var (
	invalidCharsRegex = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
	htmlTagRegex      = regexp.MustCompile(`<[^>]*>`)
)

// Write writes documents as plain text transcript files to the specified output directory.
func Write(documents []api.Document, outputDir string, fs afero.Fs) error {
	if err := fs.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	usedFilenames := make(map[string]bool)

	for _, doc := range documents {
		// Check if file needs updating
		filename := generateFilename(doc.Title, doc.ID)
		filename = sanitizeFilename(filename)
		filename = makeUnique(filename, usedFilenames)
		usedFilenames[filename] = true

		filePath := filepath.Join(outputDir, filename+".txt")

		if !shouldUpdateFile(doc, filePath, fs) {
			continue
		}

		// Format transcript with metadata header
		content := formatTranscript(doc)

		// Skip if no content
		if content == "" {
			continue
		}

		if err := afero.WriteFile(fs, filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write transcript file %s: %w", filePath, err)
		}
	}

	return nil
}

// formatTranscript creates a plain text document with a metadata header.
func formatTranscript(doc api.Document) string {
	// Get content with priority: NotesPlain > Notes (ProseMirror→plain text) > OriginalContent (HTML stripped) > Content
	var content string

	if doc.NotesPlain != "" {
		content = strings.TrimSpace(doc.NotesPlain)
	}
	if content == "" && doc.Notes != nil {
		content = strings.TrimSpace(prosemirror.ConvertToPlainText(doc.Notes))
	}
	if content == "" && doc.LastViewedPanel != nil && doc.LastViewedPanel.OriginalContent != "" {
		// Strip HTML tags for plain text output
		content = stripHTML(doc.LastViewedPanel.OriginalContent)
	}
	if content == "" && doc.Content != "" {
		content = doc.Content
	}

	// Skip if no content at all
	if content == "" {
		return ""
	}

	var builder strings.Builder

	// Header
	builder.WriteString(strings.Repeat("=", 80))
	builder.WriteString("\n")

	if doc.Title != "" {
		builder.WriteString(doc.Title)
		builder.WriteString("\n")
	}

	builder.WriteString("ID: ")
	builder.WriteString(doc.ID)
	builder.WriteString("\n")

	builder.WriteString("Created: ")
	builder.WriteString(doc.CreatedAt)
	builder.WriteString("\n")

	builder.WriteString("Updated: ")
	builder.WriteString(doc.UpdatedAt)
	builder.WriteString("\n")

	if len(doc.Tags) > 0 {
		builder.WriteString("Tags: ")
		builder.WriteString(strings.Join(doc.Tags, ", "))
		builder.WriteString("\n")
	}

	builder.WriteString(strings.Repeat("=", 80))
	builder.WriteString("\n\n")

	// Content
	builder.WriteString(content)
	if !strings.HasSuffix(content, "\n") {
		builder.WriteString("\n")
	}

	return builder.String()
}

// stripHTML removes HTML tags from a string to produce plain text.
func stripHTML(html string) string {
	// Remove HTML tags
	text := htmlTagRegex.ReplaceAllString(html, "")

	// Decode common HTML entities
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")
	text = strings.ReplaceAll(text, "&nbsp;", " ")

	// Clean up whitespace
	text = strings.TrimSpace(text)

	return text
}

// generateFilename creates a filename from the title or ID.
func generateFilename(title, id string) string {
	if title != "" {
		return title
	}
	return id
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
func shouldUpdateFile(doc api.Document, filePath string, fs afero.Fs) bool {
	fileInfo, err := fs.Stat(filePath)
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

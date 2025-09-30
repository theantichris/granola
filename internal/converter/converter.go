package converter

import (
	"fmt"
	"strings"

	"github.com/theantichris/granola/internal/api"
	"github.com/theantichris/granola/internal/prosemirror"
	"gopkg.in/yaml.v3"
)

// Metadata represents the YAML frontmatter for a Markdown file.
type Metadata struct {
	ID        string   `yaml:"id"`
	CreatedAt string   `yaml:"created"`
	UpdatedAt string   `yaml:"updated"`
	Tags      []string `yaml:"tags,omitempty"`
}

// ToMarkdown converts a Document to Markdown format with YAML frontmatter.
// It extracts content from the ProseMirror document structure if available.
func ToMarkdown(doc api.Document) (string, error) {
	metadata := Metadata{
		ID:        doc.ID,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		Tags:      doc.Tags,
	}

	yamlBytes, err := yaml.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %w", err)
	}

	var builder strings.Builder

	// Write YAML frontmatter
	builder.WriteString("---\n")
	builder.Write(yamlBytes)
	builder.WriteString("---\n\n")

	// Write title as heading
	if doc.Title != "" {
		builder.WriteString("# ")
		builder.WriteString(doc.Title)
		builder.WriteString("\n\n")
	}

	// Write content from ProseMirror if available, otherwise fall back to plain content field
	var content string
	if doc.LastViewedPanel != nil && doc.LastViewedPanel.Content != nil {
		content = prosemirror.ConvertToMarkdown(doc.LastViewedPanel.Content)
	} else if doc.Content != "" {
		content = doc.Content
	}

	if content != "" {
		builder.WriteString(content)
		if !strings.HasSuffix(content, "\n") {
			builder.WriteString("\n")
		}
	}

	return builder.String(), nil
}
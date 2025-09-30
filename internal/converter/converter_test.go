package converter

import (
	"strings"
	"testing"

	"github.com/theantichris/granola/internal/api"
)

func TestToMarkdown(t *testing.T) {
	t.Run("converts document to markdown with frontmatter", func(t *testing.T) {
		t.Parallel()

		doc := api.Document{
			ID:        "test-id-123",
			Title:     "Test Meeting",
			Content:   "This is the meeting content.\n\nWith multiple paragraphs.",
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-02T00:00:00Z",
			Tags:      []string{"work", "planning"},
		}

		result, err := ToMarkdown(doc)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		t.Logf("Generated markdown:\n%s", result)

		if !strings.Contains(result, "---") {
			t.Error("expected markdown to contain frontmatter delimiters")
		}

		if !strings.Contains(result, "id: test-id-123") {
			t.Error("expected markdown to contain document ID")
		}

		if !strings.Contains(result, "created:") || !strings.Contains(result, "2024-01-01T00:00:00Z") {
			t.Error("expected markdown to contain created timestamp")
		}

		if !strings.Contains(result, "updated:") || !strings.Contains(result, "2024-01-02T00:00:00Z") {
			t.Error("expected markdown to contain updated timestamp")
		}

		if !strings.Contains(result, "- work") || !strings.Contains(result, "- planning") {
			t.Error("expected markdown to contain tags")
		}

		if !strings.Contains(result, "# Test Meeting") {
			t.Error("expected markdown to contain title as heading")
		}

		if !strings.Contains(result, "This is the meeting content.") {
			t.Error("expected markdown to contain content")
		}
	})

	t.Run("handles document with empty content", func(t *testing.T) {
		t.Parallel()

		doc := api.Document{
			ID:        "test-id-456",
			Title:     "Empty Note",
			Content:   "",
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
			Tags:      []string{},
		}

		result, err := ToMarkdown(doc)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !strings.Contains(result, "# Empty Note") {
			t.Error("expected markdown to contain title")
		}

		if !strings.Contains(result, "id: test-id-456") {
			t.Error("expected markdown to contain document ID")
		}
	})

	t.Run("handles document with no tags", func(t *testing.T) {
		t.Parallel()

		doc := api.Document{
			ID:        "test-id-789",
			Title:     "No Tags Note",
			Content:   "Some content",
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
			Tags:      nil,
		}

		result, err := ToMarkdown(doc)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !strings.Contains(result, "id: test-id-789") {
			t.Error("expected markdown to contain document ID")
		}

		if !strings.Contains(result, "# No Tags Note") {
			t.Error("expected markdown to contain title")
		}
	})
}
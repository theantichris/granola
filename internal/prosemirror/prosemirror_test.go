package prosemirror

import (
	"strings"
	"testing"

	"github.com/theantichris/granola/internal/api"
)

func TestConvertToMarkdown(t *testing.T) {
	t.Run("converts ProseMirror doc with heading and paragraph", func(t *testing.T) {
		t.Parallel()

		doc := &api.ProseMirrorDoc{
			Type: "doc",
			Content: []api.ProseMirrorNode{
				{
					Type: "heading",
					Attrs: map[string]interface{}{
						"level": float64(1),
					},
					Content: []api.ProseMirrorNode{
						{Type: "text", Text: "Meeting Notes"},
					},
				},
				{
					Type: "paragraph",
					Content: []api.ProseMirrorNode{
						{Type: "text", Text: "This is a paragraph."},
					},
				},
			},
		}

		result := ConvertToMarkdown(doc)

		if !strings.Contains(result, "# Meeting Notes") {
			t.Error("expected markdown to contain heading")
		}

		if !strings.Contains(result, "This is a paragraph.") {
			t.Error("expected markdown to contain paragraph")
		}
	})

	t.Run("converts bullet lists", func(t *testing.T) {
		t.Parallel()

		doc := &api.ProseMirrorDoc{
			Type: "doc",
			Content: []api.ProseMirrorNode{
				{
					Type: "bulletList",
					Content: []api.ProseMirrorNode{
						{
							Type: "listItem",
							Content: []api.ProseMirrorNode{
								{
									Type: "paragraph",
									Content: []api.ProseMirrorNode{
										{Type: "text", Text: "First item"},
									},
								},
							},
						},
						{
							Type: "listItem",
							Content: []api.ProseMirrorNode{
								{
									Type: "paragraph",
									Content: []api.ProseMirrorNode{
										{Type: "text", Text: "Second item"},
									},
								},
							},
						},
					},
				},
			},
		}

		result := ConvertToMarkdown(doc)

		if !strings.Contains(result, "- First item") {
			t.Error("expected markdown to contain first bullet item")
		}

		if !strings.Contains(result, "- Second item") {
			t.Error("expected markdown to contain second bullet item")
		}
	})

	t.Run("handles nil or empty doc", func(t *testing.T) {
		t.Parallel()

		result := ConvertToMarkdown(nil)
		if result != "" {
			t.Errorf("expected empty string for nil doc, got %q", result)
		}

		emptyDoc := &api.ProseMirrorDoc{
			Type:    "doc",
			Content: nil,
		}
		result = ConvertToMarkdown(emptyDoc)
		if result != "" {
			t.Errorf("expected empty string for empty doc, got %q", result)
		}
	})

	t.Run("converts nested bullet lists", func(t *testing.T) {
		t.Parallel()

		doc := &api.ProseMirrorDoc{
			Type: "doc",
			Content: []api.ProseMirrorNode{
				{
					Type: "bulletList",
					Content: []api.ProseMirrorNode{
						{
							Type: "listItem",
							Content: []api.ProseMirrorNode{
								{
									Type: "paragraph",
									Content: []api.ProseMirrorNode{
										{Type: "text", Text: "Parent item"},
									},
								},
								{
									Type: "bulletList",
									Content: []api.ProseMirrorNode{
										{
											Type: "listItem",
											Content: []api.ProseMirrorNode{
												{
													Type: "paragraph",
													Content: []api.ProseMirrorNode{
														{Type: "text", Text: "Nested item"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		result := ConvertToMarkdown(doc)

		if !strings.Contains(result, "- Parent item") {
			t.Error("expected markdown to contain parent bullet item")
		}

		if !strings.Contains(result, "\t- Nested item") {
			t.Error("expected markdown to contain indented nested bullet item")
		}
	})
}
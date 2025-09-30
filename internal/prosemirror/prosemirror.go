package prosemirror

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/theantichris/granola/internal/api"
)

// ConvertToMarkdown converts a ProseMirror document to Markdown format.
func ConvertToMarkdown(doc *api.ProseMirrorDoc) string {
	if doc == nil || doc.Type != "doc" || doc.Content == nil {
		return ""
	}

	var output []string
	for _, node := range doc.Content {
		output = append(output, processNode(node, 0, true))
	}

	result := strings.Join(output, "")
	// Replace multiple consecutive newlines with double newlines
	re := regexp.MustCompile(`\n{3,}`)
	result = re.ReplaceAllString(result, "\n\n")
	return strings.TrimSpace(result) + "\n"
}

// processNode recursively processes a ProseMirror node and converts it to Markdown.
func processNode(node api.ProseMirrorNode, indentLevel int, isTopLevel bool) string {
	var textContent string

	if node.Content != nil && len(node.Content) > 0 {
		switch node.Type {
		case "bulletList":
			var items []string
			for _, child := range node.Content {
				items = append(items, processNode(child, indentLevel, false))
			}
			textContent = strings.Join(items, "")
		case "listItem":
			var childContents []string
			for _, child := range node.Content {
				if child.Type == "bulletList" {
					childContents = append(childContents, processNode(child, indentLevel+1, false))
				} else {
					childContents = append(childContents, processNode(child, indentLevel, false))
				}
			}
			textContent = strings.Join(childContents, "")
		default:
			var contents []string
			for _, child := range node.Content {
				contents = append(contents, processNode(child, indentLevel, false))
			}
			textContent = strings.Join(contents, "")
		}
	} else if node.Text != "" {
		textContent = node.Text
	}

	switch node.Type {
	case "heading":
		level := 1
		if node.Attrs != nil {
			if lvl, ok := node.Attrs["level"].(float64); ok {
				level = int(lvl)
			}
		}
		suffix := "\n\n"
		if !isTopLevel {
			suffix = "\n"
		}
		return strings.Repeat("#", level) + " " + strings.TrimSpace(textContent) + suffix

	case "paragraph":
		suffix := ""
		if isTopLevel {
			suffix = "\n\n"
		}
		return textContent + suffix

	case "bulletList":
		if node.Content == nil {
			return ""
		}
		var items []string
		for _, itemNode := range node.Content {
			if itemNode.Type == "listItem" {
				var childContents []string
				var nestedLists []string

				for _, child := range itemNode.Content {
					if child.Type == "bulletList" {
						nestedLists = append(nestedLists, "\n"+processNode(child, indentLevel+1, false))
					} else {
						childContents = append(childContents, processNode(child, indentLevel, false))
					}
				}

				// Find the first non-bulletList content as the main item text
				firstText := ""
				for _, c := range childContents {
					if !strings.HasPrefix(c, "\n") {
						firstText = c
						break
					}
				}

				indent := strings.Repeat("\t", indentLevel)
				rest := strings.Join(nestedLists, "")
				items = append(items, fmt.Sprintf("%s- %s%s", indent, strings.TrimSpace(firstText), rest))
			}
		}

		suffix := ""
		if isTopLevel {
			suffix = "\n\n"
		}
		return strings.Join(items, "\n") + suffix

	case "text":
		return node.Text

	default:
		return textContent
	}
}
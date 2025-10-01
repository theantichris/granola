// Package transcript provides functionality for formatting and writing transcript files.
package transcript

import (
	"fmt"
	"strings"
	"time"

	"github.com/theantichris/granola/internal/cache"
)

// FormatTranscript formats transcript segments into a readable text format.
func FormatTranscript(doc cache.Document, segments []cache.TranscriptSegment) string {
	if len(segments) == 0 {
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

	if doc.CreatedAt != "" {
		builder.WriteString("Created: ")
		builder.WriteString(doc.CreatedAt)
		builder.WriteString("\n")
	}

	if doc.UpdatedAt != "" {
		builder.WriteString("Updated: ")
		builder.WriteString(doc.UpdatedAt)
		builder.WriteString("\n")
	}

	builder.WriteString("Segments: ")
	builder.WriteString(fmt.Sprintf("%d", len(segments)))
	builder.WriteString("\n")

	builder.WriteString(strings.Repeat("=", 80))
	builder.WriteString("\n\n")

	// Transcript segments
	for _, segment := range segments {
		// Parse timestamp
		startTime := parseTimestamp(segment.StartTimestamp)

		// Format speaker
		speaker := "System"
		if segment.Source == "microphone" {
			speaker = "You"
		}

		// Format: [HH:MM:SS] Speaker: Text
		builder.WriteString(fmt.Sprintf("[%s] %s: %s\n",
			startTime,
			speaker,
			segment.Text))
	}

	return builder.String()
}

// parseTimestamp converts ISO 8601 timestamp to HH:MM:SS format.
func parseTimestamp(timestamp string) string {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return timestamp // Return as-is if parsing fails
	}

	return t.Format("15:04:05")
}

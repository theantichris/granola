package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	// userAgent is the User-Agent header value sent with API requests.
	userAgent = "Granola/5.354.0"
	// xClientVersion is the X-Client-Version header value sent with API requests.
	xClientVersion = "5.354.0"
)

var (
	ErrDocumentAPI  = errors.New("failed to get documents")
	ErrDocumentJSON = errors.New("failed to unmarshal document JSON")
	ErrResponseBody = errors.New("failed to read response body")
	ErrHTTPRequest  = errors.New("failed to create HTTP request")
)

// GranolaResponse contains the documents retrieved from Granola.
type GranolaResponse struct {
	Documents []Document `json:"docs"`
}

// ProseMirrorNode represents a node in the ProseMirror document structure.
type ProseMirrorNode struct {
	Type    string                 `json:"type"`
	Content []ProseMirrorNode      `json:"content,omitempty"`
	Text    string                 `json:"text,omitempty"`
	Attrs   map[string]interface{} `json:"attrs,omitempty"`
}

// ProseMirrorDoc represents the ProseMirror document structure.
type ProseMirrorDoc struct {
	Type    string            `json:"type"`
	Content []ProseMirrorNode `json:"content,omitempty"`
}

// LastViewedPanel contains the ProseMirror content and metadata.
type LastViewedPanel struct {
	DocumentID        string          `json:"document_id"`
	ID                string          `json:"id"`
	CreatedAt         string          `json:"created_at"`
	Title             string          `json:"title"`
	Content           *ProseMirrorDoc `json:"-"` // Handled by custom unmarshaler
	DeletedAt         *string         `json:"deleted_at"`
	TemplateSlug      string          `json:"template_slug"`
	LastViewedAt      string          `json:"last_viewed_at"`
	UpdatedAt         string          `json:"updated_at"`
	ContentUpdatedAt  string          `json:"content_updated_at"`
	AffinityNoteID    *string         `json:"affinity_note_id"`
	OriginalContent   string          `json:"original_content"`   // HTML content
	SuggestedQuestions interface{}    `json:"suggested_questions"` // Can be null or array
	GeneratedLines    []interface{}   `json:"generated_lines"`     // Array of objects
}

// UnmarshalJSON implements custom JSON unmarshaling for LastViewedPanel.
// The API returns content as either a JSON object or a JSON string that needs to be parsed.
func (lvp *LastViewedPanel) UnmarshalJSON(data []byte) error {
	// First, unmarshal all fields except content
	aux := &struct {
		DocumentID         string        `json:"document_id"`
		ID                 string        `json:"id"`
		CreatedAt          string        `json:"created_at"`
		Title              string        `json:"title"`
		ContentRaw         json.RawMessage `json:"content"`
		DeletedAt          *string       `json:"deleted_at"`
		TemplateSlug       string        `json:"template_slug"`
		LastViewedAt       string        `json:"last_viewed_at"`
		UpdatedAt          string        `json:"updated_at"`
		ContentUpdatedAt   string        `json:"content_updated_at"`
		AffinityNoteID     *string       `json:"affinity_note_id"`
		OriginalContent    string        `json:"original_content"`
		SuggestedQuestions interface{}   `json:"suggested_questions"`
		GeneratedLines     []interface{} `json:"generated_lines"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("LastViewedPanel unmarshal aux failed: %w", err)
	}

	// Copy all the simple fields
	lvp.DocumentID = aux.DocumentID
	lvp.ID = aux.ID
	lvp.CreatedAt = aux.CreatedAt
	lvp.Title = aux.Title
	lvp.DeletedAt = aux.DeletedAt
	lvp.TemplateSlug = aux.TemplateSlug
	lvp.LastViewedAt = aux.LastViewedAt
	lvp.UpdatedAt = aux.UpdatedAt
	lvp.ContentUpdatedAt = aux.ContentUpdatedAt
	lvp.AffinityNoteID = aux.AffinityNoteID
	lvp.OriginalContent = aux.OriginalContent
	lvp.SuggestedQuestions = aux.SuggestedQuestions
	lvp.GeneratedLines = aux.GeneratedLines

	// Handle the content field if present
	if len(aux.ContentRaw) > 0 && string(aux.ContentRaw) != "null" {
		// Check if it starts with a quote (meaning it's a JSON string)
		if aux.ContentRaw[0] == '"' {
			// Unmarshal the JSON string first
			var contentStr string
			if err := json.Unmarshal(aux.ContentRaw, &contentStr); err != nil {
				return err
			}

			// Check if it's HTML (starts with '<') - if so, skip it as it's not ProseMirrorDoc
			if len(contentStr) > 0 && contentStr[0] != '<' {
				// Then unmarshal the string content into ProseMirrorDoc
				var doc ProseMirrorDoc
				if err := json.Unmarshal([]byte(contentStr), &doc); err != nil {
					return fmt.Errorf("LastViewedPanel unmarshal content string to ProseMirrorDoc failed: %w", err)
				}
				lvp.Content = &doc
			}
			// If it's HTML, skip parsing as ProseMirrorDoc (HTML is in OriginalContent field)
		} else {
			// It's already a JSON object, unmarshal directly
			var doc ProseMirrorDoc
			if err := json.Unmarshal(aux.ContentRaw, &doc); err != nil {
				return fmt.Errorf("LastViewedPanel unmarshal content object to ProseMirrorDoc failed: %w", err)
			}
			lvp.Content = &doc
		}
	}

	return nil
}

// Document contains the meeting documents from Granola.
type Document struct {
	ID              string           `json:"id"`
	Title           string           `json:"title"`
	Content         string           `json:"content"`
	CreatedAt       string           `json:"created_at"`
	UpdatedAt       string           `json:"updated_at"`
	Tags            []string         `json:"tags"`
	LastViewedPanel *LastViewedPanel `json:"last_viewed_panel,omitempty"`
	Notes           *ProseMirrorDoc  `json:"-"`      // New API structure - handled by custom unmarshaler
	NotesPlain      string           `json:"notes_plain"` // Plain text version of notes
}

// UnmarshalJSON implements custom JSON unmarshaling for Document.
// The API may return notes as either a JSON object or a JSON string that needs to be parsed.
func (d *Document) UnmarshalJSON(data []byte) error {
	// Define a struct with the same fields but without the custom unmarshaler
	aux := &struct {
		ID              string           `json:"id"`
		Title           string           `json:"title"`
		Content         string           `json:"content"`
		CreatedAt       string           `json:"created_at"`
		UpdatedAt       string           `json:"updated_at"`
		Tags            []string         `json:"tags"`
		LastViewedPanel *LastViewedPanel `json:"last_viewed_panel,omitempty"`
		NotesRaw        json.RawMessage  `json:"notes"`
		NotesPlain      string           `json:"notes_plain"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("Document unmarshal aux failed: %w", err)
	}

	// Copy the basic fields
	d.ID = aux.ID
	d.Title = aux.Title
	d.Content = aux.Content
	d.CreatedAt = aux.CreatedAt
	d.UpdatedAt = aux.UpdatedAt
	d.Tags = aux.Tags
	d.LastViewedPanel = aux.LastViewedPanel
	d.NotesPlain = aux.NotesPlain

	// Handle the notes field if present
	if len(aux.NotesRaw) > 0 && string(aux.NotesRaw) != "null" {
		// Check if it starts with a quote (meaning it's a JSON string)
		if aux.NotesRaw[0] == '"' {
			// Unmarshal the JSON string first
			var notesStr string
			if err := json.Unmarshal(aux.NotesRaw, &notesStr); err != nil {
				return err
			}

			// Then unmarshal the string content into ProseMirrorDoc
			var doc ProseMirrorDoc
			if err := json.Unmarshal([]byte(notesStr), &doc); err != nil {
				return err
			}
			d.Notes = &doc
		} else {
			// It's already a JSON object, unmarshal directly
			var doc ProseMirrorDoc
			if err := json.Unmarshal(aux.NotesRaw, &doc); err != nil {
				return err
			}
			d.Notes = &doc
		}
	}

	return nil
}

// GetDocuments gets the response from the Granola API and returns a slice of Documents.
// It automatically handles pagination to fetch all documents.
func GetDocuments(url string, file []byte, httpClient *http.Client) ([]Document, error) {
	accessToken, err := getAccessToken(file)
	if err != nil {
		return []Document{}, err
	}

	var allDocuments []Document
	offset := 0
	limit := 100

	for {
		requestBody := map[string]interface{}{
			"limit":                     limit,
			"offset":                    offset,
			"include_last_viewed_panel": true,
		}
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return []Document{}, fmt.Errorf("%w: %s", ErrHTTPRequest, err)
		}

		httpRequest, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(bodyBytes)))
		if err != nil {
			return []Document{}, fmt.Errorf("%w: %s", ErrHTTPRequest, err)
		}

		httpRequest.Header.Set("Authorization", "Bearer "+accessToken)
		httpRequest.Header.Set("Accept", "*/*")
		httpRequest.Header.Set("User-Agent", userAgent)
		httpRequest.Header.Set("X-Client-Version", xClientVersion)
		httpRequest.Header.Set("Content-Type", "application/json")

		response, err := httpClient.Do(httpRequest)
		if err != nil {
			return []Document{}, fmt.Errorf("%w: %s", ErrDocumentAPI, err)
		}

		if response.StatusCode/100 != 2 {
			// Read body for error details
			body, _ := io.ReadAll(response.Body)
			_ = response.Body.Close()
			preview := string(body)
			if len(preview) > 200 {
				preview = preview[:200] + "..."
			}
			return []Document{}, fmt.Errorf("%w: status=%s, body=%s", ErrDocumentAPI, response.Status, preview)
		}

		responseBody, err := io.ReadAll(response.Body)
		_ = response.Body.Close()
		if err != nil {
			return []Document{}, fmt.Errorf("%w: %s", ErrResponseBody, err)
		}

		var granolaResponse GranolaResponse
		if err = json.Unmarshal(responseBody, &granolaResponse); err != nil {
			// Try to parse as generic JSON to find where the error is
			var raw interface{}
			if jsonErr := json.Unmarshal(responseBody, &raw); jsonErr != nil {
				return []Document{}, fmt.Errorf("%w: raw JSON parse failed: %s", ErrDocumentJSON, jsonErr)
			}

			return []Document{}, fmt.Errorf("%w: %s", ErrDocumentJSON, err)
		}

		// Add documents from this page to the result
		allDocuments = append(allDocuments, granolaResponse.Documents...)

		// If we got fewer documents than the limit, we've reached the end
		if len(granolaResponse.Documents) < limit {
			break
		}

		// Move to the next page
		offset += limit
	}

	return allDocuments, nil
}

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
	userAgent      = "Granola/5.354.0"
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

// LastViewedPanel contains the ProseMirror content.
type LastViewedPanel struct {
	Content *ProseMirrorDoc `json:"content,omitempty"`
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
			_ = response.Body.Close()
			return []Document{}, fmt.Errorf("%w: status=%s", ErrDocumentAPI, response.Status)
		}

		responseBody, err := io.ReadAll(response.Body)
		_ = response.Body.Close()
		if err != nil {
			return []Document{}, fmt.Errorf("%w: %s", ErrResponseBody, err)
		}

		var granolaResponse GranolaResponse
		if err = json.Unmarshal(responseBody, &granolaResponse); err != nil {
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

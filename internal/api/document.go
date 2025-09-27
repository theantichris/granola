package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

// Document contains the meeting documents from Granola.
type Document struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// GetDocuments gets the respons from the Granola API and returns a slice of Documents.
func GetDocuments(url string, file []byte, httpClient *http.Client) ([]Document, error) {
	// TODO: Get access token, check for err.
	accessToken, err := getAccessToken(file)
	if err != nil {
		return []Document{}, err
	}

	// TODO: Create HTTP request.
	httpRequest, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []Document{}, fmt.Errorf("%w: %s", ErrHTTPRequest, err)
	}

	// TODO: Set headers.
	httpRequest.Header.Set("Authorization", "Bearer "+accessToken)
	httpRequest.Header.Set("Accept", "*/*")
	httpRequest.Header.Set("User-Agent", userAgent)
	httpRequest.Header.Set("X-Client-Version", xClientVersion)
	httpRequest.Header.Set("Content-Type", "application/json")

	// TODO: Refactor to use request.
	response, err := httpClient.Do(httpRequest)
	if err != nil {
		return []Document{}, fmt.Errorf("%w: %s", ErrDocumentAPI, err)
	}

	// TODO: Check status code

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return []Document{}, fmt.Errorf("%w: %s", ErrResponseBody, err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	var granolaResponse GranolaResponse
	if err = json.Unmarshal(responseBody, &granolaResponse); err != nil {
		return []Document{}, fmt.Errorf("%w: %s", ErrDocumentJSON, err)
	}

	return granolaResponse.Documents, nil
}

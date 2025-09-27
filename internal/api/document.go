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

// GetDocuments gets the response from the Granola API and returns a slice of Documents.
func GetDocuments(url string, file []byte, httpClient *http.Client) ([]Document, error) {
	accessToken, err := getAccessToken(file)
	if err != nil {
		return []Document{}, err
	}

	httpRequest, err := http.NewRequest(http.MethodGet, url, nil)
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

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode/100 != 2 {
		return []Document{}, fmt.Errorf("%w: status=%s", ErrDocumentAPI, response.Status)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return []Document{}, fmt.Errorf("%w: %s", ErrResponseBody, err)
	}

	var granolaResponse GranolaResponse
	if err = json.Unmarshal(responseBody, &granolaResponse); err != nil {
		return []Document{}, fmt.Errorf("%w: %s", ErrDocumentJSON, err)
	}

	return granolaResponse.Documents, nil
}

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
)

const (
	userAgent      = "Granola/5.354.0"
	xClientVersion = "5.354.0"
)

var (
	ErrDocumentAPI  = errors.New("failed to get documents")
	ErrDocumentJSON = errors.New("failed to unmarshal document JSON")
	ErrResponseBody = errors.New("failed to read response body")
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
func GetDocuments(url string, httpClient *http.Client, logger *log.Logger) ([]Document, error) {
	// TODO: Get access token, check for err.
	// TODO: Create HTTP request.
	// TODO: Set headers.

	// TODO: Refactor to use request.
	response, err := httpClient.Get(url)
	if err != nil {
		logger.Error(ErrDocumentAPI.Error(), "error", err)

		return []Document{}, fmt.Errorf("%w: %s", ErrDocumentAPI, err)
	}

	// TODO: Check status code

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error(ErrResponseBody.Error(), "error", err)

		return []Document{}, fmt.Errorf("%w: %s", ErrResponseBody, err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	var granolaResponse GranolaResponse
	if err = json.Unmarshal(responseBody, &granolaResponse); err != nil {
		logger.Error(ErrDocumentJSON, "error", err)

		return []Document{}, fmt.Errorf("%w: %s", ErrDocumentJSON, err)
	}

	return granolaResponse.Documents, nil
}

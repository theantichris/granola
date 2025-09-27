package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/google/go-cmp/cmp"
)

type errorTransport struct{}

func (e *errorTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("forced transport error")
}

func TestGetDocuments(t *testing.T) {
	logger := log.New(io.Discard)

	t.Run("gets the Granola documents", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{\"docs\":[{\"id\":\"abc123\",\"title\":\"Test Meeting\"}]}"))
		}))
		defer testServer.Close()

		httpClient := &http.Client{Transport: testServer.Client().Transport}

		actual, err := GetDocuments(testServer.URL, httpClient, logger)
		if err != nil {
			t.Fatalf("expected no error getting documents, got %v", err)
		}

		expected := []Document{
			{ID: "abc123", Title: "Test Meeting"},
		}

		if !cmp.Equal(actual, expected) {
			t.Errorf("expected response %v, got %v", expected, actual)
		}
	})

	t.Run("returns error for bad HTTP request", func(t *testing.T) {
		t.Parallel()

		httpClient := &http.Client{Transport: &errorTransport{}}

		_, err := GetDocuments("http://test.dev", httpClient, logger)
		if err == nil {
			t.Fatal("expected error getting documents, got nil")
		}

		if !errors.Is(err, ErrDocumentAPI) {
			t.Errorf("expected error %v, got %v", ErrDocumentAPI, err)
		}
	})

	t.Run("returns error for bad JSON", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("`invalid JSON`"))
		}))
		defer testServer.Close()

		httpClient := &http.Client{Transport: testServer.Client().Transport}

		_, err := GetDocuments(testServer.URL, httpClient, logger)
		if err == nil {
			t.Fatal("expected error getting documents, got nil")
		}

		if !errors.Is(err, ErrDocumentJSON) {
			t.Errorf("expected error %v, got %v", ErrDocumentJSON, err)
		}
	})
}

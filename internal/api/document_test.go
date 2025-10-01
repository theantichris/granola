package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type errorTransport struct{}

func (e *errorTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("forced transport error")
}

func TestGetDocuments(t *testing.T) {
	t.Run("gets the Granola documents", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify it's a POST request
			if r.Method != http.MethodPost {
				t.Errorf("expected POST request, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{\"docs\":[{\"id\":\"abc123\",\"title\":\"Test Meeting\",\"content\":\"Meeting notes\",\"created_at\":\"2024-01-01T00:00:00Z\",\"updated_at\":\"2024-01-02T00:00:00Z\",\"tags\":[\"work\",\"planning\"]}]}"))
		}))
		defer testServer.Close()

		httpClient := &http.Client{Transport: testServer.Client().Transport}

		actual, err := GetDocuments(testServer.URL, []byte(accessTokenJSON), httpClient)
		if err != nil {
			t.Fatalf("expected no error getting documents, got %v", err)
		}

		expected := []Document{
			{
				ID:              "abc123",
				Title:           "Test Meeting",
				Content:         "Meeting notes",
				CreatedAt:       "2024-01-01T00:00:00Z",
				UpdatedAt:       "2024-01-02T00:00:00Z",
				Tags:            []string{"work", "planning"},
				LastViewedPanel: nil,
			},
		}

		if !cmp.Equal(actual, expected) {
			t.Errorf("expected response %v, got %v", expected, actual)
		}
	})

	t.Run("returns error for bad HTTP request", func(t *testing.T) {
		t.Parallel()

		httpClient := &http.Client{Transport: &errorTransport{}}

		_, err := GetDocuments("http://test.dev", []byte(accessTokenJSON), httpClient)
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

		_, err := GetDocuments(testServer.URL, []byte(badTokenJSON), httpClient)
		if err == nil {
			t.Fatal("expected error getting documents, got nil")
		}

		if !errors.Is(err, ErrTokensJSON) {
			t.Errorf("expected error %v, got %v", ErrDocumentJSON, err)
		}
	})

	t.Run("returns error for HTTP failure", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		testServer.Close()

		_, err := GetDocuments(testServer.URL, []byte(accessTokenJSON), http.DefaultClient)
		if err == nil {
			t.Fatal("expected error getting documents, got nil")
		}

		if !errors.Is(err, ErrDocumentAPI) {
			t.Errorf("expected error %v, got %v", ErrDocumentAPI, err)
		}
	})

	t.Run("returns error for non-2xx status", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer testServer.Close()

		httpClient := &http.Client{Transport: testServer.Client().Transport}

		_, err := GetDocuments(testServer.URL, []byte(accessTokenJSON), httpClient)
		if err == nil {
			t.Fatal("expected error getting documents, got nil")
		}

		if !errors.Is(err, ErrDocumentAPI) {
			t.Errorf("expected error %v, got %v", ErrDocumentAPI, err)
		}

		if !strings.Contains(err.Error(), "401 Unauthorized") {
			t.Errorf("expected error containing %q, got %q", "401 Unauthorized", err.Error())
		}
	})

	t.Run("handles last_viewed_panel content as JSON string", func(t *testing.T) {
		t.Parallel()

		// Simulate API returning content as a JSON string (escaped)
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"docs":[{"id":"doc1","title":"Test","content":"text","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z","tags":[],"last_viewed_panel":{"content":"{\"type\":\"doc\",\"content\":[{\"type\":\"paragraph\",\"content\":[{\"type\":\"text\",\"text\":\"Hello\"}]}]}"}}]}`))
		}))
		defer testServer.Close()

		httpClient := &http.Client{Transport: testServer.Client().Transport}

		docs, err := GetDocuments(testServer.URL, []byte(accessTokenJSON), httpClient)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(docs) != 1 {
			t.Fatalf("expected 1 document, got %d", len(docs))
		}

		if docs[0].LastViewedPanel == nil {
			t.Fatal("expected last_viewed_panel to be present")
		}

		if docs[0].LastViewedPanel.Content == nil {
			t.Fatal("expected content to be present")
		}

		if docs[0].LastViewedPanel.Content.Type != "doc" {
			t.Errorf("expected type 'doc', got %q", docs[0].LastViewedPanel.Content.Type)
		}
	})

	t.Run("handles last_viewed_panel content as JSON object", func(t *testing.T) {
		t.Parallel()

		// Simulate API returning content as a JSON object (not escaped)
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"docs":[{"id":"doc2","title":"Test","content":"text","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z","tags":[],"last_viewed_panel":{"content":{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"World"}]}]}}}]}`))
		}))
		defer testServer.Close()

		httpClient := &http.Client{Transport: testServer.Client().Transport}

		docs, err := GetDocuments(testServer.URL, []byte(accessTokenJSON), httpClient)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(docs) != 1 {
			t.Fatalf("expected 1 document, got %d", len(docs))
		}

		if docs[0].LastViewedPanel == nil {
			t.Fatal("expected last_viewed_panel to be present")
		}

		if docs[0].LastViewedPanel.Content == nil {
			t.Fatal("expected content to be present")
		}

		if docs[0].LastViewedPanel.Content.Type != "doc" {
			t.Errorf("expected type 'doc', got %q", docs[0].LastViewedPanel.Content.Type)
		}
	})

	t.Run("handles notes field as JSON object", func(t *testing.T) {
		t.Parallel()

		// Simulate new API returning notes as a JSON object
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"docs":[{"id":"doc3","title":"Test","content":"text","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z","tags":[],"notes":{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Notes content"}]}]}}]}`))
		}))
		defer testServer.Close()

		httpClient := &http.Client{Transport: testServer.Client().Transport}

		docs, err := GetDocuments(testServer.URL, []byte(accessTokenJSON), httpClient)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(docs) != 1 {
			t.Fatalf("expected 1 document, got %d", len(docs))
		}

		if docs[0].Notes == nil {
			t.Fatal("expected notes to be present")
		}

		if docs[0].Notes.Type != "doc" {
			t.Errorf("expected type 'doc', got %q", docs[0].Notes.Type)
		}
	})

	t.Run("handles notes field as JSON string", func(t *testing.T) {
		t.Parallel()

		// Simulate new API returning notes as a JSON string (escaped)
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"docs":[{"id":"doc4","title":"Test","content":"text","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z","tags":[],"notes":"{\"type\":\"doc\",\"content\":[{\"type\":\"paragraph\",\"content\":[{\"type\":\"text\",\"text\":\"Notes string\"}]}]}"}]}`))
		}))
		defer testServer.Close()

		httpClient := &http.Client{Transport: testServer.Client().Transport}

		docs, err := GetDocuments(testServer.URL, []byte(accessTokenJSON), httpClient)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(docs) != 1 {
			t.Fatalf("expected 1 document, got %d", len(docs))
		}

		if docs[0].Notes == nil {
			t.Fatal("expected notes to be present")
		}

		if docs[0].Notes.Type != "doc" {
			t.Errorf("expected type 'doc', got %q", docs[0].Notes.Type)
		}
	})
}

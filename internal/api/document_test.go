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
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{\"docs\":[{\"id\":\"abc123\",\"title\":\"Test Meeting\"}]}"))
		}))
		defer testServer.Close()

		httpClient := &http.Client{Transport: testServer.Client().Transport}

		actual, err := GetDocuments(testServer.URL, []byte(accessTokenJSON), httpClient)
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
}

package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/google/go-cmp/cmp"
)

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
}

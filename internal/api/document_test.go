package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetDocuments(t *testing.T) {
	t.Run("gets the Granola documents", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{\"docs\":[{\"id\":\"abc123\",\"title\":\"Test Meeting\"}]}"))
		}))
		defer testServer.Close()

		httpClient := &http.Client{Transport: testServer.Client().Transport}

		response, err := httpClient.Get(testServer.URL)
		if err != nil {
			t.Fatalf("expected no error sending HTTP request, got %v", err)
		}

		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("expected no error reading response body, got %v", err)
		}
		defer response.Body.Close()

		var actual GranolaResponse
		err = json.Unmarshal(responseBody, &actual)
		if err != nil {
			t.Fatalf("expected no error unmarshaling response, got %v", err)
		}

		expected := GranolaResponse{
			Documents: []Document{
				{ID: "abc123", Title: "Test Meeting"},
			},
		}

		if !cmp.Equal(actual, expected) {
			t.Errorf("expected response %v, got %v", expected, actual)
		}
	})
}

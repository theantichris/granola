package api

import (
	"errors"
	"testing"
)

var accessTokenJSON = `{"workos_tokens": "{\"access_token\":\"access_token_123\"}"}`
var badWrapperJSON = "{"
var badTokenJSON = `{"workos_tokens": "{","session_id": "session_id_123",  "user_info": "{}"}`
var sessionIDJSON = `{"workos_tokens": "{\"session_id\":\"session_id_123\"}"}`

func TestGetAccessToken(t *testing.T) {
	t.Run("returns the access token", func(t *testing.T) {
		t.Parallel()

		actual, err := getAccessToken([]byte(accessTokenJSON))
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}

		expected := "access_token_123"

		if actual != expected {
			t.Errorf("expected access token %q, got %q", expected, actual)
		}
	})

	t.Run("returns error for bad wrapper JSON", func(t *testing.T) {
		t.Parallel()

		_, err := getAccessToken([]byte(badWrapperJSON))
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if !errors.Is(err, ErrWrapperJSON) {
			t.Errorf("expected error %v, got %v", ErrWrapperJSON, err)
		}
	})

	t.Run("returns error for bad token JSON", func(t *testing.T) {
		t.Parallel()

		_, err := getAccessToken([]byte(badTokenJSON))
		if err == nil {
			t.Errorf("expected error, not nil")
		}

		if !errors.Is(err, ErrTokensJSON) {
			t.Errorf("expected error %v, got %v", ErrTokensJSON, err)
		}
	})

	t.Run("returns error if no access token is found", func(t *testing.T) {
		t.Parallel()

		_, err := getAccessToken([]byte(sessionIDJSON))
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrAccessTokenNotFound) {
			t.Errorf("expected error %v, got %v", ErrAccessTokenNotFound, err)
		}
	})
}

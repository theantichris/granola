package api

import (
	"errors"
	"io"
	"testing"

	"github.com/charmbracelet/log"
)

var logger *log.Logger = log.New(io.Discard)

const testJSON = `{
  "workos_tokens": "{\"external_id\":\"external_id_123\",\"access_token\":\"access_token_123\",\"expires_in\":43199,\"refresh_token\":\"refresh_token_123\",\"token_type\":\"Bearer\",\"obtained_at\":1758926490172,\"session_id\":\"session_id_123\"}",
  "session_id": "session_id_123",
  "user_info": "{\"id\":\"user_id_123\",\"email\":\"email_123@example.com\",\"user_metadata\":{\"name\":\"name_123\",\"picture\":\"picture_url_123\",\"hd\":\"domain_123.com\"},\"signed_in_on_platforms\":{\"macos\":true,\"windows\":true,\"ios\":true},\"dub_id\":null,\"signup_platform\":null,\"google_ads_click_id\":null,\"facebook_ads_click_id\":null}"
}`

func TestGetSupabase(t *testing.T) {
	t.Run("gets the supabase information", func(t *testing.T) {
		t.Parallel()

		actual, err := getTokens([]byte(testJSON), logger)
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}

		expected := Tokens{ExternalID: "external_id_123"}

		if actual.ExternalID != expected.ExternalID {
			t.Errorf("expected external ID %q, got %q", expected.ExternalID, actual.ExternalID)
		}
	})

	t.Run("returns error for bad wrapper JSON", func(t *testing.T) {
		t.Parallel()

		badWrapperJSON := "{"

		_, err := getTokens([]byte(badWrapperJSON), logger)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if !errors.Is(err, ErrWrapperJSON) {
			t.Errorf("expected error %v, got %v", ErrWrapperJSON, err)
		}
	})

	t.Run("returns error for bad token JSON", func(t *testing.T) {
		t.Parallel()

		badTokenJSON := `{"workos_tokens": "{","session_id": "session_id_123",  "user_info": "{}"}`

		_, err := getTokens([]byte(badTokenJSON), logger)
		if err == nil {
			t.Errorf("expected error, not nil")
		}

		if !errors.Is(err, ErrTokensJSON) {
			t.Errorf("expected error %v, got %v", ErrTokensJSON, err)
		}
	})
}

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
)

var (
	ErrWrapperJSON         = errors.New("failed to unmarshal wrapper JSON")
	ErrTokensJSON          = errors.New("failed to unmarshal token JSON")
	ErrAccessTokenNotFound = errors.New("access token not found")
)

// Wrapper holds the data from Granola's supabase.json file.
type Wrapper struct {
	Tokens string `json:"workos_tokens"`
}

// Tokens holds the access token and related information.
type Tokens struct {
	AccessToken string `json:"access_token"`
}

// getAccessToken takes the JSON from supabase.json and returns the Granola access token.
func getAccessToken(file []byte, logger *log.Logger) (string, error) {
	var wrapper Wrapper
	if err := json.Unmarshal(file, &wrapper); err != nil {
		logger.Error(ErrWrapperJSON.Error(), "error", err)

		return "", fmt.Errorf("%w: %s", ErrWrapperJSON, err)
	}

	var tokens Tokens
	if err := json.Unmarshal([]byte(wrapper.Tokens), &tokens); err != nil {
		logger.Error(ErrTokensJSON.Error(), "error", err)

		return "", fmt.Errorf("%w: %s", ErrTokensJSON, err)
	}

	if strings.TrimSpace(tokens.AccessToken) == "" {
		logger.Error(ErrAccessTokenNotFound.Error())

		return "", ErrAccessTokenNotFound
	}

	return tokens.AccessToken, nil
}

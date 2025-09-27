package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/charmbracelet/log"
)

var (
	ErrWrapperJSON = errors.New("couldn't unmarshal wrapper JSON")
	ErrTokensJSON  = errors.New("couldn't unmarshal token JSON")
)

// Wrapper holds the data from Granola's supabase.json file.
type Wrapper struct {
	Tokens string `json:"workos_tokens"`
}

// Tokens holds the access token and related information.
type Tokens struct {
	ExternalID string `json:"external_id"`
}

// getTokens takes the JSON data and returns the Granola token information.
func getTokens(file []byte, logger *log.Logger) (Tokens, error) {
	var wrapper Wrapper
	if err := json.Unmarshal(file, &wrapper); err != nil {
		logger.Error(ErrWrapperJSON.Error(), "error", err)

		return Tokens{}, fmt.Errorf("%w: %s", ErrWrapperJSON, err)
	}

	var tokens Tokens

	if err := json.Unmarshal([]byte(wrapper.Tokens), &tokens); err != nil {
		logger.Error("couldn't unmarshal token JSON", "error", err)

		return Tokens{}, fmt.Errorf("%w: %s", ErrTokensJSON, err)
	}

	return tokens, nil
}

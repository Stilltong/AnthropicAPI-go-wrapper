package anthropic

import (
	"context"
	"net/http"
)

type CompleteRequest struct {
	Model             string `json:"model"`
	Prompt            string `json:"prompt"`
	MaxTokensToSample int    `json:"max_tokens_to_sample"`

	StopSequences []string       `json:"
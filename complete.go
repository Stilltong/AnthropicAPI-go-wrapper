package anthropic

import (
	"context"
	"net/http"
)

type CompleteRequest struct {
	Model             string `json:"model"`
	Prompt            string `json:"prompt"`
	MaxTokensToSample int    `json:"max_tokens_to_sample"`

	StopSequences []string       `json:"stop_sequences,omitempty"`
	Temperature   *float32       `json:"temperature,omitempty"`
	TopP          *float32       `json:"top_p,
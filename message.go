package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
)

type MessagesResponseType string

const (
	MessagesResponseTypeMessage MessagesResponseType = "message"
	MessagesResponseTypeError   MessagesResponseType = "error"
)

type MessagesContentType string

const (
	MessagesContentTypeText           MessagesContentType = "text"
	MessagesContentTypeTextDelta      MessagesContentType = "text_delta"
	MessagesContentTypeImage          MessagesContentType = "image"
	MessagesContentTypeToolResult     MessagesContentType = "tool_result"
	MessagesContentTypeToolUse        MessagesContentType = "tool_use"
	MessagesContentTypeInputJsonDelta MessagesContentType = "input_json_delta"
)

type MessagesStopReason string

const (
	MessagesStopReasonEndTurn      MessagesStopReason = "end_turn"
	MessagesStopReasonMaxTokens    MessagesStopReason = "max_tokens"
	MessagesStopReasonStopSequence MessagesStopReason = "stop_sequence"
	MessagesStopReasonToolUse      MessagesStopReason = "tool_use"
)

type MessagesRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`

	System        string           `json:"system,omitempty"`
	Metadata      map[string]any   `json:"metadata,omitempty"`
	StopSequences []string         `json:"stop_sequences,omitempty"`
	Stream        bool             `json:"stream,omitempty"`
	Temperature   *float32         `json:"temperature,omitempty"`
	TopP          *float32         `json:"top_p,omitempty"`
	TopK          *int             `json:"top_k,omitempty"`
	Tools         []ToolDefinition `json:"tools,omitempty"`
	ToolChoice
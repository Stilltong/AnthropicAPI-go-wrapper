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
	ToolChoice    *ToolChoice      `json:"tool_choice,omitempty"`
}

func (m *MessagesRequest) SetTemperature(t float32) {
	m.Temperature = &t
}

func (m *MessagesRequest) SetTopP(p float32) {
	m.TopP = &p
}

func (m *MessagesRequest) SetTopK(k int) {
	m.TopK = &k
}

type Message struct {
	Role    string           `json:"role"`
	Content []MessageContent `json:"content"`
}

func NewUserTextMessage(text string) Message {
	return Message{
		Role:    RoleUser,
		Content: []MessageContent{NewTextMessageContent(text)},
	}
}

func NewAssistantTextMessage(text string) Message {
	return Message{
		Role:    RoleAssistant,
		Content: []MessageContent{NewTextMessageContent(text)},
	}
}

func NewToolResultsMessage(toolUseID, content string, isError bool) Message {
	return Message{
		Role:    RoleUser,
		Content: []MessageContent{NewToolResultMessageContent(toolUseID, content, isError)},
	}
}

func (m Message) GetFirstContent() MessageContent {
	if len(m.Content) == 0 {
		return MessageContent{}
	}
	return m.Content[0]
}

type MessageContent struct {
	Type MessagesContentType `json:"type"`

	Text *string `json:"text,omitempty"`

	Source *MessageContentImageSource `json:"source,omitempty"`

	*MessageContentToolResult

	*MessageContentToolUse

	PartialJson *string `json:"partial_json,omitempty"`
}

func NewTextMessageContent(text string) MessageContent {
	return MessageContent{
		Type: MessagesContentTypeText,
		Text: &text,
	}
}

func NewImageMessageContent(source MessageContentImageSource) MessageContent {
	return MessageContent{
		Type:   MessagesContentTypeImage,
		Source: &source,
	}
}

func NewToolResultMessageContent(toolUseID, content string, isError bool) MessageContent {
	return MessageContent{
		Type:                     MessagesContentTypeToolResult,
		MessageContentToolResult: NewMessageContentToolResult(toolUseID, content, isError),
	}
}

func NewToolUseMessageContent(toolUseID, name string, input json.RawMessage) MessageContent {
	return MessageContent{
		Type: MessagesContentTypeToolUse,
		MessageContentToolUse: &MessageContentToolUse{
			ID:    toolUseID,
			Name:  name,
			Input: input,
		},
	}
}

func (m *MessageContent) GetText() string {
	if m.Text != nil {
		return *m.Text
	}
	return ""
}

func (m *Messa
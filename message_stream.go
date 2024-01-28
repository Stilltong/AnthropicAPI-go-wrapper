
package anthropic

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
)

var (
	eventPrefix                   = []byte("event:")
	dataPrefix                    = []byte("data:")
	ErrTooManyEmptyStreamMessages = errors.New("stream has sent too many empty messages")
)

type (
	// MessagesEvent docs: https://docs.anthropic.com/claude/reference/messages-streaming
	MessagesEvent string
)

const (
	MessagesEventError             MessagesEvent = "error"
	MessagesEventMessageStart      MessagesEvent = "message_start"
	MessagesEventContentBlockStart MessagesEvent = "content_block_start"
	MessagesEventPing              MessagesEvent = "ping"
	MessagesEventContentBlockDelta MessagesEvent = "content_block_delta"
	MessagesEventContentBlockStop  MessagesEvent = "content_block_stop"
	MessagesEventMessageDelta      MessagesEvent = "message_delta"
	MessagesEventMessageStop       MessagesEvent = "message_stop"
)

type MessagesStreamRequest struct {
	MessagesRequest

	OnError             func(ErrorResponse)                                     `json:"-"`
	OnPing              func(MessagesEventPingData)                             `json:"-"`
	OnMessageStart      func(MessagesEventMessageStartData)                     `json:"-"`
	OnContentBlockStart func(MessagesEventContentBlockStartData)                `json:"-"`
	OnContentBlockDelta func(MessagesEventContentBlockDeltaData)                `json:"-"`
	OnContentBlockStop  func(MessagesEventContentBlockStopData, MessageContent) `json:"-"`
	OnMessageDelta      func(MessagesEventMessageDeltaData)                     `json:"-"`
	OnMessageStop       func(MessagesEventMessageStopData)                      `json:"-"`
}

type MessagesEventMessageStartData struct {
	Type    MessagesEvent    `json:"type"`
	Message MessagesResponse `json:"message"`
}

type MessagesEventContentBlockStartData struct {
	Type         MessagesEvent  `json:"type"`
	Index        int            `json:"index"`
	ContentBlock MessageContent `json:"content_block"`
}

type MessagesEventPingData struct {
	Type string `json:"type"`
}

type MessagesEventContentBlockDeltaData struct {
	Type  string         `json:"type"`
	Index int            `json:"index"`
	Delta MessageContent `json:"delta"`
}

type MessagesEventContentBlockStopData struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
}

type MessagesEventMessageDeltaData struct {
	Type  string           `json:"type"`
	Delta MessagesResponse `json:"delta"`
	Usage MessagesUsage    `json:"usage"`
}

type MessagesEventMessageStopData struct {
	Type string `json:"type"`
}

func (c *Client) CreateMessagesStream(ctx context.Context, request MessagesStreamRequest) (response MessagesResponse, err error) {
	request.Stream = true

	var setters []requestSetter
	if len(request.Tools) > 0 {
		setters = append(setters, withBetaVersion(c.config.BetaVersion))
	}

	urlSuffix := "/messages"
	req, err := c.newStreamRequest(ctx, http.MethodPost, urlSuffix, request, setters...)
	if err != nil {
		return
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	response.SetHeader(resp.Header)

	if err := c.handlerRequestError(resp); err != nil {
		return response, err
	}

	reader := bufio.NewReader(resp.Body)
	var (
		event             []byte
		emptyMessageCount uint
	)
	for {
		rawLine, readErr := reader.ReadBytes('\n')
		if readErr != nil {
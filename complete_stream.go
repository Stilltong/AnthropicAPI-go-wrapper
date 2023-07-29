
package anthropic

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type CompleteEvent string

const (
	CompleteEventError      CompleteEvent = "error"
	CompleteEventCompletion CompleteEvent = "completion"
	CompleteEventPing       CompleteEvent = "ping"
)

type CompleteStreamRequest struct {
	CompleteRequest

	OnCompletion func(CompleteResponse)       `json:"-"`
	OnPing       func(CompleteStreamPingData) `json:"-"`
	OnError      func(ErrorResponse)          `json:"-"`
}

type CompleteStreamPingData struct {
	Type string `json:"type"`
}

func (c *Client) CreateCompleteStream(ctx context.Context, request CompleteStreamRequest) (response CompleteResponse, err error) {
	request.Stream = true

	urlSuffix := "/complete"
	req, err := c.newStreamRequest(ctx, http.MethodPost, urlSuffix, request)
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
			if errors.Is(readErr, io.EOF) {
				break
			}
			return response, readErr
		}

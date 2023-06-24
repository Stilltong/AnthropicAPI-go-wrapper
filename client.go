package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	config ClientConfig
}

type Response interface {
	SetHeader(http.Header)
}

type httpHeader http.Header

func (h *httpHeader) SetHeader(header http.Header) {
	*h = httpHeader(header)
}

func (h *httpHeader) Header() http.Header {
	return http.Header(*h)
}

// NewClient create new Anthropic API client
func NewClient(apikey string, opts ...ClientOption) *Client {
	return &Client{
		config: newConfig(apikey, opts...),
	}
}

func (c *Client) sendRequest(req *http.Request, v Response) error {
	res, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	v.SetHeader(res.Header)

	if err := c.handlerRequestError(res); err != nil {
		return err
	}

	if err = json.NewDecoder(res.Body).Decode(v); err != nil {
		return e
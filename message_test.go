package anthropic_test

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/liushuangls/go-anthropic/v2/internal/test"
	"github.com/liushuangls/go-anthropic/v2/internal/test/checks"
	"github.com/liushuangls/go-anthropic/v2/jsonschema"
)

//go:embed internal/test/sources/*
var sources embed.FS

func TestMessages(t *testing.T) {
	server := test.NewTestServer()
	server.RegisterHandler("/v1/messages", handleMessagesEndpoint)

	ts := server.AnthropicTestServer()
	ts.Start()
	defer ts.Close()

	baseUrl := ts.URL + "/v1"
	client := anthropic.NewClient(
		test.GetTestToken(),
		anthropic.WithBaseURL(baseUrl),
		anthropic.WithAPIVersion(anthropic.APIVersion20230601),
		anthropic.WithEmptyMessagesLimit(100),
		anthropic.WithHTTPClient(http.DefaultClient),
	)
	resp, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.ModelClaudeInstant1Dot2,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage("What is your name?"),
		},
		MaxTokens: 1000,
	})
	if err != nil {
		t.Fatalf("CreateMessages error: %v", err)
	}

	t.Logf("CreateMessages resp: %+v", resp)
}

func TestMessagesTokenError(t *testing.T) {
	server := test.NewTestServer()
	server.RegisterHandler("/v1/messages", handleMessagesEndpoint)

	ts := server.AnthropicTestServer()
	ts.Start()
	defer ts.Close()

	baseUrl := ts.URL + "/v1"
	client := anthropic.NewClient(
		test.GetTestToken()+"1",
		anthropic.WithBaseURL(baseUrl),
	)
	_, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.ModelClaudeIn
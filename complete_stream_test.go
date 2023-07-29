
package anthropic_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/liushuangls/go-anthropic/v2/internal/test"
	"github.com/liushuangls/go-anthropic/v2/internal/test/checks"
)

var (
	testCompletionStreamContent = []string{"My", " name", " is", " Claude", "."}
)

func TestCompleteStream(t *testing.T) {
	server := test.NewTestServer()
	server.RegisterHandler("/v1/complete", handlerCompleteStream)

	ts := server.AnthropicTestServer()
	ts.Start()
	defer ts.Close()

	baseUrl := ts.URL + "/v1"
	client := anthropic.NewClient(
		test.GetTestToken(),
		anthropic.WithBaseURL(baseUrl),
	)
	var receivedContent string
	resp, err := client.CreateCompleteStream(context.Background(), anthropic.CompleteStreamRequest{
		CompleteRequest: anthropic.CompleteRequest{
			Model:             anthropic.ModelClaudeInstant1Dot2,
			Prompt:            "\n\nHuman: What is your name?\n\nAssistant:",
			MaxTokensToSample: 1000,
		},
		OnCompletion: func(data anthropic.CompleteResponse) {
			receivedContent += data.Completion
			//t.Logf("CreateCompleteStream OnCompletion data: %+v", data)
		},
		OnPing:  func(data anthropic.CompleteStreamPingData) {},
		OnError: func(response anthropic.ErrorResponse) {},
	})
	if err != nil {
		t.Fatalf("CreateCompleteStream error: %s", err)
	}

	expected := strings.Join(testCompletionStreamContent, "")
	if receivedContent != expected {
		t.Fatalf("CreateCompleteStream content not match expected: %s, got: %s", expected, receivedContent)
	}
	if resp.Completion != expected {
		t.Fatalf("CreateCompleteStream content not match expected: %s, got: %s", expected, resp.Completion)
	}
	t.Logf("CreateCompleteStream resp: %+v", resp)
}

func TestCompleteStreamError(t *testing.T) {
	server := test.NewTestServer()
	server.RegisterHandler("/v1/complete", handlerCompleteStream)

	ts := server.AnthropicTestServer()
	ts.Start()
	defer ts.Close()

	baseUrl := ts.URL + "/v1"
	client := anthropic.NewClient(
		test.GetTestToken(),
		anthropic.WithBaseURL(baseUrl),
	)
	var receivedContent string
	param := anthropic.CompleteStreamRequest{
		CompleteRequest: anthropic.CompleteRequest{
			Model:             anthropic.ModelClaudeInstant1Dot2,
			Prompt:            "\n\nHuman: What is your name?\n\nAssistant:",
			MaxTokensToSample: 1000,
			//Temperature:       &temperature,
		},
		OnCompletion: func(data anthropic.CompleteResponse) {
			receivedContent += data.Completion
			//t.Logf("CreateCompleteStream OnCompletion data: %+v", data)
		},
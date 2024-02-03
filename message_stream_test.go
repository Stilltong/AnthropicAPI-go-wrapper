
package anthropic_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/liushuangls/go-anthropic/v2/internal/test"
	"github.com/liushuangls/go-anthropic/v2/jsonschema"
)

var (
	testMessagesStreamContent    = []string{"My", " name", " is", " Claude", "."}
	testMessagesJsonDeltaContent = []string{`{\"location\":`, `\"San Francisco, CA\"}`}
)

func TestMessagesStream(t *testing.T) {
	server := test.NewTestServer()
	server.RegisterHandler("/v1/messages", handlerMessagesStream)

	ts := server.AnthropicTestServer()
	ts.Start()
	defer ts.Close()

	baseUrl := ts.URL + "/v1"
	client := anthropic.NewClient(
		test.GetTestToken(),
		anthropic.WithBaseURL(baseUrl),
	)
	var received string
	resp, err := client.CreateMessagesStream(context.Background(), anthropic.MessagesStreamRequest{
		MessagesRequest: anthropic.MessagesRequest{
			Model: anthropic.ModelClaudeInstant1Dot2,
			Messages: []anthropic.Message{
				anthropic.NewUserTextMessage("What is your name?"),
			},
			MaxTokens: 1000,
		},
		OnContentBlockDelta: func(data anthropic.MessagesEventContentBlockDeltaData) {
			received += data.Delta.GetText()
			//t.Logf("CreateMessagesStream delta resp: %+v", data)
		},
		OnError:             func(response anthropic.ErrorResponse) {},
		OnPing:              func(data anthropic.MessagesEventPingData) {},
		OnMessageStart:      func(data anthropic.MessagesEventMessageStartData) {},
		OnContentBlockStart: func(data anthropic.MessagesEventContentBlockStartData) {},
		OnContentBlockStop:  func(data anthropic.MessagesEventContentBlockStopData, content anthropic.MessageContent) {},
		OnMessageDelta:      func(data anthropic.MessagesEventMessageDeltaData) {},
		OnMessageStop:       func(data anthropic.MessagesEventMessageStopData) {},
	})
	if err != nil {
		t.Fatalf("CreateMessagesStream error: %s", err)
	}

	expectedContent := strings.Join(testMessagesStreamContent, "")
	if received != expectedContent {
		t.Fatalf("CreateMessagesStream content not match expected: %s, got: %s", expectedContent, received)
	}
	if resp.GetFirstContentText() != expectedContent {
		t.Fatalf("CreateMessagesStream content not match expected: %s, got: %s", expectedContent, resp.GetFirstContentText())
	}

	t.Logf("CreateMessagesStream resp: %+v", resp)
}

func TestMessagesStreamError(t *testing.T) {
	server := test.NewTestServer()
	server.RegisterHandler("/v1/messages", handlerMessagesStream)

	ts := server.AnthropicTestServer()
	ts.Start()
	defer ts.Close()

	baseUrl := ts.URL + "/v1"
	client := anthropic.NewClient(
		test.GetTestToken(),
		anthropic.WithBaseURL(baseUrl),
	)
	param := anthropic.MessagesStreamRequest{
		MessagesRequest: anthropic.MessagesRequest{
			Model: anthropic.ModelClaudeInstant1Dot2,
			Messages: []anthropic.Message{
				anthropic.NewUserTextMessage("What is your name?"),
			},
			MaxTokens: 1000,
		},
		OnContentBlockDelta: func(data anthropic.MessagesEventContentBlockDeltaData) {
			t.Logf("CreateMessagesStream delta resp: %+v", data)
		},
		OnError: func(response anthropic.ErrorResponse) {},
	}
	param.SetTemperature(2)
	param.SetTopP(2)
	param.SetTopK(1)
	_, err := client.CreateMessagesStream(context.Background(), param)
	if err == nil {
		t.Fatalf("CreateMessagesStream expect error, but not")
	}

	t.Logf("CreateMessagesStream error: %s", err)
}

func TestMessagesStreamToolUse(t *testing.T) {
	server := test.NewTestServer()
	server.RegisterHandler("/v1/messages", handlerMessagesStreamToolUse)

	ts := server.AnthropicTestServer()
	ts.Start()
	defer ts.Close()

	baseUrl := ts.URL + "/v1"
	cli := anthropic.NewClient(
		test.GetTestToken(),
		anthropic.WithBaseURL(baseUrl),
	)

	request := anthropic.MessagesStreamRequest{
		MessagesRequest: anthropic.MessagesRequest{
			Model: anthropic.ModelClaude3Opus20240229,
			Messages: []anthropic.Message{
				anthropic.NewUserTextMessage("What is the weather like in San Francisco?"),
			},
			MaxTokens: 1000,
			Tools: []anthropic.ToolDefinition{
				{
					Name:        "get_weather",
					Description: "Get the current weather in a given location",
					InputSchema: jsonschema.Definition{
						Type: jsonschema.Object,
						Properties: map[string]jsonschema.Definition{
							"location": {
								Type:        jsonschema.String,
								Description: "The city and state, e.g. San Francisco, CA",
							},
							"unit": {
								Type:        jsonschema.String,
								Enum:        []string{"celsius", "fahrenheit"},
								Description: "The unit of temperature, either 'celsius' or 'fahrenheit'",
							},
						},
						Required: []string{"location"},
					},
				},
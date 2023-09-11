package test

import (
	"log"
	"net/http"
	"net/http/httptest"
)

const testAPI = "this-is-my-secure-token-do-not-steal!!"

func GetTestToken() string {
	return testAPI
}

type ServerTest struct {
	handlers map[string]handler
}

type handler func(w http.ResponseWriter, r *http.Request)

func NewTestServer() *ServerTest {
	return &ServerTest{handlers: make(map[string]handler)}
}

func (ts *ServerTest) RegisterHandler(path string, handler handler) {
	ts.handlers[path] = handler
}

// AnthropicTestServer Creates a mocked OpenAI server which can pretend to handle 
package anthropic_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/liushuangls/go-anthropic/v2/internal/test"
)

func TestComplete(t *testing.T) {
	server := test.NewTestServer()
	server.RegisterHandler("/v1/complete", handleCompleteEndpoint)

	ts := server.A
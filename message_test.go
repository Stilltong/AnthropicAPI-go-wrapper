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
	server.RegisterHandler("/v1/mes
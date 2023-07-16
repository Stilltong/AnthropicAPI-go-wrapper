package anthropic

import (
	"context"
	"net/http"
)

type CompleteRequest struct {
	Model             string `json:"model"`
	Prompt           
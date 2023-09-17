package jsonschema_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/liushuangls/go-anthropic/v2/jsonschema"
)

func TestDefinition_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		def
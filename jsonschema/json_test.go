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
		def  jsonschema.Definition
		want string
	}{
		{
			name: "Test with empty Definition",
			def:  jsonschema.Definition{},
			want: `{"propertie
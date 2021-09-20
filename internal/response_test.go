package internal

import (
	"encoding/json"
	"errors"
	"github.com/techcraftt/base/api"
	"os"
	"testing"
)

func TestResponse_MarshalJSON(t *testing.T) {
	resp := api.DivResponse{
		Answer:  10,
		Err:     "error occurred",
		Message: "error message",
	}

	headers := map[string]string{
		"Content-Type":   "cTypeJson",
		"Content-Length": "15",
	}

	response := &Response{
		StatusCode:  200,
		Payload:     resp,
		PayloadType: JsonPayload,
		Headers:     headers,
		Error:       errors.New("testing error"),
	}
	_ = json.NewEncoder(os.Stdout).Encode(response)
}

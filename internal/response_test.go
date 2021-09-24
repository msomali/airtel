/*
 * MIT License
 *
 * Copyright (c) 2021 TECHCRAFT TECHNOLOGIES CO LTD.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package internal

import (
	"encoding/json"
	"errors"
	"github.com/techcraftlabs/airtel/api"
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

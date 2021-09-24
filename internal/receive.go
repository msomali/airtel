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
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

//
//type Receipt struct {
//	Method       string
//	OriginalIP   string
//	ForwardingIP string
//	Headers      map[string]string
//	Params       map[string]string
//	Payload      interface{}
//}

// ReceivePayload takes *http.Request from clients like during then unmarshal the provided
// request into given interface v
// The expected Content-Type should also be declared. If its cTypeJson or
// application/xml
func ReceivePayload(r *http.Request, v interface{}) error {

	contentType := r.Header.Get("Content-Type")
	payloadType := categorizeContentType(contentType)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if v == nil {
		return fmt.Errorf("v can not be nil")
	}
	// restore request body
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	switch payloadType {
	case JsonPayload:
		err := json.NewDecoder(r.Body).Decode(v)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		return err

	case XmlPayload:
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		return xml.Unmarshal(body, v)
	}

	return err
}

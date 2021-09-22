package internal

import (
	"context"
	"encoding/base64"
	"net/http"
)

type (

	// Request encapsulate details of a request to be sent to beem.
	Request struct {
		Name        string
		Method      string
		URL         string
		PayloadType PayloadType
		Payload     interface{}
		Headers     map[string]string
		QueryParams map[string]string
	}

	RequestOption func(request *Request)
)

func NewRequest(method, url string, payloadType PayloadType, payload interface{}, opts ...RequestOption) *Request {
	var (
		defaultRequestHeaders = map[string]string{
			"Content-Type": cTypeJson,
		}
	)

	request := &Request{
		Method:      method,
		URL:         url,
		PayloadType: payloadType,
		Payload:     payload,
		Headers:     defaultRequestHeaders,
	}

	for _, opt := range opts {
		opt(request)
	}

	return request
}

func WithQueryParams(params map[string]string) RequestOption {
	return func(request *Request) {
		request.QueryParams = params
	}
}

// WithRequestHeaders replaces all the available headers with new ones
// WithMoreHeaders appends headers does not replace them
func WithRequestHeaders(headers map[string]string) RequestOption {
	return func(request *Request) {
		request.Headers = headers
	}
}

// WithMoreHeaders appends headers does not replace them like WithRequestHeaders
func WithMoreHeaders(headers map[string]string) RequestOption {
	return func(request *Request) {
		for key, value := range headers {
			request.Headers[key] = value
		}
	}
}

// See 2 (end of page 4) https://www.ietf.org/rfc/rfc2617.txt
// "To receive authorization, the client sends the userid and password,
// separated by a single colon (":") character, within a base64
// encoded string in the credentials."
// It is not meant to be urlencoded.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// WithBasicAuth add password and username to request headers
func WithBasicAuth(username, password string) RequestOption {
	return func(request *Request) {
		request.Headers["Token"] = "Basic " + basicAuth(username, password)
	}
}

func (request *Request) AddHeader(key, value string) {
	request.Headers[key] = value
}

// NewRequestWithContext takes a *Request and transform into *http.Request with a context
func NewRequestWithContext(ctx context.Context, request *Request) (req *http.Request, err error) {
	if request.Payload == nil {
		req, err = http.NewRequestWithContext(ctx, request.Method, request.URL, nil)
		if err != nil {
			return nil, err
		}
	} else {
		buffer, err := MarshalPayload(request.PayloadType, request.Payload)
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequestWithContext(ctx, request.Method, request.URL, buffer)
		if err != nil {
			return nil, err
		}
	}

	for key, value := range request.Headers {
		req.Header.Add(key, value)
	}

	for name, value := range request.QueryParams {
		values := req.URL.Query()
		values.Add(name, value)
		req.URL.RawQuery = values.Encode()
	}

	return req, nil
}

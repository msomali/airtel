package http

import (
	"context"
	"fmt"
	"github.com/techcraftt/base/api"
	"github.com/techcraftt/base/internal"
	"net/http"
)

var _ api.Service = (*Client)(nil)

type (
	Client struct {
		BaseURL string
		Port    uint64
		base    *internal.BaseClient
	}
)

func NewClient(base string, port uint64) *Client {
	return &Client{
		BaseURL: base,
		Port:    port,
		base:    internal.NewBaseClient(internal.WithDebugMode(true)),
	}
}

func (c *Client) requestURL(endpoint string) string {
	return fmt.Sprintf("http://%s:%d/%s", c.BaseURL, c.Port, endpoint)
}

func (c *Client) Add(a, b int64) (int64, error) {
	req := api.Request{
		A: a,
		B: b,
	}
	request := internal.NewRequest(http.MethodGet, c.requestURL("add"), internal.JsonPayload, req)

	var response api.Response
	err := c.base.Send(context.TODO(), "add", request, &response)
	if err != nil {
		return 0, err
	}

	return response.Answer, nil
}

func (c *Client) Divide(a, b int64) (int64, error) {
	req := api.Request{
		A: a,
		B: b,
	}
	request := internal.NewRequest(http.MethodGet, c.requestURL("div"), internal.JsonPayload, req)

	response := new(api.DivResponse)
	do, err := c.base.Do(context.TODO(), "divide", request, response)
	if err != nil {
		return 0, err
	}
	//jsonString, _ := do.MarshalJSON()
	//log.Printf("the response returned by Do is %s\n", string(jsonString))
	if do.Error != nil {
		return 0, fmt.Errorf("%s:%s", response.Err, response.Message)
	}

	return response.Answer, nil
}

package http

import (
	"context"
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/api"
	"github.com/techcraftlabs/airtel/internal"
)

var _ api.Service = (*Client)(nil)

type (
	Config struct {
		BaseURL    string
		Port       uint64
		DebugMode bool
	}
	Client struct {
		conf *Config
		reqAdapter api.RequestAdapter
		resAdapter api.ResponseAdapter
		base       *internal.BaseClient
		airtel     *airtel.Client
	}
)


func NewClient(conf *Config, client *airtel.Client) *Client {
	return &Client{
		conf : conf,
		airtel: client,
		base:    internal.NewBaseClient(internal.WithDebugMode(conf.DebugMode)),
	}
}

func (c *Client) Push(ctx context.Context, request api.PushPayRequest) (api.PushPayResponse, error) {
	pushRequest := c.reqAdapter.ToPushPayRequest(request)
	pushResponse, err := c.airtel.Push(ctx,pushRequest)
	if err != nil{
		return api.PushPayResponse{}, err
	}
	response := c.resAdapter.ToPushPayResponse(pushResponse)
	return response,nil
}

func (c *Client) Disburse(ctx context.Context, request api.DisburseRequest) (api.DisburseResponse, error) {
	disburseRequest, err := c.reqAdapter.ToDisburseRequest(request)
	if err != nil {
		return api.DisburseResponse{}, err
	}

	disburseResponse, err := c.airtel.Disburse(ctx, disburseRequest)
	if err != nil {
		return api.DisburseResponse{}, err
	}
	response := c.resAdapter.ToDisburseResponse(disburseResponse)

	return response,nil
}


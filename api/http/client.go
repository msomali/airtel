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

package http

import (
	"context"
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/api"
	"github.com/techcraftlabs/airtel/internal/models"
)

var _ api.Service = (*Client)(nil)

type (
	Client struct {
		reqAdapter api.RequestAdapter
		resAdapter api.ResponseAdapter
		airtel     *airtel.Client
	}
)

func (c *Client) Token(ctx context.Context) (models.AirtelAuthResponse, error) {
	return c.airtel.Token(ctx)
}

func (c *Client) Summary(ctx context.Context, params airtel.Params) (models.ListTransactionsResponse, error) {
	return c.airtel.Summary(ctx, params)
}

func NewClient(client *airtel.Client) *Client {

	airtelConf := client.Conf
	return &Client{
		reqAdapter: &api.ReqAdapter{Conf: airtelConf},
		resAdapter: &api.ResAdapter{},
		airtel:     client,
	}
}

func (c *Client) Push(ctx context.Context, request api.PushPayRequest) (api.PushPayResponse, error) {
	pushRequest := c.reqAdapter.ToPushPayRequest(request)
	pushResponse, err := c.airtel.Push(ctx, pushRequest)
	if err != nil {
		return api.PushPayResponse{}, err
	}
	response := c.resAdapter.ToPushPayResponse(pushResponse)
	return response, nil
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

	return response, nil
}

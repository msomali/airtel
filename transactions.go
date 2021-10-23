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

package airtel

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel/models"
	"github.com/techcraftlabs/base"
)

type (
	// Params
	//From	query	integer(int64)	true	Date from which transactions are to be fetched.
	//To	query	integer(int64)	true	Date until transactions are to be fetched.
	//Limit	query	integer(int64)	true	The number of transactions to be fetched on a page.
	//Offset	query	integer(int64)	true	Page number from which transactions are to be fetched.
	Params struct {
		From   int64 `json:"from"`
		To     int64 `json:"to"`
		Limit  int64 `json:"limit"`
		Offset int64 `json:"offset"`
	}
	TransactionService interface {
		Summary(ctx context.Context, params Params) (models.ListTransactionsResponse, error)
	}
)

func queryParamsOptions(params Params, m map[string]string) base.RequestOption {
	from, to, limit, offset := params.From, params.To, params.Limit, params.Offset
	if from > 0 {
		m["from"] = fmt.Sprintf("%d", from)
	}
	if to > 0 {
		m["to"] = fmt.Sprintf("%d", to)
	}
	if limit > 0 {
		m["limit"] = fmt.Sprintf("%d", limit)
	}
	if offset > 0 {
		m["offset"] = fmt.Sprintf("%d", offset)
	}

	return base.WithQueryParams(m)
}

func (c *Client) Summary(ctx context.Context, params Params) (models.ListTransactionsResponse, error) {

	token, err := c.checkToken(ctx)
	if err != nil {
		return models.ListTransactionsResponse{}, err
	}

	var opts []base.RequestOption

	hs := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "*/*",
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	queryMap := make(map[string]string, 4)
	queryMapOpt := queryParamsOptions(params, queryMap)
	headersOpt := base.WithRequestHeaders(hs)
	opts = append(opts, headersOpt, queryMapOpt)
	req := c.makeInternalRequest(TransactionSummary, nil, opts...)

	res := new(models.ListTransactionsResponse)

	_, err = c.base.Do(ctx, "transaction summary", req, res)
	if err != nil {
		return models.ListTransactionsResponse{}, err
	}
	return *res, nil
}

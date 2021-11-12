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
	"github.com/techcraftlabs/airtel/pkg/countries"
	"github.com/techcraftlabs/base"
)

type (
	BalanceRequest struct {
		MSISDN               string
		CountryOfTransaction string
	}
	BalanceResponse struct {
		Data struct {
			Balance       string `json:"balance"`
			Currency      string `json:"currency"`
			AccountStatus string `json:"account_status"`
		} `json:"data"`
		Status struct {
			Code       string `json:"code"`
			Message    string `json:"message"`
			ResultCode string `json:"result_code"`
			Success    bool   `json:"success"`
		} `json:"status"`
		ErrorDescription string `json:"error_description,omitempty"`
		Error            string `json:"error,omitempty"`
		StatusMessage    string `json:"status_message,omitempty"`
		StatusCode       string `json:"status_code,omitempty"`
	}

	AccountService interface {
		Balance(ctx context.Context, request BalanceRequest) (BalanceResponse, error)
	}
)

func (c *Client) Balance(ctx context.Context, request BalanceRequest) (BalanceResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return BalanceResponse{}, err
	}

	countryName := request.CountryOfTransaction
	country, err := countries.Get(countryName)
	if err != nil {
		return BalanceResponse{}, err
	}
	var opts []base.RequestOption

	hs := map[string]string{
		"Content-Type":   "application/json",
		"Accept":         "*/*",
		"X-Country":      country.CodeName,
		"X-CurrencyName": country.CurrencyCode,
		"Authorization":  fmt.Sprintf("Bearer %s", token),
	}
	headersOpt := base.WithRequestHeaders(hs)
	endpointOption := base.WithEndpoint(request.MSISDN)
	opts = append(opts, headersOpt, endpointOption)
	req := c.makeInternalRequest(BalanceEnquiry, request, opts...)
	res := new(BalanceResponse)
	_, err = c.base.Do(ctx, "balance enquiry", req, res)
	if err != nil {
		return BalanceResponse{}, err
	}
	return *res, nil
}

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
	"github.com/techcraftlabs/airtel/internal"
	"github.com/techcraftlabs/airtel/internal/models"
	"github.com/techcraftlabs/airtel/pkg/countries"
)

type DisbursementService interface {
	disburse(ctx context.Context, request models.DisburseRequest) (models.DisburseResponse, error)
	TransactionEnquiry(ctx context.Context, response models.DisburseEnquiryRequest) (models.DisburseEnquiryResponse, error)
}

func (c *Client) Disburse(ctx context.Context, request DisburseRequest) (DisburseResponse, error) {
	disburseRequest, err := c.reqAdapter.ToDisburseRequest(request)
	if err != nil {
		return DisburseResponse{}, err
	}

	disburseResponse, err := c.disburse(ctx, disburseRequest)
	if err != nil {
		return DisburseResponse{}, err
	}
	response := c.resAdapter.ToDisburseResponse(disburseResponse)

	return response, nil
}

func (c *Client) disburse(ctx context.Context, request models.DisburseRequest) (models.DisburseResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.DisburseResponse{}, err
	}

	countryName := request.CountryOfTransaction
	country, err := countries.GetByName(countryName)
	if err != nil {
		return models.DisburseResponse{}, err
	}
	var opts []internal.RequestOption

	hs := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "*/*",
		"X-Country":     country.CodeName,
		"X-Currency":    country.CurrencyCode,
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	headersOpt := internal.WithRequestHeaders(hs)
	opts = append(opts, headersOpt)

	req := c.makeInternalRequest(Disbursement, request, opts...)
	res := new(models.DisburseResponse)
	env := c.Conf.Environment
	rn := fmt.Sprintf("%v: %s: %s", env, Disbursement.Group(), Disbursement.Name())
	_, err = c.base.Do(ctx, rn, req, res)
	if err != nil {
		return models.DisburseResponse{}, err
	}
	return *res, nil
}

func (c *Client) TransactionEnquiry(ctx context.Context, request models.DisburseEnquiryRequest) (models.DisburseEnquiryResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.DisburseEnquiryResponse{}, err
	}

	countryName := request.CountryOfTransaction
	country, err := countries.GetByName(countryName)
	if err != nil {
		return models.DisburseEnquiryResponse{}, err
	}
	var opts []internal.RequestOption

	hs := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "*/*",
		"X-Country":     country.CodeName,
		"X-Currency":    country.CurrencyCode,
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	headersOpt := internal.WithRequestHeaders(hs)
	endpointOption := internal.WithEndpoint(request.ID)
	opts = append(opts, headersOpt, endpointOption)
	req := c.makeInternalRequest(DisbursementEnquiry, request, opts...)
	res := new(models.DisburseEnquiryResponse)
	_, err = c.base.Do(ctx, "disbursement enquiry", req, res)
	if err != nil {
		return models.DisburseEnquiryResponse{}, err
	}
	return *res, nil
}

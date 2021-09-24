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
	"net/http"
)

type DisbursementService interface {
	Disburse(ctx context.Context, request models.AirtelDisburseRequest) (models.AirtelDisburseResponse, error)
	TransactionEnquiry(ctx context.Context, response models.AirtelDisburseEnquiryRequest) (models.AirtelDisburseEnquiryResponse, error)
}

func (c *Client) Disburse(ctx context.Context, request models.AirtelDisburseRequest) (models.AirtelDisburseResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelDisburseResponse{}, err
	}

	countryName := request.CountryOfTransaction
	country, err := countries.GetByName(countryName)
	if err != nil {
		return models.AirtelDisburseResponse{}, err
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
	reqUrl := requestURL(c.Conf.Environment, Disbursement)
	req := internal.NewRequest(http.MethodPost, reqUrl, request, opts...)
	res := new(models.AirtelDisburseResponse)
	_, err = c.base.Do(ctx, "disbursement", req, res)
	if err != nil {
		return models.AirtelDisburseResponse{}, err
	}
	return *res, nil
}

func (c *Client) TransactionEnquiry(ctx context.Context, request models.AirtelDisburseEnquiryRequest) (models.AirtelDisburseEnquiryResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelDisburseEnquiryResponse{}, err
	}

	countryName := request.CountryOfTransaction
	country, err := countries.GetByName(countryName)
	if err != nil {
		return models.AirtelDisburseEnquiryResponse{}, err
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
	reqUrl := requestURL(c.Conf.Environment, DisbursementEnquiry)
	req := internal.NewRequest(http.MethodGet, reqUrl, request, opts...)
	res := new(models.AirtelDisburseEnquiryResponse)
	_, err = c.base.Do(ctx, "disbursement enquiry", req, res)
	if err != nil {
		return models.AirtelDisburseEnquiryResponse{}, err
	}
	return *res, nil
}

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

type (
	AccountService interface {
		Balance(ctx context.Context, request models.AirtelBalanceEnquiryRequest) (models.AirtelBalanceEnquiryResponse, error)
	}
)

func (c *Client) Balance(ctx context.Context, request models.AirtelBalanceEnquiryRequest) (models.AirtelBalanceEnquiryResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelBalanceEnquiryResponse{}, err
	}

	countryName := request.CountryOfTransaction
	country, err := countries.GetByName(countryName)
	if err != nil {
		return models.AirtelBalanceEnquiryResponse{}, err
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
	endpointOption := internal.WithEndpoint(request.MSISDN)
	opts = append(opts, headersOpt, endpointOption)
	reqUrl := c.requestURL(BalanceEnquiry)
	req := internal.NewRequest(http.MethodGet, reqUrl, request, opts...)
	res := new(models.AirtelBalanceEnquiryResponse)
	_, err = c.base.Do(ctx, "balance enquiry", req, res)
	if err != nil {
		return models.AirtelBalanceEnquiryResponse{}, err
	}
	return *res, nil
}

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

type CollectionService interface {
	Push(ctx context.Context, request models.AirtelPushRequest) (models.AirtelPushResponse, error)
	Refund(ctx context.Context, request models.AirtelRefundRequest) (models.AirtelRefundResponse, error)
	Enquiry(ctx context.Context, request models.AirtelPushEnquiryRequest) (models.AirtelPushEnquiryResponse, error)
	CallbackServeHTTP(writer http.ResponseWriter, request *http.Request)
}

func (c *Client) Push(ctx context.Context, request models.AirtelPushRequest) (models.AirtelPushResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelPushResponse{}, err
	}

	transaction := request.Transaction
	countryCodeName := transaction.Country
	currencyCodeName := transaction.Currency

	var opts []internal.RequestOption

	hs := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "*/*",
		"X-Country":     countryCodeName,
		"X-Currency":    currencyCodeName,
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	headersOpt := internal.WithRequestHeaders(hs)
	opts = append(opts, headersOpt)

	req := c.makeInternalRequest(UssdPush,request, opts...)

	if err != nil {
		return models.AirtelPushResponse{}, err
	}

	res := new(models.AirtelPushResponse)
	_, err = c.base.Do(ctx, "ussd push", req, res)
	if err != nil {
		return models.AirtelPushResponse{}, err
	}
	return *res, nil
}

func (c *Client) Refund(ctx context.Context, request models.AirtelRefundRequest) (models.AirtelRefundResponse, error) {
	country, err := countries.GetByName(request.CountryOfTransaction)
	if err != nil {
		return models.AirtelRefundResponse{}, err
	}
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelRefundResponse{}, err
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

	req := c.makeInternalRequest(Refund, request, opts...)

	if err != nil {
		return models.AirtelRefundResponse{}, err
	}

	res := new(models.AirtelRefundResponse)
	env := c.Conf.Environment
	rn := fmt.Sprintf("%v: %s: %s",env,Refund.Group(), Refund.Name())
	_, err = c.base.Do(ctx, rn, req, res)
	if err != nil {
		return models.AirtelRefundResponse{}, err
	}
	return *res, nil

}

func (c *Client) Enquiry(ctx context.Context, request models.AirtelPushEnquiryRequest) (models.AirtelPushEnquiryResponse, error) {

	country, err := countries.GetByName(request.CountryOfTransaction)
	if err != nil {
		return models.AirtelPushEnquiryResponse{}, err
	}
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelPushEnquiryResponse{}, err
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
	endpointOpt := internal.WithEndpoint(request.ID)
	opts = append(opts, headersOpt, endpointOpt)
	req := c.makeInternalRequest(PushEnquiry, request, opts...)
	if err != nil {
		return models.AirtelPushEnquiryResponse{}, err
	}
	reqName := PushEnquiry.Name()
	res := new(models.AirtelPushEnquiryResponse)
	_, err = c.base.Do(ctx, reqName, req, res)
	if err != nil {
		return models.AirtelPushEnquiryResponse{}, err
	}
	return *res, nil
}

func (c *Client) CallbackServeHTTP(writer http.ResponseWriter, request *http.Request) {
	body := new(models.AirtelCallbackRequest)
	err := internal.ReceivePayload(request, body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	reqBody := *body

	//todo: check the hash if it is OK
	err = c.pushCallbackFunc.Handle(reqBody)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

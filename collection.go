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
	"github.com/techcraftlabs/airtel/internal/models"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"github.com/techcraftlabs/base"
	"net/http"
)

type CollectionService interface {
	push(ctx context.Context, request models.PushRequest) (models.PushResponse, error)
	Refund(ctx context.Context, request models.RefundRequest) (models.RefundResponse, error)
	PushEnquiry(ctx context.Context, request models.PushEnquiryRequest) (models.PushEnquiryResponse, error)
	CallbackServeHTTP(writer http.ResponseWriter, request *http.Request)
}

func (c *Client) Push(ctx context.Context, request PushPayRequest) (PushPayResponse, error) {

	pushRequest := c.reqAdapter.ToPushPayRequest(request)
	pushResponse, err := c.push(ctx, pushRequest)
	if err != nil {
		return PushPayResponse{}, err
	}
	response := c.resAdapter.ToPushPayResponse(pushResponse)
	return response, nil
}

func (c *Client) push(ctx context.Context, request models.PushRequest) (models.PushResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.PushResponse{}, err
	}

	transaction := request.Transaction
	countryCodeName := transaction.Country
	currencyCodeName := transaction.Currency

	var opts []base.RequestOption

	hs := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "*/*",
		"X-Country":     countryCodeName,
		"X-Currency":    currencyCodeName,
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	headersOpt := base.WithRequestHeaders(hs)
	opts = append(opts, headersOpt)

	req := c.makeInternalRequest(UssdPush, request, opts...)

	if err != nil {
		return models.PushResponse{}, err
	}

	res := new(models.PushResponse)
	_, err = c.base.Do(ctx, "ussd push", req, res)
	if err != nil {
		return models.PushResponse{}, err
	}
	return *res, nil
}

func (c *Client) Refund(ctx context.Context, request models.RefundRequest) (models.RefundResponse, error) {
	country, err := countries.GetByName(request.CountryOfTransaction)
	if err != nil {
		return models.RefundResponse{}, err
	}
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.RefundResponse{}, err
	}
	var opts []base.RequestOption
	hs := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "*/*",
		"X-Country":     country.CodeName,
		"X-Currency":    country.CurrencyCode,
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	headersOpt := base.WithRequestHeaders(hs)
	opts = append(opts, headersOpt)

	req := c.makeInternalRequest(Refund, request, opts...)

	if err != nil {
		return models.RefundResponse{}, err
	}

	res := new(models.RefundResponse)
	env := c.Conf.Environment
	rn := fmt.Sprintf("%v: %s: %s", env, Refund.Group(), Refund.name())
	_, err = c.base.Do(ctx, rn, req, res)
	if err != nil {
		return models.RefundResponse{}, err
	}
	return *res, nil

}

func (c *Client) PushEnquiry(ctx context.Context, request models.PushEnquiryRequest) (models.PushEnquiryResponse, error) {

	country, err := countries.GetByName(request.CountryOfTransaction)
	if err != nil {
		return models.PushEnquiryResponse{}, err
	}
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.PushEnquiryResponse{}, err
	}
	var opts []base.RequestOption
	hs := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "*/*",
		"X-Country":     country.CodeName,
		"X-Currency":    country.CurrencyCode,
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	headersOpt := base.WithRequestHeaders(hs)
	endpointOpt := base.WithEndpoint(request.ID)
	opts = append(opts, headersOpt, endpointOpt)
	req := c.makeInternalRequest(PushEnquiry, request, opts...)
	if err != nil {
		return models.PushEnquiryResponse{}, err
	}
	reqName := PushEnquiry.name()
	res := new(models.PushEnquiryResponse)
	_, err = c.base.Do(ctx, reqName, req, res)
	if err != nil {
		return models.PushEnquiryResponse{}, err
	}
	return *res, nil
}

func (c *Client) CallbackServeHTTP(writer http.ResponseWriter, request *http.Request) {
	body := new(models.CallbackRequest)
	err := base.ReceivePayload(request, body)
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

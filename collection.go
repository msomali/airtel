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
	"net/http"
)

var _ CollectionAdapter = (*collectAdapter)(nil)

type (
	collectAdapter struct{}

	CollectionAdapter interface {
		ToInternalRequest(request PushPayRequest) InternalPushRequest
		ToModifiedResponse(response InternalPushResponse) PushPayResponse
	}

	InternalPushRequest struct {
		Reference  string `json:"reference"`
		Subscriber struct {
			Country  string `json:"country"`
			Currency string `json:"currency"`
			Msisdn   string `json:"msisdn"`
		} `json:"subscriber"`
		Transaction struct {
			Amount   float64 `json:"amount"`
			Country  string  `json:"country"`
			Currency string  `json:"currency"`
			ID       string  `json:"id"`
		} `json:"transaction"`
	}
	InternalPushEnquiryRequest struct {
		ID                   string `json:"id"`
		CountryOfTransaction string `json:"country"`
	}

	InternalPushResponse struct {
		Data struct {
			Transaction struct {
				ID     string `json:"id,omitempty"`
				Status string `json:"status,omitempty"`
			} `json:"transaction,omitempty"`
		} `json:"data,omitempty"`
		Status struct {
			Code       string `json:"code,omitempty"`
			Message    string `json:"message,omitempty"`
			ResultCode string `json:"result_code,omitempty"`
			Success    bool   `json:"success,omitempty"`
		} `json:"status,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
		Error            string `json:"error,omitempty"`
		StatusMessage    string `json:"status_message,omitempty"`
		StatusCode       string `json:"status_code,omitempty"`
	}

	InternalPushEnquiryResponse struct {
		Data struct {
			Transaction struct {
				AirtelMoneyID string `json:"airtel_money_id,omitempty"`
				ID            string `json:"id,omitempty"`
				Message       string `json:"message,omitempty"`
				Status        string `json:"status,omitempty"`
			} `json:"transaction,omitempty"`
		} `json:"data,omitempty"`
		Status struct {
			Code       string `json:"code,omitempty"`
			Message    string `json:"message,omitempty"`
			ResultCode string `json:"result_code,omitempty"`
			Success    bool   `json:"success,omitempty"`
		} `json:"status,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
		Error            string `json:"error,omitempty"`
		StatusMessage    string `json:"status_message,omitempty"`
		StatusCode       string `json:"status_code,omitempty"`
	}

	InternalRefundRequest struct {
		CountryOfTransaction string `json:"-"`
		Transaction          struct {
			AirtelMoneyID string `json:"airtel_money_id"`
		} `json:"transaction"`
	}

	CallbackRequest struct {
		Transaction struct {
			ID            string `json:"id"`
			Message       string `json:"message"`
			StatusCode    string `json:"status_code"`
			AirtelMoneyID string `json:"airtel_money_id"`
		} `json:"transaction"`
		Hash string `json:"hash"`
	}

	InternalRefundResponse struct {
		Data struct {
			Transaction struct {
				AirtelMoneyID string `json:"airtel_money_id,omitempty"`
				Status        string `json:"status,omitempty"`
			} `json:"transaction,omitempty"`
		} `json:"data,omitempty"`
		Status struct {
			Code       string `json:"code,omitempty"`
			Message    string `json:"message,omitempty"`
			ResultCode string `json:"result_code,omitempty"`
			Success    bool   `json:"success,omitempty"`
		} `json:"status,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
		Error            string `json:"error,omitempty"`
		StatusMessage    string `json:"status_message,omitempty"`
		StatusCode       string `json:"status_code,omitempty"`
	}
	CollectionService interface {
		push(ctx context.Context, request InternalPushRequest) (InternalPushResponse, error)
		Refund(ctx context.Context, request InternalRefundRequest) (InternalRefundResponse, error)
		PushEnquiry(ctx context.Context, request InternalPushEnquiryRequest) (InternalPushEnquiryResponse, error)
		CallbackServeHTTP(writer http.ResponseWriter, request *http.Request)
	}
)

func (c *Client) Push(ctx context.Context, request PushPayRequest) (PushPayResponse, error) {

	pushRequest := c.collectionAdapter.ToInternalRequest(request)
	pushResponse, err := c.push(ctx, pushRequest)
	if err != nil {
		return PushPayResponse{}, err
	}
	response := c.collectionAdapter.ToModifiedResponse(pushResponse)
	return response, nil
}

func (c *Client) push(ctx context.Context, request InternalPushRequest) (InternalPushResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return InternalPushResponse{}, err
	}

	transaction := request.Transaction
	countryCodeName := transaction.Country
	currencyCodeName := transaction.Currency

	var opts []base.RequestOption

	hs := map[string]string{
		"Content-Type":   "application/json",
		"Accept":         "*/*",
		"X-Country":      countryCodeName,
		"X-CurrencyName": currencyCodeName,
		"Authorization":  fmt.Sprintf("Bearer %s", token),
	}

	headersOpt := base.WithRequestHeaders(hs)
	opts = append(opts, headersOpt)

	req := c.makeInternalRequest(UssdPush, request, opts...)

	if err != nil {
		return InternalPushResponse{}, err
	}

	res := new(InternalPushResponse)
	_, err = c.base.Do(ctx, req, res)
	if err != nil {
		return InternalPushResponse{}, err
	}
	return *res, nil
}

func (c *Client) Refund(ctx context.Context, request InternalRefundRequest) (InternalRefundResponse, error) {
	country, err := countries.Get(request.CountryOfTransaction)
	if err != nil {
		return InternalRefundResponse{}, err
	}
	token, err := c.checkToken(ctx)
	if err != nil {
		return InternalRefundResponse{}, err
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
	opts = append(opts, headersOpt)

	req := c.makeInternalRequest(Refund, request, opts...)

	if err != nil {
		return InternalRefundResponse{}, err
	}

	res := new(InternalRefundResponse)
	//env := c.Conf.Environment
	//rn := fmt.Sprintf("%v: %s: %s", env, Refund.Group(), Refund.name())
	_, err = c.base.Do(ctx, req, res)
	if err != nil {
		return InternalRefundResponse{}, err
	}
	return *res, nil

}

func (c *Client) PushEnquiry(ctx context.Context, request InternalPushEnquiryRequest) (InternalPushEnquiryResponse, error) {

	country, err := countries.Get(request.CountryOfTransaction)
	if err != nil {
		return InternalPushEnquiryResponse{}, err
	}
	token, err := c.checkToken(ctx)
	if err != nil {
		return InternalPushEnquiryResponse{}, err
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
	endpointOpt := base.WithEndpoint(request.ID)
	opts = append(opts, headersOpt, endpointOpt)
	req := c.makeInternalRequest(PushEnquiry, request, opts...)
	if err != nil {
		return InternalPushEnquiryResponse{}, err
	}
	//reqName := PushEnquiry.name()
	res := new(InternalPushEnquiryResponse)
	_, err = c.base.Do(ctx, req, res)
	if err != nil {
		return InternalPushEnquiryResponse{}, err
	}
	return *res, nil
}

func (c *Client) CallbackServeHTTP(writer http.ResponseWriter, request *http.Request) {
	body := new(CallbackRequest)
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

func (p *collectAdapter) ToInternalRequest(request PushPayRequest) InternalPushRequest {
	subCountry, _ := countries.Get(request.SubscriberCountry)
	transCountry, _ := countries.Get(request.TransactionCountry)
	return InternalPushRequest{
		Reference: request.Reference,
		Subscriber: struct {
			Country  string `json:"country"`
			Currency string `json:"currency"`
			Msisdn   string `json:"msisdn"`
		}{
			Country:  subCountry.CodeName,
			Currency: subCountry.CurrencyCode,
			Msisdn:   request.SubscriberMsisdn,
		},
		Transaction: struct {
			Amount   float64 `json:"amount"`
			Country  string  `json:"country"`
			Currency string  `json:"currency"`
			ID       string  `json:"id"`
		}{
			Amount:   request.TransactionAmount,
			Country:  transCountry.CodeName,
			Currency: transCountry.CurrencyCode,
			ID:       request.TransactionID,
		},
	}
}

func (p *collectAdapter) ToModifiedResponse(response InternalPushResponse) PushPayResponse {
	transaction := response.Data.Transaction
	status := response.Status

	if !status.Success {
		return PushPayResponse{
			ResultCode:    status.ResultCode,
			Success:       status.Success,
			StatusMessage: status.Message,
			StatusCode:    status.Code,
		}
	}
	isErr := response.Error != "" && response.ErrorDescription != ""
	if isErr {
		resp := PushPayResponse{
			ErrorDescription: response.ErrorDescription,
			Error:            response.Error,
			StatusMessage:    response.StatusMessage,
			StatusCode:       response.StatusCode,
		}
		return resp
	}

	return PushPayResponse{
		ID:               transaction.ID,
		Status:           transaction.Status,
		ResultCode:       status.ResultCode,
		Success:          status.Success,
		ErrorDescription: response.ErrorDescription,
		Error:            response.Error,
		StatusMessage:    response.StatusMessage,
		StatusCode:       response.StatusCode,
	}
}

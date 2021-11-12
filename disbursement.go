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

var _ DisbursementAdapter = (*disburseAdapter)(nil)

type (
	disburseAdapter struct {
		Conf *Config
	}

	DisbursementAdapter interface {
		ToModifiedResponse(response InternalDisburseResponse) DisburseResponse
		ToInternalRequest(request DisburseRequest) (InternalDisburseRequest, error)
	}

	DisburseEnquiryResponse struct {
		Data struct {
			Transaction struct {
				ID      string `json:"id"`
				Message string `json:"message"`
				Status  string `json:"status"`
			} `json:"transaction"`
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

	DisburseEnquiryRequest struct {
		CountryOfTransaction string
		ID                   string `json:"id"`
	}

	DisbursementService interface {
		Disburse(ctx context.Context, request DisburseRequest) (DisburseResponse, error)
		DisburseEnquiry(ctx context.Context, response DisburseEnquiryRequest) (DisburseEnquiryResponse, error)
	}

	InternalDisburseRequest struct {
		CountryOfTransaction string `json:"-"`
		Payee                struct {
			Msisdn string `json:"msisdn"`
		} `json:"payee"`
		Reference   string `json:"reference"`
		Pin         string `json:"pin"`
		Transaction struct {
			Amount float64 `json:"amount"`
			ID     string  `json:"id"`
		} `json:"transaction"`
	}

	InternalDisburseResponse struct {
		Data struct {
			Transaction struct {
				ReferenceID   string `json:"reference_id,omitempty"`
				AirtelMoneyID string `json:"airtel_money_id,omitempty"`
				ID            string `json:"id,omitempty"`
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
)

func (c *Client) Disburse(ctx context.Context, request DisburseRequest) (DisburseResponse, error) {
	disburseRequest, err := c.disburseAdapter.ToInternalRequest(request)
	if err != nil {
		return DisburseResponse{}, err
	}

	disburseResponse, err := c.disburse(ctx, disburseRequest)
	if err != nil {
		return DisburseResponse{}, err
	}
	response := c.disburseAdapter.ToModifiedResponse(disburseResponse)

	return response, nil
}

func (c *Client) disburse(ctx context.Context, request InternalDisburseRequest) (InternalDisburseResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return InternalDisburseResponse{}, err
	}

	countryName := request.CountryOfTransaction
	country, err := countries.Get(countryName)
	if err != nil {
		return InternalDisburseResponse{}, err
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

	req := c.makeInternalRequest(Disbursement, request, opts...)
	res := new(InternalDisburseResponse)
	env := c.Conf.Environment
	rn := fmt.Sprintf("%v: %s: %s", env, Disbursement.Group(), Disbursement.name())
	_, err = c.base.Do(ctx, rn, req, res)
	if err != nil {
		return InternalDisburseResponse{}, err
	}
	return *res, nil
}

func (c *Client) DisburseEnquiry(ctx context.Context, request DisburseEnquiryRequest) (DisburseEnquiryResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return DisburseEnquiryResponse{}, err
	}

	countryName := request.CountryOfTransaction
	country, err := countries.Get(countryName)
	if err != nil {
		return DisburseEnquiryResponse{}, err
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
	endpointOption := base.WithEndpoint(request.ID)
	opts = append(opts, headersOpt, endpointOption)
	req := c.makeInternalRequest(DisbursementEnquiry, request, opts...)
	res := new(DisburseEnquiryResponse)
	_, err = c.base.Do(ctx, "disbursement enquiry", req, res)
	if err != nil {
		return DisburseEnquiryResponse{}, err
	}
	return *res, nil
}

func (d *disburseAdapter) ToModifiedResponse(response InternalDisburseResponse) DisburseResponse {
	isErr := response.Error != "" && response.ErrorDescription != ""
	if isErr {
		resp := DisburseResponse{
			ErrorDescription: response.ErrorDescription,
			Error:            response.Error,
			StatusMessage:    response.StatusMessage,
			StatusCode:       response.StatusCode,
		}

		return resp
	}
	transaction := response.Data.Transaction
	status := response.Status

	return DisburseResponse{
		ID:            transaction.ID,
		Reference:     transaction.ReferenceID,
		AirtelMoneyID: transaction.AirtelMoneyID,
		ResultCode:    status.ResultCode,
		Success:       status.Success,
		StatusMessage: status.Message,
		StatusCode:    status.Code,
	}

}

func (d *disburseAdapter) ToInternalRequest(request DisburseRequest) (InternalDisburseRequest, error) {
	encryptedPin, err := pinEncryption(d.Conf.DisbursePIN, d.Conf.PublicKey)
	if err != nil {
		return InternalDisburseRequest{}, fmt.Errorf("could not encrypt key: %w", err)
	}
	req := InternalDisburseRequest{
		CountryOfTransaction: request.CountryOfTransaction,
		Payee: struct {
			Msisdn string `json:"msisdn"`
		}{
			Msisdn: request.MSISDN,
		},
		Reference: request.Reference,
		Pin:       encryptedPin,
		Transaction: struct {
			Amount float64 `json:"amount"`
			ID     string  `json:"id"`
		}{
			Amount: request.Amount,
			ID:     request.ID,
		},
	}
	return req, nil
}

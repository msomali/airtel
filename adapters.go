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
	"fmt"
	"github.com/techcraftlabs/airtel/internal/models"
	"github.com/techcraftlabs/airtel/pkg/countries"
)

var (
	_ ResponseAdapter = (*adapter)(nil)
	_ RequestAdapter  = (*adapter)(nil)
)

type (
	ResponseAdapter interface {
		ToDisburseResponse(response models.DisburseResponse) DisburseResponse
		ToPushPayResponse(response models.PushResponse) PushPayResponse
	}

	// adapter acts as a default RequestAdapter and ResponseAdapter. This can be replaced
	// by other user defined adapters. and then injected to the client using the
	// Client.SetRequestAdapter or Client.SetResponseAdapter, it is not recommended as of now
	// 26th September 2021
	adapter struct {
		Conf *Config
	}
	RequestAdapter interface {
		ToPushPayRequest(request PushPayRequest) models.PushRequest
		ToDisburseRequest(request DisburseRequest) (models.DisburseRequest, error)
	}

)

func (r *adapter) ToPushPayRequest(request PushPayRequest) models.PushRequest {

	subCountry, _ := countries.GetByName(request.SubscriberCountry)
	transCountry, _ := countries.GetByName(request.TransactionCountry)
	return models.PushRequest{
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
			Amount   int64  `json:"amount"`
			Country  string `json:"country"`
			Currency string `json:"currency"`
			ID       string `json:"id"`
		}{
			Amount:   request.TransactionAmount,
			Country:  transCountry.CodeName,
			Currency: transCountry.CurrencyCode,
			ID:       request.TransactionID,
		},
	}
}

func (r *adapter) ToDisburseRequest(request DisburseRequest) (models.DisburseRequest, error) {
	encryptedPin, err := PinEncryption(r.Conf.DisbursePIN, r.Conf.PublicKey)
	if err != nil {
		return models.DisburseRequest{}, fmt.Errorf("could not encrypt key: %w", err)
	}
	req := models.DisburseRequest{
		CountryOfTransaction: request.CountryOfTransaction,
		Payee: struct {
			Msisdn string `json:"msisdn"`
		}{
			Msisdn: request.MSISDN,
		},
		Reference: request.Reference,
		Pin:       encryptedPin,
		Transaction: struct {
			Amount int64  `json:"amount"`
			ID     string `json:"id"`
		}{
			Amount: request.Amount,
			ID:     request.ID,
		},
	}
	return req, nil
}

func (r *adapter) ToPushPayResponse(response models.PushResponse) PushPayResponse {
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

func (r *adapter) ToDisburseResponse(response models.DisburseResponse) DisburseResponse {

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

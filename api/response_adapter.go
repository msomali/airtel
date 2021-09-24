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

package api

import (
	"github.com/techcraftlabs/airtel/internal/models"
)

var (
	_ ResponseAdapter = (*ResAdapter)(nil)
)

type ResponseAdapter interface {
	ToDisburseResponse(response models.AirtelDisburseResponse) DisburseResponse
	ToPushPayResponse(response models.AirtelPushResponse) PushPayResponse
}

type ResAdapter struct {
}

func (r *ResAdapter) ToPushPayResponse(response models.AirtelPushResponse) PushPayResponse {
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

func (r *ResAdapter) ToDisburseResponse(response models.AirtelDisburseResponse) DisburseResponse {

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

package api

import (
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/pkg/models"
)

var (
	_ ResponseAdapter = (*responseAdapter)(nil)
)

type ResponseAdapter interface {
	ToDisburseResponse(response models.AirtelDisburseResponse) DisburseResponse
	ToPushPayResponse(response models.AirtelPushResponse) PushPayResponse
}

type responseAdapter struct {
	conf *airtel.Config
}

func (r *responseAdapter) ToPushPayResponse(response models.AirtelPushResponse) PushPayResponse {

	isErr := response.Error == "" && response.ErrorDescription == ""
	if isErr {
		resp := PushPayResponse{
			ErrorDescription: response.ErrorDescription,
			Error:            response.Error,
			StatusMessage:    response.StatusMessage,
			StatusCode:       response.StatusCode,
		}
		return resp
	}

	transaction := response.Data.Transaction
	status := response.Status
	return PushPayResponse{
		ID:               transaction.ID,
		Status:           transaction.Status,
		ResultCode:       status.ResultCode,
		Success:          status.Success,
		ErrorDescription:  response.ErrorDescription,
		Error:            response.Error,
		StatusMessage:    response.StatusMessage,
		StatusCode:       response.StatusCode,
	}
}

func (r *responseAdapter) ToDisburseResponse(response models.AirtelDisburseResponse) DisburseResponse {

	isErr := response.Error == "" && response.ErrorDescription == ""
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

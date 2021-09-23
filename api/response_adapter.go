package api

import (
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/pkg/models"
)

var (
	_ ResponseAdapter = (*responseAdapter)(nil)
)

type ResponseAdapter interface {
	ToDisburseResponse(response models.AirtelDisburseResponse) (DisburseResponse, error)
}

type responseAdapter struct {
	conf *airtel.Config
}

func (r *responseAdapter) ToDisburseResponse(response models.AirtelDisburseResponse) (DisburseResponse, error) {

	if response.Data.Transaction.AirtelMoneyID != "" {
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
		}, nil
	}
	resp := DisburseResponse{
		ErrorDescription: response.ErrorDescription,
		Error:            response.Error,
		StatusMessage:    response.StatusMessage,
		StatusCode:       response.StatusCode,
	}

	return resp, nil
}


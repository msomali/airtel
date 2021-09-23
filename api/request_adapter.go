package api

import (
	"fmt"
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/pkg/models"
)

var _ RequestAdapter = (*requestAdapter)(nil)

type RequestAdapter interface {
	ToDisburseRequest(request DisburseRequest) (models.AirtelDisburseRequest, error)
}


type requestAdapter struct {
	conf *airtel.Config
}

func (r *requestAdapter) ToDisburseRequest(request DisburseRequest) (models.AirtelDisburseRequest, error) {
	encryptedPin, err := airtel.PinEncryption(r.conf.DisbursePIN, r.conf.PublicKey)
	if err != nil {
		return models.AirtelDisburseRequest{}, fmt.Errorf("could not encrypt key: %w", err)
	}
	req := models.AirtelDisburseRequest{
		Payee: struct {
			Msisdn string `json:"msisdn"`
		}{
			Msisdn: request.MSISDN,
		},
		Reference: request.Reference,
		Pin:       encryptedPin,
		Transaction: struct {
			Amount int    `json:"amount"`
			ID     string `json:"id"`
		}{
			Amount: request.Amount,
			ID:     request.ID,
		},
	}

	return req, nil
}

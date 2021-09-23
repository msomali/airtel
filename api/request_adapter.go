package api

import (
	"fmt"
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/pkg/models"
)

var _ RequestAdapter = (*ReqAdapter)(nil)

type RequestAdapter interface {
	ToPushPayRequest(request PushPayRequest)models.AirtelPushRequest
	ToDisburseRequest(request DisburseRequest) (models.AirtelDisburseRequest, error)
}


type ReqAdapter struct {
	Conf *airtel.Config
}


func (r *ReqAdapter) ToPushPayRequest(request PushPayRequest) models.AirtelPushRequest {
	return models.AirtelPushRequest{
		Reference: request.Reference,
		Subscriber: struct {
			Country  string `json:"country"`
			Currency string `json:"currency"`
			Msisdn   string `json:"msisdn"`
		}{
			Country:  request.SubscriberCountry,
			Currency: request.SubscriberCurrency,
			Msisdn:   request.SubscriberMsisdn,
		},
		Transaction: struct {
			Amount   int    `json:"amount"`
			Country  string `json:"country"`
			Currency string `json:"currency"`
			ID       string `json:"id"`
		}{
			Amount:   request.TransactionAmount,
			Country:  request.TransactionCountry,
			Currency: request.TransactionCurrency,
			ID:       request.TransactionID,
		},
	}
}

func (r *ReqAdapter) ToDisburseRequest(request DisburseRequest) (models.AirtelDisburseRequest, error) {
	encryptedPin, err := airtel.PinEncryption(r.Conf.DisbursePIN, r.Conf.PublicKey)
	if err != nil {
		return models.AirtelDisburseRequest{}, fmt.Errorf("could not encrypt key: %w", err)
	}
	req := models.AirtelDisburseRequest{
		CountryOfTransaction: request.CountryOfTransaction,
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

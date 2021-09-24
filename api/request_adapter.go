package api

import (
	"fmt"
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/internal/models"
	"github.com/techcraftlabs/airtel/pkg/countries"
)

var _ RequestAdapter = (*ReqAdapter)(nil)

type RequestAdapter interface {
	ToPushPayRequest(request PushPayRequest) models.AirtelPushRequest
	ToDisburseRequest(request DisburseRequest) (models.AirtelDisburseRequest, error)
}

type ReqAdapter struct {
	Conf *airtel.Config
}

func (r *ReqAdapter) ToPushPayRequest(request PushPayRequest) models.AirtelPushRequest {

	subCountry, _ := countries.GetByName(request.SubscriberCountry)
	transCountry, _ := countries.GetByName(request.TransactionCountry)
	return models.AirtelPushRequest{
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
			Amount   int64    `json:"amount"`
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

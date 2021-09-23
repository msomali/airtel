package api

import (
	"context"
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/internal/models"
)

type (
	PushPayRequest struct {
		Reference          string
		SubscriberCountry  string
		SubscriberMsisdn   string
		TransactionAmount  int
		TransactionCountry string
		TransactionID      string
	}

	PushPayResponse struct {
		ID               string `json:"id,omitempty"`
		Status           string `json:"status,omitempty"`
		ResultCode       string `json:"result_code,omitempty"`
		Success          bool   `json:"success,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
		Error            string `json:"error,omitempty"`
		StatusMessage    string `json:"status_message,omitempty"`
		StatusCode       string `json:"status_code,omitempty"`
	}
	DisburseRequest struct {
		ID                   string
		MSISDN               string
		Amount               int
		Reference            string
		CountryOfTransaction string
	}

	DisburseResponse struct {
		ID               string `json:"id,omitempty"`
		Reference        string `json:"reference,omitempty"`
		AirtelMoneyID    string `json:"airtel_money_id,omitempty"`
		ResultCode       string `json:"result_code,omitempty"`
		Success          bool   `json:"success,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
		Error            string `json:"error,omitempty"`
		StatusMessage    string `json:"status_message,omitempty"`
		StatusCode       string `json:"status_code,omitempty"`
	}

	Service interface {
		Push(ctx context.Context, request PushPayRequest) (PushPayResponse, error)
		Disburse(ctx context.Context, request DisburseRequest) (DisburseResponse, error)
		Summary(ctx context.Context, params airtel.Params) (models.ListTransactionsResponse, error)
	}
)

package api

import (
	"context"
)

type (
	DisburseRequest struct {
		ID        string
		MSISDN    string
		Amount    int
		Reference string
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

	Service2 interface {
		Disburse(ctx context.Context, request DisburseRequest) (DisburseResponse, error)
	}
)


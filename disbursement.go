package airtel

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"github.com/techcraftlabs/airtel/pkg/models"
	"time"
)

type DisbursementService interface {
	Disburse(ctx context.Context, request models.AirtelDisburseRequest) (models.AirtelDisburseResponse, error)
	TransactionEnquiry(ctx context.Context, response models.AirtelDisburseEnquiryRequest) (models.AirtelDisburseEnquiryResponse, error)
}

func (c *Client) Disburse(ctx context.Context, request models.AirtelDisburseRequest) (models.AirtelDisburseResponse, error) {
	var token string
	if *c.token == "" {
		str, err := c.Token(ctx)
		if err != nil {
			return models.AirtelDisburseResponse{}, err
		}
		token = fmt.Sprintf("%s", str.AccessToken)
	}
	//Add Auth Header
	if *c.token != "" {
		if !c.tokenExpiresAt.IsZero() && time.Until(c.tokenExpiresAt) < (60*time.Second) {
			if _, err := c.Token(ctx); err != nil {
				return models.AirtelDisburseResponse{}, err
			}
		}
		token = *c.token
	}

	req, err := createInternalRequest(countries.TANZANIA, c.conf.Environment, Disbursement, token, request, "")
	if err != nil {
		return models.AirtelDisburseResponse{}, err
	}

	res := new(models.AirtelDisburseResponse)

	_, err = c.base.Do(ctx, "ussd push", req, res)
	if err != nil {
		return models.AirtelDisburseResponse{}, err
	}
	return *res, nil
}

func (c *Client) TransactionEnquiry(ctx context.Context, request models.AirtelDisburseEnquiryRequest) (models.AirtelDisburseEnquiryResponse, error) {
	var token string
	if *c.token == "" {
		str, err := c.Token(ctx)
		if err != nil {
			return models.AirtelDisburseEnquiryResponse{}, err
		}
		token = fmt.Sprintf("%s", str.AccessToken)
	}
	//Add Auth Header
	if *c.token != "" {
		if !c.tokenExpiresAt.IsZero() && time.Until(c.tokenExpiresAt) < (60*time.Second) {
			if _, err := c.Token(ctx); err != nil {
				return models.AirtelDisburseEnquiryResponse{}, err
			}
		}
		token = *c.token
	}

	req, err := createInternalRequest(countries.TANZANIA, c.conf.Environment, DisbursementEnquiry, token, nil, request.ID)
	if err != nil {
		return models.AirtelDisburseEnquiryResponse{}, err
	}

	res := new(models.AirtelDisburseEnquiryResponse)

	_, err = c.base.Do(ctx, "ussd push", req, res)
	if err != nil {
		return models.AirtelDisburseEnquiryResponse{}, err
	}
	return *res, nil
}

package airtel

import (
	"context"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"github.com/techcraftlabs/airtel/pkg/models"
)

type DisbursementService interface {
	Disburse(ctx context.Context, request models.AirtelDisburseRequest) (models.AirtelDisburseResponse, error)
	TransactionEnquiry(ctx context.Context, response models.AirtelDisburseEnquiryRequest) (models.AirtelDisburseEnquiryResponse, error)
}

func (c *Client) Disburse(ctx context.Context, request models.AirtelDisburseRequest) (models.AirtelDisburseResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil{
		return models.AirtelDisburseResponse{}, err
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
	token, err := c.checkToken(ctx)
	if err != nil{
		return models.AirtelDisburseEnquiryResponse{}, err
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

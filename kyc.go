package airtel

import (
	"context"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"github.com/techcraftlabs/airtel/pkg/models"
)

type (
	KYCService interface {
		UserEnquiry(ctx context.Context, msisdn string) (models.AirtelUserEnquiryResponse, error)
	}
)

func (c *Client) UserEnquiry(ctx context.Context, msisdn string) (models.AirtelUserEnquiryResponse, error) {

	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelUserEnquiryResponse{}, err
	}

	req, err := createInternalRequest(countries.TANZANIA, c.Conf.Environment, UserEnquiry, token, nil, msisdn)
	if err != nil {
		return models.AirtelUserEnquiryResponse{}, err
	}

	res := new(models.AirtelUserEnquiryResponse)

	_, err = c.base.Do(ctx, "user enquiry", req, res)
	if err != nil {
		return models.AirtelUserEnquiryResponse{}, err
	}
	return *res, nil
}

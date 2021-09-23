package airtel

import (
	"context"
	"github.com/techcraftlabs/airtel/internal/models"
)

type (
	AccountService interface {
		Balance(ctx context.Context) (models.AirtelBalanceEnquiryResponse, error)
	}
)

func (c *Client) Balance(ctx context.Context) (models.AirtelBalanceEnquiryResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelBalanceEnquiryResponse{}, err
	}
	req, err := createInternalRequest("", c.Conf.Environment, AccountBalance, token, nil, "")
	if err != nil {
		return models.AirtelBalanceEnquiryResponse{}, err
	}

	res := new(models.AirtelBalanceEnquiryResponse)

	_, err = c.base.Do(ctx, "check balance", req, res)
	if err != nil {
		return models.AirtelBalanceEnquiryResponse{}, err
	}
	return *res, nil
}

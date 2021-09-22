package airtel

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel/pkg/models"
	"time"
)

type (
	AccountService interface {
		Balance(ctx context.Context) (models.AirtelBalanceEnquiryResponse, error)
	}
)

func (c *Client) Balance(ctx context.Context) (models.AirtelBalanceEnquiryResponse, error) {
	var token string
	if *c.token == "" {
		str, err := c.Token(ctx)
		if err != nil {
			return models.AirtelBalanceEnquiryResponse{}, err
		}
		token = fmt.Sprintf("%s", str.AccessToken)
	}
	//Add Auth Header
	if *c.token != "" {
		if !c.tokenExpiresAt.IsZero() && time.Until(c.tokenExpiresAt) < (60*time.Second) {
			if _, err := c.Token(ctx); err != nil {
				return models.AirtelBalanceEnquiryResponse{}, err
			}
		}
		token = *c.token
	}
	req, err := createInternalRequest("", c.conf.Environment, AccountBalance, token, nil, "")
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

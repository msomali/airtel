package airtel

import (
	"context"
	"github.com/techcraftlabs/airtel/pkg/models"
)

type(
	KYCService interface {
		UserEnquiry(ctx context.Context, msisdn string)(models.AirtelUserEnquiryResponse,error)
	}
)

func (c *Client) UserEnquiry(ctx context.Context, msisdn string) (models.AirtelUserEnquiryResponse, error) {
	panic("implement me")
}

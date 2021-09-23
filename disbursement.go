package airtel

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel/internal"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"github.com/techcraftlabs/airtel/pkg/models"
	"net/http"
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

	countryName := request.CountryOfTransaction
	country, err := countries.GetByName(countryName)
	if err != nil {
		return models.AirtelDisburseResponse{}, err
	}
	var opts []internal.RequestOption

	hs := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "*/*",
		"X-Country":    country.Code,
		"X-Currency":   country.CurrencyCode,
		"Authorization":        fmt.Sprintf("Bearer %s", token),  
	}

	headersOpt := internal.WithRequestHeaders(hs)
	opts = append(opts,headersOpt)
	reqUrl := requestURL(c.Conf.Environment,Disbursement)
	req := internal.NewRequest(http.MethodPost,reqUrl,request,opts...)
	res := new(models.AirtelDisburseResponse)
	_, err = c.base.Do(ctx, "disbursement", req, res)
	if err != nil {
		return models.AirtelDisburseResponse{}, err
	}
	return *res, nil
}

func (c *Client) TransactionEnquiry(ctx context.Context, request models.AirtelDisburseEnquiryRequest) (models.AirtelDisburseEnquiryResponse, error) {
	panic("")
}

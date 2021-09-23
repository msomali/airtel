package airtel

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel/internal"
	"github.com/techcraftlabs/airtel/internal/models"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"net/http"
)

type DisbursementService interface {
	Disburse(ctx context.Context, request models.AirtelDisburseRequest) (models.AirtelDisburseResponse, error)
	TransactionEnquiry(ctx context.Context, response models.AirtelDisburseEnquiryRequest) (models.AirtelDisburseEnquiryResponse, error)
}

func (c *Client) Disburse(ctx context.Context, request models.AirtelDisburseRequest) (models.AirtelDisburseResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelDisburseResponse{}, err
	}

	countryName := request.CountryOfTransaction
	country, err := countries.GetByName(countryName)
	if err != nil {
		return models.AirtelDisburseResponse{}, err
	}
	var opts []internal.RequestOption

	hs := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "*/*",
		"X-Country":     country.CodeName,
		"X-Currency":    country.CurrencyCode,
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	headersOpt := internal.WithRequestHeaders(hs)
	opts = append(opts, headersOpt)
	reqUrl := requestURL(c.Conf.Environment, Disbursement)
	req := internal.NewRequest(http.MethodPost, reqUrl, request, opts...)
	res := new(models.AirtelDisburseResponse)
	_, err = c.base.Do(ctx, "disbursement", req, res)
	if err != nil {
		return models.AirtelDisburseResponse{}, err
	}
	return *res, nil
}

func (c *Client) TransactionEnquiry(ctx context.Context, request models.AirtelDisburseEnquiryRequest) (models.AirtelDisburseEnquiryResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelDisburseEnquiryResponse{}, err
	}

	countryName := request.CountryOfTransaction
	country, err := countries.GetByName(countryName)
	if err != nil {
		return models.AirtelDisburseEnquiryResponse{}, err
	}
	var opts []internal.RequestOption

	hs := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "*/*",
		"X-Country":     country.CodeName,
		"X-Currency":    country.CurrencyCode,
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	headersOpt := internal.WithRequestHeaders(hs)
	endpointOption := internal.WithEndpoint(request.ID)
	opts = append(opts, headersOpt,endpointOption)
	reqUrl := requestURL(c.Conf.Environment, DisbursementEnquiry)
	req := internal.NewRequest(http.MethodGet, reqUrl, request, opts...)
	res := new(models.AirtelDisburseEnquiryResponse)
	_, err = c.base.Do(ctx, "disbursement enquiry", req, res)
	if err != nil {
		return models.AirtelDisburseEnquiryResponse{}, err
	}
	return *res, nil
}

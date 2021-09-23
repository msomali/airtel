package airtel

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel/internal"
	"github.com/techcraftlabs/airtel/internal/models"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"net/http"
)

type (
	KYCService interface {
		UserEnquiry(ctx context.Context, request models.AirtelUserEnquiryRequest) (models.AirtelUserEnquiryResponse, error)
	}
)

func (c *Client) UserEnquiry(ctx context.Context, request models.AirtelUserEnquiryRequest) (models.AirtelUserEnquiryResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelUserEnquiryResponse{}, err
	}

	countryName := request.CountryOfTransaction
	country, err := countries.GetByName(countryName)
	if err != nil {
		return models.AirtelUserEnquiryResponse{}, err
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
	endpointOption := internal.WithEndpoint(request.MSISDN)
	opts = append(opts, headersOpt,endpointOption)
	reqUrl := requestURL(c.Conf.Environment, UserEnquiry)
	req := internal.NewRequest(http.MethodGet, reqUrl, request, opts...)

	res := new(models.AirtelUserEnquiryResponse)
	_, err = c.base.Do(ctx, "user enquiry", req, res)
	if err != nil {
		return models.AirtelUserEnquiryResponse{}, err
	}
	return *res, nil
}

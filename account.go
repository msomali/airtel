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
	AccountService interface {
		Balance(ctx context.Context, request models.AirtelBalanceEnquiryRequest) (models.AirtelBalanceEnquiryResponse, error)
	}
)

func (c *Client) Balance(ctx context.Context,request models.AirtelBalanceEnquiryRequest) (models.AirtelBalanceEnquiryResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelBalanceEnquiryResponse{}, err
	}

	countryName := request.CountryOfTransaction
	country, err := countries.GetByName(countryName)
	if err != nil {
		return models.AirtelBalanceEnquiryResponse{}, err
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
	reqUrl := requestURL(c.Conf.Environment, BalanceEnquiry)
	req := internal.NewRequest(http.MethodGet, reqUrl, request, opts...)
	res := new(models.AirtelBalanceEnquiryResponse)
	_, err = c.base.Do(ctx, "balance enquiry", req, res)
	if err != nil {
		return models.AirtelBalanceEnquiryResponse{}, err
	}
	return *res, nil
}

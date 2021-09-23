package airtel

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel/internal"
	"github.com/techcraftlabs/airtel/internal/models"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"net/http"
)

type CollectionService interface {
	Push(ctx context.Context, request models.AirtelPushRequest) (models.AirtelPushResponse, error)
	Refund(ctx context.Context, request models.AirtelRefundRequest) (models.AirtelRefundResponse, error)
	Enquiry(ctx context.Context, request models.AirtelPushEnquiryRequest) (models.AirtelPushEnquiryResponse, error)
	CallbackServeHTTP(writer http.ResponseWriter, request *http.Request)
}

func (c *Client) Push(ctx context.Context, request models.AirtelPushRequest) (models.AirtelPushResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelPushResponse{}, err
	}

	transaction := request.Transaction
	countryCodeName := transaction.Country
	currencyCodeName := transaction.Currency

	var opts []internal.RequestOption

	hs := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "*/*",
		"X-Country":     countryCodeName,
		"X-Currency":    currencyCodeName,
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	headersOpt := internal.WithRequestHeaders(hs)
	opts = append(opts, headersOpt)

	reqUrl := requestURL(c.Conf.Environment, USSDPush)

	req := internal.NewRequest(http.MethodPost, reqUrl, request, opts...)

	if err != nil {
		return models.AirtelPushResponse{}, err
	}

	res := new(models.AirtelPushResponse)
	_, err = c.base.Do(ctx, "ussd push", req, res)
	if err != nil {
		return models.AirtelPushResponse{}, err
	}
	return *res, nil
}

func (c *Client) Refund(ctx context.Context, request models.AirtelRefundRequest) (models.AirtelRefundResponse, error) {
	country, err := countries.GetByName(request.CountryOfTransaction)
	if err != nil {
		return models.AirtelRefundResponse{}, err
	}
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelRefundResponse{}, err
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
	reqUrl := requestURL(c.Conf.Environment, Refund)

	req := internal.NewRequest(http.MethodPost, reqUrl, request, opts...)

	if err != nil {
		return models.AirtelRefundResponse{}, err
	}

	res := new(models.AirtelRefundResponse)
	_, err = c.base.Do(ctx, "refund", req, res)
	if err != nil {
		return models.AirtelRefundResponse{}, err
	}
	return *res, nil

}

func (c *Client) Enquiry(ctx context.Context, request models.AirtelPushEnquiryRequest) (models.AirtelPushEnquiryResponse, error) {

	country, err := countries.GetByName(request.CountryOfTransaction)
	if err != nil {
		return models.AirtelPushEnquiryResponse{}, err
	}
	token, err := c.checkToken(ctx)
	if err != nil {
		return models.AirtelPushEnquiryResponse{}, err
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
	endpointOpt := internal.WithEndpoint(request.ID)
	opts = append(opts, headersOpt, endpointOpt)

	reqUrl := requestURL(c.Conf.Environment, PushEnquiry)

	req := internal.NewRequest(http.MethodGet, reqUrl, request, opts...)

	if err != nil {
		return models.AirtelPushEnquiryResponse{}, err
	}

	res := new(models.AirtelPushEnquiryResponse)
	_, err = c.base.Do(ctx, "ussd push enquiry", req, res)
	if err != nil {
		return models.AirtelPushEnquiryResponse{}, err
	}
	return *res, nil
}

func (c *Client) CallbackServeHTTP(writer http.ResponseWriter, request *http.Request) {
	body := new(models.AirtelCallbackRequest)
	err := internal.ReceivePayload(request, body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	reqBody := *body

	//todo: check the hash if it is OK
	err = c.pushCallbackFunc.Handle(reqBody)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

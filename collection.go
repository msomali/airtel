package airtel

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel/internal"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"github.com/techcraftlabs/airtel/pkg/models"
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
	if err != nil{
		return models.AirtelPushResponse{}, err
	}
	countryName := request.Transaction.Country
	country, err := countries.GetByCode(countryName)
	if err != nil {
		return models.AirtelPushResponse{}, err
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

	reqUrl := requestURL(c.Conf.Environment,USSDPush)

	req := internal.NewRequest(http.MethodPost,reqUrl,request,opts...)

	if err != nil {
		return models.AirtelPushResponse{}, err
	}
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
	panic("")
}

func (c *Client) Enquiry(ctx context.Context, request models.AirtelPushEnquiryRequest) (models.AirtelPushEnquiryResponse, error) {
	panic("")
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

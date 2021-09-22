package airtel

import (
	"context"
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
	req, err := createInternalRequest(countries.TANZANIA, c.conf.Environment, USSDPush, token, request, "")
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
	token, err := c.checkToken(ctx)
	if err != nil{
		return models.AirtelRefundResponse{}, err
	}

	req, err := createInternalRequest(countries.TANZANIA, c.conf.Environment, Refund, token, request, "")
	if err != nil {
		return models.AirtelRefundResponse{}, err
	}

	res := new(models.AirtelRefundResponse)

	_, err = c.base.Do(ctx, "ussd push", req, res)
	if err != nil {
		return models.AirtelRefundResponse{}, err
	}
	return *res, nil
}

func (c *Client) Enquiry(ctx context.Context, request models.AirtelPushEnquiryRequest) (models.AirtelPushEnquiryResponse, error) {
	token, err := c.checkToken(ctx)
	if err != nil{
		return models.AirtelPushEnquiryResponse{}, err
	}

	req, err := createInternalRequest(countries.TANZANIA, c.conf.Environment, PushEnquiry, token, nil, request.ID)
	if err != nil {
		return models.AirtelPushEnquiryResponse{}, err
	}

	res := new(models.AirtelPushEnquiryResponse)

	_, err = c.base.Do(ctx, "ussd push", req, res)
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
	err = c.pushCallbackFunc(reqBody)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

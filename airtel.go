package airtel

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/techcraftlabs/airtel/internal"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"github.com/techcraftlabs/airtel/pkg/models"
	"net/http"
	"time"
)

var _ CollectionService = (*Client)(nil)
var _ Authenticator = (*Client)(nil)
var _ DisbursementService = (*Client)(nil)

const (
	PRODUCTION                 Environment = "production"
	STAGING                    Environment = "staging"
	BaseURLProduction                      = "https://openapi.airtel.africa"
	BaseURLStaging                         = "https://openapiuat.airtel.africa"
	AuthEndpoint                           = "/auth/oauth2/token"
	PushEndpoint                           = "/merchant/v1/payments/"
	RefundEndpoint                         = "/standard/v1/payments/refund"
	PushEnquiryEndpoint                    = "/standard/v1/payments/"
	DisbursmentEndpoint                    = "/standard/v1/disbursements/"
	DisbursmentEnquiryEndpoint             = "/standard/v1/disbursements/"

	defaultGrantType = "client_credentials"
)

const (
	Authorization RequestType = iota
	USSDPush
	Refund
	PushEnquiry
	Disbursment
	DisbursmentEnquiry
)

type (
	Environment string
	RequestType uint
	Config      struct {
		PublicKey   string
		Environment Environment
		ClientID    string
		Secret      string
	}

	Client struct {
		baseURL        string
		conf           *Config
		base           *internal.BaseClient
		token          *string
		tokenExpiresAt time.Time
		callbacker func(request models.AirtelCallbackRequest)error
	}

	Request struct {
	}

	Authenticator interface {
		Token(ctx context.Context) (models.AirtelAuthResponse, error)
	}

	CollectionService interface {
		Push(ctx context.Context, request models.AirtelPushRequest) (models.AirtelPushResponse, error)
		Refund(ctx context.Context, request models.AirtelRefundRequest)(models.AirtelRefundResponse,error)
		Enquiry(ctx context.Context, request models.AirtelPushEnquiryRequest)(models.AirtelPushEnquiryResponse,error)
		CallbackServeHTTP(writer http.ResponseWriter, request *http.Request)
	}

	DisbursementService interface {
		Disburse(ctx context.Context, request models.AirtelDisburseRequest)(models.AirtelDisburseResponse,error)
		TransactionEnquiry(ctx context.Context, response models.AirtelDisburseEnquiryRequest)(models.AirtelDisburseEnquiryResponse,error)
	}
)

func (c *Client) Disburse(ctx context.Context, request models.AirtelDisburseRequest) (models.AirtelDisburseResponse, error) {
	panic("implement me")
}

func (c *Client) TransactionEnquiry(ctx context.Context, response models.AirtelDisburseEnquiryRequest) (models.AirtelDisburseEnquiryResponse, error) {
	panic("implement me")
}

func generateEncryptedKey(apiKey, pubKey string) (string, error) {

	decodedBase64, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return "", fmt.Errorf("could not decode pub key to Base64 string: %w", err)
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(decodedBase64)
	if err != nil {
		return "", fmt.Errorf("could not parse encoded public key (encryption key) : %w", err)
	}

	//check if the public key is RSA public key
	publicKey, isRSAPublicKey := publicKeyInterface.(*rsa.PublicKey)
	if !isRSAPublicKey {
		return "", fmt.Errorf("public key parsed is not an RSA public key : %w", err)
	}

	msg := []byte(apiKey)

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg)

	if err != nil {
		return "", fmt.Errorf("could not encrypt api key using generated public key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil

}

func (c *Client) Token(ctx context.Context) (models.AirtelAuthResponse, error) {
	body := models.AirtelAuthRequest{
		ClientID:     c.conf.ClientID,
		ClientSecret: c.conf.Secret,
		GrantType:    defaultGrantType,
	}
	req, err := createInternalRequest("", c.conf.Environment, Authorization, "", body, "")
	if err != nil {
		return models.AirtelAuthResponse{}, err
	}

	res := new(models.AirtelAuthResponse)

	_, err = c.base.Do(ctx, "Token", req, res)
	if err != nil {
		return models.AirtelAuthResponse{}, err
	}
	//fmt.Printf("status code: %v\nheaders: %v\npayload: %v\nerror: %v\n", do.StatusCode, do.Headers, do.Payload, do.Error)
	*c.token = res.AccessToken
	return *res, nil
}

func (c *Client) Push(ctx context.Context, request models.AirtelPushRequest) (models.AirtelPushResponse, error) {
	var token string
	if *c.token == "" {
		str, err := c.Token(ctx)
		if err != nil {
			return models.AirtelPushResponse{}, err
		}
		token = fmt.Sprintf("%s", str.AccessToken)
	}
	//Add Auth Header
	if *c.token != "" {
		if !c.tokenExpiresAt.IsZero() && time.Until(c.tokenExpiresAt) < (60*time.Second) {
			if _, err := c.Token(ctx); err != nil {
				return models.AirtelPushResponse{}, err
			}
		}
		token = *c.token
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

func (c *Client) Refund(ctx context.Context, request models.AirtelRefundRequest)(models.AirtelRefundResponse,error) {
	var token string
	if *c.token == "" {
		str, err := c.Token(ctx)
		if err != nil {
			return models.AirtelRefundResponse{}, err
		}
		token = fmt.Sprintf("%s", str.AccessToken)
	}
	//Add Auth Header
	if *c.token != "" {
		if !c.tokenExpiresAt.IsZero() && time.Until(c.tokenExpiresAt) < (60*time.Second) {
			if _, err := c.Token(ctx); err != nil {
				return models.AirtelRefundResponse{}, err
			}
		}
		token = *c.token
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

func (c *Client) Enquiry(ctx context.Context, request models.AirtelPushEnquiryRequest)(models.AirtelPushEnquiryResponse,error) {
	var token string
	if *c.token == "" {
		str, err := c.Token(ctx)
		if err != nil {
			return models.AirtelPushEnquiryResponse{}, err
		}
		token = fmt.Sprintf("%s", str.AccessToken)
	}
	//Add Auth Header
	if *c.token != "" {
		if !c.tokenExpiresAt.IsZero() && time.Until(c.tokenExpiresAt) < (60*time.Second) {
			if _, err := c.Token(ctx); err != nil {
				return models.AirtelPushEnquiryResponse{}, err
			}
		}
		token = *c.token
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
		http.Error(writer,err.Error(),http.StatusInternalServerError)
		return
	}
	reqBody := *body

	//todo: check the hash if it is OK
	err = c.callbacker(reqBody)
	if err != nil {
		http.Error(writer,err.Error(),http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func NewClient(config *Config, environment Environment, debugMode bool) *Client {
	token := new(string)
	base := internal.NewBaseClient(internal.WithDebugMode(debugMode))
	baseURL := new(string)
	switch environment {
	case STAGING:
		*baseURL = BaseURLStaging

	case PRODUCTION:
		*baseURL = BaseURLProduction

	default:
		*baseURL = BaseURLStaging
	}

	url := *baseURL
	return &Client{
		baseURL: url,
		conf:    config,
		base:    base,
		token:   token,
	}
}

func getRequestURL(env Environment, requestType RequestType, id ...string) string {

	switch requestType {
	case Authorization:
		if env == STAGING {
			return fmt.Sprintf("%s%s", BaseURLStaging, AuthEndpoint)
		}
		return fmt.Sprintf("%s%s", BaseURLProduction, AuthEndpoint)

	case USSDPush:
		if env == STAGING {
			return fmt.Sprintf("%s%s", BaseURLStaging, PushEndpoint)
		}
		return fmt.Sprintf("%s%s", BaseURLProduction, PushEndpoint)

	case Refund:
		if env == STAGING {
			return fmt.Sprintf("%s%s", BaseURLStaging, RefundEndpoint)
		}
		return fmt.Sprintf("%s%s", BaseURLProduction, RefundEndpoint)

	case PushEnquiry:
		if env == STAGING {
			return fmt.Sprintf("%s%s%s", BaseURLStaging, PushEnquiryEndpoint, id[0])
		}
		return fmt.Sprintf("%s%s%s", BaseURLProduction, PushEnquiryEndpoint, id[0])
	}

	return ""

}

func createInternalRequest(countryName string, env Environment, requestType RequestType, token string, body interface{}, id string) (*internal.Request, error) {
	var (
		country countries.Country
		err     error
	)

	if requestType != Authorization {
		country, err = countries.Get(countryName)
		if err != nil {
			return nil, err
		}
	}

	switch requestType {
	case Authorization:
		reqURL := getRequestURL(env, Authorization)
		hs := map[string]string{
			"Content-Type": "application/json",
			"Accept":       "*/*",
		}
		return internal.NewRequest(http.MethodPost, reqURL, internal.JsonPayload, body, internal.WithRequestHeaders(hs)), nil

	case USSDPush:
		fmt.Printf("case ussdpush: the token is %v\n",token)

		reqURL := getRequestURL(env, USSDPush)
		hs := map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "*/*",
			"X-Country":     country.Code,
			"X-Currency":    country.CurrencyCode,
			"Token": fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL, internal.JsonPayload, body, internal.WithRequestHeaders(hs)), nil

	case Refund:

		reqURL := getRequestURL(env, USSDPush)
		hs := map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "*/*",
			"X-Country":     country.Code,
			"X-Currency":    country.CurrencyCode,
			"Token": fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL, internal.JsonPayload, body, internal.WithRequestHeaders(hs)), nil

	case PushEnquiry:
		reqURL := getRequestURL(env, PushEnquiry, id)
		hs := map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "*/*",
			"X-Country":     country.Code,
			"X-Currency":    country.CurrencyCode,
			"Token": fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL, internal.JsonPayload, body, internal.WithRequestHeaders(hs)), nil
	}

	return nil, nil
}

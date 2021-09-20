package airtel

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel/internal"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"net/http"
)

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
		Environment Environment
		ClientID    string
		Secret      string
	}

	Client struct {
		conf *Config
		base *internal.BaseClient
	}

	Request struct {
	}

	AuthZRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		GrantType    string `json:"grant_type"`
	}

	AuthZResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	Service interface {
		Authorization(ctx context.Context, request AuthZRequest) (AuthZResponse, error)
		Push()
		Refund()
		Enquiry()
		Callback()
		Disburse()
	}
)

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
			return fmt.Sprintf("%s%s/%s", BaseURLStaging, PushEnquiryEndpoint, id)
		}
		return fmt.Sprintf("%s%s/%s", BaseURLProduction, PushEnquiryEndpoint, id)
	}

	return ""

}

func createInternalRequest(countryName string, env Environment, requestType RequestType, token string, body interface{}, id string) (*internal.Request, error) {

	country, err := countries.Get(countryName)
	if err != nil {
		return nil, err
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
		reqURL := getRequestURL(env, USSDPush)
		hs := map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "*/*",
			"X-Country":     country.Code,
			"X-Currency":    country.CurrencyCode,
			"Authorization": fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL, internal.JsonPayload, body, internal.WithRequestHeaders(hs)), nil

	case Refund:

		reqURL := getRequestURL(env, USSDPush)
		hs := map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "*/*",
			"X-Country":     country.Code,
			"X-Currency":    country.CurrencyCode,
			"Authorization": fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL, internal.JsonPayload, body, internal.WithRequestHeaders(hs)), nil

	case PushEnquiry:
		reqURL := getRequestURL(env, PushEnquiry, id)
		hs := map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "*/*",
			"X-Country":     country.Code,
			"X-Currency":    country.CurrencyCode,
			"Authorization": fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL, internal.JsonPayload, body, internal.WithRequestHeaders(hs)), nil
	}

	return nil, nil
}

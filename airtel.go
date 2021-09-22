package airtel

import (
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
		tokenExpiresAt   time.Time
		pushCallbackFunc PushCallbackFunc
	}

	PushCallbackFunc func(request models.AirtelCallbackRequest)error

	Request struct {
	}
)

func NewClient(config *Config, pushCallbackFunc PushCallbackFunc, debugMode bool) *Client {
	token := new(string)
	base := internal.NewBaseClient(internal.WithDebugMode(debugMode))
	baseURL := new(string)
	switch config.Environment {
	case STAGING:
		*baseURL = BaseURLStaging

	case PRODUCTION:
		*baseURL = BaseURLProduction

	default:
		*baseURL = BaseURLStaging
	}

	url := *baseURL
	return &Client{
		baseURL:          url,
		conf:             config,
		base:             base,
		token:            token,
		pushCallbackFunc: pushCallbackFunc,
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

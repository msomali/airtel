package airtel

import (
	"fmt"
	"github.com/techcraftlabs/airtel/internal"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"github.com/techcraftlabs/airtel/pkg/models"
	"net/http"
	"strings"
	"time"
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

	defaultGrantType = "client_credentials"
)

const (
	Authorization RequestType = iota
	USSDPush
	Refund
	PushEnquiry
	PushCallback
	Disbursement
	AccountBalance
	DisbursementEnquiry
	TransactionSummary
	UserEnquiry
)

const (
	CollectionAPIName   = "collection"
	DisbursementAPIName = "disbursement"
	AccountAPIName      = "account"
	KYCAPIName          = "kyc"
	TransactionAPIName  = "transaction"
)

type (
	Environment string
	RequestType uint
	Config      struct {
		AllowedCountries     map[string][]string
		DisbursePIN          string
		CallbackPrivateKey   string
		CallbackAuth         bool
		PublicKey            string
		Environment          Environment
		ClientID             string
		Secret               string
	}

	Client struct {
		baseURL          string
		conf             *Config
		base             *internal.BaseClient
		token            *string
		tokenExpiresAt   time.Time
		pushCallbackFunc PushCallbackFunc
	}

	PushCallbackFunc func(request models.AirtelCallbackRequest) error

	Request struct {
	}
)

func (config *Config)SetAllowedCountries(apiName string, countries []string)  {
	if config.AllowedCountries == nil{
		m := make(map[string][]string)
		config.AllowedCountries = m
	}

	config.AllowedCountries[apiName] = countries
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if strings.ToLower(v) == strings.ToLower(str) {
			return true
		}
	}

	return false
}

func checkIfCountryIsAllowed(api string, country string, allCountries map[string][]string) bool {
	a := allCountries[api]
	if a == nil || len(a) <= 0 {
		return false
	}

	return contains(a, country)
}

func NewClient(config *Config, pushCallbackFunc PushCallbackFunc, debugMode bool) *Client {
	if config.AllowedCountries == nil{
		m := make(map[string][]string)
		config.AllowedCountries = m
		config.SetAllowedCountries(CollectionAPIName,[]string{"Tanzania"})
		config.SetAllowedCountries(DisbursementAPIName,[]string{"Tanzania"})
		config.SetAllowedCountries(AccountAPIName,[]string{"Tanzania"})
		config.SetAllowedCountries(KYCAPIName,[]string{"Tanzania"})
		config.SetAllowedCountries(TransactionAPIName,[]string{"Tanzania"})

	}
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
		fmt.Printf("case ussdpush: the token is %v\n", token)

		reqURL := getRequestURL(env, USSDPush)
		hs := map[string]string{
			"Content-Type": "application/json",
			"Accept":       "*/*",
			"X-Country":    country.Code,
			"X-Currency":   country.CurrencyCode,
			"Token":        fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL, internal.JsonPayload, body, internal.WithRequestHeaders(hs)), nil

	case Refund:

		reqURL := getRequestURL(env, USSDPush)
		hs := map[string]string{
			"Content-Type": "application/json",
			"Accept":       "*/*",
			"X-Country":    country.Code,
			"X-Currency":   country.CurrencyCode,
			"Token":        fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL, internal.JsonPayload, body, internal.WithRequestHeaders(hs)), nil

	case PushEnquiry:
		reqURL := getRequestURL(env, PushEnquiry, id)
		hs := map[string]string{
			"Content-Type": "application/json",
			"Accept":       "*/*",
			"X-Country":    country.Code,
			"X-Currency":   country.CurrencyCode,
			"Token":        fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL, internal.JsonPayload, body, internal.WithRequestHeaders(hs)), nil

	case AccountBalance:
		return nil, err

	case Disbursement:
		return nil, err

	case DisbursementEnquiry:
		return nil, err
	}

	return nil, nil
}

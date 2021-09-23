package airtel

import (
	"github.com/techcraftlabs/airtel/internal"
	"github.com/techcraftlabs/airtel/pkg/models"
	"strings"
	"time"
)

const (
	PRODUCTION                 Environment = "production"
	STAGING                    Environment = "staging"
	BaseURLProduction                      = "https://openapi.airtel.africa"
	BaseURLStaging                         = "https://openapiuat.airtel.africa"
)

const (
	TransactionAPIName  = "transaction"
)

type (
	Environment string

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
	env := config.Environment
	switch env{
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

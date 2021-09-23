package airtel

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel/internal"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"net/http"
)

const (
	AuthEndpoint               = "/auth/oauth2/token"
	PushEndpoint               = "/merchant/v1/payments/"
	RefundEndpoint             = "/standard/v1/payments/refund"
	PushEnquiryEndpoint        = "/standard/v1/payments/"
	DisbursementEndpoint       = "/standard/v1/disbursements/"
	DisbursmentEnquiryEndpoint = "/standard/v1/disbursements/"
	defaultGrantType           = "client_credentials"
	CollectionAPIName          = "collection"
	DisbursementAPIName        = "disbursement"
	AccountAPIName             = "account"
	KYCAPIName                 = "kyc"
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

type (
	RequestType uint
)

func (c *Client) request(ctx context.Context,requestType RequestType, body interface{}, opts... internal.RequestOption)  (*internal.Request,error){

	reqUrl := requestURL(c.conf.Environment,requestType)

	switch requestType {
	case USSDPush:

		return internal.NewRequest(http.MethodPost,reqUrl,body, opts...), nil

	default:
		return nil, nil
	}
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
		reqURL := requestURL(env, Authorization)
		hs := map[string]string{
			"Content-Type": "application/json",
			"Accept":       "*/*",
		}
		return internal.NewRequest(http.MethodPost, reqURL,   body, internal.WithRequestHeaders(hs)), nil

	case USSDPush:
		fmt.Printf("case ussdpush: the token is %v\n", token)

		reqURL := requestURL(env, USSDPush)
		hs := map[string]string{
			"Content-Type": "application/json",
			"Accept":       "*/*",
			"X-Country":    country.Code,
			"X-Currency":   country.CurrencyCode,
			"Token":        fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL,   body, internal.WithRequestHeaders(hs)), nil

	case Refund:

		reqURL := requestURL(env, USSDPush)
		hs := map[string]string{
			"Content-Type": "application/json",
			"Accept":       "*/*",
			"X-Country":    country.Code,
			"X-Currency":   country.CurrencyCode,
			"Token":        fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL,   body, internal.WithRequestHeaders(hs)), nil

	case PushEnquiry:
		reqURL := requestURL(env, PushEnquiry)
		hs := map[string]string{
			"Content-Type": "application/json",
			"Accept":       "*/*",
			"X-Country":    country.Code,
			"X-Currency":   country.CurrencyCode,
			"Token":        fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL,   body, internal.WithRequestHeaders(hs)), nil

	case AccountBalance:
		return nil, err

	case Disbursement:
		return nil, err

	case DisbursementEnquiry:
		return nil, err
	}

	return nil, nil
}


func requestURL(env Environment, requestType RequestType) string {

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
			return fmt.Sprintf("%s%s", BaseURLStaging, PushEnquiryEndpoint)
		}
		return fmt.Sprintf("%s%s", BaseURLProduction, PushEnquiryEndpoint)

	case Disbursement:
		if env == STAGING {
			return fmt.Sprintf("%s%s", BaseURLStaging, DisbursementEndpoint)
		}
		return fmt.Sprintf("%s%s", BaseURLProduction, DisbursementEndpoint)
	}
	return ""

}

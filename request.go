/*
 * MIT License
 *
 * Copyright (c) 2021 TECHCRAFT TECHNOLOGIES CO LTD.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

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
	TransactionSummaryEndpoint = "/merchant/v1/transactions"
	BalanceEnquiryEndpoint     = "/standard/v1/users/balance"
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
	BalanceEnquiry
	DisbursementEnquiry
	TransactionSummary
	UserEnquiry
)

type (
	RequestType uint
)

func (c *Client) request(ctx context.Context, requestType RequestType, body interface{}, opts ...internal.RequestOption) (*internal.Request, error) {

	reqUrl := requestURL(c.Conf.Environment, requestType)

	switch requestType {
	case USSDPush:

		return internal.NewRequest(http.MethodPost, reqUrl, body, opts...), nil

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
		country, err = countries.GetByName(countryName)
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
		return internal.NewRequest(http.MethodPost, reqURL, body, internal.WithRequestHeaders(hs)), nil

	case USSDPush:
		fmt.Printf("case ussdpush: the token is %v\n", token)

		reqURL := requestURL(env, USSDPush)
		hs := map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "*/*",
			"X-Country":     country.CodeName,
			"X-Currency":    country.CurrencyCode,
			"Authorization": fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL, body, internal.WithRequestHeaders(hs)), nil

	case Refund:

		reqURL := requestURL(env, USSDPush)
		hs := map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "*/*",
			"X-Country":     country.CodeName,
			"X-Currency":    country.CurrencyCode,
			"Authorization": fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL, body, internal.WithRequestHeaders(hs)), nil

	case PushEnquiry:
		reqURL := requestURL(env, PushEnquiry)
		hs := map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "*/*",
			"X-Country":     country.CodeName,
			"X-Currency":    country.CurrencyCode,
			"Authorization": fmt.Sprintf("Bearer %s", token),
		}

		return internal.NewRequest(http.MethodPost, reqURL, body, internal.WithRequestHeaders(hs)), nil

	case BalanceEnquiry:
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

	case TransactionSummary:
		if env == STAGING {
			return fmt.Sprintf("%s%s", BaseURLStaging, TransactionSummaryEndpoint)
		}
		return fmt.Sprintf("%s%s", BaseURLProduction, TransactionSummaryEndpoint)

	case BalanceEnquiry:
		if env == STAGING {
			return fmt.Sprintf("%s%s", BaseURLStaging, BalanceEnquiryEndpoint)
		}
		return fmt.Sprintf("%s%s", BaseURLProduction, BalanceEnquiryEndpoint)

	}
	return ""

}

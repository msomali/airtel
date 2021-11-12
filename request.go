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
	"fmt"
	"github.com/techcraftlabs/base"
	"net/http"
	"strings"
)

const (
	defaultGrantType     = "client_credentials"
	AuthApiGroup         = "authorization"
	CollectionApiGroup   = "collection"
	DisbursementApiGroup = "disbursement"
	AccountApiGroup      = "account"
	KycApiGroup          = "kyc"
	TransactionApiGroup  = "transaction"
)

const (
	Authorization RequestType = iota
	UssdPush
	Refund
	PushEnquiry
	PushCallback
	Disbursement
	DisbursementEnquiry
	BalanceEnquiry
	TransactionSummary
	UserEnquiry
)

func (t RequestType) httpMethod() string {
	switch t {
	case Authorization, UssdPush, Refund, PushCallback, Disbursement:
		return http.MethodPost

	case PushEnquiry, DisbursementEnquiry, UserEnquiry, BalanceEnquiry,
		TransactionSummary:
		return http.MethodGet

	default:
		return ""
	}
}

func (t RequestType) name() string {
	return []string{"authorization", "ussd push", "refund", "push enquiry", "push callback",
		"disbursement", "disbursement enquiry", "balance enquiry", "transaction summary",
		"user enquiry"}[t]
}

func (t RequestType) Group() string {
	switch t {
	case Authorization:
		return AuthApiGroup

	case PushCallback, Refund, PushEnquiry, UssdPush:
		return CollectionApiGroup

	case Disbursement, DisbursementEnquiry:
		return DisbursementApiGroup

	case BalanceEnquiry:
		return AccountApiGroup

	case UserEnquiry:
		return KycApiGroup

	case TransactionSummary:
		return TransactionApiGroup

	default:
		return "unknown/unsupported api group"
	}
}

func (t RequestType) endpoint(es Endpoints) string {
	switch t {
	case Authorization:
		return es.AuthEndpoint

	case UssdPush:
		return es.PushEndpoint

	case PushEnquiry:
		return es.PushEnquiryEndpoint

	case Refund:
		return es.RefundEndpoint

	case Disbursement:
		return es.DisbursementEndpoint

	case DisbursementEnquiry:
		return es.DisbursementEnquiryEndpoint

	case UserEnquiry:
		return es.UserEnquiryEndpoint

	case BalanceEnquiry:
		return es.BalanceEnquiryEndpoint

	case TransactionSummary:
		return es.TransactionSummaryEndpoint

	default:
		return ""
	}
}

type (
	RequestType uint
	Endpoints   struct {
		AuthEndpoint                string
		PushEndpoint                string
		RefundEndpoint              string
		PushEnquiryEndpoint         string
		DisbursementEndpoint        string
		DisbursementEnquiryEndpoint string
		TransactionSummaryEndpoint  string
		BalanceEnquiryEndpoint      string
		UserEnquiryEndpoint         string
	}
)

func (c *Client) makeInternalRequest(requestType RequestType, payload interface{}, opts ...base.RequestOption) *base.Request {
	endpoints := c.endpoints
	edps := *endpoints
	url := appendEndpoint(c.baseURL, requestType.endpoint(edps))
	method := requestType.httpMethod()
	return base.NewRequest(method, url, payload, opts...)
}

func appendEndpoint(url string, endpoint string) string {
	url, endpoint = strings.TrimSpace(url), strings.TrimSpace(endpoint)
	urlHasSuffix, endpointHasPrefix := strings.HasSuffix(url, "/"), strings.HasPrefix(endpoint, "/")

	bothTrue := urlHasSuffix == true && endpointHasPrefix == true
	bothFalse := urlHasSuffix == false && endpointHasPrefix == false
	notEqual := urlHasSuffix != endpointHasPrefix

	if notEqual {
		return fmt.Sprintf("%s%s", url, endpoint)
	}

	if bothFalse {
		return fmt.Sprintf("%s/%s", url, endpoint)
	}

	if bothTrue {
		endp := strings.TrimPrefix(endpoint, "/")
		return fmt.Sprintf("%s%s", url, endp)
	}

	return ""
}

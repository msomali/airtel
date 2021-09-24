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

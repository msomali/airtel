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

type (
	PushPayRequest struct {
		Reference          string
		SubscriberCountry  string
		SubscriberMsisdn   string
		TransactionAmount  int64
		TransactionCountry string
		TransactionID      string
	}

	PushPayResponse struct {
		ID               string `json:"id,omitempty"`
		Status           string `json:"status,omitempty"`
		ResultCode       string `json:"result_code,omitempty"`
		Success          bool   `json:"success,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
		Error            string `json:"error,omitempty"`
		StatusMessage    string `json:"status_message,omitempty"`
		StatusCode       string `json:"status_code,omitempty"`
	}
	DisburseRequest struct {
		ID                   string
		MSISDN               string
		Amount               int64
		Reference            string
		CountryOfTransaction string
	}

	DisburseResponse struct {
		ID               string `json:"id,omitempty"`
		Reference        string `json:"reference,omitempty"`
		AirtelMoneyID    string `json:"airtel_money_id,omitempty"`
		ResultCode       string `json:"result_code,omitempty"`
		Success          bool   `json:"success,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
		Error            string `json:"error,omitempty"`
		StatusMessage    string `json:"status_message,omitempty"`
		StatusCode       string `json:"status_code,omitempty"`
	}

	Service interface {
		Authenticator
		CollectionService
		DisbursementService
		AccountService
		TransactionService
		KYCService
	}
)

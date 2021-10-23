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

package models

type TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type PushRequest struct {
	Reference  string `json:"reference"`
	Subscriber struct {
		Country  string `json:"country"`
		Currency string `json:"currency"`
		Msisdn   string `json:"msisdn"`
	} `json:"subscriber"`
	Transaction struct {
		Amount   int64  `json:"amount"`
		Country  string `json:"country"`
		Currency string `json:"currency"`
		ID       string `json:"id"`
	} `json:"transaction"`
}

type RefundRequest struct {
	CountryOfTransaction string `json:"-"`
	Transaction          struct {
		AirtelMoneyID string `json:"airtel_money_id"`
	} `json:"transaction"`
}

type PushEnquiryRequest struct {
	ID                   string `json:"id"`
	CountryOfTransaction string `json:"country"`
}

type CallbackRequest struct {
	Transaction struct {
		ID            string `json:"id"`
		Message       string `json:"message"`
		StatusCode    string `json:"status_code"`
		AirtelMoneyID string `json:"airtel_money_id"`
	} `json:"transaction"`
	Hash string `json:"hash"`
}

type UserEnquiryRequest struct {
	MSISDN               string
	CountryOfTransaction string
}

type BalanceEnquiryRequest struct {
	MSISDN               string
	CountryOfTransaction string
}

type DisburseRequest struct {
	CountryOfTransaction string `json:"-"`
	Payee                struct {
		Msisdn string `json:"msisdn"`
	} `json:"payee"`
	Reference   string `json:"reference"`
	Pin         string `json:"pin"`
	Transaction struct {
		Amount int64  `json:"amount"`
		ID     string `json:"id"`
	} `json:"transaction"`
}

type DisburseEnquiryRequest struct {
	CountryOfTransaction string
	ID                   string `json:"id"`
}

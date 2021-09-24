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

package main

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/api/http"
	"github.com/techcraftlabs/airtel/internal/models"
	"time"
)

func callbacker() airtel.PushCallbackFunc {
	return func(request models.AirtelCallbackRequest) error {
		return nil
	}
}

func main() {

	pubKey := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCkq3XbDI1s8Lu7SpUBP+bqOs/MC6PKWz6n/0UkqTiOZqKqaoZClI3BUDTrSIJsrN1Qx7ivBzsaAYfsB0CygSSWay4iyUcnMVEDrNVOJwtWvHxpyWJC5RfKBrweW9b8klFa/CfKRtkK730apy0Kxjg+7fF0tB4O3Ic9Gxuv4pFkbQIDAQAB"
	//	pubKey2 := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCkq3XbDI1s8Lu7SpUBP+bqOs/MC6PKWz6n/0UkqTiOZqKqaoZClI3BUDTrSIJsrN1Qx7ivBzsaAYfsB0CygSSWay4iyUcnMVEDrNVOJwtWvHxpyWJC5RfKBrweW9b8klFa/CfKRtkK730apy0Kxjg+7fF0tB4O3Ic9Gxuv4pFkbQIDAQAB"
	config := &airtel.Config{
		AllowedCountries:   nil,
		DisbursePIN:        "4094",
		CallbackPrivateKey: "",
		CallbackAuth:       false,
		PublicKey:          pubKey,
		Environment:        airtel.STAGING,
		ClientID:           "747b6063-5eea-4464-b27c-a8f89c2e1fe3",
		Secret:             "9c8ded86-f45a-48f4-a9ee-8063cf8f43a0",
	}

	fn := callbacker()

	airtelClient := airtel.NewClient(config, fn, true)

	apiConfig := &http.Config{
		BaseURL:   "",
		Port:      0,
		DebugMode: false,
	}

	apiClient := http.NewClient(apiConfig, airtelClient)

	//req := api.PushPayRequest{
	//	Reference:          "this is a test",
	//	SubscriberCountry:  countries.TANZANIA,
	//	SubscriberMsisdn:   "1783839412",
	//	TransactionAmount:  6000000,
	//	TransactionCountry: countries.TANZANIA,
	//	TransactionID:      fmt.Sprintf("%v", time.Now().UnixNano()),
	//}
	//fmt.Printf("perform push pay\n")
	//pushPayResponse, err := apiClient.Push(context.TODO(), req)
	//if err != nil {
	//	fmt.Printf("error %v\n", err)
	//	return
	//}
	//
	//fmt.Printf("pushpay response: %v\n", pushPayResponse)

	//req2 := api.DisburseRequest{
	//	ID:                   fmt.Sprintf("%v", time.Now().UnixNano()),
	//	MSISDN:               "783839412",
	//	Amount:               3000,
	//	Reference:            "test request",
	//	CountryOfTransaction: countries.TANZANIA,
	//}
	//
	//fmt.Printf("Performing Disbursement")
	//disburseResponse, err := apiClient.Disburse(context.TODO(), req2)
	//if err != nil {
	//	fmt.Printf("error %v\n", err)
	//	return
	//}
	//
	//fmt.Printf("disburse response: %v\n", disburseResponse)
	params := airtel.Params{
		From:   time.Now().UnixNano(),
		To:     time.Now().UnixNano(),
		Limit:  4,
		Offset: 1,
	}
	resp, err := apiClient.Summary(context.TODO(), params)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return
	}

	fmt.Printf("summary response: %v\n", resp)

}

package main

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/api"
	"github.com/techcraftlabs/airtel/api/http"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"github.com/techcraftlabs/airtel/pkg/models"
	"time"
)


func callbacker()airtel.PushCallbackFunc{
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
	
	apiClient := http.NewClient(apiConfig,airtelClient)

	req := api.PushPayRequest{
		Reference:           "this is a test transaction",
		SubscriberCountry:   countries.TANZANIA_CODE,
		SubscriberCurrency:  countries.TANZANIA_CURRENCY_CODE,
		SubscriberMsisdn:    "765992153",
		TransactionAmount:   500,
		TransactionCountry:  countries.TANZANIA_CODE,
		TransactionCurrency: countries.TANZANIA_CURRENCY_CODE,
		TransactionID:       fmt.Sprintf("%v",time.Now().UnixNano()),
	}
	pushPayResponse, err := apiClient.Push(context.TODO(), req)
	if err != nil {
		return 
	}

	fmt.Printf("pushpay response: %v\n",pushPayResponse)

	req2 := api.DisburseRequest{
		ID:                   fmt.Sprintf("%v",time.Now().UnixNano()),
		MSISDN:               "765992153",
		Amount:               500,
		Reference:            "test request",
		CountryOfTransaction: countries.TANZANIA,
	}

	disburseResponse, err := apiClient.Disburse(context.TODO(),req2)
	if err != nil {
		return
	}

	fmt.Printf("disburse response: %v\n",disburseResponse)

}

package main

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/pkg/countries"
	"github.com/techcraftlabs/airtel/pkg/models"
	"os"
)

func main() {

	pubKey := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCkq3XbDI1s8Lu7SpUBP+bqOs/MC6PKWz6n/0UkqTiOZqKqaoZClI3BUDTrSIJsrN1Qx7ivBzsaAYfsB0CygSSWay4iyUcnMVEDrNVOJwtWvHxpyWJC5RfKBrweW9b8klFa/CfKRtkK730apy0Kxjg+7fF0tB4O3Ic9Gxuv4pFkbQIDAQAB"

	config := &airtel.Config{
		PublicKey:   pubKey,
		Environment: airtel.STAGING,
		ClientID:    "747b6063-5eea-4464-b27c-a8f89c2e1fe3",
		Secret:      "9c8ded86-f45a-48f4-a9ee-8063cf8f43a0",
	}

	fn := func(request models.AirtelCallbackRequest) error {
		return nil
	}

	client := airtel.NewClient(config, fn, true)

	response, err := client.Token(context.TODO())
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("response: %v\n", response)

	country, _ := countries.Get(countries.TANZANIA)
	pushRequest := models.AirtelPushRequest{
		Reference: "This is just a testing transaction",
		Subscriber: struct {
			Country  string `json:"country"`
			Currency string `json:"currency"`
			Msisdn   string `json:"msisdn"`
		}{
			Country:  country.Name,
			Currency: country.Currency,
			Msisdn:   "784956141",
		},
		Transaction: struct {
			Amount   int    `json:"amount"`
			Country  string `json:"country"`
			Currency string `json:"currency"`
			ID       string `json:"id"`
		}{
			Amount:   500,
			Country:  country.Name,
			Currency: country.Currency,
			ID:       "somesuperrandomid19101jsjs",
		},
	}
	pushResponse, err := client.Push(context.TODO(), pushRequest)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("push response: %v\n", pushResponse)

	refundRequest := models.AirtelRefundRequest{
		Transaction: struct {
			AirtelMoneyID string `json:"airtel_money_id"`
		}(struct {
			AirtelMoneyID string
		}{AirtelMoneyID: "heyspecialid"}),
	}

	refundResponse, err := client.Refund(context.TODO(), refundRequest)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("refund response: %v\n", refundResponse)
}

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
	_ "github.com/joho/godotenv/autoload"

	"fmt"
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/cli"
	"github.com/techcraftlabs/airtel/internal/models"
	"github.com/techcraftlabs/airtel/pkg/env"
	"os"
	"strings"
)

const (
	envBaseURL                     = "AIRTEL_MONEY_BASE_URL"
	envPublicKey                   = "AIRTEL_MONEY_PUBKEY"
	envDisbursePin                 = "AIRTEL_MONEY_DISBURSE_PIN"
	envClientId                    = "AIRTEL_MONEY_CLIENT_ID"
	envClientSecret                = "AIRTEL_MONEY_CLIENT_SECRET"
	envDebugMode                   = "AIRTEL_MONEY_DEBUG_MODE"
	envDeploymentEnv               = "AIRTEL_MONEY_DEPLOYMENT"
	envCallbackAuth                = "AIRTEL_MONEY_CALLBACK_AUTH"
	envCallbackPrivKey             = "AIRTEL_MONEY_CALLBACK_PRIVKEY"
	envCountries                   = "AIRTEL_MONEY_COUNTRIES"
	envAuthEndpoint                = "AIRTEL_MONEY_AUTH_ENDPOINT"
	envPushEndpoint                = "AIRTEL_MONEY_PUSH_ENDPOINT"
	envRefundEndpoint              = "AIRTEL_MONEY_REFUND_ENDPOINT"
	envPushEnquiryEndpoint         = "AIRTEL_MONEY_PUSH_ENQUIRY_ENDPOINT"
	envDisbursementEndpoint        = "AIRTEL_MONEY_DISBURSE_ENDPOINT"
	envDisbursementEnquiryEndpoint = "AIRTEL_MONEY_DISBURSE_ENQUIRY_ENDPOINT"
	envTransactionSummaryEndpoint  = "AIRTEL_MONEY_SUMMARY_ENDPOINT"
	envBalanceEnquiryEndpoint      = "AIRTEL_MONEY_BALANCE_ENDPOINT"
	envUserEnquiryEndpoint         = "AIRTEL_MONEY_USER_ENDPOINT"
	defBaseURL                     = "https://openapi.airtel.africa/"
	defPublicKey                   = ""
	defDisbursePin                 = "4094"
	defClientId                    = "747b6063-5eea-4464-b27c-a8f89c2e1fe3"
	defClientSecret                = "9c8ded86-f45a-48f4-a9ee-8063cf8f43a0"
	defDebugMode                   = true
	defDeploymentEnv               = "staging"
	defCallbackAuth                = false
	defCallbackPrivKey             = "zITVAAGYSlzl1WkUQJn81kbpT5drH3koffT8jCkcJJA="
	defCountries                   = "tanzania"
	defAuthEndpoint                = "/auth/oauth2/token"
	defPushEndpoint                = "/merchant/v1/payments/"
	defRefundEndpoint              = "/standard/v1/payments/refund"
	defPushEnquiryEndpoint         = "/standard/v1/payments/"
	defDisbursementEndpoint        = "/standard/v1/disbursements/"
	defDisbursementEnquiryEndpoint = "/standard/v1/disbursements/"
	defTransactionSummaryEndpoint  = "/merchant/v1/transactions"
	defBalanceEnquiryEndpoint      = "/standard/v1/users/balance"
	defUserEnquiryEndpoint         = "/standard/v1/users/"
)

func callbacker() airtel.PushCallbackFunc {
	return func(request models.CallbackRequest) error {
		return nil
	}
}

func main() {

	var (
		baseURL         = env.String(envBaseURL, defBaseURL)
		pubKey          = env.String(envPublicKey, defPublicKey)
		disbursePin     = env.String(envDisbursePin, defDisbursePin)
		callbackPrivKey = env.String(envCallbackPrivKey, defCallbackPrivKey)
		clientID        = env.String(envClientId, defClientId)
		secret          = env.String(envClientSecret, defClientSecret)
		debugMode       = env.Bool(envDebugMode, defDebugMode)
		callbackAuth    = env.Bool(envCallbackAuth, defCallbackAuth)
		deployEnv       = strings.ToLower(env.String(envDeploymentEnv, defDeploymentEnv))
		countries       = strings.Split(env.String(envCountries, defCountries), " ")
		endpoints       = &airtel.Endpoints{
			AuthEndpoint:                env.String(envAuthEndpoint, defAuthEndpoint),
			PushEndpoint:                env.String(envPushEndpoint, defPushEndpoint),
			RefundEndpoint:              env.String(envRefundEndpoint, defRefundEndpoint),
			PushEnquiryEndpoint:         env.String(envPushEnquiryEndpoint, defPushEnquiryEndpoint),
			DisbursementEndpoint:        env.String(envDisbursementEndpoint, defDisbursementEndpoint),
			DisbursementEnquiryEndpoint: env.String(envDisbursementEnquiryEndpoint, defDisbursementEnquiryEndpoint),
			TransactionSummaryEndpoint:  env.String(envTransactionSummaryEndpoint, defTransactionSummaryEndpoint),
			BalanceEnquiryEndpoint:      env.String(envBalanceEnquiryEndpoint, defBalanceEnquiryEndpoint),
			UserEnquiryEndpoint:         env.String(envUserEnquiryEndpoint, defUserEnquiryEndpoint),
		}
	)

	config := &airtel.Config{
		BaseURL:   baseURL,
		Endpoints: endpoints,
		AllowedCountries: map[string][]string{
			airtel.TransactionApiGroup:  countries,
			airtel.CollectionApiGroup:   countries,
			airtel.DisbursementApiGroup: countries,
			airtel.AccountApiGroup:      countries,
			airtel.KycApiGroup:          countries,
		},
		DisbursePIN:        disbursePin,
		CallbackPrivateKey: callbackPrivKey,
		CallbackAuth:       callbackAuth,
		PublicKey:          pubKey,
		Environment:        airtel.Environment(deployEnv),
		ClientID:           clientID,
		Secret:             secret,
	}

	fn := callbacker()

	airtelClient := airtel.NewClient(config, fn, debugMode)

	app := cli.New(airtelClient)

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

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
	"github.com/techcraftlabs/airtel/internal"
	"github.com/techcraftlabs/airtel/internal/models"
	"time"
)

const (
	PRODUCTION        Environment = "production"
	STAGING           Environment = "staging"
	BaseURLProduction             = "https://openapi.airtel.africa"
	BaseURLStaging                = "https://openapiuat.airtel.africa"
)

type (
	Environment string

	Config struct {
		Endpoints          *Endpoints
		AllowedCountries   map[string][]string
		DisbursePIN        string
		CallbackPrivateKey string
		CallbackAuth       bool
		PublicKey          string
		Environment        Environment
		BaseURL            string
		ClientID           string
		Secret             string
	}

	Client struct {
		//baseURL          string
		Conf             *Config
		base             *internal.BaseClient
		token            *string
		tokenExpiresAt   time.Time
		pushCallbackFunc PushCallbackHandler
		reqAdapter       RequestAdapter
		resAdapter       ResponseAdapter
	}

	PushCallbackHandler interface {
		Handle(request models.CallbackRequest) error
	}
	PushCallbackFunc func(request models.CallbackRequest) error
)

func (pf PushCallbackFunc) Handle(request models.CallbackRequest) error {
	return pf(request)
}

func (config *Config) SetAllowedCountries(apiName string, countries []string) {
	if config.AllowedCountries == nil {
		m := make(map[string][]string)
		config.AllowedCountries = m
	}

	config.AllowedCountries[apiName] = countries
}

func NewClient(config *Config, pushCallbackFunc PushCallbackHandler, debugMode bool) *Client {
	if config.AllowedCountries == nil {
		m := make(map[string][]string)
		config.AllowedCountries = m
		config.SetAllowedCountries(CollectionApiGroup, []string{"Tanzania"})
		config.SetAllowedCountries(DisbursementApiGroup, []string{"Tanzania"})
		config.SetAllowedCountries(AccountApiGroup, []string{"Tanzania"})
		config.SetAllowedCountries(KycApiGroup, []string{"Tanzania"})
		config.SetAllowedCountries(TransactionApiGroup, []string{"Tanzania"})

	}
	token := new(string)
	base := internal.NewBaseClient(internal.WithDebugMode(debugMode))

	return &Client{
		Conf:             config,
		base:             base,
		token:            token,
		resAdapter:       &ResAdapter{},
		reqAdapter:       &ReqAdapter{Conf: config},
		pushCallbackFunc: pushCallbackFunc,
	}
}

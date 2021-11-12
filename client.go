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
	"github.com/techcraftlabs/base"
	"time"
)

const (
	PRODUCTION     Environment = "production"
	STAGING        Environment = "staging"
	prodBaseURL                = "https://openapi.airtel.africa"
	stagingBaseURL             = "https://openapiuat.airtel.africa"
)

type (
	Environment string

	Config struct {
		Endpoints          *Endpoints
		AllowedCountries   map[ApiGroup][]string
		DisbursePIN        string
		CallbackPrivateKey string
		CallbackAuth       bool
		PublicKey          string
		Environment        Environment
		ClientID           string
		Secret             string
	}

	Client struct {
		baseURL           string
		rv                base.Receiver
		rp                base.Replier
		Conf              *Config
		base              *base.Client
		token             *string
		tokenExpiresAt    time.Time
		pushCallbackFunc  PushCallbackHandler
		disburseAdapter   DisbursementAdapter
		collectionAdapter CollectionAdapter
	}

	PushCallbackHandler interface {
		Handle(request InternalCallbackRequest) error
	}
	PushCallbackFunc func(request InternalCallbackRequest) error
)

func (pf PushCallbackFunc) Handle(request InternalCallbackRequest) error {
	return pf(request)
}

func (config *Config) SetAllowedCountries(apiName ApiGroup, countries []string) {
	if config.AllowedCountries == nil {
		m := make(map[ApiGroup][]string)
		config.AllowedCountries = m
	}

	config.AllowedCountries[apiName] = countries
}

func (c *Client) SetCollectionAdapter(adapter CollectionAdapter) {
	c.collectionAdapter = adapter
}

func (c *Client) SetDisburseAdapter(adapter DisbursementAdapter) {
	c.disburseAdapter = adapter
}

func NewClient(config *Config, pushCallbackFunc PushCallbackHandler, debugMode bool) *Client {
	var (
		baseURL string
	)

	if config.Environment == STAGING {
		baseURL = stagingBaseURL
	}

	if config.Environment == PRODUCTION {
		baseURL = prodBaseURL
	}

	if config.AllowedCountries == nil {
		m := make(map[ApiGroup][]string)
		config.AllowedCountries = m
		config.SetAllowedCountries(Collection, []string{"Tanzania"})
		config.SetAllowedCountries(Disburse, []string{"Tanzania"})
		config.SetAllowedCountries(Account, []string{"Tanzania"})
		config.SetAllowedCountries(KYC, []string{"Tanzania"})
		config.SetAllowedCountries(Transaction, []string{"Tanzania"})

	}
	token := new(string)
	newClient := base.NewClient(base.WithDebugMode(debugMode))

	c := &Client{
		baseURL: baseURL,
		Conf:    config,
		base:    newClient,
		token:   token,
		disburseAdapter: &disburseAdapter{
			Conf: config,
		},
		pushCallbackFunc:  pushCallbackFunc,
		collectionAdapter: &collectAdapter{},
	}
	logger := c.base.Logger
	dm := c.base.DebugMode

	rv := base.NewReceiver(logger, dm)
	rp := base.NewReplier(logger, dm)
	c.rv = rv
	c.rp = rp
	return c
}

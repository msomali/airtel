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
	"context"
	"github.com/techcraftlabs/airtel/internal"
	"net/http"
	"testing"
)

func TestConfig_SetAllowedCountries(t *testing.T) {
	type fields struct {
		AllowedCountries   map[string][]string
		DisbursePIN        string
		CallbackPrivateKey string
		CallbackAuth       bool
		PublicKey          string
		Environment        Environment
		ClientID           string
		Secret             string
	}
	type args struct {
		apiName   string
		countries []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "",
			fields: fields{
				AllowedCountries:   nil,
				DisbursePIN:        "",
				CallbackPrivateKey: "",
				CallbackAuth:       false,
				PublicKey:          "",
				Environment:        "",
				ClientID:           "",
				Secret:             "",
			},
			args: args{
				apiName:   CollectionApiGroup,
				countries: []string{"Kenya", "Uganda"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				AllowedCountries:   tt.fields.AllowedCountries,
				DisbursePIN:        tt.fields.DisbursePIN,
				CallbackPrivateKey: tt.fields.CallbackPrivateKey,
				CallbackAuth:       tt.fields.CallbackAuth,
				PublicKey:          tt.fields.PublicKey,
				Environment:        tt.fields.Environment,
				ClientID:           tt.fields.ClientID,
				Secret:             tt.fields.Secret,
			}

			config.SetAllowedCountries(tt.args.apiName, tt.args.countries)

			t.Logf("allowed countries for %s api are %v\n", tt.args.apiName, config.AllowedCountries[tt.args.apiName])
		})

	}
}

func TestWithEndpoint(t *testing.T) {
	type args struct {
		url      string
		endpoint string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal endpoint",
			args: args{
				url:      "www.server.com/users",
				endpoint: "0001",
			},
			want: "www.server.com/users/0001",
		},
		{
			name: "uel has /",
			args: args{
				url:      "www.server.com/users/",
				endpoint: "0001",
			},
			want: "www.server.com/users/0001",
		},
		{
			name: "endpoint has /",
			args: args{
				url:      "www.server.com/users",
				endpoint: "/0001",
			},
			want: "www.server.com/users/0001",
		},
		{
			name: "both have /",
			args: args{
				url:      "www.server.com/users/",
				endpoint: "/0001",
			},
			want: "www.server.com/users/0001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, endpoint := tt.args.url, tt.args.endpoint
			req := internal.NewRequest(http.MethodHead, url, nil, internal.WithEndpoint(endpoint))
			reqWithCtx, err := internal.NewRequestWithContext(context.TODO(), req)
			if err != nil {
				t.Errorf("%v\n", err)
			}
			got := reqWithCtx.URL.String()
			if got != tt.want {
				t.Errorf("WithEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSomeTypeImplementsSomeInterface(t *testing.T) {
	// won't compile if SomeType does not implement SomeInterface
	var (
		_ AccountService      = (*Client)(nil)
		_ CollectionService   = (*Client)(nil)
		_ Authenticator       = (*Client)(nil)
		_ DisbursementService = (*Client)(nil)
		_ TransactionService  = (*Client)(nil)
		_ KYCService          = (*Client)(nil)
	)
}

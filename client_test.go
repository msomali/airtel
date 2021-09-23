package airtel

import (
	"context"
	"github.com/techcraftlabs/airtel/internal"
	"net/http"
	"testing"
)

func Test_getRequestURL(t *testing.T) {
	type args struct {
		env         Environment
		requestType RequestType
		id          []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "authorization request url staging environment",
			args: args{
				env:         STAGING,
				requestType: Authorization,
				id:          []string{""},
			},
			want: "https://openapiuat.airtel.africa/auth/oauth2/token",
		},
		{
			name: "authorization request url production environment",
			args: args{
				env:         PRODUCTION,
				requestType: Authorization,
				id:          []string{""},
			},
			want: "https://openapi.airtel.africa/auth/oauth2/token",
		},
		{
			name: "staging url for push pay transaction enquiry",
			args: args{
				env:         STAGING,
				requestType: PushEnquiry,
				id:          []string{"ID001"},
			},
			want: "https://openapiuat.airtel.africa/standard/v1/payments/ID001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := requestURL(tt.args.env, tt.args.requestType, tt.args.id...); got != tt.want {
				t.Errorf("requestURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
				apiName:   CollectionAPIName,
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

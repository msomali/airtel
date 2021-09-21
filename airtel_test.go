package airtel

import "testing"

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
			if got := getRequestURL(tt.args.env, tt.args.requestType, tt.args.id...); got != tt.want {
				t.Errorf("getRequestURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

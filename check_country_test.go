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
	"testing"
)

func Test_CheckCountry(t *testing.T) {

	all := map[string][]string{
		"collection":   {"Tanzania", "Kenya", "Uganda", "Rwanda"},
		"disbursement": {"Tanzania"},
		"account":      {"Rwanda", "Uganda"},
	}
	type args struct {
		api          string
		country      string
		allCountries map[string][]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "check if not allowed is false",
			args: args{
				api:          Collection.String(),
				country:      "Burundi",
				allCountries: all,
			},
			want: false,
		},
		{
			name: "checking if it returns true when the country is allowed",
			args: args{
				api:          Collection.String(),
				country:      "tanzania",
				allCountries: all,
			},
			want: true,
		},
		{
			name: "passing empty collection name",
			args: args{
				api:          "",
				country:      "burundi",
				allCountries: all,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckCountry(tt.args.api, tt.args.country, tt.args.allCountries); got != tt.want {
				t.Errorf("CheckCountry() = %v, want %v", got, tt.want)
			}
		})
	}
}

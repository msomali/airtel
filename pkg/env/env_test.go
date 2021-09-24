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

package env_test

import (
	"github.com/techcraftlabs/airtel/pkg/env"
	"os"
	"reflect"
	"testing"
)

const (
	envTechcraftBaseName = "TECHCRAFT_BASENAME"
	defTechcraftBaseName = "default name"
	setTechcraftBaseName = "setenv name"

	envTechcraftBaseAge = "TECHCRAFT_BASEAGE"
	defTechcraftBaseAge = 10
	setTechcraftBaseAge = 50

	envTechcraftBaseSalary = "TECHCRAFT_BASE_SALARY"
	defTechcraftBaseSalary = 10.0000
	setTechcraftBaseSalary = 50.0000

	envTechcraftBaseStatus = "TECHCRAFT_BASE_STATUS"
	defTechcraftBaseStatus = true
	setTechcraftBaseStatus = true
)

func TestGet(t *testing.T) {
	type args struct {
		key          string
		valueInEnv   string
		defaultValue interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "test normal string",
			args: args{
				key:          envTechcraftBaseName,
				valueInEnv:   setTechcraftBaseName,
				defaultValue: defTechcraftBaseName,
			},
			want: setTechcraftBaseName,
		},
		{
			name: "testing integer vars",
			args: args{
				key:          envTechcraftBaseAge,
				valueInEnv:   "50",
				defaultValue: defTechcraftBaseAge,
			},
			want: setTechcraftBaseAge,
		},
		{
			name: "testing float vars",
			args: args{
				key:          envTechcraftBaseSalary,
				valueInEnv:   "50.0000",
				defaultValue: defTechcraftBaseSalary,
			},
			want: setTechcraftBaseSalary,
		},
		{
			name: "testing bool vars",
			args: args{
				key:          envTechcraftBaseStatus,
				valueInEnv:   "true",
				defaultValue: defTechcraftBaseStatus,
			},
			want: setTechcraftBaseStatus,
		},
		{
			name: "testing bool vars",
			args: args{
				key:          envTechcraftBaseStatus,
				valueInEnv:   "",
				defaultValue: defTechcraftBaseStatus,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		key, valueInEnv, defValue, want := tt.args.key, tt.args.valueInEnv, tt.args.defaultValue, tt.want
		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv(key, valueInEnv)
			if err != nil {
				t.Errorf("error while set the env %s to %s: %v\n", key, valueInEnv, err)
			}

			defer func(key, value string) {
				err := os.Setenv(key, value)
				if err != nil {
					t.Logf("error in usetting the env %s to %s,%v\n", key, value, err)
				}
			}(key, "")

			loaded := os.Getenv(key)

			t.Logf("loaded \"%s\" from os.Getenv is \"%s\"", key, loaded)

			got := env.Get(key, defValue)

			if !reflect.DeepEqual(got, want) {
				t.Errorf("\nGet() = \"%v\"\nwant = \"%v\"\n", got, tt.want)
			}
		})
	}
}

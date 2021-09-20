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

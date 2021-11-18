# airtel

specification: https://developer.airtel.africa/documentation
- [usage](#usage)
- [coverage](#coverage)
    - [client](#client)
    - [authorization](#authorization)
    - [collection](#collection)
    - [disbursement](#disbursement)
    - [account](#account)
    - [transaction](#transaction)
    - [KYC](#KYC)
    

## usage

    $ go get github.com/airteldata/airtel

```go

package main

import (
	"github.com/techcraftlabs/airtel"
)

```

## client

This is the client for the airtel api, that will be used to make the calls to the server



```go
package main

import (
	"context"
	"github.com/techcraftlabs/airtel"
)

func main() {
	pushCallbackFunc := func(ctx context.Context)airtel.PushCallbackFunc{
		return func(callbackFunc airtel.CallbackRequest) error {
			return nil
        }
    }
	config := &airtel.Config{
		AllowedCountries:   nil,
		DisbursePIN:        "",
		CallbackPrivateKey: "",
		CallbackAuth:       false,
		PublicKey:          "",
		Environment:        airtel.PRODUCTION, //airtel.STAGING,
		ClientID:           "", //given by airtel
		Secret:             "", //given by airtel
	}
	client := airtel.NewClient(config, pushCallbackFunc(context.TODO()), true)
}


```
## authorization

ask for token that the client will use to make subsequent calls to the server

```go
package main

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel"
)

func main() {
	pushCallbackFunc := func(ctx context.Context) airtel.PushCallbackFunc {
		return func(callbackFunc airtel.CallbackRequest) error {
			return nil
		}
	}
	config := &airtel.Config{
		AllowedCountries:   nil,
		DisbursePIN:        "",
		CallbackPrivateKey: "",
		CallbackAuth:       false,
		PublicKey:          "",
		Environment:        airtel.PRODUCTION, //airtel.STAGING,
		ClientID:           "",                //given by airtel
		Secret:             "",                //given by airtel
	}
	client := airtel.NewClient(config, pushCallbackFunc(context.TODO()), true)

	response, err := client.Token(context.TODO())

	if err != nil {
		//handle error
	}

	fmt.Printf("%+v", response)
}

```
## collection

send push pay request to user

```go
package main

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel"
)

func main() {
	pushCallbackFunc := func(ctx context.Context) airtel.PushCallbackFunc {
		return func(callbackFunc airtel.CallbackRequest) error {
			return nil
		}
	}
	config := &airtel.Config{
		AllowedCountries:   nil,
		DisbursePIN:        "",
		CallbackPrivateKey: "",
		CallbackAuth:       false,
		PublicKey:          "",
		Environment:        airtel.PRODUCTION, //airtel.STAGING,
		ClientID:           "",                //given by airtel
		Secret:             "",                //given by airtel
	}
	client := airtel.NewClient(config, pushCallbackFunc(context.TODO()), true)
	
	request := airtel.PushPayRequest{
		Reference:          "",
		SubscriberCountry:  "",
		SubscriberMsisdn:   "",
		TransactionAmount:  0,
		TransactionCountry: "",
		TransactionID:      "",
	}

	response, err := client.Push(context.TODO(), request)

	if err != nil {
		//handle error
	}

	fmt.Printf("%+v", response)
}

```

## disbursement
## account
## transaction
## KYC
## coverage:

- [x] **Authorization**
- [x] **Encryption**
- [x] **Collection**
- [x] **Disbursement**
- [x] **Account**
- [x] **Transactions**
- [x] **KYC**
- [ ] Billers Callback
- [ ] IMT Callback
- [ ] Remittance


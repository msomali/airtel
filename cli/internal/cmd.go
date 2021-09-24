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

package internal

import (
	"fmt"
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/api"
	"github.com/techcraftlabs/airtel/api/http"
	"github.com/techcraftlabs/airtel/internal/models"
	"github.com/techcraftlabs/airtel/pkg/countries"
	clix "github.com/urfave/cli/v2"
	"time"
)

var _ Commander = (*Cmd)(nil)

const (
	json outFormat = iota + 1
	text
)

type (
	Cmd struct {
		ApiClient   *http.Client
		RequestType airtel.RequestType
		Name        string
		Usage       string
		Description string
		Flags       []clix.Flag
		SubCommands []*clix.Command
	}
	outFormat int
	Commander interface {
		Command() *clix.Command
		Before(ctx *clix.Context) error
		After(ctx *clix.Context) error
		Action(ctx *clix.Context) error
		OnError(ctx *clix.Context, err error, isSubcommand bool) error
		PrintOut(payload interface{}, format outFormat) error
	}
)

func (c *Cmd) Command() *clix.Command {
	cmd := &clix.Command{
		Name:        c.Name,
		Usage:       c.Usage,
		Description: c.Description,
		Before: func(ctx *clix.Context) error {
			return c.Before(ctx)
		},
		After: func(ctx *clix.Context) error {
			return c.After(ctx)
		},
		Action: func(ctx *clix.Context) error {
			return c.Action(ctx)
		},
		OnUsageError: func(ctx *clix.Context, err error, isSubcommand bool) error {
			return c.OnError(ctx, err, isSubcommand)
		},
		Subcommands: c.SubCommands,
		Flags:       c.Flags,
	}
	return cmd
}

func (c *Cmd) Before(ctx *clix.Context) error {
	return nil
}

func (c *Cmd) After(ctx *clix.Context) error {
	return nil
}

func (c *Cmd) Action(ctx *clix.Context) error {
	switch c.RequestType {
	case airtel.Authorization:
		authResponse, err := c.ApiClient.Token(ctx.Context)
		if err != nil {
			return err
		}

		err = c.PrintOut(authResponse, text)
		if err != nil {
			return err
		}

		return nil
//0784956141
	case airtel.USSDPush:
		ref := ctx.String("reference")
		phone := ctx.String("phone")
		amount := ctx.Int64("amount")
		enquiry := ctx.Bool("enquiry")
		if enquiry{
			id := ctx.String("reference")
			fmt.Printf("enquire about transaction of id %v\n",id)
			return nil
		}
		req := api.PushPayRequest{
			Reference:          ref,
			SubscriberCountry:  countries.TANZANIA,
			SubscriberMsisdn:   phone,
			TransactionAmount:  amount,
			TransactionCountry: countries.TANZANIA,
			TransactionID:      fmt.Sprintf("%d", time.Now().UnixNano()),
		}
		pushPayResponse, err := c.ApiClient.Push(ctx.Context, req)
		if err != nil {
			return err
		}

		err = c.PrintOut(pushPayResponse, text)
		if err != nil {
			return err
		}

		return nil

	default:
		return nil
	}
}

func (c *Cmd) OnError(ctx *clix.Context, err error, isSubcommand bool) error {
	return nil
}

func (c *Cmd) PrintOut(payload interface{}, format outFormat) error {
	switch c.RequestType {
	case airtel.Authorization:
		resp, ok := payload.(models.AirtelAuthResponse)
		if !ok {
			return fmt.Errorf("bad request expected models.AirtelAuthResponse")
		}
		tokenResponsePrintOut(resp)
		return nil

	case airtel.USSDPush:
		resp, ok := payload.(api.PushPayResponse)
		if !ok {
			return fmt.Errorf("bad request expected models.AirtelAuthResponse")
		}
		pushResponsePrintOut(resp)
		return nil

	}
	return nil
}

package airtel

import (
	"context"
	"fmt"
	"github.com/techcraftlabs/airtel/pkg/models"
	"time"
)

type (
	// Params
	//From	query	integer(int64)	true	Date from which transactions are to be fetched.
	//To	query	integer(int64)	true	Date until transactions are to be fetched.
	//Limit	query	integer(int64)	true	The number of transactions to be fetched on a page.
	//Offset	query	integer(int64)	true	Page number from which transactions are to be fetched.
	Params struct {
		From   int64 `json:"from"`
		To     int64 `json:"to"`
		Limit  int64 `json:"limit"`
		Offset int64 `json:"offset"`
	}
	TransactionService interface {
		Summary(ctx context.Context, params Params) (models.ListTransactionsResponse, error)
	}
)

func (c *Client) Summary(ctx context.Context, params Params) (models.ListTransactionsResponse, error) {

	//?from=4&to=7&offset=4&limit=8
	queryString := fmt.Sprintf("?from=%d&to=%d&limit=%d&offset=%d",params.From, params.To,params.Limit,params.Offset)
	fmt.Printf("query string: %s\n",queryString)

	var token string
	if *c.token == "" {
		str, err := c.Token(ctx)
		if err != nil {
			return models.ListTransactionsResponse{}, err
		}
		token = fmt.Sprintf("%s", str.AccessToken)
	}
	//Add Auth Header
	if *c.token != "" {
		if !c.tokenExpiresAt.IsZero() && time.Until(c.tokenExpiresAt) < (60*time.Second) {
			if _, err := c.Token(ctx); err != nil {
				return models.ListTransactionsResponse{}, err
			}
		}
		token = *c.token
	}
	req, err := createInternalRequest("", c.Conf.Environment, AccountBalance, token, nil, "")
	if err != nil {
		return models.ListTransactionsResponse{}, err
	}

	res := new(models.ListTransactionsResponse)

	_, err = c.base.Do(ctx, "check balance", req, res)
	if err != nil {
		return models.ListTransactionsResponse{}, err
	}
	return *res, nil
}

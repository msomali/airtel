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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/techcraftlabs/airtel/internal"
	"github.com/techcraftlabs/airtel/internal/models"
	"time"
)

type Authenticator interface {
	Token(ctx context.Context) (models.TokenResponse, error)
}

func (c *Client) checkToken(ctx context.Context) (string, error) {
	var token string
	if *c.token == "" {
		str, err := c.Token(ctx)
		if err != nil {
			return "", err
		}
		token = fmt.Sprintf("%s", str.AccessToken)
	}
	//Add Auth Header
	if *c.token != "" {
		if !c.tokenExpiresAt.IsZero() && time.Until(c.tokenExpiresAt) < (60*time.Second) {
			if _, err := c.Token(ctx); err != nil {
				return "", err
			}
		}
		token = *c.token
	}

	return token, nil
}

func (c *Client) Token(ctx context.Context) (models.TokenResponse, error) {
	body := models.TokenRequest{
		ClientID:     c.Conf.ClientID,
		ClientSecret: c.Conf.Secret,
		GrantType:    defaultGrantType,
	}

	var opts []internal.RequestOption
	hs := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "*/*",
	}

	opts = append(opts, internal.WithRequestHeaders(hs))
	req := c.makeInternalRequest(Authorization, body, opts...)

	res := new(models.TokenResponse)
	reqName := Authorization.Name()
	_, err := c.base.Do(ctx, reqName, req, res)
	if err != nil {
		return models.TokenResponse{}, err
	}
	duration := time.Duration(res.ExpiresIn)
	now := time.Now()
	later := now.Add(time.Second * duration)
	c.tokenExpiresAt = later
	*c.token = res.AccessToken
	return *res, nil
}

func PinEncryption(pin string, pubKey string) (string, error) {

	decodedBase64, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return "", fmt.Errorf("could not decode pub key to Base64 string: %w", err)
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(decodedBase64)
	if err != nil {
		return "", fmt.Errorf("could not parse encoded public key (encryption key) : %w", err)
	}

	//check if the public key is RSA public key
	publicKey, isRSAPublicKey := publicKeyInterface.(*rsa.PublicKey)
	if !isRSAPublicKey {
		return "", fmt.Errorf("public key parsed is not an RSA public key : %w", err)
	}

	msg := []byte(pin)

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg)

	if err != nil {
		return "", fmt.Errorf("could not encrypt api key using generated public key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil

}

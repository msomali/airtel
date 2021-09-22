package airtel

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/techcraftlabs/airtel/pkg/models"
)

type Authenticator interface {
	Token(ctx context.Context) (models.AirtelAuthResponse, error)
}

func (c *Client) Token(ctx context.Context) (models.AirtelAuthResponse, error) {
	body := models.AirtelAuthRequest{
		ClientID:     c.conf.ClientID,
		ClientSecret: c.conf.Secret,
		GrantType:    defaultGrantType,
	}
	req, err := createInternalRequest("", c.conf.Environment, Authorization, "", body, "")
	if err != nil {
		return models.AirtelAuthResponse{}, err
	}

	res := new(models.AirtelAuthResponse)

	_, err = c.base.Do(ctx, "Token", req, res)
	if err != nil {
		return models.AirtelAuthResponse{}, err
	}
	//fmt.Printf("status code: %v\nheaders: %v\npayload: %v\nerror: %v\n", do.StatusCode, do.Headers, do.Payload, do.Error)
	*c.token = res.AccessToken
	return *res, nil
}


func generateEncryptedKey(apiKey, pubKey string) (string, error) {

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

	msg := []byte(apiKey)

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg)

	if err != nil {
		return "", fmt.Errorf("could not encrypt api key using generated public key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil

}

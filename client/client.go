package client

import (
	"cryptohelper"
	"errors"
)

type Client struct {
	accountID string
	authToken  string
}

// NewClient creates a new warranted client instance
func NewClient(accountID, authToken string) (*Client, error) {
	if accountID == "" {
		return nil, errors.New("no accountId provided")
	}
	if authToken == "" {
		return nil, errors.New("no authToken provided")
	}
	return &Client{accountID: accountID, authToken: authToken}, nil
}

// ValidateRequest validates the signature of a request
func (c *Client) ValidateRequest(signature, url, jsonData string) bool {
	hmac, err := cryptohelper.CreateHMAC(url, jsonData, c.authToken, "sha256")
	if err != nil {
		return false
	}
	return cryptohelper.TimeSafeCompare(signature, hmac)
}

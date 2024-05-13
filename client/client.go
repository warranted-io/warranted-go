package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const jsonContentType = "application/json"
const libraryVersion = "1.0.0"

func extractContentTypeHeader(headers map[string]interface{}) string {
	headerType, ok := headers["Content-Type"]
	if !ok {
		return jsonContentType
	}
	return headerType.(string)
}

type Client struct {
	accountID  string
	authToken  string
	host       string
	headers    http.Header
	HTTPClient *http.Client
}

// NewClient creates a new warranted client instance
func NewClient(accountID, authToken string) (*Client, error) {
	if accountID == "" {
		return nil, errors.New("no accountId provided")
	}
	if authToken == "" {
		return nil, errors.New("no authToken provided")
	}
	return &Client{
		accountID: accountID,
		authToken: authToken,
		host:      "https://app.warranted.io",
		headers:   make(http.Header),
	}, nil
}

// default http Client should not follow redirects and return the most recent response.
func defaultHTTPClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 10,
	}
}

func (c *Client) basicAuth() (string, string) {
	return c.accountID, c.authToken
}

func (c *Client) doWithErr(req *http.Request) (*http.Response, error) {
	client := c.HTTPClient

	if client == nil {
		client = defaultHTTPClient()
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Note that 3XX response codes are not allowed for fetches
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		err = &WarrantedError{}
		if decodeErr := json.NewDecoder(res.Body).Decode(err); decodeErr != nil {
			err = errors.Wrap(decodeErr, "error decoding the response for an HTTP error code: "+strconv.Itoa(res.StatusCode))
			return nil, err
		}

		return nil, err
	}
	return res, nil
}

// SendRequest verifies, constructs, and authorizes an HTTP request.
func (c *Client) SendRequest(method string, rawURL string, data url.Values,
	headers map[string]interface{}, body ...byte) (*http.Response, error) {

	contentType := extractContentTypeHeader(headers)

	u, err := url.Parse(c.host + rawURL)
	if err != nil {
		return nil, err
	}

	valueReader := &strings.Reader{}
	goVersion := runtime.Version()
	var req *http.Request

	//data is already processed and information will be added to u(the url) in the
	//previous step. Now body will solely contain json payload
	if contentType == jsonContentType {
		req, err = http.NewRequest(method, u.String(), bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
	} else {
		//Here the HTTP POST methods which is not having json content type are processed
		//All the values will be added in data and encoded (all body, query, path parameters)
		if method == http.MethodPost {
			valueReader = strings.NewReader(data.Encode())
		}

		req, err = http.NewRequest(method, u.String(), valueReader)
		if err != nil {
			return nil, err
		}

	}

	if contentType == jsonContentType {
		req.Header.Add("Content-Type", jsonContentType)
	}

	req.SetBasicAuth(c.basicAuth())

	// E.g. "User-Agent": "warranted-go/1.0.0 (darwin amd64) go/go1.22.0"
	userAgent := fmt.Sprintf("warranted-go/%s (%s %s) go/%s", libraryVersion, runtime.GOOS, runtime.GOARCH, goVersion)

	req.Header.Add("User-Agent", userAgent)

	for k, v := range headers {
		req.Header.Add(k, fmt.Sprint(v))
	}

	// Add custom headers
	for k, v := range c.headers {
		req.Header.Add(k, fmt.Sprint(v))
	}

	return c.doWithErr(req)
}

// ValidateRequest validates the signature of a request
func (c *Client) ValidateRequest(signature, url, jsonData string) bool {
	hmac, err := CreateHMAC(url, jsonData, c.authToken, "sha256")
	if err != nil {
		return false
	}
	return TimeSafeCompare(signature, hmac)
}

func (c *Client) SetHeader(headers http.Header) {
	c.headers = headers
}

func (c *Client) SetHost(host string) {
	c.host = host
}

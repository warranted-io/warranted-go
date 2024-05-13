// Package client provides internal utilities for the twilio-go client library.
package client

import (
	"net/http"
	"net/url"
)

type RequestHandler struct {
	Client Client
	Edge   string
	Region string
}

func NewRequestHandler(client Client) *RequestHandler {
	return &RequestHandler{
		Client: client,
	}
}

func (c *RequestHandler) sendRequest(method string, url string, data url.Values,
	headers map[string]interface{}, body ...byte) (*http.Response, error) {
	return c.Client.SendRequest(method, url, data, headers, body...)
}

func (c *RequestHandler) Post(path string, bodyData url.Values, headers map[string]interface{}, body ...byte) (*http.Response, error) {
	return c.sendRequest(http.MethodPost, path, bodyData, headers, body...)
}

func (c *RequestHandler) Put(path string, bodyData url.Values, headers map[string]interface{}, body ...byte) (*http.Response, error) {
	return c.sendRequest(http.MethodPut, path, bodyData, headers, body...)
}

func (c *RequestHandler) Get(path string, queryData url.Values, headers map[string]interface{}) (*http.Response, error) {
	return c.sendRequest(http.MethodGet, path, queryData, headers)
}

func (c *RequestHandler) Delete(path string, nothing url.Values, headers map[string]interface{}) (*http.Response, error) {
	return c.sendRequest(http.MethodDelete, path, nil, headers)
}

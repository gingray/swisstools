package sentry

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
)

type client struct {
	baseUrl  string
	apiToken string
	http     *http.Client
}

var apiVersion = "api/0"

func newClient(sentry *Sentry) *client {
	baseUrl, _ := url.JoinPath(sentry.Url, apiVersion)
	return &client{baseUrl: baseUrl, http: &http.Client{}, apiToken: sentry.ApiToken}
}

func (c *client) makeRequest(method string, apiSegment string, body []byte) (*http.Response, error) {
	reqUrl, _ := url.JoinPath(c.baseUrl, apiSegment)
	return c.do(method, reqUrl, body)
}

func (c *client) do(method string, url string, body []byte) (*http.Response, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, _ := http.NewRequest(method, url, reader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))
	req.Header.Add("Accept", "application/json")
	req.Close = true

	resp, err := c.http.Do(req)

	if err != nil {
		log.Error(err)
	}

	return resp, err
}

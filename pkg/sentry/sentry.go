package sentry

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
	"net/http"
	"net/url"
)

type Sentry struct {
	ApiToken string
	Url      string
}
type client struct {
	baseUrl  string
	apiToken string
	http     *http.Client
}

func newClient(sentry *Sentry) *client {
	return &client{baseUrl: sentry.Url, http: &http.Client{}, apiToken: sentry.ApiToken}
}

func NewSentry(cfg *common.Config) *Sentry {
	return &Sentry{ApiToken: cfg.Sentry.ApiToken, Url: cfg.Sentry.Url}
}

func (c *client) makeRequest() {
	url, _ := url.JoinPath(c.baseUrl, "api/0/projects/")
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))
	req.Header.Add("Accept", "application/json")
	req.Close = true

	resp, err := c.http.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Error(err)
	}

	fmt.Println(resp.StatusCode)

}

func (s *Sentry) GetIssues() {
	client := newClient(s)
	client.makeRequest()
}

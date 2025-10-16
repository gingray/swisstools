package sentry

import (
	"fmt"

	"github.com/gingray/swisstools/pkg/common"
)

type Sentry struct {
	ApiToken string
	Url      string
}

func NewSentry(cfg *common.Config) *Sentry {
	return &Sentry{ApiToken: cfg.Sentry.ApiToken, Url: cfg.Sentry.Url}
}

func (s *Sentry) GetTagValues(organization string, project string, tag string) {
	client := newClient(s)
	tags, _ := client.GetTagValues(organization, project, tag)
	fmt.Println(tags)
}

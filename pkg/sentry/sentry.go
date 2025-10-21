package sentry

import (
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
)

type Sentry struct {
	ApiToken string
	Url      string
	View     common.ViewRecords
}

func NewSentry(cfg *common.Config, view common.ViewRecords) *Sentry {
	return &Sentry{ApiToken: cfg.Sentry.ApiToken, Url: cfg.Sentry.Url, View: view}
}

func (s *Sentry) GetTagValues(organization string, project string, tag string) {
	client := newClient(s)
	tags, _ := client.GetTagValues(organization, project, tag)
	dataView := common.NewDataView()
	for _, key := range []string{"Idx", "Name", "Link", "LastSeen"} {
		dataView.AddKey(key)
	}

	for idx, item := range tags {
		row := map[string]string{"Idx": strconv.Itoa(idx + 1), "Name": item.Name, "Link": generateLinkToDashboard(s.Url, organization, project, tag, item.Value), "LastSeen": item.LastSeen.Format("2006-01-02 15:04:05")}
		dataView.AddRow(row)
	}
	err := s.View.Show(dataView)
	if err != nil {
		log.Error(err)
	}

}

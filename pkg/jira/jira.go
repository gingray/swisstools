package jira

import (
	"context"
	"fmt"
	jira2 "github.com/andygrunwald/go-jira/v2/onpremise"
	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
)

const DEFAULT_SEARCH_QUERY = "project = \"%s\" AND assignee = currentuser() AND status NOT IN (\"Done\", \"Canceled\") order by created"

type Jira struct {
	ApiToken string
	Url      string
	Project  string
}

func NewJira(cfg *common.Config) *Jira {
	return &Jira{
		ApiToken: cfg.Jira.ApiToken,
		Url:      cfg.Jira.Url,
		Project:  cfg.Jira.Project,
	}
}

func (j *Jira) GetIssues() {
	jiraURL := j.Url

	tp := jira2.BearerAuthTransport{
		Token: j.ApiToken,
	}
	client, err := jira2.NewClient(jiraURL, tp.Client())
	if err != nil {
		log.Error(err)
	}
	query := fmt.Sprintf(DEFAULT_SEARCH_QUERY, j.Project)
	opt := &jira2.SearchOptions{
		MaxResults: 100, // Max results can go up to 1000
		StartAt:    0,
	}
	issues, _, err := client.Issue.Search(context.TODO(), query, opt)
	for _, issue := range issues {
		fmt.Println(issue.Key)
	}
}

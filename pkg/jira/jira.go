package jira

import (
	"context"
	"fmt"
	jira2 "github.com/andygrunwald/go-jira/v2/onpremise"
	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
	"net/url"
	"time"
)

const DEFAULT_SEARCH_QUERY = "project = \"%s\" AND assignee = currentuser() AND status NOT IN (\"Done\", \"Canceled\") order by created"

type Jira struct {
	ApiToken string
	Url      string
	Project  string
	View     common.ViewRecords
}

func NewJira(cfg *common.Config, view common.ViewRecords) *Jira {
	return &Jira{
		ApiToken: cfg.Jira.ApiToken,
		Url:      cfg.Jira.Url,
		Project:  cfg.Jira.Project,
		View:     view,
	}
}

type issue struct {
	Url       string
	Title     string
	Status    string
	CreatedAt time.Time
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
	var tableRows [][]string
	for _, item := range issues {
		issue := issue{
			Url:       getJiraUrl(j.Url, item.Key),
			Title:     item.Fields.Summary,
			Status:    item.Fields.Status.Name,
			CreatedAt: time.Time(item.Fields.Created),
		}
		tableRows = append(tableRows, []string{issue.Url, issue.Title, issue.Status, issue.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	err = j.View.Show(tableRows)
	if err != nil {
		log.Error(err)
	}
}

func getJiraUrl(baseUrl, key string) string {
	issueUrl, err := url.JoinPath(baseUrl, fmt.Sprintf("browse/%s", key))
	if err != nil {
		log.Error(err)
	}
	return issueUrl
}

package jira

import (
	"context"
	"fmt"
	jira2 "github.com/andygrunwald/go-jira/v2/onpremise"
	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"net/url"
	"os"
	"time"
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

	table := tablewriter.NewTable(os.Stdout, tablewriter.WithStreaming(tw.StreamConfig{Enable: true}))
	defer table.Close()
	table.Header([]string{"Url", "Title", "Status", "Created"})
	for _, item := range issues {
		jiraUrl, err := url.JoinPath(j.Url, fmt.Sprintf("browse/%s", j.Url, item.Key))
		if err != nil {
			log.Error(err)
		}
		issue := issue{
			Url:       jiraUrl,
			Title:     item.Fields.Summary,
			Status:    item.Fields.Status.Name,
			CreatedAt: time.Time(item.Fields.Created),
		}
		table.Append([]string{issue.Url, issue.Title, issue.Status, issue.CreatedAt.Format("2006-01-02 15:04:05")})
		//fmt.Printf("%v\n", issue)
	}
}

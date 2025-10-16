package mr

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Gitlab struct {
	Url      string
	Token    string
	Authors  []string
	Projects []string
	view     common.ViewRecords
}

type gitLabUser struct {
	Id   int
	Name string
}
type mergeRequest struct {
	Url       string
	Title     string
	Author    string
	UpdatedAt time.Time
}

func NewGitlab(cfg *common.Config, view common.ViewRecords) *Gitlab {
	return &Gitlab{
		Url:      cfg.GitLab.Url,
		Token:    cfg.GitLab.ApiToken,
		Authors:  cfg.GitLab.Authors,
		Projects: cfg.GitLab.Projects,
		view:     view,
	}
}

func (g *Gitlab) FetchMrs() {
	git, err := gitlab.NewClient(g.Token, gitlab.WithBaseURL(g.Url))
	if err != nil {
		log.Error(err)
	}
	var gitLabUsers []gitLabUser
	for _, user := range g.Authors {
		remoteUsers, _, err := git.Users.ListUsers(&gitlab.ListUsersOptions{Username: &user})
		if err != nil {
			log.Error(err)
		}
		for _, remoteUser := range remoteUsers {
			gitLabUsers = append(gitLabUsers, gitLabUser{Id: remoteUser.ID, Name: remoteUser.Username})
		}
	}
	var mergeRequests []mergeRequest
	for _, repo := range g.Projects {
		state := "opened"
		for _, user := range gitLabUsers {
			mrs, _, err := git.MergeRequests.ListProjectMergeRequests(repo, &gitlab.ListProjectMergeRequestsOptions{State: &state, AuthorID: &user.Id})
			if err != nil {
				log.Error(err)
				continue
			}
			for _, mr := range mrs {
				mergeRequests = append(mergeRequests, mergeRequest{Url: mr.WebURL, Title: mr.Title, Author: mr.Author.Username, UpdatedAt: *mr.UpdatedAt})
			}

		}
	}
	dataView := common.NewDataView()
	for _, key := range []string{"Url", "Title", "Author", "Updated"} {
		dataView.AddKey(key)
	}
	for _, item := range mergeRequests {
		row := map[string]string{"Url": item.Url, "Title": item.Title, "Author": item.Author, "Updated": item.UpdatedAt.Format("2006-01-02 15:04:05")}
		dataView.AddRow(row)
	}
	err = g.view.Show(dataView)
}

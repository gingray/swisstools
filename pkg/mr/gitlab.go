package mr

import (
	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Gitlab struct {
	Url      string
	Token    string
	Authors  []string
	Projects []string
}

type gitLabUser struct {
	Id   int
	Name string
}
type mergeRequest struct {
	Title string
}

func NewGitlab(cfg *common.Config) *Gitlab {
	return &Gitlab{
		Url:   cfg.GitLab.Url,
		Token: cfg.GitLab.ApiToken,
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
				mergeRequests = append(mergeRequests, mergeRequest{Title: mr.Title})
			}

		}
	}
	//git.Search.
}

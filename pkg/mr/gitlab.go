package mr

import (
	"sort"

	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
	"github.com/gingray/swisstools/pkg/service"
)

type client interface {
	Users(filter service.UserFilter) ([]service.GitlabUser, error)
	ProjectMRs(filter service.ProjectFilter) ([]service.GitlabProject, error)
}
type Gitlab struct {
	Authors      []string
	Projects     []string
	gitlabClient client
	view         common.ViewRecords
}

func NewGitlab(cfg *common.Config, gitlabClient client, view common.ViewRecords) *Gitlab {
	return &Gitlab{
		Authors:      cfg.GitLab.Authors,
		Projects:     cfg.GitLab.Projects,
		view:         view,
		gitlabClient: gitlabClient,
	}
}

func (g *Gitlab) FetchMrs() {
	var gitlabUserIDs []int
	for _, user := range g.Authors {
		remoteUsers, err := g.gitlabClient.Users(service.UserFilter{Name: user})
		if err != nil {
			log.Error(err)
		}
		for _, remoteUser := range remoteUsers {
			gitlabUserIDs = append(gitlabUserIDs, remoteUser.ID)
		}
	}
	mergeRequests, err := g.gitlabClient.ProjectMRs(service.ProjectFilter{Projects: g.Projects, UserIDs: gitlabUserIDs})
	if err != nil {
		log.Error(err)
	}
	dataView := common.NewDataView()
	for _, key := range []string{"Url", "Title", "Author", "Updated"} {
		dataView.AddKey(key)
	}
	sort.Slice(mergeRequests, func(i, j int) bool {
		if mergeRequests[i].Author != mergeRequests[j].Author {
			return mergeRequests[i].Author < mergeRequests[j].Author
		}
		return mergeRequests[i].UpdatedAt.After(mergeRequests[j].UpdatedAt)
	})
	for _, item := range mergeRequests {
		row := map[string]string{"Url": item.Url, "Title": item.Title, "Author": item.Author, "Updated": item.UpdatedAt.Format("2006-01-02 15:04:05")}
		dataView.AddRow(row)
	}
	err = g.view.Show(dataView)
}

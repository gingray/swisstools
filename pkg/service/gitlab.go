package service

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Gitlab struct {
	cfg    *common.GitLabConfig
	logger *log.Logger
	client *gitlab.Client
}

type UserFilter struct {
	Name string
}

type ProjectFilter struct {
	Projects []string
	UserIDs  []int
	State    string
}

type BranchFilter struct {
	User   string
	Branch string
}

type GitlabUser struct {
	ID   int
	Name string
}

type GitlabProject struct {
	Title     string
	Url       string
	Author    string
	UpdatedAt time.Time
}

type BranchChanges struct {
	Branch       string
	ChangedFiles []string
}

func NewGitlab(cfg *common.GitLabConfig, logger *log.Logger) (*Gitlab, error) {
	client, err := gitlab.NewClient(cfg.ApiToken, gitlab.WithBaseURL(cfg.Url))
	if err != nil {
		return nil, err
	}
	return &Gitlab{cfg: cfg, logger: logger, client: client}, nil
}

func (g *Gitlab) Users(filter UserFilter) ([]GitlabUser, error) {
	remoteUsers, _, err := g.client.Users.ListUsers(&gitlab.ListUsersOptions{Username: &filter.Name})
	if err != nil {
		return nil, err
	}
	var gitlabUsers []GitlabUser
	for _, remoteUser := range remoteUsers {
		gitlabUsers = append(gitlabUsers, GitlabUser{ID: remoteUser.ID, Name: remoteUser.Username})
	}
	return gitlabUsers, nil
}

func (g *Gitlab) ProjectMRs(filter ProjectFilter) ([]GitlabProject, error) {
	var mergeRequests []GitlabProject
	for _, repo := range filter.Projects {
		state := "opened"
		for _, user := range filter.UserIDs {
			mrs, _, err := g.client.MergeRequests.ListProjectMergeRequests(repo, &gitlab.ListProjectMergeRequestsOptions{State: &state, AuthorID: &user})
			if err != nil {
				log.Error(err)
				continue
			}
			for _, mr := range mrs {
				mergeRequests = append(mergeRequests, GitlabProject{Url: mr.WebURL, Title: mr.Title, Author: mr.Author.Username, UpdatedAt: *mr.UpdatedAt})
			}
		}
	}
	return mergeRequests, nil
}

func (g *Gitlab) ChangedFilesForBranch(filter BranchFilter) (*BranchChanges, error) {
	state := "all"
	orderBy := "updated_at"
	sort := "desc"

	gitLabUsers, err := g.Users(UserFilter{Name: filter.User})
	if err != nil {
		log.Error(err)
	}
	var mr gitlab.BasicMergeRequest
	for _, user := range gitLabUsers {
		mrs, _, err := g.client.MergeRequests.ListMergeRequests(&gitlab.ListMergeRequestsOptions{
			AuthorID:     &user.ID,
			ListOptions:  gitlab.ListOptions{Page: 1, PerPage: 1},
			State:        gitlab.Ptr(state),
			SourceBranch: gitlab.Ptr(filter.Branch),
			OrderBy:      gitlab.Ptr(orderBy),
			Sort:         gitlab.Ptr(sort),
		})
		if err != nil {
			return nil, fmt.Errorf("list merge requests for project %q: %w", user, err)
		}
		if len(mrs) == 0 {
			return nil, fmt.Errorf("no merge request found for user %q", user)
		}
		mr = *mrs[0]
		break
	}

	filesChanged, err := g.mergeRequestFilePaths(&mr)
	if err != nil {
		return nil, err
	}
	return &BranchChanges{Branch: mr.SourceBranch, ChangedFiles: filesChanged}, nil
}

func (g *Gitlab) mergeRequestFilePaths(mr *gitlab.BasicMergeRequest) ([]string, error) {
	page := 1
	const perPage = 100
	var paths []string
	seen := make(map[string]struct{})

	for {
		diffs, resp, err := g.client.MergeRequests.ListMergeRequestDiffs(mr.ProjectID, mr.IID, &gitlab.ListMergeRequestDiffsOptions{
			ListOptions: gitlab.ListOptions{Page: page, PerPage: perPage},
		})
		if err != nil {
			return nil, fmt.Errorf("merge request diffs: %w", err)
		}
		for _, d := range diffs {
			p := g.diffFilePath(d)
			if p == "" {
				continue
			}
			if _, ok := seen[p]; ok {
				continue
			}
			seen[p] = struct{}{}
			paths = append(paths, p)
		}
		if resp.NextPage == 0 || len(diffs) < perPage {
			break
		}
		page = resp.NextPage
	}

	return paths, nil
}

func (g *Gitlab) diffFilePath(d *gitlab.MergeRequestDiff) string {
	if d == nil {
		return ""
	}
	if d.DeletedFile {
		return ""
	}
	return d.NewPath
}

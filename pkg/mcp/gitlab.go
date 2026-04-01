package mcp

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Gitlab struct {
	cfg *common.Config
}

type gitLabUser struct {
	Id   int
	Name string
}

// BranchMergeRequestFiles is the merge request GitLab associates with a source branch and the paths it changes.
type BranchMergeRequestFiles struct {
	Project      int      `json:"projectID"`
	IID          int      `json:"iid"`
	WebURL       string   `json:"web_url"`
	Title        string   `json:"title"`
	State        string   `json:"state"`
	ChangedFiles []string `json:"changed_files"`
}

func NewGitlab(cfg *common.Config) *Gitlab {
	return &Gitlab{cfg: cfg}
}

// ChangedFilesForBranch finds the most recently updated merge request whose source branch matches branch
// across cfg.Projects, then lists all changed file paths from that MR’s diffs.
func (g *Gitlab) ChangedFilesForBranch(branch string) (*BranchMergeRequestFiles, error) {
	branch = strings.TrimSpace(branch)
	if branch == "" {
		return nil, errors.New("branch name is empty")
	}
	if g.cfg.GitLab.Url == "" || g.cfg.GitLab.ApiToken == "" {
		return nil, errors.New("gitlab url and apiToken must be configured")
	}

	client, err := gitlab.NewClient(g.cfg.GitLab.ApiToken, gitlab.WithBaseURL(g.cfg.GitLab.Url))
	if err != nil {
		return nil, fmt.Errorf("gitlab client: %w", err)
	}

	project, mr, err := findMergeRequestForBranch(client, g.cfg.MCP.GitLabUser, branch)
	if err != nil {
		return nil, err
	}

	files, err := mergeRequestFilePaths(client, project, mr.IID)
	if err != nil {
		return nil, err
	}

	return &BranchMergeRequestFiles{
		Project:      project,
		IID:          mr.IID,
		WebURL:       mr.WebURL,
		Title:        mr.Title,
		State:        mr.State,
		ChangedFiles: files,
	}, nil
}

func findMergeRequestForBranch(client *gitlab.Client, author string, branch string) (int, *gitlab.BasicMergeRequest, error) {
	state := "all"
	orderBy := "updated_at"
	sort := "desc"

	var bestProject int
	var best *gitlab.BasicMergeRequest
	var bestUpdated time.Time
	var gitLabUsers []gitLabUser

	remoteUsers, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{Username: &author})
	if err != nil {
		log.Error(err)
	}
	for _, remoteUser := range remoteUsers {
		gitLabUsers = append(gitLabUsers, gitLabUser{Id: remoteUser.ID, Name: remoteUser.Username})
	}

	for _, user := range gitLabUsers {
		mrs, _, err := client.MergeRequests.ListMergeRequests(&gitlab.ListMergeRequestsOptions{
			AuthorID:     &user.Id,
			ListOptions:  gitlab.ListOptions{Page: 1, PerPage: 1},
			State:        gitlab.Ptr(state),
			SourceBranch: gitlab.Ptr(branch),
			OrderBy:      gitlab.Ptr(orderBy),
			Sort:         gitlab.Ptr(sort),
		})
		if err != nil {
			return 0, nil, fmt.Errorf("list merge requests for project %q: %w", user, err)
		}
		if len(mrs) == 0 {
			continue
		}
		mr := mrs[0]
		updated := time.Time{}
		if mr.UpdatedAt != nil {
			updated = *mr.UpdatedAt
		}
		if best == nil || updated.After(bestUpdated) {
			best = mr
			bestUpdated = updated
			bestProject = mr.ProjectID
		}
	}

	if best == nil {
		return 0, nil, fmt.Errorf("no merge request found for source branch %q in configured projects", branch)
	}
	return bestProject, best, nil
}

func mergeRequestFilePaths(client *gitlab.Client, project int, mergeRequestIID int) ([]string, error) {
	page := 1
	const perPage = 100
	var paths []string
	seen := make(map[string]struct{})

	for {
		diffs, resp, err := client.MergeRequests.ListMergeRequestDiffs(project, mergeRequestIID, &gitlab.ListMergeRequestDiffsOptions{
			ListOptions: gitlab.ListOptions{Page: page, PerPage: perPage},
		})
		if err != nil {
			return nil, fmt.Errorf("merge request diffs: %w", err)
		}
		for _, d := range diffs {
			p := diffFilePath(d)
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

func diffFilePath(d *gitlab.MergeRequestDiff) string {
	if d == nil {
		return ""
	}
	if d.DeletedFile {
		return ""
	}
	return d.NewPath
}

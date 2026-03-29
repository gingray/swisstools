package mcp

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gingray/swisstools/pkg/common"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// BranchMergeRequestFiles is the merge request GitLab associates with a source branch and the paths it changes.
type BranchMergeRequestFiles struct {
	Project string   `json:"project"`
	IID     int      `json:"iid"`
	WebURL  string   `json:"web_url"`
	Title   string   `json:"title"`
	State   string   `json:"state"`
	Files   []string `json:"files"`
}

// ChangedFilesForBranch finds the most recently updated merge request whose source branch matches branch
// across cfg.Projects, then lists all changed file paths from that MR’s diffs.
func ChangedFilesForBranch(cfg *common.GitLabConfig, branch string) (*BranchMergeRequestFiles, error) {
	branch = strings.TrimSpace(branch)
	if branch == "" {
		return nil, errors.New("branch name is empty")
	}
	if cfg.Url == "" || cfg.ApiToken == "" {
		return nil, errors.New("gitlab url and apiToken must be configured")
	}
	if len(cfg.Projects) == 0 {
		return nil, errors.New("no gitlab projects configured")
	}

	client, err := gitlab.NewClient(cfg.ApiToken, gitlab.WithBaseURL(cfg.Url))
	if err != nil {
		return nil, fmt.Errorf("gitlab client: %w", err)
	}

	project, mr, err := findMergeRequestForBranch(client, cfg.Projects, branch)
	if err != nil {
		return nil, err
	}

	files, err := mergeRequestFilePaths(client, project, mr.IID)
	if err != nil {
		return nil, err
	}

	return &BranchMergeRequestFiles{
		Project: project,
		IID:     mr.IID,
		WebURL:  mr.WebURL,
		Title:   mr.Title,
		State:   mr.State,
		Files:   files,
	}, nil
}

func findMergeRequestForBranch(client *gitlab.Client, projects []string, branch string) (string, *gitlab.BasicMergeRequest, error) {
	state := "all"
	orderBy := "updated_at"
	sort := "desc"

	var bestProject string
	var best *gitlab.BasicMergeRequest
	var bestUpdated time.Time

	for _, pid := range projects {
		pid = strings.TrimSpace(pid)
		if pid == "" {
			continue
		}
		mrs, _, err := client.MergeRequests.ListProjectMergeRequests(pid, &gitlab.ListProjectMergeRequestsOptions{
			ListOptions: gitlab.ListOptions{Page: 1, PerPage: 1},
			State:       gitlab.Ptr(state),
			SourceBranch: gitlab.Ptr(branch),
			OrderBy:     gitlab.Ptr(orderBy),
			Sort:        gitlab.Ptr(sort),
		})
		if err != nil {
			return "", nil, fmt.Errorf("list merge requests for project %q: %w", pid, err)
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
			bestProject = pid
		}
	}

	if best == nil {
		return "", nil, fmt.Errorf("no merge request found for source branch %q in configured projects", branch)
	}
	return bestProject, best, nil
}

func mergeRequestFilePaths(client *gitlab.Client, project string, mergeRequestIID int) ([]string, error) {
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
		return d.OldPath
	}
	return d.NewPath
}

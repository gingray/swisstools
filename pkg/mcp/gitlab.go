package mcp

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
	"github.com/gingray/swisstools/pkg/service"
)

type GitlabChanges interface {
	ChangedFilesForBranch(filter service.BranchFilter) (*service.BranchChanges, error)
}
type Gitlab struct {
	cfg    *common.Config
	gitlab GitlabChanges
}

func NewGitlab(cfg *common.Config) (*Gitlab, error) {
	gitlab, err := service.NewGitlab(&cfg.GitLab, log.New(os.Stdout))
	if err != nil {
		return nil, err
	}
	return &Gitlab{cfg: cfg, gitlab: gitlab}, nil
}

// ChangedFilesForBranch finds the most recently updated merge request whose source branch matches branch
// across cfg.Projects, then lists all changed file paths from that MR’s diffs.
func (g *Gitlab) ChangedFilesForBranch(branch string) (*service.BranchChanges, error) {
	return g.gitlab.ChangedFilesForBranch(service.BranchFilter{Branch: branch, User: g.cfg.MCP.GitLabUser})
}

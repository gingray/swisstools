package workflow

import (
	"fmt"
	"github.com/go-git/go-git/v6"
	"os"
)

type TemplateVars struct {
	GitCurrentBranch string
	GitRepoUrl       string
	CurrentDir       string
}

func CreateTemplateVars() (*TemplateVars, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	templateVars := &TemplateVars{CurrentDir: currentDir}
	err = fillGitVars(templateVars)
	if err != nil {
		return nil, err
	}
	return templateVars, nil
}

func fillGitVars(templateVars *TemplateVars) error {
	repo, err := git.PlainOpenWithOptions(templateVars.CurrentDir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return fmt.Errorf("failed to open git repo: %w", err)
	}
	ref, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	templateVars.GitCurrentBranch = ref.Name().Short()
	remote, err := repo.Remote("origin")
	if err != nil {
		return fmt.Errorf("failed to get remote origin: %w", err)
	}
	cfg := remote.Config()
	if len(cfg.URLs) == 0 {
		return fmt.Errorf("no remote origin found")
	}
	templateVars.GitRepoUrl = cfg.URLs[0]
	return nil
}

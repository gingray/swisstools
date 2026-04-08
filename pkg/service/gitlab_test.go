package service

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func TestGitlab_ChangedFilesForBranch(t *testing.T) {
	type fields struct {
		cfg    *common.GitLabConfig
		logger *log.Logger
		client *gitlab.Client
	}
	type args struct {
		filter BranchFilter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *BranchChanges
		wantErr bool
	}{
		{
			name: "Test success request",
			fields: fields{
				cfg: &common.GitLabConfig{Url: "https://gitlab.com", ApiToken: "",
					Projects: []string{"gingray/swisstools"}, Authors: []string{"gingray"}},
				logger: log.Default(),
				client: nil,
			},
			args: args{
				BranchFilter{Branch: "master"},
			},
			want: &BranchChanges{
				Branch:       "",
				ChangedFiles: []string{"README.md", "CHANGELOG.md"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Gitlab{
				cfg:    tt.fields.cfg,
				logger: tt.fields.logger,
				client: tt.fields.client,
			}
			got, err := g.ChangedFilesForBranch(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangedFilesForBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChangedFilesForBranch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

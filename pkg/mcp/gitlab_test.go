package mcp

import (
	"reflect"
	"testing"

	"github.com/gingray/swisstools/pkg/common"
)

func TestGitlab_ChangedFilesForBranch(t *testing.T) {
	type fields struct {
		cfg *common.Config
	}
	type args struct {
		branch string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *BranchMergeRequestFiles
		wantErr bool
	}{
		{
			name: "fetch changed files for branch",
			fields: fields{
				cfg: &common.Config{
					MCP: common.MCPConfig{
						GitLabUser: "test-user",
					},
				},
			},
			args: args{
				branch: "test-branch",
			},
			want: &BranchMergeRequestFiles{
				Files: []string{"README.md"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Gitlab{
				cfg: tt.fields.cfg,
			}
			got, err := g.ChangedFilesForBranch(tt.args.branch)
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

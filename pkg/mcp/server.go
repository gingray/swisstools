package mcp

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gingray/swisstools/pkg/common"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPServer struct {
	cfg *common.Config
}

type BranchParams struct {
	Branch string `json:"branch" jsonschema:"git branch name"`
}

func NewMCPServer(cfg *common.Config) *MCPServer {
	return &MCPServer{cfg: cfg}
}

func (s *MCPServer) Run() error {
	opts := mcp.ServerOptions{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "DevTool",
		Version: "v0.0.1",
	}, &opts)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetFilesChanged",
		Description: "get changed files for particular branch",
	}, s.GitLabChanges)
	handler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		return server
	}, nil)
	return http.ListenAndServe(":8080", handler)
}

func (s *MCPServer) GitLabChanges(ctx context.Context, req *mcp.CallToolRequest, args BranchParams) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "file1,file2,file3"},
		},
	}, nil, nil
}

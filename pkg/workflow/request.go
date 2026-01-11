package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gingray/swisstools/pkg/common"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type WorkflowRequest struct {
	WorkflowName string
	config       *common.Config
}

func NewWorkflowRequest(workflowName string, config *common.Config) *WorkflowRequest {
	return &WorkflowRequest{WorkflowName: workflowName, config: config}
}

func (r *WorkflowRequest) MakeRequest(keyValues []string) error {
	ctx := context.Background()
	workflow := r.config.Workflows[r.WorkflowName]
	payload := make(map[string]string)
	templateVars, err := CreateTemplateVars()
	if err != nil {
		return err
	}
	for _, arg := range workflow.PredefinedArgs {
		templ, err := template.New(arg.Key).Parse(arg.Value)
		if err != nil {
			return err
		}
		var buffer bytes.Buffer
		err = templ.Execute(&buffer, templateVars)
		if err != nil {
			return err
		}
		payload[arg.Key] = buffer.String()
	}

	for _, item := range keyValues {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid data %q", item)
		}
		payload[parts[0]] = parts[1]
	}
	log.Infof("Template vars: %+v", payload)
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, workflow.Method, workflow.Endpoint, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for _, item := range workflow.Headers {
		req.Header.Set(item.Key, item.Value)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("http %d: %s", resp.StatusCode, b)
	}
	return nil
}

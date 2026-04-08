package initialize

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gingray/swisstools/pkg/common"
	"github.com/gingray/swisstools/test"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)
import "github.com/stretchr/testify/assert"

func TestConfigPathCombine(t *testing.T) {
	assertion := assert.New(t)
	home := "/home/user"
	configPath := getConfigPath(home)
	assertion.Equal("/home/user/.swisstools", configPath.ConfigDir)
	assertion.Equal("/home/user/.swisstools/config.yaml", configPath.FullPath)
}

func TestPartialConfigUnmarshall(t *testing.T) {
	assertion := assert.New(t)
	configFile := "partial_config.yaml"
	configPath := filepath.Join(test.RootPath, "fixtures", configFile)
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	assertion.NoError(err)
	cfg := &common.Config{}
	cfg, err = unmarshallConfig(cfg)
	assertion.NoError(err)
	assertion.NotNil(cfg)
	assertion.Equal("CODE", cfg.Jira.Project)
}
func TestFullConfigUnmarshall(t *testing.T) {
	assertion := assert.New(t)
	configFile := "full_config.yaml"
	configPath := filepath.Join(test.RootPath, "fixtures", configFile)
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	assertion.NoError(err)
	cfg := &common.Config{}
	cfg, err = unmarshallConfig(cfg)
	assertion.NoError(err)
	assertion.NotNil(cfg)

	assertion.Equal("CODE", cfg.Jira.Project)
	assertion.Equal([]string{"alex.rows"}, cfg.GitLab.Authors)

}

func TestFullConfigInFixtureDirectory(t *testing.T) {
	t.Skip("skip config file creation")

	assertion := assert.New(t)
	configFile := "full_config.yaml"
	configPath := filepath.Join(test.RootPath, "fixtures", configFile)
	err := viper.SafeWriteConfigAs(configPath)
	assertion.NoError(err)
	emptyConfig := common.Config{Workflows: map[string]common.WorkflowConfig{
		"workflow-1": {
			Endpoint:       "http://localhost:8080",
			Method:         "POST",
			Headers:        []common.WorkflowHeader{{Key: "X-Token", Value: "my-secret-token"}},
			PredefinedArgs: []common.PredefinedArg{{Key: "branch", Value: "{{.BranchName}}"}}},
		"workflow-2": {
			Endpoint:       "http://localhost:8080",
			Method:         "POST",
			PredefinedArgs: []common.PredefinedArg{{Key: "currentDir", Value: "{{.CurrentDir}}"}},
		}}}
	data, _ := yaml.Marshal(emptyConfig)
	err = os.WriteFile(configPath, data, 0644)
	assertion.NoError(err)
}

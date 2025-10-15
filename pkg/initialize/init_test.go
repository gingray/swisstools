package initialize

import "testing"
import "github.com/stretchr/testify/assert"

func TestConfigPathCombine(t *testing.T) {
	assertion := assert.New(t)
	home := "/home/user"
	configPath := getConfigPath(home)
	assertion.Equal("/home/user/.swisstools", configPath.ConfigDir)
	assertion.Equal("/home/user/.swisstools/config.yaml", configPath.FullPath)
}

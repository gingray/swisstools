package common

const ConfigFile = "config"
const ConfigType = "yaml"
const ConfigDir = ".swisstools"

type Config struct {
	Jira   JiraConfig   `yaml:"jira" mapstructure:"jira"`
	GitLab GitLabConfig `yaml:"gitlab" mapstructure:"gitlab"`
}

type JiraConfig struct {
	ApiToken string `yaml:"apiToken" mapstructure:"apiToken"`
	Url      string `yaml:"url" mapstructure:"url"`
	Project  string `yaml:"project" mapstructure:"project"`
}

type GitLabConfig struct {
	Token    string   `yaml:"token" mapstructure:"token"`
	Url      string   `yaml:"url" mapstructure:"url"`
	Authors  []string `yaml:"authors" mapstructure:"authors"`
	Projects []string `yaml:"projects" mapstructure:"projects"`
}

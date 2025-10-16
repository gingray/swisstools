package common

const ConfigFile = "config"
const ConfigType = "yaml"
const ConfigDir = ".swisstools"

type Config struct {
	Jira   JiraConfig   `yaml:"jira" mapstructure:"jira"`
	GitLab GitLabConfig `yaml:"gitlab" mapstructure:"gitlab"`
	Sentry SentryConfig `yaml:"sentry" mapstructure:"sentry"`
}

type JiraConfig struct {
	ApiToken string `yaml:"apiToken" mapstructure:"apiToken"`
	Url      string `yaml:"url" mapstructure:"url"`
	Project  string `yaml:"project" mapstructure:"project"`
}

type GitLabConfig struct {
	ApiToken string   `yaml:"apiToken" mapstructure:"apiToken"`
	Url      string   `yaml:"url" mapstructure:"url"`
	Authors  []string `yaml:"authors" mapstructure:"authors"`
	Projects []string `yaml:"projects" mapstructure:"projects"`
}

type SentryConfig struct {
	ApiToken     string `yaml:"apiToken" mapstructure:"apiToken"`
	Url          string `yaml:"url" mapstructure:"url"`
	Organization string `yaml:"organization" mapstructure:"organization"`
	Project      string `yaml:"project" mapstructure:"project"`
}

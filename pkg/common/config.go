package common

const ConfigFile = "config"
const ConfigType = "yaml"
const ConfigDir = ".swisstools"

type Config struct {
	Jira JiraConfig `yaml:"jira" mapstructure:"jira"`
}

type JiraConfig struct {
	ApiToken string `yaml:"apiToken" mapstructure:"apiToken"`
	Url      string `yaml:"url" mapstructure:"url"`
	Project  string `yaml:"project" mapstructure:"project"`
}

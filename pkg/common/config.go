package common

const ConfigFile = "config"
const ConfigType = "yaml"
const ConfigDir = ".swisstools"

type Config struct {
	Jira      JiraConfig                `yaml:"jira" mapstructure:"jira"`
	GitLab    GitLabConfig              `yaml:"gitlab" mapstructure:"gitlab"`
	Sentry    SentryConfig              `yaml:"sentry" mapstructure:"sentry"`
	Api       ApiConfig                 `yaml:"api" mapstructure:"api"`
	Workflows map[string]WorkflowConfig `yaml:"workflow" mapstructure:"workflow"`
	MCP       MCPConfig                 `yaml:"mcp" mapstructure:"mcp"`
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
	Tag          string `yaml:"tag" mapstructure:"tag"`
}

type ApiConfig struct {
	Port      int    `yaml:"port" mapstructure:"port"`
	SecretKey string `yaml:"secretKey" mapstructure:"secretKey"`
}

type WorkflowConfig struct {
	Endpoint       string           `yaml:"endpoint" mapstructure:"endpoint"`
	Method         string           `yaml:"method" mapstructure:"method"`
	Headers        []WorkflowHeader `yaml:"headers" mapstructure:"headers"`
	PredefinedArgs []PredefinedArg  `yaml:"predefinedArgs" mapstructure:"predefinedArgs"`
}

type PredefinedArg struct {
	Key   string `yaml:"key" mapstructure:"key"`
	Value string `yaml:"value" mapstructure:"value"`
}

type WorkflowHeader struct {
	Key   string `json:"key" mapstructure:"key"`
	Value string `json:"value" mapstructure:"value"`
}

type MCPConfig struct {
	Port int `yaml:"port" mapstructure:"port"`
}

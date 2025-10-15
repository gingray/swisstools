package initialize

import (
	"errors"
	"fmt"
	"github.com/gingray/swisstools/pkg/common"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type configPath struct {
	FullPath  string
	ConfigDir string
}

func CreateConfigIfNotExists() (string, error) {
	home, err := os.UserHomeDir()
	configPath := getConfigPath(home)
	viper.AddConfigPath(configPath.ConfigDir)
	viper.SetConfigName(common.ConfigFile)
	viper.SetConfigType(common.ConfigType)

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			err := os.MkdirAll(configPath.ConfigDir, 0755)
			if err != nil {
				return "", err
			}
			err = viper.SafeWriteConfigAs(configPath.FullPath)
			emptyConfig := common.Config{}
			data, _ := yaml.Marshal(emptyConfig)
			err = os.WriteFile(configPath.FullPath, data, 0644)
			return "", err
		} else {
			return "", err
		}
	}

	cfg := common.Config{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return "", err
	}
	out, err := yaml.Marshal(&cfg)
	if err != nil {
		return "", err
	}
	return string(out), err
}

func getConfigPath(home string) configPath {
	configDir := filepath.Join(home, common.ConfigDir)
	configFullPath := filepath.Join(configDir, fmt.Sprintf("%s.%s", common.ConfigFile, common.ConfigType))
	return configPath{
		FullPath:  configFullPath,
		ConfigDir: configDir,
	}
}

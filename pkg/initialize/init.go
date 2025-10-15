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

func CreateConfigIfNotExists() (string, error) {
	home, err := os.UserHomeDir()
	configPath := filepath.Join(home, common.ConfigDir)
	viper.AddConfigPath(configPath)
	viper.SetConfigName(common.ConfigFile)
	viper.SetConfigType(common.ConfigType)

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			err := os.MkdirAll(configPath, 0755)
			if err != nil {
				return "", err
			}
			configFullPath := filepath.Join(configPath, fmt.Sprintf("%s.%s", common.ConfigFile, common.ConfigType))
			err = viper.SafeWriteConfigAs(configFullPath)
			emptyConfig := common.Config{}
			data, _ := yaml.Marshal(emptyConfig)
			err = os.WriteFile(configFullPath, data, 0644)
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

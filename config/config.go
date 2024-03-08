package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Target      string
	AccessToken string
	CodeMacro   map[string]string
}

func LoadConfig() (Config, error) {
	var config Config
	var configPath string

	dir, err := os.UserConfigDir()
	if err != nil {
		return config, err
	}

	configPath = filepath.Join(dir, "go-confluence")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err = viper.ReadInConfig(); err != nil {
		return config, err
	}

	if err = viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil

}

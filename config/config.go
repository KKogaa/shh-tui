package config

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"server"`
	Client struct {
		Chat struct {
			Width  int `mapstructure:"width"`
			Height int `mapstructure:"height"`
		} `mapstructure:"chat"`
		Chatbox struct {
			Width  int `mapstructure:"width"`
			Height int `mapstructure:"height"`
		} `mapstructure:"chatbox"`
		Username string `mapstructure:"username"`
		Chatroom string `mapstructure:"chatroom"`
	} `mapstructure:"client"`
}

func LoadConfig() (*Config, error) {
	homeDir, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %w", err)
	}

	configFile := filepath.Join(homeDir.HomeDir, ".config/shh/config.yaml")

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configFile)
	viper.AutomaticEnv()
	viper.SetDefault("server.url", "ws://localhost:8080/ws")
	viper.SetDefault("client.username", "noname")
	viper.SetDefault("client.chatroom", "default")
	viper.SetDefault("client.chat.width", 200)
	viper.SetDefault("client.chat.height", 5)
	viper.SetDefault("client.chatbox.height", 200)
	viper.SetDefault("client.chatbox.height", 3)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &config, nil
}

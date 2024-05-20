package config

import (
	"errors"
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
	"golang.org/x/term"
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

func GetTerminalWidth() (int, error) {
	if !term.IsTerminal(0) {
		return 0, errors.New("error not executing on a terminal")
	}

	width, _, err := term.GetSize(0)
	if err != nil {
		return 0, errors.New("error obtaining width of terminal")
	}

	return width, nil
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

	defaultWidth, err := GetTerminalWidth()
	if err != nil {
		return nil, err
	}

	viper.SetDefault("client.chat.width", defaultWidth)
	viper.SetDefault("client.chat.height", 5)
	viper.SetDefault("client.chatbox.width", defaultWidth)
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

package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Client ClientConfig `yaml:"client"`
}

type ServerConfig struct {
	//Host string `yaml:"host"`
	//Port string `yaml:"port"`
	Address string `yaml:"address"`
	//ReadTimeoutSec int `yaml:"read_timeout_sec"`
	//WriteTimeoutSec int `yaml:"write_timeout_sec"`
}

type ClientConfig struct {
	//Host string `yaml:"host"`
	//Port string `yaml:"port"`
	Address string `yaml:"address"`
	//ReadTimeoutSec int `yaml:"read_timeout_sec"`
	//WriteTimeoutSec int `yaml:"write_timeout_sec"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

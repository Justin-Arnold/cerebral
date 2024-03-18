package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Services []Service `toml:"services"`
}

type Service struct {
	Name string `toml:"name"`
	URL  string `toml:"url"`
}

func LoadConfig(filename string) (*Config, error) {
	var config Config
	_, decodeError := toml.DecodeFile(filename, &config)
	if decodeError != nil {
		return nil, decodeError
	}

	return &config, nil
}

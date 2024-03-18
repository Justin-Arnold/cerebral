package config

import (
	"os"

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

func UpdateConfig(filename, name, url string) (*Config, error) {
	// 1. Load existing configuration
	config, err := LoadConfig(filename)
	if err != nil {
		return nil, err
	}

	// 2. Create the new service
	newService := Service{Name: name, URL: url}

	// 3. Append the new service to the slice
	config.Services = append(config.Services, newService)

	// 4. Open the file in write mode (overwriting existing content)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 5. Encode and write the updated configuration directly
	encoder := toml.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

package config

import (
	"os"
	"slices"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Services []Service `toml:"services"`
}

type Service struct {
	Name string `toml:"name"`
	URL  string `toml:"url"`
}

type EditServiceData struct {
	OldName string
	Name    string
	URL     string
}

func CreateNewConfig(filePath string) (*os.File, error) {
	// Create File
	configFile, createFileError := os.Create(filePath)
	if createFileError != nil {
		return nil, createFileError
	}

	// Define Blank Config
	config := Config{
		Services: []Service{},
	}

	// Write initial configuration to file using TOML encoding
	if err := toml.NewEncoder(configFile).Encode(config); err != nil {
		return nil, err
	}

	// Ensure to close the file after writing
	if err := configFile.Close(); err != nil {
		return nil, err
	}

	return configFile, nil
}

func LoadConfig(filePath string) (*Config, error) {
	var config Config
	_, decodeError := toml.DecodeFile(filePath, &config)
	if decodeError != nil {
		return nil, decodeError
	}

	return &config, nil
}

func WriteConfig(filePath string, config *Config) error {
	// 4. Open the file in write mode (overwriting existing content)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 5. Encode and write the updated configuration directly
	encoder := toml.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}

func AddServiceToConfig(filename, name, url string) (*Config, error) {
	// 1. Load existing configuration
	config, err := LoadConfig(filename)
	if err != nil {
		return nil, err
	}

	// 2. Create the new service
	newService := Service{Name: name, URL: url}

	// 3. Append the new service to the slice
	config.Services = append(config.Services, newService)

	// Write
	WriteConfig(filename, config)

	return config, nil
}

func EditServiceInConfig(filePath string, data EditServiceData) (*Config, error) {
	// Load in current config
	config, configError := LoadConfig(filePath)
	if configError != nil {
		return nil, configError
	}

	// Get matching service
	serviceIndex := slices.IndexFunc(config.Services, func(service Service) bool { return service.Name == data.OldName })
	if serviceIndex == -1 {
		return nil, nil
	}

	// Update Existing Service Info
	updatedService := config.Services[serviceIndex]
	updatedService.Name = data.Name
	updatedService.URL = data.URL

	config.Services[serviceIndex] = updatedService

	writeError := WriteConfig(filePath, config)
	if writeError != nil {
		return nil, writeError
	}

	return config, nil
}

func DeleteServiceFromConfig(fileName string, name string) (*Config, error) {
	// Load in current config
	config, configError := LoadConfig(fileName)
	if configError != nil {
		return nil, configError
	}

	// Get matching service
	serviceIndex := slices.IndexFunc(config.Services, func(service Service) bool { return service.Name == name })
	if serviceIndex == -1 {
		return nil, nil
	}

	// Remove the service
	config.Services = append(config.Services[:serviceIndex], config.Services[serviceIndex+1:]...)

	// Open the file in write mode (overwriting existing content)
	WriteConfig(fileName, config)

	return config, nil
}

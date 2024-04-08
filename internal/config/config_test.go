package config

import (
	"os"
	"testing"
)

func TestCreateNewConfig(t *testing.T) {
	// Create the new file
	const fileName = "config-test.toml"
	configFile, creationError := CreateNewConfig(fileName)
	defer os.Remove(configFile.Name())
	if creationError != nil {
		t.Fatalf("Error creating new config: %v", creationError)
	}

	// Load the config from file
	config, err := LoadConfig(fileName)
	if err != nil {
		t.Fatalf("Error loading new config")
	}

	// Test that services are empty
	if len(config.Services) != 0 {
		t.Fatalf("New config services wrong length")
	}
}

func TestAddServiceToBlankFile(t *testing.T) {
	// Create test config file
	fileName := "config-test.toml"
	configFile, creationError := CreateNewConfig(fileName)
	defer os.Remove(configFile.Name())
	if creationError != nil {
		t.Fatalf("Error creating new config: %v", creationError)
	}

	// Define test service data
	const newServiceName = "New Service"
	const newServiceUrl = "0.0.0.0"
	updatedConfig, addError := AddServiceToConfig(fileName, newServiceName, newServiceUrl)
	if addError != nil {
		t.Fatalf("Error adding to blank config: %v", addError)
	}
	_ = updatedConfig

	// Load config from file
	config, loadError := LoadConfig(fileName)
	if loadError != nil {
		t.Fatalf("Error loading config: %v", loadError)
	}

	// Check that the new service has the correct name
	if config.Services[0].Name != newServiceName {
		t.Fatalf("Service Name is incorrect")
	}

	// Check that the new service has the correct URL
	if config.Services[0].URL != newServiceUrl {
		t.Fatalf("Serice URL is incorrect")
	}

	// Check that the new service is the only service
	if len(config.Services) != 1 {
		t.Fatalf("Unexpected config services length: %v", len(config.Services))
	}
}

func TestAddSeriviceToPopulatedFile(t *testing.T) {
	fileName := "config-test.toml"
	configFile, creationError := CreateNewConfig(fileName)
	defer os.Remove(configFile.Name())
	if creationError != nil {
		t.Fatalf("Error creating new config: %v", creationError)
	}
	config := Config{
		Services: []Service{
			{
				Name: "Old Service",
				URL:  "0.0.0.0",
			},
		},
	}
	writeError := WriteConfig(fileName, &config)
	if writeError != nil {
		t.Fatalf("Error writing initial data to config: %v", writeError)
	}

	// Define test service data
	const newServiceName = "New Service"
	const newServiceUrl = "1.1.1.1"

	// Add service to config
	updatedConfig, addError := AddServiceToConfig(fileName, newServiceName, newServiceUrl)
	if addError != nil {
		t.Fatalf("Error adding to populated config: %v", addError)
	}
	_ = updatedConfig

	// Load config from file
	loadedConfig, loadError := LoadConfig(fileName)
	if loadError != nil {
		t.Fatalf("Error loading config: %v", loadError)
	}

	// Check that the new service has the correct name
	if loadedConfig.Services[1].Name != newServiceName {
		t.Fatalf("Service Name is incorrect")
	}

	// Check that the new service has the correct URL
	if loadedConfig.Services[1].URL != newServiceUrl {
		t.Fatalf("Serice URL is incorrect")
	}

	// Check that there are two services
	if len(loadedConfig.Services) != 2 {
		t.Fatalf("Unexpected config services length: %v", len(config.Services))
	}
}

func TestEditServiceInConfig(t *testing.T) {
	// Setup: Create a temp file and write initial config to it
	fileName := "config-test.toml"
	tempFile := CreateTempFile(t, fileName)
	defer os.Remove(tempFile.Name())

	// Initial configuration to write to the temp file
	config := Config{
		Services: []Service{
			{
				Name: "test service before",
				URL:  "0.0.0.0",
			},
		},
	}

	writeError := WriteConfig(fileName, &config)
	if writeError != nil {
		t.Fatalf("Error while writing initial config to file: %v", writeError)
	}

	// Test: Call EditServiceInConfig with the path to the temp file
	editData := EditServiceData{
		OldName: config.Services[0].Name,
		Name:    "test service after",
		URL:     "1.1.1.1",
	}

	editedConfig, editError := EditServiceInConfig(fileName, editData)
	if editError != nil {
		t.Fatalf("failed to edit temp config %v", editError)
	}

	if len(editedConfig.Services) != 1 {
		t.Fatalf("Invalid Service Count")
	}

	if editedConfig.Services[0].Name != editData.Name {
		t.Fatalf("Failed to edit name")
	}

	if editedConfig.Services[0].URL != editData.URL {
		t.Fatalf("Failed to edit Url")
	}
}

func TestDeleteServiceFromConfig(t *testing.T) {
	// Setup: Create a temp file and write initial config to it
	fileName := "config-test.toml"
	tempFile := CreateTempFile(t, fileName)
	defer os.Remove(tempFile.Name())

	// Initial configuration to write to the temp file
	config := Config{
		Services: []Service{
			{
				Name: "test service 1",
				URL:  "0.0.0.0",
			},
			{
				Name: "test service 2",
				URL:  "1.1.1.1",
			},
		},
	}

	writeError := WriteConfig(fileName, &config)
	if writeError != nil {
		t.Fatalf("Error while writing initial config to file: %v", writeError)
	}

	updatedConfig, deleteError := DeleteServiceFromConfig(fileName, config.Services[0].Name)
	if deleteError != nil {
		t.Fatalf("Error deleting service from config: %v", deleteError)
	}

	if len(updatedConfig.Services) != 1 {
		t.Fatalf("Failed to delete service")
	}

	if updatedConfig.Services[0].Name != config.Services[1].Name {
		t.Fatalf("Deleted incorrect service")
	}
}

// -----------------
// Helper Functions
// -----------------

func CreateTempFile(t *testing.T, fileName string) *os.File {
	tempFile, createTempFileError := os.Create(fileName)

	if createTempFileError != nil {
		t.Fatalf("Failed to create temp file: %v", createTempFileError)
	}

	return tempFile
}

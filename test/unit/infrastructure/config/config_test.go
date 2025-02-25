package config_test

import (
	"os"
	"testing"

	"assessment/infrastructure/config"
)

func TestNewConfig(t *testing.T) {

	cfg := config.NewConfig()

	if len(cfg.DisabledSorters) != 0 {
		t.Errorf("Default DisabledSorters not empty: got %v", cfg.DisabledSorters)
	}

	if cfg.DefaultPageSize != 10 {
		t.Errorf("Default DefaultPageSize mismatch: got %d, want %d", cfg.DefaultPageSize, 10)
	}
}

func TestConfigLoadAndSave(t *testing.T) {
	// Create a temporary file in the current directory
	tempFile := "temp_config_test.json"

	// Ensure cleanup after test
	defer os.Remove(tempFile)

	// Create a new config
	cfg := config.NewConfig()
	cfg.DisabledSorters = []string{"Sorter1", "Sorter2"}
	cfg.DefaultPageSize = 20

	// Save the config to file
	err := cfg.SaveToFile(tempFile)
	if err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	// Verify the file was created
	if _, statErr := os.Stat(tempFile); os.IsNotExist(statErr) {
		t.Fatal("Config file was not created")
	}

	// Create a new config to load into
	loadedCfg := config.NewConfig()

	// Load the config from file
	err = loadedCfg.LoadFromFile(tempFile)
	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	// Verify the loaded config matches the original
	if len(loadedCfg.DisabledSorters) != 2 {
		t.Errorf("Loaded DisabledSorters length mismatch: got %d, want %d", len(loadedCfg.DisabledSorters), 2)
	}

	if loadedCfg.DisabledSorters[0] != "Sorter1" || loadedCfg.DisabledSorters[1] != "Sorter2" {
		t.Errorf("Loaded DisabledSorters mismatch: got %v, want %v", loadedCfg.DisabledSorters, []string{"Sorter1", "Sorter2"})
	}

	if loadedCfg.DefaultPageSize != 20 {
		t.Errorf("Loaded DefaultPageSize mismatch: got %d, want %d", loadedCfg.DefaultPageSize, 20)
	}
}

func TestConfigLoadNonExistentFile(t *testing.T) {

	cfg := config.NewConfig()

	err := cfg.LoadFromFile("non-existent-file.json")

	if err == nil {
		t.Error("LoadFromFile did not return error for non-existent file")
	}
}

func TestConfigSaveToInvalidPath(t *testing.T) {

	cfg := config.NewConfig()

	err := cfg.SaveToFile("/invalid/path/config.json")

	if err == nil {
		t.Error("SaveToFile did not return error for invalid path")
	}
}

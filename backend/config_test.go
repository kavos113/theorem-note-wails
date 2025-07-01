package backend

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadAndSaveProjectConfig(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test loading from a non-existent config
	config, err := LoadProjectConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadProjectConfig failed for non-existent config: %v", err)
	}
	if !reflect.DeepEqual(config, getDefaultProjectConfig()) {
		t.Errorf("Expected default config for non-existent config, but got %v", config)
	}

	// Define the config data to be saved
	expectedConfig := ProjectConfig{
		FontSettings: FontSettings{
			EditorFontFamily:  "test-font",
			EditorFontSize:    20,
			PreviewFontFamily: "test-preview-font",
			PreviewFontSize:   22,
		},
	}

	// Test saving the config
	err = SaveProjectConfig(tmpDir, expectedConfig)
	if err != nil {
		t.Fatalf("SaveProjectConfig failed: %v", err)
	}

	// Verify that the config file was created
	configFilePath, _ := getProjectConfigPath(tmpDir)
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created")
	}

	// Test loading the config
	loadedConfig, err := LoadProjectConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadProjectConfig failed: %v", err)
	}

	// Compare the loaded data with the original data
	if !reflect.DeepEqual(loadedConfig, expectedConfig) {
		t.Errorf("Loaded config data does not match expected data.\nGot:  %v\nWant: %v", loadedConfig, expectedConfig)
	}
}

func TestGetProjectConfigPath(t *testing.T) {
	// Test with a valid root directory
	rootDir := "/some/dir"
	expectedPath := filepath.Join(rootDir, projectConfigDir, projectConfigFile)
	actualPath, err := getProjectConfigPath(rootDir)
	if err != nil {
		t.Errorf("getProjectConfigPath failed with valid rootDir: %v", err)
	}
	if actualPath != expectedPath {
		t.Errorf("Expected path %s, but got %s", expectedPath, actualPath)
	}

	// Test with an empty root directory
	_, err = getProjectConfigPath("")
	if err != os.ErrInvalid {
		t.Errorf("Expected os.ErrInvalid for empty rootDir, but got %v", err)
	}
}

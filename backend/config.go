package backend

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	LastOpened string `json:"last_opened"`
}

type ConfigManager struct {
	config Config
	path   string
}

const (
	configDirName  = "theorem-note-wails"
	configFileName = "config.json"
)

func NewConfigManager() *ConfigManager {
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		panic("Failed to get user config directory: " + err.Error())
	}

	path := filepath.Join(appDataDir, configDirName, configFileName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			panic("Failed to create config directory: " + err.Error())
		}

		// Initialize default config
		defaultConfig := Config{
			LastOpened: "",
		}
		file, err := os.Create(path)
		if err != nil {
			panic("Failed to create config file: " + err.Error())
		}
		defer file.Close()
		b, err := json.Marshal(defaultConfig)
		if err != nil {
			panic("Failed to marshal default config: " + err.Error())
		}
		_, err = file.Write(b)
		if err != nil {
			panic("Failed to write default config: " + err.Error())
		}
	} else if err != nil {
		panic("Failed to check config file: " + err.Error())
	}

	// Load existing config
	configManager := &ConfigManager{}

	file, err := os.ReadFile(path)
	if err != nil {
		panic("Failed to read config file: " + err.Error())
	}
	if err := json.Unmarshal(file, &configManager.config); err != nil {
		panic("Failed to unmarshal config file: " + err.Error())
	}
	configManager.path = path
	return configManager
}

func (cm *ConfigManager) GetConfig() Config {
	return cm.config
}

func (cm *ConfigManager) saveConfig() error {
	file, err := os.Create(cm.path)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := json.Marshal(cm.config)
	if err != nil {
		return err
	}

	_, err = file.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (cm *ConfigManager) SetLastOpened(path string) {
	cm.config.LastOpened = path
	if err := cm.saveConfig(); err != nil {
		panic("Failed to save config: " + err.Error())
	}
}

func (cm *ConfigManager) GetLastOpened() string {
	return cm.config.LastOpened
}

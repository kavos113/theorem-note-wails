package backend

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// =================================================================
// グローバル設定 (最後に開いたディレクトリなど)
// =================================================================

const (
	globalConfigDirName  = "theorem-note-wails"
	globalConfigFileName = "global_config.json"
)

// GlobalConfig はアプリケーション全体のグローバル設定を保持します
type GlobalConfig struct {
	LastOpenedPath string `json:"last_opened_path"`
}

// ConfigManager はグローバル設定を管理します
type ConfigManager struct {
	config GlobalConfig
	path   string
}

// NewConfigManager は新しいConfigManagerを初期化します
func NewConfigManager() *ConfigManager {
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		panic("Failed to get user config directory: " + err.Error())
	}

	path := filepath.Join(appDataDir, globalConfigDirName, globalConfigFileName)
	cm := &ConfigManager{path: path}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		cm.config = GlobalConfig{LastOpenedPath: ""}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			panic("Failed to create global config directory: " + err.Error())
		}
		cm.save()
	} else {
		file, err := os.ReadFile(path)
		if err != nil {
			panic("Failed to read global config file: " + err.Error())
		}
		if json.Unmarshal(file, &cm.config) != nil {
			// ファイルが壊れている場合などはデフォルトで初期化
			cm.config = GlobalConfig{LastOpenedPath: ""}
		}
	}
	return cm
}

func (cm *ConfigManager) save() {
	b, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		panic("Failed to marshal global config: " + err.Error())
	}
	if err := os.WriteFile(cm.path, b, 0644); err != nil {
		panic("Failed to write global config: " + err.Error())
	}
}

// GetLastOpened は最後に開いたディレクトリのパスを返します
func (cm *ConfigManager) GetLastOpened() string {
	return cm.config.LastOpenedPath
}

// SetLastOpened は最後に開いたディレクトリのパスを設定します
func (cm *ConfigManager) SetLastOpened(path string) {
	cm.config.LastOpenedPath = path
	cm.save()
}

// =================================================================
// プロジェクト設定 (フォント設定など)
// =================================================================

const (
	projectConfigDir  = ".theorem-note"
	projectConfigFile = "config.json"
)

// FontSettings はフォント設定を保持します
type FontSettings struct {
	EditorFontFamily  string `json:"editor_font_family"`
	EditorFontSize    int    `json:"editor_font_size"`
	PreviewFontFamily string `json:"preview_font_family"`
	PreviewFontSize   int    `json:"preview_font_size"`
}

// ProjectConfig はプロジェクトごとの設定を保持します
type ProjectConfig struct {
	FontSettings FontSettings `json:"font_settings"`
}

func getProjectConfigPath(rootDir string) (string, error) {
	if rootDir == "" {
		return "", os.ErrInvalid
	}
	return filepath.Join(rootDir, projectConfigDir, projectConfigFile), nil
}

func getDefaultProjectConfig() ProjectConfig {
	return ProjectConfig{
		FontSettings: FontSettings{
			EditorFontFamily:  "Consolas, Monaco, 'Courier New', monospace",
			EditorFontSize:    14,
			PreviewFontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif",
			PreviewFontSize:   14,
		},
	}
}

// LoadProjectConfig はプロジェクトの設定を読み込みます
func LoadProjectConfig(rootDir string) (ProjectConfig, error) {
	defaultConfig := getDefaultProjectConfig()
	if rootDir == "" {
		return defaultConfig, nil
	}

	path, err := getProjectConfigPath(rootDir)
	if err != nil {
		return defaultConfig, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return defaultConfig, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return defaultConfig, err
	}

	var config ProjectConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return defaultConfig, err
	}
	return config, nil
}

// SaveProjectConfig はプロジェクトの設定を保存します
func SaveProjectConfig(rootDir string, config ProjectConfig) error {
	if rootDir == "" {
		return os.ErrInvalid
	}

	// ディレクトリが存在することを確認
	if err := os.MkdirAll(filepath.Join(rootDir, projectConfigDir), 0755); err != nil {
		return err
	}

	path, err := getProjectConfigPath(rootDir)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

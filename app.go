package main

import (
	"context"
	"fmt"

	"github.com/kavos113/theorem-note-wails/backend"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	configManager *backend.ConfigManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.configManager = backend.NewConfigManager()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetNewDirectoryFileTree() ([]backend.FileItem, error) {
	path, items, err := backend.GetNewDirectoryFileTree(a.ctx)
	if err != nil {
		return nil, err
	}
	a.configManager.SetLastOpened(path)
	return items, nil
}

func (a *App) GetFileTree(path string) ([]backend.FileItem, error) {
	// Note: SetLastOpenedはGetNewDirectoryFileTreeでのみ呼び出すように変更
	return backend.GetFileTree(path)
}

func (a *App) ReadFile(path string) (string, error) {
	return backend.ReadFile(path)
}

func (a *App) WriteFile(path string, content string) error {
	return backend.WriteFile(path, content)
}

func (a *App) CreateFile(path string) error {
	return backend.CreateFile(path)
}

func (a *App) CreateDirectory(path string) error {
	return backend.CreateDirectory(path)
}

// GetLastOpened はグローバル設定から最後に開いたパスを取得します
func (a *App) GetLastOpened() string {
	return a.configManager.GetLastOpened()
}

// SetLastOpened はグローバル設定に最後に開いたパスを保存します
func (a *App) SetLastOpened(path string) {
	a.configManager.SetLastOpened(path)
}

// --- Session Management ---

func (a *App) SaveSession(rootDir string, filePaths []string) error {
	return backend.SaveSession(rootDir, filePaths)
}

func (a *App) LoadSession(rootDir string) ([]string, error) {
	return backend.LoadSession(rootDir)
}

// --- Project Settings ---

func (a *App) GetFontSettings(rootDir string) (backend.FontSettings, error) {
	config, err := backend.LoadProjectConfig(rootDir)
	if err != nil {
		return backend.FontSettings{}, err
	}
	return config.FontSettings, nil
}

func (a *App) SaveFontSettings(rootDir string, settings backend.FontSettings) error {
	// 現在の設定を読み込み、フォント設定のみを更新して保存する
	config, err := backend.LoadProjectConfig(rootDir)
	if err != nil {
		return err
	}

	config.FontSettings = settings

	if err := backend.SaveProjectConfig(rootDir, config); err != nil {
		return err
	}

	// フロントエンドに設定が更新されたことを通知
	runtime.EventsEmit(a.ctx, "font-settings-updated", settings)

	return nil
}

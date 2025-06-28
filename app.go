package main

import (
	"context"
	"fmt"

	"github.com/kavos113/theorem-note-wails/backend"
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
	a.configManager.SetLastOpened(path)
	return backend.GetFileTree(path)
}

func (a *App) ReadFile(path string) (string, error) {
	return backend.ReadFile(path)
}

func (a *App) WriteFile(path string, content string) error {
	return backend.WriteFile(path, content)
}

func (a *App) GetLastOpened() string {
	return a.configManager.GetLastOpened()
}

func (a *App) SetLastOpened(path string) {
	a.configManager.SetLastOpened(path)
}

func (a *App) SaveSession(rootDir string, filePaths []string) error {
	return backend.SaveSession(rootDir, filePaths)
}

func (a *App) LoadSession(rootDir string) ([]string, error) {
	return backend.LoadSession(rootDir)
}

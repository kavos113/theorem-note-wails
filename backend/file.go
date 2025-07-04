package backend

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"path/filepath"
	"regexp"
	"slices"
)

type FileItem struct {
	Name        string
	Path        string
	IsDirectory bool
	Children    []FileItem
}

const (
	sessionDirPath   = ".theorem-note"
	sessionFileName  = "session.json"
	theoremsFileName = "theorems.json"
)

func GetNewDirectoryFileTree(ctx context.Context) (string, []FileItem, error) {
	path, err := runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{
		Title: "Select Directory",
	})
	if err != nil {
		return "", nil, err
	}

	items, err := GetFileTree(path)
	if err != nil {
		return "", nil, err
	}
	return path, items, nil
}

func GetFileTree(path string) ([]FileItem, error) {
	fileInfo, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var items []FileItem
	for _, fi := range fileInfo {
		item := FileItem{
			Name:        fi.Name(),
			Path:        filepath.Join(path, fi.Name()),
			IsDirectory: fi.IsDir(),
		}
		if fi.IsDir() {
			children, err := GetFileTree(item.Path)
			if err != nil {
				return nil, err
			}
			item.Children = children
		}
		items = append(items, item)
	}

	slices.SortFunc(items, func(a, b FileItem) int {
		if a.IsDirectory && !b.IsDirectory {
			return -1
		} else if !a.IsDirectory && b.IsDirectory {
			return 1
		}
		return cmp.Compare(a.Name, b.Name)
	})

	return items, nil
}

func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func WriteFile(path string, content string, rootDir string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}
	return extractAndSaveTheorems(path, content, rootDir)
}

func CreateFile(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return os.ErrExist
	}
	return os.WriteFile(path, []byte(""), 0644)
}

func CreateDirectory(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return os.ErrExist
	}
	return os.Mkdir(path, 0755)
}

func getTheoremsFilePath(rootDir string) (string, error) {
	if rootDir == "" {
		return "", os.ErrInvalid
	}
	return filepath.Join(rootDir, sessionDirPath, theoremsFileName), nil
}

func extractAndSaveTheorems(path string, content string, rootDir string) error {
	re := regexp.MustCompile(`<theorem name="([^"]+)">`)
	matches := re.FindAllStringSubmatch(content, -1)

	if len(matches) == 0 {
		return nil
	}

	if err := ensureSessionDirExists(rootDir); err != nil {
		return err
	}

	theoremsFilePath, err := getTheoremsFilePath(rootDir)
	if err != nil {
		return err
	}

	theorems := make(map[string]string)
	if _, err := os.Stat(theoremsFilePath); !os.IsNotExist(err) {
		data, err := os.ReadFile(theoremsFilePath)
		if err != nil {
			return err
		}
		if len(data) > 0 {
			if err := json.Unmarshal(data, &theorems); err != nil {
				return err
			}
		}
	}

	for k, v := range theorems {
		if v == path {
			delete(theorems, k)
		}
	}

	for _, match := range matches {
		if len(match) > 1 {
			theoremName := match[1]
			theorems[theoremName] = path
		}
	}

	data, err := json.MarshalIndent(theorems, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("Saving theorems to %s: %v\n", theoremsFilePath, theorems)

	return os.WriteFile(theoremsFilePath, data, 0644)
}

func getSessionFilePath(rootDir string) (string, error) {
	if rootDir == "" {
		return "", os.ErrInvalid
	}
	return filepath.Join(rootDir, sessionDirPath, sessionFileName), nil
}

func ensureSessionDirExists(rootDir string) error {
	if rootDir == "" {
		return os.ErrInvalid
	}
	return os.MkdirAll(filepath.Join(rootDir, sessionDirPath), 0755)
}

func SaveSession(rootDir string, filePaths []string) error {
	if err := ensureSessionDirExists(rootDir); err != nil {
		return err
	}

	path, err := getSessionFilePath(rootDir)
	if err != nil {
		return err
	}

	data, err := json.Marshal(filePaths)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func LoadSession(rootDir string) ([]string, error) {
	if rootDir == "" {
		return []string{}, nil
	}

	path, err := getSessionFilePath(rootDir)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []string{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var filePaths []string
	if err := json.Unmarshal(data, &filePaths); err != nil {
		return nil, err
	}

	return filePaths, nil
}

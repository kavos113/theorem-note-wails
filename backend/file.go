package backend

import (
	"cmp"
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"path/filepath"
	"slices"
)

type FileItem struct {
	Name        string
	Path        string
	IsDirectory bool
	Children    []FileItem
}

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

func WriteFile(path string, content string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

package backend

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetFileTree(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a nested directory structure
	err = os.MkdirAll(filepath.Join(tmpDir, "dir1", "dir2"), 0755)
	if err != nil {
		t.Fatalf("Failed to create nested dir: %v", err)
	}

	// Create some test files
	err = os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("file1"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file1.txt: %v", err)
	}
	err = os.WriteFile(filepath.Join(tmpDir, "dir1", "file2.txt"), []byte("file2"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file2.txt: %v", err)
	}
	err = os.WriteFile(filepath.Join(tmpDir, "dir1", "dir2", "file3.txt"), []byte("file3"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file3.txt: %v", err)
	}

	// Call the function to be tested
	items, err := GetFileTree(tmpDir)
	if err != nil {
		t.Fatalf("GetFileTree failed: %v", err)
	}

	// Define the expected structure
	expected := []FileItem{
		{
			Name:        "dir1",
			Path:        filepath.Join(tmpDir, "dir1"),
			IsDirectory: true,
			Children: []FileItem{
				{
					Name:        "dir2",
					Path:        filepath.Join(tmpDir, "dir1", "dir2"),
					IsDirectory: true,
					Children: []FileItem{
						{
							Name:        "file3.txt",
							Path:        filepath.Join(tmpDir, "dir1", "dir2", "file3.txt"),
							IsDirectory: false,
							Children:    nil,
						},
					},
				},
				{
					Name:        "file2.txt",
					Path:        filepath.Join(tmpDir, "dir1", "file2.txt"),
					IsDirectory: false,
					Children:    nil,
				},
			},
		},
		{
			Name:        "file1.txt",
			Path:        filepath.Join(tmpDir, "file1.txt"),
			IsDirectory: false,
			Children:    nil,
		},
	}

	// Compare the actual and expected results
	if !reflect.DeepEqual(items, expected) {
		t.Errorf("GetFileTree returned unexpected structure.\nGot:  %v\nWant: %v", items, expected)
	}
}

func TestSaveAndLoadSession(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Define the session data to be saved
	expectedFilePaths := []string{"/path/to/file1", "/path/to/file2"}

	// Test saving the session
	err = SaveSession(tmpDir, expectedFilePaths)
	if err != nil {
		t.Fatalf("SaveSession failed: %v", err)
	}

	// Verify that the session file was created
	sessionFilePath, _ := getSessionFilePath(tmpDir)
	if _, err := os.Stat(sessionFilePath); os.IsNotExist(err) {
		t.Fatalf("Session file was not created")
	}

	// Test loading the session
	loadedFilePaths, err := LoadSession(tmpDir)
	if err != nil {
		t.Fatalf("LoadSession failed: %v", err)
	}

	// Compare the loaded data with the original data
	if !reflect.DeepEqual(loadedFilePaths, expectedFilePaths) {
		t.Errorf("Loaded session data does not match expected data.\nGot:  %v\nWant: %v", loadedFilePaths, expectedFilePaths)
	}

	// Test loading from a non-existent session
	err = os.RemoveAll(filepath.Join(tmpDir, sessionDirPath))
	if err != nil {
		t.Fatalf("Failed to remove session dir: %v", err)
	}
	loadedFilePaths, err = LoadSession(tmpDir)
	if err != nil {
		t.Fatalf("LoadSession failed for non-existent session: %v", err)
	}
	if len(loadedFilePaths) != 0 {
		t.Errorf("Expected empty slice for non-existent session, but got %v", loadedFilePaths)
	}
}

func TestGetSessionFilePath(t *testing.T) {
	// Test with a valid root directory
	rootDir := "/some/dir"
	expectedPath := filepath.Join(rootDir, sessionDirPath, sessionFileName)
	actualPath, err := getSessionFilePath(rootDir)
	if err != nil {
		t.Errorf("getSessionFilePath failed with valid rootDir: %v", err)
	}
	if actualPath != expectedPath {
		t.Errorf("Expected path %s, but got %s", expectedPath, actualPath)
	}

	// Test with an empty root directory
	_, err = getSessionFilePath("")
	if err != os.ErrInvalid {
		t.Errorf("Expected os.ErrInvalid for empty rootDir, but got %v", err)
	}
}

func TestEnsureSessionDirExists(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test creating the session directory
	err = ensureSessionDirExists(tmpDir)
	if err != nil {
		t.Fatalf("ensureSessionDirExists failed: %v", err)
	}

	// Verify that the session directory was created
	sessionDirPath := filepath.Join(tmpDir, ".theorem-note")
	if _, err := os.Stat(sessionDirPath); os.IsNotExist(err) {
		t.Fatalf("Session directory was not created")
	}

	// Test with an empty root directory
	err = ensureSessionDirExists("")
	if err != os.ErrInvalid {
		t.Errorf("Expected os.ErrInvalid for empty rootDir, but got %v", err)
	}
}

func TestReadFile(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some content to the file
	expectedContent := "hello world"
	_, err = tmpFile.WriteString(expectedContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test reading the file
	actualContent, err := ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	// Compare the read content with the expected content
	if actualContent != expectedContent {
		t.Errorf("ReadFile returned unexpected content.\nGot:  %s\nWant: %s", actualContent, expectedContent)
	}
}

func TestWriteFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp(tmpDir, "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Define the content to be written
	expectedContent := "hello world"

	// Test writing to the file
	err = WriteFile(tmpFile.Name(), expectedContent, "")
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Read the file to verify the content
	actualContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read back file: %v", err)
	}

	// Compare the written content with the expected content
	if string(actualContent) != expectedContent {
		t.Errorf("WriteFile wrote unexpected content.\nGot:  %s\nWant: %s", string(actualContent), expectedContent)
	}
}

func TestWriteFile_WithTheoremTag(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a dummy file path
	filePath := filepath.Join(tmpDir, "test.md")

	// Define content with a <theorem> tag
	content := "This is a test file with a theorem.\n<theorem name=\"Test Theorem\">Some theorem content.</theorem>"

	// Call WriteFile
	err = WriteFile(filePath, content, "")
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Check if .theorem-note/theorems.json is created
	theoremsFilePath := filepath.Join(tmpDir, ".theorem-note", "theorems.json")
	if _, err := os.Stat(theoremsFilePath); os.IsNotExist(err) {
		t.Fatalf("theorems.json was not created")
	}

	// Check the content of theorems.json
	file, err := os.ReadFile(theoremsFilePath)
	if err != nil {
		t.Fatalf("Failed to read theorems.json: %v", err)
	}

	var theorems map[string]string
	err = json.Unmarshal(file, &theorems)
	if err != nil {
		t.Fatalf("Failed to unmarshal theorems.json: %v", err)
	}

	expectedTheorems := map[string]string{
		"Test Theorem": filePath,
	}

	if !reflect.DeepEqual(theorems, expectedTheorems) {
		t.Errorf("Theorems map is incorrect.\nGot:  %v\nWant: %v", theorems, expectedTheorems)
	}
}

func TestJsonMarshalFileItem(t *testing.T) {
	// Create a sample FileItem structure
	fileItem := FileItem{
		Name:        "dir1",
		Path:        "/path/to/dir1",
		IsDirectory: true,
		Children: []FileItem{
			{
				Name:        "file1.txt",
				Path:        "/path/to/dir1/file1.txt",
				IsDirectory: false,
				Children:    nil,
			},
		},
	}

	// Marshal the FileItem to JSON
	jsonData, err := json.Marshal(fileItem)
	if err != nil {
		t.Fatalf("Failed to marshal FileItem to JSON: %v", err)
	}

	// Unmarshal the JSON back to a FileItem
	var newFileItem FileItem
	err = json.Unmarshal(jsonData, &newFileItem)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to FileItem: %v", err)
	}

	// Compare the original and unmarshaled FileItems
	if !reflect.DeepEqual(fileItem, newFileItem) {
		t.Errorf("Unmarshaled FileItem does not match original.\nGot:  %v\nWant: %v", newFileItem, fileItem)
	}
}

func TestCreateFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "testfile.txt")

	// Test creating a new file
	err = CreateFile(filePath)
	if err != nil {
		t.Fatalf("CreateFile failed: %v", err)
	}

	// Verify that the file was created
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("File was not created")
	}

	// Test creating a file that already exists
	err = CreateFile(filePath)
	if err == nil {
		t.Fatalf("Expected an error when creating a file that already exists, but got nil")
	}
}

func TestCreateDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dirPath := filepath.Join(tmpDir, "testdir")

	// Test creating a new directory
	err = CreateDirectory(dirPath)
	if err != nil {
		t.Fatalf("CreateDirectory failed: %v", err)
	}

	// Verify that the directory was created
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Fatalf("Directory was not created")
	}

	// Test creating a directory that already exists
	err = CreateDirectory(dirPath)
	if err == nil {
		t.Fatalf("Expected an error when creating a directory that already exists, but got nil")
	}
}

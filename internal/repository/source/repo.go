package source

import (
	"fmt"
	"os"
	"path/filepath"
)

type SourceRepo struct{}

func New() *SourceRepo {
	return &SourceRepo{}
}

func (r SourceRepo) DirsExist(paths []string) error {
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("path %s does not exist or cannot be accessed: %w", path, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("path %s is not a directory", path)
		}
	}
	return nil
}

func (r SourceRepo) WriteFile(path string, source string) (string, error) {
	dirPath := filepath.Dir(path)

	err := os.MkdirAll(dirPath, 0o755)
	if err != nil {
		return "", fmt.Errorf("error creating directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(source)
	if err != nil {
		return "", fmt.Errorf("error writing to file: %w", err)
	}

	return path, nil
}

func (r SourceRepo) FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

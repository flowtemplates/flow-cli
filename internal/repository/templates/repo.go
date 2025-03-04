package templates

import (
	"fmt"
	"os"
	"path/filepath"
)

type TemplatesRepo struct {
	baseDir string
}

func New(baseDir string) *TemplatesRepo {
	return &TemplatesRepo{
		baseDir: baseDir,
	}
}

func (r TemplatesRepo) GetTemplatesNames() ([]string, error) {
	var directories []string
	files, err := os.ReadDir(r.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, file.Name())
		}
	}
	return directories, nil
}

type File struct {
	Name string
	Path string
}

// Dir represents a directory containing files and subdirectories.
type Dir struct {
	Name  string
	Path  string
	Files []File
	Dirs  []Dir
}

func readDirTree(dirPath string) (Dir, error) {
	root := Dir{Name: filepath.Base(dirPath), Path: dirPath}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return Dir{}, err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			subDir, err := readDirTree(fullPath)
			if err != nil {
				return Dir{}, err
			}
			root.Dirs = append(root.Dirs, subDir)
		} else {
			root.Files = append(root.Files, File{Name: entry.Name(), Path: fullPath})
		}
	}

	return root, nil
}

func (r TemplatesRepo) GetTemplate(templateName string) (Dir, error) {
	templatePath := filepath.Join(r.baseDir, templateName)
	return readDirTree(templatePath)
}

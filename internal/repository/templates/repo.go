package templates

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flowtemplates/cli/pkg/fs"
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

func (r TemplatesRepo) readDirTree(templateName string, relBaseDirPath string) (fs.Dir, error) {
	fullDirPath := filepath.Join(r.baseDir, templateName, relBaseDirPath)

	root := fs.Dir{Name: filepath.Base(relBaseDirPath), Path: filepath.Dir(relBaseDirPath)}

	entries, err := os.ReadDir(fullDirPath)
	if err != nil {
		return fs.Dir{}, err
	}

	for _, entry := range entries {
		relPath := filepath.Join(relBaseDirPath, entry.Name())

		if entry.IsDir() {
			subDir, err := r.readDirTree(templateName, relPath)
			if err != nil {
				return fs.Dir{}, err
			}
			root.Dirs = append(root.Dirs, subDir)
		} else {
			full := filepath.Join(r.baseDir, templateName, relPath)
			source, err := os.ReadFile(full)
			if err != nil {
				return fs.Dir{}, fmt.Errorf("failed to open file %s: %w", relPath, err)
			}

			root.Files = append(root.Files, fs.File{
				Name:   entry.Name(),
				Path:   relPath,
				Source: string(source),
			})
		}
	}

	return root, nil
}

func (r TemplatesRepo) GetTemplate(templateName string) (fs.Dir, error) {
	return r.readDirTree(templateName, "")
}

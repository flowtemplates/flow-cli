package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/flowtemplates/cli/internal/repository/templates"
)

type templatesRepo interface {
	GetTemplatesNames() ([]string, error)
	GetTemplate(templateName string) (templates.Dir, error)
}

type Service struct {
	tr templatesRepo
}

func New(tr templatesRepo) *Service {
	return &Service{
		tr: tr,
	}
}

func (s Service) ListTemplates() ([]string, error) {
	templateNames, err := s.tr.GetTemplatesNames()
	if err != nil {
		return nil, fmt.Errorf("failed to get templates names: %w", err)
	}

	return templateNames, nil
}

func checkDirectories(paths []string) error {
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

func (s Service) Add(templateName string, dests ...string) error {
	if len(dests) < 1 {
		return errors.New("at least one dest required")
	}

	if err := checkDirectories(dests); err != nil {
		return err
	}

	templateDir, err := s.tr.GetTemplate(templateName)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	j, _ := json.MarshalIndent(templateDir, "", "  ")
	fmt.Printf("%s\n", j)
	return nil
}

func (s Service) Get(templateName string) error {
	templateDir, err := s.tr.GetTemplate(templateName)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	j, _ := json.MarshalIndent(templateDir, "", "  ")
	fmt.Printf("%s\n", j)

	return nil
}

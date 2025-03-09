package service

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/flowtemplates/cli/pkg/flow-go/analyzer"
	"github.com/flowtemplates/cli/pkg/flow-go/renderer"
	"github.com/flowtemplates/cli/pkg/fs"
)

type templatesRepo interface {
	GetTemplatesNames() ([]string, error)
	GetTemplate(templateName string) (fs.Dir, error)
}

type sourceRepo interface {
	DirsExist(paths []string) error
	WriteFile(path string, source string) (string, error)
	FileExists(path string) bool
}

type Service struct {
	tr templatesRepo
	sr sourceRepo
}

func New(tr templatesRepo, sr sourceRepo) *Service {
	return &Service{
		tr: tr,
		sr: sr,
	}
}

func (s Service) ListTemplates() ([]string, error) {
	templateNames, err := s.tr.GetTemplatesNames()
	if err != nil {
		return nil, fmt.Errorf("failed to get templates names: %w", err)
	}

	return templateNames, nil
}

func (s Service) Add(
	templateName string,
	scope map[string]*string,
	overwriteFn func(files []string) ([]string, error),
	dests ...string,
) error {
	if len(dests) < 1 {
		return errors.New("at least one dest required")
	}

	if err := s.sr.DirsExist(dests); err != nil {
		return err // nolint: wrapcheck
	}

	templateDir, err := s.tr.GetTemplate(templateName)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	tm := make(analyzer.TypeMap)
	if err := getTypeMapFromDir(templateDir, tm); err != nil {
		return err
	}

	sc := renderer.Scope{}
	for n, v := range scope {
		if v == nil {
			sc[n] = "true"
		} else {
			sc[n] = *v
		}
	}

	if err := analyzer.Typecheck(sc, tm); err != nil {
		return fmt.Errorf("TypeErrors: %s", err)
	}

	rendered, err := s.renderDir(templateDir, sc)
	if err != nil {
		return err
	}

	filesToWrite := make(map[string]string)
	overwriteRequest := []string{}

	for _, dest := range dests {
		for path, source := range rendered {
			destPath := filepath.Join(dest, path)
			if s.sr.FileExists(destPath) {
				overwriteRequest = append(overwriteRequest, destPath)
			}
			filesToWrite[destPath] = source
		}
	}

	if len(overwriteRequest) > 0 {
		overwrite, err := overwriteFn(overwriteRequest)
		if err != nil {
			return err
		}

		for _, initOverwrite := range overwriteRequest {
			if !slices.Contains(overwrite, initOverwrite) {
				delete(filesToWrite, initOverwrite)
			}
		}
	}

	for path, source := range filesToWrite {
		_, err := s.sr.WriteFile(path, source)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	return nil
}

func (s Service) GetTemplate(templateName string) (analyzer.TypeMap, error) {
	templateDir, err := s.tr.GetTemplate(templateName)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	tm := make(analyzer.TypeMap)
	if err := getTypeMapFromDir(templateDir, tm); err != nil {
		return nil, err
	}

	return tm, nil
}

const templateFileExt = ".ft"

func isTemplateFile(file fs.File) bool {
	return strings.HasSuffix(file.Name, templateFileExt)
}

func (s Service) renderDir(dir fs.Dir, scope renderer.Scope) (map[string]string, error) {
	f := make(map[string]string)
	if err := s.renderDirRecursive(dir, scope, f); err != nil {
		return nil, err
	}

	return f, nil
}

func (s Service) renderDirRecursive(dir fs.Dir, scope renderer.Scope, out map[string]string) error {
	for _, d := range dir.Dirs {
		dirName, err := renderer.RenderString(d.Name, scope)
		if err != nil {
			return fmt.Errorf("failed to render dirName: %w", err)
		}
		d.Name = dirName

		if err := s.renderDirRecursive(d, scope, out); err != nil {
			return err
		}
	}

	for _, file := range dir.Files {
		filename, err := renderer.RenderString(file.Name, scope)
		if err != nil {
			return fmt.Errorf("failed to render filename: %w", err)
		}

		content := file.Source
		if isTemplateFile(file) {
			content, err = renderer.RenderString(file.Source, scope)
			if err != nil {
				return fmt.Errorf("failed to render file %s: %w", filename, err)
			}
			filename = strings.TrimSuffix(filename, templateFileExt)
		}

		out[filepath.Join(dir.Path, dir.Name, filename)] = content
	}

	return nil
}

func getTypeMapFromDir(dir fs.Dir, tm analyzer.TypeMap) error {
	for _, file := range dir.Files {
		if err := analyzer.GetTypeMapFromString(file.Name, tm); err != nil {
			return fmt.Errorf("failed to parse types in filename: %w", err)
		}

		if isTemplateFile(file) {
			if err := analyzer.GetTypeMapFromString(file.Source, tm); err != nil {
				return fmt.Errorf("failed to parse types in file: %w", err)
			}
		}
	}

	for _, d := range dir.Dirs {
		if err := getTypeMapFromDir(d, tm); err != nil {
			return err
		}
	}

	return nil
}

package service

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/flowtemplates/cli/pkg/flow-go/analyzer"
	"github.com/flowtemplates/cli/pkg/flow-go/lexer"
	"github.com/flowtemplates/cli/pkg/flow-go/parser"
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
	overwriteFn func(path string) bool,
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

	sc := renderer.Scope{} // TODO: rename
	for n, v := range scope {
		if v == nil {
			sc[n] = "true"
		} else {
			sc[n] = *v
		}
	}

	if err := analyzer.Typecheck(sc, tm); err != nil {
		return fmt.Errorf("TypeErrors: %w", err)
	}

	f := []F{}
	if err := s.renderDir(templateDir, sc, &f); err != nil {
		return err
	}
	fmt.Println(f)

	list := []F{}

	for _, dest := range dests {
		for _, file := range f {
			destPath := filepath.Join(dest, file.Path)
			fmt.Printf("dest %s\n", destPath)
			if s.sr.FileExists(destPath) {
				if !overwriteFn(destPath) {
					continue
				}
			}
			list = append(list, F{
				Path:   destPath,
				Source: file.Source,
			})
		}
	}

	for _, file := range list {
		_, err := s.sr.WriteFile(file.Path, file.Source)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	return nil
}

func (s Service) Get(templateName string) (analyzer.TypeMap, error) {
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

type F struct {
	Path   string
	Source string
}

func (s Service) renderDir(dir fs.Dir, scope renderer.Scope, out *[]F) error {
	for _, d := range dir.Dirs {
		dirName, err := render(d.Name, scope)
		if err != nil {
			return err
		}
		d.Name = dirName

		if err := s.renderDir(d, scope, out); err != nil {
			return err
		}
	}

	for _, file := range dir.Files {
		filename, err := render(file.Name, scope)
		if err != nil {
			return err
		}

		content := file.Source
		if isTemplateFile(file) {
			content, err = render(file.Source, scope)
			if err != nil {
				return err
			}
		}

		*out = append(*out, F{
			Path:   filepath.Join(dir.Path, dir.Name, filename),
			Source: content,
		})
	}

	return nil
}

func getTypeMapFromDir(dir fs.Dir, tm analyzer.TypeMap) error {
	for _, file := range dir.Files {
		if err := getTypeMap(file.Name, tm); err != nil {
			return err
		}

		if isTemplateFile(file) {
			if err := getTypeMap(file.Source, tm); err != nil {
				return err
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

func render(input string, tm renderer.Scope) (string, error) {
	tokens := lexer.TokensFromString(input)
	ast, errs := parser.New(tokens).Parse()
	if len(errs) != 0 {
		return "", errs[0]
	}

	res, err := renderer.RenderAst(ast, tm)
	if err != nil {
		return "", fmt.Errorf("failed to render: %w", err)
	}

	return res, nil
}

func getTypeMap(input string, tm analyzer.TypeMap) error {
	tokens := lexer.TokensFromString(input)
	ast, errs := parser.New(tokens).Parse()
	if len(errs) != 0 {
		return errs[0]
	}

	if errs := analyzer.GetTypeMap(ast, tm); len(errs) != 0 {
		return errs[0]
	}

	return nil
}

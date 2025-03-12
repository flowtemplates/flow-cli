package cli

import (
	"log/slog"

	"github.com/flowtemplates/flow-go/analyzer"
)

type iService interface {
	ListTemplates() ([]string, error)
	Create(
		templateName string,
		scope map[string]*string,
		overwriteFn func(paths []string) ([]string, error),
		dests ...string,
	) error
	GetTemplateContext(templateName string) (analyzer.TypeMap, error)
}

type CliController struct {
	service iService
	logger  *slog.Logger
}

func New(service iService, logger *slog.Logger) *CliController {
	return &CliController{
		service: service,
		logger:  logger,
	}
}
